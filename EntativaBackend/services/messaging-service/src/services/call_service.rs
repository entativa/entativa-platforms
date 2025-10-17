use anyhow::{anyhow, Result};
use sqlx::PgPool;
use redis::AsyncCommands;
use uuid::Uuid;
use chrono::Utc;

use crate::models::message::{Call, CallType, CallStatus};

/// Call Service
/// WebRTC signaling for audio/video calls
pub struct CallService {
    db: PgPool,
    redis: redis::Client,
}

impl CallService {
    pub fn new(db: PgPool, redis: redis::Client) -> Self {
        Self { db, redis }
    }
    
    /// Initiate call
    pub async fn initiate_call(
        &self,
        caller_id: Uuid,
        conversation_id: Uuid,
        call_type: CallType,
        sdp_offer: String,
    ) -> Result<Call> {
        let call_id = Uuid::new_v4();
        
        // Create call record
        let call = Call {
            id: call_id,
            conversation_id,
            caller_id,
            call_type: call_type.clone(),
            status: CallStatus::Ringing,
            sdp_offer: Some(sdp_offer.clone()),
            sdp_answer: None,
            ice_candidates: Vec::new(),
            started_at: None,
            ended_at: None,
            duration_seconds: None,
            created_at: Utc::now(),
        };
        
        // Store in database
        sqlx::query!(
            r#"
            INSERT INTO calls (id, conversation_id, caller_id, call_type, status, sdp_offer, created_at)
            VALUES ($1, $2, $3, $4, $5, $6, NOW())
            "#,
            call.id,
            call.conversation_id,
            call.caller_id,
            serde_json::to_value(&call.call_type)?,
            serde_json::to_value(&call.status)?,
            &sdp_offer,
        )
        .execute(&self.db)
        .await?;
        
        // Publish call event
        self.publish_call_event(&call, "call_initiated").await?;
        
        Ok(call)
    }
    
    /// Answer call
    pub async fn answer_call(
        &self,
        call_id: Uuid,
        user_id: Uuid,
        sdp_answer: String,
    ) -> Result<()> {
        // Update call
        sqlx::query!(
            r#"
            UPDATE calls
            SET status = $1, sdp_answer = $2, started_at = NOW()
            WHERE id = $3
            "#,
            serde_json::to_value(&CallStatus::Answered)?,
            &sdp_answer,
            call_id,
        )
        .execute(&self.db)
        .await?;
        
        // Get call for event
        let call = self.get_call(call_id).await?;
        
        // Publish answer event
        self.publish_call_event(&call, "call_answered").await?;
        
        Ok(())
    }
    
    /// Decline call
    pub async fn decline_call(&self, call_id: Uuid, user_id: Uuid) -> Result<()> {
        sqlx::query!(
            r#"
            UPDATE calls
            SET status = $1, ended_at = NOW()
            WHERE id = $2
            "#,
            serde_json::to_value(&CallStatus::Declined)?,
            call_id,
        )
        .execute(&self.db)
        .await?;
        
        let call = self.get_call(call_id).await?;
        self.publish_call_event(&call, "call_declined").await?;
        
        Ok(())
    }
    
    /// End call
    pub async fn end_call(&self, call_id: Uuid, user_id: Uuid) -> Result<()> {
        // Calculate duration
        let call = self.get_call(call_id).await?;
        
        let duration_seconds = if let Some(started_at) = call.started_at {
            Some((Utc::now() - started_at).num_seconds() as i32)
        } else {
            None
        };
        
        sqlx::query!(
            r#"
            UPDATE calls
            SET status = $1, ended_at = NOW(), duration_seconds = $2
            WHERE id = $3
            "#,
            serde_json::to_value(&CallStatus::Ended)?,
            duration_seconds,
            call_id,
        )
        .execute(&self.db)
        .await?;
        
        let call = self.get_call(call_id).await?;
        self.publish_call_event(&call, "call_ended").await?;
        
        Ok(())
    }
    
    /// Add ICE candidate
    pub async fn add_ice_candidate(
        &self,
        call_id: Uuid,
        user_id: Uuid,
        candidate: String,
    ) -> Result<()> {
        // Store ICE candidate
        sqlx::query!(
            r#"
            INSERT INTO call_ice_candidates (id, call_id, user_id, candidate, created_at)
            VALUES ($1, $2, $3, $4, NOW())
            "#,
            Uuid::new_v4(),
            call_id,
            user_id,
            &candidate,
        )
        .execute(&self.db)
        .await?;
        
        // Publish ICE candidate
        self.publish_ice_candidate(call_id, user_id, candidate).await?;
        
        Ok(())
    }
    
    /// Get ICE candidates for call
    pub async fn get_ice_candidates(&self, call_id: Uuid) -> Result<Vec<String>> {
        let rows = sqlx::query!(
            "SELECT candidate FROM call_ice_candidates WHERE call_id = $1 ORDER BY created_at ASC",
            call_id
        )
        .fetch_all(&self.db)
        .await?;
        
        Ok(rows.into_iter().map(|row| row.candidate).collect())
    }
    
    /// Get call by ID
    pub async fn get_call(&self, call_id: Uuid) -> Result<Call> {
        let row = sqlx::query!(
            "SELECT * FROM calls WHERE id = $1",
            call_id
        )
        .fetch_one(&self.db)
        .await?;
        
        Ok(Call {
            id: row.id,
            conversation_id: row.conversation_id,
            caller_id: row.caller_id,
            call_type: serde_json::from_value(row.call_type)?,
            status: serde_json::from_value(row.status)?,
            sdp_offer: row.sdp_offer,
            sdp_answer: row.sdp_answer,
            ice_candidates: Vec::new(), // Load separately if needed
            started_at: row.started_at,
            ended_at: row.ended_at,
            duration_seconds: row.duration_seconds,
            created_at: row.created_at,
        })
    }
    
    /// Get active call for conversation
    pub async fn get_active_call(&self, conversation_id: Uuid) -> Result<Option<Call>> {
        let row = sqlx::query!(
            r#"
            SELECT * FROM calls
            WHERE conversation_id = $1
            AND status IN ('Ringing', 'Answered')
            ORDER BY created_at DESC
            LIMIT 1
            "#,
            conversation_id
        )
        .fetch_optional(&self.db)
        .await?;
        
        if let Some(row) = row {
            Ok(Some(Call {
                id: row.id,
                conversation_id: row.conversation_id,
                caller_id: row.caller_id,
                call_type: serde_json::from_value(row.call_type)?,
                status: serde_json::from_value(row.status)?,
                sdp_offer: row.sdp_offer,
                sdp_answer: row.sdp_answer,
                ice_candidates: Vec::new(),
                started_at: row.started_at,
                ended_at: row.ended_at,
                duration_seconds: row.duration_seconds,
                created_at: row.created_at,
            }))
        } else {
            Ok(None)
        }
    }
    
    /// Get call history for conversation
    pub async fn get_call_history(
        &self,
        conversation_id: Uuid,
        limit: i32,
    ) -> Result<Vec<Call>> {
        let rows = sqlx::query!(
            "SELECT * FROM calls WHERE conversation_id = $1 ORDER BY created_at DESC LIMIT $2",
            conversation_id,
            limit as i64
        )
        .fetch_all(&self.db)
        .await?;
        
        Ok(rows
            .into_iter()
            .map(|row| Call {
                id: row.id,
                conversation_id: row.conversation_id,
                caller_id: row.caller_id,
                call_type: serde_json::from_value(row.call_type).unwrap(),
                status: serde_json::from_value(row.status).unwrap(),
                sdp_offer: row.sdp_offer,
                sdp_answer: row.sdp_answer,
                ice_candidates: Vec::new(),
                started_at: row.started_at,
                ended_at: row.ended_at,
                duration_seconds: row.duration_seconds,
                created_at: row.created_at,
            })
            .collect())
    }
    
    // Helper methods
    
    async fn publish_call_event(&self, call: &Call, event_type: &str) -> Result<()> {
        let mut conn = self.redis.get_async_connection().await?;
        
        let channel = format!("calls:{}", call.conversation_id);
        let payload = serde_json::json!({
            "type": event_type,
            "call": call,
        });
        
        let _: () = conn.publish(&channel, payload.to_string()).await?;
        
        Ok(())
    }
    
    async fn publish_ice_candidate(&self, call_id: Uuid, user_id: Uuid, candidate: String) -> Result<()> {
        let mut conn = self.redis.get_async_connection().await?;
        
        let channel = format!("ice:{}", call_id);
        let payload = serde_json::json!({
            "user_id": user_id,
            "candidate": candidate,
        });
        
        let _: () = conn.publish(&channel, payload.to_string()).await?;
        
        Ok(())
    }
}
