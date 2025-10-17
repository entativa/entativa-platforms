use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use sqlx::FromRow;
use uuid::Uuid;
use validator::Validate;

#[derive(Debug, Clone, Serialize, Deserialize, FromRow)]
pub struct Media {
    pub id: Uuid,
    pub user_id: Uuid,
    pub filename: String,
    pub original_filename: String,
    pub mime_type: String,
    pub media_type: MediaType,
    pub file_size: i64,
    pub storage_path: String,
    pub storage_provider: String,
    pub url: String,
    pub cdn_url: Option<String>,
    pub thumbnail_url: Option<String>,
    pub width: Option<i32>,
    pub height: Option<i32>,
    pub duration: Option<f64>,
    pub aspect_ratio: Option<f64>,
    pub orientation: Option<i16>,
    pub hash: String,
    pub blurhash: Option<String>,
    pub processing_status: ProcessingStatus,
    pub variants: sqlx::types::Json<Vec<MediaVariant>>,
    pub metadata: sqlx::types::Json<MediaMetadata>,
    pub is_processed: bool,
    pub is_public: bool,
    pub is_deleted: bool,
    pub view_count: i64,
    pub download_count: i64,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
    pub deleted_at: Option<DateTime<Utc>>,
}

#[derive(Debug, Clone, Serialize, Deserialize, sqlx::Type, PartialEq)]
#[sqlx(type_name = "media_type", rename_all = "lowercase")]
pub enum MediaType {
    Image,
    Video,
    Audio,
    Document,
}

#[derive(Debug, Clone, Serialize, Deserialize, sqlx::Type, PartialEq)]
#[sqlx(type_name = "processing_status", rename_all = "lowercase")]
pub enum ProcessingStatus {
    Pending,
    Processing,
    Completed,
    Failed,
    Cancelled,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MediaVariant {
    pub variant_type: String, // thumbnail, preview, hd, sd, etc.
    pub url: String,
    pub width: Option<u32>,
    pub height: Option<u32>,
    pub file_size: Option<i64>,
    pub format: String,
    pub bitrate: Option<u32>,
}

#[derive(Debug, Clone, Serialize, Deserialize, Default)]
pub struct MediaMetadata {
    // Image metadata
    pub exif: Option<ExifData>,
    pub color_profile: Option<String>,
    pub dominant_colors: Option<Vec<String>>,
    
    // Video metadata
    pub codec: Option<String>,
    pub bitrate: Option<u32>,
    pub frame_rate: Option<f64>,
    pub audio_codec: Option<String>,
    pub audio_bitrate: Option<u32>,
    pub audio_sample_rate: Option<u32>,
    
    // Audio metadata
    pub artist: Option<String>,
    pub album: Option<String>,
    pub title: Option<String>,
    pub genre: Option<String>,
    pub year: Option<i32>,
    
    // Common metadata
    pub camera: Option<String>,
    pub lens: Option<String>,
    pub iso: Option<u32>,
    pub aperture: Option<f64>,
    pub shutter_speed: Option<String>,
    pub focal_length: Option<f64>,
    pub gps_latitude: Option<f64>,
    pub gps_longitude: Option<f64>,
    pub gps_altitude: Option<f64>,
    
    // Processing metadata
    pub filters_applied: Option<Vec<String>>,
    pub edited: bool,
    pub edit_history: Option<Vec<EditOperation>>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ExifData {
    pub make: Option<String>,
    pub model: Option<String>,
    pub software: Option<String>,
    pub date_time: Option<String>,
    pub exposure_time: Option<String>,
    pub f_number: Option<f64>,
    pub iso_speed: Option<u32>,
    pub focal_length: Option<f64>,
    pub lens_make: Option<String>,
    pub lens_model: Option<String>,
    pub flash: Option<String>,
    pub white_balance: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct EditOperation {
    pub operation: String,
    pub parameters: serde_json::Value,
    pub timestamp: DateTime<Utc>,
}

// Request/Response DTOs

#[derive(Debug, Clone, Serialize, Deserialize, Validate)]
pub struct UploadRequest {
    #[validate(length(min = 1, max = 255))]
    pub filename: Option<String>,
    pub media_type: Option<MediaType>,
    pub is_public: Option<bool>,
    pub user_id: Uuid,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct UploadResponse {
    pub media_id: Uuid,
    pub url: String,
    pub thumbnail_url: Option<String>,
    pub cdn_url: Option<String>,
    pub media_type: MediaType,
    pub file_size: i64,
    pub width: Option<i32>,
    pub height: Option<i32>,
    pub duration: Option<f64>,
    pub processing_status: ProcessingStatus,
}

#[derive(Debug, Clone, Serialize, Deserialize, Validate)]
pub struct ProcessMediaRequest {
    pub media_id: Uuid,
    pub operations: Vec<ProcessingOperation>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(tag = "type", rename_all = "snake_case")]
pub enum ProcessingOperation {
    Resize {
        width: u32,
        height: u32,
        maintain_aspect_ratio: bool,
    },
    Crop {
        x: u32,
        y: u32,
        width: u32,
        height: u32,
    },
    Rotate {
        degrees: i32,
    },
    Flip {
        horizontal: bool,
        vertical: bool,
    },
    Filter {
        filter_name: String,
        intensity: f32,
    },
    Compress {
        quality: u8,
        format: Option<String>,
    },
    Watermark {
        text: Option<String>,
        image_url: Option<String>,
        position: WatermarkPosition,
        opacity: f32,
    },
    Transcode {
        codec: String,
        bitrate: Option<u32>,
        resolution: Option<String>,
    },
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "snake_case")]
pub enum WatermarkPosition {
    TopLeft,
    TopRight,
    BottomLeft,
    BottomRight,
    Center,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MediaQuery {
    pub user_id: Option<Uuid>,
    pub media_type: Option<MediaType>,
    pub processing_status: Option<ProcessingStatus>,
    pub is_public: Option<bool>,
    pub limit: Option<i64>,
    pub offset: Option<i64>,
    pub sort_by: Option<String>,
    pub sort_order: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MediaListResponse {
    pub media: Vec<Media>,
    pub total: i64,
    pub limit: i64,
    pub offset: i64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MediaStats {
    pub total_files: i64,
    pub total_size_bytes: i64,
    pub by_type: std::collections::HashMap<String, MediaTypeStats>,
    pub processing_queue_length: i64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MediaTypeStats {
    pub count: i64,
    pub total_size_bytes: i64,
    pub average_size_bytes: f64,
}

impl Media {
    pub fn new(
        user_id: Uuid,
        filename: String,
        original_filename: String,
        mime_type: String,
        media_type: MediaType,
        file_size: i64,
        storage_path: String,
        storage_provider: String,
    ) -> Self {
        let now = Utc::now();
        Self {
            id: Uuid::new_v4(),
            user_id,
            filename,
            original_filename,
            mime_type,
            media_type,
            file_size,
            storage_path,
            storage_provider,
            url: String::new(),
            cdn_url: None,
            thumbnail_url: None,
            width: None,
            height: None,
            duration: None,
            aspect_ratio: None,
            orientation: None,
            hash: String::new(),
            blurhash: None,
            processing_status: ProcessingStatus::Pending,
            variants: sqlx::types::Json(Vec::new()),
            metadata: sqlx::types::Json(MediaMetadata::default()),
            is_processed: false,
            is_public: false,
            is_deleted: false,
            view_count: 0,
            download_count: 0,
            created_at: now,
            updated_at: now,
            deleted_at: None,
        }
    }

    pub fn calculate_aspect_ratio(&self) -> Option<f64> {
        match (self.width, self.height) {
            (Some(w), Some(h)) if h > 0 => Some(w as f64 / h as f64),
            _ => None,
        }
    }

    pub fn is_image(&self) -> bool {
        self.media_type == MediaType::Image
    }

    pub fn is_video(&self) -> bool {
        self.media_type == MediaType::Video
    }

    pub fn is_audio(&self) -> bool {
        self.media_type == MediaType::Audio
    }
}

impl From<&Media> for UploadResponse {
    fn from(media: &Media) -> Self {
        Self {
            media_id: media.id,
            url: media.url.clone(),
            thumbnail_url: media.thumbnail_url.clone(),
            cdn_url: media.cdn_url.clone(),
            media_type: media.media_type.clone(),
            file_size: media.file_size,
            width: media.width,
            height: media.height,
            duration: media.duration,
            processing_status: media.processing_status.clone(),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_media_creation() {
        let media = Media::new(
            Uuid::new_v4(),
            "test.jpg".to_string(),
            "original.jpg".to_string(),
            "image/jpeg".to_string(),
            MediaType::Image,
            1024,
            "/path/to/file".to_string(),
            "local".to_string(),
        );

        assert_eq!(media.processing_status, ProcessingStatus::Pending);
        assert!(!media.is_processed);
        assert!(!media.is_deleted);
    }

    #[test]
    fn test_aspect_ratio_calculation() {
        let mut media = Media::new(
            Uuid::new_v4(),
            "test.jpg".to_string(),
            "original.jpg".to_string(),
            "image/jpeg".to_string(),
            MediaType::Image,
            1024,
            "/path/to/file".to_string(),
            "local".to_string(),
        );

        media.width = Some(1920);
        media.height = Some(1080);

        let ratio = media.calculate_aspect_ratio();
        assert!(ratio.is_some());
        assert!((ratio.unwrap() - 1.777).abs() < 0.01);
    }
}
