package users_test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/poomipat-k/running-fund/pkg/mock"
	"github.com/poomipat-k/running-fund/pkg/users"
)

func TestActivateEmail(t *testing.T) {

	tests := []struct {
		name                 string
		payload              users.ActivateUserRequest
		store                *mock.MockUserStore
		expectedStatus       int
		expectedError        error
		expectedEffectedRows *expectedEffectedRowsExpect
	}{
		{
			name:           "should error when activate code length is not equal to 24",
			payload:        users.ActivateUserRequest{ActivateCode: "abc"},
			store:          &mock.MockUserStore{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.InvalidActivateCodeError{},
		},
		{
			name:    "should error when activate code not found",
			payload: users.ActivateUserRequest{ActivateCode: "abcdabcdabcdabcdabcdabcd"},
			store: &mock.MockUserStore{
				ActivateUserFunc: func(activateCode string) (int64, error) {
					return 0, &users.UserToActivateNotFoundError{}
				},
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  &users.UserToActivateNotFoundError{},
		},
		{
			name:    "should error when exceed activate before",
			payload: users.ActivateUserRequest{ActivateCode: "abcdabcdabcdabcdabcdabcd"},
			store: &mock.MockUserStore{
				ActivateUserFunc: func(activateCode string) (int64, error) {
					return 0, &users.UserToActivateNotFoundError{}
				},
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  &users.UserToActivateNotFoundError{},
		},
		{
			name:    "should activate an account successfully",
			payload: users.ActivateUserRequest{ActivateCode: "abcdabcdabcdabcdabcdabcd"},
			store: &mock.MockUserStore{
				ActivateUserFunc: func(activateCode string) (int64, error) {
					return 1, nil
				},
			},
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

			reqPayload := activateUserPayloadToJSON(tt.payload)
			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/auth/activate-email", reqPayload)

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

func activateUserPayloadToJSON(payload users.ActivateUserRequest) *strings.Reader {
	activateRequest, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}
	return strings.NewReader(string(activateRequest))
}
