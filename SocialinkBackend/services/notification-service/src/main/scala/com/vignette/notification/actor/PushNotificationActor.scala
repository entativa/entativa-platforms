package com.socialink.notification.actor

import akka.actor.typed.{ActorRef, Behavior}
import akka.actor.typed.scaladsl.{ActorContext, Behaviors}
import java.util.UUID

import com.socialink.notification.model._
import com.socialink.notification.service.{FCMService, APNService}

/**
 * Push notification actor
 * Sends push notifications via FCM (Android) and APN (iOS)
 */
object PushNotificationActor {
  
  // Commands
  sealed trait Command
  
  case class SendPush(
    userId: UUID,
    notification: Notification
  ) extends Command
  
  case class SendPushToDevice(
    device: Device,
    payload: PushNotificationPayload
  ) extends Command
  
  case class BatchSendPush(
    userIds: Seq[UUID],
    payload: PushNotificationPayload
  ) extends Command
  
  private case class PushResult(success: Boolean, deviceId: UUID, error: Option[String]) extends Command
  
  /**
   * Create push notification actor
   */
  def apply(
    fcmService: FCMService,
    apnService: APNService
  ): Behavior[Command] = {
    Behaviors.setup { context =>
      context.log.info("PushNotificationActor started")
      active(fcmService, apnService, context)
    }
  }
  
  private def active(
    fcmService: FCMService,
    apnService: APNService,
    context: ActorContext[Command]
  ): Behavior[Command] = {
    Behaviors.receiveMessage {
      
      case SendPush(userId, notification) =>
        context.log.debug(s"Sending push notification to user $userId")
        
        // Get user's devices (would come from DeviceRepository)
        // For now, we'll simulate
        
        // Create push payload
        val payload = PushNotificationPayload(
          title = notification.title,
          body = notification.message,
          imageUrl = notification.imageUrl,
          deepLink = notification.deepLink,
          data = notification.data.getOrElse(Map.empty).map { case (k, v) => k -> v.toString },
          badge = Some(1), // Would calculate unread count
          sound = "default",
          priority = notification.priority match {
            case NotificationPriority.Urgent => "high"
            case NotificationPriority.High => "high"
            case _ => "normal"
          }
        )
        
        // Send to all user devices (async)
        context.log.info(s"Push notification queued for user $userId")
        
        Behaviors.same
      
      case SendPushToDevice(device, payload) =>
        context.log.debug(s"Sending push to device ${device.id} (${device.platform})")
        
        // Route to appropriate service
        device.platform match {
          case DevicePlatform.Android =>
            fcmService.sendPush(device.deviceToken, payload) match {
              case Right(_) =>
                context.log.info(s"FCM push sent successfully to device ${device.id}")
                context.self ! PushResult(success = true, device.id, None)
              case Left(error) =>
                context.log.error(s"FCM push failed: $error")
                context.self ! PushResult(success = false, device.id, Some(error))
            }
          
          case DevicePlatform.iOS =>
            apnService.sendPush(device.deviceToken, payload) match {
              case Right(_) =>
                context.log.info(s"APN push sent successfully to device ${device.id}")
                context.self ! PushResult(success = true, device.id, None)
              case Left(error) =>
                context.log.error(s"APN push failed: $error")
                context.self ! PushResult(success = false, device.id, Some(error))
            }
          
          case DevicePlatform.Web =>
            // Web push notifications (Web Push API)
            context.log.info(s"Web push not yet implemented for device ${device.id}")
        }
        
        Behaviors.same
      
      case BatchSendPush(userIds, payload) =>
        context.log.debug(s"Batch sending push to ${userIds.size} users")
        
        // Send to each user (would be parallelized in production)
        userIds.foreach { userId =>
          // Queue push for each user
          context.log.info(s"Queued push for user $userId")
        }
        
        Behaviors.same
      
      case PushResult(success, deviceId, error) =>
        if (success) {
          context.log.debug(s"Push delivered successfully to device $deviceId")
        } else {
          context.log.warning(s"Push failed for device $deviceId: ${error.getOrElse("Unknown error")}")
          // Could mark device as inactive if multiple failures
        }
        
        Behaviors.same
    }
  }
}
