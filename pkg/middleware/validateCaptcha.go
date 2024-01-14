package mw

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/poomipat-k/running-fund/pkg/captcha"
	"github.com/poomipat-k/running-fund/pkg/utils"
)

func ValidateCaptcha(next http.HandlerFunc, captchaStore captcha.CaptchaStore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var captchaPayload captcha.CheckCaptchaRequest

		// middleware do read the body
		bodyBytes, _ := io.ReadAll(r.Body)
		r.Body.Close() //  must close
		// put back the bodyBytes to r.Body
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		err := json.Unmarshal(bodyBytes, &captchaPayload)
		if err != nil {
			utils.ErrorJSON(w, err, "unmarshal payload")
			return
		}

		errName, err := captcha.CheckCaptcha(captchaPayload.CaptchaId, captchaPayload.CaptchaValue, captchaStore)
		if err != nil {
			// Delete captcha when captcha answer is not valid
			captchaStore.Delete(captchaPayload.CaptchaId)
			utils.ErrorJSON(w, err, errName)
			return
		}
		next(w, r)
	})
}
