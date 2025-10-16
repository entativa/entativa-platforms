mod models;
mod crypto;
mod services;
mod handlers;
mod websocket;

use actix_web::{web, App, HttpServer, HttpResponse, middleware};
use actix_cors::Cors;
use sqlx::postgres::PgPoolOptions;
use std::sync::Arc;
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    // Load environment
    dotenv::dotenv().ok();
    
    // Initialize tracing
    tracing_subscriber::registry()
        .with(tracing_subscriber::EnvFilter::new(
            std::env::var("RUST_LOG").unwrap_or_else(|_| "info".into()),
        ))
        .with(tracing_subscriber::fmt::layer())
        .init();
    
    tracing::info!("🔐 Starting Vignette Messaging Service (Signal-level E2EE)...");
    
    // Database connection
    let database_url = std::env::var("DATABASE_URL")
        .expect("DATABASE_URL must be set");
    
    tracing::info!("📊 Connecting to PostgreSQL...");
    let pool = PgPoolOptions::new()
        .max_connections(20)
        .connect(&database_url)
        .await
        .expect("Failed to connect to database");
    
    tracing::info!("✅ PostgreSQL connected");
    
    // Redis connection
    let redis_url = std::env::var("REDIS_URL")
        .unwrap_or_else(|_| "redis://localhost:6379".to_string());
    
    tracing::info!("💾 Connecting to Redis...");
    let redis_client = redis::Client::open(redis_url)
        .expect("Failed to connect to Redis");
    
    tracing::info!("✅ Redis connected");
    
    // Initialize WebSocket server
    let ws_server = Arc::new(websocket::ws_server::WsServer::new(redis_client.clone()));
    
    // Start Redis subscription for real-time delivery
    let ws_server_clone = ws_server.clone();
    tokio::spawn(async move {
        ws_server_clone.start_redis_subscription().await;
    });
    
    tracing::info!("🌐 WebSocket server initialized");
    
    // Server configuration
    let host = std::env::var("HOST").unwrap_or_else(|_| "0.0.0.0".to_string());
    let port = std::env::var("PORT").unwrap_or_else(|_| "8091".to_string());
    let bind_addr = format!("{}:{}", host, port);
    
    tracing::info!("🚀 Starting HTTP server on {}...", bind_addr);
    
    HttpServer::new(move || {
        let cors = Cors::default()
            .allow_any_origin()
            .allow_any_method()
            .allow_any_header()
            .max_age(3600);
        
        App::new()
            .app_data(web::Data::new(pool.clone()))
            .app_data(web::Data::new(redis_client.clone()))
            .app_data(web::Data::new(ws_server.clone()))
            .wrap(cors)
            .wrap(middleware::Logger::default())
            .wrap(middleware::Compress::default())
            
            // Health check
            .route("/health", web::get().to(health_check))
            .route("/", web::get().to(root))
            
            // API v1
            .service(
                web::scope("/api/v1")
                    // Key Management
                    .service(
                        web::scope("/keys")
                            .route("/register/{user_id}", web::post().to(handlers::key_handler::register_device))
                            .route("/bundle/{user_id}", web::get().to(handlers::key_handler::get_prekey_bundle))
                            .route("/rotate/{user_id}/{device_id}", web::put().to(handlers::key_handler::rotate_signed_prekey))
                            .route("/prekeys/{user_id}/{device_id}", web::post().to(handlers::key_handler::upload_onetime_prekeys))
                            .route("/deactivate/{user_id}/{device_id}", web::delete().to(handlers::key_handler::deactivate_device))
                            .route("/devices/{user_id}", web::get().to(handlers::key_handler::get_user_devices))
                            .route("/stats/{user_id}/{device_id}", web::get().to(handlers::key_handler::get_key_stats))
                    )
                    // Messages
                    .service(
                        web::scope("/messages")
                            .route("/send/{sender_id}", web::post().to(handlers::message_handler::send_message))
                            .route("/conversation/{user_id}", web::get().to(handlers::message_handler::get_messages))
                            .route("/delivered/{user_id}/{message_id}", web::put().to(handlers::message_handler::mark_delivered))
                            .route("/read/{user_id}/{message_id}", web::put().to(handlers::message_handler::mark_read))
                            .route("/delete/{user_id}/{message_id}", web::delete().to(handlers::message_handler::delete_message))
                            .route("/queue/{user_id}/{device_id}", web::get().to(handlers::message_handler::get_offline_queue))
                            .route("/queue/{user_id}/{device_id}", web::delete().to(handlers::message_handler::clear_offline_queue))
                    )
                    // Groups
                    .service(
                        web::scope("/groups")
                            .route("/create/{creator_id}", web::post().to(handlers::group_handler::create_group))
                            .route("/{group_id}/members/{added_by}", web::post().to(handlers::group_handler::add_member))
                            .route("/{group_id}/members/{user_id}/{removed_by}", web::delete().to(handlers::group_handler::remove_member))
                            .route("/{group_id}/send/{sender_id}", web::post().to(handlers::group_handler::send_group_message))
                            .route("/{group_id}", web::get().to(handlers::group_handler::get_group_info))
                            .route("/{group_id}/members", web::get().to(handlers::group_handler::get_group_members))
                    )
                    // Presence
                    .service(
                        web::scope("/presence")
                            .route("/online/{user_id}", web::put().to(handlers::presence_handler::set_online))
                            .route("/offline/{user_id}", web::put().to(handlers::presence_handler::set_offline))
                            .route("/status/{user_id}", web::put().to(handlers::presence_handler::set_custom_status))
                            .route("/{user_id}", web::get().to(handlers::presence_handler::get_presence))
                            .route("/bulk", web::post().to(handlers::presence_handler::get_bulk_presence))
                            .route("/online-count", web::get().to(handlers::presence_handler::get_online_count))
                    )
                    // Typing
                    .service(
                        web::scope("/typing")
                            .route("/{conversation_id}/{user_id}", web::put().to(handlers::typing_handler::set_typing))
                            .route("/{conversation_id}/{user_id}", web::delete().to(handlers::typing_handler::clear_typing))
                            .route("/{conversation_id}", web::get().to(handlers::typing_handler::get_typing_users))
                    )
                    // Calls
                    .service(
                        web::scope("/calls")
                            .route("/initiate/{caller_id}", web::post().to(handlers::call_handler::initiate_call))
                            .route("/{call_id}/answer/{user_id}", web::put().to(handlers::call_handler::answer_call))
                            .route("/{call_id}/decline/{user_id}", web::put().to(handlers::call_handler::decline_call))
                            .route("/{call_id}/end/{user_id}", web::put().to(handlers::call_handler::end_call))
                            .route("/{call_id}/ice/{user_id}", web::post().to(handlers::call_handler::add_ice_candidate))
                            .route("/{call_id}/ice", web::get().to(handlers::call_handler::get_ice_candidates))
                            .route("/history/{conversation_id}", web::get().to(handlers::call_handler::get_call_history))
                    )
            )
            
            // WebSocket
            .route("/ws/{user_id}/{device_id}", web::get().to(websocket::ws_server::ws_route))
    })
    .bind(&bind_addr)?
    .run()
    .await
}

async fn health_check() -> HttpResponse {
    HttpResponse::Ok().json(serde_json::json!({
        "status": "healthy",
        "service": "Vignette Messaging Service",
        "version": "1.0.0",
        "encryption": "Signal Protocol + MLS"
    }))
}

async fn root() -> HttpResponse {
    HttpResponse::Ok().json(serde_json::json!({
        "service": "Vignette Messaging Service",
        "version": "1.0.0",
        "description": "Signal-level E2EE messaging with libsignal + MLS",
        "features": [
            "End-to-end encryption (server cannot decrypt)",
            "1:1 messaging with Signal Protocol",
            "Group chats up to 1,500 members with MLS",
            "Perfect Forward Secrecy",
            "Post-Compromise Security",
            "Offline message queue",
            "Real-time WebSocket delivery",
            "Read receipts",
            "Self-destructing messages",
            "Multi-device support",
            "Presence tracking (online/offline)",
            "Typing indicators",
            "Audio/video calls with WebRTC",
            "ICE candidate exchange"
        ],
        "endpoints": {
            "keys": "/api/v1/keys/*",
            "messages": "/api/v1/messages/*",
            "groups": "/api/v1/groups/*",
            "presence": "/api/v1/presence/*",
            "typing": "/api/v1/typing/*",
            "calls": "/api/v1/calls/*",
            "websocket": "/ws/{user_id}/{device_id}"
        },
        "documentation": "https://docs.vignette.com/messaging"
    }))
}
