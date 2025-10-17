use actix_web::{web, HttpResponse, Result as ActixResult};
use sqlx::PgPool;
use redis::Client as RedisClient;
use uuid::Uuid;
use serde::Deserialize;

use crate::models::message::CallType;
use crate::services::call_service::CallService;

#[derive(Deserialize)]
pub struct InitiateCallRequest {
    pub conversation_id: Uuid,
    pub call_type: String, // Audio or Video
    pub sdp_offer: String,
}

#[derive(Deserialize)]
pub struct AnswerCallRequest {
    pub sdp_answer: String,
}

#[derive(Deserialize)]
pub struct AddIceCandidateRequest {
    pub candidate: String,
}

/// Initiate call
pub async fn initiate_call(
    pool: web::Data<PgPool>,
    redis: web::Data<RedisClient>,
    caller_id: web::Path<Uuid>,
    request: web::Json<InitiateCallRequest>,
) -> ActixResult<HttpResponse> {
    let call_service = CallService::new(pool.get_ref().clone(), redis.get_ref().clone());
    
    let call_type = match request.call_type.as_str() {
        "Audio" => CallType::Audio,
        "Video" => CallType::Video,
        _ => return Ok(HttpResponse::BadRequest().json(serde_json::json!({
            "error": "Invalid call type"
        }))),
    };
    
    match call_service.initiate_call(
        *caller_id,
        request.conversation_id,
        call_type,
        request.sdp_offer.clone(),
    ).await {
        Ok(call) => Ok(HttpResponse::Created().json(call)),
        Err(e) => Ok(HttpResponse::BadRequest().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}

/// Answer call
pub async fn answer_call(
    pool: web::Data<PgPool>,
    redis: web::Data<RedisClient>,
    path: web::Path<(Uuid, Uuid)>,
    request: web::Json<AnswerCallRequest>,
) -> ActixResult<HttpResponse> {
    let (call_id, user_id) = path.into_inner();
    let call_service = CallService::new(pool.get_ref().clone(), redis.get_ref().clone());
    
    match call_service.answer_call(call_id, user_id, request.sdp_answer.clone()).await {
        Ok(_) => Ok(HttpResponse::Ok().json(serde_json::json!({
            "success": true,
            "status": "answered"
        }))),
        Err(e) => Ok(HttpResponse::BadRequest().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}

/// Decline call
pub async fn decline_call(
    pool: web::Data<PgPool>,
    redis: web::Data<RedisClient>,
    path: web::Path<(Uuid, Uuid)>,
) -> ActixResult<HttpResponse> {
    let (call_id, user_id) = path.into_inner();
    let call_service = CallService::new(pool.get_ref().clone(), redis.get_ref().clone());
    
    match call_service.decline_call(call_id, user_id).await {
        Ok(_) => Ok(HttpResponse::Ok().json(serde_json::json!({
            "success": true,
            "status": "declined"
        }))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}

/// End call
pub async fn end_call(
    pool: web::Data<PgPool>,
    redis: web::Data<RedisClient>,
    path: web::Path<(Uuid, Uuid)>,
) -> ActixResult<HttpResponse> {
    let (call_id, user_id) = path.into_inner();
    let call_service = CallService::new(pool.get_ref().clone(), redis.get_ref().clone());
    
    match call_service.end_call(call_id, user_id).await {
        Ok(_) => Ok(HttpResponse::Ok().json(serde_json::json!({
            "success": true,
            "status": "ended"
        }))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}

/// Add ICE candidate
pub async fn add_ice_candidate(
    pool: web::Data<PgPool>,
    redis: web::Data<RedisClient>,
    path: web::Path<(Uuid, Uuid)>,
    request: web::Json<AddIceCandidateRequest>,
) -> ActixResult<HttpResponse> {
    let (call_id, user_id) = path.into_inner();
    let call_service = CallService::new(pool.get_ref().clone(), redis.get_ref().clone());
    
    match call_service.add_ice_candidate(call_id, user_id, request.candidate.clone()).await {
        Ok(_) => Ok(HttpResponse::Ok().json(serde_json::json!({
            "success": true
        }))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}

/// Get ICE candidates
pub async fn get_ice_candidates(
    pool: web::Data<PgPool>,
    redis: web::Data<RedisClient>,
    call_id: web::Path<Uuid>,
) -> ActixResult<HttpResponse> {
    let call_service = CallService::new(pool.get_ref().clone(), redis.get_ref().clone());
    
    match call_service.get_ice_candidates(*call_id).await {
        Ok(candidates) => Ok(HttpResponse::Ok().json(serde_json::json!({
            "candidates": candidates,
            "count": candidates.len()
        }))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}

/// Get call history
pub async fn get_call_history(
    pool: web::Data<PgPool>,
    redis: web::Data<RedisClient>,
    conversation_id: web::Path<Uuid>,
) -> ActixResult<HttpResponse> {
    let call_service = CallService::new(pool.get_ref().clone(), redis.get_ref().clone());
    
    match call_service.get_call_history(*conversation_id, 20).await {
        Ok(calls) => Ok(HttpResponse::Ok().json(serde_json::json!({
            "calls": calls,
            "count": calls.len()
        }))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}
