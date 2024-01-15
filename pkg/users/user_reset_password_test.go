package users_test

import (
	"encoding/json"
	"errors"
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
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.PasswordAndConfirmPasswordNotMatchError{},
		},
		{
			name: "should error when password is not provided",
			resetPasswordPayload: users.ResetPasswordRequest{
				Password:        "",
				ConfirmPassword: "",
			},
			store:          &MockUserStore{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.PasswordRequiredError{},
		},
		{
			name: "should error when password is too short",
			resetPasswordPayload: users.ResetPasswordRequest{
				Password:        "ab",
				ConfirmPassword: "ab",
			},
			store:          &MockUserStore{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.PasswordTooShortError{},
		},
		{
			name: "should error when password is too long",
			resetPasswordPayload: users.ResetPasswordRequest{
				Password:        "abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345a",
				ConfirmPassword: "abcde12345abcde12345abcde12345abcde12345abcde12345abcde12345a",
			},
			store:          &MockUserStore{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.PasswordTooLongError{},
		},
		{
			name: "should error when reset password code is not valid",
			resetPasswordPayload: users.ResetPasswordRequest{
				Password:          "abcd1234",
				ConfirmPassword:   "abcd1234",
				ResetPasswordCode: "code",
			},
			store:          &MockUserStore{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.ResetPasswordCodeNotValidError{},
		},
		{
			name: "should error when reset password code is not valid",
			resetPasswordPayload: users.ResetPasswordRequest{
				Password:          "abcd1234",
				ConfirmPassword:   "abcd1234",
				ResetPasswordCode: "abcdefghabcdefghabcdefgh",
			},
			store: &MockUserStore{
				ResetPasswordFunc: func(resetPasswordCode, newPassword string) (int64, error) {
					return 0, errors.New("abc")
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  errors.New("abc"),
		},
		{
			name: "should error when reset password failed and 0 row updated",
			resetPasswordPayload: users.ResetPasswordRequest{
				Password:          "abcd1234",
				ConfirmPassword:   "abcd1234",
				ResetPasswordCode: "abcdefghabcdefghabcdefgh",
			},
			store: &MockUserStore{
				ResetPasswordFunc: func(resetPasswordCode, newPassword string) (int64, error) {
					return 0, nil
				},
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  &users.ResetPasswordCodeNotFound{},
		},
		{
			name: "should reset password successfully",
			resetPasswordPayload: users.ResetPasswordRequest{
				Password:          "abcd1234",
				ConfirmPassword:   "abcd1234",
				ResetPasswordCode: "abcdefghabcdefghabcdefgh",
			},
			store: &MockUserStore{
				ResetPasswordFunc: func(resetPasswordCode, newPassword string) (int64, error) {
					return 1, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := tt.store
			handler := users.NewUserHandler(store)

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
