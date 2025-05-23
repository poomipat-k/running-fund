package users_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/poomipat-k/running-fund/pkg/mock"
	"github.com/poomipat-k/running-fund/pkg/users"
)

func TestSignUp(t *testing.T) {

	tests := []struct {
		name             string
		payload          users.SignUpRequest
		store            *mock.MockUserStore
		expectedStatus   int
		expectedError    error
		expectedReturnId int
	}{
		// Basic email validation
		{
			name: "should get an error for missing email",
			payload: users.SignUpRequest{
				Email:             "",
				Password:          "password",
				FirstName:         "abc",
				LastName:          "ab",
				TermsAndCondition: true,
				Privacy:           true,
			},
			store:          &mock.MockUserStore{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.EmailRequiredError{},
		},
		{
			name: "should get an error for too long email",
			payload: users.SignUpRequest{
				Email: `abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeab
				cdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdea
				bcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde
				abcde12345123451234512345123451234512345123451234512@test.com`,
				Password:          "bad-example",
				FirstName:         "x",
				LastName:          "y",
				TermsAndCondition: true,
				Privacy:           true,
			},
			store:          &mock.MockUserStore{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.EmailTooLongError{},
		},
		{
			name: "should get an error for invalid email",
			payload: users.SignUpRequest{
				Email:             "abc@",
				Password:          "bad-example",
				FirstName:         "x",
				LastName:          "y",
				TermsAndCondition: true,
				Privacy:           true,
			},
			store: &mock.MockUserStore{
				GetUserByEmailFunc: func(email string) (users.User, error) {
					return users.User{}, sql.ErrNoRows
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.InvalidEmailError{},
		},
		// Password validation
		{
			name: "should get an error for missing password",
			payload: users.SignUpRequest{
				Email:             "a@a.com",
				Password:          "",
				FirstName:         "abc",
				LastName:          "ab",
				TermsAndCondition: true,
				Privacy:           true,
			},
			store:          &mock.MockUserStore{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.PasswordRequiredError{},
		},
		{
			name: "should get an error for too short password",
			payload: users.SignUpRequest{
				Email:             "a@a.com",
				Password:          "x",
				FirstName:         "x",
				LastName:          "y",
				TermsAndCondition: true,
				Privacy:           true,
			},
			store: &mock.MockUserStore{
				GetUserByEmailFunc: func(email string) (users.User, error) {
					return users.User{}, sql.ErrNoRows
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.PasswordTooShortError{},
		},
		{
			name: "should get an error for too long password",
			payload: users.SignUpRequest{
				Email:             "a@a.com",
				Password:          "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxw",
				FirstName:         "x",
				LastName:          "y",
				TermsAndCondition: true,
				Privacy:           true,
			},
			store: &mock.MockUserStore{
				GetUserByEmailFunc: func(email string) (users.User, error) {
					return users.User{}, sql.ErrNoRows
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.PasswordTooLongError{},
		},
		// Validate first name
		{
			name: "should get an error for missing first name",
			payload: users.SignUpRequest{
				Email:             "a@a.com",
				Password:          "password",
				FirstName:         "",
				LastName:          "ab",
				TermsAndCondition: true,
				Privacy:           true,
			},
			store: &mock.MockUserStore{
				GetUserByEmailFunc: func(email string) (users.User, error) {
					return users.User{}, sql.ErrNoRows
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.FirstNameRequiredError{},
		},
		{
			name: "should get an error for too long first name",
			payload: users.SignUpRequest{
				Email:    "abc@test.com",
				Password: "password",
				FirstName: `Lorem Ipsum is simply dummy text of the printing and typesetting industry.
				Lorem Ipsum has been the industry's standard dummy text ever since the 1500s,
				when an unknown printer took a galley of type and scrambled it to make a type specimen book.
				It has survived not only five centuries, but also the leap into electronic typesetting,
				remaining essentially unchanged. It was popularised in the 1960s with
				the release of Letraset sheets containing Lorem Ipsum passages, and more recently
				with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.
				Lorem Ipsum is simply dummy text of the printing and typesetting industry.
				Lorem Ipsum has been the industry's standard dummy text ever since the 1500s,
				when an unknown printer took a galley of type and scrambled it to make a type specimen book.`,
				LastName:          "test",
				TermsAndCondition: true,
				Privacy:           true,
			},
			store: &mock.MockUserStore{
				GetUserByEmailFunc: func(email string) (users.User, error) {
					return users.User{}, sql.ErrNoRows
				},
				AddUserFunc: func(user users.User, toBeDeletedId int) (int, string, error) {
					return 1, "", nil
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.FirstNameTooLongError{},
		},
		// Validate last name
		{
			name: "should get an error for missing last name",
			payload: users.SignUpRequest{
				Email:             "a@a.com",
				Password:          "password",
				FirstName:         "x",
				LastName:          "",
				TermsAndCondition: true,
				Privacy:           true,
			},
			store: &mock.MockUserStore{
				GetUserByEmailFunc: func(email string) (users.User, error) {
					return users.User{}, sql.ErrNoRows
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.LastNameRequiredError{},
		},
		{
			name: "should get an error for too long last name",
			payload: users.SignUpRequest{
				Email:     "abc@test.com",
				Password:  "password",
				FirstName: "last",
				LastName: `Lorem Ipsum is simply dummy text of the printing and typesetting industry.
				Lorem Ipsum has been the industry's standard dummy text ever since the 1500s,
				when an unknown printer took a galley of type and scrambled it to make a type specimen book.
				It has survived not only five centuries, but also the leap into electronic typesetting,
				remaining essentially unchanged. It was popularised in the 1960s with
				the release of Letraset sheets containing Lorem Ipsum passages, and more recently
				with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.
				Lorem Ipsum is simply dummy text of the printing and typesetting industry.
				Lorem Ipsum has been the industry's standard dummy text ever since the 1500s,
				when an unknown printer took a galley of type and scrambled it to make a type specimen book.`,
				TermsAndCondition: true,
				Privacy:           true,
			},
			store: &mock.MockUserStore{
				GetUserByEmailFunc: func(email string) (users.User, error) {
					return users.User{}, sql.ErrNoRows
				},
				AddUserFunc: func(user users.User, toBeDeletedId int) (int, string, error) {
					return 1, "", nil
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.LastNameTooLongError{},
		},
		// Validate duplicated email
		{
			name: "should get an error for duplicated email - found activated email",
			payload: users.SignUpRequest{
				Email:             "a@a.com",
				Password:          "password",
				FirstName:         "x",
				LastName:          "l",
				TermsAndCondition: true,
				Privacy:           true,
			},
			store: &mock.MockUserStore{
				GetUserByEmailFunc: func(email string) (users.User, error) {
					return users.User{Activated: true}, nil
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.DuplicatedEmailError{},
		},
		{
			name: "should get an error for duplicated email - not activated but before activate_before ends",
			payload: users.SignUpRequest{
				Email:             "a@a.com",
				Password:          "password",
				FirstName:         "x",
				LastName:          "l",
				TermsAndCondition: true,
				Privacy:           true,
			},
			store: &mock.MockUserStore{
				GetUserByEmailFunc: func(email string) (users.User, error) {
					return users.User{Activated: false, ActivatedBefore: time.Now().Local().Add(time.Duration(24 * time.Hour))}, nil
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.DuplicatedEmailError{},
		},
		{
			name: "should get an error for not acknowledging terms and condition",
			payload: users.SignUpRequest{
				Email:             "a@a.com",
				Password:          "password",
				FirstName:         "x",
				LastName:          "b",
				TermsAndCondition: false,
				Privacy:           true,
			},
			store: &mock.MockUserStore{
				GetUserByEmailFunc: func(email string) (users.User, error) {
					return users.User{}, sql.ErrNoRows
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.MissingTermsAndConditionError{},
		},
		{
			name: "should get an error for not acknowledging privacy policy",
			payload: users.SignUpRequest{
				Email:             "a@a.com",
				Password:          "password",
				FirstName:         "x",
				LastName:          "b",
				TermsAndCondition: true,
				Privacy:           false,
			},
			store: &mock.MockUserStore{
				GetUserByEmailFunc: func(email string) (users.User, error) {
					return users.User{}, sql.ErrNoRows
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &users.MissingPrivacyError{},
		},
		// Success sign up
		{
			name: "should sign up successfully",
			payload: users.SignUpRequest{
				Email:             "a@a.com",
				Password:          "password",
				FirstName:         "x",
				LastName:          "l",
				TermsAndCondition: true,
				Privacy:           true,
			},
			store: &mock.MockUserStore{
				GetUserByEmailFunc: func(email string) (users.User, error) {
					return users.User{}, sql.ErrNoRows
				},
				AddUserFunc: func(user users.User, toBeDeletedId int) (int, string, error) {
					return 1, "", nil
				},
			},
			expectedStatus:   http.StatusCreated,
			expectedReturnId: 1,
		},
		{
			name: "should sign up successfully when email already exist but that user is not activated and activate_before is less than now",
			payload: users.SignUpRequest{
				Email:             "a@a.com",
				Password:          "password",
				FirstName:         "x",
				LastName:          "l",
				TermsAndCondition: true,
				Privacy:           true,
			},
			store: &mock.MockUserStore{
				GetUserByEmailFunc: func(email string) (users.User, error) {
					return users.User{Id: 1, Activated: false, ActivatedBefore: time.Now().Local().Add(time.Duration(-24 * time.Hour))}, nil
				},
				AddUserFunc: func(user users.User, toBeDeletedId int) (int, string, error) {
					return 2, "", nil
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

			if tt.expectedError != nil {
				errBody := getErrorResponse(t, res)
				assertErrorMessage(t, errBody.Message, tt.expectedError.Error())
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
