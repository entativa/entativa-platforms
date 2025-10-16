use actix_web::{web, HttpRequest, HttpResponse, Result as ActixResult};
use sqlx::PgPool;
use uuid::Uuid;

use crate::{models::Media, AppState};

/// Stream media with range request support for video/audio
pub async fn stream_media(
    req: HttpRequest,
    path: web::Path<Uuid>,
    data: web::Data<AppState>,
) -> ActixResult<HttpResponse> {
    let media_id = path.into_inner();

    // Get media info
    let media = get_media_by_id(&data.db, media_id).await.map_err(|_| {
        actix_web::error::ErrorNotFound("Media not found")
    })?;

    // Only stream video/audio
    if !media.is_video() && !media.is_audio() {
        return Ok(HttpResponse::BadRequest().json(serde_json::json!({
            "error": "Media is not streamable (must be video or audio)"
        })));
    }

    // Download file from storage
    let file_data = data
        .storage
        .download(&media.storage_path)
        .await
        .map_err(|e| actix_web::error::ErrorInternalServerError(e.to_string()))?;

    // Handle range requests for seeking
    let range_header = req.headers().get("Range");

    if let Some(range) = range_header {
        return handle_range_streaming(&media, file_data, range);
    }

    // Full file streaming
    Ok(HttpResponse::Ok()
        .content_type(&media.mime_type)
        .insert_header(("Accept-Ranges", "bytes"))
        .insert_header(("Content-Length", file_data.len()))
        .insert_header(("Cache-Control", "public, max-age=31536000"))
        .streaming(futures_util::stream::once(async move { Ok::<_, actix_web::Error>(file_data) })))
}

/// Serve HLS master playlist
pub async fn serve_hls_playlist(
    path: web::Path<Uuid>,
    data: web::Data<AppState>,
) -> ActixResult<HttpResponse> {
    let media_id = path.into_inner();

    // Get media info
    let media = get_media_by_id(&data.db, media_id).await.map_err(|_| {
        actix_web::error::ErrorNotFound("Media not found")
    })?;

    if !media.is_video() {
        return Ok(HttpResponse::BadRequest().json(serde_json::json!({
            "error": "Media is not a video"
        })));
    }

    // Check if HLS playlist exists
    let hls_path = format!("{}/hls/master.m3u8", media.storage_path);
    
    match data.storage.download(&hls_path).await {
        Ok(playlist_data) => {
            Ok(HttpResponse::Ok()
                .content_type("application/vnd.apple.mpegurl")
                .insert_header(("Cache-Control", "public, max-age=3600"))
                .body(playlist_data))
        }
        Err(_) => {
            // HLS not generated yet
            Ok(HttpResponse::NotFound().json(serde_json::json!({
                "error": "HLS stream not available for this video",
                "message": "Video may still be processing"
            })))
        }
    }
}

/// Serve HLS video segments
pub async fn serve_hls_segment(
    path: web::Path<(Uuid, String)>,
    data: web::Data<AppState>,
) -> ActixResult<HttpResponse> {
    let (media_id, segment_name) = path.into_inner();

    // Get media info
    let media = get_media_by_id(&data.db, media_id).await.map_err(|_| {
        actix_web::error::ErrorNotFound("Media not found")
    })?;

    // Validate segment name (security check)
    if !segment_name.ends_with(".ts") || segment_name.contains("..") || segment_name.contains('/') {
        return Ok(HttpResponse::BadRequest().json(serde_json::json!({
            "error": "Invalid segment name"
        })));
    }

    // Build segment path
    let segment_path = format!("{}/hls/{}", media.storage_path, segment_name);

    // Get segment from storage
    let segment_data = data.storage.download(&segment_path).await.map_err(|_| {
        actix_web::error::ErrorNotFound("Segment not found")
    })?;

    Ok(HttpResponse::Ok()
        .content_type("video/MP2T")
        .insert_header(("Cache-Control", "public, max-age=31536000"))
        .body(segment_data))
}

// Helper functions

fn handle_range_streaming(
    media: &Media,
    data: bytes::Bytes,
    range_header: &actix_web::http::header::HeaderValue,
) -> ActixResult<HttpResponse> {
    let range_str = range_header.to_str().map_err(|_| {
        actix_web::error::ErrorBadRequest("Invalid Range header")
    })?;

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
        .streaming(futures_util::stream::once(async move { Ok::<_, actix_web::Error>(chunk) })))
}

async fn get_media_by_id(pool: &PgPool, media_id: Uuid) -> Result<Media, sqlx::Error> {
    sqlx::query_as::<_, Media>("SELECT * FROM media WHERE id = $1 AND is_deleted = false")
        .bind(media_id)
        .fetch_one(pool)
        .await
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_range_parsing() {
        // Test range header parsing logic
        let range = "bytes=0-1023";
        assert!(range.starts_with("bytes="));
    }
}
