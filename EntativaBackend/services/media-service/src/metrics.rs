use once_cell::sync::Lazy;
use prometheus::{
    register_counter_vec, register_gauge_vec, register_histogram_vec, CounterVec, Encoder,
    GaugeVec, HistogramVec, TextEncoder,
};

// Upload metrics
pub static UPLOAD_TOTAL: Lazy<CounterVec> = Lazy::new(|| {
    register_counter_vec!(
        "media_upload_total",
        "Total number of media uploads",
        &["media_type", "status"]
    )
    .unwrap()
});

pub static UPLOAD_SIZE_BYTES: Lazy<HistogramVec> = Lazy::new(|| {
    register_histogram_vec!(
        "media_upload_size_bytes",
        "Size of uploaded media in bytes",
        &["media_type"],
        vec![
            1024.0,      // 1KB
            102400.0,    // 100KB
            1048576.0,   // 1MB
            10485760.0,  // 10MB
            104857600.0, // 100MB
            1073741824.0 // 1GB
        ]
    )
    .unwrap()
});

pub static UPLOAD_DURATION_SECONDS: Lazy<HistogramVec> = Lazy::new(|| {
    register_histogram_vec!(
        "media_upload_duration_seconds",
        "Duration of media upload in seconds",
        &["media_type"],
        vec![0.1, 0.5, 1.0, 2.5, 5.0, 10.0, 30.0, 60.0]
    )
    .unwrap()
});

// Processing metrics
pub static PROCESSING_TOTAL: Lazy<CounterVec> = Lazy::new(|| {
    register_counter_vec!(
        "media_processing_total",
        "Total number of media processing operations",
        &["media_type", "operation", "status"]
    )
    .unwrap()
});

pub static PROCESSING_DURATION_SECONDS: Lazy<HistogramVec> = Lazy::new(|| {
    register_histogram_vec!(
        "media_processing_duration_seconds",
        "Duration of media processing operations",
        &["media_type", "operation"],
        vec![0.1, 0.5, 1.0, 5.0, 10.0, 30.0, 60.0, 300.0]
    )
    .unwrap()
});

pub static PROCESSING_QUEUE_SIZE: Lazy<GaugeVec> = Lazy::new(|| {
    register_gauge_vec!(
        "media_processing_queue_size",
        "Number of items in processing queue",
        &["media_type"]
    )
    .unwrap()
});

// Download/Streaming metrics
pub static DOWNLOAD_TOTAL: Lazy<CounterVec> = Lazy::new(|| {
    register_counter_vec!(
        "media_download_total",
        "Total number of media downloads",
        &["media_type", "status"]
    )
    .unwrap()
});

pub static DOWNLOAD_BYTES: Lazy<CounterVec> = Lazy::new(|| {
    register_counter_vec!(
        "media_download_bytes_total",
        "Total bytes downloaded",
        &["media_type"]
    )
    .unwrap()
});

pub static STREAM_TOTAL: Lazy<CounterVec> = Lazy::new(|| {
    register_counter_vec!(
        "media_stream_total",
        "Total number of stream requests",
        &["media_type", "protocol"]
    )
    .unwrap()
});

// Storage metrics
pub static STORAGE_OPERATIONS: Lazy<CounterVec> = Lazy::new(|| {
    register_counter_vec!(
        "storage_operations_total",
        "Total storage operations",
        &["operation", "backend", "status"]
    )
    .unwrap()
});

pub static STORAGE_DURATION_SECONDS: Lazy<HistogramVec> = Lazy::new(|| {
    register_histogram_vec!(
        "storage_operation_duration_seconds",
        "Duration of storage operations",
        &["operation", "backend"],
        vec![0.01, 0.05, 0.1, 0.5, 1.0, 5.0, 10.0]
    )
    .unwrap()
});

// Cache metrics
pub static CACHE_HITS: Lazy<CounterVec> = Lazy::new(|| {
    register_counter_vec!(
        "cache_hits_total",
        "Total cache hits",
        &["cache_type"]
    )
    .unwrap()
});

pub static CACHE_MISSES: Lazy<CounterVec> = Lazy::new(|| {
    register_counter_vec!(
        "cache_misses_total",
        "Total cache misses",
        &["cache_type"]
    )
    .unwrap()
});

// Database metrics
pub static DB_QUERIES: Lazy<CounterVec> = Lazy::new(|| {
    register_counter_vec!(
        "db_queries_total",
        "Total database queries",
        &["operation", "status"]
    )
    .unwrap()
});

pub static DB_QUERY_DURATION_SECONDS: Lazy<HistogramVec> = Lazy::new(|| {
    register_histogram_vec!(
        "db_query_duration_seconds",
        "Database query duration",
        &["operation"],
        vec![0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1.0]
    )
    .unwrap()
});

// Error metrics
pub static ERRORS_TOTAL: Lazy<CounterVec> = Lazy::new(|| {
    register_counter_vec!(
        "errors_total",
        "Total errors",
        &["error_type", "severity"]
    )
    .unwrap()
});

// Active connections
pub static ACTIVE_CONNECTIONS: Lazy<GaugeVec> = Lazy::new(|| {
    register_gauge_vec!(
        "active_connections",
        "Number of active connections",
        &["type"]
    )
    .unwrap()
});

// Media inventory
pub static MEDIA_TOTAL: Lazy<GaugeVec> = Lazy::new(|| {
    register_gauge_vec!(
        "media_total",
        "Total number of media items",
        &["media_type", "status"]
    )
    .unwrap()
});

pub static MEDIA_STORAGE_BYTES: Lazy<GaugeVec> = Lazy::new(|| {
    register_gauge_vec!(
        "media_storage_bytes_total",
        "Total storage used by media",
        &["media_type"]
    )
    .unwrap()
});

/// Collect and encode all metrics for Prometheus
pub fn collect_metrics() -> Result<String, Box<dyn std::error::Error>> {
    let encoder = TextEncoder::new();
    let metric_families = prometheus::gather();
    let mut buffer = Vec::new();
    encoder.encode(&metric_families, &mut buffer)?;
    Ok(String::from_utf8(buffer)?)
}

/// Record upload metrics
pub fn record_upload(media_type: &str, size_bytes: u64, duration_secs: f64, success: bool) {
    let status = if success { "success" } else { "failure" };
    
    UPLOAD_TOTAL
        .with_label_values(&[media_type, status])
        .inc();
    
    if success {
        UPLOAD_SIZE_BYTES
            .with_label_values(&[media_type])
            .observe(size_bytes as f64);
        
        UPLOAD_DURATION_SECONDS
            .with_label_values(&[media_type])
            .observe(duration_secs);
    }
}

/// Record processing metrics
pub fn record_processing(media_type: &str, operation: &str, duration_secs: f64, success: bool) {
    let status = if success { "success" } else { "failure" };
    
    PROCESSING_TOTAL
        .with_label_values(&[media_type, operation, status])
        .inc();
    
    if success {
        PROCESSING_DURATION_SECONDS
            .with_label_values(&[media_type, operation])
            .observe(duration_secs);
    }
}

/// Record download metrics
pub fn record_download(media_type: &str, size_bytes: u64, success: bool) {
    let status = if success { "success" } else { "failure" };
    
    DOWNLOAD_TOTAL
        .with_label_values(&[media_type, status])
        .inc();
    
    if success {
        DOWNLOAD_BYTES
            .with_label_values(&[media_type])
            .inc_by(size_bytes);
    }
}

/// Record streaming metrics
pub fn record_stream(media_type: &str, protocol: &str) {
    STREAM_TOTAL
        .with_label_values(&[media_type, protocol])
        .inc();
}

/// Record storage operation
pub fn record_storage_operation(operation: &str, backend: &str, duration_secs: f64, success: bool) {
    let status = if success { "success" } else { "failure" };
    
    STORAGE_OPERATIONS
        .with_label_values(&[operation, backend, status])
        .inc();
    
    STORAGE_DURATION_SECONDS
        .with_label_values(&[operation, backend])
        .observe(duration_secs);
}

/// Record cache operation
pub fn record_cache_hit(cache_type: &str) {
    CACHE_HITS
        .with_label_values(&[cache_type])
        .inc();
}

pub fn record_cache_miss(cache_type: &str) {
    CACHE_MISSES
        .with_label_values(&[cache_type])
        .inc();
}

/// Record database query
pub fn record_db_query(operation: &str, duration_secs: f64, success: bool) {
    let status = if success { "success" } else { "failure" };
    
    DB_QUERIES
        .with_label_values(&[operation, status])
        .inc();
    
    DB_QUERY_DURATION_SECONDS
        .with_label_values(&[operation])
        .observe(duration_secs);
}

/// Record error
pub fn record_error(error_type: &str, severity: &str) {
    ERRORS_TOTAL
        .with_label_values(&[error_type, severity])
        .inc();
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_record_upload() {
        record_upload("image", 1024000, 1.5, true);
        // Verify metrics are recorded (would check actual values in integration test)
    }

    #[test]
    fn test_collect_metrics() {
        let metrics = collect_metrics();
        assert!(metrics.is_ok());
        assert!(metrics.unwrap().contains("media_upload_total"));
    }
}
