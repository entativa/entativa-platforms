package com.vignette.notification.service

import com.google.firebase.FirebaseApp
import com.google.firebase.messaging._
import scala.util.{Try, Success, Failure}

import com.vignette.notification.model.PushNotificationPayload

/**
 * Firebase Cloud Messaging service
 * Sends push notifications to Android devices
 */
class FCMService {
  
  /**
   * Send push notification via FCM
   */
  def sendPush(deviceToken: String, payload: PushNotificationPayload): Either[String, String] = {
    Try {
      // Build FCM message
      val notification = com.google.firebase.messaging.Notification.builder()
        .setTitle(payload.title)
        .setBody(payload.body)
        .setImage(payload.imageUrl.orNull)
        .build()
      
      // Build Android config
      val androidConfig = AndroidConfig.builder()
        .setPriority(if (payload.priority == "high") AndroidConfig.Priority.HIGH else AndroidConfig.Priority.NORMAL)
        .setNotification(AndroidNotification.builder()
          .setSound(payload.sound)
          .setClickAction(payload.deepLink.getOrElse("FLUTTER_NOTIFICATION_CLICK"))
          .build())
        .build()
      
      // Build message
      val messageBuilder = Message.builder()
        .setToken(deviceToken)
        .setNotification(notification)
        .setAndroidConfig(androidConfig)
      
      // Add custom data
      if (payload.data.nonEmpty) {
        val javaMap: java.util.Map[String, String] = scala.collection.JavaConverters.mapAsJavaMap(payload.data)
        messageBuilder.putAllData(javaMap)
      }
      
      val message = messageBuilder.build()
      
      // Send message
      val messageId = FirebaseMessaging.getInstance().send(message)
      messageId
      
    } match {
      case Success(messageId) => Right(messageId)
      case Failure(ex) => Left(s"FCM error: ${ex.getMessage}")
    }
  }
  
  /**
   * Send push to multiple devices (batch)
   */
  def sendMulticast(deviceTokens: Seq[String], payload: PushNotificationPayload): Either[String, MulticastMessage] = {
    Try {
      // Build notification
      val notification = com.google.firebase.messaging.Notification.builder()
        .setTitle(payload.title)
        .setBody(payload.body)
        .setImage(payload.imageUrl.orNull)
        .build()
      
      // Build Android config
      val androidConfig = AndroidConfig.builder()
        .setPriority(AndroidConfig.Priority.HIGH)
        .setNotification(AndroidNotification.builder()
          .setSound(payload.sound)
          .build())
        .build()
      
      // Build multicast message
      val messageBuilder = MulticastMessage.builder()
        .addAllTokens(scala.collection.JavaConverters.seqAsJavaList(deviceTokens))
        .setNotification(notification)
        .setAndroidConfig(androidConfig)
      
      // Add custom data
      if (payload.data.nonEmpty) {
        val javaMap: java.util.Map[String, String] = scala.collection.JavaConverters.mapAsJavaMap(payload.data)
        messageBuilder.putAllData(javaMap)
      }
      
      messageBuilder.build()
      
    } match {
      case Success(message) => Right(message)
      case Failure(ex) => Left(s"FCM multicast error: ${ex.getMessage}")
    }
  }
  
  /**
   * Validate FCM token
   */
  def validateToken(token: String): Boolean = {
    // Basic validation
    token.nonEmpty && token.length > 50
  }
}
