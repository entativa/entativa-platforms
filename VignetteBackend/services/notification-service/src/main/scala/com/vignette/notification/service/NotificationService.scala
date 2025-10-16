package com.vignette.notification.service

import java.time.Instant
import java.time.temporal.ChronoUnit
import java.util.UUID
import scala.concurrent.{ExecutionContext, Future}

import com.vignette.notification.model._
import com.vignette.notification.repository.NotificationRepository

/**
 * Core notification business logic
 */
class NotificationService(
  repository: NotificationRepository
)(implicit ec: ExecutionContext) {
  
  /**
   * Save a notification to database
   */
  def saveNotification(notification: Notification): Future[Notification] = {
    repository.create(notification)
  }
  
  /**
   * Get notifications for a user
   */
  def getNotifications(
    userId: UUID,
    unreadOnly: Boolean,
    limit: Int,
    offset: Int
  ): (Seq[Notification], Long) = {
    if (unreadOnly) {
      repository.getUnread(userId, limit, offset)
    } else {
      repository.getAll(userId, limit, offset)
    }
  }
  
  /**
   * Get single notification
   */
  def getNotification(notificationId: UUID): Option[Notification] = {
    repository.getById(notificationId)
  }
  
  /**
   * Mark notification as read
   */
  def markAsRead(notificationId: UUID, userId: UUID): Boolean = {
    repository.markAsRead(notificationId, userId)
  }
  
  /**
   * Mark all notifications as read for user
   */
  def markAllAsRead(userId: UUID): Int = {
    repository.markAllAsRead(userId)
  }
  
  /**
   * Delete a notification
   */
  def deleteNotification(notificationId: UUID, userId: UUID): Boolean = {
    repository.delete(notificationId, userId)
  }
  
  /**
   * Get unread count
   */
  def getUnreadCount(userId: UUID): Long = {
    repository.getUnreadCount(userId)
  }
  
  /**
   * Get notification statistics
   */
  def getStats(userId: UUID): NotificationStats = {
    repository.getStats(userId)
  }
  
  /**
   * Find existing grouped notification
   */
  def findGroupedNotification(
    userId: UUID,
    groupKey: String,
    windowMinutes: Int
  ): Option[Notification] = {
    if (groupKey.isEmpty) return None
    
    val since = Instant.now().minus(windowMinutes, ChronoUnit.MINUTES)
    repository.findByGroupKey(userId, groupKey, since)
  }
  
  /**
   * Update grouped notification
   */
  def updateGroupedNotification(
    notificationId: UUID,
    newActorId: UUID,
    newActorUsername: Option[String]
  ): Notification = {
    val notification = repository.getById(notificationId).get
    
    // Build updated message
    val groupCount = notification.groupCount + 1
    val message = buildGroupedMessage(
      notification.notificationType,
      groupCount,
      notification.actorUsername,
      newActorUsername
    )
    
    val updated = notification.copy(
      groupCount = groupCount,
      message = message,
      updatedAt = Instant.now()
    )
    
    repository.update(updated)
    updated
  }
  
  /**
   * Get user preferences
   */
  def getPreferences(userId: UUID): Option[NotificationPreferences] = {
    repository.getPreferences(userId)
  }
  
  /**
   * Update user preferences
   */
  def updatePreferences(prefs: NotificationPreferences): Future[NotificationPreferences] = {
    repository.updatePreferences(prefs)
  }
  
  /**
   * Clean up expired notifications
   */
  def cleanupExpired(): Int = {
    val now = Instant.now()
    repository.deleteExpired(now)
  }
  
  /**
   * Build grouped message
   */
  private def buildGroupedMessage(
    notifType: NotificationType.NotificationType,
    groupCount: Int,
    firstActor: Option[String],
    newActor: Option[String]
  ): String = {
    val action = notifType match {
      case NotificationType.Like => "liked"
      case NotificationType.Comment => "commented on"
      case NotificationType.Follow => "followed"
      case NotificationType.Share => "shared"
      case NotificationType.TakeRemix => "remixed"
      case NotificationType.TaggedInPost => "tagged you in"
      case _ => "interacted with"
    }
    
    val entity = notifType match {
      case NotificationType.TakeRemix => "your Take"
      case NotificationType.Comment => "your post"
      case NotificationType.TaggedInPost => "a post"
      case NotificationType.TaggedInTake => "a Take"
      case _ => "your content"
    }
    
    if (groupCount == 2) {
      s"${firstActor.getOrElse("Someone")} and ${newActor.getOrElse("someone else")} $action $entity"
    } else if (groupCount <= 3) {
      s"${firstActor.getOrElse("Someone")} and ${groupCount - 1} others $action $entity"
    } else {
      s"${firstActor.getOrElse("Someone")} and ${groupCount - 1} others $action $entity"
    }
  }
}
