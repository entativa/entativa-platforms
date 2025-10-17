use tonic::{Request, Response, Status};
use uuid::Uuid;
use std::sync::Arc;
use sqlx::PgPool;
use bytes::Bytes;

use crate::models::{Media, MediaType as ModelMediaType, ProcessingStatus as ModelProcessingStatus};
use crate::storage::StorageBackend;
use crate::services::{ImageProcessor, ThumbnailGenerator};
use crate::utils::{crypto, mime_types, validation::FileValidator};
use crate::config::Config;

// Include generated protobuf code
pub mod media {
    tonic::include_proto!("media");
}

use media::{
    media_service_server::{MediaService, MediaServiceServer},
    *,
};

pub struct MediaServiceImpl {
    db: PgPool,
    redis: redis::Client,
    storage: Arc<dyn StorageBackend>,
    config: Config,
}

impl MediaServiceImpl {
    pub fn new(
        db: PgPool,
        redis: redis::Client,
        storage: Arc<dyn StorageBackend>,
        config: Config,
    ) -> Self {
        Self {
            db,
            redis,
            storage,
            config,
        }
    }

    pub fn into_service(self) -> MediaServiceServer<Self> {
        MediaServiceServer::new(self)
    }
}

#[tonic::async_trait]
impl MediaService for MediaServiceImpl {
    async fn upload_media(
        &self,
        request: Request<UploadMediaRequest>,
    ) -> Result<Response<UploadMediaResponse>, Status> {
        let req = request.into_inner();

        // Parse user_id
        let user_id = Uuid::parse_str(&req.user_id)
            .map_err(|_| Status::invalid_argument("Invalid user_id"))?;

        // Validate file
        let validator = FileValidator::new(self.config.limits.clone());
        let file_data = req.data;
        
        // Detect MIME type
        let detected_mime = mime_types::detect_mime_from_bytes(&file_data)
            .or_else(|| mime_types::detect_mime_type(&req.filename))
            .ok_or_else(|| Status::invalid_argument("Unable to determine file type"))?;

        let media_type = mime_types::get_media_type_from_mime(&detected_mime);

        // Validate based on type
        match media_type {
            ModelMediaType::Image => {
                validator.validate_image(file_data.len() as u64, detected_mime.as_ref())
                    .map_err(|e| Status::invalid_argument(e.to_string()))?;
            }
            ModelMediaType::Video => {
                validator.validate_video(file_data.len() as u64, detected_mime.as_ref())
                    .map_err(|e| Status::invalid_argument(e.to_string()))?;
            }
            ModelMediaType::Audio => {
                validator.validate_audio(file_data.len() as u64, detected_mime.as_ref())
                    .map_err(|e| Status::invalid_argument(e.to_string()))?;
            }
            _ => {}
        }

        // Generate unique filename
        let extension = mime_types::get_extension_from_mime(&detected_mime).unwrap_or("bin");
        let unique_filename = format!("{}.{}", Uuid::new_v4(), extension);

        // Calculate file hash
        let file_hash = crypto::generate_file_hash(&file_data);

        // Check for duplicates
        if let Ok(Some(existing)) = self.check_duplicate_hash(&file_hash, user_id).await {
            // Return existing media
            return Ok(Response::new(self.media_to_response(&existing)));
        }

        // Build storage path
        let now = chrono::Utc::now();
        let storage_path = format!(
            "{}/{}/{}/{}",
            user_id,
            now.format("%Y"),
            now.format("%m"),
            unique_filename
        );

        // Upload to storage
        self.storage
            .upload(&storage_path, Bytes::from(file_data.clone()), detected_mime.as_ref())
            .await
            .map_err(|e| Status::internal(format!("Storage upload failed: {}", e)))?;

        // Get URL
        let url = self.storage
            .get_url(&storage_path)
            .await
            .map_err(|e| Status::internal(format!("Failed to get URL: {}", e)))?;

        // Create media record
        let mut media = Media::new(
            user_id,
            unique_filename.clone(),
            req.filename,
            detected_mime.to_string(),
            media_type.clone(),
            file_data.len() as i64,
            storage_path.clone(),
            self.config.storage.provider.clone(),
        );

        media.url = url.clone();
        media.hash = file_hash;

        // Process image if needed
        if media.is_image() {
            if let Ok(img) = image::load_from_memory(&file_data) {
                let (width, height) = img.dimensions();
                media.width = Some(width as i32);
                media.height = Some(height as i32);
                media.aspect_ratio = Some(width as f64 / height as f64);

                // Generate thumbnails for profile pictures and cover photos
                let should_generate_thumbs = matches!(
                    req.purpose(),
                    MediaPurpose::ProfilePicture | MediaPurpose::CoverPhoto
                );

                if should_generate_thumbs {
                    let processor = ImageProcessor::default();
                    let generator = ThumbnailGenerator::new(processor.clone());

                    // Generate thumbnails
                    let thumbnail_sizes = if req.purpose() == MediaPurpose::ProfilePicture as i32 {
                        vec![(150, 150, "thumb"), (300, 300, "medium"), (600, 600, "large")]
                    } else {
                        vec![(400, 150, "thumb"), (800, 300, "medium"), (1600, 600, "large")]
                    };

                    if let Ok(thumbnails) = generator.generate_thumbnails(&file_data, &thumbnail_sizes).await {
                        if let Some(thumb) = thumbnails.first() {
                            let thumb_path = format!("{}/thumbnails/{}", user_id, unique_filename);
                            if self.storage.upload(&thumb_path, Bytes::from(thumb.data.clone()), "image/jpeg").await.is_ok() {
                                if let Ok(thumb_url) = self.storage.get_url(&thumb_path).await {
                                    media.thumbnail_url = Some(thumb_url);
                                }
                            }
                        }
                    }

                    // Generate blurhash
                    if let Some(blurhash) = processor.generate_blurhash(&img) {
                        media.blurhash = Some(blurhash);
                    }
                }

                media.processing_status = ModelProcessingStatus::Completed;
                media.is_processed = true;
            }
        }

        // Save to database
        self.save_media(&media).await
            .map_err(|e| Status::internal(format!("Database error: {}", e)))?;

        // Cache in Redis
        self.cache_media(&media).await.ok();

        Ok(Response::new(self.media_to_response(&media)))
    }

    async fn upload_media_stream(
        &self,
        request: Request<tonic::Streaming<UploadChunk>>,
    ) -> Result<Response<UploadMediaResponse>, Status> {
        let mut stream = request.into_inner();
        let mut chunks = Vec::new();
        let mut metadata: Option<UploadChunk> = None;

        // Collect all chunks
        while let Some(chunk) = stream.message().await? {
            if metadata.is_none() {
                metadata = Some(chunk.clone());
            }
            chunks.push(chunk);
        }

        let meta = metadata.ok_or_else(|| Status::invalid_argument("No chunks received"))?;

        // Assemble file data
        let file_data: Vec<u8> = chunks.iter().flat_map(|c| c.data.clone()).collect();

        // Create upload request and reuse upload_media logic
        let upload_req = UploadMediaRequest {
            data: file_data,
            filename: meta.filename,
            content_type: meta.content_type,
            user_id: meta.user_id,
            purpose: meta.purpose,
            processing_options: None,
        };

        self.upload_media(Request::new(upload_req)).await
    }

    async fn get_media(
        &self,
        request: Request<GetMediaRequest>,
    ) -> Result<Response<MediaResponse>, Status> {
        let req = request.into_inner();
        
        let media_id = Uuid::parse_str(&req.media_id)
            .map_err(|_| Status::invalid_argument("Invalid media_id"))?;

        // Try cache first
        if let Ok(Some(media)) = self.get_media_from_cache(media_id).await {
            return Ok(Response::new(self.media_to_full_response(&media)));
        }

        // Query database
        let media = self.get_media_by_id(media_id).await
            .map_err(|_| Status::not_found("Media not found"))?;

        // Cache it
        self.cache_media(&media).await.ok();

        Ok(Response::new(self.media_to_full_response(&media)))
    }

    async fn delete_media(
        &self,
        request: Request<DeleteMediaRequest>,
    ) -> Result<Response<DeleteMediaResponse>, Status> {
        let req = request.into_inner();

        let media_id = Uuid::parse_str(&req.media_id)
            .map_err(|_| Status::invalid_argument("Invalid media_id"))?;
        let user_id = Uuid::parse_str(&req.user_id)
            .map_err(|_| Status::invalid_argument("Invalid user_id"))?;

        // Get media and verify ownership
        let media = self.get_media_by_id(media_id).await
            .map_err(|_| Status::not_found("Media not found"))?;

        if media.user_id != user_id {
            return Err(Status::permission_denied("Not authorized to delete this media"));
        }

        // Soft delete
        self.soft_delete_media(media_id).await
            .map_err(|e| Status::internal(format!("Delete failed: {}", e)))?;

        // Delete from storage (async)
        let storage = self.storage.clone();
        let storage_path = media.storage_path.clone();
        tokio::spawn(async move {
            storage.delete(&storage_path).await.ok();
        });

        // Invalidate cache
        self.invalidate_cache(media_id).await.ok();

        Ok(Response::new(DeleteMediaResponse {
            success: true,
            message: "Media deleted successfully".to_string(),
        }))
    }

    async fn process_media(
        &self,
        _request: Request<ProcessMediaRequest>,
    ) -> Result<Response<ProcessMediaResponse>, Status> {
        // Processing operations will be queued for async processing
        Err(Status::unimplemented("Processing via gRPC not yet implemented"))
    }

    async fn get_signed_url(
        &self,
        request: Request<GetSignedUrlRequest>,
    ) -> Result<Response<GetSignedUrlResponse>, Status> {
        let req = request.into_inner();

        let media_id = Uuid::parse_str(&req.media_id)
            .map_err(|_| Status::invalid_argument("Invalid media_id"))?;

        let media = self.get_media_by_id(media_id).await
            .map_err(|_| Status::not_found("Media not found"))?;

        let expiry = if req.expiry_seconds > 0 {
            req.expiry_seconds as u64
        } else {
            3600 // Default 1 hour
        };

        let signed_url = self.storage
            .get_signed_url(&media.storage_path, expiry)
            .await
            .map_err(|e| Status::internal(format!("Failed to generate signed URL: {}", e)))?;

        let expires_at = chrono::Utc::now().timestamp() + expiry as i64;

        Ok(Response::new(GetSignedUrlResponse {
            signed_url,
            expires_at,
        }))
    }

    async fn batch_get_media(
        &self,
        request: Request<BatchGetMediaRequest>,
    ) -> Result<Response<BatchGetMediaResponse>, Status> {
        let req = request.into_inner();
        let mut media_list = Vec::new();

        for media_id_str in req.media_ids {
            if let Ok(media_id) = Uuid::parse_str(&media_id_str) {
                if let Ok(media) = self.get_media_by_id(media_id).await {
                    media_list.push(self.media_to_full_response(&media));
                }
            }
        }

        Ok(Response::new(BatchGetMediaResponse { media: media_list }))
    }

    async fn batch_delete_media(
        &self,
        request: Request<BatchDeleteMediaRequest>,
    ) -> Result<Response<BatchDeleteMediaResponse>, Status> {
        let req = request.into_inner();
        
        let user_id = Uuid::parse_str(&req.user_id)
            .map_err(|_| Status::invalid_argument("Invalid user_id"))?;

        let mut deleted_count = 0;
        let mut failed_ids = Vec::new();

        for media_id_str in req.media_ids {
            if let Ok(media_id) = Uuid::parse_str(&media_id_str) {
                if let Ok(media) = self.get_media_by_id(media_id).await {
                    if media.user_id == user_id {
                        if self.soft_delete_media(media_id).await.is_ok() {
                            deleted_count += 1;
                            self.invalidate_cache(media_id).await.ok();
                        } else {
                            failed_ids.push(media_id_str);
                        }
                    } else {
                        failed_ids.push(media_id_str);
                    }
                } else {
                    failed_ids.push(media_id_str);
                }
            } else {
                failed_ids.push(media_id_str);
            }
        }

        Ok(Response::new(BatchDeleteMediaResponse {
            deleted_count,
            failed_ids,
        }))
    }
}

// Helper methods
impl MediaServiceImpl {
    async fn check_duplicate_hash(&self, hash: &str, user_id: Uuid) -> Result<Option<Media>, sqlx::Error> {
        sqlx::query_as::<_, Media>(
            "SELECT * FROM media WHERE hash = $1 AND user_id = $2 AND is_deleted = false LIMIT 1"
        )
        .bind(hash)
        .bind(user_id)
        .fetch_optional(&self.db)
        .await
    }

    async fn save_media(&self, media: &Media) -> Result<(), sqlx::Error> {
        sqlx::query(
            r#"
            INSERT INTO media (
                id, user_id, filename, original_filename, mime_type, media_type,
                file_size, storage_path, storage_provider, url, cdn_url, thumbnail_url,
                width, height, duration, aspect_ratio, orientation, hash, blurhash,
                processing_status, variants, metadata, is_processed, is_public,
                created_at, updated_at
            ) VALUES (
                $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
                $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26
            )
            "#
        )
        .bind(media.id)
        .bind(media.user_id)
        .bind(&media.filename)
        .bind(&media.original_filename)
        .bind(&media.mime_type)
        .bind(&media.media_type)
        .bind(media.file_size)
        .bind(&media.storage_path)
        .bind(&media.storage_provider)
        .bind(&media.url)
        .bind(&media.cdn_url)
        .bind(&media.thumbnail_url)
        .bind(media.width)
        .bind(media.height)
        .bind(media.duration)
        .bind(media.aspect_ratio)
        .bind(media.orientation)
        .bind(&media.hash)
        .bind(&media.blurhash)
        .bind(&media.processing_status)
        .bind(&media.variants)
        .bind(&media.metadata)
        .bind(media.is_processed)
        .bind(media.is_public)
        .bind(media.created_at)
        .bind(media.updated_at)
        .execute(&self.db)
        .await?;

        Ok(())
    }

    async fn get_media_by_id(&self, media_id: Uuid) -> Result<Media, sqlx::Error> {
        sqlx::query_as::<_, Media>("SELECT * FROM media WHERE id = $1 AND is_deleted = false")
            .bind(media_id)
            .fetch_one(&self.db)
            .await
    }

    async fn soft_delete_media(&self, media_id: Uuid) -> Result<(), sqlx::Error> {
        sqlx::query("UPDATE media SET is_deleted = true, deleted_at = $1 WHERE id = $2")
            .bind(chrono::Utc::now())
            .bind(media_id)
            .execute(&self.db)
            .await?;
        Ok(())
    }

    async fn cache_media(&self, media: &Media) -> Result<(), redis::RedisError> {
        let mut conn = self.redis.get_async_connection().await?;
        let cache_key = format!("media:{}", media.id);
        let media_json = serde_json::to_string(media).unwrap_or_default();

        redis::cmd("SET")
            .arg(cache_key)
            .arg(media_json)
            .arg("EX")
            .arg(3600)
            .query_async(&mut conn)
            .await
    }

    async fn get_media_from_cache(&self, media_id: Uuid) -> Result<Option<Media>, redis::RedisError> {
        let mut conn = self.redis.get_async_connection().await?;
        let cache_key = format!("media:{}", media_id);
        
        let data: Option<String> = redis::cmd("GET")
            .arg(cache_key)
            .query_async(&mut conn)
            .await?;

        Ok(data.and_then(|d| serde_json::from_str(&d).ok()))
    }

    async fn invalidate_cache(&self, media_id: Uuid) -> Result<(), redis::RedisError> {
        let mut conn = self.redis.get_async_connection().await?;
        redis::cmd("DEL")
            .arg(format!("media:{}", media_id))
            .query_async(&mut conn)
            .await
    }

    fn media_to_response(&self, media: &Media) -> UploadMediaResponse {
        UploadMediaResponse {
            media_id: media.id.to_string(),
            url: media.url.clone(),
            thumbnail_url: media.thumbnail_url.clone().unwrap_or_default(),
            width: media.width.unwrap_or(0),
            height: media.height.unwrap_or(0),
            file_size: media.file_size,
            mime_type: media.mime_type.clone(),
            media_type: self.convert_media_type(&media.media_type) as i32,
            processing_status: self.convert_processing_status(&media.processing_status) as i32,
            blurhash: media.blurhash.clone().unwrap_or_default(),
        }
    }

    fn media_to_full_response(&self, media: &Media) -> MediaResponse {
        MediaResponse {
            id: media.id.to_string(),
            user_id: media.user_id.to_string(),
            filename: media.filename.clone(),
            original_filename: media.original_filename.clone(),
            mime_type: media.mime_type.clone(),
            media_type: self.convert_media_type(&media.media_type) as i32,
            file_size: media.file_size,
            url: media.url.clone(),
            cdn_url: media.cdn_url.clone().unwrap_or_default(),
            thumbnail_url: media.thumbnail_url.clone().unwrap_or_default(),
            width: media.width.unwrap_or(0),
            height: media.height.unwrap_or(0),
            duration: media.duration.unwrap_or(0.0),
            hash: media.hash.clone(),
            blurhash: media.blurhash.clone().unwrap_or_default(),
            processing_status: self.convert_processing_status(&media.processing_status) as i32,
            is_processed: media.is_processed,
            is_public: media.is_public,
            created_at: media.created_at.to_rfc3339(),
            updated_at: media.updated_at.to_rfc3339(),
        }
    }

    fn convert_media_type(&self, mt: &ModelMediaType) -> MediaType {
        match mt {
            ModelMediaType::Image => MediaType::Image,
            ModelMediaType::Video => MediaType::Video,
            ModelMediaType::Audio => MediaType::Audio,
            ModelMediaType::Document => MediaType::Document,
        }
    }

    fn convert_processing_status(&self, ps: &ModelProcessingStatus) -> ProcessingStatus {
        match ps {
            ModelProcessingStatus::Pending => ProcessingStatus::Pending,
            ModelProcessingStatus::Processing => ProcessingStatus::Processing,
            ModelProcessingStatus::Completed => ProcessingStatus::Completed,
            ModelProcessingStatus::Failed => ProcessingStatus::Failed,
            ModelProcessingStatus::Cancelled => ProcessingStatus::Cancelled,
        }
    }
}
