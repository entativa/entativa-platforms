use blake3;
use bytes::Bytes;
use hex;
use sha2::{Digest, Sha256};

pub fn compute_sha256(data: &[u8]) -> String {
    let mut hasher = Sha256::new();
    hasher.update(data);
    hex::encode(hasher.finalize())
}

pub fn compute_blake3(data: &[u8]) -> String {
    let hash = blake3::hash(data);
    hash.to_hex().to_string()
}

pub fn compute_checksum(data: &Bytes) -> String {
    compute_blake3(data)
}

pub fn verify_checksum(data: &Bytes, expected: &str) -> bool {
    let computed = compute_checksum(data);
    computed == expected
}

pub fn generate_file_hash(data: &[u8]) -> String {
    compute_blake3(data)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_sha256() {
        let data = b"hello world";
        let hash = compute_sha256(data);
        assert_eq!(hash.len(), 64); // SHA256 produces 64 hex characters
    }

    #[test]
    fn test_blake3() {
        let data = b"hello world";
        let hash = compute_blake3(data);
        assert_eq!(hash.len(), 64); // BLAKE3 produces 64 hex characters
    }

    #[test]
    fn test_verify_checksum() {
        let data = Bytes::from("test data");
        let checksum = compute_checksum(&data);
        assert!(verify_checksum(&data, &checksum));
        assert!(!verify_checksum(&data, "invalid"));
    }
}
