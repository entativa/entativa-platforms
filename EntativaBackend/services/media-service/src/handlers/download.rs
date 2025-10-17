use actix_web::{web, HttpRequest, HttpResponse, Result as ActixResult};
use futures_util::StreamExt;
use sqlx::PgPool;
use uuid::Uuid;

use crate::{
    models::{Media, MediaQuery, MediaListResponse},
    AppState,
};

/// Get media information
pub async fn get_media(
    req: HttpRequest,
    path: web::Path<Uuid>,
    data: web::Data<AppState>,
) -> ActixResult<HttpResponse> {
    let media_id = path.into_inner();

    // Try cache first
    if let Ok(cached_media) = get_media_from_cache(&data.redis, media_id).await {
        return Ok(HttpResponse::Ok().json(cached_media));
    }

    // Query database
    let media = get_media_by_id(&data.db, media_id).await.map_err(|_| {
        actix_web::error::ErrorNotFound("Media not found")
    })?;

    // Increment view count asynchronously
    let pool = data.db.clone();
    tokio::spawn(async move {
        increment_view_count(&pool, media_id).await.ok();
    });

    // Cache the result
    cache_media(&data.redis, &media).await.ok();

    Ok(HttpResponse::Ok().json(media))
}

/// Download media file with streaming support
pub async fn download_media(
    req: HttpRequest,
    path: web::Path<Uuid>,
    data: web::Data<AppState>,
) -> ActixResult<HttpResponse> {
    let media_id = path.into_inner();

    // Get media info
    let media = get_media_by_id(&data.db, media_id).await.map_err(|_| {
        actix_web::error::ErrorNotFound("Media not found")
    })?;

    // Check if range request (for streaming/resume)
    let range_header = req.headers().get("Range");

    // Download file from storage
    let file_data = data
        .storage
        .download(&media.storage_path)
        .await
        .map_err(|e| {
            actix_web::error::ErrorInternalServerError(format!("Storage error: {}", e))
        })?;

    // Increment download count
    let pool = data.db.clone();
    tokio::spawn(async move {
        increment_download_count(&pool, media_id).await.ok();
    });

    // Handle range requests for partial content
    if let Some(range) = range_header {
        return handle_range_request(&media, file_data, range);
    }

    // Full file download
    Ok(HttpResponse::Ok()
        .content_type(&media.mime_type)
        .insert_header(("Content-Disposition", format!("attachment; filename=\"{}\"", media.original_filename)))
        .insert_header(("Content-Length", file_data.len()))
        .insert_header(("Accept-Ranges", "bytes"))
        .insert_header(("Cache-Control", "public, max-age=31536000")) // 1 year
        .insert_header(("ETag", format!("\"{}\"", media.hash)))
        .body(file_data))
}

/// Get media metadata
pub async fn get_metadata(
    path: web::Path<Uuid>,
    data: web::Data<AppState>,
) -> ActixResult<HttpResponse> {
    let media_id = path.into_inner();

    let media = get_media_by_id(&data.db, media_id).await.map_err(|_| {
        actix_web::error::ErrorNotFound("Media not found")
    })?;

    Ok(HttpResponse::Ok().json(serde_json::json!({
        "id": media.id,
        "filename": media.original_filename,
        "mime_type": media.mime_type,
        "media_type": media.media_type,
        "file_size": media.file_size,
        "width": media.width,
        "height": media.height,
        "duration": media.duration,
        "aspect_ratio": media.aspect_ratio,
        "metadata": media.metadata,
        "created_at": media.created_at,
        "view_count": media.view_count,
        "download_count": media.download_count,
    })))
}

/// List user's media with pagination and filtering
pub async fn list_media(
    req: HttpRequest,
    query: web::Query<MediaQuery>,
    data: web::Data<AppState>,
) -> ActixResult<HttpResponse> {
    let user_id = extract_user_id(&req)?;

    let limit = query.limit.unwrap_or(20).min(100);
    let offset = query.offset.unwrap_or(0);

    // Build query
    let mut sql = String::from(
        "SELECT * FROM media WHERE user_id = $1 AND is_deleted = false"
    );
    let mut params_count = 1;

    // Add filters
    if let Some(media_type) = &query.media_type {
        params_count += 1;
        sql.push_str(&format!(" AND media_type = ${}", params_count));
    }

    if let Some(status) = &query.processing_status {
        params_count += 1;
        sql.push_str(&format!(" AND processing_status = ${}", params_count));
    }

    // Add sorting
    let sort_by = query.sort_by.as_deref().unwrap_or("created_at");
    let sort_order = query.sort_order.as_deref().unwrap_or("desc");
    sql.push_str(&format!(" ORDER BY {} {}", sort_by, sort_order));

    // Add pagination
    sql.push_str(&format!(" LIMIT {} OFFSET {}", limit, offset));

    // Execute query (simplified - in production, use query builder)
    let media_list = sqlx::query_as::<_, Media>(&sql)
        .bind(user_id)
        .fetch_all(&data.db)
        .await
        .map_err(|e| actix_web::error::ErrorInternalServerError(e.to_string()))?;

    // Get total count
    let total: i64 = sqlx::query_scalar(
        "SELECT COUNT(*) FROM media WHERE user_id = $1 AND is_deleted = false"
    )
    .bind(user_id)
    .fetch_one(&data.db)
    .await
    .unwrap_or(0);

    Ok(HttpResponse::Ok().json(MediaListResponse {
        media: media_list,
        total,
        limit,
        offset,
    }))
}

// Helper functions

fn extract_user_id(req: &HttpRequest) -> ActixResult<Uuid> {
    req.headers()
        .get("X-User-ID")
        .and_then(|h| h.to_str().ok())
        .and_then(|s| Uuid::parse_str(s).ok())
        .ok_or_else(|| actix_web::error::ErrorUnauthorized("Missing or invalid user ID"))
}

fn handle_range_request(
    media: &Media,
    data: bytes::Bytes,
    range_header: &actix_web::http::header::HeaderValue,
) -> ActixResult<HttpResponse> {
    let range_str = range_header.to_str().map_err(|_| {
        actix_web::error::ErrorBadRequest("Invalid Range header")
    })?;

    // Parse range: "bytes=start-end"
    let range_parts: Vec<&str> = range_str
        .strip_prefix("bytes=")
        .unwrap_or("")
        .split('-')
        .collect();

    if range_parts.len() != 2 {
        return Ok(HttpResponse::BadRequest().finish());
    }

    let file_size = data.len();
    let start: usize = range_parts[0].parse().unwrap_or(0);
    let end: usize = if range_parts[1].is_empty() {
        file_size - 1
    } else {
        range_parts[1].parse().unwrap_or(file_size - 1).min(file_size - 1)
    };

    if start > end || start >= file_size {
        return Ok(HttpResponse::RangeNotSatisfiable()
            .insert_header(("Content-Range", format!("bytes */{}", file_size)))
            .finish());
    }

    let content_length = end - start + 1;
    let chunk = data.slice(start..=end);

    Ok(HttpResponse::PartialContent()
        .content_type(&media.mime_type)
        .insert_header(("Content-Range", format!("bytes {}-{}/{}", start, end, file_size)))
        .insert_header(("Content-Length", content_length))
        .insert_header(("Accept-Ranges", "bytes"))
        .insert_header(("Cache-Control", "public, max-age=31536000"))
        .body(chunk))
}

async fn get_media_by_id(pool: &PgPool, media_id: Uuid) -> Result<Media, sqlx::Error> {
    sqlx::query_as::<_, Media>("SELECT * FROM media WHERE id = $1 AND is_deleted = false")
        .bind(media_id)
        .fetch_one(pool)
        .await
}

async fn increment_view_count(pool: &PgPool, media_id: Uuid) -> Result<(), sqlx::Error> {
    sqlx::query("UPDATE media SET view_count = view_count + 1 WHERE id = $1")
        .bind(media_id)
        .execute(pool)
        .await?;
    Ok(())
}

async fn increment_download_count(pool: &PgPool, media_id: Uuid) -> Result<(), sqlx::Error> {
    sqlx::query("UPDATE media SET download_count = download_count + 1 WHERE id = $1")
        .bind(media_id)
        .execute(pool)
        .await?;
    Ok(())
}

async fn soft_delete_media(pool: &PgPool, media_id: Uuid) -> Result<(), sqlx::Error> {
    sqlx::query("UPDATE media SET is_deleted = true, deleted_at = $1 WHERE id = $2")
        .bind(chrono::Utc::now())
        .bind(media_id)
        .execute(pool)
        .await?;
    Ok(())
}

async fn get_media_from_cache(
    redis_client: &redis::Client,
    media_id: Uuid,
) -> Result<Media, redis::RedisError> {
    let mut conn = redis_client.get_async_connection().await?;
    let cache_key = format!("media:{}", media_id);
    let data: String = redis::cmd("GET").arg(cache_key).query_async(&mut conn).await?;
    serde_json::from_str(&data).map_err(|e| {
        redis::RedisError::from(std::io::Error::new(std::io::ErrorKind::InvalidData, e))
    })
}

async fn cache_media(redis_client: &redis::Client, media: &Media) -> Result<(), redis::RedisError> {
    let mut conn = redis_client.get_async_connection().await?;
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

async fn invalidate_media_cache(
    redis_client: &redis::Client,
    media_id: Uuid,
) -> Result<(), redis::RedisError> {
    let mut conn = redis_client.get_async_connection().await?;
    redis::cmd("DEL")
        .arg(format!("media:{}", media_id))
        .query_async(&mut conn)
        .await
}
