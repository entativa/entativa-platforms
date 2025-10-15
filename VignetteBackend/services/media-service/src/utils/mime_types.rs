use mime::Mime;
use mime_guess;
use std::path::Path;

use crate::models::MediaType;

pub fn detect_mime_type(filename: &str) -> Option<Mime> {
    mime_guess::from_path(filename).first()
}

pub fn detect_mime_from_bytes(data: &[u8]) -> Option<Mime> {
    // Try to detect from magic bytes
    if data.len() < 12 {
        return None;
    }

    match &data[0..8] {
        // PNG signature
        [0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A] => {
            Some(mime::IMAGE_PNG)
        }
        // JPEG signature
        [0xFF, 0xD8, 0xFF, _, _, _, _, _] => {
            Some(mime::IMAGE_JPEG)
        }
        // GIF signature
        [0x47, 0x49, 0x46, 0x38, 0x37, 0x61, _, _] | 
        [0x47, 0x49, 0x46, 0x38, 0x39, 0x61, _, _] => {
            Some(mime::IMAGE_GIF)
        }
        // WebP signature
        [0x52, 0x49, 0x46, 0x46, _, _, _, _] if &data[8..12] == b"WEBP" => {
            Some("image/webp".parse().ok()?)
        }
        // MP4 signature
        [_, _, _, _, 0x66, 0x74, 0x79, 0x70] => {
            Some("video/mp4".parse().ok()?)
        }
        // MP3 signature
        [0x49, 0x44, 0x33, _, _, _, _, _] | [0xFF, 0xFB, _, _, _, _, _, _] => {
            Some("audio/mpeg".parse().ok()?)
        }
        _ => None,
    }
}

pub fn get_media_type_from_mime(mime: &Mime) -> MediaType {
    match mime.type_().as_str() {
        "image" => MediaType::Image,
        "video" => MediaType::Video,
        "audio" => MediaType::Audio,
        _ => MediaType::Document,
    }
}

pub fn is_image(mime: &Mime) -> bool {
    mime.type_() == mime::IMAGE
}

pub fn is_video(mime: &Mime) -> bool {
    mime.type_() == mime::VIDEO
}

pub fn is_audio(mime: &Mime) -> bool {
    mime.type_() == mime::AUDIO
}

pub fn get_extension_from_mime(mime: &Mime) -> Option<&'static str> {
    match mime.as_ref() {
        "image/jpeg" => Some("jpg"),
        "image/png" => Some("png"),
        "image/gif" => Some("gif"),
        "image/webp" => Some("webp"),
        "video/mp4" => Some("mp4"),
        "video/quicktime" => Some("mov"),
        "video/webm" => Some("webm"),
        "audio/mpeg" => Some("mp3"),
        "audio/mp4" => Some("m4a"),
        "audio/wav" => Some("wav"),
        "audio/ogg" => Some("ogg"),
        _ => None,
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_detect_mime_from_filename() {
        assert_eq!(
            detect_mime_type("test.jpg"),
            Some(mime::IMAGE_JPEG)
        );
        assert_eq!(
            detect_mime_type("test.png"),
            Some(mime::IMAGE_PNG)
        );
    }

    #[test]
    fn test_detect_mime_from_png_bytes() {
        let png_signature = [0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x00];
        assert_eq!(detect_mime_from_bytes(&png_signature), Some(mime::IMAGE_PNG));
    }

    #[test]
    fn test_media_type_from_mime() {
        let jpeg = mime::IMAGE_JPEG;
        assert!(matches!(get_media_type_from_mime(&jpeg), MediaType::Image));
    }

    #[test]
    fn test_extension_from_mime() {
        assert_eq!(get_extension_from_mime(&mime::IMAGE_JPEG), Some("jpg"));
        assert_eq!(get_extension_from_mime(&mime::IMAGE_PNG), Some("png"));
    }
}
