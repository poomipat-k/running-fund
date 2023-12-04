package users

type EmailRequiredError struct{}

func (e *EmailRequiredError) Error() string {
	return "email is required"
}

type DuplicatedEmailError struct{}

func (e *DuplicatedEmailError) Error() string {
	return "email is already exist"
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
