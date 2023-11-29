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

func (m *MockUserStore) AddUser() (int, error) {
	return 1, nil
}

type ErrorBody struct {
	Error   bool
	Message string
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
	t.Run("should get an error for duplicated email", func(t *testing.T) {
		store := &MockUserStore{
			GetUserByEmailFunc: func(email string) (users.User, error) {
				return users.User{}, nil
			},
		}
		handler := users.NewUserHandler(store)

		payload := users.SignUpRequest{
			Email:     "a@a.com",
			Password:  "password",
			FirstName: "x",
			LastName:  "l",
		}
		userJson, err := json.Marshal(payload)
		if err != nil {
			log.Fatal(err)
			return
		}
		req := httptest.NewRequest(http.MethodPost, "/user/register", strings.NewReader(string(userJson)))
		res := httptest.NewRecorder()

		handler.SignUp(res, req)
		assertStatus(t, res.Code, http.StatusBadRequest)

		errBody := getErrorResponse(t, res)
		assertErrorMessage(t, errBody.Message, "email is already exist")

	})

	t.Run("should get an error for missing last name", func(t *testing.T) {
		store := &MockUserStore{
			GetUserByEmailFunc: func(email string) (users.User, error) {
				return users.User{}, sql.ErrNoRows
			},
		}
		handler := users.NewUserHandler(store)

		payload := users.SignUpRequest{
			Email:     "a@a.com",
			Password:  "x",
			FirstName: "x",
			LastName:  "",
		}
		reqPayload := signUpPayloadToJSON(payload)

		req := httptest.NewRequest(http.MethodPost, "/user/register", reqPayload)
		res := httptest.NewRecorder()

		handler.SignUp(res, req)
		assertStatus(t, res.Code, http.StatusBadRequest)

		errBody := getErrorResponse(t, res)
		assertErrorMessage(t, errBody.Message, "last name is required")
	})

	t.Run("should sign up successfully", func(t *testing.T) {
		store := &MockUserStore{
			GetUserByEmailFunc: func(email string) (users.User, error) {
				return users.User{}, sql.ErrNoRows
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
