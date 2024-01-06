package captcha

import (
	"net/http"

	"github.com/poomipat-k/running-fund/pkg/utils"
)

type CaptchaStore interface {
	GenerateCaptchaId() (string, int)
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
	captchaId, _ := h.store.GenerateCaptchaId()
	utils.WriteJSON(w, http.StatusOK, Captcha{
		CaptchaId: captchaId,
	})
}
