// MinIO client implementation
// MinIO is S3-compatible, so we can use the S3Client with custom endpoint

use super::s3_client::S3Client;

pub type MinioClient = S3Client;

// MinIO-specific helper functions can be added here
pub async fn create_minio_client(
    endpoint: &str,
    access_key: &str,
    secret_key: &str,
    bucket: &str,
    region: &str,
) -> MinioClient {
    std::env::set_var("AWS_ACCESS_KEY_ID", access_key);
    std::env::set_var("AWS_SECRET_ACCESS_KEY", secret_key);
    
    S3Client::new(region, bucket, Some(endpoint.to_string())).await
}
