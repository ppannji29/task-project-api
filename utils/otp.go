package utils

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"os"
	"time"
)

func GenerateOtp(length int) string {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	otp := make([]byte, length)
	for i := range otp {
		otp[i] = digits[rand.Intn(len(digits))]
	}
	return string(otp)
}

func SendEmail(to, subject, body string) error {
	from := os.Getenv("ADMIN_EMAIL")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("PASSWORD_EMAIL")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	if from == "" || username == "" || password == "" || smtpHost == "" || smtpPort == "" {
		return fmt.Errorf("SMTP config not set in environment variables")
	}

	// Header email
	header := "MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		fmt.Sprintf("From: %s\r\n", from) +
		fmt.Sprintf("To: %s\r\n", to) +
		fmt.Sprintf("Subject: %s\r\n\r\n", subject)

	// HTML body
	htmlBody := fmt.Sprintf(`
        <html>
            <body style="font-family: Arial, sans-serif;">
                <p>Hello,</p>
                <p>Your OTP code is:</p>
                <h2 style="color: #4CAF50;">%s</h2>
                <p>This code is valid for 5 minutes.</p>
            </body>
        </html>
    `, body)

	message := []byte(header + htmlBody)

	auth := smtp.PlainAuth("", username, password, "smtp-relay.brevo.com")

	// Kirim email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// func SendEmail(to, subject, body string) error {
// 	from := os.Getenv("ADMIN_EMAIL")
// 	username := os.Getenv("SMTP_USERNAME")
// 	password := os.Getenv("PASSWORD_EMAIL")
// 	smtpHost := os.Getenv("SMTP_HOST")
// 	smtpPort := os.Getenv("SMTP_PORT")

// 	if username == "" || password == "" || smtpHost == "" || smtpPort == "" {
// 		return fmt.Errorf("SMTP config not set in environment variables")
// 	}

// 	// Header email (MIME + HTML)
// 	header := "MIME-Version: 1.0\r\n" +
// 		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
// 		fmt.Sprintf("From: %s\r\n", from) +
// 		fmt.Sprintf("To: %s\r\n", to) +
// 		fmt.Sprintf("Subject: %s\r\n\r\n", subject)

// 	// HTML body dengan template profesional
// 	htmlBody := fmt.Sprintf(`
// 		<html>
// 			<body style="font-family: Arial, sans-serif; color: #333;">
// 				<p>Hello,</p>
// 				<p>We received a request to log in to your account using a One-Time Password (OTP).</p>

// 				<div style="padding: 20px; margin: 20px 0; background-color: #f9f9f9; border-left: 5px solid #4CAF50;">
// 					<p style="margin: 0; font-size: 18px;">🔐 <strong>Your OTP Code:</strong></p>
// 					<p style="margin: 5px 0; font-size: 28px; font-weight: bold; color: #4CAF50;">%s</p>
// 				</div>

// 				<p>This code is valid for 5 minutes. <strong>Do not share this code</strong> with anyone.</p>
// 				<p>If you did not request this code, you can safely ignore this email.</p>

// 				<p>Best regards,<br><strong>Developer Testing Team</strong></p>
// 			</body>
// 		</html>
// 		`, body)

// 	// Gabungkan header dan body
// 	message := []byte(header + htmlBody)

// 	// Auth menggunakan kredensial SMTP
// 	auth := smtp.PlainAuth("", username, password, smtpHost)

// 	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
// 	if err != nil {
// 		return fmt.Errorf("failed to send email: %w", err)
// 	}

// 	return nil
// }
