package com.socialink.notification.model

import java.time.Instant
import java.util.UUID
import spray.json._
import NotificationJsonProtocol._

/**
 * Device platform
 */
object DevicePlatform extends Enumeration {
  type DevicePlatform = Value
  val iOS, Android, Web = Value
}

/**
 * Device model for push notifications
 */
case class Device(
  id: UUID,
  userId: UUID,
  deviceToken: String, // FCM token or APN token
  platform: DevicePlatform.DevicePlatform,
  deviceName: Option[String] = None,
  deviceModel: Option[String] = None,
  osVersion: Option[String] = None,
  appVersion: Option[String] = None,
  isActive: Boolean = true,
  lastUsedAt: Instant = Instant.now(),
  createdAt: Instant = Instant.now(),
  updatedAt: Instant = Instant.now()
)

/**
 * Device registration request
 */
case class DeviceRegistration(
  userId: UUID,
  deviceToken: String,
  platform: DevicePlatform.DevicePlatform,
  deviceName: Option[String] = None,
  deviceModel: Option[String] = None,
  osVersion: Option[String] = None,
  appVersion: Option[String] = None
)

/**
 * Push notification payload
 */
case class PushNotificationPayload(
  title: String,
  body: String,
  imageUrl: Option[String] = None,
  deepLink: Option[String] = None,
  data: Map[String, String] = Map.empty,
  badge: Option[Int] = None,
  sound: String = "default",
  priority: String = "high"
)

/**
 * JSON formatters
 */
object DeviceJsonProtocol extends DefaultJsonProtocol {
  
  implicit object DevicePlatformFormat extends RootJsonFormat[DevicePlatform.DevicePlatform] {
    def write(obj: DevicePlatform.DevicePlatform): JsValue = JsString(obj.toString)
    def read(json: JsValue): DevicePlatform.DevicePlatform = json match {
      case JsString(s) => DevicePlatform.withName(s)
      case _ => throw DeserializationException("DevicePlatform expected")
    }
  }
  
  implicit val deviceFormat: RootJsonFormat[Device] = jsonFormat12(Device.apply)
  implicit val deviceRegistrationFormat: RootJsonFormat[DeviceRegistration] = jsonFormat7(DeviceRegistration.apply)
  implicit val pushNotificationPayloadFormat: RootJsonFormat[PushNotificationPayload] = jsonFormat8(PushNotificationPayload.apply)
}
