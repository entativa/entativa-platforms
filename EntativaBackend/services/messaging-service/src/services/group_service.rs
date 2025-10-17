use anyhow::{anyhow, Result};
use sqlx::PgPool;
use redis::AsyncCommands;
use uuid::Uuid;
use chrono::Utc;

use crate::models::message::*;
use crate::models::keys::*;
use crate::crypto::mls::{MLSProtocol, MLSGroup};

/// Group Message Service
/// Handles MLS group management and group messaging
pub struct GroupService {
    db: PgPool,
    redis: redis::Client,
}

impl GroupService {
    pub fn new(db: PgPool, redis: redis::Client) -> Self {
        Self { db, redis }
    }
    
    /// Create group chat
    pub async fn create_group(
        &self,
        creator_id: Uuid,
        name: String,
        description: Option<String>,
        member_ids: Vec<Uuid>,
    ) -> Result<(GroupChat, MLSGroup)> {
        // Validate
        if name.is_empty() || name.len() > 100 {
            return Err(anyhow!("Group name must be 1-100 characters"));
        }
        
        if member_ids.len() > 1499 {
            return Err(anyhow!("Max 1,499 members (plus creator = 1,500 total)"));
        }
        
        // Get creator's identity key
        let creator_device = sqlx::query!(
            "SELECT identity_key FROM devices WHERE user_id = $1 AND is_active = true ORDER BY last_seen DESC LIMIT 1",
            creator_id
        )
        .fetch_one(&self.db)
        .await?;
        
        // Create MLS group
        let group_id = Uuid::new_v4();
        let mls_group = MLSProtocol::create_group(
            group_id,
            creator_id,
            &creator_device.identity_key,
        )?;
        
        // Create group in database
        let mut tx = self.db.begin().await?;
        
        let group = GroupChat {
            id: group_id,
            name: name.clone(),
            description: description.clone(),
            avatar_url: None,
            created_by: creator_id,
            member_count: 1, // Just creator for now
            max_members: 1500,
            is_public: false,
            invite_link: None,
            mls_group_id: group_id.as_bytes().to_vec(),
            current_epoch: mls_group.epoch,
            created_at: Utc::now(),
            updated_at: Utc::now(),
        };
        
        sqlx::query!(
            r#"
            INSERT INTO group_chats (
                id, name, description, created_by, member_count, max_members,
                is_public, mls_group_id, current_epoch, created_at, updated_at
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
            "#,
            group.id,
            group.name,
            group.description,
            group.created_by,
            group.member_count,
            group.max_members,
            group.is_public,
            &group.mls_group_id,
            group.current_epoch as i64,
        )
        .execute(&mut *tx)
        .await?;
        
        // Add creator as owner
        sqlx::query!(
            r#"
            INSERT INTO group_members (group_id, user_id, role, joined_at, added_by)
            VALUES ($1, $2, $3, NOW(), $2)
            "#,
            group.id,
            creator_id,
            serde_json::to_value(&GroupMemberRole::Owner)?,
        )
        .execute(&mut *tx)
        .await?;
        
        // Create conversation
        let conversation_id = Uuid::new_v4();
        sqlx::query!(
            r#"
            INSERT INTO conversations (id, conversation_type, name, created_by, created_at, updated_at)
            VALUES ($1, $2, $3, $4, NOW(), NOW())
            "#,
            conversation_id,
            serde_json::to_value(&ConversationType::Group)?,
            name,
            creator_id,
        )
        .execute(&mut *tx)
        .await?;
        
        // Link conversation to group
        sqlx::query!(
            "UPDATE group_chats SET conversation_id = $1 WHERE id = $2",
            conversation_id,
            group.id
        )
        .execute(&mut *tx)
        .await?;
        
        // Add creator to conversation
        sqlx::query!(
            "INSERT INTO conversation_participants (conversation_id, user_id, joined_at) VALUES ($1, $2, NOW())",
            conversation_id,
            creator_id
        )
        .execute(&mut *tx)
        .await?;
        
        // Store MLS group state
        self.store_mls_group_state(&mut tx, &mls_group).await?;
        
        tx.commit().await?;
        
        // Add members
        let mut updated_group = mls_group;
        for member_id in member_ids {
            if let Ok(_) = self.add_member_internal(&mut updated_group, group_id, member_id, creator_id).await {
                // Member added successfully
            }
        }
        
        Ok((group, updated_group))
    }
    
    /// Add member to group
    pub async fn add_member(
        &self,
        group_id: Uuid,
        user_id: Uuid,
        added_by: Uuid,
    ) -> Result<()> {
        // Verify adder is admin or owner
        self.verify_admin_access(group_id, added_by).await?;
        
        // Load MLS group
        let mut mls_group = self.load_mls_group_state(group_id).await?;
        
        // Add member
        self.add_member_internal(&mut mls_group, group_id, user_id, added_by).await?;
        
        Ok(())
    }
    
    async fn add_member_internal(
        &self,
        mls_group: &mut MLSGroup,
        group_id: Uuid,
        user_id: Uuid,
        added_by: Uuid,
    ) -> Result<()> {
        // Check if already member
        let exists = sqlx::query!(
            "SELECT 1 FROM group_members WHERE group_id = $1 AND user_id = $2",
            group_id,
            user_id
        )
        .fetch_optional(&self.db)
        .await?;
        
        if exists.is_some() {
            return Err(anyhow!("User already in group"));
        }
        
        // Check group size
        if mls_group.member_map.len() >= 1500 {
            return Err(anyhow!("Group is full (max 1,500 members)"));
        }
        
        // Get user's identity key
        let user_device = sqlx::query!(
            "SELECT identity_key FROM devices WHERE user_id = $1 AND is_active = true ORDER BY last_seen DESC LIMIT 1",
            user_id
        )
        .fetch_one(&self.db)
        .await?;
        
        // Add to MLS group
        let welcome_message = MLSProtocol::add_member(mls_group, user_id, &user_device.identity_key)?;
        
        // Update database
        let mut tx = self.db.begin().await?;
        
        // Add member
        sqlx::query!(
            r#"
            INSERT INTO group_members (group_id, user_id, role, joined_at, added_by)
            VALUES ($1, $2, $3, NOW(), $4)
            "#,
            group_id,
            user_id,
            serde_json::to_value(&GroupMemberRole::Member)?,
            added_by,
        )
        .execute(&mut *tx)
        .await?;
        
        // Update member count and epoch
        sqlx::query!(
            "UPDATE group_chats SET member_count = member_count + 1, current_epoch = $1, updated_at = NOW() WHERE id = $2",
            mls_group.epoch as i64,
            group_id
        )
        .execute(&mut *tx)
        .await?;
        
        // Add to conversation
        let conversation_id = sqlx::query!(
            "SELECT conversation_id FROM group_chats WHERE id = $1",
            group_id
        )
        .fetch_one(&mut *tx)
        .await?;
        
        if let Some(conv_id) = conversation_id.conversation_id {
            sqlx::query!(
                "INSERT INTO conversation_participants (conversation_id, user_id, joined_at) VALUES ($1, $2, NOW())",
                conv_id,
                user_id
            )
            .execute(&mut *tx)
            .await?;
        }
        
        // Store updated MLS state
        self.store_mls_group_state(&mut *tx, mls_group).await?;
        
        // Store welcome message for new member
        self.store_welcome_message(&mut *tx, group_id, user_id, &welcome_message).await?;
        
        tx.commit().await?;
        
        // Send system message
        self.send_system_message(group_id, format!("User added to group")).await?;
        
        Ok(())
    }
    
    /// Remove member from group
    pub async fn remove_member(
        &self,
        group_id: Uuid,
        user_id: Uuid,
        removed_by: Uuid,
    ) -> Result<()> {
        // Verify remover is admin or owner
        self.verify_admin_access(group_id, removed_by).await?;
        
        // Can't remove owner
        let member = sqlx::query!(
            "SELECT role FROM group_members WHERE group_id = $1 AND user_id = $2",
            group_id,
            user_id
        )
        .fetch_one(&self.db)
        .await?;
        
        let role: GroupMemberRole = serde_json::from_value(member.role)?;
        if role == GroupMemberRole::Owner {
            return Err(anyhow!("Cannot remove group owner"));
        }
        
        // Load MLS group
        let mut mls_group = self.load_mls_group_state(group_id).await?;
        
        // Remove from MLS group
        MLSProtocol::remove_member(&mut mls_group, user_id)?;
        
        // Update database
        let mut tx = self.db.begin().await?;
        
        // Remove member
        sqlx::query!(
            "DELETE FROM group_members WHERE group_id = $1 AND user_id = $2",
            group_id,
            user_id
        )
        .execute(&mut *tx)
        .await?;
        
        // Update member count and epoch
        sqlx::query!(
            "UPDATE group_chats SET member_count = member_count - 1, current_epoch = $1, updated_at = NOW() WHERE id = $2",
            mls_group.epoch as i64,
            group_id
        )
        .execute(&mut *tx)
        .await?;
        
        // Remove from conversation
        let conversation_id = sqlx::query!(
            "SELECT conversation_id FROM group_chats WHERE id = $1",
            group_id
        )
        .fetch_one(&mut *tx)
        .await?;
        
        if let Some(conv_id) = conversation_id.conversation_id {
            sqlx::query!(
                "DELETE FROM conversation_participants WHERE conversation_id = $1 AND user_id = $2",
                conv_id,
                user_id
            )
            .execute(&mut *tx)
            .await?;
        }
        
        // Store updated MLS state
        self.store_mls_group_state(&mut *tx, &mls_group).await?;
        
        tx.commit().await?;
        
        // Send system message
        self.send_system_message(group_id, format!("User removed from group")).await?;
        
        Ok(())
    }
    
    /// Send group message
    pub async fn send_group_message(
        &self,
        sender_id: Uuid,
        group_id: Uuid,
        ciphertext: Vec<u8>,
    ) -> Result<MessageResponse> {
        // Verify sender is member
        self.verify_member_access(group_id, sender_id).await?;
        
        // Get conversation
        let conversation_id = sqlx::query!(
            "SELECT conversation_id FROM group_chats WHERE id = $1",
            group_id
        )
        .fetch_one(&self.db)
        .await?
        .conversation_id
        .ok_or_else(|| anyhow!("Group has no conversation"))?;
        
        // Get current epoch
        let epoch = sqlx::query!(
            "SELECT current_epoch FROM group_chats WHERE id = $1",
            group_id
        )
        .fetch_one(&self.db)
        .await?
        .current_epoch as u64;
        
        // Get next sequence number
        let sequence_number = sqlx::query!(
            "SELECT COALESCE(MAX(sequence_number), 0) + 1 as next_seq FROM messages WHERE conversation_id = $1",
            conversation_id
        )
        .fetch_one(&self.db)
        .await?
        .next_seq
        .unwrap_or(1);
        
        // Create message
        let message_id = Uuid::new_v4();
        let message = Message {
            id: message_id,
            conversation_id,
            sender_id,
            sender_device_id: "".to_string(), // TODO: Get from request
            recipient_id: None,
            message_type: MessageType::Text,
            ciphertext,
            ephemeral_key: vec![],
            sequence_number,
            timestamp: Utc::now(),
            is_group: true,
            group_epoch: Some(epoch),
            status: MessageStatus::Sent,
            delivered_at: None,
            read_at: None,
            is_self_destructing: false,
            expires_at: None,
            is_edited: false,
            edited_at: None,
            created_at: Utc::now(),
        };
        
        // Store message
        sqlx::query!(
            r#"
            INSERT INTO messages (
                id, conversation_id, sender_id, message_type, ciphertext,
                sequence_number, timestamp, is_group, group_epoch, status, created_at
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
            "#,
            message.id,
            message.conversation_id,
            message.sender_id,
            serde_json::to_value(&message.message_type)?,
            &message.ciphertext,
            message.sequence_number,
            message.timestamp,
            message.is_group,
            message.group_epoch.map(|e| e as i64),
            serde_json::to_value(&message.status)?,
            message.created_at,
        )
        .execute(&self.db)
        .await?;
        
        // Queue for all members
        self.queue_group_message(group_id, &message).await?;
        
        Ok(MessageResponse {
            message_id,
            conversation_id,
            sequence_number,
            timestamp: message.timestamp,
            status: message.status,
        })
    }
    
    // Helper methods
    
    async fn verify_admin_access(&self, group_id: Uuid, user_id: Uuid) -> Result<()> {
        let member = sqlx::query!(
            "SELECT role FROM group_members WHERE group_id = $1 AND user_id = $2",
            group_id,
            user_id
        )
        .fetch_one(&self.db)
        .await?;
        
        let role: GroupMemberRole = serde_json::from_value(member.role)?;
        if role != GroupMemberRole::Owner && role != GroupMemberRole::Admin {
            return Err(anyhow!("Insufficient permissions"));
        }
        
        Ok(())
    }
    
    async fn verify_member_access(&self, group_id: Uuid, user_id: Uuid) -> Result<()> {
        let exists = sqlx::query!(
            "SELECT 1 FROM group_members WHERE group_id = $1 AND user_id = $2",
            group_id,
            user_id
        )
        .fetch_optional(&self.db)
        .await?;
        
        if exists.is_none() {
            return Err(anyhow!("User not in group"));
        }
        
        Ok(())
    }
    
    async fn load_mls_group_state(&self, group_id: Uuid) -> Result<MLSGroup> {
        // Load from cache first
        let mut conn = self.redis.get_async_connection().await?;
        let cache_key = format!("mls_group:{}", group_id);
        
        if let Ok(cached) = conn.get::<_, Vec<u8>>(&cache_key).await {
            if let Ok(group) = bincode::deserialize(&cached) {
                return Ok(group);
            }
        }
        
        // Load from database
        let row = sqlx::query!(
            "SELECT mls_group_state FROM mls_group_states WHERE group_id = $1",
            group_id
        )
        .fetch_one(&self.db)
        .await?;
        
        let group: MLSGroup = bincode::deserialize(&row.mls_group_state)?;
        
        // Cache for 1 hour
        let serialized = bincode::serialize(&group)?;
        let _: () = conn.set_ex(&cache_key, serialized, 3600).await?;
        
        Ok(group)
    }
    
    async fn store_mls_group_state<'a, E>(&self, executor: E, mls_group: &MLSGroup) -> Result<()>
    where
        E: sqlx::Executor<'a, Database = sqlx::Postgres>,
    {
        let serialized = bincode::serialize(mls_group)?;
        
        sqlx::query!(
            r#"
            INSERT INTO mls_group_states (group_id, epoch, mls_group_state, updated_at)
            VALUES ($1, $2, $3, NOW())
            ON CONFLICT (group_id) DO UPDATE SET
                epoch = $2,
                mls_group_state = $3,
                updated_at = NOW()
            "#,
            mls_group.group_id,
            mls_group.epoch as i64,
            &serialized,
        )
        .execute(executor)
        .await?;
        
        // Update cache
        let mut conn = self.redis.get_async_connection().await?;
        let cache_key = format!("mls_group:{}", mls_group.group_id);
        let _: () = conn.set_ex(&cache_key, serialized, 3600).await?;
        
        Ok(())
    }
    
    async fn store_welcome_message<'a, E>(
        &self,
        executor: E,
        group_id: Uuid,
        user_id: Uuid,
        welcome_message: &[u8],
    ) -> Result<()>
    where
        E: sqlx::Executor<'a, Database = sqlx::Postgres>,
    {
        sqlx::query!(
            "INSERT INTO mls_welcome_messages (group_id, user_id, welcome_message, created_at) VALUES ($1, $2, $3, NOW())",
            group_id,
            user_id,
            welcome_message,
        )
        .execute(executor)
        .await?;
        
        Ok(())
    }
    
    async fn queue_group_message(&self, group_id: Uuid, message: &Message) -> Result<()> {
        // Get all group members
        let members = sqlx::query!(
            "SELECT user_id FROM group_members WHERE group_id = $1",
            group_id
        )
        .fetch_all(&self.db)
        .await?;
        
        let mut conn = self.redis.get_async_connection().await?;
        
        for member in members {
            if member.user_id == message.sender_id {
                continue; // Don't queue for sender
            }
            
            // Get member devices
            let devices = sqlx::query!(
                "SELECT device_id FROM devices WHERE user_id = $1 AND is_active = true",
                member.user_id
            )
            .fetch_all(&self.db)
            .await?;
            
            // Queue for each device
            for device in devices {
                let queue_key = format!("offline_queue:{}:{}", member.user_id, device.device_id);
                let _: () = conn.rpush(&queue_key, message.id.to_string()).await?;
            }
        }
        
        Ok(())
    }
    
    async fn send_system_message(&self, group_id: Uuid, text: String) -> Result<()> {
        // Get conversation
        let conversation_id = sqlx::query!(
            "SELECT conversation_id FROM group_chats WHERE id = $1",
            group_id
        )
        .fetch_one(&self.db)
        .await?
        .conversation_id
        .ok_or_else(|| anyhow!("Group has no conversation"))?;
        
        // Create system message
        let message_id = Uuid::new_v4();
        sqlx::query!(
            r#"
            INSERT INTO messages (
                id, conversation_id, sender_id, message_type, ciphertext,
                sequence_number, timestamp, is_group, status, created_at
            ) VALUES ($1, $2, $3, $4, $5, 0, NOW(), true, $6, NOW())
            "#,
            message_id,
            conversation_id,
            Uuid::nil(), // System message
            serde_json::to_value(&MessageType::System)?,
            text.as_bytes(),
            serde_json::to_value(&MessageStatus::Sent)?,
        )
        .execute(&self.db)
        .await?;
        
        Ok(())
    }
}
