use anyhow::{anyhow, Result};
use sqlx::PgPool;
use uuid::Uuid;
use chrono::Utc;
use std::collections::HashMap;

use crate::models::keys::*;
use crate::crypto::signal::SignalProtocol;

/// Key Management Service
/// Handles device registration, pre-key distribution, and key rotation
pub struct KeyService {
    db: PgPool,
}

impl KeyService {
    pub fn new(db: PgPool) -> Self {
        Self { db }
    }
    
    /// Register device with keys
    pub async fn register_device(
        &self,
        user_id: Uuid,
        request: KeyRegistrationRequest,
    ) -> Result<Device> {
        // Validate request
        request.validate()?;
        
        // Decode identity key
        let identity_key = base64::decode(&request.identity_key)?;
        
        // Check if device already exists
        let existing = sqlx::query!(
            "SELECT id FROM devices WHERE user_id = $1 AND device_id = $2",
            user_id,
            request.device_id
        )
        .fetch_optional(&self.db)
        .await?;
        
        if existing.is_some() {
            return Err(anyhow!("Device already registered"));
        }
        
        // Start transaction
        let mut tx = self.db.begin().await?;
        
        // Insert device
        let device_id_uuid = Uuid::new_v4();
        sqlx::query!(
            r#"
            INSERT INTO devices (id, user_id, device_id, device_name, registration_id, identity_key, is_active, last_seen, created_at)
            VALUES ($1, $2, $3, $4, $5, $6, true, NOW(), NOW())
            "#,
            device_id_uuid,
            user_id,
            request.device_id,
            request.device_name,
            request.registration_id,
            &identity_key,
        )
        .execute(&mut *tx)
        .await?;
        
        // Insert signed pre-key
        let signed_prekey_public = base64::decode(&request.signed_prekey.public_key)?;
        let signed_prekey_signature = base64::decode(&request.signed_prekey.signature)?;
        
        // Verify signature
        if !SignalProtocol::verify_prekey_signature(
            &signed_prekey_public,
            &signed_prekey_signature,
            &identity_key,
        )? {
            return Err(anyhow!("Invalid signed pre-key signature"));
        }
        
        sqlx::query!(
            r#"
            INSERT INTO signed_prekeys (id, user_id, device_id, prekey_id, public_key, signature, created_at)
            VALUES ($1, $2, $3, $4, $5, $6, NOW())
            "#,
            Uuid::new_v4(),
            user_id,
            request.device_id,
            request.signed_prekey.id,
            &signed_prekey_public,
            &signed_prekey_signature,
        )
        .execute(&mut *tx)
        .await?;
        
        // Insert one-time pre-keys
        for prekey in &request.onetime_prekeys {
            let public_key = base64::decode(&prekey.public_key)?;
            
            sqlx::query!(
                r#"
                INSERT INTO onetime_prekeys (id, user_id, device_id, prekey_id, public_key, is_used, created_at)
                VALUES ($1, $2, $3, $4, $5, false, NOW())
                "#,
                Uuid::new_v4(),
                user_id,
                request.device_id,
                prekey.id,
                &public_key,
            )
            .execute(&mut *tx)
            .await?;
        }
        
        // Commit transaction
        tx.commit().await?;
        
        Ok(Device {
            id: device_id_uuid,
            user_id,
            device_id: request.device_id,
            device_name: request.device_name,
            registration_id: request.registration_id,
            identity_key,
            is_active: true,
            last_seen: Utc::now(),
            created_at: Utc::now(),
        })
    }
    
    /// Get pre-key bundle for starting conversation
    pub async fn get_prekey_bundle(
        &self,
        user_id: Uuid,
        device_id: Option<String>,
    ) -> Result<PreKeyBundle> {
        // Get device (specific or any active)
        let device = if let Some(dev_id) = device_id {
            sqlx::query!(
                "SELECT * FROM devices WHERE user_id = $1 AND device_id = $2 AND is_active = true",
                user_id,
                dev_id
            )
            .fetch_one(&self.db)
            .await?
        } else {
            sqlx::query!(
                "SELECT * FROM devices WHERE user_id = $1 AND is_active = true ORDER BY last_seen DESC LIMIT 1",
                user_id
            )
            .fetch_one(&self.db)
            .await?
        };
        
        // Get signed pre-key
        let signed_prekey = sqlx::query!(
            r#"
            SELECT prekey_id, public_key, signature
            FROM signed_prekeys
            WHERE user_id = $1 AND device_id = $2
            ORDER BY created_at DESC
            LIMIT 1
            "#,
            user_id,
            device.device_id
        )
        .fetch_one(&self.db)
        .await?;
        
        // Get and mark one-time pre-key as used
        let mut tx = self.db.begin().await?;
        
        let onetime_prekey = sqlx::query!(
            r#"
            SELECT id, prekey_id, public_key
            FROM onetime_prekeys
            WHERE user_id = $1 AND device_id = $2 AND is_used = false
            ORDER BY created_at ASC
            LIMIT 1
            FOR UPDATE
            "#,
            user_id,
            device.device_id
        )
        .fetch_optional(&mut *tx)
        .await?;
        
        let (onetime_prekey_id, onetime_prekey_public) = if let Some(otpk) = onetime_prekey {
            // Mark as used
            sqlx::query!(
                "UPDATE onetime_prekeys SET is_used = true WHERE id = $1",
                otpk.id
            )
            .execute(&mut *tx)
            .await?;
            
            (Some(otpk.prekey_id), Some(otpk.public_key))
        } else {
            (None, None)
        };
        
        tx.commit().await?;
        
        // Check if we need to alert about low pre-key count
        let remaining_count = sqlx::query!(
            "SELECT COUNT(*) as count FROM onetime_prekeys WHERE user_id = $1 AND device_id = $2 AND is_used = false",
            user_id,
            device.device_id
        )
        .fetch_one(&self.db)
        .await?;
        
        if remaining_count.count.unwrap_or(0) < 20 {
            // TODO: Send notification to user to upload more pre-keys
            tracing::warn!(
                "Low pre-key count for user {} device {}: {}",
                user_id,
                device.device_id,
                remaining_count.count.unwrap_or(0)
            );
        }
        
        Ok(PreKeyBundle {
            user_id,
            device_id: device.device_id,
            registration_id: device.registration_id,
            identity_key: device.identity_key,
            signed_prekey_id: signed_prekey.prekey_id,
            signed_prekey: signed_prekey.public_key,
            signed_prekey_signature: signed_prekey.signature,
            onetime_prekey_id,
            onetime_prekey: onetime_prekey_public,
        })
    }
    
    /// Rotate signed pre-key
    pub async fn rotate_signed_prekey(
        &self,
        user_id: Uuid,
        device_id: String,
        new_prekey: SignedPreKeyUpload,
    ) -> Result<()> {
        // Get device
        let device = sqlx::query!(
            "SELECT identity_key FROM devices WHERE user_id = $1 AND device_id = $2",
            user_id,
            device_id
        )
        .fetch_one(&self.db)
        .await?;
        
        // Decode and verify
        let public_key = base64::decode(&new_prekey.public_key)?;
        let signature = base64::decode(&new_prekey.signature)?;
        
        if !SignalProtocol::verify_prekey_signature(&public_key, &signature, &device.identity_key)? {
            return Err(anyhow!("Invalid signature"));
        }
        
        // Insert new signed pre-key
        sqlx::query!(
            r#"
            INSERT INTO signed_prekeys (id, user_id, device_id, prekey_id, public_key, signature, created_at)
            VALUES ($1, $2, $3, $4, $5, $6, NOW())
            "#,
            Uuid::new_v4(),
            user_id,
            device_id,
            new_prekey.id,
            &public_key,
            &signature,
        )
        .execute(&self.db)
        .await?;
        
        // Delete old signed pre-keys (keep last 3)
        sqlx::query!(
            r#"
            DELETE FROM signed_prekeys
            WHERE user_id = $1 AND device_id = $2
            AND id NOT IN (
                SELECT id FROM signed_prekeys
                WHERE user_id = $1 AND device_id = $2
                ORDER BY created_at DESC
                LIMIT 3
            )
            "#,
            user_id,
            device_id
        )
        .execute(&self.db)
        .await?;
        
        Ok(())
    }
    
    /// Upload new batch of one-time pre-keys
    pub async fn upload_onetime_prekeys(
        &self,
        user_id: Uuid,
        device_id: String,
        prekeys: Vec<OneTimePreKeyUpload>,
    ) -> Result<()> {
        if prekeys.is_empty() {
            return Err(anyhow!("Must upload at least one pre-key"));
        }
        
        let mut tx = self.db.begin().await?;
        
        for prekey in prekeys {
            let public_key = base64::decode(&prekey.public_key)?;
            
            sqlx::query!(
                r#"
                INSERT INTO onetime_prekeys (id, user_id, device_id, prekey_id, public_key, is_used, created_at)
                VALUES ($1, $2, $3, $4, $5, false, NOW())
                ON CONFLICT (user_id, device_id, prekey_id) DO NOTHING
                "#,
                Uuid::new_v4(),
                user_id,
                device_id.clone(),
                prekey.id,
                &public_key,
            )
            .execute(&mut *tx)
            .await?;
        }
        
        tx.commit().await?;
        
        Ok(())
    }
    
    /// Deactivate device
    pub async fn deactivate_device(
        &self,
        user_id: Uuid,
        device_id: String,
    ) -> Result<()> {
        sqlx::query!(
            "UPDATE devices SET is_active = false WHERE user_id = $1 AND device_id = $2",
            user_id,
            device_id
        )
        .execute(&self.db)
        .await?;
        
        Ok(())
    }
    
    /// Get user's devices
    pub async fn get_user_devices(&self, user_id: Uuid) -> Result<Vec<Device>> {
        let rows = sqlx::query!(
            "SELECT * FROM devices WHERE user_id = $1 ORDER BY last_seen DESC",
            user_id
        )
        .fetch_all(&self.db)
        .await?;
        
        Ok(rows
            .into_iter()
            .map(|row| Device {
                id: row.id,
                user_id: row.user_id,
                device_id: row.device_id,
                device_name: row.device_name,
                registration_id: row.registration_id,
                identity_key: row.identity_key,
                is_active: row.is_active,
                last_seen: row.last_seen,
                created_at: row.created_at,
            })
            .collect())
    }
    
    /// Update device last seen
    pub async fn update_device_last_seen(
        &self,
        user_id: Uuid,
        device_id: String,
    ) -> Result<()> {
        sqlx::query!(
            "UPDATE devices SET last_seen = NOW() WHERE user_id = $1 AND device_id = $2",
            user_id,
            device_id
        )
        .execute(&self.db)
        .await?;
        
        Ok(())
    }
    
    /// Get statistics
    pub async fn get_key_stats(&self, user_id: Uuid, device_id: String) -> Result<HashMap<String, i64>> {
        let unused_count = sqlx::query!(
            "SELECT COUNT(*) as count FROM onetime_prekeys WHERE user_id = $1 AND device_id = $2 AND is_used = false",
            user_id,
            device_id
        )
        .fetch_one(&self.db)
        .await?;
        
        let total_count = sqlx::query!(
            "SELECT COUNT(*) as count FROM onetime_prekeys WHERE user_id = $1 AND device_id = $2",
            user_id,
            device_id
        )
        .fetch_one(&self.db)
        .await?;
        
        let mut stats = HashMap::new();
        stats.insert("unused_prekeys".to_string(), unused_count.count.unwrap_or(0));
        stats.insert("total_prekeys".to_string(), total_count.count.unwrap_or(0));
        stats.insert("used_prekeys".to_string(), total_count.count.unwrap_or(0) - unused_count.count.unwrap_or(0));
        
        Ok(stats)
    }
}
