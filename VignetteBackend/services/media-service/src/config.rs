use config::{Config as ConfigLoader, ConfigError, Environment, File};
use serde::Deserialize;
use std::env;

#[derive(Debug, Clone, Deserialize)]
pub struct Config {
    pub server: ServerConfig,
    pub storage: StorageConfig,
    pub database: DatabaseConfig,
    pub redis: RedisConfig,
    pub processing: ProcessingConfig,
    pub limits: LimitsConfig,
    pub cdn: CdnConfig,
}

#[derive(Debug, Clone, Deserialize)]
pub struct ServerConfig {
    pub host: String,
    pub port: u16,
    #[serde(default = "default_grpc_port")]
    pub grpc_port: u16,
    pub workers: usize,
    pub max_connections: usize,
    pub request_timeout_secs: u64,
}

fn default_grpc_port() -> u16 {
    50051
}

#[derive(Debug, Clone, Deserialize)]
pub struct StorageConfig {
    pub provider: String, // "s3", "minio", "local"
    pub s3: S3Config,
    pub minio: MinioConfig,
    pub local: LocalConfig,
    pub bucket_name: String,
    pub region: String,
}

#[derive(Debug, Clone, Deserialize)]
pub struct S3Config {
    pub access_key_id: String,
    pub secret_access_key: String,
    pub region: String,
    pub endpoint: Option<String>,
    pub use_path_style: bool,
}

#[derive(Debug, Clone, Deserialize)]
pub struct MinioConfig {
    pub endpoint: String,
    pub access_key: String,
    pub secret_key: String,
    pub use_ssl: bool,
}

#[derive(Debug, Clone, Deserialize)]
pub struct LocalConfig {
    pub base_path: String,
    pub serve_url: String,
}

#[derive(Debug, Clone, Deserialize)]
pub struct DatabaseConfig {
    pub url: String,
    pub max_connections: u32,
    pub min_connections: u32,
    pub connect_timeout_secs: u64,
    pub idle_timeout_secs: u64,
}

#[derive(Debug, Clone, Deserialize)]
pub struct RedisConfig {
    pub url: String,
    pub pool_size: u32,
    pub connection_timeout_secs: u64,
    pub cache_ttl_secs: u64,
}

#[derive(Debug, Clone, Deserialize)]
pub struct ProcessingConfig {
    pub image: ImageProcessingConfig,
    pub video: VideoProcessingConfig,
    pub audio: AudioProcessingConfig,
    pub thumbnail: ThumbnailConfig,
    pub compression: CompressionConfig,
}

#[derive(Debug, Clone, Deserialize)]
pub struct ImageProcessingConfig {
    pub max_dimension: u32,
    pub quality: u8,
    pub formats: Vec<String>,
    pub webp_quality: u8,
    pub thumbnail_sizes: Vec<u32>,
    pub enable_watermark: bool,
    pub watermark_path: Option<String>,
}

#[derive(Debug, Clone, Deserialize)]
pub struct VideoProcessingConfig {
    pub max_duration_secs: u64,
    pub max_resolution: String,
    pub codecs: Vec<String>,
    pub bitrate_kbps: u32,
    pub frame_rate: u8,
    pub enable_hls: bool,
    pub hls_segment_duration: u8,
    pub enable_dash: bool,
}

#[derive(Debug, Clone, Deserialize)]
pub struct AudioProcessingConfig {
    pub max_duration_secs: u64,
    pub sample_rate: u32,
    pub bitrate_kbps: u32,
    pub channels: u8,
    pub formats: Vec<String>,
}

#[derive(Debug, Clone, Deserialize)]
pub struct ThumbnailConfig {
    pub sizes: Vec<ThumbnailSize>,
    pub format: String,
    pub quality: u8,
}

#[derive(Debug, Clone, Deserialize)]
pub struct ThumbnailSize {
    pub name: String,
    pub width: u32,
    pub height: u32,
}

#[derive(Debug, Clone, Deserialize)]
pub struct CompressionConfig {
    pub enable_brotli: bool,
    pub enable_gzip: bool,
    pub enable_zstd: bool,
    pub level: u8,
}

#[derive(Debug, Clone, Deserialize)]
pub struct LimitsConfig {
    pub max_file_size_mb: u64,
    pub max_image_size_mb: u64,
    pub max_video_size_mb: u64,
    pub max_audio_size_mb: u64,
    pub allowed_image_types: Vec<String>,
    pub allowed_video_types: Vec<String>,
    pub allowed_audio_types: Vec<String>,
    pub rate_limit_per_minute: u32,
}

#[derive(Debug, Clone, Deserialize)]
pub struct CdnConfig {
    pub enabled: bool,
    pub base_url: String,
    pub cache_control: String,
    pub max_age_secs: u64,
}

impl Config {
    pub fn from_env() -> Result<Self, ConfigError> {
        let env_name = env::var("RUN_ENV").unwrap_or_else(|_| "development".to_string());
        
        let mut builder = ConfigLoader::builder()
            .add_source(File::with_name("config/default").required(false))
            .add_source(File::with_name(&format!("config/{}", env_name)).required(false))
            .add_source(Environment::with_prefix("MEDIA_SERVICE").separator("__"));

        let config = builder.build()?;
        config.try_deserialize()
    }

    pub fn from_file(path: &str) -> Result<Self, ConfigError> {
        let config = ConfigLoader::builder()
            .add_source(File::with_name(path))
            .build()?;
        
        config.try_deserialize()
    }
}

impl Default for Config {
    fn default() -> Self {
        Self {
            server: ServerConfig {
                host: "0.0.0.0".to_string(),
                port: 8083,
                grpc_port: 50051,
                workers: num_cpus::get(),
                max_connections: 25000,
                request_timeout_secs: 300,
            },
            storage: StorageConfig {
                provider: "local".to_string(),
                s3: S3Config {
                    access_key_id: String::new(),
                    secret_access_key: String::new(),
                    region: "us-east-1".to_string(),
                    endpoint: None,
                    use_path_style: false,
                },
                minio: MinioConfig {
                    endpoint: "localhost:9000".to_string(),
                    access_key: "minioadmin".to_string(),
                    secret_key: "minioadmin".to_string(),
                    use_ssl: false,
                },
                local: LocalConfig {
                    base_path: "./media_storage".to_string(),
                    serve_url: "http://localhost:8083/media".to_string(),
                },
                bucket_name: "socialink-media".to_string(),
                region: "us-east-1".to_string(),
            },
            database: DatabaseConfig {
                url: "postgresql://postgres:postgres@localhost/socialink_media".to_string(),
                max_connections: 100,
                min_connections: 10,
                connect_timeout_secs: 30,
                idle_timeout_secs: 600,
            },
            redis: RedisConfig {
                url: "redis://localhost:6379".to_string(),
                pool_size: 100,
                connection_timeout_secs: 5,
                cache_ttl_secs: 3600,
            },
            processing: ProcessingConfig {
                image: ImageProcessingConfig {
                    max_dimension: 4096,
                    quality: 85,
                    formats: vec!["jpeg".to_string(), "png".to_string(), "webp".to_string()],
                    webp_quality: 80,
                    thumbnail_sizes: vec![150, 300, 600, 1200],
                    enable_watermark: false,
                    watermark_path: None,
                },
                video: VideoProcessingConfig {
                    max_duration_secs: 3600,
                    max_resolution: "1920x1080".to_string(),
                    codecs: vec!["h264".to_string(), "h265".to_string(), "vp9".to_string()],
                    bitrate_kbps: 5000,
                    frame_rate: 30,
                    enable_hls: true,
                    hls_segment_duration: 10,
                    enable_dash: false,
                },
                audio: AudioProcessingConfig {
                    max_duration_secs: 7200,
                    sample_rate: 48000,
                    bitrate_kbps: 192,
                    channels: 2,
                    formats: vec!["mp3".to_string(), "aac".to_string(), "opus".to_string()],
                },
                thumbnail: ThumbnailConfig {
                    sizes: vec![
                        ThumbnailSize {
                            name: "small".to_string(),
                            width: 150,
                            height: 150,
                        },
                        ThumbnailSize {
                            name: "medium".to_string(),
                            width: 300,
                            height: 300,
                        },
                        ThumbnailSize {
                            name: "large".to_string(),
                            width: 600,
                            height: 600,
                        },
                    ],
                    format: "jpeg".to_string(),
                    quality: 85,
                },
                compression: CompressionConfig {
                    enable_brotli: true,
                    enable_gzip: true,
                    enable_zstd: false,
                    level: 6,
                },
            },
            limits: LimitsConfig {
                max_file_size_mb: 100,
                max_image_size_mb: 10,
                max_video_size_mb: 100,
                max_audio_size_mb: 50,
                allowed_image_types: vec![
                    "image/jpeg".to_string(),
                    "image/png".to_string(),
                    "image/gif".to_string(),
                    "image/webp".to_string(),
                ],
                allowed_video_types: vec![
                    "video/mp4".to_string(),
                    "video/quicktime".to_string(),
                    "video/webm".to_string(),
                ],
                allowed_audio_types: vec![
                    "audio/mpeg".to_string(),
                    "audio/mp4".to_string(),
                    "audio/wav".to_string(),
                    "audio/ogg".to_string(),
                ],
                rate_limit_per_minute: 60,
            },
            cdn: CdnConfig {
                enabled: false,
                base_url: String::new(),
                cache_control: "public, max-age=31536000".to_string(),
                max_age_secs: 31536000,
            },
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_default_config() {
        let config = Config::default();
        assert_eq!(config.server.port, 8083);
        assert_eq!(config.storage.provider, "local");
    }

    #[test]
    fn test_config_validation() {
        let config = Config::default();
        assert!(config.processing.image.quality <= 100);
        assert!(config.processing.video.frame_rate > 0);
        assert!(config.limits.max_file_size_mb > 0);
    }
}
