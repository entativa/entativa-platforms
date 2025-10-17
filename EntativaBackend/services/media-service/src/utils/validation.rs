use crate::config::LimitsConfig;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum ValidationError {
    #[error("File too large: {0} bytes (max: {1} bytes)")]
    FileTooLarge(u64, u64),
    
    #[error("Invalid file type: {0}")]
    InvalidFileType(String),
    
    #[error("File extension mismatch: expected {0}, got {1}")]
    ExtensionMismatch(String, String),
    
    #[error("Invalid filename: {0}")]
    InvalidFilename(String),
    
    #[error("Validation failed: {0}")]
    General(String),
}

pub type ValidationResult<T> = Result<T, ValidationError>;

pub struct FileValidator {
    limits: LimitsConfig,
}

impl FileValidator {
    pub fn new(limits: LimitsConfig) -> Self {
        Self { limits }
    }

    pub fn validate_image(&self, size: u64, mime_type: &str) -> ValidationResult<()> {
        let max_size = self.limits.max_image_size_mb * 1024 * 1024;
        
        if size > max_size {
            return Err(ValidationError::FileTooLarge(size, max_size));
        }

        if !self.limits.allowed_image_types.contains(&mime_type.to_string()) {
            return Err(ValidationError::InvalidFileType(mime_type.to_string()));
        }

        Ok(())
    }

    pub fn validate_video(&self, size: u64, mime_type: &str) -> ValidationResult<()> {
        let max_size = self.limits.max_video_size_mb * 1024 * 1024;
        
        if size > max_size {
            return Err(ValidationError::FileTooLarge(size, max_size));
        }

        if !self.limits.allowed_video_types.contains(&mime_type.to_string()) {
            return Err(ValidationError::InvalidFileType(mime_type.to_string()));
        }

        Ok(())
    }

    pub fn validate_audio(&self, size: u64, mime_type: &str) -> ValidationResult<()> {
        let max_size = self.limits.max_audio_size_mb * 1024 * 1024;
        
        if size > max_size {
            return Err(ValidationError::FileTooLarge(size, max_size));
        }

        if !self.limits.allowed_audio_types.contains(&mime_type.to_string()) {
            return Err(ValidationError::InvalidFileType(mime_type.to_string()));
        }

        Ok(())
    }

    pub fn validate_filename(&self, filename: &str) -> ValidationResult<()> {
        if filename.is_empty() {
            return Err(ValidationError::InvalidFilename("Filename cannot be empty".to_string()));
        }

        if filename.len() > 255 {
            return Err(ValidationError::InvalidFilename("Filename too long (max 255 characters)".to_string()));
        }

        // Check for invalid characters
        let invalid_chars = ['/', '\\', '\0', '<', '>', ':', '"', '|', '?', '*'];
        if filename.chars().any(|c| invalid_chars.contains(&c)) {
            return Err(ValidationError::InvalidFilename("Filename contains invalid characters".to_string()));
        }

        Ok(())
    }

    pub fn sanitize_filename(&self, filename: &str) -> String {
        filename
            .chars()
            .map(|c| {
                if c.is_alphanumeric() || c == '.' || c == '-' || c == '_' {
                    c
                } else {
                    '_'
                }
            })
            .collect::<String>()
            .trim_matches('.')
            .to_string()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    fn get_test_limits() -> LimitsConfig {
        LimitsConfig {
            max_file_size_mb: 100,
            max_image_size_mb: 10,
            max_video_size_mb: 100,
            max_audio_size_mb: 50,
            allowed_image_types: vec!["image/jpeg".to_string(), "image/png".to_string()],
            allowed_video_types: vec!["video/mp4".to_string()],
            allowed_audio_types: vec!["audio/mpeg".to_string()],
            rate_limit_per_minute: 60,
        }
    }

    #[test]
    fn test_validate_image_success() {
        let validator = FileValidator::new(get_test_limits());
        assert!(validator.validate_image(5 * 1024 * 1024, "image/jpeg").is_ok());
    }

    #[test]
    fn test_validate_image_too_large() {
        let validator = FileValidator::new(get_test_limits());
        assert!(validator.validate_image(15 * 1024 * 1024, "image/jpeg").is_err());
    }

    #[test]
    fn test_sanitize_filename() {
        let validator = FileValidator::new(get_test_limits());
        assert_eq!(validator.sanitize_filename("test file!@#.jpg"), "test_file___.jpg");
    }
}
