package com.socialink.notification.api

import akka.actor.typed.{ActorRef, ActorSystem}
import akka.actor.typed.scaladsl.AskPattern._
import akka.http.scaladsl.model.StatusCodes
import akka.http.scaladsl.server.Directives._
import akka.http.scaladsl.server.Route
import akka.util.Timeout
import scala.concurrent.duration._
import scala.concurrent.{ExecutionContext, Future}
import java.util.UUID

import com.socialink.notification.actor.NotificationActor
import com.socialink.notification.model._
import NotificationJsonProtocol._

/**
 * Notification API routes
 */
class NotificationRoutes(
  notificationActor: ActorRef[NotificationActor.Command]
)(implicit system: ActorSystem[_], ec: ExecutionContext) {
  
  implicit val timeout: Timeout = 30.seconds
  
  val routes: Route = pathPrefix("api" / "v1" / "notifications") {
    concat(
      // Get notifications
      pathEnd {
        get {
          parameters("user_id".as[UUID], "unread_only".as[Boolean].optional, "limit".as[Int].optional, "offset".as[Int].optional) {
            (userId, unreadOnlyOpt, limitOpt, offsetOpt) =>
              val unreadOnly = unreadOnlyOpt.getOrElse(false)
              val limit = limitOpt.getOrElse(20)
              val offset = offsetOpt.getOrElse(0)
              
              val result: Future[NotificationActor.NotificationList] = notificationActor.ask(ref =>
                NotificationActor.GetNotifications(userId, unreadOnly, limit, offset, ref)
              )
              
              onSuccess(result) { list =>
                complete(StatusCodes.OK, Map(
                  "notifications" -> list.notifications,
                  "total" -> list.total,
                  "count" -> list.notifications.size
                ))
              }
          }
        }
      },
      
      // Send notification (internal endpoint)
      pathEnd {
        post {
          entity(as[NotificationRequest]) { request =>
            val result: Future[NotificationActor.NotificationResponse] = notificationActor.ask(ref =>
              NotificationActor.SendNotification(request, ref)
            )
            
            onSuccess(result) { response =>
              complete(StatusCodes.Created, response.notification)
            }
          }
        }
      },
      
      // Send batch notifications
      path("batch") {
        post {
          entity(as[NotificationBatch]) { batch =>
            val result: Future[NotificationActor.BatchResponse] = notificationActor.ask(ref =>
              NotificationActor.SendBatch(batch, ref)
            )
            
            onSuccess(result) { response =>
              complete(StatusCodes.OK, Map(
                "success" -> response.success,
                "failed" -> response.failed,
                "total" -> batch.notifications.size
              ))
            }
          }
        }
      },
      
      // Mark as read
      path(JavaUUID / "read") { notificationId =>
        put {
          parameter("user_id".as[UUID]) { userId =>
            val result: Future[NotificationActor.OperationResult] = notificationActor.ask(ref =>
              NotificationActor.MarkAsRead(notificationId, userId, ref)
            )
            
            onSuccess(result) { opResult =>
              complete(StatusCodes.OK, Map(
                "success" -> opResult.success,
                "message" -> opResult.message
              ))
            }
          }
        }
      },
      
      // Mark all as read
      path("read-all") {
        put {
          parameter("user_id".as[UUID]) { userId =>
            val result: Future[NotificationActor.OperationResult] = notificationActor.ask(ref =>
              NotificationActor.MarkAllAsRead(userId, ref)
            )
            
            onSuccess(result) { opResult =>
              complete(StatusCodes.OK, Map(
                "success" -> opResult.success,
                "message" -> opResult.message
              ))
            }
          }
        }
      },
      
      // Delete notification
      path(JavaUUID) { notificationId =>
        delete {
          parameter("user_id".as[UUID]) { userId =>
            val result: Future[NotificationActor.OperationResult] = notificationActor.ask(ref =>
              NotificationActor.DeleteNotification(notificationId, userId, ref)
            )
            
            onSuccess(result) { opResult =>
              if (opResult.success) {
                complete(StatusCodes.OK, Map("message" -> "Notification deleted"))
              } else {
                complete(StatusCodes.NotFound, Map("error" -> "Notification not found"))
              }
            }
          }
        }
      }
    )
  }
}
