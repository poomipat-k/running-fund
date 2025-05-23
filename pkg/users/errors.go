package users

type EmailRequiredError struct{}

func (e *EmailRequiredError) Error() string {
	return "email is required"
}

type DuplicatedEmailError struct{}

func (e *DuplicatedEmailError) Error() string {
	return "ที่อยู่อีเมลถูกใช้ไปแล้ว"
}

type InvalidEmailError struct{}

func (e *InvalidEmailError) Error() string {
	return "email is invalid"
}

type EmailTooLongError struct{}

func (e *EmailTooLongError) Error() string {
	return "email is too long"
}

type FirstNameRequiredError struct{}

func (e *FirstNameRequiredError) Error() string {
	return "first name is required"
}

type FirstNameTooLongError struct{}

func (e *FirstNameTooLongError) Error() string {
	return "first name is too long"
}

type LastNameTooLongError struct{}

func (e *LastNameTooLongError) Error() string {
	return "last name is too long"
}

type LastNameRequiredError struct{}

func (e *LastNameRequiredError) Error() string {
	return "last name is required"
}

type PasswordTooShortError struct{}

func (e *PasswordTooShortError) Error() string {
	return "password minimum length are 8 characters"
}

type PasswordRequiredError struct{}

func (e *PasswordRequiredError) Error() string {
	return "password is required"
}

type PasswordTooLongError struct{}

func (e *PasswordTooLongError) Error() string {
	return "password maximum length are 60 characters"
}

type UserNotActivatedError struct{}

func (e *UserNotActivatedError) Error() string {
	return "user is not activated"
}

type InvalidLoginCredentialError struct{}

func (e *InvalidLoginCredentialError) Error() string {
	return "login credential is not valid"
}

type InvalidActivateCodeError struct{}

func (e *InvalidActivateCodeError) Error() string {
	return "invalid activate code"
}

type UserToActivateNotFoundError struct{}

func (e *UserToActivateNotFoundError) Error() string {
	return "email to activate is not found"
}

type UserIsNotActivatedError struct{}

func (e *UserIsNotActivatedError) Error() string {
	return "user is not activated"
}

type PasswordAndConfirmPasswordNotMatchError struct{}

func (e *PasswordAndConfirmPasswordNotMatchError) Error() string {
	return "password and confirmation password is not identical"
}

type ResetPasswordCodeNotValidError struct{}

func (e *ResetPasswordCodeNotValidError) Error() string {
	return "reset password code is not valid"
}

type ResetPasswordCodeNotFound struct{}

func (e *ResetPasswordCodeNotFound) Error() string {
	return "reset password code is not found"
}

type MissingTermsAndConditionError struct{}

func (e *MissingTermsAndConditionError) Error() string {
	return "กรุณายอมรับข้อตกลงและเงื่อนไขการใข้งาน"
}

type MissingPrivacyError struct{}

func (e *MissingPrivacyError) Error() string {
	return "กรุณายอมรับนโยบายคุ้มครองความเป็นส่วนตัว"
}
