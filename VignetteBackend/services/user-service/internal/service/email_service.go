package service

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
)

// EmailService handles sending emails
type EmailService struct {
	smtpHost     string
	smtpPort     string
	smtpUsername string
	smtpPassword string
	fromEmail    string
	fromName     string
}

// NewEmailService creates a new email service instance
func NewEmailService() *EmailService {
	return &EmailService{
		smtpHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		smtpPort:     getEnv("SMTP_PORT", "587"),
		smtpUsername: getEnv("SMTP_USERNAME", ""),
		smtpPassword: getEnv("SMTP_PASSWORD", ""),
		fromEmail:    getEnv("FROM_EMAIL", "noreply@entativa.com"),
		fromName:     getEnv("FROM_NAME", "Entativa"),
	}
}

// SendPasswordResetEmail sends a password reset email
func (s *EmailService) SendPasswordResetEmail(toEmail, firstName, resetLink string) error {
	subject := "Reset Your Entativa Password"
	
	// HTML template for the email
	htmlTemplate := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #007CFC 0%, #6F3EFB 50%, #FC30E1 100%); color: white; padding: 30px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: white; padding: 30px; border: 1px solid #e4e6eb; border-top: none; border-radius: 0 0 8px 8px; }
        .button { display: inline-block; background: #007CFC; color: white; padding: 14px 32px; text-decoration: none; border-radius: 8px; margin: 20px 0; font-weight: 600; }
        .footer { text-align: center; color: #65676b; font-size: 12px; margin-top: 20px; }
        .warning { background: #fff3cd; border: 1px solid #ffc107; padding: 15px; border-radius: 6px; margin: 20px 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1 style="margin: 0; font-size: 36px; font-style: italic;">entativa</h1>
        </div>
        <div class="content">
            <h2>Hi {{.FirstName}},</h2>
            <p>We received a request to reset your Entativa password. Click the button below to create a new password:</p>
            <p style="text-align: center;">
                <a href="{{.ResetLink}}" class="button">Reset Password</a>
            </p>
            <div class="warning">
                <strong>⚠️ Important:</strong> This link will expire in 1 hour and can only be used once.
            </div>
            <p>If you didn't request this, you can safely ignore this email. Your password won't change unless you click the link above.</p>
            <p>For security reasons, we never send your password via email.</p>
            <hr style="border: none; border-top: 1px solid #e4e6eb; margin: 30px 0;">
            <p style="color: #65676b; font-size: 14px;">
                If the button doesn't work, copy and paste this link into your browser:<br>
                <a href="{{.ResetLink}}" style="color: #007CFC; word-break: break-all;">{{.ResetLink}}</a>
            </p>
        </div>
        <div class="footer">
            <p>© 2025 Entativa. All rights reserved.</p>
            <p>This is an automated email. Please do not reply.</p>
        </div>
    </div>
</body>
</html>
`
	
	// Parse template
	tmpl, err := template.New("passwordReset").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}
	
	// Execute template
	var body bytes.Buffer
	data := struct {
		FirstName string
		ResetLink string
	}{
		FirstName: firstName,
		ResetLink: resetLink,
	}
	
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}
	
	// Send email
	return s.sendEmail(toEmail, subject, body.String())
}

// SendWelcomeEmail sends a welcome email to new users
func (s *EmailService) SendWelcomeEmail(toEmail, firstName string) error {
	subject := "Welcome to Entativa!"
	
	htmlTemplate := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #007CFC 0%, #6F3EFB 50%, #FC30E1 100%); color: white; padding: 40px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: white; padding: 30px; border: 1px solid #e4e6eb; border-top: none; border-radius: 0 0 8px 8px; }
        .button { display: inline-block; background: #007CFC; color: white; padding: 14px 32px; text-decoration: none; border-radius: 8px; margin: 20px 0; font-weight: 600; }
        .footer { text-align: center; color: #65676b; font-size: 12px; margin-top: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1 style="margin: 0; font-size: 48px; font-style: italic;">entativa</h1>
            <p style="font-size: 18px; margin-top: 10px;">Connect with friends and the world</p>
        </div>
        <div class="content">
            <h2>Welcome, {{.FirstName}}! 🎉</h2>
            <p>Your Entativa account has been successfully created. We're excited to have you join our community!</p>
            <p>Here's what you can do next:</p>
            <ul>
                <li>Complete your profile</li>
                <li>Connect with friends</li>
                <li>Share your first post</li>
                <li>Explore groups and communities</li>
            </ul>
            <p style="text-align: center;">
                <a href="https://app.entativa.com" class="button">Get Started</a>
            </p>
        </div>
        <div class="footer">
            <p>© 2025 Entativa. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`
	
	tmpl, err := template.New("welcome").Parse(htmlTemplate)
	if err != nil {
		return err
	}
	
	var body bytes.Buffer
	data := struct {
		FirstName string
	}{
		FirstName: firstName,
	}
	
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}
	
	return s.sendEmail(toEmail, subject, body.String())
}

// sendEmail sends an email using SMTP
func (s *EmailService) sendEmail(to, subject, htmlBody string) error {
	// In development, just log the email
	if s.smtpUsername == "" || s.smtpPassword == "" {
		log.Printf("📧 [DEV MODE] Email would be sent to: %s\nSubject: %s\n", to, subject)
		return nil
	}
	
	// Build email message
	from := fmt.Sprintf("%s <%s>", s.fromName, s.fromEmail)
	message := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", from, to, subject, htmlBody))
	
	// Set up authentication
	auth := smtp.PlainAuth("", s.smtpUsername, s.smtpPassword, s.smtpHost)
	
	// Send email
	addr := fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort)
	err := smtp.SendMail(addr, auth, s.fromEmail, []string{to}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	
	log.Printf("✅ Email sent to: %s", to)
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
