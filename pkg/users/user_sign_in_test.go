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

func TestSignIn(t *testing.T) {

	tests := []struct {
		name           string
		payload        users.SignInRequest
		store          *MockUserStore
		expectedStatus int
		expectedError  error
	}{
		// Validate email
		{
			name: "should fail when email is empty",
			payload: users.SignInRequest{
				Email:    "",
				Password: "password",
			},
			store:          &MockUserStore{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.EmailRequiredError{},
		},
		{
			name: "should fail when email is too long",
			payload: users.SignInRequest{
				Email: `abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeab
				cdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdea
				bcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde
				abcde12345123451234512345123451234512345123451234512@test.com`,
				Password: "password",
			},
			store:          &MockUserStore{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.EmailTooLongError{},
		},
		{
			name: "should fail when email is invalid",
			payload: users.SignInRequest{
				Email:    `abc@`,
				Password: "password",
			},
			store:          &MockUserStore{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.InvalidEmailError{},
		},
		// Password validation
		{
			name: "should get an error for missing password",
			payload: users.SignInRequest{
				Email:    "a@a.com",
				Password: "",
			},
			store:          &MockUserStore{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.PasswordRequiredError{},
		},
		{
			name: "should get an error for too short password",
			payload: users.SignInRequest{
				Email:    "a@a.com",
				Password: "x",
			},
			store:          &MockUserStore{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.PasswordTooShortError{},
		},
		{
			name: "should get an error for too long password",
			payload: users.SignInRequest{
				Email:    "a@a.com",
				Password: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxw",
			},
			store:          &MockUserStore{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.PasswordTooLongError{},
		},
		// End of payload validation
		{
			name: "should fail to login when user not found",
			payload: users.SignInRequest{
				Email:    "not-exist@test.com",
				Password: "password",
			},
			store: &MockUserStore{
				GetUserByEmailFunc: func(email string) (users.User, error) {
					return users.User{}, sql.ErrNoRows
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "should fail to login when user is not activated",
			payload: users.SignInRequest{
				Email:    "not-activated@test.com",
				Password: "password",
			},
			store: &MockUserStore{
				GetUserByEmailFunc: func(email string) (users.User, error) {
					return users.User{
						Email:     "a@a.com",
						Password:  "$2a$10$sC6PANC9sIqpQWGVHku7Fu9vw4En4fGHLAioOkHPbJ7lZxOeKdB8G_testSalt",
						Activated: false,
					}, nil
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.UserNotActivatedError{},
		},
		{
			name: "should fail to login when password doesn't match",
			payload: users.SignInRequest{
				Email:    "a@a.com",
				Password: "password2", // "password" is the correct password
			},
			store: &MockUserStore{
				GetUserByEmailFunc: func(email string) (users.User, error) {
					return users.User{
						Email:     "a@a.com",
						Password:  "$2a$10$sC6PANC9sIqpQWGVHku7Fu9vw4En4fGHLAioOkHPbJ7lZxOeKdB8G_testSalt",
						Activated: true,
					}, nil
				},
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "should login successfully",
			payload: users.SignInRequest{
				Email:    "a@a.com",
				Password: "password",
			},
			store: &MockUserStore{
				GetUserByEmailFunc: func(email string) (users.User, error) {
					return users.User{
						Email:     "a@a.com",
						Password:  "$2a$10$sC6PANC9sIqpQWGVHku7Fu9vw4En4fGHLAioOkHPbJ7lZxOeKdB8G_testSalt",
						Activated: true,
					}, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := users.NewUserHandler(tt.store)
			reqPayload := signInPayloadToJSON(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", reqPayload)
			res := httptest.NewRecorder()
			handler.SignIn(res, req)
			assertStatus(t, res.Code, tt.expectedStatus)

			if tt.expectedError != nil {
				errBody := getErrorResponse(t, res)
				assertErrorMessage(t, errBody.Message, tt.expectedError.Error())
			}

		})
	}
}

func signInPayloadToJSON(payload users.SignInRequest) *strings.Reader {
	userJson, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}
	return strings.NewReader(string(userJson))
}
