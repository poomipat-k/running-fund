package users_test

import "github.com/poomipat-k/running-fund/pkg/users"

type MockUserStore struct {
	GetReviewerByIdFunc func(id int) (users.User, error)
	Users               map[int]users.User
	UsersMapByEmail     map[string]users.User
	GetUserByEmailFunc  func(email string) (users.User, error)
	AddUserFunc         func(user users.User, toBeDeletedId int) (int, error)
}

func (m *MockUserStore) GetReviewerById(id int) (users.User, error) {
	return m.GetReviewerByIdFunc(id)
}

func (m *MockUserStore) GetReviewers() ([]users.User, error) {
	return []users.User{}, nil
}

func (m *MockUserStore) GetUserByEmail(email string) (users.User, error) {
	return m.GetUserByEmailFunc(email)
}

func (m *MockUserStore) AddUser(user users.User, toBeDeletedId int) (int, error) {
	return m.AddUserFunc(user, toBeDeletedId)
}

type ErrorBody struct {
	Error   bool
	Message string
}
