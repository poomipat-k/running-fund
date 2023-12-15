package users_test

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jordan-wright/email"
	"github.com/poomipat-k/running-fund/pkg/users"
)

func TestEmailForgetPassword(t *testing.T) {

	tests := []struct {
		name                  string
		forgotPasswordPayload users.ForgotPasswordRequest
		store                 *MockUserStore
		emailService          *MockEmailService
		expectedStatus        int
		expectedError         error
	}{
		{
			name:           "should error when email is not provided",
			store:          &MockUserStore{},
			emailService:   &MockEmailService{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.EmailRequiredError{},
		},
		{
			name: "should error when email is not found",
			store: &MockUserStore{
				GetUserByEmailFunc: func(email string) (users.User, error) {
					return users.User{}, sql.ErrNoRows
				},
			},
			forgotPasswordPayload: users.ForgotPasswordRequest{
				Email: "abc@test.com",
			},
			emailService:   &MockEmailService{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  sql.ErrNoRows,
		},
		{
			name: "should error when user is not activated",
			store: &MockUserStore{
				GetUserByEmailFunc: func(email string) (users.User, error) {
					return users.User{Activated: false}, nil
				},
			},
			forgotPasswordPayload: users.ForgotPasswordRequest{
				Email: "abc@test.com",
			},
			emailService:   &MockEmailService{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.UserIsNotActivated{},
		},
		{
			name: "should send create a new password successfully",
			store: &MockUserStore{
				GetUserByEmailFunc: func(email string) (users.User, error) {
					return users.User{Activated: true}, nil
				},
				ForgotPasswordActionFunc: func(resetPasswordCode, email, resetPasswordLink string) (int64, error) {
					return 1, nil
				},
			},
			forgotPasswordPayload: users.ForgotPasswordRequest{
				Email: "abc@test.com",
			},
			emailService: &MockEmailService{
				BuildResetPasswordEmailFunc: func(to string, activateLink string) email.Email {
					return *email.NewEmail()
				},
				SendEmailFunc: func(e email.Email) error {
					return nil
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := tt.store
			handler := users.NewUserHandler(store, tt.emailService)

			res := httptest.NewRecorder()
			payload := forgotPasswordPayloadToJSON(tt.forgotPasswordPayload)
			req := httptest.NewRequest(http.MethodPost, "/user/forgot-password", payload)

			handler.ForgotPassword(res, req)

			assertStatus(t, res.Code, tt.expectedStatus)
			if tt.expectedError != nil {
				errBody := getErrorResponse(t, res)
				assertErrorMessage(t, errBody.Message, tt.expectedError.Error())
			}
		})
	}

}

func forgotPasswordPayloadToJSON(payload users.ForgotPasswordRequest) *strings.Reader {
	emailJson, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}
	return strings.NewReader(string(emailJson))
}
