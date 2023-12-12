package appEmail

import (
	"fmt"
	"log/slog"
	"net/smtp"
	"os"

	"github.com/jordan-wright/email"
)

func SendEmail(em email.Email) error {
	e := email.NewEmail()
	e.From = em.From
	e.To = em.To
	e.Subject = em.Subject
	e.Text = em.Text
	e.HTML = em.HTML
	err := e.Send(
		fmt.Sprintf("%s:%s", os.Getenv("EMAIL_HOST"), os.Getenv("EMAIL_PORT")),
		smtp.PlainAuth(
			"",
			os.Getenv("EMAIL_AUTH_USERNAME"),
			os.Getenv("EMAIL_AUTH_PASSWORD"),
			os.Getenv("EMAIL_HOST"),
		))
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil

}
