use actix_multipart::Multipart;
use actix_web::{web, HttpRequest, HttpResponse, Result as ActixResult};
use bytes::Bytes;
use futures_util::StreamExt;
use sqlx::PgPool;
use std::sync::Arc;
use uuid::Uuid;

use crate::{
    models::{Media, MediaType, ProcessingStatus, UploadResponse},
    services::{ImageProcessor, ThumbnailGenerator},
    storage::StorageBackend,
    utils::{crypto, mime_types, validation::FileValidator},
    AppState,
};

const MAX_CHUNK_SIZE: usize = 10 * 1024 * 1024; // 10MB chunks

/// Upload single file with automatic processing
pub async fn upload_media(
    req: HttpRequest,
    mut payload: Multipart,
    data: web::Data<AppState>,
) -> ActixResult<HttpResponse> {
    // Extract user ID from Authorization header
    let user_id = extract_user_id(&req)?;

    let mut file_data = Vec::new();
    let mut original_filename = String::new();
    let mut content_type = String::new();

    // Parse multipart form data
    while let Some(item) = payload.next().await {
        let mut field = item.map_err(|e| {
            actix_web::error::ErrorBadRequest(format!("Multipart error: {}", e))
        })?;

        let content_disposition = field.content_disposition();
        
        if let Some(name) = content_disposition.get_name() {
            if name == "file" {
                original_filename = content_disposition
                    .get_filename()
                    .unwrap_or("unnamed")
                    .to_string();
                
                content_type = field
                    .content_type()
                    .map(|m| m.to_string())
                    .unwrap_or_else(|| "application/octet-stream".to_string());

                // Read file data
                while let Some(chunk) = field.next().await {
                    let chunk_data = chunk.map_err(|e| {
                        actix_web::error::ErrorBadRequest(format!("Chunk read error: {}", e))
                    })?;
                    file_data.extend_from_slice(&chunk_data);
                }
            }
        }
    }

    if file_data.is_empty() {
        return Ok(HttpResponse::BadRequest().json(serde_json::json!({
            "error": "No file data received"
        })));
    }

    // Detect MIME type from content
    let detected_mime = mime_types::detect_mime_from_bytes(&file_data)
        .or_else(|| mime_types::detect_mime_type(&original_filename))
        .ok_or_else(|| actix_web::error::ErrorBadRequest("Unable to determine file type"))?;

    let media_type = mime_types::get_media_type_from_mime(&detected_mime);

    // Validate file
    let validator = FileValidator::new(data.config.limits.clone());
    match media_type {
        MediaType::Image => {
            validator.validate_image(file_data.len() as u64, &content_type)
                .map_err(|e| actix_web::error::ErrorBadRequest(e.to_string()))?;
        }
        MediaType::Video => {
            validator.validate_video(file_data.len() as u64, &content_type)
                .map_err(|e| actix_web::error::ErrorBadRequest(e.to_string()))?;
        }
        MediaType::Audio => {
            validator.validate_audio(file_data.len() as u64, &content_type)
                .map_err(|e| actix_web::error::ErrorBadRequest(e.to_string()))?;
        }
        _ => {}
    }

    // Generate unique filename with extension
    let extension = mime_types::get_extension_from_mime(&detected_mime)
        .unwrap_or("bin");
    let unique_filename = format!("{}.{}", Uuid::new_v4(), extension);

    // Calculate file hash for deduplication
    let file_hash = crypto::generate_file_hash(&file_data);

    // Check for duplicate (optional optimization)
    let existing_media = check_duplicate_hash(&data.db, &file_hash, user_id).await;
    if let Some(existing) = existing_media {
        tracing::info!("Duplicate file detected, returning existing media");
        return Ok(HttpResponse::Ok().json(UploadResponse::from(&existing)));
    }

    // Build storage path: {user_id}/{year}/{month}/{filename}
    let now = chrono::Utc::now();
    let storage_path = format!(
        "{}/{}/{}/{}",
        user_id,
        now.format("%Y"),
        now.format("%m"),
        unique_filename
    );

    // Upload to storage
    let upload_result = data
        .storage
        .upload(&storage_path, Bytes::from(file_data.clone()), &content_type)
        .await
        .map_err(|e| {
            actix_web::error::ErrorInternalServerError(format!("Storage upload failed: {}", e))
        })?;

    // Get public URL
    let url = data.storage.get_url(&storage_path).await.map_err(|e| {
        actix_web::error::ErrorInternalServerError(format!("Failed to get URL: {}", e))
    })?;

    // Create media record
    let mut media = Media::new(
        user_id,
        unique_filename.clone(),
        original_filename,
        content_type,
        media_type.clone(),
        file_data.len() as i64,
        storage_path.clone(),
        data.config.storage.provider.clone(),
    );

    media.url = url.clone();
    media.hash = file_hash;

    // Process based on media type
    match media_type {
        MediaType::Image => {
            // Extract image dimensions and metadata
            if let Ok(img) = image::load_from_memory(&file_data) {
                let (width, height) = img.dimensions();
                media.width = Some(width as i32);
                media.height = Some(height as i32);
                media.aspect_ratio = Some(width as f64 / height as f64);

                // Generate thumbnails asynchronously
                let processor = ImageProcessor::default();
                let generator = ThumbnailGenerator::new(processor);
                
                if let Ok(thumbnails) = generator.generate_thumbnails(
                    &file_data,
                    &[(150, 150, "thumb"), (600, 600, "preview")]
                ).await {
                    if let Some(thumb) = thumbnails.first() {
                        // Upload thumbnail
                        let thumb_path = format!("{}/thumbnails/{}", user_id, unique_filename);
                        if data.storage.upload(&thumb_path, Bytes::from(thumb.data.clone()), "image/jpeg").await.is_ok() {
                            if let Ok(thumb_url) = data.storage.get_url(&thumb_path).await {
                                media.thumbnail_url = Some(thumb_url);
                            }
                        }
                    }
                }

                media.processing_status = ProcessingStatus::Completed;
                media.is_processed = true;
            }
        }
        MediaType::Video => {
            // Video processing will be handled asynchronously
            media.processing_status = ProcessingStatus::Pending;
            tracing::info!("Video queued for processing: {}", media.id);
        }
        MediaType::Audio => {
            // Audio processing
            media.processing_status = ProcessingStatus::Pending;
            tracing::info!("Audio queued for processing: {}", media.id);
        }
        _ => {}
    }

    // Save to database
    save_media_to_db(&data.db, &media).await.map_err(|e| {
        actix_web::error::ErrorInternalServerError(format!("Database error: {}", e))
    })?;

    // Cache media info in Redis
    cache_media_info(&data.redis, &media).await.ok();

    tracing::info!(
        "Media uploaded successfully: id={}, type={:?}, size={}",
        media.id, media.media_type, media.file_size
    );

    Ok(HttpResponse::Created().json(UploadResponse::from(&media)))
}

/// Initialize multipart upload for large files
pub async fn init_multipart_upload(
    req: HttpRequest,
    body: web::Json<serde_json::Value>,
    data: web::Data<AppState>,
) -> ActixResult<HttpResponse> {
    let user_id = extract_user_id(&req)?;

    let filename = body.get("filename")
        .and_then(|v| v.as_str())
        .ok_or_else(|| actix_web::error::ErrorBadRequest("filename required"))?;

    let file_size = body.get("file_size")
        .and_then(|v| v.as_i64())
        .ok_or_else(|| actix_web::error::ErrorBadRequest("file_size required"))?;

    let chunk_size = body.get("chunk_size")
        .and_then(|v| v.as_i64())
        .unwrap_or(5 * 1024 * 1024); // Default 5MB

    let total_chunks = ((file_size as f64) / (chunk_size as f64)).ceil() as i32;

    let upload_id = Uuid::new_v4();
    let media_id = Uuid::new_v4();

    // Store upload session in Redis
    let session_key = format!("upload_session:{}", upload_id);
    let session_data = serde_json::json!({
        "upload_id": upload_id,
        "media_id": media_id,
        "filename": filename,
        "file_size": file_size,
        "chunk_size": chunk_size,
        "total_chunks": total_chunks,
        "uploaded_chunks": [],
        "created_at": chrono::Utc::now(),
    });

    // Cache session for 24 hours
    let mut conn = data.redis.get_async_connection().await.map_err(|e| {
        actix_web::error::ErrorInternalServerError(format!("Redis error: {}", e))
    })?;

    redis::cmd("SET")
        .arg(&session_key)
        .arg(serde_json::to_string(&session_data).unwrap())
        .arg("EX")
        .arg(86400) // 24 hours
        .query_async::<_, ()>(&mut conn)
        .await
        .map_err(|e| actix_web::error::ErrorInternalServerError(e.to_string()))?;

    Ok(HttpResponse::Ok().json(serde_json::json!({
        "upload_id": upload_id,
        "media_id": media_id,
        "chunk_size": chunk_size,
        "total_chunks": total_chunks,
    })))
}

/// Upload individual chunk
pub async fn upload_chunk(
    req: HttpRequest,
    body: web::Bytes,
    query: web::Query<ChunkQuery>,
    data: web::Data<AppState>,
) -> ActixResult<HttpResponse> {
    let user_id = extract_user_id(&req)?;

    // Validate chunk
    if body.len() > MAX_CHUNK_SIZE {
        return Ok(HttpResponse::PayloadTooLarge().json(serde_json::json!({
            "error": "Chunk size exceeds maximum"
        })));
    }

    // Store chunk in temporary storage
    let chunk_path = format!(
        "temp/{}/chunk_{}",
        query.upload_id,
        query.chunk_number
    );

    data.storage
        .upload(&chunk_path, body, "application/octet-stream")
        .await
        .map_err(|e| actix_web::error::ErrorInternalServerError(e.to_string()))?;

    // Update session in Redis
    let session_key = format!("upload_session:{}", query.upload_id);
    let mut conn = data.redis.get_async_connection().await.map_err(|e| {
        actix_web::error::ErrorInternalServerError(e.to_string())
    })?;

    // Add chunk to uploaded list
    redis::cmd("SADD")
        .arg(format!("{}:chunks", session_key))
        .arg(query.chunk_number)
        .query_async::<_, ()>(&mut conn)
        .await
        .ok();

    Ok(HttpResponse::Ok().json(serde_json::json!({
        "chunk_number": query.chunk_number,
        "status": "uploaded"
    })))
}

/// Complete multipart upload by assembling chunks
pub async fn complete_multipart_upload(
    req: HttpRequest,
    body: web::Json<CompleteUploadRequest>,
    data: web::Data<AppState>,
) -> ActixResult<HttpResponse> {
    let user_id = extract_user_id(&req)?;

    // Retrieve session from Redis
    let session_key = format!("upload_session:{}", body.upload_id);
    let mut conn = data.redis.get_async_connection().await.map_err(|e| {
        actix_web::error::ErrorInternalServerError(e.to_string())
    })?;

    let session_data: String = redis::cmd("GET")
        .arg(&session_key)
        .query_async(&mut conn)
        .await
        .map_err(|e| actix_web::error::ErrorNotFound("Upload session not found"))?;

    let session: serde_json::Value = serde_json::from_str(&session_data)
        .map_err(|e| actix_web::error::ErrorInternalServerError(e.to_string()))?;

    let media_id: Uuid = serde_json::from_value(session["media_id"].clone())
        .map_err(|_| actix_web::error::ErrorBadRequest("Invalid media_id"))?;

    let total_chunks = session["total_chunks"].as_i64().unwrap_or(0);

    // Assemble chunks
    let mut assembled_data = Vec::new();
    for chunk_num in 0..total_chunks {
        let chunk_path = format!("temp/{}/chunk_{}", body.upload_id, chunk_num);
        
        match data.storage.download(&chunk_path).await {
            Ok(chunk_data) => {
                assembled_data.extend_from_slice(&chunk_data);
                // Delete chunk after reading
                data.storage.delete(&chunk_path).await.ok();
            }
            Err(e) => {
                tracing::error!("Failed to read chunk {}: {}", chunk_num, e);
                return Ok(HttpResponse::InternalServerError().json(serde_json::json!({
                    "error": format!("Missing chunk: {}", chunk_num)
                })));
            }
        }
    }

    // Process the assembled file (similar to single upload)
    let filename = session["filename"].as_str().unwrap_or("unnamed");
    let detected_mime = mime_types::detect_mime_from_bytes(&assembled_data)
        .or_else(|| mime_types::detect_mime_type(filename))
        .ok_or_else(|| actix_web::error::ErrorBadRequest("Unable to determine file type"))?;

    let media_type = mime_types::get_media_type_from_mime(&detected_mime);
    let file_hash = crypto::generate_file_hash(&assembled_data);

    // Generate unique filename
    let extension = mime_types::get_extension_from_mime(&detected_mime).unwrap_or("bin");
    let unique_filename = format!("{}.{}", Uuid::new_v4(), extension);

    // Build storage path
    let now = chrono::Utc::now();
    let storage_path = format!(
        "{}/{}/{}/{}",
        user_id,
        now.format("%Y"),
        now.format("%m"),
        unique_filename
    );

    // Upload assembled file
    data.storage
        .upload(&storage_path, Bytes::from(assembled_data.clone()), detected_mime.as_ref())
        .await
        .map_err(|e| actix_web::error::ErrorInternalServerError(e.to_string()))?;

    let url = data.storage.get_url(&storage_path).await.map_err(|e| {
        actix_web::error::ErrorInternalServerError(e.to_string())
    })?;

    // Create media record
    let mut media = Media::new(
        user_id,
        unique_filename,
        filename.to_string(),
        detected_mime.to_string(),
        media_type,
        assembled_data.len() as i64,
        storage_path,
        data.config.storage.provider.clone(),
    );

    media.url = url;
    media.hash = file_hash;
    media.id = media_id;

    // Process image if applicable
    if media.is_image() {
        if let Ok(img) = image::load_from_memory(&assembled_data) {
            let (width, height) = img.dimensions();
            media.width = Some(width as i32);
            media.height = Some(height as i32);
            media.aspect_ratio = Some(width as f64 / height as f64);
            media.processing_status = ProcessingStatus::Completed;
            media.is_processed = true;
        }
    }

    // Save to database
    save_media_to_db(&data.db, &media).await.map_err(|e| {
        actix_web::error::ErrorInternalServerError(e.to_string())
    })?;

    // Clean up session
    redis::cmd("DEL")
        .arg(&session_key)
        .query_async::<_, ()>(&mut conn)
        .await
        .ok();

    tracing::info!("Multipart upload completed: media_id={}", media.id);

    Ok(HttpResponse::Created().json(UploadResponse::from(&media)))
}

/// Delete media file
pub async fn delete_media(
    req: HttpRequest,
    path: web::Path<Uuid>,
    data: web::Data<AppState>,
) -> ActixResult<HttpResponse> {
    let user_id = extract_user_id(&req)?;
    let media_id = path.into_inner();

    // Get media from database
    let media = get_media_by_id(&data.db, media_id).await.map_err(|_| {
        actix_web::error::ErrorNotFound("Media not found")
    })?;

    // Verify ownership
    if media.user_id != user_id {
        return Ok(HttpResponse::Forbidden().json(serde_json::json!({
            "error": "You don't have permission to delete this media"
        })));
    }

    // Soft delete in database
    soft_delete_media(&data.db, media_id).await.map_err(|e| {
        actix_web::error::ErrorInternalServerError(e.to_string())
    })?;

    // Delete from storage (async)
    let storage = data.storage.clone();
    let storage_path = media.storage_path.clone();
    tokio::spawn(async move {
        if let Err(e) = storage.delete(&storage_path).await {
            tracing::error!("Failed to delete file from storage: {}", e);
        }
    });

    // Invalidate cache
    invalidate_media_cache(&data.redis, media_id).await.ok();

    Ok(HttpResponse::Ok().json(serde_json::json!({
        "success": true,
        "message": "Media deleted successfully"
    })))
}

// Helper functions

fn extract_user_id(req: &HttpRequest) -> ActixResult<Uuid> {
    req.headers()
        .get("X-User-ID")
        .and_then(|h| h.to_str().ok())
        .and_then(|s| Uuid::parse_str(s).ok())
        .ok_or_else(|| actix_web::error::ErrorUnauthorized("Missing or invalid user ID"))
}

async fn check_duplicate_hash(pool: &PgPool, hash: &str, user_id: Uuid) -> Option<Media> {
    sqlx::query_as::<_, Media>(
        "SELECT * FROM media WHERE hash = $1 AND user_id = $2 AND is_deleted = false LIMIT 1"
    )
    .bind(hash)
    .bind(user_id)
    .fetch_optional(pool)
    .await
    .ok()
    .flatten()
}

async fn save_media_to_db(pool: &PgPool, media: &Media) -> Result<(), sqlx::Error> {
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
    .execute(pool)
    .await?;

    Ok(())
}

async fn get_media_by_id(pool: &PgPool, media_id: Uuid) -> Result<Media, sqlx::Error> {
    sqlx::query_as::<_, Media>("SELECT * FROM media WHERE id = $1 AND is_deleted = false")
        .bind(media_id)
        .fetch_one(pool)
        .await
}

async fn soft_delete_media(pool: &PgPool, media_id: Uuid) -> Result<(), sqlx::Error> {
    sqlx::query(
        "UPDATE media SET is_deleted = true, deleted_at = $1 WHERE id = $2"
    )
    .bind(chrono::Utc::now())
    .bind(media_id)
    .execute(pool)
    .await?;

    Ok(())
}

async fn cache_media_info(redis_client: &redis::Client, media: &Media) -> Result<(), redis::RedisError> {
    let mut conn = redis_client.get_async_connection().await?;
    let cache_key = format!("media:{}", media.id);
    let media_json = serde_json::to_string(media).unwrap_or_default();

    redis::cmd("SET")
        .arg(cache_key)
        .arg(media_json)
        .arg("EX")
        .arg(3600) // 1 hour
        .query_async(&mut conn)
        .await
}

async fn invalidate_media_cache(redis_client: &redis::Client, media_id: Uuid) -> Result<(), redis::RedisError> {
    let mut conn = redis_client.get_async_connection().await?;
    let cache_key = format!("media:{}", media_id);
    redis::cmd("DEL").arg(cache_key).query_async(&mut conn).await
}

#[derive(Debug, Deserialize)]
struct ChunkQuery {
    upload_id: Uuid,
    chunk_number: i32,
}

#[derive(Debug, Deserialize)]
struct CompleteUploadRequest {
    upload_id: Uuid,
}

use serde::Deserialize;
