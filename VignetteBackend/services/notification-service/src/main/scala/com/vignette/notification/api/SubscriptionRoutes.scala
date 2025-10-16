package com.vignette.notification.api

import akka.actor.typed.ActorSystem
import akka.http.scaladsl.model.StatusCodes
import akka.http.scaladsl.server.Directives._
import akka.http.scaladsl.server.Route
import scala.concurrent.ExecutionContext
import java.util.UUID

import com.vignette.notification.model._
import com.vignette.notification.repository.DeviceRepository
import DeviceJsonProtocol._

/**
 * Subscription (device) API routes
 */
class SubscriptionRoutes(
  deviceRepository: DeviceRepository
)(implicit system: ActorSystem[_], ec: ExecutionContext) {
  
  val routes: Route = pathPrefix("api" / "v1" / "devices") {
    concat(
      // Register device
      pathEnd {
        post {
          entity(as[DeviceRegistration]) { registration =>
            val device = Device(
              id = UUID.randomUUID(),
              userId = registration.userId,
              deviceToken = registration.deviceToken,
              platform = registration.platform,
              deviceName = registration.deviceName,
              deviceModel = registration.deviceModel,
              osVersion = registration.osVersion,
              appVersion = registration.appVersion
            )
            
            val result = deviceRepository.register(device)
            
            onSuccess(result) { registeredDevice =>
              complete(StatusCodes.Created, Map(
                "success" -> true,
                "device_id" -> registeredDevice.id.toString,
                "message" -> "Device registered successfully"
              ))
            }
          }
        }
      },
      
      // Get user devices
      pathEnd {
        get {
          parameter("user_id".as[UUID]) { userId =>
            val devices = deviceRepository.getUserDevices(userId)
            
            onSuccess(devices) { deviceList =>
              complete(StatusCodes.OK, Map(
                "devices" -> deviceList,
                "count" -> deviceList.size
              ))
            }
          }
        }
      },
      
      // Delete device
      path(JavaUUID) { deviceId =>
        delete {
          parameter("user_id".as[UUID]) { userId =>
            val result = deviceRepository.delete(deviceId, userId)
            
            onSuccess(result) { success =>
              if (success) {
                complete(StatusCodes.OK, Map("message" -> "Device removed"))
              } else {
                complete(StatusCodes.NotFound, Map("error" -> "Device not found"))
              }
            }
          }
        }
      },
      
      // Deactivate device
      path(JavaUUID / "deactivate") { deviceId =>
        put {
          parameter("user_id".as[UUID]) { userId =>
            val result = deviceRepository.deactivate(deviceId, userId)
            
            onSuccess(result) { success =>
              complete(StatusCodes.OK, Map("message" -> "Device deactivated"))
            }
          }
        }
      }
    )
  }
}
