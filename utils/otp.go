package utils

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
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
	apiKey := os.Getenv("BREVO_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("BREVO_API_KEY not set")
	}

	client := resty.New()

	payload := map[string]interface{}{
		"sender": map[string]string{
			"name":  "Admin-noreply",
			"email": os.Getenv("ADMIN_EMAIL"),
		},
		"to": []map[string]string{
			{"email": to},
		},
		"subject": subject,
		"htmlContent": fmt.Sprintf(`
			<html>
				<body style="font-family: Arial, sans-serif; color: #333; line-height: 1.6;">
					<p>Hello,</p>
					<p>We received a request to log in using a One-Time Password (OTP).</p>

					<div style="padding: 20px; margin: 20px 0; background-color: #f4f4f4; border-left: 5px solid #4CAF50;">
						<p style="margin: 0; font-size: 16px;">🔐 <strong>Your OTP Code:</strong></p>
						<p style="margin: 5px 0; font-size: 28px; font-weight: bold; color: #4CAF50;">%s</p>
					</div>

					<p>This code is valid for <strong>2 minutes</strong>. Please use it immediately to complete your login.</p>
					<p><strong>Do not share this code</strong> with anyone. If you did not request this code, you can safely ignore this email.</p>

					<p>Best regards,<br><strong>Admin Team</strong></p>
				</body>
			</html>
		`, body),
	}
	brevoApiSendEmail := os.Getenv("BREVO_API_SEND_EMAIL")
	resp, err := client.R().
		SetHeader("accept", "application/json").
		SetHeader("api-key", apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post(brevoApiSendEmail)

	if err != nil {
		log.Printf("❌ Email send failed: %v", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	if resp.StatusCode() >= 300 {
		log.Printf("❌ Email failed with status: %d, body: %s", resp.StatusCode(), resp.String())
		return fmt.Errorf("email failed with status: %d", resp.StatusCode())
	}

	log.Printf("✅ Email sent successfully to %s", to)
	return nil
}

// SEND OLD OTP IN LOCAL
// func SendEmail(to, subject, body string) error {
// 	from := os.Getenv("ADMIN_EMAIL")
// 	username := os.Getenv("SMTP_USERNAME")
// 	password := os.Getenv("PASSWORD_EMAIL")
// 	smtpHost := os.Getenv("SMTP_HOST")
// 	smtpPort := os.Getenv("SMTP_PORT")

// 	if from == "" || username == "" || password == "" || smtpHost == "" || smtpPort == "" {
// 		return fmt.Errorf("SMTP config not set in environment variables")
// 	}

// 	// Header email
// 	header := "MIME-Version: 1.0\r\n" +
// 		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
// 		fmt.Sprintf("From: %s\r\n", from) +
// 		fmt.Sprintf("To: %s\r\n", to) +
// 		fmt.Sprintf("Subject: %s\r\n\r\n", subject)

// 	// HTML body
// 	htmlBody := fmt.Sprintf(`
//         <html>
//             <body style="font-family: Arial, sans-serif;">
//                 <p>Hello,</p>
//                 <p>Your OTP code is:</p>
//                 <h2 style="color: #4CAF50;">%s</h2>
//                 <p>This code is valid for 5 minutes.</p>
//             </body>
//         </html>
//     `, body)

// 	message := []byte(header + htmlBody)

// 	auth := smtp.PlainAuth("", username, password, "smtp-relay.brevo.com")

// 	// Kirim email
// 	log.Printf("📧 Sending email to %s via %s:%s", to, smtpHost, smtpPort)
// 	log.Printf("📧 Using username: %s", username)

// 	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
// 	if err != nil {
// 		log.Printf("❌ Email send failed: %v", err)
// 		return fmt.Errorf("failed to send email: %w", err)
// 	}

// 	return nil
// }
