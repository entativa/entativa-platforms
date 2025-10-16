use actix_web::{web, HttpResponse, Result as ActixResult};
use sqlx::PgPool;
use redis::Client as RedisClient;
use uuid::Uuid;
use serde::Deserialize;

use crate::services::group_service::GroupService;

#[derive(Deserialize)]
pub struct CreateGroupRequest {
    pub name: String,
    pub description: Option<String>,
    pub member_ids: Vec<Uuid>,
}

#[derive(Deserialize)]
pub struct AddMemberRequest {
    pub user_id: Uuid,
}

#[derive(Deserialize)]
pub struct SendGroupMessageRequest {
    pub ciphertext: String, // Base64 encoded
}

/// Create group
pub async fn create_group(
    pool: web::Data<PgPool>,
    redis: web::Data<RedisClient>,
    creator_id: web::Path<Uuid>,
    request: web::Json<CreateGroupRequest>,
) -> ActixResult<HttpResponse> {
    let group_service = GroupService::new(pool.get_ref().clone(), redis.get_ref().clone());
    let req = request.into_inner();
    
    match group_service.create_group(
        *creator_id,
        req.name,
        req.description,
        req.member_ids,
    ).await {
        Ok((group, _mls_group)) => Ok(HttpResponse::Created().json(group)),
        Err(e) => Ok(HttpResponse::BadRequest().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}

/// Add member to group
pub async fn add_member(
    pool: web::Data<PgPool>,
    redis: web::Data<RedisClient>,
    path: web::Path<(Uuid, Uuid)>,
    request: web::Json<AddMemberRequest>,
) -> ActixResult<HttpResponse> {
    let (group_id, added_by) = path.into_inner();
    let group_service = GroupService::new(pool.get_ref().clone(), redis.get_ref().clone());
    
    match group_service.add_member(group_id, request.user_id, added_by).await {
        Ok(_) => Ok(HttpResponse::Ok().json(serde_json::json!({
            "success": true,
            "message": "Member added"
        }))),
        Err(e) => Ok(HttpResponse::BadRequest().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}

/// Remove member from group
pub async fn remove_member(
    pool: web::Data<PgPool>,
    redis: web::Data<RedisClient>,
    path: web::Path<(Uuid, Uuid, Uuid)>,
) -> ActixResult<HttpResponse> {
    let (group_id, user_id, removed_by) = path.into_inner();
    let group_service = GroupService::new(pool.get_ref().clone(), redis.get_ref().clone());
    
    match group_service.remove_member(group_id, user_id, removed_by).await {
        Ok(_) => Ok(HttpResponse::Ok().json(serde_json::json!({
            "success": true,
            "message": "Member removed"
        }))),
        Err(e) => Ok(HttpResponse::BadRequest().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}

/// Send group message
pub async fn send_group_message(
    pool: web::Data<PgPool>,
    redis: web::Data<RedisClient>,
    path: web::Path<(Uuid, Uuid)>,
    request: web::Json<SendGroupMessageRequest>,
) -> ActixResult<HttpResponse> {
    let (group_id, sender_id) = path.into_inner();
    let group_service = GroupService::new(pool.get_ref().clone(), redis.get_ref().clone());
    
    // Decode ciphertext
    let ciphertext = match base64::decode(&request.ciphertext) {
        Ok(ct) => ct,
        Err(e) => return Ok(HttpResponse::BadRequest().json(serde_json::json!({
            "error": format!("Invalid base64: {}", e)
        }))),
    };
    
    match group_service.send_group_message(sender_id, group_id, ciphertext).await {
        Ok(response) => Ok(HttpResponse::Created().json(response)),
        Err(e) => Ok(HttpResponse::BadRequest().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}

/// Get group info
pub async fn get_group_info(
    pool: web::Data<PgPool>,
    group_id: web::Path<Uuid>,
) -> ActixResult<HttpResponse> {
    let row = sqlx::query!(
        r#"
        SELECT g.*, COUNT(gm.user_id) as member_count
        FROM group_chats g
        LEFT JOIN group_members gm ON g.id = gm.group_id
        WHERE g.id = $1
        GROUP BY g.id
        "#,
        *group_id
    )
    .fetch_one(pool.get_ref())
    .await;
    
    match row {
        Ok(group) => Ok(HttpResponse::Ok().json(serde_json::json!({
            "id": group.id,
            "name": group.name,
            "description": group.description,
            "created_by": group.created_by,
            "member_count": group.member_count.unwrap_or(0),
            "max_members": group.max_members,
            "current_epoch": group.current_epoch,
            "created_at": group.created_at,
        }))),
        Err(e) => Ok(HttpResponse::NotFound().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}

/// Get group members
pub async fn get_group_members(
    pool: web::Data<PgPool>,
    group_id: web::Path<Uuid>,
) -> ActixResult<HttpResponse> {
    let members = sqlx::query!(
        "SELECT * FROM group_members WHERE group_id = $1 ORDER BY joined_at ASC",
        *group_id
    )
    .fetch_all(pool.get_ref())
    .await;
    
    match members {
        Ok(rows) => Ok(HttpResponse::Ok().json(serde_json::json!({
            "members": rows,
            "count": rows.len()
        }))),
        Err(e) => Ok(HttpResponse::InternalServerError().json(serde_json::json!({
            "error": e.to_string()
        }))),
    }
}
