package users_test

import (
	"github.com/jordan-wright/email"
	"github.com/poomipat-k/running-fund/pkg/users"
)

type MockUserStore struct {
	Users              map[int]users.User
	UsersMapByEmail    map[string]users.User
	GetUserByEmailFunc func(email string) (users.User, error)
	AddUserFunc        func(user users.User, toBeDeletedId int) (int, error)
	GetUserByIdFunc    func(id int) (users.User, error)
}

func (m *MockUserStore) GetUserById(id int) (users.User, error) {
	return m.GetUserByIdFunc(id)
}

func (m *MockUserStore) GetUserByEmail(email string) (users.User, error) {
	return m.GetUserByEmailFunc(email)
}

func (m *MockUserStore) AddUser(user users.User, toBeDeletedId int) (int, error) {
	return m.AddUserFunc(user, toBeDeletedId)
}

type MockEmailService struct {
	SendEmailFunc                    func(e email.Email) error
	BuildSignUpConfirmationEmailFunc func(email string) email.Email
}

func (m *MockEmailService) SendEmail(em email.Email) error {
	return m.SendEmailFunc(em)
}

func (m *MockEmailService) BuildSignUpConfirmationEmail(email string) email.Email {
	return m.BuildSignUpConfirmationEmailFunc(email)
}

type ErrorBody struct {
	Error   bool
	Message string
}
