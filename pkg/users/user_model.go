package users

import "time"

type User struct {
	Id        int        `json:"id"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Email     string     `json:"email"`
	UserRole  string     `json:"userRole"`
	CreatedAt *time.Time `json:"createdAt"`
}
