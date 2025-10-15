use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

#[derive(Debug, Clone, Serialize, Deserialize, Default)]
pub struct ImageMetadata {
    pub width: u32,
    pub height: u32,
    pub format: String,
    pub color_space: Option<String>,
    pub bit_depth: Option<u8>,
    pub has_alpha: bool,
    pub exif: Option<ExifMetadata>,
    pub dominant_colors: Vec<Color>,
    pub average_color: Option<Color>,
    pub histogram: Option<ColorHistogram>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ExifMetadata {
    pub camera_make: Option<String>,
    pub camera_model: Option<String>,
    pub lens_make: Option<String>,
    pub lens_model: Option<String>,
    pub focal_length: Option<f64>,
    pub aperture: Option<f64>,
    pub shutter_speed: Option<String>,
    pub iso: Option<u32>,
    pub flash: Option<bool>,
    pub white_balance: Option<String>,
    pub exposure_program: Option<String>,
    pub metering_mode: Option<String>,
    pub orientation: Option<u32>,
    pub date_time_original: Option<DateTime<Utc>>,
    pub gps: Option<GpsMetadata>,
    pub copyright: Option<String>,
    pub artist: Option<String>,
    pub software: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct GpsMetadata {
    pub latitude: f64,
    pub longitude: f64,
    pub altitude: Option<f64>,
    pub timestamp: Option<DateTime<Utc>>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Color {
    pub r: u8,
    pub g: u8,
    pub b: u8,
    pub a: Option<u8>,
    pub hex: String,
    pub percentage: Option<f32>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ColorHistogram {
    pub red: Vec<u32>,
    pub green: Vec<u32>,
    pub blue: Vec<u32>,
}

#[derive(Debug, Clone, Serialize, Deserialize, Default)]
pub struct VideoMetadata {
    pub duration: f64,
    pub width: u32,
    pub height: u32,
    pub codec: String,
    pub bitrate: u32,
    pub frame_rate: f64,
    pub aspect_ratio: String,
    pub has_audio: bool,
    pub audio_codec: Option<String>,
    pub audio_bitrate: Option<u32>,
    pub audio_sample_rate: Option<u32>,
    pub audio_channels: Option<u8>,
    pub keyframe_interval: Option<u32>,
    pub total_frames: Option<u64>,
}

#[derive(Debug, Clone, Serialize, Deserialize, Default)]
pub struct AudioMetadata {
    pub duration: f64,
    pub codec: String,
    pub bitrate: u32,
    pub sample_rate: u32,
    pub channels: u8,
    pub bits_per_sample: Option<u8>,
    pub id3: Option<Id3Metadata>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Id3Metadata {
    pub title: Option<String>,
    pub artist: Option<String>,
    pub album: Option<String>,
    pub album_artist: Option<String>,
    pub genre: Option<String>,
    pub year: Option<i32>,
    pub track_number: Option<u32>,
    pub disc_number: Option<u32>,
    pub duration: Option<u32>,
    pub cover_art: Option<Vec<u8>>,
}

impl Color {
    pub fn new(r: u8, g: u8, b: u8, a: Option<u8>) -> Self {
        let hex = if let Some(alpha) = a {
            format!("#{:02x}{:02x}{:02x}{:02x}", r, g, b, alpha)
        } else {
            format!("#{:02x}{:02x}{:02x}", r, g, b)
        };

        Self {
            r,
            g,
            b,
            a,
            hex,
            percentage: None,
        }
    }

    pub fn from_hex(hex: &str) -> Option<Self> {
        let hex = hex.trim_start_matches('#');
        if hex.len() == 6 {
            let r = u8::from_str_radix(&hex[0..2], 16).ok()?;
            let g = u8::from_str_radix(&hex[2..4], 16).ok()?;
            let b = u8::from_str_radix(&hex[4..6], 16).ok()?;
            Some(Self::new(r, g, b, None))
        } else if hex.len() == 8 {
            let r = u8::from_str_radix(&hex[0..2], 16).ok()?;
            let g = u8::from_str_radix(&hex[2..4], 16).ok()?;
            let b = u8::from_str_radix(&hex[4..6], 16).ok()?;
            let a = u8::from_str_radix(&hex[6..8], 16).ok()?;
            Some(Self::new(r, g, b, Some(a)))
        } else {
            None
        }
    }

    pub fn luminance(&self) -> f32 {
        (0.299 * self.r as f32 + 0.587 * self.g as f32 + 0.114 * self.b as f32) / 255.0
    }
}
