use super::{FileMetadata, StorageBackend, StorageError, StorageResult};
use async_trait::async_trait;
use bytes::Bytes;
use std::path::{Path, PathBuf};
use tokio::fs;
use tokio::io::AsyncWriteExt;

pub struct LocalStorage {
    base_path: PathBuf,
}

impl LocalStorage {
    pub fn new(base_path: &str) -> Self {
        let path = PathBuf::from(base_path);
        
        // Create base directory if it doesn't exist
        std::fs::create_dir_all(&path).ok();
        
        Self { base_path: path }
    }

    fn get_full_path(&self, path: &str) -> PathBuf {
        let clean_path = path.trim_start_matches('/');
        self.base_path.join(clean_path)
    }

    async fn ensure_parent_dir(&self, path: &Path) -> StorageResult<()> {
        if let Some(parent) = path.parent() {
            fs::create_dir_all(parent).await?;
        }
        Ok(())
    }
}

#[async_trait]
impl StorageBackend for LocalStorage {
    async fn upload(&self, path: &str, data: Bytes, _content_type: &str) -> StorageResult<String> {
        let full_path = self.get_full_path(path);
        self.ensure_parent_dir(&full_path).await?;

        let mut file = fs::File::create(&full_path).await?;
        file.write_all(&data).await?;
        file.sync_all().await?;

        Ok(path.to_string())
    }

    async fn download(&self, path: &str) -> StorageResult<Bytes> {
        let full_path = self.get_full_path(path);

        if !full_path.exists() {
            return Err(StorageError::NotFound(path.to_string()));
        }

        let data = fs::read(&full_path).await?;
        Ok(Bytes::from(data))
    }

    async fn delete(&self, path: &str) -> StorageResult<()> {
        let full_path = self.get_full_path(path);

        if full_path.exists() {
            fs::remove_file(&full_path).await?;
        }

        Ok(())
    }

    async fn exists(&self, path: &str) -> StorageResult<bool> {
        let full_path = self.get_full_path(path);
        Ok(full_path.exists())
    }

    async fn get_url(&self, path: &str) -> StorageResult<String> {
        // For local storage, return a file:// URL or relative path
        Ok(format!("/media/{}", path.trim_start_matches('/')))
    }

    async fn get_signed_url(&self, path: &str, _expires_in_secs: u64) -> StorageResult<String> {
        // Local storage doesn't support signed URLs, return regular URL
        self.get_url(path).await
    }

    async fn list(&self, prefix: &str) -> StorageResult<Vec<String>> {
        let full_path = self.get_full_path(prefix);
        let mut results = Vec::new();

        if !full_path.exists() {
            return Ok(results);
        }

        let mut entries = fs::read_dir(&full_path).await?;

        while let Some(entry) = entries.next_entry().await? {
            if let Ok(relative) = entry.path().strip_prefix(&self.base_path) {
                results.push(relative.to_string_lossy().to_string());
            }
        }

        Ok(results)
    }

    async fn copy(&self, from: &str, to: &str) -> StorageResult<()> {
        let from_path = self.get_full_path(from);
        let to_path = self.get_full_path(to);

        if !from_path.exists() {
            return Err(StorageError::NotFound(from.to_string()));
        }

        self.ensure_parent_dir(&to_path).await?;
        fs::copy(&from_path, &to_path).await?;

        Ok(())
    }

    async fn get_metadata(&self, path: &str) -> StorageResult<FileMetadata> {
        let full_path = self.get_full_path(path);

        if !full_path.exists() {
            return Err(StorageError::NotFound(path.to_string()));
        }

        let metadata = fs::metadata(&full_path).await?;
        let content_type = mime_guess::from_path(&full_path)
            .first_or_octet_stream()
            .to_string();

        Ok(FileMetadata {
            size: metadata.len(),
            content_type,
            last_modified: metadata.modified().ok().map(|t| {
                chrono::DateTime::from(t)
            }),
            etag: None,
        })
    }
}
