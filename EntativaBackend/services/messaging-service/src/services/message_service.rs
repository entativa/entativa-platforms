use anyhow::{anyhow, Result};
use sqlx::PgPool;
use redis::AsyncCommands;
use uuid::Uuid;
use chrono::{Utc, Duration};

use crate::models::message::*;
use crate::crypto::signal::SignalProtocol;

/// Message Service
/// Handles message routing, offline queue, and delivery tracking
pub struct MessageService {
    db: PgPool,
    redis: redis::Client,
}

impl MessageService {
    pub fn new(db: PgPool, redis: redis::Client) -> Self {
        Self { db, redis }
    }
    
    /// Send 1:1 message
    pub async fn send_message(
        &self,
        sender_id: Uuid,
        request: SendMessageRequest,
    ) -> Result<MessageResponse> {
        // Validate
        if request.recipient_id.is_none() && request.conversation_id == Uuid::nil() {
            return Err(anyhow!("Must specify recipient_id or conversation_id"));
        }
        
        // Decode ciphertext and ephemeral key
        let ciphertext = base64::decode(&request.ciphertext)?;
        let ephemeral_key = base64::decode(&request.ephemeral_key)?;
        
        // Get or create conversation
        let conversation_id = if request.conversation_id == Uuid::nil() {
            self.get_or_create_conversation(sender_id, request.recipient_id.unwrap()).await?
        } else {
            request.conversation_id
        };
        
        // Get next sequence number
        let sequence_number = self.get_next_sequence_number(conversation_id).await?;
        
        // Calculate expiry time
        let expires_at = if request.is_self_destructing {
            request.expires_in_seconds.map(|secs| Utc::now() + Duration::seconds(secs))
        } else {
            None
        };
        
        // Create message
        let message_id = Uuid::new_v4();
        let message = Message {
            id: message_id,
            conversation_id,
            sender_id,
            sender_device_id: request.device_id.clone(),
            recipient_id: request.recipient_id,
            message_type: request.message_type.clone(),
            ciphertext,
            ephemeral_key,
            sequence_number,
            timestamp: Utc::now(),
            is_group: false,
            group_epoch: None,
            status: MessageStatus::Sent,
            delivered_at: None,
            read_at: None,
            is_self_destructing: request.is_self_destructing,
            expires_at,
            is_edited: false,
            edited_at: None,
            created_at: Utc::now(),
        };
        
        // Store message
        sqlx::query!(
            r#"
            INSERT INTO messages (
                id, conversation_id, sender_id, sender_device_id, recipient_id,
                message_type, ciphertext, ephemeral_key, sequence_number, timestamp,
                is_group, status, is_self_destructing, expires_at, created_at
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
            "#,
            message.id,
            message.conversation_id,
            message.sender_id,
            message.sender_device_id,
            message.recipient_id,
            serde_json::to_value(&message.message_type)?,
            &message.ciphertext,
            &message.ephemeral_key,
            message.sequence_number,
            message.timestamp,
            message.is_group,
            serde_json::to_value(&message.status)?,
            message.is_self_destructing,
            message.expires_at,
            message.created_at,
        )
        .execute(&self.db)
        .await?;
        
        // Queue for delivery
        self.queue_message_for_delivery(&message).await?;
        
        // Update conversation last message
        self.update_conversation_last_message(conversation_id, message_id).await?;
        
        // Publish event for real-time delivery
        self.publish_message_event(&message).await?;
        
        Ok(MessageResponse {
            message_id,
            conversation_id,
            sequence_number,
            timestamp: message.timestamp,
            status: message.status,
        })
    }
    
    /// Get messages for conversation
    pub async fn get_messages(
        &self,
        user_id: Uuid,
        conversation_id: Uuid,
        before_sequence: Option<i64>,
        limit: i32,
    ) -> Result<Vec<Message>> {
        // Verify user is in conversation
        self.verify_conversation_access(user_id, conversation_id).await?;
        
        let limit = limit.min(100); // Max 100 messages per request
        
        let rows = if let Some(before) = before_sequence {
            sqlx::query!(
                r#"
                SELECT * FROM messages
                WHERE conversation_id = $1 AND sequence_number < $2
                ORDER BY sequence_number DESC
                LIMIT $3
                "#,
                conversation_id,
                before,
                limit as i64
            )
            .fetch_all(&self.db)
            .await?
        } else {
            sqlx::query!(
                r#"
                SELECT * FROM messages
                WHERE conversation_id = $1
                ORDER BY sequence_number DESC
                LIMIT $2
                "#,
                conversation_id,
                limit as i64
            )
            .fetch_all(&self.db)
            .await?
        };
        
        Ok(rows
            .into_iter()
            .map(|row| Message {
                id: row.id,
                conversation_id: row.conversation_id,
                sender_id: row.sender_id,
                sender_device_id: row.sender_device_id,
                recipient_id: row.recipient_id,
                message_type: serde_json::from_value(row.message_type).unwrap_or(MessageType::Text),
                ciphertext: row.ciphertext,
                ephemeral_key: row.ephemeral_key,
                sequence_number: row.sequence_number,
                timestamp: row.timestamp,
                is_group: row.is_group,
                group_epoch: row.group_epoch.map(|e| e as u64),
                status: serde_json::from_value(row.status).unwrap_or(MessageStatus::Sent),
                delivered_at: row.delivered_at,
                read_at: row.read_at,
                is_self_destructing: row.is_self_destructing,
                expires_at: row.expires_at,
                is_edited: row.is_edited,
                edited_at: row.edited_at,
                created_at: row.created_at,
            })
            .collect())
    }
    
    /// Mark message as delivered
    pub async fn mark_delivered(
        &self,
        message_id: Uuid,
        user_id: Uuid,
    ) -> Result<()> {
        sqlx::query!(
            r#"
            UPDATE messages
            SET status = $1, delivered_at = NOW()
            WHERE id = $2 AND recipient_id = $3 AND status != $4
            "#,
            serde_json::to_value(&MessageStatus::Delivered)?,
            message_id,
            user_id,
            serde_json::to_value(&MessageStatus::Read)?,
        )
        .execute(&self.db)
        .await?;
        
        // Publish delivery receipt
        self.publish_delivery_receipt(message_id, user_id, false).await?;
        
        Ok(())
    }
    
    /// Mark message as read
    pub async fn mark_read(
        &self,
        message_id: Uuid,
        user_id: Uuid,
    ) -> Result<()> {
        sqlx::query!(
            r#"
            UPDATE messages
            SET status = $1, read_at = NOW()
            WHERE id = $2 AND recipient_id = $3
            "#,
            serde_json::to_value(&MessageStatus::Read)?,
            message_id,
            user_id,
        )
        .execute(&self.db)
        .await?;
        
        // Publish read receipt
        self.publish_delivery_receipt(message_id, user_id, true).await?;
        
        Ok(())
    }
    
    /// Delete message (for self)
    pub async fn delete_message(
        &self,
        message_id: Uuid,
        user_id: Uuid,
    ) -> Result<()> {
        // In E2EE, we can't truly delete for everyone
        // Just mark as deleted for this user
        sqlx::query!(
            "INSERT INTO deleted_messages (user_id, message_id, deleted_at) VALUES ($1, $2, NOW())",
            user_id,
            message_id
        )
        .execute(&self.db)
        .await?;
        
        Ok(())
    }
    
    /// Get offline queue for user
    pub async fn get_offline_queue(&self, user_id: Uuid, device_id: String) -> Result<Vec<Message>> {
        let mut conn = self.redis.get_async_connection().await?;
        
        // Get message IDs from queue
        let queue_key = format!("offline_queue:{}:{}", user_id, device_id);
        let message_ids: Vec<String> = conn.lrange(&queue_key, 0, 100).await?;
        
        if message_ids.is_empty() {
            return Ok(Vec::new());
        }
        
        // Fetch messages from database
        let uuids: Vec<Uuid> = message_ids
            .iter()
            .filter_map(|id| Uuid::parse_str(id).ok())
            .collect();
        
        let rows = sqlx::query!(
            "SELECT * FROM messages WHERE id = ANY($1) ORDER BY created_at ASC",
            &uuids
        )
        .fetch_all(&self.db)
        .await?;
        
        Ok(rows
            .into_iter()
            .map(|row| Message {
                id: row.id,
                conversation_id: row.conversation_id,
                sender_id: row.sender_id,
                sender_device_id: row.sender_device_id,
                recipient_id: row.recipient_id,
                message_type: serde_json::from_value(row.message_type).unwrap_or(MessageType::Text),
                ciphertext: row.ciphertext,
                ephemeral_key: row.ephemeral_key,
                sequence_number: row.sequence_number,
                timestamp: row.timestamp,
                is_group: row.is_group,
                group_epoch: row.group_epoch.map(|e| e as u64),
                status: serde_json::from_value(row.status).unwrap_or(MessageStatus::Sent),
                delivered_at: row.delivered_at,
                read_at: row.read_at,
                is_self_destructing: row.is_self_destructing,
                expires_at: row.expires_at,
                is_edited: row.is_edited,
                edited_at: row.edited_at,
                created_at: row.created_at,
            })
            .collect())
    }
    
    /// Clear offline queue after delivery
    pub async fn clear_offline_queue(
        &self,
        user_id: Uuid,
        device_id: String,
        message_ids: Vec<Uuid>,
    ) -> Result<()> {
        let mut conn = self.redis.get_async_connection().await?;
        
        let queue_key = format!("offline_queue:{}:{}", user_id, device_id);
        
        for message_id in message_ids {
            let _: () = conn.lrem(&queue_key, 1, message_id.to_string()).await?;
        }
        
        Ok(())
    }
    
    // Helper methods
    
    async fn get_or_create_conversation(&self, user1_id: Uuid, user2_id: Uuid) -> Result<Uuid> {
        // Check if conversation exists
        let existing = sqlx::query!(
            r#"
            SELECT c.id FROM conversations c
            JOIN conversation_participants cp1 ON c.id = cp1.conversation_id AND cp1.user_id = $1
            JOIN conversation_participants cp2 ON c.id = cp2.conversation_id AND cp2.user_id = $2
            WHERE c.conversation_type = $3
            LIMIT 1
            "#,
            user1_id,
            user2_id,
            serde_json::to_value(&ConversationType::OneToOne)?
        )
        .fetch_optional(&self.db)
        .await?;
        
        if let Some(row) = existing {
            return Ok(row.id);
        }
        
        // Create new conversation
        let conversation_id = Uuid::new_v4();
        let mut tx = self.db.begin().await?;
        
        sqlx::query!(
            r#"
            INSERT INTO conversations (id, conversation_type, created_by, created_at, updated_at)
            VALUES ($1, $2, $3, NOW(), NOW())
            "#,
            conversation_id,
            serde_json::to_value(&ConversationType::OneToOne)?,
            user1_id,
        )
        .execute(&mut *tx)
        .await?;
        
        // Add participants
        for user_id in &[user1_id, user2_id] {
            sqlx::query!(
                "INSERT INTO conversation_participants (conversation_id, user_id, joined_at) VALUES ($1, $2, NOW())",
                conversation_id,
                user_id
            )
            .execute(&mut *tx)
            .await?;
        }
        
        tx.commit().await?;
        
        Ok(conversation_id)
    }
    
    async fn get_next_sequence_number(&self, conversation_id: Uuid) -> Result<i64> {
        let row = sqlx::query!(
            "SELECT COALESCE(MAX(sequence_number), 0) + 1 as next_seq FROM messages WHERE conversation_id = $1",
            conversation_id
        )
        .fetch_one(&self.db)
        .await?;
        
        Ok(row.next_seq.unwrap_or(1))
    }
    
    async fn queue_message_for_delivery(&self, message: &Message) -> Result<()> {
        if let Some(recipient_id) = message.recipient_id {
            let mut conn = self.redis.get_async_connection().await?;
            
            // Get recipient devices
            let devices = sqlx::query!(
                "SELECT device_id FROM devices WHERE user_id = $1 AND is_active = true",
                recipient_id
            )
            .fetch_all(&self.db)
            .await?;
            
            // Queue for each device
            for device in devices {
                let queue_key = format!("offline_queue:{}:{}", recipient_id, device.device_id);
                let _: () = conn.rpush(&queue_key, message.id.to_string()).await?;
                
                // Set TTL (30 days)
                let _: () = conn.expire(&queue_key, 30 * 24 * 60 * 60).await?;
            }
        }
        
        Ok(())
    }
    
    async fn update_conversation_last_message(&self, conversation_id: Uuid, message_id: Uuid) -> Result<()> {
        sqlx::query!(
            "UPDATE conversations SET last_message_id = $1, updated_at = NOW() WHERE id = $2",
            message_id,
            conversation_id
        )
        .execute(&self.db)
        .await?;
        
        Ok(())
    }
    
    async fn publish_message_event(&self, message: &Message) -> Result<()> {
        let mut conn = self.redis.get_async_connection().await?;
        
        // Publish to Redis pub/sub for WebSocket delivery
        let channel = format!("messages:user:{}", message.recipient_id.unwrap_or(Uuid::nil()));
        let payload = serde_json::to_string(message)?;
        
        let _: () = conn.publish(&channel, payload).await?;
        
        Ok(())
    }
    
    async fn publish_delivery_receipt(&self, message_id: Uuid, user_id: Uuid, is_read: bool) -> Result<()> {
        let mut conn = self.redis.get_async_connection().await?;
        
        // Get message to find sender
        let message = sqlx::query!("SELECT sender_id FROM messages WHERE id = $1", message_id)
            .fetch_one(&self.db)
            .await?;
        
        // Publish receipt to sender
        let channel = format!("receipts:user:{}", message.sender_id);
        let receipt = serde_json::json!({
            "message_id": message_id,
            "user_id": user_id,
            "is_read": is_read,
            "timestamp": Utc::now(),
        });
        
        let _: () = conn.publish(&channel, receipt.to_string()).await?;
        
        Ok(())
    }
    
    async fn verify_conversation_access(&self, user_id: Uuid, conversation_id: Uuid) -> Result<()> {
        let exists = sqlx::query!(
            "SELECT 1 FROM conversation_participants WHERE conversation_id = $1 AND user_id = $2",
            conversation_id,
            user_id
        )
        .fetch_optional(&self.db)
        .await?;
        
        if exists.is_none() {
            return Err(anyhow!("User not in conversation"));
        }
        
        Ok(())
    }
}
