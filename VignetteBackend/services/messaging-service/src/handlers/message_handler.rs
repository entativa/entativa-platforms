use actix_web::{web, HttpResponse, Result as ActixResult};
use sqlx::PgPool;
use redis::Client as RedisClient;
use uuid::Uuid;

use crate::models::message::*;
use crate::services::message_service::MessageService;

/// Send message
pub async fn send_message(
    pool: web::Data<PgPool>,
    redis: web::Data<RedisClient>,
    sender_id: web::Path<Uuid>,
    request: web::Json<SendMessageRequest>,
) -> ActixResult<HttpResponse> {
    let message_service = MessageService::new(pool.get_ref().clone(), redis.get_ref().clone());
    
    match message_service.send_message(*sender_id, request.into_inner()).await {
        Ok(response) => Ok(HttpResponse::Created().json(response)),
        Err(e) => Ok(HttpResponse::BadRequest().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}

/// Get messages for conversation
pub async fn get_messages(
    pool: web::Data<PgPool>,
    redis: web::Data<RedisClient>,
    user_id: web::Path<Uuid>,
    query: web::Query<GetMessagesRequest>,
) -> ActixResult<HttpResponse> {
    let message_service = MessageService::new(pool.get_ref().clone(), redis.get_ref().clone());
    
    match message_service.get_messages(
        *user_id,
        query.conversation_id,
        query.before_sequence,
        query.limit,
    ).await {
        Ok(messages) => Ok(HttpResponse::Ok().json(serde_json::json!({
            "messages": messages,
            "count": messages.len()
        }))),
        Err(e) => Ok(HttpResponse::BadRequest().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}

/// Mark message as delivered
pub async fn mark_delivered(
    pool: web::Data<PgPool>,
    redis: web::Data<RedisClient>,
    path: web::Path<(Uuid, Uuid)>,
) -> ActixResult<HttpResponse> {
    let (user_id, message_id) = path.into_inner();
    let message_service = MessageService::new(pool.get_ref().clone(), redis.get_ref().clone());
    
    match message_service.mark_delivered(message_id, user_id).await {
        Ok(_) => Ok(HttpResponse::Ok().json(serde_json::json!({
            "success": true
        }))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}

/// Mark message as read
pub async fn mark_read(
    pool: web::Data<PgPool>,
    redis: web::Data<RedisClient>,
    path: web::Path<(Uuid, Uuid)>,
) -> ActixResult<HttpResponse> {
    let (user_id, message_id) = path.into_inner();
    let message_service = MessageService::new(pool.get_ref().clone(), redis.get_ref().clone());
    
    match message_service.mark_read(message_id, user_id).await {
        Ok(_) => Ok(HttpResponse::Ok().json(serde_json::json!({
            "success": true
        }))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}

/// Delete message
pub async fn delete_message(
    pool: web::Data<PgPool>,
    redis: web::Data<RedisClient>,
    path: web::Path<(Uuid, Uuid)>,
) -> ActixResult<HttpResponse> {
    let (user_id, message_id) = path.into_inner();
    let message_service = MessageService::new(pool.get_ref().clone(), redis.get_ref().clone());
    
    match message_service.delete_message(message_id, user_id).await {
        Ok(_) => Ok(HttpResponse::Ok().json(serde_json::json!({
            "success": true,
            "message": "Message deleted"
        }))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}

/// Get offline queue
pub async fn get_offline_queue(
    pool: web::Data<PgPool>,
    redis: web::Data<RedisClient>,
    path: web::Path<(Uuid, String)>,
) -> ActixResult<HttpResponse> {
    let (user_id, device_id) = path.into_inner();
    let message_service = MessageService::new(pool.get_ref().clone(), redis.get_ref().clone());
    
    match message_service.get_offline_queue(user_id, device_id.clone()).await {
        Ok(messages) => Ok(HttpResponse::Ok().json(serde_json::json!({
            "messages": messages,
            "count": messages.len()
        }))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}

/// Clear offline queue
pub async fn clear_offline_queue(
    pool: web::Data<PgPool>,
    redis: web::Data<RedisClient>,
    path: web::Path<(Uuid, String)>,
    message_ids: web::Json<Vec<Uuid>>,
) -> ActixResult<HttpResponse> {
    let (user_id, device_id) = path.into_inner();
    let message_service = MessageService::new(pool.get_ref().clone(), redis.get_ref().clone());
    
    match message_service.clear_offline_queue(user_id, device_id, message_ids.into_inner()).await {
        Ok(_) => Ok(HttpResponse::Ok().json(serde_json::json!({
            "success": true
        }))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}
