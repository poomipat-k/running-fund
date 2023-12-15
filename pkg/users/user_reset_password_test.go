package users_test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/poomipat-k/running-fund/pkg/users"
)

func TestResetPassword(t *testing.T) {

	tests := []struct {
		name                 string
		resetPasswordPayload users.ResetPasswordRequest
		store                *MockUserStore
		emailService         *MockEmailService
		expectedStatus       int
		expectedError        error
	}{
		{
			name: "should error when password and confirmPassword is not identical",
			resetPasswordPayload: users.ResetPasswordRequest{
				Password:        "abc",
				ConfirmPassword: "abb",
			},
			store:          &MockUserStore{},
			emailService:   &MockEmailService{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.PasswordAndConfirmPasswordNotMatch{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := tt.store
			handler := users.NewUserHandler(store, tt.emailService)

			res := httptest.NewRecorder()
			payload := resetPasswordPayloadToJSON(tt.resetPasswordPayload)
			req := httptest.NewRequest(http.MethodPost, "/user/reset-password", payload)

			handler.ResetPassword(res, req)

			assertStatus(t, res.Code, tt.expectedStatus)
			if tt.expectedError != nil {
				errBody := getErrorResponse(t, res)
				assertErrorMessage(t, errBody.Message, tt.expectedError.Error())
			}
		})
	}

}

func resetPasswordPayloadToJSON(payload users.ResetPasswordRequest) *strings.Reader {
	emailJson, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}
	return strings.NewReader(string(emailJson))
}
