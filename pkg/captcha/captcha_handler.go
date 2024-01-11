package captcha

import (
	"log/slog"
	"net/http"

	"github.com/poomipat-k/running-fund/pkg/utils"
)

type CaptchaStore interface {
	GenerateCaptcha() (Captcha, error)
	Get(captchaId string) (int, bool)
}

type CaptchaHandler struct {
	store CaptchaStore
}

func NewCaptchaHandler(s CaptchaStore) *CaptchaHandler {
	return &CaptchaHandler{
		store: s,
	}
}

func (h *CaptchaHandler) GenerateCaptcha(w http.ResponseWriter, r *http.Request) {
	captcha, err := h.store.GenerateCaptcha()
	if err != nil {
		fail(w, err, "captcha")
		return
	}
	utils.WriteJSON(w, http.StatusOK, Captcha{
		CaptchaId:    captcha.CaptchaId,
		Background64: captcha.Background64,
		Puzzle64:     captcha.Puzzle64,
	})
}

func (h *CaptchaHandler) CheckCaptcha(w http.ResponseWriter, r *http.Request) {
	var payload CheckCaptchaRequest
	err := utils.ReadJSON(w, r, &payload)
	if err != nil {
		fail(w, err, "")
		return
	}
	name, err := validateCheckCaptchaPayload(payload)
	if err != nil {
		fail(w, err, name)
		return
	}

	storeValue, found := h.store.Get(payload.CaptchaId)
	if !found {
		fail(w, &CaptchaNotFoundError{}, "captchaId", http.StatusNotFound)
		return
	}
	if payload.CaptchaValue != storeValue {
		fail(w, &CaptchaValueNotValidError{}, "captchaValue")
		return
	}
}

func validateCheckCaptchaPayload(payload CheckCaptchaRequest) (string, error) {
	if payload.CaptchaId == "" {
		return "captchaId", &CaptchaIdRequiredError{}
	}
	if payload.CaptchaValue == 0 {
		return "captchaValue", &CaptchaValueRequiredError{}
	}
	return "", nil
}

func fail(w http.ResponseWriter, err error, name string, status ...int) {
	slog.Error(err.Error())
	s := http.StatusBadRequest
	if len(status) > 0 {
		s = status[0]
	}
	utils.ErrorJSON(w, err, name, s)
}
