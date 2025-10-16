use serde::{Deserialize, Serialize};
use uuid::Uuid;
use validator::Validate;
use chrono::{DateTime, Utc};

#[derive(Debug, Clone, Serialize, Deserialize, Validate)]
pub struct MultipartUploadInit {
    #[validate(length(min = 1, max = 255))]
    pub filename: String,
    #[validate(range(min = 1))]
    pub file_size: i64,
    pub mime_type: String,
    pub chunk_size: Option<i64>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MultipartUploadSession {
    pub upload_id: Uuid,
    pub media_id: Uuid,
    pub filename: String,
    pub total_size: i64,
    pub chunk_size: i64,
    pub total_chunks: i32,
    pub uploaded_chunks: Vec<i32>,
    pub status: UploadStatus,
    pub created_at: DateTime<Utc>,
    pub expires_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "lowercase")]
pub enum UploadStatus {
    Pending,
    InProgress,
    Completed,
    Failed,
    Cancelled,
    Expired,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ChunkUpload {
    pub upload_id: Uuid,
    pub chunk_number: i32,
    pub data: Vec<u8>,
    pub checksum: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct UploadProgress {
    pub upload_id: Uuid,
    pub uploaded_chunks: i32,
    pub total_chunks: i32,
    pub bytes_uploaded: i64,
    pub total_bytes: i64,
    pub percentage: f64,
    pub status: UploadStatus,
}

impl UploadProgress {
    pub fn new(
        upload_id: Uuid,
        uploaded_chunks: i32,
        total_chunks: i32,
        bytes_uploaded: i64,
        total_bytes: i64,
    ) -> Self {
        let percentage = if total_bytes > 0 {
            (bytes_uploaded as f64 / total_bytes as f64) * 100.0
        } else {
            0.0
        };

        Self {
            upload_id,
            uploaded_chunks,
            total_chunks,
            bytes_uploaded,
            total_bytes,
            percentage,
            status: UploadStatus::InProgress,
        }
    }
}
