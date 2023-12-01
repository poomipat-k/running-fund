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

func TestSignIn(t *testing.T) {

	tests := []struct {
		name                 string
		payload              users.SignInRequest
		store                *MockUserStore
		expectedStatus       int
		expectedErrorMessage string
		expectedReturnId     int
	}{
		{
			name: "should login successfully",
			payload: users.SignInRequest{
				Email:    "a@a.com",
				Password: "password",
			},
			store: &MockUserStore{
				GetUserByEmailFunc: func(email string) (users.User, error) {
					return users.User{
						Email:    "a@a.com",
						Password: "$2a$10$sC6PANC9sIqpQWGVHku7Fu9vw4En4fGHLAioOkHPbJ7lZxOeKdB8G_testSalt",
					}, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "should fail to login when password doesn't match",
			payload: users.SignInRequest{
				Email:    "a@a.com",
				Password: "password2", // password is the right password
			},
			store: &MockUserStore{
				GetUserByEmailFunc: func(email string) (users.User, error) {
					return users.User{
						Email:    "a@a.com",
						Password: "$2a$10$sC6PANC9sIqpQWGVHku7Fu9vw4En4fGHLAioOkHPbJ7lZxOeKdB8G_testSalt",
					}, nil
				},
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "should fail when email is empty",
			payload: users.SignInRequest{
				Email:    "",
				Password: "password", // password is the right password
			},
			store: &MockUserStore{
				GetUserByEmailFunc: func(email string) (users.User, error) {
					return users.User{
						Email:    "a@a.com",
						Password: "$2a$10$sC6PANC9sIqpQWGVHku7Fu9vw4En4fGHLAioOkHPbJ7lZxOeKdB8G_testSalt",
					}, nil
				},
			},
			expectedStatus:       http.StatusBadRequest,
			expectedErrorMessage: "x",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := users.NewUserHandler(tt.store)
			reqPayload := signInPayloadToJSON(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/user/login", reqPayload)
			res := httptest.NewRecorder()
			handler.SignIn(res, req)
			assertStatus(t, res.Code, tt.expectedStatus)
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
