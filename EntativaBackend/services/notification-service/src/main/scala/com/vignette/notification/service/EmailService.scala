package com.socialink.notification.service

import java.util.Properties
import javax.mail._
import javax.mail.internet._
import scala.util.{Try, Success, Failure}

/**
 * Email notification service
 * Sends email notifications with HTML templates
 */
class EmailService(
  smtpHost: String,
  smtpPort: Int,
  username: String,
  password: String,
  fromAddress: String
) {
  
  private val props = new Properties()
  props.put("mail.smtp.host", smtpHost)
  props.put("mail.smtp.port", smtpPort.toString)
  props.put("mail.smtp.auth", "true")
  props.put("mail.smtp.starttls.enable", "true")
  
  private val session = Session.getInstance(props, new Authenticator {
    override def getPasswordAuthentication: PasswordAuthentication = {
      new PasswordAuthentication(username, password)
    }
  })
  
  /**
   * Send email
   */
  def sendEmail(to: String, subject: String, body: String, isHTML: Boolean = true): Either[String, Unit] = {
    Try {
      val message = new MimeMessage(session)
      message.setFrom(new InternetAddress(fromAddress))
      message.setRecipients(Message.RecipientType.TO, to)
      message.setSubject(subject)
      
      if (isHTML) {
        message.setContent(body, "text/html; charset=utf-8")
      } else {
        message.setText(body)
      }
      
      Transport.send(message)
      
    } match {
      case Success(_) => Right(())
      case Failure(ex) => Left(s"Email error: ${ex.getMessage}")
    }
  }
  
  /**
   * Send templated email
   */
  def sendTemplatedEmail(
    to: String,
    templateName: String,
    data: Map[String, String],
    subject: String
  ): Either[String, Unit] = {
    val body = renderTemplate(templateName, data)
    sendEmail(to, subject, body, isHTML = true)
  }
  
  /**
   * Send batch emails
   */
  def sendBatch(recipients: Seq[String], subject: String, body: String): Either[String, Int] = {
    var successCount = 0
    var errors = List.empty[String]
    
    recipients.foreach { recipient =>
      sendEmail(recipient, subject, body) match {
        case Right(_) => successCount += 1
        case Left(error) => errors = error :: errors
      }
    }
    
    if (errors.isEmpty) {
      Right(successCount)
    } else {
      Left(s"$successCount/${recipients.size} sent. Errors: ${errors.take(5).mkString(", ")}")
    }
  }
  
  /**
   * Render email template
   */
  private def renderTemplate(templateName: String, data: Map[String, String]): String = {
    templateName match {
      case "notification" =>
        s"""
          <!DOCTYPE html>
          <html>
          <head>
            <style>
              body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; margin: 0; padding: 0; background: #f5f5f5; }
              .container { max-width: 600px; margin: 40px auto; background: white; border-radius: 12px; overflow: hidden; box-shadow: 0 4px 6px rgba(0,0,0,0.1); }
              .header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 30px; text-align: center; }
              .content { padding: 30px; }
              .actor { display: flex; align-items: center; margin-bottom: 20px; }
              .avatar { width: 48px; height: 48px; border-radius: 50%; margin-right: 12px; background: #667eea; }
              .message { font-size: 18px; color: #333; line-height: 1.6; }
              .button { display: inline-block; background: #667eea; color: white; padding: 14px 28px; text-decoration: none; border-radius: 8px; margin-top: 20px; font-weight: 600; }
              .footer { background: #f9f9f9; padding: 20px; text-align: center; color: #666; font-size: 14px; }
            </style>
          </head>
          <body>
            <div class="container">
              <div class="header">
                <h1>📱 ${data.getOrElse("title", "New Notification")}</h1>
              </div>
              <div class="content">
                ${data.get("actorUsername").map(u => s"""
                <div class="actor">
                  <div class="avatar"></div>
                  <div><strong>@$u</strong></div>
                </div>
                """).getOrElse("")}
                <div class="message">
                  ${data.getOrElse("message", "You have a new notification")}
                </div>
                ${data.get("deepLink").map(link => s"""
                <a href="$link" class="button">View Now</a>
                """).getOrElse("")}
              </div>
              <div class="footer">
                <p><strong>Socialink</strong> - Share Your Moments</p>
                <p><small>Manage notification settings in your app</small></p>
              </div>
            </div>
          </body>
          </html>
        """
      
      case "weekly_digest" =>
        s"""
          <!DOCTYPE html>
          <html>
          <head>
            <style>
              body { font-family: Arial, sans-serif; }
              .container { max-width: 600px; margin: 0 auto; }
              .stat { background: #f0f0f0; padding: 15px; margin: 10px 0; border-radius: 8px; }
            </style>
          </head>
          <body>
            <div class="container">
              <h2>Your Weekly Highlights</h2>
              <div class="stat">
                <h3>${data.getOrElse("likes", "0")} Likes</h3>
              </div>
              <div class="stat">
                <h3>${data.getOrElse("comments", "0")} Comments</h3>
              </div>
              <div class="stat">
                <h3>${data.getOrElse("followers", "0")} New Followers</h3>
              </div>
            </div>
          </body>
          </html>
        """
      
      case _ =>
        s"<html><body><p>${data.getOrElse("message", "")}</p></body></html>"
    }
  }
}
