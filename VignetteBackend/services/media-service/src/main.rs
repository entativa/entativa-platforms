mod config;
mod grpc_server;
mod handlers;
mod metrics;
mod models;
mod services;
mod storage;
mod utils;

use actix_cors::Cors;
use actix_web::{middleware, web, App, HttpServer};
use sqlx::postgres::PgPoolOptions;
use std::sync::Arc;
use tracing::{info, warn};
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

use config::Config;
use storage::{LocalStorage, S3Client, StorageBackend};

pub struct AppState {
    pub config: Config,
    pub db: sqlx::PgPool,
    pub redis: redis::Client,
    pub storage: Arc<dyn StorageBackend>,
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    // Initialize tracing
    tracing_subscriber::registry()
        .with(
            tracing_subscriber::EnvFilter::try_from_default_env()
                .unwrap_or_else(|_| "vignette_media_service=info,actix_web=info".into()),
        )
        .with(tracing_subscriber::fmt::layer().json())
        .init();

    info!("Starting Entativa Media Service");

    // Load configuration
    let config = Config::from_env().unwrap_or_else(|_| {
        warn!("Failed to load config from environment, using defaults");
        Config::default()
    });

    info!(
        "Configuration loaded: storage_provider={}, server={}:{}",
        config.storage.provider, config.server.host, config.server.port
    );

    // Initialize database pool
    let db_pool = PgPoolOptions::new()
        .max_connections(config.database.max_connections)
        .min_connections(config.database.min_connections)
        .connect(&config.database.url)
        .await
        .expect("Failed to create database pool");

    info!("Database connection pool established");

    // Run migrations
    sqlx::migrate!("./migrations")
        .run(&db_pool)
        .await
        .expect("Failed to run migrations");

    info!("Database migrations completed");

    // Initialize Redis client
    let redis_client = redis::Client::open(config.redis.url.clone())
        .expect("Failed to create Redis client");

    // Test Redis connection
    let mut redis_conn = redis_client
        .get_async_connection()
        .await
        .expect("Failed to connect to Redis");
    
    redis::cmd("PING")
        .query_async::<_, String>(&mut redis_conn)
        .await
        .expect("Redis ping failed");

    info!("Redis connection established");

    // Initialize storage backend
    let storage: Arc<dyn StorageBackend> = match config.storage.provider.as_str() {
        "s3" => {
            info!("Initializing S3 storage backend");
            Arc::new(
                S3Client::new(
                    &config.storage.s3.region,
                    &config.storage.bucket_name,
                    config.storage.s3.endpoint.clone(),
                )
                .await,
            )
        }
        "local" => {
            info!("Initializing local storage backend at: {}", config.storage.local.base_path);
            Arc::new(LocalStorage::new(&config.storage.local.base_path))
        }
        provider => {
            panic!("Unsupported storage provider: {}", provider);
        }
    };

    // Create application state
    let app_state = web::Data::new(AppState {
        config: config.clone(),
        db: db_pool,
        redis: redis_client,
        storage,
    });

    // Start gRPC server in background
    let grpc_bind_addr = format!("{}:{}", config.server.host, config.server.grpc_port);
    let grpc_service = grpc_server::MediaServiceImpl::new(
        db_pool.clone(),
        redis_client.clone(),
        storage.clone(),
        config.clone(),
    ).into_service();
    
    info!("Starting gRPC server on {}", grpc_bind_addr);
    
    let grpc_addr = grpc_bind_addr.parse()
        .expect("Invalid gRPC bind address");
    
    tokio::spawn(async move {
        tonic::transport::Server::builder()
            .add_service(grpc_service)
            .serve(grpc_addr)
            .await
            .expect("gRPC server failed");
    });

    let bind_addr = format!("{}:{}", config.server.host, config.server.port);
    info!("Starting HTTP server on {}", bind_addr);

    // Start HTTP server
    HttpServer::new(move || {
        let cors = Cors::default()
            .allow_any_origin()
            .allow_any_method()
            .allow_any_header()
            .max_age(3600);

        App::new()
            .app_data(app_state.clone())
            .wrap(middleware::Logger::default())
            .wrap(middleware::Compress::default())
            .wrap(cors)
            .wrap(tracing_actix_web::TracingLogger::default())
            .configure(configure_routes)
    })
    .workers(config.server.workers)
    .max_connections(config.server.max_connections)
    .bind(&bind_addr)?
    .run()
    .await
}

fn configure_routes(cfg: &mut web::ServiceConfig) {
    cfg.service(
        web::scope("/api/v1")
            .service(
                web::scope("/media")
                    .route("/upload", web::post().to(handlers::upload::upload_media))
                    .route("/upload/multipart/init", web::post().to(handlers::upload::init_multipart_upload))
                    .route("/upload/multipart/chunk", web::post().to(handlers::upload::upload_chunk))
                    .route("/upload/multipart/complete", web::post().to(handlers::upload::complete_multipart_upload))
                    .route("/{media_id}", web::get().to(handlers::download::get_media))
                    .route("/{media_id}/download", web::get().to(handlers::download::download_media))
                    .route("/{media_id}/metadata", web::get().to(handlers::download::get_metadata))
                    .route("/{media_id}", web::delete().to(handlers::upload::delete_media))
                    .route("", web::get().to(handlers::download::list_media))
            )
            .service(
                web::scope("/process")
                    .route("/{media_id}", web::post().to(handlers::processing::process_media))
                    .route("/{media_id}/status", web::get().to(handlers::processing::get_processing_status))
                    .route("/batch", web::post().to(handlers::processing::batch_process))
            )
            .service(
                web::scope("/stream")
                    .route("/{media_id}", web::get().to(handlers::streaming::stream_media))
                    .route("/{media_id}/hls/playlist.m3u8", web::get().to(handlers::streaming::serve_hls_playlist))
                    .route("/{media_id}/hls/{segment}", web::get().to(handlers::streaming::serve_hls_segment))
            )
            .route("/health", web::get().to(health_check))
            .route("/metrics", web::get().to(metrics))
    );
}

async fn health_check() -> actix_web::Result<actix_web::HttpResponse> {
    Ok(actix_web::HttpResponse::Ok().json(serde_json::json!({
        "status": "healthy",
        "service": "vignette-media-service",
        "version": env!("CARGO_PKG_VERSION")
    })))
}

async fn metrics() -> actix_web::Result<actix_web::HttpResponse> {
    match metrics::collect_metrics() {
        Ok(metrics_output) => Ok(actix_web::HttpResponse::Ok()
            .content_type("text/plain; version=0.0.4")
            .body(metrics_output)),
        Err(e) => {
            tracing::error!("Failed to collect metrics: {}", e);
            Ok(actix_web::HttpResponse::InternalServerError()
                .body(format!("Error collecting metrics: {}", e)))
        }
    }
}
