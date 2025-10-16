package com.vignette.notification.actor

import akka.actor.typed.{ActorRef, Behavior}
import akka.actor.typed.scaladsl.{ActorContext, Behaviors}
import java.util.UUID

import com.vignette.notification.model._
import com.vignette.notification.service.EmailService

/**
 * Email notification actor
 * Sends email notifications with templates
 */
object EmailActor {
  
  // Commands
  sealed trait Command
  
  case class SendEmail(
    userId: UUID,
    notification: Notification
  ) extends Command
  
  case class SendTemplatedEmail(
    userId: UUID,
    templateName: String,
    data: Map[String, String],
    subject: String
  ) extends Command
  
  case class SendBatchEmail(
    userIds: Seq[UUID],
    subject: String,
    body: String
  ) extends Command
  
  private case class EmailResult(success: Boolean, userId: UUID, error: Option[String]) extends Command
  
  /**
   * Create email actor
   */
  def apply(emailService: EmailService): Behavior[Command] = {
    Behaviors.setup { context =>
      context.log.info("EmailActor started")
      active(emailService, context)
    }
  }
  
  private def active(
    emailService: EmailService,
    context: ActorContext[Command]
  ): Behavior[Command] = {
    Behaviors.receiveMessage {
      
      case SendEmail(userId, notification) =>
        context.log.debug(s"Sending email to user $userId")
        
        // Build email content
        val subject = notification.title
        val body = buildEmailBody(notification)
        
        // Get user email (would come from user service)
        val userEmail = s"user-$userId@example.com" // Placeholder
        
        // Send email asynchronously
        emailService.sendEmail(userEmail, subject, body) match {
          case Right(_) =>
            context.log.info(s"Email sent successfully to user $userId")
            context.self ! EmailResult(success = true, userId, None)
          case Left(error) =>
            context.log.error(s"Email failed: $error")
            context.self ! EmailResult(success = false, userId, Some(error))
        }
        
        Behaviors.same
      
      case SendTemplatedEmail(userId, templateName, data, subject) =>
        context.log.debug(s"Sending templated email to user $userId")
        
        val userEmail = s"user-$userId@example.com" // Placeholder
        
        emailService.sendTemplatedEmail(userEmail, templateName, data, subject) match {
          case Right(_) =>
            context.log.info(s"Templated email sent to user $userId")
            context.self ! EmailResult(success = true, userId, None)
          case Left(error) =>
            context.log.error(s"Templated email failed: $error")
            context.self ! EmailResult(success = false, userId, Some(error))
        }
        
        Behaviors.same
      
      case SendBatchEmail(userIds, subject, body) =>
        context.log.debug(s"Batch sending email to ${userIds.size} users")
        
        // Send to each user
        userIds.foreach { userId =>
          val userEmail = s"user-$userId@example.com"
          emailService.sendEmail(userEmail, subject, body)
        }
        
        Behaviors.same
      
      case EmailResult(success, userId, error) =>
        if (success) {
          context.log.debug(s"Email delivered successfully to user $userId")
        } else {
          context.log.warning(s"Email failed for user $userId: ${error.getOrElse("Unknown error")}")
        }
        
        Behaviors.same
    }
  }
  
  /**
   * Build HTML email body from notification
   */
  private def buildEmailBody(notification: Notification): String = {
    s"""
      <!DOCTYPE html>
      <html>
      <head>
        <style>
          body { font-family: Arial, sans-serif; }
          .container { max-width: 600px; margin: 0 auto; padding: 20px; }
          .header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 20px; border-radius: 8px 8px 0 0; }
          .content { background: #f9f9f9; padding: 20px; }
          .footer { background: #333; color: #999; padding: 15px; text-align: center; border-radius: 0 0 8px 8px; }
          .button { background: #667eea; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; display: inline-block; }
        </style>
      </head>
      <body>
        <div class="container">
          <div class="header">
            <h2>${notification.title}</h2>
          </div>
          <div class="content">
            <p>${notification.message}</p>
            ${notification.actorUsername.map(u => s"<p><strong>From:</strong> @$u</p>").getOrElse("")}
            ${notification.deepLink.map(link => s"""<p><a href="$link" class="button">View Now</a></p>""").getOrElse("")}
          </div>
          <div class="footer">
            <p>Vignette - Share Your Moments</p>
            <p><small>Manage notification settings in your app</small></p>
          </div>
        </div>
      </body>
      </html>
    """
  }
}
