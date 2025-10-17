use actix_web::{web, HttpResponse, Result as ActixResult};
use redis::Client as RedisClient;
use uuid::Uuid;

use crate::services::typing_service::TypingService;

/// Set typing
pub async fn set_typing(
    redis: web::Data<RedisClient>,
    path: web::Path<(Uuid, Uuid)>,
) -> ActixResult<HttpResponse> {
    let (conversation_id, user_id) = path.into_inner();
    let typing_service = TypingService::new(redis.get_ref().clone());
    
    match typing_service.set_typing(conversation_id, user_id).await {
        Ok(_) => Ok(HttpResponse::Ok().json(serde_json::json!({
            "success": true
        }))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}

/// Clear typing
pub async fn clear_typing(
    redis: web::Data<RedisClient>,
    path: web::Path<(Uuid, Uuid)>,
) -> ActixResult<HttpResponse> {
    let (conversation_id, user_id) = path.into_inner();
    let typing_service = TypingService::new(redis.get_ref().clone());
    
    match typing_service.clear_typing(conversation_id, user_id).await {
        Ok(_) => Ok(HttpResponse::Ok().json(serde_json::json!({
            "success": true
        }))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}

/// Get typing users
pub async fn get_typing_users(
    redis: web::Data<RedisClient>,
    conversation_id: web::Path<Uuid>,
) -> ActixResult<HttpResponse> {
    let typing_service = TypingService::new(redis.get_ref().clone());
    
    match typing_service.get_typing_users(*conversation_id).await {
        Ok(users) => Ok(HttpResponse::Ok().json(serde_json::json!({
            "typing_users": users,
            "count": users.len()
        }))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}
