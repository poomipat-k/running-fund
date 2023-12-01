package users_test

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/poomipat-k/running-fund/pkg/users"
)

func TestSignUp(t *testing.T) {

	tests := []struct {
		name                 string
		payload              users.SignUpRequest
		store                *MockUserStore
		expectedStatus       int
		expectedErrorMessage string
		expectedReturnId     int
	}{
		{
			name: "should get an error for duplicated email - activated",
			payload: users.SignUpRequest{
				Email:     "a@a.com",
				Password:  "password",
				FirstName: "x",
				LastName:  "l",
			},
			store: &MockUserStore{
				GetUserByEmailFunc: func(email string) (users.User, error) {
					return users.User{Activated: true}, nil
				},
			},
			expectedStatus:       http.StatusBadRequest,
			expectedErrorMessage: "email is already exist",
		},
		{
			name: "should get an error for duplicated email - not activated but before activate_before ends",
			payload: users.SignUpRequest{
				Email:     "a@a.com",
				Password:  "password",
				FirstName: "x",
				LastName:  "l",
			},
			store: &MockUserStore{
				GetUserByEmailFunc: func(email string) (users.User, error) {
					return users.User{Activated: false, ActivatedBefore: time.Now().Local().Add(time.Duration(24 * time.Hour))}, nil
				},
			},
			expectedStatus:       http.StatusBadRequest,
			expectedErrorMessage: "email is already exist",
		},
		{
			name: "should get an error for invalid email",
			payload: users.SignUpRequest{
				Email:     "abc@",
				Password:  "bad-example",
				FirstName: "x",
				LastName:  "y",
			},
			store: &MockUserStore{
				GetUserByEmailFunc: func(email string) (users.User, error) {
					return users.User{}, sql.ErrNoRows
				},
			},
			expectedStatus:       http.StatusBadRequest,
			expectedErrorMessage: "email is invalid",
		},
		{
			name: "should get an error for missing last name",
			payload: users.SignUpRequest{
				Email:     "a@a.com",
				Password:  "password",
				FirstName: "x",
				LastName:  "",
			},
			store: &MockUserStore{
				GetUserByEmailFunc: func(email string) (users.User, error) {
					return users.User{}, sql.ErrNoRows
				},
			},
			expectedStatus:       http.StatusBadRequest,
			expectedErrorMessage: "last name is required",
		},
		{
			name: "should get an error for too short password",
			payload: users.SignUpRequest{
				Email:     "a@a.com",
				Password:  "x",
				FirstName: "x",
				LastName:  "y",
			},
			store: &MockUserStore{
				GetUserByEmailFunc: func(email string) (users.User, error) {
					return users.User{}, sql.ErrNoRows
				},
			},
			expectedStatus:       http.StatusBadRequest,
			expectedErrorMessage: "password minimum length are 8 characters",
		},
		{
			name: "should get an error for too long password",
			payload: users.SignUpRequest{
				Email:     "a@a.com",
				Password:  "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxw",
				FirstName: "x",
				LastName:  "y",
			},
			store: &MockUserStore{
				GetUserByEmailFunc: func(email string) (users.User, error) {
					return users.User{}, sql.ErrNoRows
				},
			},
			expectedStatus:       http.StatusBadRequest,
			expectedErrorMessage: "password maximum length are 60 characters",
		},
		{
			name: "should sign up successfully",
			payload: users.SignUpRequest{
				Email:     "a@a.com",
				Password:  "password",
				FirstName: "x",
				LastName:  "l",
			},
			store: &MockUserStore{
				GetUserByEmailFunc: func(email string) (users.User, error) {
					return users.User{}, sql.ErrNoRows
				},
				AddUserFunc: func(user users.User, toBeDeletedId int) (int, error) {
					return 1, nil
				},
			},
			expectedStatus:   http.StatusCreated,
			expectedReturnId: 1,
		},
		{
			name: "should sign up successfully when email already exist but that user is not activated and activate_before is less than now",
			payload: users.SignUpRequest{
				Email:     "a@a.com",
				Password:  "password",
				FirstName: "x",
				LastName:  "l",
			},
			store: &MockUserStore{
				GetUserByEmailFunc: func(email string) (users.User, error) {
					return users.User{Id: 1, Activated: false, ActivatedBefore: time.Now().Local().Add(time.Duration(-24 * time.Hour))}, nil
				},
				AddUserFunc: func(user users.User, toBeDeletedId int) (int, error) {
					return 2, nil
				},
			},
			expectedStatus:   http.StatusCreated,
			expectedReturnId: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := tt.store
			handler := users.NewUserHandler(store)

			reqPayload := signUpPayloadToJSON(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/user/register", reqPayload)
			res := httptest.NewRecorder()

			handler.SignUp(res, req)
			assertStatus(t, res.Code, tt.expectedStatus)

			if tt.expectedErrorMessage != "" {
				errBody := getErrorResponse(t, res)
				assertErrorMessage(t, errBody.Message, tt.expectedErrorMessage)
			}
			if tt.expectedReturnId > 0 {
				var got int
				err := json.Unmarshal(res.Body.Bytes(), &got)
				if err != nil {
					t.Errorf("fail to unmarshal err: %+v", err)
				}
				if got != tt.expectedReturnId {
					t.Errorf("user id did not match, got %d, want %d", got, tt.expectedReturnId)
				}
			}
		})
	}

}

func getErrorResponse(t testing.TB, res *httptest.ResponseRecorder) ErrorBody {
	t.Helper()
	var body ErrorBody
	err := json.Unmarshal(res.Body.Bytes(), &body)
	if err != nil {
		t.Errorf("Error unmarshal ErrorResponse")
	}
	return body
}

func signUpPayloadToJSON(payload users.SignUpRequest) *strings.Reader {
	userJson, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}
	return strings.NewReader(string(userJson))
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertErrorMessage(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct error message, got %v, want %v", got, want)
	}
}
