package assist

type EmailRequiredError struct{}

func (e *EmailRequiredError) Error() string {
	return "email is required"
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
