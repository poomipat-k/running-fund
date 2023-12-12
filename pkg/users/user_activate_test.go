package users_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jordan-wright/email"
	"github.com/poomipat-k/running-fund/pkg/users"
)

func TestActivateEmail(t *testing.T) {
	store := &MockUserStore{
		GetUserByEmailFunc: func(email string) (users.User, error) {
			return users.User{}, sql.ErrNoRows
		},
		AddUserFunc: func(user users.User, toBeDeletedId int) (int, error) {
			return 1, nil
		},
	}
	es := &MockEmailService{
		BuildSignUpConfirmationEmailFunc: func(em, activateLink string) email.Email {
			return *email.NewEmail()
		},
		SendEmailFunc: func(e email.Email) error {
			return nil
		}}
	handler := users.NewUserHandler(store, es)

	req := httptest.NewRequest(http.MethodGet, "/auth/activate-email", nil)
	res := httptest.NewRecorder()

	handler.ActivateEmail(res, req)

	want := http.StatusOK
	assertStatus(t, res.Code, want)
}
