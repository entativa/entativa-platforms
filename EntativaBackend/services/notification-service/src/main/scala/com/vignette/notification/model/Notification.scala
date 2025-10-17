package com.socialink.notification.model

import java.time.Instant
import java.util.UUID
import spray.json._

/**
 * Notification types
 */
object NotificationType extends Enumeration {
  type NotificationType = Value
  val Like, Comment, Follow, Mention, Share, TakeRemix, TrendJoin, 
      ReplyToStory, ReactionToStory, NewFollower, TaggedInPost, 
      TaggedInTake, CommentReply, QuizAnswer, PollVote, 
      CountdownReminder, BTTCreated, TemplateUsed = Value
}

/**
 * Notification priority
 */
object NotificationPriority extends Enumeration {
  type NotificationPriority = Value
  val Low, Normal, High, Urgent = Value
}

/**
 * Notification delivery channels
 */
object DeliveryChannel extends Enumeration {
  type DeliveryChannel = Value
  val InApp, Push, Email, SMS, WebSocket = Value
}

/**
 * Core notification model
 */
case class Notification(
  id: UUID,
  userId: UUID, // Recipient
  notificationType: NotificationType.NotificationType,
  title: String,
  message: String,
  actorId: Option[UUID], // Who triggered the notification
  actorUsername: Option[String],
  actorAvatarUrl: Option[String],
  
  // Related entities
  postId: Option[UUID] = None,
  takeId: Option[UUID] = None,
  commentId: Option[UUID] = None,
  storyId: Option[UUID] = None,
  trendId: Option[UUID] = None,
  
  // Metadata
  data: Option[Map[String, String]] = None,
  imageUrl: Option[String] = None,
  deepLink: Option[String] = None,
  
  // Status
  isRead: Boolean = false,
  isDelivered: Boolean = false,
  deliveryChannels: Set[DeliveryChannel.DeliveryChannel] = Set(DeliveryChannel.InApp),
  priority: NotificationPriority.NotificationPriority = NotificationPriority.Normal,
  
  // Grouping
  groupKey: Option[String] = None, // For notification grouping
  groupCount: Int = 1,
  
  // Timestamps
  createdAt: Instant = Instant.now(),
  readAt: Option[Instant] = None,
  deliveredAt: Option[Instant] = None,
  expiresAt: Option[Instant] = None
)

/**
 * Notification create request
 */
case class NotificationRequest(
  userId: UUID,
  notificationType: NotificationType.NotificationType,
  title: String,
  message: String,
  actorId: Option[UUID] = None,
  actorUsername: Option[String] = None,
  actorAvatarUrl: Option[String] = None,
  postId: Option[UUID] = None,
  takeId: Option[UUID] = None,
  commentId: Option[UUID] = None,
  storyId: Option[UUID] = None,
  trendId: Option[UUID] = None,
  data: Option[Map[String, String]] = None,
  imageUrl: Option[String] = None,
  deepLink: Option[String] = None,
  deliveryChannels: Set[DeliveryChannel.DeliveryChannel] = Set(DeliveryChannel.InApp),
  priority: NotificationPriority.NotificationPriority = NotificationPriority.Normal,
  groupKey: Option[String] = None
)

/**
 * Notification batch request (for grouping)
 */
case class NotificationBatch(
  notifications: Seq[NotificationRequest]
)

/**
 * Notification response
 */
case class NotificationResponse(
  id: UUID,
  userId: UUID,
  notificationType: String,
  title: String,
  message: String,
  actorId: Option[UUID],
  actorUsername: Option[String],
  actorAvatarUrl: Option[String],
  postId: Option[UUID],
  takeId: Option[UUID],
  commentId: Option[UUID],
  storyId: Option[UUID],
  trendId: Option[UUID],
  data: Option[Map[String, String]],
  imageUrl: Option[String],
  deepLink: Option[String],
  isRead: Boolean,
  groupKey: Option[String],
  groupCount: Int,
  createdAt: Instant,
  readAt: Option[Instant]
)

/**
 * Notification preferences
 */
case class NotificationPreferences(
  userId: UUID,
  enablePush: Boolean = true,
  enableEmail: Boolean = true,
  enableSMS: Boolean = false,
  
  // Fine-grained preferences
  notifyOnLike: Boolean = true,
  notifyOnComment: Boolean = true,
  notifyOnFollow: Boolean = true,
  notifyOnMention: Boolean = true,
  notifyOnShare: Boolean = true,
  notifyOnTakeRemix: Boolean = true,
  notifyOnStoryReply: Boolean = true,
  notifyOnTagged: Boolean = true,
  
  // Quiet hours
  quietHoursEnabled: Boolean = false,
  quietHoursStart: Option[String] = None, // HH:MM format
  quietHoursEnd: Option[String] = None,
  
  updatedAt: Instant = Instant.now()
)

/**
 * Notification statistics
 */
case class NotificationStats(
  totalNotifications: Long,
  unreadCount: Long,
  readCount: Long,
  deliveredCount: Long,
  notificationsByType: Map[String, Long]
)

/**
 * JSON formatters
 */
object NotificationJsonProtocol extends DefaultJsonProtocol {
  
  implicit object NotificationTypeFormat extends RootJsonFormat[NotificationType.NotificationType] {
    def write(obj: NotificationType.NotificationType): JsValue = JsString(obj.toString)
    def read(json: JsValue): NotificationType.NotificationType = json match {
      case JsString(s) => NotificationType.withName(s)
      case _ => throw DeserializationException("NotificationType expected")
    }
  }
  
  implicit object NotificationPriorityFormat extends RootJsonFormat[NotificationPriority.NotificationPriority] {
    def write(obj: NotificationPriority.NotificationPriority): JsValue = JsString(obj.toString)
    def read(json: JsValue): NotificationPriority.NotificationPriority = json match {
      case JsString(s) => NotificationPriority.withName(s)
      case _ => throw DeserializationException("NotificationPriority expected")
    }
  }
  
  implicit object DeliveryChannelFormat extends RootJsonFormat[DeliveryChannel.DeliveryChannel] {
    def write(obj: DeliveryChannel.DeliveryChannel): JsValue = JsString(obj.toString)
    def read(json: JsValue): DeliveryChannel.DeliveryChannel = json match {
      case JsString(s) => DeliveryChannel.withName(s)
      case _ => throw DeserializationException("DeliveryChannel expected")
    }
  }
  
  implicit object InstantFormat extends RootJsonFormat[Instant] {
    def write(obj: Instant): JsValue = JsString(obj.toString)
    def read(json: JsValue): Instant = json match {
      case JsString(s) => Instant.parse(s)
      case _ => throw DeserializationException("Instant expected")
    }
  }
  
  implicit object UUIDFormat extends RootJsonFormat[UUID] {
    def write(obj: UUID): JsValue = JsString(obj.toString)
    def read(json: JsValue): UUID = json match {
      case JsString(s) => UUID.fromString(s)
      case _ => throw DeserializationException("UUID expected")
    }
  }
  
  implicit val notificationRequestFormat: RootJsonFormat[NotificationRequest] = jsonFormat17(NotificationRequest.apply)
  implicit val notificationBatchFormat: RootJsonFormat[NotificationBatch] = jsonFormat1(NotificationBatch.apply)
  implicit val notificationResponseFormat: RootJsonFormat[NotificationResponse] = jsonFormat20(NotificationResponse.apply)
  implicit val notificationPreferencesFormat: RootJsonFormat[NotificationPreferences] = jsonFormat16(NotificationPreferences.apply)
  implicit val notificationStatsFormat: RootJsonFormat[NotificationStats] = jsonFormat5(NotificationStats.apply)
}
