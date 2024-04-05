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
		{
			name: "should error when projectStatusPrimary is missing",
			payload: projects.AdminUpdateProjectRequest{
				ProjectStatusPrimary: "CurrentBeforeApprove",
			},
			store:          &mock.MockProjectStore{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &projects.ProjectStatusSecondaryRequiredError{},
		},
		{
			name: "should error when adminScore is less than 0",
			payload: projects.AdminUpdateProjectRequest{
				ProjectStatusPrimary:   "CurrentBeforeApprove",
				ProjectStatusSecondary: "Reviewing",
				AdminScore:             newInt(-10),
			},
			store:          &mock.MockProjectStore{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &projects.AdminScoreOutOfRangeError{},
		},
		{
			name: "should error when adminScore is greater than 100",
			payload: projects.AdminUpdateProjectRequest{
				ProjectStatusPrimary:   "CurrentBeforeApprove",
				ProjectStatusSecondary: "Reviewing",
				AdminScore:             newInt(110),
			},
			store:          &mock.MockProjectStore{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &projects.AdminScoreOutOfRangeError{},
		},
		{
			name: "should error when fundApprovedAmount is less than 0",
			payload: projects.AdminUpdateProjectRequest{
				ProjectStatusPrimary:   "CurrentBeforeApprove",
				ProjectStatusSecondary: "Reviewing",
				AdminScore:             newInt(60),
				FundApprovedAmount:     newInt64(-20),
			},
			store:          &mock.MockProjectStore{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &projects.FundApprovedAmountNegativeError{},
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
			if tt.expectedError != nil {
				errBody := getErrorResponse(t, res)
				assertErrorMessage(t, errBody.Message, tt.expectedError.Error())
			}
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
