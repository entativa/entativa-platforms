use serde::{Deserialize, Serialize};
use uuid::Uuid;
use chrono::{DateTime, Utc};

/// Identity key pair (long-term Ed25519)
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct IdentityKeyPair {
    pub user_id: Uuid,
    pub public_key: Vec<u8>,  // Ed25519 public key (32 bytes)
    #[serde(skip_serializing)]
    pub private_key: Vec<u8>, // Ed25519 private key (64 bytes) - NEVER send over network!
    pub created_at: DateTime<Utc>,
}

/// Signed pre-key (medium-term X25519)
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SignedPreKey {
    pub id: i32,
    pub user_id: Uuid,
    pub public_key: Vec<u8>,  // X25519 public key (32 bytes)
    #[serde(skip_serializing)]
    pub private_key: Vec<u8>, // X25519 private key (32 bytes)
    pub signature: Vec<u8>,   // Signed by identity key
    pub created_at: DateTime<Utc>,
}

/// One-time pre-key (single-use X25519)
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct OneTimePreKey {
    pub id: i32,
    pub user_id: Uuid,
    pub public_key: Vec<u8>,  // X25519 public key (32 bytes)
    #[serde(skip_serializing)]
    pub private_key: Vec<u8>, // X25519 private key (32 bytes)
    pub is_used: bool,
    pub created_at: DateTime<Utc>,
}

/// Pre-key bundle (what client fetches to start conversation)
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PreKeyBundle {
    pub user_id: Uuid,
    pub device_id: String,
    pub registration_id: i32,
    pub identity_key: Vec<u8>,        // Identity public key
    pub signed_prekey_id: i32,
    pub signed_prekey: Vec<u8>,       // Signed pre-key public
    pub signed_prekey_signature: Vec<u8>,
    pub onetime_prekey_id: Option<i32>,
    pub onetime_prekey: Option<Vec<u8>>, // One-time pre-key public (if available)
}

/// Device registration info
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Device {
    pub id: Uuid,
    pub user_id: Uuid,
    pub device_id: String,           // Unique device identifier
    pub device_name: String,         // User-friendly name
    pub registration_id: i32,        // Random 14-bit number
    pub identity_key: Vec<u8>,       // Device's identity key
    pub is_active: bool,
    pub last_seen: DateTime<Utc>,
    pub created_at: DateTime<Utc>,
}

/// Session state (Double Ratchet)
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SessionState {
    pub id: Uuid,
    pub user_id: Uuid,              // Local user
    pub peer_user_id: Uuid,         // Remote user
    pub peer_device_id: String,     // Remote device
    pub root_key: Vec<u8>,          // 32 bytes (encrypted in DB)
    pub chain_key: Vec<u8>,         // 32 bytes (encrypted in DB)
    pub receive_chain_key: Vec<u8>, // 32 bytes (encrypted in DB)
    pub send_ratchet_public: Vec<u8>,
    pub send_ratchet_private: Vec<u8>,
    pub receive_ratchet_public: Vec<u8>,
    pub message_number: i32,
    pub previous_chain_length: i32,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

/// MLS Group key material (for group chats)
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MLSGroupState {
    pub group_id: Uuid,
    pub epoch: u64,                    // MLS epoch number
    pub tree_hash: Vec<u8>,            // Hash of ratchet tree
    pub group_context: Vec<u8>,        // Serialized MLS group context
    pub encryption_key: Vec<u8>,       // Current epoch encryption key (encrypted)
    pub sender_data_key: Vec<u8>,      // Sender data key (encrypted)
    pub member_count: i32,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

/// MLS Group member
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MLSGroupMember {
    pub group_id: Uuid,
    pub user_id: Uuid,
    pub device_id: String,
    pub leaf_index: u32,               // Position in ratchet tree
    pub credential: Vec<u8>,           // MLS credential
    pub key_package: Vec<u8>,          // MLS key package
    pub role: GroupRole,
    pub joined_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub enum GroupRole {
    Owner,
    Admin,
    Member,
}

/// Key registration request (from client)
#[derive(Debug, Deserialize)]
pub struct KeyRegistrationRequest {
    pub device_id: String,
    pub device_name: String,
    pub registration_id: i32,
    pub identity_key: String,           // Base64 encoded
    pub signed_prekey: SignedPreKeyUpload,
    pub onetime_prekeys: Vec<OneTimePreKeyUpload>,
}

#[derive(Debug, Deserialize)]
pub struct SignedPreKeyUpload {
    pub id: i32,
    pub public_key: String,             // Base64 encoded
    pub signature: String,              // Base64 encoded
}

#[derive(Debug, Deserialize)]
pub struct OneTimePreKeyUpload {
    pub id: i32,
    pub public_key: String,             // Base64 encoded
}

/// Pre-key bundle request
#[derive(Debug, Deserialize)]
pub struct PreKeyBundleRequest {
    pub user_id: Uuid,
    pub device_id: Option<String>,      // Specific device, or any device
}

/// Key rotation request
#[derive(Debug, Deserialize)]
pub struct KeyRotationRequest {
    pub device_id: String,
    pub new_signed_prekey: SignedPreKeyUpload,
    pub new_onetime_prekeys: Vec<OneTimePreKeyUpload>,
}

impl KeyRegistrationRequest {
    /// Validate registration request
    pub fn validate(&self) -> Result<(), String> {
        // Registration ID should be 14-bit (0-16383)
        if self.registration_id < 0 || self.registration_id > 16383 {
            return Err("Invalid registration_id (must be 0-16383)".to_string());
        }
        
        // Identity key should be 32 bytes (Ed25519)
        if let Ok(decoded) = base64::decode(&self.identity_key) {
            if decoded.len() != 32 {
                return Err("Identity key must be 32 bytes".to_string());
            }
        } else {
            return Err("Invalid base64 for identity_key".to_string());
        }
        
        // Signed prekey validation
        if let Ok(decoded) = base64::decode(&self.signed_prekey.public_key) {
            if decoded.len() != 32 {
                return Err("Signed prekey must be 32 bytes".to_string());
            }
        } else {
            return Err("Invalid base64 for signed_prekey".to_string());
        }
        
        // Should have at least 50 one-time prekeys
        if self.onetime_prekeys.len() < 50 {
            return Err("Must upload at least 50 one-time prekeys".to_string());
        }
        
        Ok(())
    }
}
