use serde::{Deserialize, Serialize};
use uuid::Uuid;
use chrono::{DateTime, Utc};

/// Encrypted message (what gets transmitted/stored)
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Message {
    pub id: Uuid,
    pub conversation_id: Uuid,      // 1:1 chat or group
    pub sender_id: Uuid,
    pub sender_device_id: String,
    pub recipient_id: Option<Uuid>, // For 1:1, None for groups
    pub message_type: MessageType,
    
    // ENCRYPTED CONTENT - Server cannot decrypt!
    pub ciphertext: Vec<u8>,        // Encrypted message body
    pub ephemeral_key: Vec<u8>,     // For key agreement
    
    // METADATA - Server CAN see this
    pub sequence_number: i64,
    pub timestamp: DateTime<Utc>,
    pub is_group: bool,
    pub group_epoch: Option<u64>,   // MLS epoch for groups
    
    // Delivery tracking
    pub status: MessageStatus,
    pub delivered_at: Option<DateTime<Utc>>,
    pub read_at: Option<DateTime<Utc>>,
    
    // Features
    pub is_self_destructing: bool,
    pub expires_at: Option<DateTime<Utc>>,
    pub is_edited: bool,
    pub edited_at: Option<DateTime<Utc>>,
    
    pub created_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub enum MessageType {
    Text,
    Media,              // Photos, videos (encrypted)
    Audio,              // Voice notes
    File,               // Documents, etc.
    Location,
    Contact,
    Poll,
    Event,
    Call,               // Call invitation
    System,             // "Alice added Bob"
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub enum MessageStatus {
    Sending,
    Sent,
    Delivered,
    Read,
    Failed,
}

/// Decrypted message content (client-side only!)
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MessageContent {
    pub text: Option<String>,
    pub media: Option<MediaAttachment>,
    pub location: Option<Location>,
    pub contact: Option<Contact>,
    pub poll: Option<Poll>,
    pub event: Option<Event>,
    pub quoted_message_id: Option<Uuid>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MediaAttachment {
    pub media_type: MediaType,
    pub encrypted_url: String,       // URL to encrypted blob
    pub encryption_key: Vec<u8>,     // Key to decrypt media
    pub file_size: i64,
    pub mime_type: String,
    pub duration_ms: Option<i32>,    // For audio/video
    pub width: Option<i32>,
    pub height: Option<i32>,
    pub thumbnail: Option<Vec<u8>>,  // Encrypted thumbnail
    pub blurhash: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub enum MediaType {
    Photo,
    Video,
    VoiceNote,
    Audio,
    Document,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Location {
    pub latitude: f64,
    pub longitude: f64,
    pub address: Option<String>,
    pub name: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Contact {
    pub name: String,
    pub phone: Option<String>,
    pub user_id: Option<Uuid>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Poll {
    pub question: String,
    pub options: Vec<String>,
    pub votes: Vec<PollVote>,
    pub is_anonymous: bool,
    pub allows_multiple: bool,
    pub expires_at: Option<DateTime<Utc>>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PollVote {
    pub user_id: Uuid,
    pub option_index: i32,
    pub voted_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Event {
    pub title: String,
    pub description: Option<String>,
    pub location: Option<String>,
    pub start_time: DateTime<Utc>,
    pub end_time: Option<DateTime<Utc>>,
    pub attendees: Vec<Uuid>,
}

/// Conversation (1:1 or group chat)
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Conversation {
    pub id: Uuid,
    pub conversation_type: ConversationType,
    pub name: Option<String>,          // For groups
    pub avatar_url: Option<String>,
    pub created_by: Uuid,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
    pub is_archived: bool,
    pub is_muted: bool,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub enum ConversationType {
    OneToOne,
    Group,
    NoteToSelf,
}

/// Group chat info
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct GroupChat {
    pub id: Uuid,
    pub name: String,
    pub description: Option<String>,
    pub avatar_url: Option<String>,
    pub created_by: Uuid,
    pub member_count: i32,
    pub max_members: i32,              // 1,500
    pub is_public: bool,
    pub invite_link: Option<String>,
    pub mls_group_id: Vec<u8>,         // MLS group ID
    pub current_epoch: u64,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

/// Group member
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct GroupMember {
    pub group_id: Uuid,
    pub user_id: Uuid,
    pub role: GroupMemberRole,
    pub joined_at: DateTime<Utc>,
    pub added_by: Option<Uuid>,
    pub is_muted: bool,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub enum GroupMemberRole {
    Owner,
    Admin,
    Member,
}

/// Read receipt
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ReadReceipt {
    pub message_id: Uuid,
    pub user_id: Uuid,
    pub read_at: DateTime<Utc>,
}

/// Typing indicator
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TypingIndicator {
    pub conversation_id: Uuid,
    pub user_id: Uuid,
    pub is_typing: bool,
    pub timestamp: DateTime<Utc>,
}

/// Presence status
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Presence {
    pub user_id: Uuid,
    pub status: PresenceStatus,
    pub last_seen: DateTime<Utc>,
    pub custom_status: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub enum PresenceStatus {
    Online,
    Away,
    Busy,
    Offline,
}

/// Call (audio/video)
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Call {
    pub id: Uuid,
    pub conversation_id: Uuid,
    pub caller_id: Uuid,
    pub call_type: CallType,
    pub status: CallStatus,
    pub sdp_offer: Option<String>,     // WebRTC SDP offer (encrypted)
    pub sdp_answer: Option<String>,    // WebRTC SDP answer (encrypted)
    pub ice_candidates: Vec<String>,   // ICE candidates
    pub started_at: Option<DateTime<Utc>>,
    pub ended_at: Option<DateTime<Utc>>,
    pub duration_seconds: Option<i32>,
    pub created_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub enum CallType {
    Audio,
    Video,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub enum CallStatus {
    Ringing,
    Answered,
    Declined,
    Missed,
    Ended,
    Failed,
}

/// Message send request (from client)
#[derive(Debug, Deserialize)]
pub struct SendMessageRequest {
    pub conversation_id: Uuid,
    pub recipient_id: Option<Uuid>,   // For 1:1
    pub device_id: String,
    pub ciphertext: String,           // Base64 encoded
    pub ephemeral_key: String,        // Base64 encoded
    pub message_type: MessageType,
    pub is_self_destructing: bool,
    pub expires_in_seconds: Option<i64>,
}

/// Message response (to client)
#[derive(Debug, Serialize)]
pub struct MessageResponse {
    pub message_id: Uuid,
    pub conversation_id: Uuid,
    pub sequence_number: i64,
    pub timestamp: DateTime<Utc>,
    pub status: MessageStatus,
}

/// Get messages request
#[derive(Debug, Deserialize)]
pub struct GetMessagesRequest {
    pub conversation_id: Uuid,
    pub before_sequence: Option<i64>,
    pub limit: i32,
}

impl Message {
    pub fn is_expired(&self) -> bool {
        if let Some(expires_at) = self.expires_at {
            return Utc::now() > expires_at;
        }
        false
    }
    
    pub fn mark_delivered(&mut self) {
        self.status = MessageStatus::Delivered;
        self.delivered_at = Some(Utc::now());
    }
    
    pub fn mark_read(&mut self) {
        self.status = MessageStatus::Read;
        self.read_at = Some(Utc::now());
    }
}

impl GroupChat {
    pub fn can_add_member(&self) -> bool {
        self.member_count < self.max_members
    }
    
    pub fn is_full(&self) -> bool {
        self.member_count >= self.max_members
    }
}
