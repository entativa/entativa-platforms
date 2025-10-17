package com.socialink.notification

import akka.actor.typed.ActorSystem
import akka.actor.typed.scaladsl.Behaviors
import akka.http.scaladsl.Http
import akka.http.scaladsl.model.StatusCodes
import akka.http.scaladsl.server.Directives._
import akka.http.scaladsl.server.Route
import slick.jdbc.PostgresProfile.api._
import scala.concurrent.{ExecutionContext, Future}
import scala.util.{Success, Failure}

import com.socialink.notification.actor._
import com.socialink.notification.api._
import com.socialink.notification.service._
import com.socialink.notification.repository._

/**
 * Main application entry point
 */
object Main extends App {
  
  // Create actor system
  implicit val system: ActorSystem[Nothing] = ActorSystem(Behaviors.empty, "NotificationSystem")
  implicit val ec: ExecutionContext = system.executionContext
  
  println("""
    ╔══════════════════════════════════════════╗
    ║                                          ║
    ║   Socialink Notification Service 🔔       ║
    ║                                          ║
    ║   Real-time notifications with Akka      ║
    ║                                          ║
    ╚══════════════════════════════════════════╝
  """)
  
  // Initialize database
  println("📊 Connecting to PostgreSQL...")
  val db = Database.forURL(
    url = Config.dbUrl,
    user = Config.dbUser,
    password = Config.dbPassword,
    driver = "org.postgresql.Driver"
  )
  
  // Initialize repositories
  val notificationRepository = new NotificationRepository(db)
  val deviceRepository = new DeviceRepository(db)
  
  // Initialize services
  val fcmService = new FCMService()
  val apnService = new APNService()
  val emailService = new EmailService(
    Config.smtpHost,
    Config.smtpPort,
    Config.smtpUsername,
    Config.smtpPassword,
    Config.emailFrom
  )
  val notificationService = new NotificationService(notificationRepository)
  
  // Create actors
  println("🎭 Starting Akka actors...")
  val pushActor = system.systemActorOf(
    PushNotificationActor(fcmService, apnService),
    "push-notification-actor"
  )
  
  val emailActor = system.systemActorOf(
    EmailActor(emailService),
    "email-actor"
  )
  
  val deviceRegistry = system.systemActorOf(
    DeviceRegistry(),
    "device-registry"
  )
  
  val notificationActor = system.systemActorOf(
    NotificationActor(notificationService, pushActor, emailActor),
    "notification-actor"
  )
  
  // Create API routes
  val notificationRoutes = new NotificationRoutes(notificationActor)
  val subscriptionRoutes = new SubscriptionRoutes(deviceRepository)
  
  // Combined routes
  val routes: Route = concat(
    // Health check
    path("health") {
      get {
        complete(StatusCodes.OK, Map(
          "status" -> "healthy",
          "service" -> "Socialink Notification Service",
          "version" -> "1.0.0"
        ))
      }
    },
    
    // Root endpoint
    pathSingleSlash {
      get {
        complete(StatusCodes.OK, Map(
          "service" -> "Socialink Notification Service",
          "version" -> "1.0.0",
          "description" -> "Real-time notification service with Akka actors",
          "features" -> List(
            "Real-time WebSocket delivery",
            "Push notifications (FCM + APN)",
            "Email notifications",
            "Smart grouping",
            "Fine-grained preferences",
            "Multi-channel delivery"
          ),
          "endpoints" -> Map(
            "notifications" -> "/api/v1/notifications",
            "devices" -> "/api/v1/devices",
            "health" -> "/health"
          )
        ))
      }
    },
    
    // API routes
    notificationRoutes.routes,
    subscriptionRoutes.routes
  )
  
  // Start HTTP server
  println(s"🚀 Starting HTTP server on ${Config.httpHost}:${Config.httpPort}...")
  
  val binding: Future[Http.ServerBinding] = Http()
    .newServerAt(Config.httpHost, Config.httpPort)
    .bind(routes)
  
  binding.onComplete {
    case Success(binding) =>
      val address = binding.localAddress
      println(s"""
        ✅ Notification service started successfully!
        
        🌐 HTTP:      http://${address.getHostString}:${address.getPort}
        🔔 API:       http://${address.getHostString}:${address.getPort}/api/v1/notifications
        📱 Devices:   http://${address.getHostString}:${address.getPort}/api/v1/devices
        💚 Health:    http://${address.getHostString}:${address.getPort}/health
        
        🎭 Actors:
           - NotificationActor ✅
           - PushNotificationActor ✅
           - EmailActor ✅
           - DeviceRegistry ✅
        
        🔥 Ready to send notifications!
      """)
    
    case Failure(exception) =>
      println(s"❌ Failed to start server: ${exception.getMessage}")
      system.terminate()
  }
  
  // Graceful shutdown
  sys.addShutdownHook {
    println("\n🛑 Shutting down notification service...")
    
    binding.flatMap(_.unbind()).onComplete { _ =>
      db.close()
      system.terminate()
      println("✅ Notification service shut down successfully")
    }
  }
}
