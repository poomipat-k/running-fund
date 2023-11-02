package users

import "time"

type User struct {
	Id        int        `json:"id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Email     string     `json:"email"`
	UserRole  string     `json:"user_role"`
	CreatedAt *time.Time `json:"created_at"`
}
