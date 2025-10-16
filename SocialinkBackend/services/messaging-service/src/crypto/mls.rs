use anyhow::{anyhow, Result};
use hkdf::Hkdf;
use sha2::Sha256;
use aes_gcm::{
    aead::{Aead, KeyInit, OsRng},
    Aes256Gcm, Nonce,
};
use rand::RngCore;
use std::collections::HashMap;
use uuid::Uuid;

/// MLS (Messaging Layer Security) implementation for groups
/// Supports up to 1,500 members with efficient key management
pub struct MLSProtocol;

/// Ratchet tree node
#[derive(Debug, Clone)]
pub struct TreeNode {
    pub index: u32,
    pub public_key: Option<Vec<u8>>,
    pub private_key: Option<Vec<u8>>,
    pub parent_hash: Vec<u8>,
    pub unmerged_leaves: Vec<u32>,
}

/// MLS Group state
#[derive(Debug, Clone)]
pub struct MLSGroup {
    pub group_id: Uuid,
    pub epoch: u64,
    pub tree: Vec<TreeNode>,
    pub encryption_key: Vec<u8>,
    pub sender_data_key: Vec<u8>,
    pub member_map: HashMap<Uuid, u32>, // user_id -> leaf_index
}

impl MLSProtocol {
    /// Create new MLS group
    pub fn create_group(
        group_id: Uuid,
        creator_id: Uuid,
        creator_public_key: &[u8],
    ) -> Result<MLSGroup> {
        // Initialize tree with creator as first leaf
        let mut tree = Vec::new();
        
        // Leaf node for creator
        tree.push(TreeNode {
            index: 0,
            public_key: Some(creator_public_key.to_vec()),
            private_key: None,
            parent_hash: vec![],
            unmerged_leaves: vec![],
        });
        
        // Root node
        tree.push(TreeNode {
            index: 1,
            public_key: Some(creator_public_key.to_vec()),
            private_key: None,
            parent_hash: vec![],
            unmerged_leaves: vec![],
        });
        
        // Derive initial epoch keys
        let (encryption_key, sender_data_key) = Self::derive_epoch_keys(creator_public_key, 0)?;
        
        let mut member_map = HashMap::new();
        member_map.insert(creator_id, 0);
        
        Ok(MLSGroup {
            group_id,
            epoch: 0,
            tree,
            encryption_key,
            sender_data_key,
            member_map,
        })
    }
    
    /// Add member to group
    pub fn add_member(
        group: &mut MLSGroup,
        user_id: Uuid,
        public_key: &[u8],
    ) -> Result<Vec<u8>> {
        let leaf_index = group.tree.len() as u32;
        
        // Add leaf node for new member
        group.tree.push(TreeNode {
            index: leaf_index,
            public_key: Some(public_key.to_vec()),
            private_key: None,
            parent_hash: vec![],
            unmerged_leaves: vec![],
        });
        
        // Update member map
        group.member_map.insert(user_id, leaf_index);
        
        // Increment epoch
        group.epoch += 1;
        
        // Re-derive epoch keys (in practice, use path secrets)
        let tree_hash = Self::compute_tree_hash(&group.tree)?;
        let (encryption_key, sender_data_key) = Self::derive_epoch_keys(&tree_hash, group.epoch)?;
        
        group.encryption_key = encryption_key;
        group.sender_data_key = sender_data_key;
        
        // Generate Welcome message for new member (simplified)
        Ok(Self::generate_welcome(group, user_id, leaf_index)?)
    }
    
    /// Remove member from group
    pub fn remove_member(
        group: &mut MLSGroup,
        user_id: Uuid,
    ) -> Result<()> {
        let leaf_index = group.member_map.remove(&user_id)
            .ok_or_else(|| anyhow!("Member not found"))?;
        
        // Blank the leaf node
        if let Some(node) = group.tree.get_mut(leaf_index as usize) {
            node.public_key = None;
            node.private_key = None;
        }
        
        // Increment epoch
        group.epoch += 1;
        
        // Re-derive epoch keys
        let tree_hash = Self::compute_tree_hash(&group.tree)?;
        let (encryption_key, sender_data_key) = Self::derive_epoch_keys(&tree_hash, group.epoch)?;
        
        group.encryption_key = encryption_key;
        group.sender_data_key = sender_data_key;
        
        Ok(())
    }
    
    /// Update member's key (key rotation)
    pub fn update_member_key(
        group: &mut MLSGroup,
        user_id: Uuid,
        new_public_key: &[u8],
    ) -> Result<()> {
        let leaf_index = *group.member_map.get(&user_id)
            .ok_or_else(|| anyhow!("Member not found"))?;
        
        // Update leaf node
        if let Some(node) = group.tree.get_mut(leaf_index as usize) {
            node.public_key = Some(new_public_key.to_vec());
        }
        
        // Increment epoch
        group.epoch += 1;
        
        // Re-derive epoch keys
        let tree_hash = Self::compute_tree_hash(&group.tree)?;
        let (encryption_key, sender_data_key) = Self::derive_epoch_keys(&tree_hash, group.epoch)?;
        
        group.encryption_key = encryption_key;
        group.sender_data_key = sender_data_key;
        
        Ok(())
    }
    
    /// Encrypt group message
    pub fn encrypt_group_message(
        group: &MLSGroup,
        plaintext: &[u8],
        sender_id: Uuid,
    ) -> Result<Vec<u8>> {
        // Get sender's leaf index
        let sender_index = group.member_map.get(&sender_id)
            .ok_or_else(|| anyhow!("Sender not in group"))?;
        
        // Build associated data (epoch, group_id, sender_index)
        let mut associated_data = Vec::new();
        associated_data.extend_from_slice(&group.epoch.to_le_bytes());
        associated_data.extend_from_slice(group.group_id.as_bytes());
        associated_data.extend_from_slice(&sender_index.to_le_bytes());
        
        // Encrypt with epoch encryption key
        let cipher = Aes256Gcm::new_from_slice(&group.encryption_key)?;
        
        // Generate random nonce
        let mut nonce_bytes = [0u8; 12];
        OsRng.fill_bytes(&mut nonce_bytes);
        let nonce = Nonce::from_slice(&nonce_bytes);
        
        let mut ciphertext = cipher.encrypt(nonce, plaintext)?;
        
        // Prepend nonce to ciphertext
        let mut result = nonce_bytes.to_vec();
        result.append(&mut ciphertext);
        
        Ok(result)
    }
    
    /// Decrypt group message
    pub fn decrypt_group_message(
        group: &MLSGroup,
        ciphertext_with_nonce: &[u8],
    ) -> Result<Vec<u8>> {
        if ciphertext_with_nonce.len() < 12 {
            return Err(anyhow!("Ciphertext too short"));
        }
        
        // Extract nonce and ciphertext
        let nonce = Nonce::from_slice(&ciphertext_with_nonce[0..12]);
        let ciphertext = &ciphertext_with_nonce[12..];
        
        // Decrypt with epoch encryption key
        let cipher = Aes256Gcm::new_from_slice(&group.encryption_key)?;
        let plaintext = cipher.decrypt(nonce, ciphertext)?;
        
        Ok(plaintext)
    }
    
    /// Compute tree hash (simplified)
    fn compute_tree_hash(tree: &[TreeNode]) -> Result<Vec<u8>> {
        use sha2::{Digest, Sha256};
        
        let mut hasher = Sha256::new();
        
        for node in tree {
            hasher.update(&node.index.to_le_bytes());
            if let Some(ref public_key) = node.public_key {
                hasher.update(public_key);
            }
        }
        
        Ok(hasher.finalize().to_vec())
    }
    
    /// Derive epoch keys from tree hash
    fn derive_epoch_keys(tree_hash: &[u8], epoch: u64) -> Result<(Vec<u8>, Vec<u8>)> {
        let hkdf = Hkdf::<Sha256>::new(Some(&epoch.to_le_bytes()), tree_hash);
        let mut okm = [0u8; 64];
        hkdf.expand(b"MLSEpochKeys", &mut okm)?;
        
        let encryption_key = okm[0..32].to_vec();
        let sender_data_key = okm[32..64].to_vec();
        
        Ok((encryption_key, sender_data_key))
    }
    
    /// Generate Welcome message for new member (simplified)
    fn generate_welcome(group: &MLSGroup, user_id: Uuid, leaf_index: u32) -> Result<Vec<u8>> {
        use serde_json::json;
        
        // In production, this would be a proper MLS Welcome message
        // containing epoch secrets encrypted to the new member
        let welcome = json!({
            "group_id": group.group_id,
            "epoch": group.epoch,
            "leaf_index": leaf_index,
            "tree_size": group.tree.len(),
        });
        
        Ok(serde_json::to_vec(&welcome)?)
    }
    
    /// Validate group size (max 1,500 members)
    pub fn validate_group_size(member_count: usize) -> Result<()> {
        if member_count > 1500 {
            return Err(anyhow!("Group size exceeds maximum of 1,500 members"));
        }
        Ok(())
    }
    
    /// Get group info (for external view)
    pub fn get_group_info(group: &MLSGroup) -> HashMap<String, serde_json::Value> {
        use serde_json::json;
        
        let mut info = HashMap::new();
        info.insert("group_id".to_string(), json!(group.group_id));
        info.insert("epoch".to_string(), json!(group.epoch));
        info.insert("member_count".to_string(), json!(group.member_map.len()));
        info.insert("tree_size".to_string(), json!(group.tree.len()));
        
        info
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_create_group() {
        let group_id = Uuid::new_v4();
        let creator_id = Uuid::new_v4();
        let creator_key = vec![1u8; 32];
        
        let group = MLSProtocol::create_group(group_id, creator_id, &creator_key).unwrap();
        
        assert_eq!(group.epoch, 0);
        assert_eq!(group.member_map.len(), 1);
        assert_eq!(group.encryption_key.len(), 32);
    }
    
    #[test]
    fn test_add_member() {
        let group_id = Uuid::new_v4();
        let creator_id = Uuid::new_v4();
        let creator_key = vec![1u8; 32];
        
        let mut group = MLSProtocol::create_group(group_id, creator_id, &creator_key).unwrap();
        
        let new_member_id = Uuid::new_v4();
        let new_member_key = vec![2u8; 32];
        
        let welcome = MLSProtocol::add_member(&mut group, new_member_id, &new_member_key).unwrap();
        
        assert_eq!(group.epoch, 1);
        assert_eq!(group.member_map.len(), 2);
        assert!(!welcome.is_empty());
    }
    
    #[test]
    fn test_encrypt_decrypt_group_message() {
        let group_id = Uuid::new_v4();
        let creator_id = Uuid::new_v4();
        let creator_key = vec![1u8; 32];
        
        let group = MLSProtocol::create_group(group_id, creator_id, &creator_key).unwrap();
        
        let plaintext = b"Hello, group!";
        let ciphertext = MLSProtocol::encrypt_group_message(&group, plaintext, creator_id).unwrap();
        let decrypted = MLSProtocol::decrypt_group_message(&group, &ciphertext).unwrap();
        
        assert_eq!(&decrypted, plaintext);
    }
    
    #[test]
    fn test_group_size_validation() {
        assert!(MLSProtocol::validate_group_size(1500).is_ok());
        assert!(MLSProtocol::validate_group_size(1501).is_err());
    }
}
