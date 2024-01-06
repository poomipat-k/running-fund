package captcha

type Captcha struct {
	CaptchaId    string `json:"captchaId,omitempty"`
	CaptchaValue int    `json:"captchaValue,omitempty"`
}
