package users_test

import (
	"encoding/json"
	"log"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jordan-wright/email"
	"github.com/poomipat-k/running-fund/pkg/users"
)

type MockUserStore struct {
	Users                    map[int]users.User
	UsersMapByEmail          map[string]users.User
	GetUserByEmailFunc       func(email string) (users.User, error)
	AddUserFunc              func(user users.User, toBeDeletedId int) (int, error)
	GetUserByIdFunc          func(id int) (users.User, error)
	ActivateUserFunc         func(activateCode string) (int64, error)
	ForgotPasswordActionFunc func(resetPasswordCode string, email string, resetPasswordLink string) (int64, error)
	ResetPasswordFunc        func(resetPasswordCode string, newPassword string) (int64, error)
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

func (m *MockUserStore) ActivateUser(activateCode string) (int64, error) {
	return m.ActivateUserFunc(activateCode)
}

func (m *MockUserStore) ForgotPasswordAction(resetPasswordCode string, email string, resetPasswordLink string) (int64, error) {
	return m.ForgotPasswordActionFunc(resetPasswordCode, email, resetPasswordLink)
}

func (m *MockUserStore) ResetPassword(resetPasswordCode string, newPassword string) (int64, error) {
	return m.ResetPasswordFunc(resetPasswordCode, newPassword)
}

type MockEmailService struct {
	SendEmailFunc                    func(e email.Email) error
	BuildSignUpConfirmationEmailFunc func(to, activateLink string) email.Email
	BuildResetPasswordEmailFunc      func(to, resetPasswordLink string) email.Email
}

func (m *MockEmailService) SendEmail(em email.Email) error {
	return m.SendEmailFunc(em)
}

func (m *MockEmailService) BuildSignUpConfirmationEmail(to, activateLink string) email.Email {
	return m.BuildSignUpConfirmationEmailFunc(to, activateLink)
}

func (m *MockEmailService) BuildResetPasswordEmail(to, resetPasswordLink string) email.Email {
	return m.BuildResetPasswordEmailFunc(to, resetPasswordLink)
}

type ErrorBody struct {
	Error   bool
	Message string
}

type expectedEffectedRowsExpect struct {
	rowEffected int
}

func getErrorResponse(t testing.TB, res *httptest.ResponseRecorder) ErrorBody {
	t.Helper()
	var body ErrorBody
	err := json.Unmarshal(res.Body.Bytes(), &body)
	if err != nil {
		t.Errorf("Error unmarshal ErrorResponse")
	}
	return body
}

func signUpPayloadToJSON(payload users.SignUpRequest) *strings.Reader {
	userJson, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}
	return strings.NewReader(string(userJson))
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertErrorMessage(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct error, got %v, want %v", got, want)
	}
}
