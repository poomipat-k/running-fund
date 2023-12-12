package appEmail

import (
	"fmt"
	"log/slog"
	"net/smtp"
	"os"

	"github.com/jordan-wright/email"
)

type EmailService struct{}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (es *EmailService) SendEmail(em email.Email) error {
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

func (es *EmailService) BuildSignUpConfirmationEmail(to, activateLink string) email.Email {
	html := fmt.Sprintf(`<p>เรียนผู้ขอทุน,</p>
	<p>เพื่อยืนยันการสมัครสมาชิก กรุณายืนยันตัวตนด้วยการคลิกที่ลิ้งด้านล่าง</p>
	<a href="%s">%s</a>
	<br><br>
	<p>บัญชีของท่านจะยังไม่สามารถใช้ได้จนกว่าจะทำการยืนยันด้วยลิ้งด้านบน ลิ้งด้านบนจะหมดอายุภายใน 24 ชั่วโมง</p>
	<p>ถ้าท่านไม่ได้ทำการสมัครสมาชิกเว็บไซต์ running-fund กรุณาอย่าคลิกลิ้งด้านบน</p>
	`, activateLink, activateLink)
	text := fmt.Sprintf(`เรียนผู้ขอทุน,
	เพื่อยืนยันการสมัครสมาชิก กรุณายืนยันตัวตนด้วยการคลิกที่ลิ้งด้านล่าง
	
	%s
	
	บัญชีของท่านจะยังไม่สามารถใช้ได้จนกว่าจะทำการยืนยันด้วยลิ้งด้านบน ลิ้งด้านบนจะหมดอายุภายใน 24 ชั่วโมง
	ถ้าท่านไม่ได้ทำการสมัครสมาชิกเว็บไซต์ running-fund กรุณาอย่าคลิกลิ้งด้านบน`, activateLink)
	mail := email.Email{
		From:    os.Getenv("EMAIL_SENDER"),
		To:      []string{to},
		Subject: fmt.Sprintf("Registration confirmation - %s", to),
		Text:    []byte(text),
		HTML:    []byte(html),
	}
	return mail
}
