package users

type EmailRequiredError struct{}

func (e *EmailRequiredError) Error() string {
	return "email is required"
}
