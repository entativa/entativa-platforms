use super::{FileMetadata, StorageBackend, StorageError, StorageResult};
use async_trait::async_trait;
use aws_sdk_s3::{
    config::Region,
    operation::get_object::GetObjectError,
    primitives::ByteStream,
    types::{Delete, ObjectIdentifier},
    Client,
};
use bytes::Bytes;
use std::time::Duration;

pub struct S3Client {
    client: Client,
    bucket: String,
}

impl S3Client {
    pub async fn new(region: &str, bucket: &str, endpoint: Option<String>) -> Self {
        let config = if let Some(endpoint_url) = endpoint {
            let shared_config = aws_config::from_env()
                .region(Region::new(region.to_string()))
                .endpoint_url(endpoint_url)
                .load()
                .await;
            aws_sdk_s3::Config::from(&shared_config)
        } else {
            let shared_config = aws_config::from_env()
                .region(Region::new(region.to_string()))
                .load()
                .await;
            aws_sdk_s3::Config::from(&shared_config)
        };

        let client = Client::from_conf(config);

        Self {
            client,
            bucket: bucket.to_string(),
        }
    }

    fn normalize_path(&self, path: &str) -> String {
        path.trim_start_matches('/').to_string()
    }
}

#[async_trait]
impl StorageBackend for S3Client {
    async fn upload(&self, path: &str, data: Bytes, content_type: &str) -> StorageResult<String> {
        let path = self.normalize_path(path);
        let body = ByteStream::from(data);

        self.client
            .put_object()
            .bucket(&self.bucket)
            .key(&path)
            .body(body)
            .content_type(content_type)
            .send()
            .await
            .map_err(|e| StorageError::S3Error(e.to_string()))?;

        Ok(path)
    }

    async fn download(&self, path: &str) -> StorageResult<Bytes> {
        let path = self.normalize_path(path);

        let resp = self
            .client
            .get_object()
            .bucket(&self.bucket)
            .key(&path)
            .send()
            .await
            .map_err(|e| {
                if matches!(e.as_service_error(), Some(GetObjectError::NoSuchKey(_))) {
                    StorageError::NotFound(path.clone())
                } else {
                    StorageError::S3Error(e.to_string())
                }
            })?;

        let data = resp
            .body
            .collect()
            .await
            .map_err(|e| StorageError::S3Error(e.to_string()))?
            .into_bytes();

        Ok(data)
    }

    async fn delete(&self, path: &str) -> StorageResult<()> {
        let path = self.normalize_path(path);

        self.client
            .delete_object()
            .bucket(&self.bucket)
            .key(&path)
            .send()
            .await
            .map_err(|e| StorageError::S3Error(e.to_string()))?;

        Ok(())
    }

    async fn exists(&self, path: &str) -> StorageResult<bool> {
        let path = self.normalize_path(path);

        match self
            .client
            .head_object()
            .bucket(&self.bucket)
            .key(&path)
            .send()
            .await
        {
            Ok(_) => Ok(true),
            Err(e) => {
                if e.to_string().contains("404") || e.to_string().contains("NotFound") {
                    Ok(false)
                } else {
                    Err(StorageError::S3Error(e.to_string()))
                }
            }
        }
    }

    async fn get_url(&self, path: &str) -> StorageResult<String> {
        let path = self.normalize_path(path);
        Ok(format!(
            "https://{}.s3.amazonaws.com/{}",
            self.bucket, path
        ))
    }

    async fn get_signed_url(&self, path: &str, expires_in_secs: u64) -> StorageResult<String> {
        let path = self.normalize_path(path);

        let presigned = self
            .client
            .get_object()
            .bucket(&self.bucket)
            .key(&path)
            .presigned(
                aws_sdk_s3::presigning::PresigningConfig::expires_in(
                    Duration::from_secs(expires_in_secs),
                )
                .map_err(|e| StorageError::S3Error(e.to_string()))?,
            )
            .await
            .map_err(|e| StorageError::S3Error(e.to_string()))?;

        Ok(presigned.uri().to_string())
    }

    async fn list(&self, prefix: &str) -> StorageResult<Vec<String>> {
        let prefix = self.normalize_path(prefix);

        let resp = self
            .client
            .list_objects_v2()
            .bucket(&self.bucket)
            .prefix(&prefix)
            .send()
            .await
            .map_err(|e| StorageError::S3Error(e.to_string()))?;

        let keys = resp
            .contents()
            .iter()
            .filter_map(|obj| obj.key().map(|k| k.to_string()))
            .collect();

        Ok(keys)
    }

    async fn copy(&self, from: &str, to: &str) -> StorageResult<()> {
        let from = self.normalize_path(from);
        let to = self.normalize_path(to);
        let copy_source = format!("{}/{}", self.bucket, from);

        self.client
            .copy_object()
            .bucket(&self.bucket)
            .copy_source(&copy_source)
            .key(&to)
            .send()
            .await
            .map_err(|e| StorageError::S3Error(e.to_string()))?;

        Ok(())
    }

    async fn get_metadata(&self, path: &str) -> StorageResult<FileMetadata> {
        let path = self.normalize_path(path);

        let resp = self
            .client
            .head_object()
            .bucket(&self.bucket)
            .key(&path)
            .send()
            .await
            .map_err(|e| StorageError::S3Error(e.to_string()))?;

        Ok(FileMetadata {
            size: resp.content_length().unwrap_or(0) as u64,
            content_type: resp
                .content_type()
                .unwrap_or("application/octet-stream")
                .to_string(),
            last_modified: resp.last_modified().map(|t| {
                chrono::DateTime::from_timestamp(t.secs(), 0)
                    .unwrap_or_else(chrono::Utc::now)
            }),
            etag: resp.e_tag().map(|e| e.to_string()),
        })
    }
}
