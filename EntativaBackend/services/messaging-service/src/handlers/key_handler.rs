use actix_web::{web, HttpResponse, Result as ActixResult};
use sqlx::PgPool;
use uuid::Uuid;

use crate::models::keys::*;
use crate::services::key_service::KeyService;

/// Register device with keys
pub async fn register_device(
    pool: web::Data<PgPool>,
    user_id: web::Path<Uuid>,
    request: web::Json<KeyRegistrationRequest>,
) -> ActixResult<HttpResponse> {
    let key_service = KeyService::new(pool.get_ref().clone());
    
    match key_service.register_device(*user_id, request.into_inner()).await {
        Ok(device) => Ok(HttpResponse::Created().json(device)),
        Err(e) => Ok(HttpResponse::BadRequest().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}

/// Get pre-key bundle for user
pub async fn get_prekey_bundle(
    pool: web::Data<PgPool>,
    user_id: web::Path<Uuid>,
    device_id: web::Query<Option<String>>,
) -> ActixResult<HttpResponse> {
    let key_service = KeyService::new(pool.get_ref().clone());
    
    match key_service.get_prekey_bundle(*user_id, device_id.0.clone()).await {
        Ok(bundle) => Ok(HttpResponse::Ok().json(bundle)),
        Err(e) => Ok(HttpResponse::NotFound().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}

/// Rotate signed pre-key
pub async fn rotate_signed_prekey(
    pool: web::Data<PgPool>,
    path: web::Path<(Uuid, String)>,
    request: web::Json<SignedPreKeyUpload>,
) -> ActixResult<HttpResponse> {
    let (user_id, device_id) = path.into_inner();
    let key_service = KeyService::new(pool.get_ref().clone());
    
    match key_service.rotate_signed_prekey(user_id, device_id, request.into_inner()).await {
        Ok(_) => Ok(HttpResponse::Ok().json(serde_json::json!({
            "success": true,
            "message": "Signed pre-key rotated"
        }))),
        Err(e) => Ok(HttpResponse::BadRequest().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}

/// Upload one-time pre-keys
pub async fn upload_onetime_prekeys(
    pool: web::Data<PgPool>,
    path: web::Path<(Uuid, String)>,
    request: web::Json<Vec<OneTimePreKeyUpload>>,
) -> ActixResult<HttpResponse> {
    let (user_id, device_id) = path.into_inner();
    let key_service = KeyService::new(pool.get_ref().clone());
    
    match key_service.upload_onetime_prekeys(user_id, device_id, request.into_inner()).await {
        Ok(_) => Ok(HttpResponse::Ok().json(serde_json::json!({
            "success": true,
            "message": "One-time pre-keys uploaded"
        }))),
        Err(e) => Ok(HttpResponse::BadRequest().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}

/// Deactivate device
pub async fn deactivate_device(
    pool: web::Data<PgPool>,
    path: web::Path<(Uuid, String)>,
) -> ActixResult<HttpResponse> {
    let (user_id, device_id) = path.into_inner();
    let key_service = KeyService::new(pool.get_ref().clone());
    
    match key_service.deactivate_device(user_id, device_id).await {
        Ok(_) => Ok(HttpResponse::Ok().json(serde_json::json!({
            "success": true,
            "message": "Device deactivated"
        }))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}

/// Get user's devices
pub async fn get_user_devices(
    pool: web::Data<PgPool>,
    user_id: web::Path<Uuid>,
) -> ActixResult<HttpResponse> {
    let key_service = KeyService::new(pool.get_ref().clone());
    
    match key_service.get_user_devices(*user_id).await {
        Ok(devices) => Ok(HttpResponse::Ok().json(serde_json::json!({
            "devices": devices,
            "count": devices.len()
        }))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}

/// Get key statistics
pub async fn get_key_stats(
    pool: web::Data<PgPool>,
    path: web::Path<(Uuid, String)>,
) -> ActixResult<HttpResponse> {
    let (user_id, device_id) = path.into_inner();
    let key_service = KeyService::new(pool.get_ref().clone());
    
    match key_service.get_key_stats(user_id, device_id).await {
        Ok(stats) => Ok(HttpResponse::Ok().json(stats)),
        Err(e) => Ok(HttpResponse::InternalServerError().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}
