use crate::models::MediaVariant;
use crate::services::image_processor::{ImageProcessor, ImageProcessingError};
use image::{DynamicImage, GenericImageView, ImageFormat};
use std::io::Cursor;

pub struct ThumbnailGenerator {
    processor: ImageProcessor,
}

impl ThumbnailGenerator {
    pub fn new(processor: ImageProcessor) -> Self {
        Self { processor }
    }

    /// Generate multiple thumbnail sizes with highest quality
    pub async fn generate_thumbnails(
        &self,
        data: &[u8],
        sizes: &[(u32, u32, &str)], // (width, height, name)
    ) -> Result<Vec<ThumbnailResult>, ImageProcessingError> {
        let mut thumbnails = Vec::new();

        for (width, height, name) in sizes {
            match self.generate_single_thumbnail(data, *width, *height, name).await {
                Ok(thumb) => thumbnails.push(thumb),
                Err(e) => {
                    tracing::warn!("Failed to generate thumbnail {}: {}", name, e);
                    continue;
                }
            }
        }

        Ok(thumbnails)
    }

    /// Generate single thumbnail with smart cropping
    async fn generate_single_thumbnail(
        &self,
        data: &[u8],
        width: u32,
        height: u32,
        name: &str,
    ) -> Result<ThumbnailResult, ImageProcessingError> {
        // Use smart crop for square thumbnails
        let thumbnail_data = if width == height {
            self.processor.smart_crop(data, width, height).await?
        } else {
            self.processor.resize(data, width, height, true).await?
        };

        Ok(ThumbnailResult {
            name: name.to_string(),
            width,
            height,
            data: thumbnail_data,
            file_size: 0, // Will be set after upload
            format: "jpeg".to_string(),
        })
    }

    /// Generate video thumbnail from first frame
    pub async fn generate_video_thumbnail(
        &self,
        video_data: &[u8],
        timestamp_secs: f64,
        width: u32,
        height: u32,
    ) -> Result<Vec<u8>, ImageProcessingError> {
        // In production, use ffmpeg to extract frame at timestamp
        // For now, return error indicating implementation needed
        Err(ImageProcessingError::ProcessingError(
            "Video thumbnail generation requires ffmpeg integration".to_string()
        ))
    }

    /// Generate progressive thumbnails (blurred placeholder)
    pub async fn generate_progressive_placeholder(
        &self,
        data: &[u8],
        width: u32,
    ) -> Result<Vec<u8>, ImageProcessingError> {
        // Create tiny blurred version for progressive loading
        let tiny = self.processor.resize(data, width, width, true).await?;
        self.processor.blur(&tiny, 10.0).await
    }

    /// Generate blurhash for lazy loading placeholders
    pub fn generate_blurhash(&self, img: &DynamicImage) -> Option<String> {
        // In production, integrate blurhash-rs crate
        // Returns base83 encoded string
        None
    }

    /// Create responsive image set (srcset)
    pub async fn generate_responsive_set(
        &self,
        data: &[u8],
        base_width: u32,
    ) -> Result<Vec<ThumbnailResult>, ImageProcessingError> {
        let sizes = vec![
            (base_width / 4, base_width / 4, "xs"),
            (base_width / 2, base_width / 2, "sm"),
            (base_width, base_width, "md"),
            (base_width * 2, base_width * 2, "lg"),
        ];

        self.generate_thumbnails(data, &sizes).await
    }
}

#[derive(Debug, Clone)]
pub struct ThumbnailResult {
    pub name: String,
    pub width: u32,
    pub height: u32,
    pub data: Vec<u8>,
    pub file_size: i64,
    pub format: String,
}

impl ThumbnailResult {
    pub fn to_media_variant(&self, url: String) -> MediaVariant {
        MediaVariant {
            variant_type: self.name.clone(),
            url,
            width: Some(self.width),
            height: Some(self.height),
            file_size: Some(self.file_size),
            format: self.format.clone(),
            bitrate: None,
        }
    }
}

impl Default for ThumbnailGenerator {
    fn default() -> Self {
        Self::new(ImageProcessor::default())
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn test_thumbnail_generation() {
        let generator = ThumbnailGenerator::default();
        // Would test with actual image data
    }
}
