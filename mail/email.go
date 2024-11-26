package mail

import (
	"github.com/Tsuzat/zipit-go-fiber/config"
	"github.com/Tsuzat/zipit-go-fiber/models"
	"github.com/gofiber/fiber/v2/log"
	"gopkg.in/gomail.v2"
)

func SendEmail(to string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "zipit@tsuzat.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(config.EMAIL_HOST, config.EMAIL_PORT, config.EMAIL_USERNAME, config.EMAIL_PASSWORD)

	if err := d.DialAndSend(m); err != nil {
		log.Error("Error sending email: ", err)
		return err
	}
	log.Info("Email Send Successfully to ", to)
	return nil
}

func SendEmailVerification(user *models.User) error {
	to := user.Email
	link := config.BACKEND_URL + "/api/v1/auth/verify?email=" + user.Email + "&verification_token=" + user.VerificationToken
	body := emailVerificationTemplate(link)
	return SendEmail(to, "Verify your email", body)
}
