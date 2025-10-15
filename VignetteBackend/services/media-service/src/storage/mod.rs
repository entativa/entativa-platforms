pub mod s3_client;
pub mod minio_client;
pub mod local_storage;

use async_trait::async_trait;
use bytes::Bytes;
use std::path::Path;
use thiserror::Error;

pub use s3_client::S3Client;
pub use local_storage::LocalStorage;

#[derive(Error, Debug)]
pub enum StorageError {
    #[error("File not found: {0}")]
    NotFound(String),
    
    #[error("Permission denied: {0}")]
    PermissionDenied(String),
    
    #[error("Storage error: {0}")]
    StorageError(String),
    
    #[error("IO error: {0}")]
    IoError(#[from] std::io::Error),
    
    #[error("S3 error: {0}")]
    S3Error(String),
}

pub type StorageResult<T> = Result<T, StorageError>;

#[async_trait]
pub trait StorageBackend: Send + Sync {
    async fn upload(&self, path: &str, data: Bytes, content_type: &str) -> StorageResult<String>;
    
    async fn download(&self, path: &str) -> StorageResult<Bytes>;
    
    async fn delete(&self, path: &str) -> StorageResult<()>;
    
    async fn exists(&self, path: &str) -> StorageResult<bool>;
    
    async fn get_url(&self, path: &str) -> StorageResult<String>;
    
    async fn get_signed_url(&self, path: &str, expires_in_secs: u64) -> StorageResult<String>;
    
    async fn list(&self, prefix: &str) -> StorageResult<Vec<String>>;
    
    async fn copy(&self, from: &str, to: &str) -> StorageResult<()>;
    
    async fn get_metadata(&self, path: &str) -> StorageResult<FileMetadata>;
}

#[derive(Debug, Clone)]
pub struct FileMetadata {
    pub size: u64,
    pub content_type: String,
    pub last_modified: Option<chrono::DateTime<chrono::Utc>>,
    pub etag: Option<String>,
}
