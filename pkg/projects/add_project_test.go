package projects_test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/poomipat-k/running-fund/pkg/mock"
	"github.com/poomipat-k/running-fund/pkg/projects"
	s3Service "github.com/poomipat-k/running-fund/pkg/upload"
)

type ErrorBody struct {
	Error   bool
	Message string
}

func TestActivateEmail(t *testing.T) {

	tests := []struct {
		name           string
		payload        projects.AddProjectRequest
		store          *mock.MockProjectStore
		expectedStatus int
		expectedError  error
	}{
		{
			name:           "should error collaborated required",
			payload:        projects.AddProjectRequest{Collaborated: false},
			store:          &mock.MockProjectStore{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &projects.CollaboratedRequiredError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := tt.store
			userStore := &mock.MockUserStore{}
			handler := projects.NewProjectHandler(store, userStore, s3Service.S3Service{})

			reqPayload := addProjectPayloadToJSON(tt.payload)
			res := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/project", reqPayload)

			handler.AddProject(res, req)

			assertStatus(t, res.Code, tt.expectedStatus)
			if tt.expectedError != nil {
				errBody := getErrorResponse(t, res)
				assertErrorMessage(t, errBody.Message, tt.expectedError.Error())
			}
		})
	}
}

func assertErrorMessage(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct error, got %v, want %v", got, want)
	}
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
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

func addProjectPayloadToJSON(payload projects.AddProjectRequest) *strings.Reader {
	addProject, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}
	return strings.NewReader(string(addProject))
}
