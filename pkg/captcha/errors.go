package captcha

type CaptchaIdRequiredError struct{}

func (e *CaptchaIdRequiredError) Error() string {
	return "captchaId is required"
}

type CaptchaValueRequiredError struct{}

func (e *CaptchaValueRequiredError) Error() string {
	return "captchaValue is required"
}

type CaptchaNotFoundError struct{}

func (e *CaptchaNotFoundError) Error() string {
	return "captcha is not found"
}

type CaptchaValueNotValidError struct{}

func (e *CaptchaValueNotValidError) Error() string {
	return "captcha value is not valid"
}
