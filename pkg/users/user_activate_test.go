package users_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/poomipat-k/running-fund/pkg/users"
)

func TestActivateEmail(t *testing.T) {

	tests := []struct {
		name                 string
		activateCode         string
		store                *MockUserStore
		emailService         *MockEmailService
		expectedStatus       int
		expectedError        error
		expectedEffectedRows *expectedEffectedRowsExpect
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
				ActivateUserFunc: func(activateCode string) (int64, error) {
					return 0, &users.UserToActivateNotFoundError{}
				},
			},
			emailService:   &MockEmailService{},
			expectedStatus: http.StatusNotFound,
			expectedError:  &users.UserToActivateNotFoundError{},
		},
		{
			name:         "should error when exceed activate before",
			activateCode: "abcdabcdabcdabcdabcdabcd",
			store: &MockUserStore{
				ActivateUserFunc: func(activateCode string) (int64, error) {
					return 0, &users.UserToActivateNotFoundError{}
				},
			},
			emailService:   &MockEmailService{},
			expectedStatus: http.StatusNotFound,
			expectedError:  &users.UserToActivateNotFoundError{},
		},
		{
			name:         "should activate an account successfully",
			activateCode: "abcdabcdabcdabcdabcdabcd",
			store: &MockUserStore{
				ActivateUserFunc: func(activateCode string) (int64, error) {
					return 1, nil
				},
			},
			emailService:   &MockEmailService{},
			expectedStatus: http.StatusOK,
			expectedEffectedRows: &expectedEffectedRowsExpect{
				rowEffected: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := tt.store
			handler := users.NewUserHandler(store)

			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/auth/activate-email?&activateCode=%s", tt.activateCode), nil)

			handler.ActivateUser(res, req)

			assertStatus(t, res.Code, tt.expectedStatus)
			if tt.expectedError != nil {
				errBody := getErrorResponse(t, res)
				assertErrorMessage(t, errBody.Message, tt.expectedError.Error())
			}

			if tt.expectedEffectedRows != nil {
				rowEffected, err := strconv.Atoi(res.Body.String())
				if err != nil {
					t.Errorf(err.Error())
				}
				if rowEffected != tt.expectedEffectedRows.rowEffected {
					t.Errorf("row effect mismatch got %d, want %d", rowEffected, tt.expectedEffectedRows.rowEffected)
				}
			}
		})
	}

}
