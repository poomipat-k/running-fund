package users_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/poomipat-k/running-fund/pkg/users"
)

func TestActivateEmail(t *testing.T) {
	tests := []struct {
		name           string
		activateCode   string
		store          *MockUserStore
		emailService   *MockEmailService
		expectedStatus int
		expectedError  error
	}{
		{
			name:           "should error when activate code length is not equal to 24",
			activateCode:   "abc",
			store:          &MockUserStore{},
			emailService:   &MockEmailService{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.InvalidActivateCodeError{},
		},
		{
			name:         "should error when activate code not found",
			activateCode: "abcdabcdabcdabcdabcdabcd",
			store: &MockUserStore{
				ActivateUserFunc: func(activateCode string) (int, error) {
					return 0, &users.UserToActivateNotFoundError{}
				},
			},
			emailService:   &MockEmailService{},
			expectedStatus: http.StatusNotFound,
			expectedError:  &users.UserToActivateNotFoundError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := tt.store
			es := &MockEmailService{}
			handler := users.NewUserHandler(store, es)

			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/auth/activate-email?&activateCode=%s", tt.activateCode), nil)

			handler.ActivateUser(res, req)

			assertStatus(t, res.Code, tt.expectedStatus)
			if tt.expectedError != nil {
				errBody := getErrorResponse(t, res)
				assertErrorMessage(t, errBody.Message, tt.expectedError.Error())
			}
		})
	}

}
