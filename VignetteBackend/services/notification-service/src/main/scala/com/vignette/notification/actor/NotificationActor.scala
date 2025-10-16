package com.vignette.notification.actor

import akka.actor.typed.{ActorRef, Behavior}
import akka.actor.typed.scaladsl.{ActorContext, Behaviors}
import java.time.Instant
import java.util.UUID

import com.vignette.notification.model._
import com.vignette.notification.service.NotificationService

/**
 * Main notification coordinator actor
 * Handles notification creation, delivery, and routing
 */
object NotificationActor {
  
  // Commands
  sealed trait Command
  
  case class SendNotification(
    notification: NotificationRequest,
    replyTo: ActorRef[NotificationResponse]
  ) extends Command
  
  case class SendBatch(
    batch: NotificationBatch,
    replyTo: ActorRef[BatchResponse]
  ) extends Command
  
  case class MarkAsRead(
    notificationId: UUID,
    userId: UUID,
    replyTo: ActorRef[OperationResult]
  ) extends Command
  
  case class MarkAllAsRead(
    userId: UUID,
    replyTo: ActorRef[OperationResult]
  ) extends Command
  
  case class GetNotifications(
    userId: UUID,
    unreadOnly: Boolean,
    limit: Int,
    offset: Int,
    replyTo: ActorRef[NotificationList]
  ) extends Command
  
  case class DeleteNotification(
    notificationId: UUID,
    userId: UUID,
    replyTo: ActorRef[OperationResult]
  ) extends Command
  
  // Responses
  case class NotificationResponse(notification: Notification)
  case class BatchResponse(success: Int, failed: Int)
  case class OperationResult(success: Boolean, message: String)
  case class NotificationList(notifications: Seq[Notification], total: Long)
  
  /**
   * Create notification actor behavior
   */
  def apply(
    notificationService: NotificationService,
    pushActor: ActorRef[PushNotificationActor.Command],
    emailActor: ActorRef[EmailActor.Command]
  ): Behavior[Command] = {
    Behaviors.setup { context =>
      context.log.info("NotificationActor started")
      active(notificationService, pushActor, emailActor, context)
    }
  }
  
  private def active(
    notificationService: NotificationService,
    pushActor: ActorRef[PushNotificationActor.Command],
    emailActor: ActorRef[EmailActor.Command],
    context: ActorContext[Command]
  ): Behavior[Command] = {
    Behaviors.receiveMessage {
      
      case SendNotification(request, replyTo) =>
        context.log.debug(s"Sending notification to user ${request.userId}")
        
        // Check user preferences
        val prefsOpt = notificationService.getPreferences(request.userId)
        
        val shouldSend = prefsOpt.exists { prefs =>
          // Check if notification type is enabled
          val typeEnabled = isNotificationTypeEnabled(request.notificationType, prefs)
          
          // Check quiet hours
          val inQuietHours = isInQuietHours(prefs)
          
          // Send if enabled and not in quiet hours (unless urgent)
          typeEnabled && (!inQuietHours || request.priority == NotificationPriority.Urgent)
        }
        
        if (shouldSend || prefsOpt.isEmpty) {
          // Check for grouping
          val existingGroup = notificationService.findGroupedNotification(
            request.userId,
            request.groupKey.getOrElse(""),
            5 // 5-minute window
          )
          
          val notification = existingGroup match {
            case Some(existing) =>
              // Update existing grouped notification
              notificationService.updateGroupedNotification(
                existing.id,
                request.actorId.get,
                request.actorUsername
              )
              
            case None =>
              // Create new notification
              val notif = Notification(
                id = UUID.randomUUID(),
                userId = request.userId,
                notificationType = request.notificationType,
                title = request.title,
                message = request.message,
                actorId = request.actorId,
                actorUsername = request.actorUsername,
                actorAvatarUrl = request.actorAvatarUrl,
                postId = request.postId,
                takeId = request.takeId,
                commentId = request.commentId,
                storyId = request.storyId,
                trendId = request.trendId,
                data = request.data,
                imageUrl = request.imageUrl,
                deepLink = request.deepLink,
                deliveryChannels = request.deliveryChannels,
                priority = request.priority,
                groupKey = request.groupKey,
                groupCount = 1
              )
              
              notificationService.saveNotification(notif)
              notif
          }
          
          // Send via appropriate channels
          if (notification.deliveryChannels.contains(DeliveryChannel.Push)) {
            pushActor ! PushNotificationActor.SendPush(notification.userId, notification)
          }
          
          if (notification.deliveryChannels.contains(DeliveryChannel.Email)) {
            emailActor ! EmailActor.SendEmail(notification.userId, notification)
          }
          
          replyTo ! NotificationResponse(notification)
          
        } else {
          context.log.debug(s"Notification not sent due to user preferences")
          // Still create notification but don't deliver
          val notif = Notification(
            id = UUID.randomUUID(),
            userId = request.userId,
            notificationType = request.notificationType,
            title = request.title,
            message = request.message,
            actorId = request.actorId,
            deliveryChannels = Set(DeliveryChannel.InApp)
          )
          notificationService.saveNotification(notif)
          replyTo ! NotificationResponse(notif)
        }
        
        Behaviors.same
      
      case SendBatch(batch, replyTo) =>
        context.log.debug(s"Sending batch of ${batch.notifications.size} notifications")
        
        var successCount = 0
        var failedCount = 0
        
        batch.notifications.foreach { request =>
          try {
            // Send each notification (could be parallelized)
            context.self ! SendNotification(request, context.system.ignoreRef)
            successCount += 1
          } catch {
            case e: Exception =>
              context.log.error(s"Failed to send notification: ${e.getMessage}")
              failedCount += 1
          }
        }
        
        replyTo ! BatchResponse(successCount, failedCount)
        Behaviors.same
      
      case MarkAsRead(notificationId, userId, replyTo) =>
        context.log.debug(s"Marking notification $notificationId as read")
        
        val result = notificationService.markAsRead(notificationId, userId)
        replyTo ! OperationResult(result, if (result) "Marked as read" else "Failed")
        
        Behaviors.same
      
      case MarkAllAsRead(userId, replyTo) =>
        context.log.debug(s"Marking all notifications as read for user $userId")
        
        val count = notificationService.markAllAsRead(userId)
        replyTo ! OperationResult(true, s"Marked $count notifications as read")
        
        Behaviors.same
      
      case GetNotifications(userId, unreadOnly, limit, offset, replyTo) =>
        context.log.debug(s"Getting notifications for user $userId")
        
        val (notifications, total) = notificationService.getNotifications(userId, unreadOnly, limit, offset)
        replyTo ! NotificationList(notifications, total)
        
        Behaviors.same
      
      case DeleteNotification(notificationId, userId, replyTo) =>
        context.log.debug(s"Deleting notification $notificationId")
        
        val result = notificationService.deleteNotification(notificationId, userId)
        replyTo ! OperationResult(result, if (result) "Deleted" else "Failed")
        
        Behaviors.same
    }
  }
  
  // Helper methods
  private def isNotificationTypeEnabled(
    notifType: NotificationType.NotificationType,
    prefs: NotificationPreferences
  ): Boolean = {
    notifType match {
      case NotificationType.Like => prefs.notifyOnLike
      case NotificationType.Comment => prefs.notifyOnComment
      case NotificationType.Follow => prefs.notifyOnFollow
      case NotificationType.Mention => prefs.notifyOnMention
      case NotificationType.Share => prefs.notifyOnShare
      case NotificationType.TakeRemix => prefs.notifyOnTakeRemix
      case NotificationType.ReplyToStory => prefs.notifyOnStoryReply
      case NotificationType.TaggedInPost | NotificationType.TaggedInTake => prefs.notifyOnTagged
      case _ => true // Default to enabled for new types
    }
  }
  
  private def isInQuietHours(prefs: NotificationPreferences): Boolean = {
    if (!prefs.quietHoursEnabled) return false
    
    (prefs.quietHoursStart, prefs.quietHoursEnd) match {
      case (Some(start), Some(end)) =>
        val now = java.time.LocalTime.now()
        val startTime = java.time.LocalTime.parse(start)
        val endTime = java.time.LocalTime.parse(end)
        
        if (startTime.isBefore(endTime)) {
          now.isAfter(startTime) && now.isBefore(endTime)
        } else {
          // Crosses midnight
          now.isAfter(startTime) || now.isBefore(endTime)
        }
      case _ => false
    }
  }
}
