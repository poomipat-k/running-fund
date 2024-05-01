package assist

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/mail"
	"os"

	"github.com/jordan-wright/email"
	"github.com/poomipat-k/running-fund/pkg/utils"
)

const ADMIN_EMAIL = "ADMIN_EMAIL"

type EmailService interface {
	SendEmail(email email.Email) error
}

type AssistHandler struct {
	emailService EmailService
}

func NewAssistHandler(es EmailService) *AssistHandler {
	return &AssistHandler{
		emailService: es,
	}
}

func (h *AssistHandler) ContactUs(w http.ResponseWriter, r *http.Request) {
	var payload ContactUsRequest
	err := utils.ReadJSON(w, r, &payload)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "")
		return
	}
	err = validateContactUs(payload)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "")
		return
	}
	// Send email from system to admin
	adminEmail := os.Getenv(ADMIN_EMAIL)
	mail := buildContactUsEmail(adminEmail, payload.Email, payload.FirstName, payload.LastName, payload.Message)
	err = h.emailService.SendEmail(mail)
	if err != nil {
		slog.Error("ContactUs: failed to send contact-us email to website admin", "error", err.Error())
		utils.ErrorJSON(w, err, "", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, CommonSuccessResponse{
		Success: true,
		Message: "message successfully sent to admin",
	})
}

func validateContactUs(payload ContactUsRequest) error {
	err := validateEmail(payload.Email)
	if err != nil {
		return err
	}
	if payload.FirstName == "" {
		return &FirstNameRequiredError{}
	}
	if payload.LastName == "" {
		return &LastNameRequiredError{}
	}
	if payload.Message == "" {
		return &MessageRequiredError{}
	}
	return nil
}

func validateEmail(email string) error {
	if email == "" {
		return &EmailRequiredError{}
	}
	if len(email) > 255 {
		return &EmailTooLongError{}
	}
	if !isValidEmail(email) {
		return &InvalidEmailError{}
	}
	return nil
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func buildContactUsEmail(to, userEmail, firstName, lastName, message string) email.Email {
	html := fmt.Sprintf(`<p>เรียน ผู้ดูแลเว็บไซต์</p>
	<br>

	<p>ที่อยู่อีเมล: %s</p>
	<p>ชื่อ-นามสกุล: %s %s</p>
	<p>ข้อความ: %s</p>

	<br>
	<p>ขอแสดงความนับถือ</p>
	<p>ระบบขอความช่วยเหลือ</p>
	<p>มูลนิธิสมาพันธ์ชมรมเดิน-วิ่งเพื่อสุขภาพไทย</p>
	`, userEmail, firstName, lastName, message)
	text := fmt.Sprintf(`เรียน ผู้ดูแลเว็บไซต์

	ที่อยู่อีเมล: %s
	ชื่อ-นามสกุล: %s %s
	ข้อความ: %s
	
	ขอแสดงความนับถือ
	ระบบขอความช่วยเหลือ
	มูลนิธิสมาพันธ์ชมรมเดิน-วิ่งเพื่อสุขภาพไทย`, userEmail, firstName, lastName, message)
	mail := email.Email{
		From:    os.Getenv("EMAIL_SENDER"),
		To:      []string{to},
		Subject: "กรุณายืนยันอีเมลเพื่อส่งข้อเสนอโครงการวิ่งเพื่อสุขภาพ",
		Text:    []byte(text),
		HTML:    []byte(html),
	}
	return mail
}
