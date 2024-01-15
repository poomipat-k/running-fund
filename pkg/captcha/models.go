package captcha

type Captcha struct {
	CaptchaId    string  `json:"captchaId,omitempty"`
	Background64 string  `json:"background,omitempty"`
	Puzzle64     string  `json:"puzzle,omitempty"`
	CaptchaValue float64 `json:"captchaValue,omitempty"`
	YPosition    float64 `json:"yPosition,omitempty"`
}

type CheckCaptchaRequest struct {
	CaptchaId    string  `json:"captchaId"`
	CaptchaValue float64 `json:"captchaValue"`
}

type Puzzle struct {
	BackgroundPath string  `json:"backgroundPath"`
	PuzzlePath     string  `json:"puzzlePath"`
	Value          float64 `json:"value"`
	YPosition      float64 `json:"yPosition"`
}

type CommonSuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
