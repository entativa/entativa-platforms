package com.socialink.notification.actor

import akka.actor.typed.{ActorRef, Behavior}
import akka.actor.typed.scaladsl.{ActorContext, Behaviors}
import akka.stream.scaladsl.SourceQueueWithComplete
import java.util.UUID

import com.socialink.notification.model._

/**
 * Device registry actor
 * Manages WebSocket connections and real-time notification delivery
 */
object DeviceRegistry {
  
  // Commands
  sealed trait Command
  
  case class RegisterConnection(
    userId: UUID,
    queue: SourceQueueWithComplete[String]
  ) extends Command
  
  case class UnregisterConnection(
    userId: UUID
  ) extends Command
  
  case class SendToUser(
    userId: UUID,
    message: String
  ) extends Command
  
  case class BroadcastToUsers(
    userIds: Seq[UUID],
    message: String
  ) extends Command
  
  case class GetConnectedUsers(
    replyTo: ActorRef[ConnectedUsers]
  ) extends Command
  
  // Responses
  case class ConnectedUsers(userIds: Set[UUID], count: Int)
  
  /**
   * Create device registry actor
   */
  def apply(): Behavior[Command] = {
    Behaviors.setup { context =>
      context.log.info("DeviceRegistry started")
      active(Map.empty, context)
    }
  }
  
  private def active(
    connections: Map[UUID, SourceQueueWithComplete[String]],
    context: ActorContext[Command]
  ): Behavior[Command] = {
    Behaviors.receiveMessage {
      
      case RegisterConnection(userId, queue) =>
        context.log.info(s"Registering WebSocket connection for user $userId")
        
        // Close existing connection if any
        connections.get(userId).foreach { existingQueue =>
          context.log.debug(s"Closing existing connection for user $userId")
          existingQueue.complete()
        }
        
        // Register new connection
        val updatedConnections = connections + (userId -> queue)
        
        // Send welcome message
        queue.offer(s"""{"type":"connected","message":"WebSocket connected","userId":"$userId"}""")
        
        context.log.info(s"Total connected users: ${updatedConnections.size}")
        active(updatedConnections, context)
      
      case UnregisterConnection(userId) =>
        context.log.info(s"Unregistering WebSocket connection for user $userId")
        
        connections.get(userId).foreach(_.complete())
        val updatedConnections = connections - userId
        
        context.log.info(s"Total connected users: ${updatedConnections.size}")
        active(updatedConnections, context)
      
      case SendToUser(userId, message) =>
        connections.get(userId) match {
          case Some(queue) =>
            context.log.debug(s"Sending WebSocket message to user $userId")
            queue.offer(message) // Non-blocking offer
          case None =>
            context.log.debug(s"User $userId not connected via WebSocket")
        }
        
        Behaviors.same
      
      case BroadcastToUsers(userIds, message) =>
        context.log.debug(s"Broadcasting to ${userIds.size} users")
        
        var deliveredCount = 0
        userIds.foreach { userId =>
          connections.get(userId).foreach { queue =>
            queue.offer(message)
            deliveredCount += 1
          }
        }
        
        context.log.info(s"Broadcast delivered to $deliveredCount/${userIds.size} users")
        Behaviors.same
      
      case GetConnectedUsers(replyTo) =>
        replyTo ! ConnectedUsers(connections.keySet, connections.size)
        Behaviors.same
    }
  }
}

/**
 * WebSocket message types
 */
object WebSocketMessages {
  
  def notificationMessage(notification: Notification): String = {
    import spray.json._
    import NotificationJsonProtocol._
    
    val response = NotificationResponse(
      id = notification.id,
      userId = notification.userId,
      notificationType = notification.notificationType.toString,
      title = notification.title,
      message = notification.message,
      actorId = notification.actorId,
      actorUsername = notification.actorUsername,
      actorAvatarUrl = notification.actorAvatarUrl,
      postId = notification.postId,
      takeId = notification.takeId,
      commentId = notification.commentId,
      storyId = notification.storyId,
      trendId = notification.trendId,
      data = notification.data,
      imageUrl = notification.imageUrl,
      deepLink = notification.deepLink,
      isRead = notification.isRead,
      groupKey = notification.groupKey,
      groupCount = notification.groupCount,
      createdAt = notification.createdAt,
      readAt = notification.readAt
    )
    
    Map(
      "type" -> "notification",
      "data" -> response.toJson
    ).toJson.compactPrint
  }
  
  def pingMessage(): String = {
    """{"type":"ping","timestamp":""" + Instant.now().toString + "\"}"
  }
  
  def pongMessage(): String = {
    """{"type":"pong","timestamp":""" + Instant.now().toString + "\"}"
  }
  
  def errorMessage(error: String): String = {
    s"""{"type":"error","message":"$error"}"""
  }
}
