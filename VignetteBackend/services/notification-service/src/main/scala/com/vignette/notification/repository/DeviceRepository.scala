package com.vignette.notification.repository

import slick.jdbc.PostgresProfile.api._
import java.time.Instant
import java.util.UUID
import scala.concurrent.{ExecutionContext, Future}

import com.vignette.notification.model._

/**
 * Device repository for managing push notification devices
 */
class DeviceRepository(db: Database)(implicit ec: ExecutionContext) {
  
  /**
   * Register device
   */
  def register(device: Device): Future[Device] = {
    val query = sqlu"""
      INSERT INTO devices (
        id, user_id, device_token, platform, device_name,
        device_model, os_version, app_version, is_active,
        last_used_at, created_at, updated_at
      ) VALUES (
        ${device.id}, ${device.userId}, ${device.deviceToken},
        ${device.platform.toString}, ${device.deviceName},
        ${device.deviceModel}, ${device.osVersion}, ${device.appVersion},
        ${device.isActive}, ${device.lastUsedAt}, ${device.createdAt}, ${device.updatedAt}
      )
      ON CONFLICT (device_token) DO UPDATE SET
        user_id = EXCLUDED.user_id,
        is_active = TRUE,
        last_used_at = EXCLUDED.last_used_at,
        updated_at = EXCLUDED.updated_at
    """
    
    db.run(query).map(_ => device)
  }
  
  /**
   * Get user devices
   */
  def getUserDevices(userId: UUID): Future[Seq[Device]] = {
    // Simplified - would use proper Slick query
    Future.successful(Seq.empty)
  }
  
  /**
   * Get active devices for user
   */
  def getActiveDevices(userId: UUID): Future[Seq[Device]] = {
    // Simplified - would filter by is_active = true
    Future.successful(Seq.empty)
  }
  
  /**
   * Get device by ID
   */
  def getById(deviceId: UUID): Future[Option[Device]] = {
    Future.successful(None)
  }
  
  /**
   * Deactivate device
   */
  def deactivate(deviceId: UUID, userId: UUID): Future[Boolean] = {
    val query = sqlu"""
      UPDATE devices
      SET is_active = FALSE, updated_at = ${Instant.now()}
      WHERE id = $deviceId AND user_id = $userId
    """
    
    db.run(query).map(_ => true)
  }
  
  /**
   * Delete device
   */
  def delete(deviceId: UUID, userId: UUID): Future[Boolean] = {
    val query = sqlu"""
      DELETE FROM devices
      WHERE id = $deviceId AND user_id = $userId
    """
    
    db.run(query).map(_ => true)
  }
  
  /**
   * Update device last used time
   */
  def updateLastUsed(deviceId: UUID): Future[Boolean] = {
    val query = sqlu"""
      UPDATE devices
      SET last_used_at = ${Instant.now()}, updated_at = ${Instant.now()}
      WHERE id = $deviceId
    """
    
    db.run(query).map(_ => true)
  }
  
  /**
   * Clean up inactive devices (not used in 90 days)
   */
  def cleanupInactive(days: Int): Future[Int] = {
    val cutoff = Instant.now().minusSeconds(days * 24 * 60 * 60)
    
    val query = sqlu"""
      DELETE FROM devices
      WHERE last_used_at < $cutoff
    """
    
    db.run(query).map(_ => 0)
  }
}
