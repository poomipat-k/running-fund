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

type MockUserStore struct {
	GetReviewerByIdFunc func(id int) (users.User, error)
	Users               map[int]users.User
	UsersMapByEmail     map[string]users.User
	GetUserByEmailFunc  func(email string) (users.User, error)
	AddUserFunc         func(user users.User) (int, error)
}

func (m *MockUserStore) GetReviewerById(id int) (users.User, error) {
	return m.GetReviewerByIdFunc(id)
}

func (m *MockUserStore) GetReviewers() ([]users.User, error) {
	return []users.User{}, nil
}

func (m *MockUserStore) GetUserByEmail(email string) (users.User, error) {
	return m.GetUserByEmailFunc(email)
}

func (m *MockUserStore) AddUser(user users.User) (int, error) {
	return m.AddUserFunc(user)
}

type ErrorBody struct {
	Error   bool
	Message string
	Field   string
}

func TestGetReviewerById(t *testing.T) {
	t.Run("Get reviewer by id", func(t *testing.T) {
		user := users.User{Id: 1, FirstName: "aa", LastName: "bb", Email: "test@test.com", UserRole: "applicant"}
		expectedEmail := "test@test.com"
		userMap := map[int]users.User{
			1: user,
		}

		store := &MockUserStore{
			GetReviewerByIdFunc: func(id int) (users.User, error) {
				if _, valid := userMap[id]; !valid {
					return users.User{}, sql.ErrNoRows
				}
				return userMap[id], nil
			},
			Users: userMap,
		}
		handler := users.NewUserHandler(store)

		req := httptest.NewRequest(http.MethodGet, "/user/reviewer", nil)
		req.Header.Add("Authorization", "Bearer 1")
		res := httptest.NewRecorder()

		handler.GetReviewerById(res, req)

		if res.Code != http.StatusOK {
			t.Errorf("got %v, want %v", res.Code, http.StatusOK)
		}

		var got users.User
		err := json.Unmarshal(res.Body.Bytes(), &got)
		if err != nil {
			log.Fatal("Error decoding a user")
		}
		if got.Email != expectedEmail {
			t.Errorf("got %v, want %v", got.Id, expectedEmail)
		}
	})

	t.Run("returns error if not found", func(t *testing.T) {
		user := users.User{Id: 1, FirstName: "aa", LastName: "bb", Email: "test@test.com", UserRole: "applicant"}
		userMap := map[int]users.User{
			1: user,
		}

		store := &MockUserStore{
			GetReviewerByIdFunc: func(id int) (users.User, error) {
				if _, valid := userMap[id]; !valid {
					return users.User{}, sql.ErrNoRows
				}
				return userMap[id], nil
			},
			Users: userMap,
		}
		handler := users.NewUserHandler(store)

		req := httptest.NewRequest(http.MethodGet, "/user/reviewer", nil)
		req.Header.Add("Authorization", "Bearer 2")
		res := httptest.NewRecorder()

		handler.GetReviewerById(res, req)

		if res.Code != http.StatusBadRequest {
			t.Errorf("got %v, want %v", res.Code, http.StatusBadRequest)
		}

		errBody := getErrorResponse(t, res)
		if !errBody.Error {
			t.Errorf("expected to get an error")
		}
	})
}

func TestSignUp(t *testing.T) {

	tests := []struct {
		name                 string
		payload              users.SignUpRequest
		store                *MockUserStore
		expectedStatus       int
		expectedErrorMessage string
	}{
		{
			name: "should get an error for duplicated email",
			payload: users.SignUpRequest{
				Email:     "a@a.com",
				Password:  "password",
				FirstName: "x",
				LastName:  "l",
			},
			store: &MockUserStore{
				GetUserByEmailFunc: func(email string) (users.User, error) {
					return users.User{}, nil
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

			errBody := getErrorResponse(t, res)
			assertErrorMessage(t, errBody.Message, tt.expectedErrorMessage)
		})
	}

	t.Run("should sign up successfully", func(t *testing.T) {
		store := &MockUserStore{
			GetUserByEmailFunc: func(email string) (users.User, error) {
				return users.User{}, sql.ErrNoRows
			},
			AddUserFunc: func(user users.User) (int, error) {
				return 1, nil
			},
		}
		handler := users.NewUserHandler(store)

		payload := users.SignUpRequest{
			Email:     "a@a.com",
			Password:  "password",
			FirstName: "x",
			LastName:  "l",
		}

		reqPayload := signUpPayloadToJSON(payload)

		req := httptest.NewRequest(http.MethodPost, "/user/register", reqPayload)
		res := httptest.NewRecorder()

		handler.SignUp(res, req)

		assertStatus(t, res.Code, http.StatusCreated)
		var got int
		err := json.Unmarshal(res.Body.Bytes(), &got)
		if err != nil {
			t.Errorf("fail to unmarshal err: %+v", err)
		}
		want := 1
		if got != want {
			t.Errorf("user id did not match, got %d, want %d", got, want)
		}
	})
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
