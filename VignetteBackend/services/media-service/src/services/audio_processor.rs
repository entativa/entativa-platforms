use crate::models::AudioMetadata;
use thiserror::Error;
use std::process::Command;

#[derive(Error, Debug)]
pub enum AudioProcessingError {
    #[error("FFmpeg not found")]
    FFmpegNotFound,
    
    #[error("Audio processing failed: {0}")]
    ProcessingFailed(String),
    
    #[error("Invalid audio format: {0}")]
    InvalidFormat(String),
    
    #[error("IO error: {0}")]
    IoError(#[from] std::io::Error),
}

pub type AudioResult<T> = Result<T, AudioProcessingError>;

pub struct AudioProcessor {
    sample_rate: u32,
    bitrate: u32,
    channels: u8,
}

impl AudioProcessor {
    pub fn new(sample_rate: u32, bitrate: u32, channels: u8) -> Self {
        Self {
            sample_rate,
            bitrate,
            channels,
        }
    }

    /// Transcode audio to AAC format (best for web)
    pub async fn transcode_to_aac(
        &self,
        input_path: &str,
        output_path: &str,
    ) -> AudioResult<String> {
        let output = Command::new("ffmpeg")
            .args(&[
                "-i", input_path,
                "-c:a", "aac",
                "-b:a", &format!("{}k", self.bitrate / 1000),
                "-ar", &self.sample_rate.to_string(),
                "-ac", &self.channels.to_string(),
                "-y",
                output_path,
            ])
            .output()
            .await?;

        if !output.status.success() {
            return Err(AudioProcessingError::ProcessingFailed(
                String::from_utf8_lossy(&output.stderr).to_string()
            ));
        }

        Ok(output_path.to_string())
    }

    /// Transcode audio to MP3
    pub async fn transcode_to_mp3(
        &self,
        input_path: &str,
        output_path: &str,
        quality: u8,
    ) -> AudioResult<String> {
        // VBR quality: 0 (best) to 9 (worst), default 4
        let output = Command::new("ffmpeg")
            .args(&[
                "-i", input_path,
                "-c:a", "libmp3lame",
                "-q:a", &quality.to_string(),
                "-ar", &self.sample_rate.to_string(),
                "-y",
                output_path,
            ])
            .output()
            .await?;

        if !output.status.success() {
            return Err(AudioProcessingError::ProcessingFailed(
                "MP3 encoding failed".to_string()
            ));
        }

        Ok(output_path.to_string())
    }

    /// Normalize audio levels for consistent volume
    pub async fn normalize_audio(
        &self,
        input_path: &str,
        output_path: &str,
        target_lufs: f32,
    ) -> AudioResult<String> {
        // Use loudnorm filter for EBU R128 loudness normalization
        let output = Command::new("ffmpeg")
            .args(&[
                "-i", input_path,
                "-af", &format!("loudnorm=I={}:TP=-1.5:LRA=11", target_lufs),
                "-c:a", "aac",
                "-b:a", &format!("{}k", self.bitrate / 1000),
                "-y",
                output_path,
            ])
            .output()
            .await?;

        if !output.status.success() {
            return Err(AudioProcessingError::ProcessingFailed(
                "Audio normalization failed".to_string()
            ));
        }

        Ok(output_path.to_string())
    }

    /// Extract audio from video file
    pub async fn extract_audio_from_video(
        &self,
        input_path: &str,
        output_path: &str,
    ) -> AudioResult<String> {
        let output = Command::new("ffmpeg")
            .args(&[
                "-i", input_path,
                "-vn",                    // No video
                "-c:a", "aac",
                "-b:a", "192k",
                "-y",
                output_path,
            ])
            .output()
            .await?;

        if !output.status.success() {
            return Err(AudioProcessingError::ProcessingFailed(
                "Audio extraction failed".to_string()
            ));
        }

        Ok(output_path.to_string())
    }

    /// Generate waveform data for visualization
    pub async fn generate_waveform(
        &self,
        input_path: &str,
        sample_points: usize,
    ) -> AudioResult<Vec<f32>> {
        // Use FFmpeg to generate waveform JSON
        let output = Command::new("ffmpeg")
            .args(&[
                "-i", input_path,
                "-af", &format!("compand,showwavespic=s={}x200", sample_points),
                "-frames:v", "1",
                "-f", "image2",
                "-",
            ])
            .output()
            .await?;

        if !output.status.success() {
            return Err(AudioProcessingError::ProcessingFailed(
                "Waveform generation failed".to_string()
            ));
        }

        // Parse waveform data (simplified)
        // In production, analyze the actual audio samples
        Ok(vec![0.0; sample_points])
    }

    /// Extract audio metadata
    pub async fn extract_metadata(&self, audio_path: &str) -> AudioResult<AudioMetadata> {
        let output = Command::new("ffprobe")
            .args(&[
                "-v", "quiet",
                "-print_format", "json",
                "-show_format",
                "-show_streams",
                audio_path,
            ])
            .output()
            .await?;

        if !output.status.success() {
            return Err(AudioProcessingError::ProcessingFailed(
                "FFprobe failed".to_string()
            ));
        }

        let json: serde_json::Value = serde_json::from_slice(&output.stdout)
            .map_err(|e| AudioProcessingError::ProcessingFailed(e.to_string()))?;

        let audio_stream = json["streams"]
            .as_array()
            .and_then(|s| s.iter().find(|stream| stream["codec_type"] == "audio"))
            .ok_or_else(|| AudioProcessingError::InvalidFormat("No audio stream".to_string()))?;

        let duration = json["format"]["duration"]
            .as_str()
            .and_then(|s| s.parse::<f64>().ok())
            .unwrap_or(0.0);

        let codec = audio_stream["codec_name"]
            .as_str()
            .unwrap_or("unknown")
            .to_string();

        let bitrate = audio_stream["bit_rate"]
            .as_str()
            .and_then(|s| s.parse::<u32>().ok())
            .unwrap_or(0);

        let sample_rate = audio_stream["sample_rate"]
            .as_str()
            .and_then(|s| s.parse::<u32>().ok())
            .unwrap_or(44100);

        let channels = audio_stream["channels"]
            .as_u64()
            .unwrap_or(2) as u8;

        Ok(AudioMetadata {
            duration,
            codec,
            bitrate,
            sample_rate,
            channels,
            bits_per_sample: None,
            id3: None, // Would parse ID3 tags using id3 crate
        })
    }

    /// Convert audio to different format
    pub async fn convert_format(
        &self,
        input_path: &str,
        output_path: &str,
        format: &str,
    ) -> AudioResult<String> {
        let codec = match format {
            "mp3" => "libmp3lame",
            "aac" => "aac",
            "opus" => "libopus",
            "flac" => "flac",
            _ => return Err(AudioProcessingError::InvalidFormat(format!("Unsupported format: {}", format))),
        };

        let output = Command::new("ffmpeg")
            .args(&[
                "-i", input_path,
                "-c:a", codec,
                "-b:a", &format!("{}k", self.bitrate / 1000),
                "-y",
                output_path,
            ])
            .output()
            .await?;

        if !output.status.success() {
            return Err(AudioProcessingError::ProcessingFailed(
                "Format conversion failed".to_string()
            ));
        }

        Ok(output_path.to_string())
    }
}

impl Default for AudioProcessor {
    fn default() -> Self {
        Self::new(48000, 192000, 2)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_audio_processor_creation() {
        let processor = AudioProcessor::default();
        assert_eq!(processor.sample_rate, 48000);
        assert_eq!(processor.channels, 2);
    }
}
