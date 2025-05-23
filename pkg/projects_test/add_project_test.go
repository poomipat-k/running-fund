package projects_test

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/poomipat-k/running-fund/pkg/mock"
	"github.com/poomipat-k/running-fund/pkg/projects"
	s3Service "github.com/poomipat-k/running-fund/pkg/s3-service"
)

type ErrorBody struct {
	Error   bool
	Message string
}

type TestCase struct {
	name                   string
	payload                projects.AddProjectRequest
	store                  *mock.MockProjectStore
	expectedStatus         int
	expectedError          error
	collaborationFilesPath string
	marketingFilesPath     string
	routeFilesPath         string
	eventMapFilesPath      string
	eventDetailsFilesPath  string
}

func TestAddProject(t *testing.T) {

	pagesCases := [][]TestCase{
		GeneralAndCollaboratedTestCases,
		ContactTestCases,
		Details,
		Experience,
		Fund,
		Attachment,
	}
	t.Setenv("APPLICANT_CRITERIA_VERSION", "1")
	for _, cases := range pagesCases {
		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				store := tt.store
				userStore := &mock.MockUserStore{}
				handler := projects.NewProjectHandler(store, userStore, s3Service.S3Service{})

				// multipart/form-data set up

				// set up a pipe avoid buffering
				pipeReader, pipeWriter := io.Pipe()

				// this writer is going to transform what we pass to it to multipart form data
				// and write it to our io.Pipe
				multipartWriter := multipart.NewWriter(pipeWriter)

				body, err := json.Marshal(tt.payload)
				if err != nil {
					t.Error("error marshal payload err:", err)
				}

				go func() {
					// close it when it done its job
					defer multipartWriter.Close()

					// create a form field writer for name
					formStr, err := multipartWriter.CreateFormField("form")
					if err != nil {
						t.Error(err)
					}

					// write string to the form field writer for form
					formStr.Write(body)

					filePath := "test.png"
					file, err := os.Open(filePath)
					if err != nil {
						t.Error(err)
					}
					defer file.Close()

					if tt.collaborationFilesPath != "" {
						mk, err := multipartWriter.CreateFormFile("collaborationFiles", filePath)
						if err != nil {
							t.Error(err)
						}
						if _, err := io.Copy(mk, file); err != nil {
							t.Error(err)
						}
					}

					if tt.marketingFilesPath != "" {
						mk, err := multipartWriter.CreateFormFile("marketingFiles", filePath)
						if err != nil {
							t.Error(err)
						}
						if _, err := io.Copy(mk, file); err != nil {
							t.Error(err)
						}
					}

					if tt.routeFilesPath != "" {
						mk, err := multipartWriter.CreateFormFile("routeFiles", filePath)
						if err != nil {
							t.Error(err)
						}
						if _, err := io.Copy(mk, file); err != nil {
							t.Error(err)
						}
					}

					if tt.eventMapFilesPath != "" {
						mk, err := multipartWriter.CreateFormFile("eventMapFiles", filePath)
						if err != nil {
							t.Error(err)
						}
						if _, err := io.Copy(mk, file); err != nil {
							t.Error(err)
						}
					}

					if tt.eventDetailsFilesPath != "" {
						mk, err := multipartWriter.CreateFormFile("eventDetailsFiles", filePath)
						if err != nil {
							t.Error(err)
						}
						if _, err := io.Copy(mk, file); err != nil {
							t.Error(err)
						}
					}
				}()

				// End multipart/form-data setup

				res := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodPost, "/project", pipeReader)

				req.Header.Set("userId", "1")
				// Set content-type to multipart
				req.Header.Add("content-type", multipartWriter.FormDataContentType())

				handler.AddProject(res, req)

				assertStatus(t, res.Code, tt.expectedStatus)
				if tt.expectedError != nil {
					errBody := getErrorResponse(t, res)
					assertErrorMessage(t, errBody.Message, tt.expectedError.Error())
				}
			})
		}
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
		t.Errorf("Error unmarshal ErrorResponse err:%v", err)
	}
	return body
}

func assertUpdatedData(t testing.TB, got, want projects.AdminUpdateParam) {
	t.Helper()
	if got.ProjectHistoryId != want.ProjectHistoryId {
		t.Errorf("ProjectHistoryId: got %d, want %d", got.ProjectHistoryId, want.ProjectHistoryId)
	}
	if got.ProjectStatus != want.ProjectStatus {
		t.Errorf("ProjectStatus: got %s, want %s", got.ProjectStatus, want.ProjectStatus)
	}
	// AdminScore
	if got.AdminScore == nil && want.AdminScore != nil {
		t.Error("AdminScore should not be nil")
	}
	if got.AdminScore != nil && want.AdminScore == nil {
		t.Error("AdminScore should be nil")
	}
	if got.AdminScore != nil && want.AdminScore != nil && *got.AdminScore != *want.AdminScore {
		t.Errorf("AdminScore: got %f, want %f", *got.AdminScore, *want.AdminScore)
	}
	// FundApprovedAmount
	if got.FundApprovedAmount == nil && want.FundApprovedAmount != nil {
		t.Error("FundApprovedAmount should not be nil")
	}
	if got.FundApprovedAmount != nil && want.FundApprovedAmount == nil {
		t.Error("FundApprovedAmount should be nil")
	}
	if got.FundApprovedAmount != nil && want.FundApprovedAmount != nil && *got.FundApprovedAmount != *want.FundApprovedAmount {
		t.Errorf("FundApprovedAmount: got %d, want %d", *got.FundApprovedAmount, *want.FundApprovedAmount)
	}
	// AdminComment
	if got.AdminComment == nil && want.AdminComment != nil {
		t.Error("AdminComment should not be nil")
	}
	if got.AdminComment != nil && want.AdminComment == nil {
		t.Error("AdminComment should be nil")
	}
	if got.AdminComment != nil && want.AdminComment != nil && *got.AdminComment != *want.AdminComment {
		t.Errorf("AdminComment: got %s, want %s", *got.AdminComment, *want.AdminComment)
	}
	// AdminApprovedAt
	if got.AdminApprovedAt == nil && want.AdminApprovedAt != nil {
		t.Error("AdminApprovedAt should not be nil")
	}
	if got.AdminApprovedAt != nil && want.AdminApprovedAt == nil {
		t.Error("AdminApprovedAt should be nil")
	}
	if got.AdminApprovedAt != nil && want.AdminApprovedAt != nil && *got.AdminApprovedAt != *want.AdminApprovedAt {
		t.Errorf("AdminApprovedAt: got %v, want %v", *got.AdminApprovedAt, *want.AdminApprovedAt)
	}
}

func newInt(val int) *int {
	v := val
	return &v
}

func newInt64(val int64) *int64 {
	v := val
	return &v
}

func newString(val string) *string {
	v := val
	return &v
}

func newTime(t time.Time) *time.Time {
	v := t
	return &v
}

func newFloat64(val float64) *float64 {
	v := val
	return &v
}

func newFalse() *bool {
	b := false
	return &b
}

func newTrue() *bool {
	b := true
	return &b
}

func addProjectSuccess(addProject projects.AddProjectRequest, userId int, criteria []projects.ApplicantSelfScoreCriteria, attachments []projects.Attachments) (int, error) {
	return 1, nil
}

func getApplicantCriteriaSuccess(version int) ([]projects.ApplicantSelfScoreCriteria, error) {
	return []projects.ApplicantSelfScoreCriteria{
		{
			CriteriaVersion: 1,
			OrderNumber:     1,
			Display:         "A",
		},
		{
			CriteriaVersion: 1,
			OrderNumber:     2,
			Display:         "B",
		},
	}, nil
}

func getApplicantCriteriaNotFound(version int) ([]projects.ApplicantSelfScoreCriteria, error) {
	return []projects.ApplicantSelfScoreCriteria{}, nil
}
