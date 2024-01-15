package users_test

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/poomipat-k/running-fund/pkg/users"
)

func TestEmailForgetPassword(t *testing.T) {

	tests := []struct {
		name                  string
		forgotPasswordPayload users.ForgotPasswordRequest
		store                 *MockUserStore
		expectedStatus        int
		expectedError         error
	}{
		{
			name:           "should error when email is not provided",
			store:          &MockUserStore{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.EmailRequiredError{},
		},
		{
			name: "should error when email is too long",
			forgotPasswordPayload: users.ForgotPasswordRequest{
				Email: `abcdeabcdeabcdeabcdeabcdeabcdeabcd
				eabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeab
				cdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcd
				eabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdea
				bcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabc
				deabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde@test.com`,
			},
			store:          &MockUserStore{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.EmailTooLongError{},
		},
		{
			name: "should error when email is not valid",
			forgotPasswordPayload: users.ForgotPasswordRequest{
				Email: `aab@`,
			},
			store:          &MockUserStore{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.InvalidEmailError{},
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
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.UserIsNotActivatedError{},
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
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := tt.store
			handler := users.NewUserHandler(store)

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
