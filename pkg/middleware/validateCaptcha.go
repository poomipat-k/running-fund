package mw

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/poomipat-k/running-fund/pkg/captcha"
	"github.com/poomipat-k/running-fund/pkg/utils"
)

func ValidateCaptcha(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var captcha captcha.CheckCaptchaRequest

		// middleware do read the body
		bodyBytes, _ := io.ReadAll(r.Body)
		r.Body.Close() //  must close
		// put back the bodyBytes to r.Body
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		err := json.Unmarshal(bodyBytes, &captcha)
		if err != nil {
			utils.ErrorJSON(w, err, "unmarshal payload")
			return
		}

		next(w, r)
	})
}
