use actix::prelude::*;
use actix_web::{web, Error, HttpRequest, HttpResponse};
use actix_ws::Message as WsMessage;
use redis::AsyncCommands;
use std::collections::HashMap;
use std::sync::Arc;
use tokio::sync::RwLock;
use uuid::Uuid;

/// WebSocket connection manager
pub struct WsServer {
    /// Connected clients: user_id -> device_id -> Addr
    connections: Arc<RwLock<HashMap<Uuid, HashMap<String, Recipient<ServerMessage>>>>>,
    redis: redis::Client,
}

impl WsServer {
    pub fn new(redis: redis::Client) -> Self {
        Self {
            connections: Arc::new(RwLock::new(HashMap::new())),
            redis,
        }
    }
    
    /// Register new connection
    pub async fn register(
        &self,
        user_id: Uuid,
        device_id: String,
        addr: Recipient<ServerMessage>,
    ) {
        let mut connections = self.connections.write().await;
        connections
            .entry(user_id)
            .or_insert_with(HashMap::new)
            .insert(device_id, addr);
        
        tracing::info!("User {} device {} connected", user_id, device_id);
    }
    
    /// Unregister connection
    pub async fn unregister(&self, user_id: Uuid, device_id: &str) {
        let mut connections = self.connections.write().await;
        if let Some(devices) = connections.get_mut(&user_id) {
            devices.remove(device_id);
            if devices.is_empty() {
                connections.remove(&user_id);
            }
        }
        
        tracing::info!("User {} device {} disconnected", user_id, device_id);
    }
    
    /// Send message to user
    pub async fn send_to_user(&self, user_id: Uuid, message: ServerMessage) {
        let connections = self.connections.read().await;
        if let Some(devices) = connections.get(&user_id) {
            for (device_id, addr) in devices {
                if let Err(e) = addr.try_send(message.clone()) {
                    tracing::warn!("Failed to send to user {} device {}: {}", user_id, device_id, e);
                }
            }
        }
    }
    
    /// Send message to specific device
    pub async fn send_to_device(&self, user_id: Uuid, device_id: &str, message: ServerMessage) {
        let connections = self.connections.read().await;
        if let Some(devices) = connections.get(&user_id) {
            if let Some(addr) = devices.get(device_id) {
                if let Err(e) = addr.try_send(message) {
                    tracing::warn!("Failed to send to device: {}", e);
                }
            }
        }
    }
    
    /// Get connected user count
    pub async fn get_connection_count(&self) -> usize {
        let connections = self.connections.read().await;
        connections.len()
    }
    
    /// Check if user is online
    pub async fn is_user_online(&self, user_id: Uuid) -> bool {
        let connections = self.connections.read().await;
        connections.contains_key(&user_id)
    }
    
    /// Start Redis subscription for message events
    pub async fn start_redis_subscription(self: Arc<Self>) {
        tokio::spawn(async move {
            let mut pubsub = match self.redis.get_async_connection().await {
                Ok(conn) => conn.into_pubsub(),
                Err(e) => {
                    tracing::error!("Failed to connect to Redis: {}", e);
                    return;
                }
            };
            
            // Subscribe to all user message channels
            if let Err(e) = pubsub.psubscribe("messages:user:*").await {
                tracing::error!("Failed to subscribe to Redis: {}", e);
                return;
            }
            
            tracing::info!("WebSocket server listening to Redis pub/sub");
            
            let mut stream = pubsub.on_message();
            while let Some(msg) = stream.next().await {
                let payload: String = match msg.get_payload() {
                    Ok(p) => p,
                    Err(e) => {
                        tracing::warn!("Invalid Redis payload: {}", e);
                        continue;
                    }
                };
                
                // Parse message
                let message: serde_json::Value = match serde_json::from_str(&payload) {
                    Ok(m) => m,
                    Err(e) => {
                        tracing::warn!("Failed to parse message: {}", e);
                        continue;
                    }
                };
                
                // Extract recipient
                if let Some(recipient_id) = message.get("recipient_id").and_then(|r| r.as_str()) {
                    if let Ok(user_id) = Uuid::parse_str(recipient_id) {
                        // Send to user
                        self.send_to_user(user_id, ServerMessage::NewMessage(payload)).await;
                    }
                }
            }
        });
    }
}

/// Messages sent from server to client
#[derive(Message, Clone)]
#[rtype(result = "()")]
pub enum ServerMessage {
    NewMessage(String),           // New encrypted message
    DeliveryReceipt(String),      // Message delivered/read
    TypingIndicator(String),      // Someone is typing
    PresenceUpdate(String),       // User went online/offline
    Ping,
}

/// WebSocket session
pub struct WsSession {
    user_id: Uuid,
    device_id: String,
    server: Arc<WsServer>,
}

impl WsSession {
    pub fn new(user_id: Uuid, device_id: String, server: Arc<WsServer>) -> Self {
        Self {
            user_id,
            device_id,
            server,
        }
    }
}

impl Actor for WsSession {
    type Context = ws::WebsocketContext<Self>;
    
    fn started(&mut self, ctx: &mut Self::Context) {
        tracing::info!("WebSocket session started for user {}", self.user_id);
        
        // Register with server
        let server = self.server.clone();
        let user_id = self.user_id;
        let device_id = self.device_id.clone();
        let addr = ctx.address().recipient();
        
        ctx.spawn(async move {
            server.register(user_id, device_id, addr).await;
        }.into_actor(self));
        
        // Start heartbeat
        self.heartbeat(ctx);
    }
    
    fn stopped(&mut self, _ctx: &mut Self::Context) {
        tracing::info!("WebSocket session stopped for user {}", self.user_id);
        
        // Unregister from server
        let server = self.server.clone();
        let user_id = self.user_id;
        let device_id = self.device_id.clone();
        
        actix::spawn(async move {
            server.unregister(user_id, &device_id).await;
        });
    }
}

impl WsSession {
    fn heartbeat(&self, ctx: &mut ws::WebsocketContext<Self>) {
        ctx.run_interval(std::time::Duration::from_secs(30), |_act, ctx| {
            ctx.ping(b"");
        });
    }
}

impl Handler<ServerMessage> for WsSession {
    type Result = ();
    
    fn handle(&mut self, msg: ServerMessage, ctx: &mut Self::Context) {
        match msg {
            ServerMessage::NewMessage(payload) => {
                ctx.text(payload);
            }
            ServerMessage::DeliveryReceipt(payload) => {
                ctx.text(payload);
            }
            ServerMessage::TypingIndicator(payload) => {
                ctx.text(payload);
            }
            ServerMessage::PresenceUpdate(payload) => {
                ctx.text(payload);
            }
            ServerMessage::Ping => {
                ctx.ping(b"");
            }
        }
    }
}

impl StreamHandler<Result<WsMessage, ws::ProtocolError>> for WsSession {
    fn handle(&mut self, msg: Result<WsMessage, ws::ProtocolError>, ctx: &mut Self::Context) {
        match msg {
            Ok(WsMessage::Ping(msg)) => {
                ctx.pong(&msg);
            }
            Ok(WsMessage::Pong(_)) => {
                // Heartbeat received
            }
            Ok(WsMessage::Text(text)) => {
                tracing::debug!("Received WebSocket text: {}", text);
                
                // Parse command
                if let Ok(command) = serde_json::from_str::<serde_json::Value>(&text) {
                    if let Some(cmd_type) = command.get("type").and_then(|t| t.as_str()) {
                        match cmd_type {
                            "typing" => {
                                // Handle typing indicator
                                tracing::debug!("User {} is typing", self.user_id);
                            }
                            "read" => {
                                // Handle read receipt
                                tracing::debug!("User {} read message", self.user_id);
                            }
                            _ => {
                                tracing::debug!("Unknown command: {}", cmd_type);
                            }
                        }
                    }
                }
            }
            Ok(WsMessage::Binary(_)) => {
                tracing::warn!("Binary WebSocket messages not supported");
            }
            Ok(WsMessage::Close(reason)) => {
                tracing::info!("WebSocket close: {:?}", reason);
                ctx.stop();
            }
            _ => ctx.stop(),
        }
    }
}

/// WebSocket route handler
pub async fn ws_route(
    req: HttpRequest,
    stream: web::Payload,
    server: web::Data<Arc<WsServer>>,
    user_id: web::Path<(Uuid, String)>,
) -> Result<HttpResponse, Error> {
    let (user_id, device_id) = user_id.into_inner();
    
    let session = WsSession::new(user_id, device_id, server.get_ref().clone());
    
    ws::start(session, &req, stream)
}
