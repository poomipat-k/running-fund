package users_test

import (
	"encoding/json"
	"log"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/poomipat-k/running-fund/pkg/users"
)

type ErrorBody struct {
	Error   bool
	Message string
}

type expectedEffectedRowsExpect struct {
	rowEffected int
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
		t.Errorf("did not get correct error, got %v, want %v", got, want)
	}
}
