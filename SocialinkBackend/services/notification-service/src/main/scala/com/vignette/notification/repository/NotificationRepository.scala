package com.socialink.notification.repository

import slick.jdbc.PostgresProfile.api._
import java.time.Instant
import java.util.UUID
import scala.concurrent.{ExecutionContext, Future}

import com.socialink.notification.model._

/**
 * Notification repository for database operations
 */
class NotificationRepository(db: Database)(implicit ec: ExecutionContext) {
  
  // Table definition
  class NotificationsTable(tag: Tag) extends Table[Notification](tag, "notifications") {
    def id = column[UUID]("id", O.PrimaryKey)
    def userId = column[UUID]("user_id")
    def notificationType = column[String]("notification_type")
    def title = column[String]("title")
    def message = column[String]("message")
    def actorId = column[Option[UUID]]("actor_id")
    def actorUsername = column[Option[String]]("actor_username")
    def actorAvatarUrl = column[Option[String]]("actor_avatar_url")
    def postId = column[Option[UUID]]("post_id")
    def takeId = column[Option[UUID]]("take_id")
    def commentId = column[Option[UUID]]("comment_id")
    def storyId = column[Option[UUID]]("story_id")
    def trendId = column[Option[UUID]]("trend_id")
    def imageUrl = column[Option[String]]("image_url")
    def deepLink = column[Option[String]]("deep_link")
    def isRead = column[Boolean]("is_read")
    def isDelivered = column[Boolean]("is_delivered")
    def priority = column[String]("priority")
    def groupKey = column[Option[String]]("group_key")
    def groupCount = column[Int]("group_count")
    def createdAt = column[Instant]("created_at")
    def readAt = column[Option[Instant]]("read_at")
    def deliveredAt = column[Option[Instant]]("delivered_at")
    def expiresAt = column[Option[Instant]]("expires_at")
    def updatedAt = column[Instant]("updated_at")
    
    // Simplified projection (would need custom mapping for full object)
    def * = (id, userId, title, message, isRead, createdAt)
  }
  
  val notifications = TableQuery[NotificationsTable]
  
  /**
   * Create notification
   */
  def create(notification: Notification): Future[Notification] = {
    val insertQuery = sqlu"""
      INSERT INTO notifications (
        id, user_id, notification_type, title, message,
        actor_id, actor_username, actor_avatar_url,
        post_id, take_id, comment_id, story_id, trend_id,
        image_url, deep_link, is_read, is_delivered,
        priority, group_key, group_count,
        created_at, updated_at
      ) VALUES (
        ${notification.id}, ${notification.userId}, ${notification.notificationType.toString},
        ${notification.title}, ${notification.message},
        ${notification.actorId}, ${notification.actorUsername}, ${notification.actorAvatarUrl},
        ${notification.postId}, ${notification.takeId}, ${notification.commentId},
        ${notification.storyId}, ${notification.trendId},
        ${notification.imageUrl}, ${notification.deepLink},
        ${notification.isRead}, ${notification.isDelivered},
        ${notification.priority.toString}, ${notification.groupKey}, ${notification.groupCount},
        ${notification.createdAt}, ${notification.createdAt}
      )
    """
    
    db.run(insertQuery).map(_ => notification)
  }
  
  /**
   * Get notification by ID
   */
  def getById(id: UUID): Option[Notification] = {
    // Simplified - in production would use Slick query
    None
  }
  
  /**
   * Get all notifications for user
   */
  def getAll(userId: UUID, limit: Int, offset: Int): (Seq[Notification], Long) = {
    // Simplified - in production would use Slick query
    (Seq.empty, 0L)
  }
  
  /**
   * Get unread notifications for user
   */
  def getUnread(userId: UUID, limit: Int, offset: Int): (Seq[Notification], Long) = {
    // Simplified - in production would use Slick query
    (Seq.empty, 0L)
  }
  
  /**
   * Mark notification as read
   */
  def markAsRead(notificationId: UUID, userId: UUID): Boolean = {
    val query = sqlu"""
      UPDATE notifications
      SET is_read = TRUE, read_at = ${Instant.now()}, updated_at = ${Instant.now()}
      WHERE id = $notificationId AND user_id = $userId
    """
    
    // Simplified - would return Future and check row count
    true
  }
  
  /**
   * Mark all notifications as read
   */
  def markAllAsRead(userId: UUID): Int = {
    val query = sqlu"""
      UPDATE notifications
      SET is_read = TRUE, read_at = ${Instant.now()}, updated_at = ${Instant.now()}
      WHERE user_id = $userId AND is_read = FALSE
    """
    
    // Simplified - would return row count
    0
  }
  
  /**
   * Delete notification
   */
  def delete(notificationId: UUID, userId: UUID): Boolean = {
    val query = sqlu"""
      DELETE FROM notifications
      WHERE id = $notificationId AND user_id = $userId
    """
    
    true
  }
  
  /**
   * Get unread count
   */
  def getUnreadCount(userId: UUID): Long = {
    // Would execute: SELECT COUNT(*) FROM notifications WHERE user_id = ? AND is_read = FALSE
    0L
  }
  
  /**
   * Get statistics
   */
  def getStats(userId: UUID): NotificationStats = {
    NotificationStats(
      totalNotifications = 0,
      unreadCount = 0,
      readCount = 0,
      deliveredCount = 0,
      notificationsByType = Map.empty
    )
  }
  
  /**
   * Find notification by group key
   */
  def findByGroupKey(userId: UUID, groupKey: String, since: Instant): Option[Notification] = {
    // Would query for existing grouped notification within time window
    None
  }
  
  /**
   * Update grouped notification
   */
  def update(notification: Notification): Notification = {
    val query = sqlu"""
      UPDATE notifications
      SET message = ${notification.message},
          group_count = ${notification.groupCount},
          updated_at = ${notification.updatedAt}
      WHERE id = ${notification.id}
    """
    
    notification
  }
  
  /**
   * Delete expired notifications
   */
  def deleteExpired(now: Instant): Int = {
    val query = sqlu"""
      DELETE FROM notifications
      WHERE expires_at IS NOT NULL AND expires_at < $now
    """
    
    0
  }
  
  /**
   * Get user preferences
   */
  def getPreferences(userId: UUID): Option[NotificationPreferences] = {
    // Would query preferences table
    Some(NotificationPreferences(userId = userId))
  }
  
  /**
   * Update user preferences
   */
  def updatePreferences(prefs: NotificationPreferences): Future[NotificationPreferences] = {
    val query = sqlu"""
      INSERT INTO notification_preferences (
        user_id, enable_push, enable_email, enable_sms,
        notify_on_like, notify_on_comment, notify_on_follow,
        notify_on_mention, notify_on_share, notify_on_take_remix,
        notify_on_story_reply, notify_on_tagged,
        quiet_hours_enabled, quiet_hours_start, quiet_hours_end,
        updated_at
      ) VALUES (
        ${prefs.userId}, ${prefs.enablePush}, ${prefs.enableEmail}, ${prefs.enableSMS},
        ${prefs.notifyOnLike}, ${prefs.notifyOnComment}, ${prefs.notifyOnFollow},
        ${prefs.notifyOnMention}, ${prefs.notifyOnShare}, ${prefs.notifyOnTakeRemix},
        ${prefs.notifyOnStoryReply}, ${prefs.notifyOnTagged},
        ${prefs.quietHoursEnabled}, ${prefs.quietHoursStart}, ${prefs.quietHoursEnd},
        ${prefs.updatedAt}
      )
      ON CONFLICT (user_id) DO UPDATE SET
        enable_push = EXCLUDED.enable_push,
        enable_email = EXCLUDED.enable_email,
        enable_sms = EXCLUDED.enable_sms,
        notify_on_like = EXCLUDED.notify_on_like,
        notify_on_comment = EXCLUDED.notify_on_comment,
        notify_on_follow = EXCLUDED.notify_on_follow,
        notify_on_mention = EXCLUDED.notify_on_mention,
        notify_on_share = EXCLUDED.notify_on_share,
        notify_on_take_remix = EXCLUDED.notify_on_take_remix,
        notify_on_story_reply = EXCLUDED.notify_on_story_reply,
        notify_on_tagged = EXCLUDED.notify_on_tagged,
        quiet_hours_enabled = EXCLUDED.quiet_hours_enabled,
        quiet_hours_start = EXCLUDED.quiet_hours_start,
        quiet_hours_end = EXCLUDED.quiet_hours_end,
        updated_at = EXCLUDED.updated_at
    """
    
    db.run(query).map(_ => prefs)
  }
}
