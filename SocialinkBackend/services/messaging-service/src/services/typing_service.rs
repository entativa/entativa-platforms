use anyhow::Result;
use redis::AsyncCommands;
use uuid::Uuid;
use chrono::{Utc, DateTime};
use serde::{Serialize, Deserialize};

/// Typing Indicator Service
/// Ephemeral typing indicators (stored in Redis only)
pub struct TypingService {
    redis: redis::Client,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct TypingIndicator {
    pub conversation_id: Uuid,
    pub user_id: Uuid,
    pub is_typing: bool,
    pub timestamp: DateTime<Utc>,
}

impl TypingService {
    pub fn new(redis: redis::Client) -> Self {
        Self { redis }
    }
    
    /// Set user typing in conversation
    pub async fn set_typing(&self, conversation_id: Uuid, user_id: Uuid) -> Result<()> {
        let mut conn = self.redis.get_async_connection().await?;
        
        // Store in Redis with 10-second TTL (typing expires quickly)
        let key = format!("typing:{}:{}", conversation_id, user_id);
        let data = serde_json::json!({
            "conversation_id": conversation_id,
            "user_id": user_id,
            "is_typing": true,
            "timestamp": Utc::now(),
        });
        
        let _: () = conn.set_ex(&key, data.to_string(), 10).await?;
        
        // Publish typing indicator
        self.publish_typing_indicator(conversation_id, user_id, true).await?;
        
        Ok(())
    }
    
    /// Clear user typing (stopped typing)
    pub async fn clear_typing(&self, conversation_id: Uuid, user_id: Uuid) -> Result<()> {
        let mut conn = self.redis.get_async_connection().await?;
        
        // Remove from Redis
        let key = format!("typing:{}:{}", conversation_id, user_id);
        let _: () = conn.del(&key).await?;
        
        // Publish typing stopped
        self.publish_typing_indicator(conversation_id, user_id, false).await?;
        
        Ok(())
    }
    
    /// Get who's typing in conversation
    pub async fn get_typing_users(&self, conversation_id: Uuid) -> Result<Vec<Uuid>> {
        let mut conn = self.redis.get_async_connection().await?;
        
        // Get all keys for this conversation
        let pattern = format!("typing:{}:*", conversation_id);
        let keys: Vec<String> = conn.keys(&pattern).await?;
        
        let mut typing_users = Vec::new();
        
        for key in keys {
            if let Ok(data) = conn.get::<_, String>(&key).await {
                if let Ok(json) = serde_json::from_str::<serde_json::Value>(&data) {
                    if let Some(user_id_str) = json.get("user_id").and_then(|v| v.as_str()) {
                        if let Ok(user_id) = Uuid::parse_str(user_id_str) {
                            typing_users.push(user_id);
                        }
                    }
                }
            }
        }
        
        Ok(typing_users)
    }
    
    /// Check if user is typing
    pub async fn is_typing(&self, conversation_id: Uuid, user_id: Uuid) -> Result<bool> {
        let mut conn = self.redis.get_async_connection().await?;
        
        let key = format!("typing:{}:{}", conversation_id, user_id);
        let exists: bool = conn.exists(&key).await?;
        
        Ok(exists)
    }
    
    /// Publish typing indicator
    async fn publish_typing_indicator(
        &self,
        conversation_id: Uuid,
        user_id: Uuid,
        is_typing: bool,
    ) -> Result<()> {
        let mut conn = self.redis.get_async_connection().await?;
        
        let indicator = TypingIndicator {
            conversation_id,
            user_id,
            is_typing,
            timestamp: Utc::now(),
        };
        
        let channel = format!("typing:{}", conversation_id);
        let payload = serde_json::to_string(&indicator)?;
        
        let _: () = conn.publish(&channel, payload).await?;
        
        Ok(())
    }
    
    /// Clean up expired typing indicators (automatic via TTL)
    pub async fn cleanup_expired(&self) -> Result<usize> {
        // Redis TTL handles this automatically
        // This method is for manual cleanup if needed
        Ok(0)
    }
}
