use actix_web::{web, HttpRequest, HttpResponse, Result as ActixResult};
use sqlx::PgPool;
use uuid::Uuid;

use crate::{
    models::{Media, ProcessMediaRequest, ProcessingOperation, ProcessingStatus},
    services::{ImageProcessor, CompressionService},
    AppState,
};

/// Process media with requested operations
pub async fn process_media(
    req: HttpRequest,
    path: web::Path<Uuid>,
    body: web::Json<ProcessMediaRequest>,
    data: web::Data<AppState>,
) -> ActixResult<HttpResponse> {
    let media_id = path.into_inner();
    let user_id = extract_user_id(&req)?;

    // Get media
    let media = get_media_by_id(&data.db, media_id).await.map_err(|_| {
        actix_web::error::ErrorNotFound("Media not found")
    })?;

    // Verify ownership
    if media.user_id != user_id {
        return Ok(HttpResponse::Forbidden().json(serde_json::json!({
            "error": "Permission denied"
        })));
    }

    // Update status to processing
    update_processing_status(&data.db, media_id, ProcessingStatus::Processing).await.ok();

    // Download original file
    let file_data = data
        .storage
        .download(&media.storage_path)
        .await
        .map_err(|e| actix_web::error::ErrorInternalServerError(e.to_string()))?;

    // Process image
    if media.is_image() {
        let processor = ImageProcessor::default();
        let mut processed_data = file_data.to_vec();

        // Apply operations sequentially
        for operation in &body.operations {
            processed_data = match operation {
                ProcessingOperation::Resize { width, height, maintain_aspect_ratio } => {
                    processor.resize(&processed_data, *width, *height, *maintain_aspect_ratio)
                        .await
                        .map_err(|e| actix_web::error::ErrorInternalServerError(e.to_string()))?
                }
                ProcessingOperation::Crop { x, y, width, height } => {
                    processor.crop(&processed_data, *x, *y, *width, *height)
                        .await
                        .map_err(|e| actix_web::error::ErrorInternalServerError(e.to_string()))?
                }
                ProcessingOperation::Rotate { degrees } => {
                    processor.rotate(&processed_data, *degrees)
                        .await
                        .map_err(|e| actix_web::error::ErrorInternalServerError(e.to_string()))?
                }
                ProcessingOperation::Flip { horizontal, vertical } => {
                    processor.flip(&processed_data, *horizontal, *vertical)
                        .await
                        .map_err(|e| actix_web::error::ErrorInternalServerError(e.to_string()))?
                }
                ProcessingOperation::Compress { quality, format } => {
                    let compression = CompressionService::default();
                    compression.compress_quality(&processed_data, *quality)
                        .await
                        .map_err(|e| actix_web::error::ErrorInternalServerError(e.to_string()))?
                }
                _ => processed_data, // Other operations not implemented yet
            };
        }

        // Upload processed file (replace original or save as variant)
        let processed_path = format!("{}.processed", media.storage_path);
        data.storage
            .upload(&processed_path, bytes::Bytes::from(processed_data.clone()), &media.mime_type)
            .await
            .map_err(|e| actix_web::error::ErrorInternalServerError(e.to_string()))?;

        // Update media record
        update_processing_status(&data.db, media_id, ProcessingStatus::Completed).await.ok();

        return Ok(HttpResponse::Ok().json(serde_json::json!({
            "success": true,
            "media_id": media_id,
            "operations_applied": body.operations.len(),
            "status": "completed"
        })));
    }

    // For video/audio, queue for background processing
    Ok(HttpResponse::Accepted().json(serde_json::json!({
        "success": true,
        "media_id": media_id,
        "status": "queued",
        "message": "Video/audio processing queued"
    })))
}

/// Get processing status
pub async fn get_processing_status(
    path: web::Path<Uuid>,
    data: web::Data<AppState>,
) -> ActixResult<HttpResponse> {
    let media_id = path.into_inner();

    let media = get_media_by_id(&data.db, media_id).await.map_err(|_| {
        actix_web::error::ErrorNotFound("Media not found")
    })?;

    Ok(HttpResponse::Ok().json(serde_json::json!({
        "media_id": media.id,
        "status": media.processing_status,
        "is_processed": media.is_processed,
        "variants": media.variants,
    })))
}

/// Batch process multiple media files
pub async fn batch_process(
    req: HttpRequest,
    body: web::Json<BatchProcessRequest>,
    data: web::Data<AppState>,
) -> ActixResult<HttpResponse> {
    let user_id = extract_user_id(&req)?;

    let mut results = Vec::new();

    for media_id in &body.media_ids {
        // Queue each media for processing
        match queue_processing(&data.db, *media_id, user_id).await {
            Ok(_) => results.push(serde_json::json!({
                "media_id": media_id,
                "status": "queued"
            })),
            Err(_) => results.push(serde_json::json!({
                "media_id": media_id,
                "status": "failed",
                "error": "Failed to queue"
            })),
        }
    }

    Ok(HttpResponse::Ok().json(serde_json::json!({
        "success": true,
        "total": body.media_ids.len(),
        "results": results
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

async fn get_media_by_id(pool: &PgPool, media_id: Uuid) -> Result<Media, sqlx::Error> {
    sqlx::query_as::<_, Media>("SELECT * FROM media WHERE id = $1 AND is_deleted = false")
        .bind(media_id)
        .fetch_one(pool)
        .await
}

async fn update_processing_status(
    pool: &PgPool,
    media_id: Uuid,
    status: ProcessingStatus,
) -> Result<(), sqlx::Error> {
    sqlx::query("UPDATE media SET processing_status = $1, updated_at = $2 WHERE id = $3")
        .bind(status)
        .bind(chrono::Utc::now())
        .bind(media_id)
        .execute(pool)
        .await?;
    Ok(())
}

async fn queue_processing(pool: &PgPool, media_id: Uuid, user_id: Uuid) -> Result<(), sqlx::Error> {
    // Verify ownership
    let media = sqlx::query_as::<_, Media>(
        "SELECT * FROM media WHERE id = $1 AND user_id = $2 AND is_deleted = false"
    )
    .bind(media_id)
    .bind(user_id)
    .fetch_one(pool)
    .await?;

    // Update status
    update_processing_status(pool, media_id, ProcessingStatus::Pending).await
}

#[derive(serde::Deserialize)]
struct BatchProcessRequest {
    media_ids: Vec<Uuid>,
    operations: Option<Vec<ProcessingOperation>>,
}
