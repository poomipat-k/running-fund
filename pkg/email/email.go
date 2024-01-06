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
	html := fmt.Sprintf(`<p>เรียน ผู้สนใจส่งข้อเสนอโครงการวิ่งเพื่อสุขภาพ</p>
	<br>
	<p>ขอบพระคุณในความสนใจส่งข้อเสนอโครงการวิ่งเพื่อสุขภาพ อันจะเป็นประโยชน์ต่อสุขภาพประชาชน</p>
	<p>ระบบลงทะเบียนได้รับคำขอใช้งาน เพื่อส่งข้อเสนอโครงการฯ ของท่านแล้ว โปรดกดปุ่มด้านล่างนี้ เพื่อยืนยันอีเมลของท่านอีกครั้ง หลังจากนั้น ท่านจะสามารถเข้าสู่ขั้นตอนต่อไปในการส่งข้อเสนอโครงการ</p>
	<br>
	<span style="padding-left: 40px;">กรุณากดลิงก์เพื่อยืนยัน <a href="%s">%s</a></span>
	<br><br>
	<p>หมายเหตุ: ท่านจะไม่สามารถใช้งานระบบส่งข้อเสนอโครงการได้ หากไม่ได้กดยืนยันอีเมลภายใน 24 ชั่วโมง (ตัวหนา) ก่อนลิ้งก์จะหมดอายุ หากลิ้งก์นี้หมดอายุ ท่านต้องลงทะเบียนใหม่อีกครั้ง</p>
	<br>
	<p>ขอแสดงความนับถือ</p>
	<p>ผู้ดูแลระบบ</p>
	<p>มูลนิธิสมาพันธ์ชมรมเดิน-วิ่งเพื่อสุขภาพไทย</p>
	`, activateLink, activateLink)
	text := fmt.Sprintf(`เรียน ผู้สนใจส่งข้อเสนอโครงการวิ่งเพื่อสุขภาพ

	ขอบพระคุณในความสนใจส่งข้อเสนอโครงการวิ่งเพื่อสุขภาพ อันจะเป็นประโยชน์ต่อสุขภาพประชาชน

	ระบบลงทะเบียนได้รับคำขอใช้งาน เพื่อส่งข้อเสนอโครงการฯ ของท่านแล้ว โปรดกดปุ่มด้านล่างนี้ เพื่อยืนยันอีเมลของท่านอีกครั้ง หลังจากนั้น ท่านจะสามารถเข้าสู่ขั้นตอนต่อไปในการส่งข้อเสนอโครงการ
	
		กรุณากดลิงก์เพื่อยืนยัน %s
	
	หมายเหตุ: ท่านจะไม่สามารถใช้งานระบบส่งข้อเสนอโครงการได้ หากไม่ได้กดยืนยันอีเมลภายใน 24 ชั่วโมง (ตัวหนา) ก่อนลิ้งก์จะหมดอายุ หากลิ้งก์นี้หมดอายุ ท่านต้องลงทะเบียนใหม่อีกครั้ง
	
	ขอแสดงความนับถือ
	ผู้ดูแลระบบ
	มูลนิธิสมาพันธ์ชมรมเดิน-วิ่งเพื่อสุขภาพไทย`, activateLink)
	mail := email.Email{
		From:    os.Getenv("EMAIL_SENDER"),
		To:      []string{to},
		Subject: "กรุณายืนยันอีเมลเพื่อส่งข้อเสนอโครงการวิ่งเพื่อสุขภาพ",
		Text:    []byte(text),
		HTML:    []byte(html),
	}
	return mail
}

func (es *EmailService) BuildResetPasswordEmail(to, resetPasswordLink string) email.Email {
	html := fmt.Sprintf(`<p>เรียน ท่านสมาชิก</p>
	<p>กรุณาคลิกที่ลิงค์ด้านล่างเพื่อตั้งรหัสผ่านใหม่</p>
	<a href="%s">%s</a>
	`, resetPasswordLink, resetPasswordLink)
	text := fmt.Sprintf(`เรียน ท่านสมาชิก
	กรุณาคลิกที่ลิงค์ด้านล่างเพื่อตั้งรหัสผ่านใหม่
	
	%s`, resetPasswordLink)
	mail := email.Email{
		From:    os.Getenv("EMAIL_SENDER"),
		To:      []string{to},
		Subject: fmt.Sprintf("Account password reset - %s", to),
		Text:    []byte(text),
		HTML:    []byte(html),
	}
	return mail
}
