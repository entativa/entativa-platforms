use actix_web::{web, HttpResponse, Result as ActixResult};
use sqlx::PgPool;
use redis::Client as RedisClient;
use uuid::Uuid;
use serde::Deserialize;

use crate::services::presence_service::PresenceService;

#[derive(Deserialize)]
pub struct SetPresenceRequest {
    pub status: String, // Online, Away, Busy, Offline
    pub device_id: String,
}

#[derive(Deserialize)]
pub struct CustomStatusRequest {
    pub custom_status: String,
}

/// Set user online
pub async fn set_online(
    pool: web::Data<PgPool>,
    redis: web::Data<RedisClient>,
    path: web::Path<Uuid>,
    request: web::Json<SetPresenceRequest>,
) -> ActixResult<HttpResponse> {
    let user_id = path.into_inner();
    let presence_service = PresenceService::new(pool.get_ref().clone(), redis.get_ref().clone());
    
    match presence_service.set_online(user_id, request.device_id.clone()).await {
        Ok(_) => Ok(HttpResponse::Ok().json(serde_json::json!({
            "success": true,
            "status": "online"
        }))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}

/// Set user offline
pub async fn set_offline(
    pool: web::Data<PgPool>,
    redis: web::Data<RedisClient>,
    user_id: web::Path<Uuid>,
) -> ActixResult<HttpResponse> {
    let presence_service = PresenceService::new(pool.get_ref().clone(), redis.get_ref().clone());
    
    match presence_service.set_offline(*user_id).await {
        Ok(_) => Ok(HttpResponse::Ok().json(serde_json::json!({
            "success": true,
            "status": "offline"
        }))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}

/// Set custom status
pub async fn set_custom_status(
    pool: web::Data<PgPool>,
    redis: web::Data<RedisClient>,
    user_id: web::Path<Uuid>,
    request: web::Json<CustomStatusRequest>,
) -> ActixResult<HttpResponse> {
    let presence_service = PresenceService::new(pool.get_ref().clone(), redis.get_ref().clone());
    
    match presence_service.set_custom_status(*user_id, request.custom_status.clone()).await {
        Ok(_) => Ok(HttpResponse::Ok().json(serde_json::json!({
            "success": true
        }))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}

/// Get user presence
pub async fn get_presence(
    pool: web::Data<PgPool>,
    redis: web::Data<RedisClient>,
    user_id: web::Path<Uuid>,
) -> ActixResult<HttpResponse> {
    let presence_service = PresenceService::new(pool.get_ref().clone(), redis.get_ref().clone());
    
    match presence_service.get_presence(*user_id).await {
        Ok(presence) => Ok(HttpResponse::Ok().json(presence)),
        Err(e) => Ok(HttpResponse::InternalServerError().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}

/// Get bulk presence
pub async fn get_bulk_presence(
    pool: web::Data<PgPool>,
    redis: web::Data<RedisClient>,
    user_ids: web::Json<Vec<Uuid>>,
) -> ActixResult<HttpResponse> {
    let presence_service = PresenceService::new(pool.get_ref().clone(), redis.get_ref().clone());
    
    match presence_service.get_bulk_presence(user_ids.into_inner()).await {
        Ok(presences) => Ok(HttpResponse::Ok().json(serde_json::json!({
            "presences": presences,
            "count": presences.len()
        }))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}

/// Get online count
pub async fn get_online_count(
    pool: web::Data<PgPool>,
    redis: web::Data<RedisClient>,
) -> ActixResult<HttpResponse> {
    let presence_service = PresenceService::new(pool.get_ref().clone(), redis.get_ref().clone());
    
    match presence_service.get_online_count().await {
        Ok(count) => Ok(HttpResponse::Ok().json(serde_json::json!({
            "online_count": count
        }))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}
