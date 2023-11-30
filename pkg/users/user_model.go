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
	CreatedAt       time.Time `json:"createdAt,omitempty"`
}

type SignUpRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
