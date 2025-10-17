package com.socialink.notification

import com.typesafe.config.ConfigFactory

/**
 * Application configuration
 */
object Config {
  
  private val config = ConfigFactory.load()
  
  // Server
  val httpHost: String = config.getString("http.host")
  val httpPort: Int = config.getInt("http.port")
  
  // Database
  val dbUrl: String = config.getString("database.url")
  val dbUser: String = config.getString("database.user")
  val dbPassword: String = config.getString("database.password")
  val dbPoolSize: Int = config.getInt("database.pool-size")
  
  // Redis
  val redisHost: String = config.getString("redis.host")
  val redisPort: Int = config.getInt("redis.port")
  val redisPassword: Option[String] = if (config.hasPath("redis.password")) {
    Some(config.getString("redis.password"))
  } else None
  
  // Firebase (FCM)
  val fcmEnabled: Boolean = config.getBoolean("fcm.enabled")
  val fcmCredentialsPath: String = config.getString("fcm.credentials-path")
  
  // Apple Push Notifications (APN)
  val apnEnabled: Boolean = config.getBoolean("apn.enabled")
  val apnKeyId: String = config.getString("apn.key-id")
  val apnTeamId: String = config.getString("apn.team-id")
  val apnBundleId: String = config.getString("apn.bundle-id")
  
  // Email (SMTP)
  val emailEnabled: Boolean = config.getBoolean("email.enabled")
  val smtpHost: String = config.getString("email.smtp.host")
  val smtpPort: Int = config.getInt("email.smtp.port")
  val smtpUsername: String = config.getString("email.smtp.username")
  val smtpPassword: String = config.getString("email.smtp.password")
  val emailFrom: String = config.getString("email.from")
  
  // Kafka
  val kafkaEnabled: Boolean = config.getBoolean("kafka.enabled")
  val kafkaBootstrapServers: String = config.getString("kafka.bootstrap-servers")
  val kafkaTopic: String = config.getString("kafka.topic")
  
  // Notification settings
  val notificationGroupingWindowMinutes: Int = config.getInt("notification.grouping-window-minutes")
  val notificationExpiryDays: Int = config.getInt("notification.expiry-days")
}
