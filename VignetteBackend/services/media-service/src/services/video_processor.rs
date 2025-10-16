use crate::models::{VideoMetadata, MediaVariant};
use thiserror::Error;
use std::process::{Command, Stdio};
use std::path::Path;
use tokio::fs;

#[derive(Error, Debug)]
pub enum VideoProcessingError {
    #[error("FFmpeg not found or not installed")]
    FFmpegNotFound,
    
    #[error("Video processing failed: {0}")]
    ProcessingFailed(String),
    
    #[error("Invalid video format: {0}")]
    InvalidFormat(String),
    
    #[error("IO error: {0}")]
    IoError(#[from] std::io::Error),
}

pub type VideoResult<T> = Result<T, VideoProcessingError>;

pub struct VideoProcessor {
    max_resolution: (u32, u32),
    target_bitrate: u32,
    target_fps: u8,
}

impl VideoProcessor {
    pub fn new(max_resolution: (u32, u32), target_bitrate: u32, target_fps: u8) -> Self {
        Self {
            max_resolution,
            target_bitrate,
            target_fps,
        }
    }

    /// Transcode video to H.264 with optimal settings for web
    pub async fn transcode_to_h264(
        &self,
        input_path: &str,
        output_path: &str,
        resolution: Option<(u32, u32)>,
    ) -> VideoResult<String> {
        let (width, height) = resolution.unwrap_or(self.max_resolution);

        // Build FFmpeg command for maximum quality
        let output = Command::new("ffmpeg")
            .args(&[
                "-i", input_path,
                "-c:v", "libx264",           // H.264 codec
                "-preset", "slow",            // Best quality (slower encoding)
                "-crf", "23",                 // Constant Rate Factor (18-28, lower=better)
                "-profile:v", "high",         // High profile for better compression
                "-level", "4.0",              // Compatibility level
                "-pix_fmt", "yuv420p",        // Pixel format for compatibility
                "-vf", &format!("scale={}:{}:flags=lanczos", width, height), // Lanczos scaling
                "-c:a", "aac",                // AAC audio codec
                "-b:a", "192k",               // Audio bitrate
                "-ar", "48000",               // Audio sample rate
                "-ac", "2",                   // Stereo audio
                "-movflags", "+faststart",    // Enable streaming
                "-y",                         // Overwrite output
                output_path,
            ])
            .stdout(Stdio::piped())
            .stderr(Stdio::piped())
            .output()
            .await?;

        if !output.status.success() {
            let error = String::from_utf8_lossy(&output.stderr);
            return Err(VideoProcessingError::ProcessingFailed(error.to_string()));
        }

        Ok(output_path.to_string())
    }

    /// Create HLS playlist for adaptive streaming
    pub async fn create_hls_stream(
        &self,
        input_path: &str,
        output_dir: &str,
        segment_duration: u8,
    ) -> VideoResult<String> {
        // Create output directory
        fs::create_dir_all(output_dir).await?;

        // Generate HLS with multiple quality variants
        let variants = vec![
            ("360p", 640, 360, 800),    // Low quality
            ("480p", 854, 480, 1400),   // Medium quality
            ("720p", 1280, 720, 2800),  // HD
            ("1080p", 1920, 1080, 5000), // Full HD
        ];

        let mut variant_playlists = Vec::new();

        for (name, width, height, bitrate) in variants {
            let variant_dir = format!("{}/{}", output_dir, name);
            fs::create_dir_all(&variant_dir).await?;

            let output = Command::new("ffmpeg")
                .args(&[
                    "-i", input_path,
                    "-c:v", "libx264",
                    "-b:v", &format!("{}k", bitrate),
                    "-s", &format!("{}x{}", width, height),
                    "-c:a", "aac",
                    "-b:a", "128k",
                    "-hls_time", &segment_duration.to_string(),
                    "-hls_playlist_type", "vod",
                    "-hls_segment_filename", &format!("{}/segment_%03d.ts", variant_dir),
                    &format!("{}/playlist.m3u8", variant_dir),
                ])
                .output()
                .await?;

            if output.status.success() {
                variant_playlists.push((name.to_string(), format!("{}/playlist.m3u8", name)));
            }
        }

        // Create master playlist
        let master_playlist = self.create_master_playlist(&variant_playlists, &variants);
        let master_path = format!("{}/master.m3u8", output_dir);
        fs::write(&master_path, master_playlist).await?;

        Ok(master_path)
    }

    fn create_master_playlist(
        &self,
        playlists: &[(String, String)],
        variants: &[(&str, u32, u32, u32)],
    ) -> String {
        let mut content = String::from("#EXTM3U\n#EXT-X-VERSION:3\n\n");

        for (i, (name, path)) in playlists.iter().enumerate() {
            if let Some((_, width, height, bitrate)) = variants.get(i) {
                content.push_str(&format!(
                    "#EXT-X-STREAM-INF:BANDWIDTH={},RESOLUTION={}x{},NAME=\"{}\"\n{}\n\n",
                    bitrate * 1000,
                    width,
                    height,
                    name,
                    path
                ));
            }
        }

        content
    }

    /// Extract frame at specific timestamp for thumbnail
    pub async fn extract_frame(
        &self,
        input_path: &str,
        timestamp_secs: f64,
        output_path: &str,
    ) -> VideoResult<String> {
        let output = Command::new("ffmpeg")
            .args(&[
                "-i", input_path,
                "-ss", &timestamp_secs.to_string(),
                "-vframes", "1",
                "-q:v", "2",              // High quality
                "-y",
                output_path,
            ])
            .output()
            .await?;

        if !output.status.success() {
            return Err(VideoProcessingError::ProcessingFailed(
                "Failed to extract frame".to_string()
            ));
        }

        Ok(output_path.to_string())
    }

    /// Extract video metadata using FFprobe
    pub async fn extract_metadata(&self, video_path: &str) -> VideoResult<VideoMetadata> {
        let output = Command::new("ffprobe")
            .args(&[
                "-v", "quiet",
                "-print_format", "json",
                "-show_format",
                "-show_streams",
                video_path,
            ])
            .output()
            .await?;

        if !output.status.success() {
            return Err(VideoProcessingError::ProcessingFailed(
                "FFprobe failed".to_string()
            ));
        }

        let json: serde_json::Value = serde_json::from_slice(&output.stdout)
            .map_err(|e| VideoProcessingError::ProcessingFailed(e.to_string()))?;

        // Parse video stream
        let video_stream = json["streams"]
            .as_array()
            .and_then(|streams| {
                streams.iter().find(|s| s["codec_type"] == "video")
            })
            .ok_or_else(|| VideoProcessingError::InvalidFormat("No video stream found".to_string()))?;

        let duration = json["format"]["duration"]
            .as_str()
            .and_then(|s| s.parse::<f64>().ok())
            .unwrap_or(0.0);

        let width = video_stream["width"].as_u64().unwrap_or(0) as u32;
        let height = video_stream["height"].as_u64().unwrap_or(0) as u32;
        let codec = video_stream["codec_name"].as_str().unwrap_or("unknown").to_string();
        
        let fps_str = video_stream["r_frame_rate"].as_str().unwrap_or("0/1");
        let fps_parts: Vec<&str> = fps_str.split('/').collect();
        let frame_rate = if fps_parts.len() == 2 {
            let num: f64 = fps_parts[0].parse().unwrap_or(0.0);
            let den: f64 = fps_parts[1].parse().unwrap_or(1.0);
            if den > 0.0 { num / den } else { 0.0 }
        } else {
            0.0
        };

        let bitrate = video_stream["bit_rate"]
            .as_str()
            .and_then(|s| s.parse::<u32>().ok())
            .unwrap_or(0);

        // Check for audio stream
        let audio_stream = json["streams"]
            .as_array()
            .and_then(|streams| streams.iter().find(|s| s["codec_type"] == "audio"));

        let (has_audio, audio_codec, audio_bitrate, audio_sample_rate) = if let Some(audio) = audio_stream {
            (
                true,
                audio["codec_name"].as_str().map(|s| s.to_string()),
                audio["bit_rate"].as_str().and_then(|s| s.parse().ok()),
                audio["sample_rate"].as_str().and_then(|s| s.parse().ok()),
            )
        } else {
            (false, None, None, None)
        };

        Ok(VideoMetadata {
            duration,
            width,
            height,
            codec,
            bitrate,
            frame_rate,
            aspect_ratio: format!("{}:{}", width, height),
            has_audio,
            audio_codec,
            audio_bitrate,
            audio_sample_rate,
            audio_channels: audio_stream.and_then(|a| a["channels"].as_u64().map(|c| c as u8)),
            keyframe_interval: None,
            total_frames: Some((duration * frame_rate) as u64),
        })
    }

    /// Optimize video for web delivery
    pub async fn optimize_for_web(
        &self,
        input_path: &str,
        output_path: &str,
    ) -> VideoResult<String> {
        let output = Command::new("ffmpeg")
            .args(&[
                "-i", input_path,
                "-c:v", "libx264",
                "-preset", "medium",
                "-crf", "23",
                "-profile:v", "main",
                "-level", "3.1",
                "-pix_fmt", "yuv420p",
                "-c:a", "aac",
                "-b:a", "128k",
                "-movflags", "+faststart",
                "-max_muxing_queue_size", "1024",
                "-y",
                output_path,
            ])
            .output()
            .await?;

        if !output.status.success() {
            return Err(VideoProcessingError::ProcessingFailed(
                String::from_utf8_lossy(&output.stderr).to_string()
            ));
        }

        Ok(output_path.to_string())
    }
}

impl Default for VideoProcessor {
    fn default() -> Self {
        Self::new((1920, 1080), 5000, 30)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_video_processor_creation() {
        let processor = VideoProcessor::default();
        assert_eq!(processor.max_resolution, (1920, 1080));
    }
}
