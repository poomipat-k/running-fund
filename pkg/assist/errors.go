package assist

type EmailRequiredError struct{}

func (e *EmailRequiredError) Error() string {
	return "email is required"
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
	return "firstName is required"
}

type LastNameRequiredError struct{}

func (e *LastNameRequiredError) Error() string {
	return "lastName is required"
}

type MessageRequiredError struct{}

func (e *MessageRequiredError) Error() string {
	return "message is required"
}
