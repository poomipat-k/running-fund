package users

import "time"

type User struct {
	Id              int       `json:"id,omitempty"`
	Email           string    `json:"email,omitempty"`
	Password        string    `json:"password,omitempty"`
	FirstName       string    `json:"firstName,omitempty"`
	LastName        string    `json:"lastName,omitempty"`
	UserRole        string    `json:"userRole,omitempty"`
	Activated       bool      `json:"activated,omitempty"`
	ActivatedBefore time.Time `json:"activatedBefore,omitempty"`
	ActivateCode    string    `json:"activate_code,omitempty"`
	CreatedAt       time.Time `json:"createdAt,omitempty"`
}

type SignUpRequest struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	TermsAndCondition bool   `json:"termsAndConditions"`
	Privacy           bool   `json:"privacy"`
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ActivateUserRequest struct {
	ActivateCode string `json:"activateCode"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	ResetPasswordCode string `json:"resetPasswordCode"`
	Password          string `json:"password"`
	ConfirmPassword   string `json:"confirmPassword"`
}

type CommonSuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
