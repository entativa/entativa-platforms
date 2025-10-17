use anyhow::Result;
use sqlx::PgPool;
use redis::AsyncCommands;
use uuid::Uuid;
use chrono::{Utc, DateTime};
use serde::{Serialize, Deserialize};

use crate::models::message::{Presence, PresenceStatus};

/// Presence Service
/// Tracks user online/offline status and last seen
pub struct PresenceService {
    db: PgPool,
    redis: redis::Client,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct PresenceUpdate {
    pub user_id: Uuid,
    pub status: PresenceStatus,
    pub custom_status: Option<String>,
    pub timestamp: DateTime<Utc>,
}

impl PresenceService {
    pub fn new(db: PgPool, redis: redis::Client) -> Self {
        Self { db, redis }
    }
    
    /// Set user online
    pub async fn set_online(&self, user_id: Uuid, device_id: String) -> Result<()> {
        let mut conn = self.redis.get_async_connection().await?;
        
        // Set in Redis (fast lookup)
        let key = format!("presence:{}", user_id);
        let data = serde_json::json!({
            "status": "Online",
            "last_seen": Utc::now(),
            "device_id": device_id,
        });
        
        let _: () = conn.set_ex(&key, data.to_string(), 300).await?; // 5-minute TTL
        
        // Update database
        sqlx::query!(
            r#"
            INSERT INTO user_presence (user_id, status, last_seen, updated_at)
            VALUES ($1, $2, NOW(), NOW())
            ON CONFLICT (user_id) DO UPDATE SET
                status = $2,
                last_seen = NOW(),
                updated_at = NOW()
            "#,
            user_id,
            "Online"
        )
        .execute(&self.db)
        .await?;
        
        // Publish presence update
        self.publish_presence_update(user_id, PresenceStatus::Online, None).await?;
        
        Ok(())
    }
    
    /// Set user offline
    pub async fn set_offline(&self, user_id: Uuid) -> Result<()> {
        let mut conn = self.redis.get_async_connection().await?;
        
        // Remove from Redis
        let key = format!("presence:{}", user_id);
        let _: () = conn.del(&key).await?;
        
        // Update database
        sqlx::query!(
            r#"
            UPDATE user_presence
            SET status = $1, last_seen = NOW(), updated_at = NOW()
            WHERE user_id = $2
            "#,
            "Offline",
            user_id
        )
        .execute(&self.db)
        .await?;
        
        // Publish presence update
        self.publish_presence_update(user_id, PresenceStatus::Offline, None).await?;
        
        Ok(())
    }
    
    /// Set user away
    pub async fn set_away(&self, user_id: Uuid) -> Result<()> {
        let mut conn = self.redis.get_async_connection().await?;
        
        // Update Redis
        let key = format!("presence:{}", user_id);
        let data = serde_json::json!({
            "status": "Away",
            "last_seen": Utc::now(),
        });
        
        let _: () = conn.set_ex(&key, data.to_string(), 300).await?;
        
        // Update database
        sqlx::query!(
            "UPDATE user_presence SET status = $1, updated_at = NOW() WHERE user_id = $2",
            "Away",
            user_id
        )
        .execute(&self.db)
        .await?;
        
        // Publish presence update
        self.publish_presence_update(user_id, PresenceStatus::Away, None).await?;
        
        Ok(())
    }
    
    /// Set custom status
    pub async fn set_custom_status(&self, user_id: Uuid, custom_status: String) -> Result<()> {
        // Update database
        sqlx::query!(
            "UPDATE user_presence SET custom_status = $1, updated_at = NOW() WHERE user_id = $2",
            custom_status,
            user_id
        )
        .execute(&self.db)
        .await?;
        
        // Update Redis
        let mut conn = self.redis.get_async_connection().await?;
        let key = format!("presence:{}", user_id);
        
        if let Ok(current) = conn.get::<_, String>(&key).await {
            if let Ok(mut data) = serde_json::from_str::<serde_json::Value>(&current) {
                data["custom_status"] = serde_json::json!(custom_status);
                let _: () = conn.set_ex(&key, data.to_string(), 300).await?;
            }
        }
        
        Ok(())
    }
    
    /// Clear custom status
    pub async fn clear_custom_status(&self, user_id: Uuid) -> Result<()> {
        sqlx::query!(
            "UPDATE user_presence SET custom_status = NULL, updated_at = NOW() WHERE user_id = $2",
            user_id
        )
        .execute(&self.db)
        .await?;
        
        Ok(())
    }
    
    /// Get user presence
    pub async fn get_presence(&self, user_id: Uuid) -> Result<Presence> {
        // Try Redis first (fast)
        let mut conn = self.redis.get_async_connection().await?;
        let key = format!("presence:{}", user_id);
        
        if let Ok(data) = conn.get::<_, String>(&key).await {
            if let Ok(json) = serde_json::from_str::<serde_json::Value>(&data) {
                return Ok(Presence {
                    user_id,
                    status: PresenceStatus::Online,
                    last_seen: json.get("last_seen")
                        .and_then(|v| v.as_str())
                        .and_then(|s| DateTime::parse_from_rfc3339(s).ok())
                        .map(|dt| dt.with_timezone(&Utc))
                        .unwrap_or_else(Utc::now),
                    custom_status: json.get("custom_status")
                        .and_then(|v| v.as_str())
                        .map(|s| s.to_string()),
                });
            }
        }
        
        // Fallback to database
        let row = sqlx::query!(
            "SELECT status, last_seen, custom_status FROM user_presence WHERE user_id = $1",
            user_id
        )
        .fetch_optional(&self.db)
        .await?;
        
        if let Some(row) = row {
            Ok(Presence {
                user_id,
                status: match row.status.as_str() {
                    "Online" => PresenceStatus::Online,
                    "Away" => PresenceStatus::Away,
                    "Busy" => PresenceStatus::Busy,
                    _ => PresenceStatus::Offline,
                },
                last_seen: row.last_seen,
                custom_status: row.custom_status,
            })
        } else {
            // Default to offline
            Ok(Presence {
                user_id,
                status: PresenceStatus::Offline,
                last_seen: Utc::now(),
                custom_status: None,
            })
        }
    }
    
    /// Get presence for multiple users (bulk)
    pub async fn get_bulk_presence(&self, user_ids: Vec<Uuid>) -> Result<Vec<Presence>> {
        let mut results = Vec::new();
        
        for user_id in user_ids {
            if let Ok(presence) = self.get_presence(user_id).await {
                results.push(presence);
            }
        }
        
        Ok(results)
    }
    
    /// Check if user is online
    pub async fn is_online(&self, user_id: Uuid) -> Result<bool> {
        let mut conn = self.redis.get_async_connection().await?;
        let key = format!("presence:{}", user_id);
        
        let exists: bool = conn.exists(&key).await?;
        Ok(exists)
    }
    
    /// Heartbeat - keep user online
    pub async fn heartbeat(&self, user_id: Uuid, device_id: String) -> Result<()> {
        let mut conn = self.redis.get_async_connection().await?;
        let key = format!("presence:{}", user_id);
        
        // Extend TTL
        let _: () = conn.expire(&key, 300).await?;
        
        // Update last seen
        sqlx::query!(
            "UPDATE user_presence SET last_seen = NOW() WHERE user_id = $1",
            user_id
        )
        .execute(&self.db)
        .await?;
        
        Ok(())
    }
    
    /// Publish presence update
    async fn publish_presence_update(
        &self,
        user_id: Uuid,
        status: PresenceStatus,
        custom_status: Option<String>,
    ) -> Result<()> {
        let mut conn = self.redis.get_async_connection().await?;
        
        let update = PresenceUpdate {
            user_id,
            status,
            custom_status,
            timestamp: Utc::now(),
        };
        
        let channel = format!("presence:updates");
        let payload = serde_json::to_string(&update)?;
        
        let _: () = conn.publish(&channel, payload).await?;
        
        Ok(())
    }
    
    /// Get online count (for stats)
    pub async fn get_online_count(&self) -> Result<usize> {
        let mut conn = self.redis.get_async_connection().await?;
        
        let keys: Vec<String> = conn.keys("presence:*").await?;
        Ok(keys.len())
    }
}
