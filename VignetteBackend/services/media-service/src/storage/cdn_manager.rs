use reqwest::Client;
use std::collections::HashMap;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum CdnError {
    #[error("CDN upload failed: {0}")]
    UploadFailed(String),
    
    #[error("CDN invalidation failed: {0}")]
    InvalidationFailed(String),
    
    #[error("HTTP error: {0}")]
    HttpError(#[from] reqwest::Error),
}

pub type CdnResult<T> = Result<T, CdnError>;

/// CDN Manager for distributing media globally with low latency
pub struct CdnManager {
    client: Client,
    base_url: String,
    api_key: Option<String>,
    enabled: bool,
}

impl CdnManager {
    pub fn new(base_url: String, api_key: Option<String>, enabled: bool) -> Self {
        Self {
            client: Client::new(),
            base_url,
            api_key,
            enabled,
        }
    }

    /// Push file to CDN
    pub async fn push_to_cdn(
        &self,
        file_path: &str,
        data: &[u8],
        content_type: &str,
    ) -> CdnResult<String> {
        if !self.enabled {
            return Ok(format!("{}/{}", self.base_url, file_path));
        }

        // Upload to CDN provider (e.g., CloudFront, Cloudflare)
        let url = format!("{}/upload", self.base_url);
        
        let mut headers = HashMap::new();
        if let Some(key) = &self.api_key {
            headers.insert("X-API-Key".to_string(), key.clone());
        }

        let response = self.client
            .post(&url)
            .header("Content-Type", content_type)
            .body(data.to_vec())
            .send()
            .await?;

        if !response.status().is_success() {
            return Err(CdnError::UploadFailed(
                format!("CDN returned status: {}", response.status())
            ));
        }

        let cdn_url = format!("{}/{}", self.base_url, file_path);
        Ok(cdn_url)
    }

    /// Invalidate CDN cache for updated files
    pub async fn invalidate_cache(&self, paths: Vec<String>) -> CdnResult<()> {
        if !self.enabled {
            return Ok(());
        }

        let url = format!("{}/invalidate", self.base_url);
        
        let payload = serde_json::json!({
            "paths": paths
        });

        let response = self.client
            .post(&url)
            .json(&payload)
            .send()
            .await?;

        if !response.status().is_success() {
            return Err(CdnError::InvalidationFailed(
                "CDN invalidation failed".to_string()
            ));
        }

        Ok(())
    }

    /// Get CDN URL for a file
    pub fn get_cdn_url(&self, file_path: &str) -> String {
        if !self.enabled {
            return file_path.to_string();
        }

        format!("{}/{}", self.base_url, file_path.trim_start_matches('/'))
    }

    /// Pre-warm CDN cache for popular content
    pub async fn prewarm_cache(&self, paths: Vec<String>) -> CdnResult<()> {
        if !self.enabled {
            return Ok(());
        }

        // Trigger CDN to fetch and cache files
        for path in paths {
            let url = self.get_cdn_url(&path);
            self.client.head(&url).send().await.ok();
        }

        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_cdn_url_generation() {
        let cdn = CdnManager::new(
            "https://cdn.vignette.com".to_string(),
            None,
            true,
        );

        let url = cdn.get_cdn_url("images/test.jpg");
        assert!(url.starts_with("https://cdn.vignette.com"));
    }

    #[test]
    fn test_cdn_disabled() {
        let cdn = CdnManager::new(
            "https://cdn.vignette.com".to_string(),
            None,
            false,
        );

        let url = cdn.get_cdn_url("images/test.jpg");
        assert_eq!(url, "images/test.jpg");
    }
}
