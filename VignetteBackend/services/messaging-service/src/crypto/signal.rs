use anyhow::{anyhow, Result};
use ed25519_dalek::{Keypair as Ed25519KeyPair, PublicKey as Ed25519PublicKey, SecretKey as Ed25519SecretKey, Signature, Signer, Verifier};
use x25519_dalek::{PublicKey as X25519PublicKey, StaticSecret as X25519PrivateKey};
use hkdf::Hkdf;
use sha2::Sha256;
use aes_gcm::{
    aead::{Aead, KeyInit, OsRng},
    Aes256Gcm, Nonce,
};
use rand::RngCore;

/// Signal Protocol implementation (simplified Double Ratchet)
pub struct SignalProtocol;

impl SignalProtocol {
    /// Generate identity key pair (Ed25519 for signing)
    pub fn generate_identity_keypair() -> Result<(Vec<u8>, Vec<u8>)> {
        let mut csprng = OsRng;
        let keypair = Ed25519KeyPair::generate(&mut csprng);
        
        Ok((
            keypair.public.as_bytes().to_vec(),
            keypair.secret.as_bytes().to_vec(),
        ))
    }
    
    /// Generate pre-key pair (X25519 for ECDH)
    pub fn generate_prekey_keypair() -> Result<(Vec<u8>, Vec<u8>)> {
        let mut csprng = OsRng;
        let private = X25519PrivateKey::random_from_rng(&mut csprng);
        let public = X25519PublicKey::from(&private);
        
        Ok((
            public.as_bytes().to_vec(),
            private.to_bytes().to_vec(),
        ))
    }
    
    /// Sign pre-key with identity key
    pub fn sign_prekey(
        prekey_public: &[u8],
        identity_private: &[u8],
    ) -> Result<Vec<u8>> {
        if identity_private.len() != 32 {
            return Err(anyhow!("Invalid identity private key length"));
        }
        
        let secret = Ed25519SecretKey::from_bytes(identity_private)?;
        let public = Ed25519PublicKey::from(&secret);
        let keypair = Ed25519KeyPair { secret, public };
        
        let signature = keypair.sign(prekey_public);
        Ok(signature.to_bytes().to_vec())
    }
    
    /// Verify pre-key signature
    pub fn verify_prekey_signature(
        prekey_public: &[u8],
        signature: &[u8],
        identity_public: &[u8],
    ) -> Result<bool> {
        if signature.len() != 64 {
            return Ok(false);
        }
        if identity_public.len() != 32 {
            return Ok(false);
        }
        
        let public_key = Ed25519PublicKey::from_bytes(identity_public)?;
        let sig = Signature::from_bytes(signature)?;
        
        Ok(public_key.verify(prekey_public, &sig).is_ok())
    }
    
    /// Perform X25519 ECDH
    pub fn ecdh(private_key: &[u8], public_key: &[u8]) -> Result<Vec<u8>> {
        if private_key.len() != 32 || public_key.len() != 32 {
            return Err(anyhow!("Invalid key length for ECDH"));
        }
        
        let private_key_bytes: [u8; 32] = private_key.try_into()?;
        let public_key_bytes: [u8; 32] = public_key.try_into()?;
        
        let private = X25519PrivateKey::from(private_key_bytes);
        let public = X25519PublicKey::from(public_key_bytes);
        
        let shared = private.diffie_hellman(&public);
        Ok(shared.as_bytes().to_vec())
    }
    
    /// Derive root key and chain key from shared secrets
    /// This is the X3DH key agreement
    pub fn derive_keys(
        dh1: &[u8],
        dh2: &[u8],
        dh3: &[u8],
        dh4: Option<&[u8]>,
    ) -> Result<(Vec<u8>, Vec<u8>)> {
        // Concatenate DH outputs
        let mut ikm = Vec::new();
        ikm.extend_from_slice(dh1);
        ikm.extend_from_slice(dh2);
        ikm.extend_from_slice(dh3);
        if let Some(dh4) = dh4 {
            ikm.extend_from_slice(dh4);
        }
        
        // HKDF to derive root key and chain key
        let hkdf = Hkdf::<Sha256>::new(None, &ikm);
        let mut okm = [0u8; 64];
        hkdf.expand(b"SignalRootAndChainKeys", &mut okm)?;
        
        let root_key = okm[0..32].to_vec();
        let chain_key = okm[32..64].to_vec();
        
        Ok((root_key, chain_key))
    }
    
    /// KDF chain step (ratchet forward)
    pub fn kdf_chain(chain_key: &[u8]) -> Result<(Vec<u8>, Vec<u8>)> {
        let hkdf = Hkdf::<Sha256>::new(None, chain_key);
        let mut okm = [0u8; 64];
        hkdf.expand(b"SignalChainKey", &mut okm)?;
        
        let message_key = okm[0..32].to_vec();
        let next_chain_key = okm[32..64].to_vec();
        
        Ok((message_key, next_chain_key))
    }
    
    /// Encrypt message with message key
    pub fn encrypt_message(
        plaintext: &[u8],
        message_key: &[u8],
        associated_data: &[u8],
    ) -> Result<Vec<u8>> {
        if message_key.len() != 32 {
            return Err(anyhow!("Message key must be 32 bytes"));
        }
        
        // Derive encryption key and auth key
        let hkdf = Hkdf::<Sha256>::new(None, message_key);
        let mut okm = [0u8; 80];
        hkdf.expand(b"SignalMessageKeys", &mut okm)?;
        
        let encryption_key = &okm[0..32];
        let auth_key = &okm[32..64];
        let iv = &okm[64..80];
        
        // AES-GCM encryption
        let cipher = Aes256Gcm::new_from_slice(encryption_key)?;
        let nonce = Nonce::from_slice(&iv[0..12]);
        
        let mut ciphertext = cipher.encrypt(nonce, plaintext)?;
        
        // Compute MAC over ciphertext + associated data
        let mac = Self::compute_mac(auth_key, &ciphertext, associated_data)?;
        
        // Append MAC
        ciphertext.extend_from_slice(&mac);
        
        Ok(ciphertext)
    }
    
    /// Decrypt message with message key
    pub fn decrypt_message(
        ciphertext_with_mac: &[u8],
        message_key: &[u8],
        associated_data: &[u8],
    ) -> Result<Vec<u8>> {
        if message_key.len() != 32 {
            return Err(anyhow!("Message key must be 32 bytes"));
        }
        if ciphertext_with_mac.len() < 32 {
            return Err(anyhow!("Ciphertext too short"));
        }
        
        // Split ciphertext and MAC
        let split_point = ciphertext_with_mac.len() - 32;
        let ciphertext = &ciphertext_with_mac[0..split_point];
        let mac = &ciphertext_with_mac[split_point..];
        
        // Derive keys
        let hkdf = Hkdf::<Sha256>::new(None, message_key);
        let mut okm = [0u8; 80];
        hkdf.expand(b"SignalMessageKeys", &mut okm)?;
        
        let encryption_key = &okm[0..32];
        let auth_key = &okm[32..64];
        let iv = &okm[64..80];
        
        // Verify MAC
        let expected_mac = Self::compute_mac(auth_key, ciphertext, associated_data)?;
        if mac != expected_mac {
            return Err(anyhow!("MAC verification failed"));
        }
        
        // AES-GCM decryption
        let cipher = Aes256Gcm::new_from_slice(encryption_key)?;
        let nonce = Nonce::from_slice(&iv[0..12]);
        
        let plaintext = cipher.decrypt(nonce, ciphertext)?;
        
        Ok(plaintext)
    }
    
    /// Compute HMAC-SHA256
    fn compute_mac(key: &[u8], ciphertext: &[u8], associated_data: &[u8]) -> Result<Vec<u8>> {
        use sha2::{Digest, Sha256};
        
        let mut hasher = Sha256::new();
        hasher.update(key);
        hasher.update(ciphertext);
        hasher.update(associated_data);
        
        Ok(hasher.finalize().to_vec())
    }
    
    /// Generate random bytes
    pub fn random_bytes(len: usize) -> Vec<u8> {
        let mut bytes = vec![0u8; len];
        OsRng.fill_bytes(&mut bytes);
        bytes
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_keypair_generation() {
        let (public, private) = SignalProtocol::generate_identity_keypair().unwrap();
        assert_eq!(public.len(), 32);
        assert_eq!(private.len(), 32);
        
        let (public, private) = SignalProtocol::generate_prekey_keypair().unwrap();
        assert_eq!(public.len(), 32);
        assert_eq!(private.len(), 32);
    }
    
    #[test]
    fn test_prekey_signing() {
        let (identity_public, identity_private) = SignalProtocol::generate_identity_keypair().unwrap();
        let (prekey_public, _) = SignalProtocol::generate_prekey_keypair().unwrap();
        
        let signature = SignalProtocol::sign_prekey(&prekey_public, &identity_private).unwrap();
        assert_eq!(signature.len(), 64);
        
        let valid = SignalProtocol::verify_prekey_signature(&prekey_public, &signature, &identity_public).unwrap();
        assert!(valid);
    }
    
    #[test]
    fn test_ecdh() {
        let (public1, private1) = SignalProtocol::generate_prekey_keypair().unwrap();
        let (public2, private2) = SignalProtocol::generate_prekey_keypair().unwrap();
        
        let shared1 = SignalProtocol::ecdh(&private1, &public2).unwrap();
        let shared2 = SignalProtocol::ecdh(&private2, &public1).unwrap();
        
        assert_eq!(shared1, shared2);
        assert_eq!(shared1.len(), 32);
    }
    
    #[test]
    fn test_encrypt_decrypt() {
        let message_key = SignalProtocol::random_bytes(32);
        let plaintext = b"Hello, Signal!";
        let associated_data = b"metadata";
        
        let ciphertext = SignalProtocol::encrypt_message(plaintext, &message_key, associated_data).unwrap();
        let decrypted = SignalProtocol::decrypt_message(&ciphertext, &message_key, associated_data).unwrap();
        
        assert_eq!(&decrypted, plaintext);
    }
}
