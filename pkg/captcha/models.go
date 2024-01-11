package captcha

type Captcha struct {
	CaptchaId    string `json:"captchaId,omitempty"`
	Background64 string `json:"background,omitempty"`
	Puzzle64     string `json:"puzzle,omitempty"`
	CaptchaValue int    `json:"captchaValue,omitempty"`
}

type CheckCaptchaRequest struct {
	CaptchaId    string `json:"captchaId"`
	CaptchaValue int    `json:"captchaValue"`
}

type Puzzle struct {
	BackgroundPath string `json:"backgroundPath"`
	PuzzlePath     string `json:"puzzlePath"`
	Value          int    `json:"value"`
}
