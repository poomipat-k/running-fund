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
	s3Service "github.com/poomipat-k/running-fund/pkg/s3-service"
)

type UpdateProjectTestCase struct {
	name           string
	payload        projects.AdminUpdateProjectRequest
	store          *mock.MockProjectStore
	expectedStatus int
	expectedError  error
}

func TestAdminUpdateProject(t *testing.T) {
	tests := []UpdateProjectTestCase{
		{
			name:           "should error when projectStatusPrimary is missing",
			payload:        projects.AdminUpdateProjectRequest{},
			store:          &mock.MockProjectStore{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &projects.ProjectStatusPrimaryRequiredError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := tt.store
			userStore := &mock.MockUserStore{}
			handler := projects.NewProjectHandler(store, userStore, s3Service.S3Service{})
			reqPayload := adminUpdateProjectPayloadToJSON(tt.payload)
			req := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/project/{projectCode}", reqPayload)
			res := httptest.NewRecorder()
			handler.AdminUpdateProject(res, req)
			assertStatus(t, res.Code, tt.expectedStatus)
		})
	}
}

func adminUpdateProjectPayloadToJSON(payload projects.AdminUpdateProjectRequest) *strings.Reader {
	userJson, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}
	return strings.NewReader(string(userJson))
}
