use crate::models::{Media, MediaVariant, ProcessingStatus};
use crate::services::{VideoProcessor, AudioProcessor};
use sqlx::PgPool;
use std::sync::Arc;
use tokio::sync::Semaphore;
use uuid::Uuid;

pub struct TranscodingService {
    video_processor: VideoProcessor,
    audio_processor: AudioProcessor,
    max_concurrent: usize,
}

impl TranscodingService {
    pub fn new(
        video_processor: VideoProcessor,
        audio_processor: AudioProcessor,
        max_concurrent: usize,
    ) -> Self {
        Self {
            video_processor,
            audio_processor,
            max_concurrent,
        }
    }

    /// Start background transcoding worker
    pub async fn start_worker(
        &self,
        pool: PgPool,
        storage: Arc<dyn crate::storage::StorageBackend>,
    ) {
        let semaphore = Arc::new(Semaphore::new(self.max_concurrent));

        loop {
            // Get pending media
            let pending_media = self.get_pending_media(&pool, 10).await;

            for media in pending_media {
                let permit = semaphore.clone().acquire_owned().await.unwrap();
                let pool_clone = pool.clone();
                let storage_clone = storage.clone();
                let video_proc = self.video_processor.clone();
                let audio_proc = self.audio_processor.clone();

                tokio::spawn(async move {
                    let _permit = permit; // Hold permit until complete

                    if media.is_video() {
                        Self::process_video(media, pool_clone, storage_clone, video_proc).await.ok();
                    } else if media.is_audio() {
                        Self::process_audio(media, pool_clone, storage_clone, audio_proc).await.ok();
                    }
                });
            }

            // Sleep before next check
            tokio::time::sleep(tokio::time::Duration::from_secs(5)).await;
        }
    }

    async fn get_pending_media(&self, pool: &PgPool, limit: i32) -> Vec<Media> {
        sqlx::query_as::<_, Media>(
            "SELECT * FROM media WHERE processing_status = $1 AND is_deleted = false ORDER BY created_at ASC LIMIT $2"
        )
        .bind(ProcessingStatus::Pending)
        .bind(limit)
        .fetch_all(pool)
        .await
        .unwrap_or_default()
    }

    async fn process_video(
        media: Media,
        pool: PgPool,
        storage: Arc<dyn crate::storage::StorageBackend>,
        processor: VideoProcessor,
    ) -> Result<(), Box<dyn std::error::Error>> {
        tracing::info!("Processing video: {}", media.id);

        // Update status
        sqlx::query("UPDATE media SET processing_status = $1 WHERE id = $2")
            .bind(ProcessingStatus::Processing)
            .bind(media.id)
            .execute(&pool)
            .await?;

        // Download original
        let original_data = storage.download(&media.storage_path).await?;

        // Save to temp file for FFmpeg
        let temp_input = format!("/tmp/{}_input", media.id);
        tokio::fs::write(&temp_input, &original_data).await?;

        // Transcode to H.264
        let temp_output = format!("/tmp/{}_output.mp4", media.id);
        processor.transcode_to_h264(&temp_input, &temp_output, Some((1280, 720))).await.ok();

        // Upload processed video
        let processed_data = tokio::fs::read(&temp_output).await?;
        let processed_path = format!("{}.h264.mp4", media.storage_path);
        storage.upload(&processed_path, bytes::Bytes::from(processed_data), "video/mp4").await?;

        // Extract thumbnail
        let thumb_path = format!("/tmp/{}_thumb.jpg", media.id);
        if processor.extract_frame(&temp_input, 1.0, &thumb_path).await.is_ok() {
            let thumb_data = tokio::fs::read(&thumb_path).await?;
            let thumb_storage_path = format!("{}/thumbnail.jpg", media.id);
            storage.upload(&thumb_storage_path, bytes::Bytes::from(thumb_data), "image/jpeg").await.ok();
        }

        // Clean up temp files
        tokio::fs::remove_file(&temp_input).await.ok();
        tokio::fs::remove_file(&temp_output).await.ok();

        // Update status to completed
        sqlx::query("UPDATE media SET processing_status = $1, is_processed = $2 WHERE id = $3")
            .bind(ProcessingStatus::Completed)
            .bind(true)
            .bind(media.id)
            .execute(&pool)
            .await?;

        tracing::info!("Video processing completed: {}", media.id);

        Ok(())
    }

    async fn process_audio(
        media: Media,
        pool: PgPool,
        storage: Arc<dyn crate::storage::StorageBackend>,
        processor: AudioProcessor,
    ) -> Result<(), Box<dyn std::error::Error>> {
        tracing::info!("Processing audio: {}", media.id);

        // Update status
        sqlx::query("UPDATE media SET processing_status = $1 WHERE id = $2")
            .bind(ProcessingStatus::Processing)
            .bind(media.id)
            .execute(&pool)
            .await?;

        // Download original
        let original_data = storage.download(&media.storage_path).await?;

        // Save to temp file
        let temp_input = format!("/tmp/{}_audio_input", media.id);
        tokio::fs::write(&temp_input, &original_data).await?;

        // Transcode to AAC
        let temp_output = format!("/tmp/{}_audio_output.m4a", media.id);
        processor.transcode_to_aac(&temp_input, &temp_output).await.ok();

        // Upload processed audio
        if tokio::fs::metadata(&temp_output).await.is_ok() {
            let processed_data = tokio::fs::read(&temp_output).await?;
            let processed_path = format!("{}.aac.m4a", media.storage_path);
            storage.upload(&processed_path, bytes::Bytes::from(processed_data), "audio/mp4").await?;
        }

        // Clean up
        tokio::fs::remove_file(&temp_input).await.ok();
        tokio::fs::remove_file(&temp_output).await.ok();

        // Update status
        sqlx::query("UPDATE media SET processing_status = $1, is_processed = $2 WHERE id = $3")
            .bind(ProcessingStatus::Completed)
            .bind(true)
            .bind(media.id)
            .execute(&pool)
            .await?;

        tracing::info!("Audio processing completed: {}", media.id);

        Ok(())
    }
}

impl Default for TranscodingService {
    fn default() -> Self {
        Self::new(
            VideoProcessor::default(),
            AudioProcessor::default(),
            num_cpus::get().max(2),
        )
    }
}

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
