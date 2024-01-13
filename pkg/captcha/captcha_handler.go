package captcha

import (
	"log/slog"
	"net/http"

	"github.com/poomipat-k/running-fund/pkg/utils"
)

const captchaErrorLimit = 5.0 // in pixels

type CaptchaStore interface {
	GenerateCaptcha() (Captcha, error)
	Get(captchaId string) (float64, bool)
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
		YPosition:    captcha.YPosition,
	})
}

func CheckCaptcha(captchaId string, captchaValue float64, store CaptchaStore) (string, error) {
	name, err := validateCheckCaptchaPayload(captchaId, captchaValue)
	if err != nil {
		return name, err
	}
	storeValue, found := store.Get(captchaId)
	if !found {
		return "captchaId", &CaptchaNotFoundError{}
	}
	if captchaValue-storeValue < -captchaErrorLimit || captchaValue-storeValue > captchaErrorLimit {
		return "captchaValue", &CaptchaValueNotValidError{}
	}
	return "", nil
}

func validateCheckCaptchaPayload(captchaId string, captchaValue float64) (string, error) {
	if captchaId == "" {
		return "captchaId", &CaptchaIdRequiredError{}
	}
	if captchaValue == 0 {
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
