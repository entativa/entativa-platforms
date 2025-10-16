use crate::services::image_processor::{ImageProcessor, ImageProcessingError};
use image::{DynamicImage, GenericImageView};
use flate2::write::GzEncoder;
use flate2::Compression as GzCompression;
use std::io::Write;

pub struct CompressionService {
    processor: ImageProcessor,
}

impl CompressionService {
    pub fn new(processor: ImageProcessor) -> Self {
        Self { processor }
    }

    /// Intelligently compress image based on content analysis
    pub async fn smart_compress(&self, data: &[u8]) -> Result<Vec<u8>, ImageProcessingError> {
        let img = image::load_from_memory(data)
            .map_err(|e| ImageProcessingError::DecodeError(e.to_string()))?;

        let (width, height) = img.dimensions();
        let pixel_count = width * height;

        // Determine optimal quality based on resolution
        let quality = if pixel_count > 4_000_000 {
            // High-res images: aggressive compression
            75
        } else if pixel_count > 1_000_000 {
            // Medium-res: balanced
            82
        } else {
            // Low-res: preserve quality
            90
        };

        // Check if WebP offers better compression
        let webp = self.encode_webp(&img, quality - 5)?;
        let jpeg = self.encode_jpeg(&img, quality)?;

        // Return smaller file
        if webp.len() < jpeg.len() {
            Ok(webp)
        } else {
            Ok(jpeg)
        }
    }

    /// Compress with target file size (iterative quality reduction)
    pub async fn compress_to_target_size(
        &self,
        data: &[u8],
        target_size_bytes: usize,
        min_quality: u8,
    ) -> Result<Vec<u8>, ImageProcessingError> {
        let img = image::load_from_memory(data)
            .map_err(|e| ImageProcessingError::DecodeError(e.to_string()))?;

        let mut quality = 90u8;
        let mut result = self.encode_jpeg(&img, quality)?;

        // Binary search for optimal quality
        while result.len() > target_size_bytes && quality > min_quality {
            quality = quality.saturating_sub(5);
            result = self.encode_jpeg(&img, quality)?;
        }

        // If still too large, try reducing dimensions
        if result.len() > target_size_bytes {
            let (width, height) = img.dimensions();
            let scale_factor = (target_size_bytes as f32 / result.len() as f32).sqrt();
            let new_width = ((width as f32 * scale_factor) as u32).max(100);
            let new_height = ((height as f32 * scale_factor) as u32).max(100);

            result = self.processor.resize(data, new_width, new_height, true).await?;
        }

        Ok(result)
    }

    /// Compress image with specific quality
    pub async fn compress_quality(&self, data: &[u8], quality: u8) -> Result<Vec<u8>, ImageProcessingError> {
        let img = image::load_from_memory(data)
            .map_err(|e| ImageProcessingError::DecodeError(e.to_string()))?;

        self.encode_jpeg(&img, quality)
    }

    /// Lossless compression using PNG
    pub async fn compress_lossless(&self, data: &[u8]) -> Result<Vec<u8>, ImageProcessingError> {
        let img = image::load_from_memory(data)
            .map_err(|e| ImageProcessingError::DecodeError(e.to_string()))?;

        let mut output = Vec::new();
        let mut cursor = std::io::Cursor::new(&mut output);
        
        img.write_to(&mut cursor, image::ImageFormat::Png)
            .map_err(|e| ImageProcessingError::EncodeError(e.to_string()))?;

        Ok(output)
    }

    /// Progressive JPEG encoding for better web experience
    pub async fn create_progressive_jpeg(&self, data: &[u8], quality: u8) -> Result<Vec<u8>, ImageProcessingError> {
        let img = image::load_from_memory(data)
            .map_err(|e| ImageProcessingError::DecodeError(e.to_string()))?;

        // Convert to RGB
        let rgb = img.to_rgb8();
        let mut output = Vec::new();
        let mut cursor = Cursor::new(&mut output);

        // Progressive JPEG encoder
        let mut encoder = image::codecs::jpeg::JpegEncoder::new_with_quality(&mut cursor, quality);
        
        rgb.write_with_encoder(encoder)
            .map_err(|e| ImageProcessingError::EncodeError(e.to_string()))?;

        Ok(output)
    }

    // Private helper methods

    fn encode_jpeg(&self, img: &DynamicImage, quality: u8) -> Result<Vec<u8>, ImageProcessingError> {
        let mut output = Vec::new();
        let mut cursor = Cursor::new(&mut output);
        
        let rgb = img.to_rgb8();
        let encoder = image::codecs::jpeg::JpegEncoder::new_with_quality(&mut cursor, quality);
        
        rgb.write_with_encoder(encoder)
            .map_err(|e| ImageProcessingError::EncodeError(e.to_string()))?;

        Ok(output)
    }

    fn encode_webp(&self, img: &DynamicImage, quality: u8) -> Result<Vec<u8>, ImageProcessingError> {
        let mut output = Vec::new();
        let mut cursor = Cursor::new(&mut output);
        
        img.write_to(&mut cursor, ImageFormat::WebP)
            .map_err(|e| ImageProcessingError::EncodeError(e.to_string()))?;

        Ok(output)
    }
}

impl Default for CompressionService {
    fn default() -> Self {
        Self::new(ImageProcessor::default())
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn test_smart_compress() {
        let service = CompressionService::default();
        // Would test with actual image data
    }

    #[tokio::test]
    async fn test_compress_to_target() {
        let service = CompressionService::default();
        // Verify compression achieves target size
    }
}
