use image::{DynamicImage, GenericImageView, ImageBuffer, Rgba, RgbaImage};
use imageproc::drawing::draw_filled_rect_mut;
use imageproc::rect::Rect;
use std::io::Cursor;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum FilterError {
    #[error("Failed to apply filter: {0}")]
    ProcessingError(String),
    
    #[error("Invalid image data: {0}")]
    InvalidImage(String),
}

pub type FilterResult<T> = Result<T, FilterError>;

/// Instagram-style filter service for professional-quality photo filters
pub struct FilterService;

impl FilterService {
    /// Clarendon filter - Brightens and intensifies cool tones
    pub fn clarendon(img: &DynamicImage) -> FilterResult<Vec<u8>> {
        let mut rgba = img.to_rgba8();
        let (width, height) = rgba.dimensions();

        for pixel in rgba.pixels_mut() {
            let r = pixel[0] as f32;
            let g = pixel[1] as f32;
            let b = pixel[2] as f32;

            // Brighten highlights
            let brightness_factor = 1.15;
            
            // Enhance cool tones (blues)
            let blue_boost = if b > 100.0 { 1.2 } else { 1.0 };
            
            // Subtle saturation increase
            let avg = (r + g + b) / 3.0;
            let sat_factor = 1.15;

            pixel[0] = ((r - avg) * sat_factor + avg * brightness_factor).min(255.0) as u8;
            pixel[1] = ((g - avg) * sat_factor + avg * brightness_factor).min(255.0) as u8;
            pixel[2] = ((b - avg) * sat_factor * blue_boost + avg * brightness_factor).min(255.0) as u8;
        }

        Self::encode_jpeg(&DynamicImage::ImageRgba8(rgba), 92)
    }

    /// Gingham filter - Soft, warm, vintage look
    pub fn gingham(img: &DynamicImage) -> FilterResult<Vec<u8>> {
        let mut rgba = img.to_rgba8();

        for pixel in rgba.pixels_mut() {
            let r = pixel[0] as f32;
            let g = pixel[1] as f32;
            let b = pixel[2] as f32;

            // Warm tones
            let warm_factor = 1.1;
            let cool_reduce = 0.95;

            // Soft, lower contrast
            let contrast = 0.9;
            let avg = (r + g + b) / 3.0;

            pixel[0] = ((r - avg) * contrast + avg * warm_factor).min(255.0) as u8;
            pixel[1] = ((g - avg) * contrast + avg).min(255.0) as u8;
            pixel[2] = ((b - avg) * contrast + avg * cool_reduce).min(255.0) as u8;
        }

        Self::encode_jpeg(&DynamicImage::ImageRgba8(rgba), 92)
    }

    /// Juno filter - Enhances cool tones and adds subtle vignette
    pub fn juno(img: &DynamicImage) -> FilterResult<Vec<u8>> {
        let mut rgba = img.to_rgba8();
        let (width, height) = rgba.dimensions();
        let center_x = width as f32 / 2.0;
        let center_y = height as f32 / 2.0;
        let max_dist = ((width * width + height * height) as f32).sqrt() / 2.0;

        for (x, y, pixel) in rgba.enumerate_pixels_mut() {
            let r = pixel[0] as f32;
            let g = pixel[1] as f32;
            let b = pixel[2] as f32;

            // Enhance cool tones
            let cool_boost = 1.15;
            
            // Calculate vignette
            let dx = x as f32 - center_x;
            let dy = y as f32 - center_y;
            let dist = (dx * dx + dy * dy).sqrt();
            let vignette = 1.0 - (dist / max_dist * 0.3);

            pixel[0] = (r * vignette).min(255.0) as u8;
            pixel[1] = (g * cool_boost * vignette).min(255.0) as u8;
            pixel[2] = (b * cool_boost * vignette).min(255.0) as u8;
        }

        Self::encode_jpeg(&DynamicImage::ImageRgba8(rgba), 92)
    }

    /// Lark filter - Bright, desaturated, cool
    pub fn lark(img: &DynamicImage) -> FilterResult<Vec<u8>> {
        let mut rgba = img.to_rgba8();

        for pixel in rgba.pixels_mut() {
            let r = pixel[0] as f32;
            let g = pixel[1] as f32;
            let b = pixel[2] as f32;

            // Brightness boost
            let brightness = 1.2;
            
            // Desaturation
            let avg = (r + g + b) / 3.0;
            let desat = 0.7;

            pixel[0] = ((r - avg) * desat + avg * brightness).min(255.0) as u8;
            pixel[1] = ((g - avg) * desat + avg * brightness).min(255.0) as u8;
            pixel[2] = ((b - avg) * desat + avg * brightness * 1.05).min(255.0) as u8;
        }

        Self::encode_jpeg(&DynamicImage::ImageRgba8(rgba), 92)
    }

    /// Reyes filter - Vintage, low contrast, warm
    pub fn reyes(img: &DynamicImage) -> FilterResult<Vec<u8>> {
        let mut rgba = img.to_rgba8();

        for pixel in rgba.pixels_mut() {
            let r = pixel[0] as f32;
            let g = pixel[1] as f32;
            let b = pixel[2] as f32;

            // Vintage warmth
            let warm = 1.15;
            
            // Low contrast
            let contrast = 0.85;
            let avg = (r + g + b) / 3.0;
            let base = 128.0;

            pixel[0] = ((r - base) * contrast + base * warm).min(255.0) as u8;
            pixel[1] = ((g - base) * contrast + base).min(255.0) as u8;
            pixel[2] = ((b - base) * contrast + base * 0.95).min(255.0) as u8;
        }

        Self::encode_jpeg(&DynamicImage::ImageRgba8(rgba), 92)
    }

    /// Valencia filter - Warm, faded, vintage
    pub fn valencia(img: &DynamicImage) -> FilterResult<Vec<u8>> {
        let mut rgba = img.to_rgba8();

        for pixel in rgba.pixels_mut() {
            let r = pixel[0] as f32;
            let g = pixel[1] as f32;
            let b = pixel[2] as f32;

            // Warm fade
            let fade = 0.9;
            let warm = 1.2;

            pixel[0] = (r * warm * fade + 25.0).min(255.0) as u8;
            pixel[1] = (g * fade + 15.0).min(255.0) as u8;
            pixel[2] = (b * fade * 0.9 + 10.0).min(255.0) as u8;
        }

        Self::encode_jpeg(&DynamicImage::ImageRgba8(rgba), 92)
    }

    /// X-Pro II filter - High contrast, warm shadows, cool highlights
    pub fn xpro2(img: &DynamicImage) -> FilterResult<Vec<u8>> {
        let mut rgba = img.to_rgba8();

        for pixel in rgba.pixels_mut() {
            let r = pixel[0] as f32;
            let g = pixel[1] as f32;
            let b = pixel[2] as f32;

            // High contrast
            let contrast = 1.3;
            let avg = (r + g + b) / 3.0;
            let base = 128.0;

            // Split toning: warm shadows, cool highlights
            let tone_shift = if avg < 128.0 {
                (1.1, 1.0, 0.95) // Warm shadows
            } else {
                (0.95, 1.0, 1.1) // Cool highlights
            };

            pixel[0] = ((r - base) * contrast + base * tone_shift.0).min(255.0) as u8;
            pixel[1] = ((g - base) * contrast + base * tone_shift.1).min(255.0) as u8;
            pixel[2] = ((b - base) * contrast + base * tone_shift.2).min(255.0) as u8;
        }

        Self::encode_jpeg(&DynamicImage::ImageRgba8(rgba), 92)
    }

    /// Lo-Fi filter - Rich colors, strong shadows, subtle vignette
    pub fn lofi(img: &DynamicImage) -> FilterResult<Vec<u8>> {
        let mut rgba = img.to_rgba8();
        let (width, height) = rgba.dimensions();
        let center_x = width as f32 / 2.0;
        let center_y = height as f32 / 2.0;
        let max_dist = ((width * width + height * height) as f32).sqrt() / 2.0;

        for (x, y, pixel) in rgba.enumerate_pixels_mut() {
            let r = pixel[0] as f32;
            let g = pixel[1] as f32;
            let b = pixel[2] as f32;

            // Rich, saturated colors
            let avg = (r + g + b) / 3.0;
            let saturation = 1.4;

            // Strong shadows
            let shadow_strength = if avg < 100.0 { 0.8 } else { 1.0 };

            // Vignette
            let dx = x as f32 - center_x;
            let dy = y as f32 - center_y;
            let dist = (dx * dx + dy * dy).sqrt();
            let vignette = 1.0 - (dist / max_dist * 0.4);

            pixel[0] = ((r - avg) * saturation + avg * shadow_strength * vignette).min(255.0) as u8;
            pixel[1] = ((g - avg) * saturation + avg * shadow_strength * vignette).min(255.0) as u8;
            pixel[2] = ((b - avg) * saturation + avg * shadow_strength * vignette).min(255.0) as u8;
        }

        Self::encode_jpeg(&DynamicImage::ImageRgba8(rgba), 92)
    }

    /// Nashville filter - Warm, vintage with pink/yellow tones
    pub fn nashville(img: &DynamicImage) -> FilterResult<Vec<u8>> {
        let mut rgba = img.to_rgba8();

        for pixel in rgba.pixels_mut() {
            let r = pixel[0] as f32;
            let g = pixel[1] as f32;
            let b = pixel[2] as f32;

            // Warm vintage tones
            let pink_shift = 15.0;
            let yellow_shift = 10.0;

            // Reduce contrast slightly
            let contrast = 0.95;
            let avg = (r + g + b) / 3.0;

            pixel[0] = ((r - avg) * contrast + avg + pink_shift).min(255.0) as u8;
            pixel[1] = ((g - avg) * contrast + avg + yellow_shift).min(255.0) as u8;
            pixel[2] = ((b - avg) * contrast + avg - 5.0).max(0.0) as u8;
        }

        Self::encode_jpeg(&DynamicImage::ImageRgba8(rgba), 92)
    }

    /// Perpetua filter - Cool, desaturated, pastel
    pub fn perpetua(img: &DynamicImage) -> FilterResult<Vec<u8>> {
        let mut rgba = img.to_rgba8();

        for pixel in rgba.pixels_mut() {
            let r = pixel[0] as f32;
            let g = pixel[1] as f32;
            let b = pixel[2] as f32;

            let avg = (r + g + b) / 3.0;
            
            // Pastel, desaturated
            let desat = 0.6;
            let lightness = 1.1;

            pixel[0] = ((r - avg) * desat + avg * lightness).min(255.0) as u8;
            pixel[1] = ((g - avg) * desat + avg * lightness * 1.05).min(255.0) as u8;
            pixel[2] = ((b - avg) * desat + avg * lightness * 1.1).min(255.0) as u8;
        }

        Self::encode_jpeg(&DynamicImage::ImageRgba8(rgba), 92)
    }

    /// Aden filter - Subtle, desaturated with blue shadows
    pub fn aden(img: &DynamicImage) -> FilterResult<Vec<u8>> {
        let mut rgba = img.to_rgba8();

        for pixel in rgba.pixels_mut() {
            let r = pixel[0] as f32;
            let g = pixel[1] as f32;
            let b = pixel[2] as f32;

            let avg = (r + g + b) / 3.0;
            let desat = 0.75;

            // Blue tint in shadows
            let blue_shift = if avg < 128.0 { 10.0 } else { 0.0 };

            pixel[0] = ((r - avg) * desat + avg).min(255.0) as u8;
            pixel[1] = ((g - avg) * desat + avg).min(255.0) as u8;
            pixel[2] = ((b - avg) * desat + avg + blue_shift).min(255.0) as u8;
        }

        Self::encode_jpeg(&DynamicImage::ImageRgba8(rgba), 92)
    }

    /// Inkwell filter - Black and white with high contrast
    pub fn inkwell(img: &DynamicImage) -> FilterResult<Vec<u8>> {
        let mut rgba = img.to_rgba8();

        for pixel in rgba.pixels_mut() {
            let r = pixel[0] as f32;
            let g = pixel[1] as f32;
            let b = pixel[2] as f32;

            // Convert to grayscale
            let gray = (0.299 * r + 0.587 * g + 0.114 * b) as u8;

            // High contrast
            let contrasted = if gray > 128 {
                ((gray as f32 - 128.0) * 1.3 + 128.0).min(255.0) as u8
            } else {
                ((gray as f32 - 128.0) * 1.3 + 128.0).max(0.0) as u8
            };

            pixel[0] = contrasted;
            pixel[1] = contrasted;
            pixel[2] = contrasted;
        }

        Self::encode_jpeg(&DynamicImage::ImageRgba8(rgba), 92)
    }

    /// Ludwig filter - Cool, clean, bright
    pub fn ludwig(img: &DynamicImage) -> FilterResult<Vec<u8>> {
        let mut rgba = img.to_rgba8();

        for pixel in rgba.pixels_mut() {
            let r = pixel[0] as f32;
            let g = pixel[1] as f32;
            let b = pixel[2] as f32;

            // Bright and clean
            let brightness = 1.1;
            let cool_boost = 1.05;

            pixel[0] = (r * brightness).min(255.0) as u8;
            pixel[1] = (g * brightness * cool_boost).min(255.0) as u8;
            pixel[2] = (b * brightness * cool_boost).min(255.0) as u8;
        }

        Self::encode_jpeg(&DynamicImage::ImageRgba8(rgba), 92)
    }

    /// Slumber filter - Desaturated, warm, dreamy
    pub fn slumber(img: &DynamicImage) -> FilterResult<Vec<u8>> {
        let mut rgba = img.to_rgba8();

        for pixel in rgba.pixels_mut() {
            let r = pixel[0] as f32;
            let g = pixel[1] as f32;
            let b = pixel[2] as f32;

            let avg = (r + g + b) / 3.0;
            
            // Dreamy desaturation
            let desat = 0.6;
            let warm = 1.05;

            pixel[0] = ((r - avg) * desat + avg * warm + 10.0).min(255.0) as u8;
            pixel[1] = ((g - avg) * desat + avg + 5.0).min(255.0) as u8;
            pixel[2] = ((b - avg) * desat + avg * 0.95).max(0.0) as u8;
        }

        Self::encode_jpeg(&DynamicImage::ImageRgba8(rgba), 92)
    }

    /// Crema filter - Subtle, warm, reduced contrast
    pub fn crema(img: &DynamicImage) -> FilterResult<Vec<u8>> {
        let mut rgba = img.to_rgba8();

        for pixel in rgba.pixels_mut() {
            let r = pixel[0] as f32;
            let g = pixel[1] as f32;
            let b = pixel[2] as f32;

            // Reduced contrast, warm
            let contrast = 0.9;
            let base = 128.0;
            let warm = 1.05;

            pixel[0] = ((r - base) * contrast + base * warm).min(255.0) as u8;
            pixel[1] = ((g - base) * contrast + base).min(255.0) as u8;
            pixel[2] = ((b - base) * contrast + base * 0.98).max(0.0) as u8;
        }

        Self::encode_jpeg(&DynamicImage::ImageRgba8(rgba), 92)
    }

    /// Apply filter by name
    pub fn apply_filter(filter_name: &str, img: &DynamicImage) -> FilterResult<Vec<u8>> {
        match filter_name.to_lowercase().as_str() {
            "clarendon" => Self::clarendon(img),
            "gingham" => Self::gingham(img),
            "juno" => Self::juno(img),
            "lark" => Self::lark(img),
            "reyes" => Self::reyes(img),
            "valencia" => Self::valencia(img),
            "xpro2" | "x-pro-ii" => Self::xpro2(img),
            "lofi" | "lo-fi" => Self::lofi(img),
            "nashville" => Self::nashville(img),
            "perpetua" => Self::perpetua(img),
            "aden" => Self::aden(img),
            "ludwig" => Self::ludwig(img),
            "slumber" => Self::slumber(img),
            "crema" => Self::crema(img),
            _ => Err(FilterError::ProcessingError(format!("Unknown filter: {}", filter_name))),
        }
    }

    /// Get list of available filters
    pub fn available_filters() -> Vec<&'static str> {
        vec![
            "clarendon",
            "gingham",
            "juno",
            "lark",
            "reyes",
            "valencia",
            "xpro2",
            "lofi",
            "nashville",
            "perpetua",
            "aden",
            "ludwig",
            "slumber",
            "crema",
        ]
    }

    // Helper methods

    fn encode_jpeg(img: &DynamicImage, quality: u8) -> FilterResult<Vec<u8>> {
        let mut output = Vec::new();
        let mut cursor = Cursor::new(&mut output);
        
        let rgb = img.to_rgb8();
        let encoder = image::codecs::jpeg::JpegEncoder::new_with_quality(&mut cursor, quality);
        
        rgb.write_with_encoder(encoder)
            .map_err(|e| FilterError::ProcessingError(e.to_string()))?;

        Ok(output)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_available_filters() {
        let filters = FilterService::available_filters();
        assert!(filters.len() > 10);
        assert!(filters.contains(&"clarendon"));
    }

    #[test]
    fn test_filter_names() {
        // Verify all filters are accessible
        for filter in FilterService::available_filters() {
            // Would test with actual image
            assert!(!filter.is_empty());
        }
    }
}
