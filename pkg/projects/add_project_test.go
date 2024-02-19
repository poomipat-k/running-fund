package projects_test

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/poomipat-k/running-fund/pkg/mock"
	"github.com/poomipat-k/running-fund/pkg/projects"
	s3Service "github.com/poomipat-k/running-fund/pkg/upload"
)

type ErrorBody struct {
	Error   bool
	Message string
}

func TestAddProject(t *testing.T) {

	tests := []struct {
		name           string
		payload        projects.AddProjectRequest
		store          *mock.MockProjectStore
		expectedStatus int
		expectedError  error
	}{
		{
			name:    "should error collaborated required",
			payload: projects.AddProjectRequest{},
			store: &mock.MockProjectStore{
				AddProjectFunc: addProjectSuccess,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &projects.CollaboratedRequiredError{},
		},
		{
			name:    "should error when has collaborated and collaborateFiles is empty",
			payload: projects.AddProjectRequest{Collaborated: newTrue()},
			store: &mock.MockProjectStore{
				AddProjectFunc: addProjectSuccess,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &projects.CollaboratedFilesRequiredError{},
		},
		{
			name:    "should error when has general.projectName is empty",
			payload: projects.AddProjectRequest{Collaborated: newFalse()},
			store: &mock.MockProjectStore{
				AddProjectFunc: addProjectSuccess,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &projects.ProjectNameRequiredError{},
		},
		{
			name: "should error when has general.eventDate.year is empty",
			payload: projects.AddProjectRequest{Collaborated: newFalse(),
				General: projects.AddProjectGeneralDetails{
					ProjectName: "A",
				}},
			store: &mock.MockProjectStore{
				AddProjectFunc: addProjectSuccess,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &projects.YearRequiredError{},
		},
	}

	for _, tt := range tests {
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
				// close it when it has done its job
				defer multipartWriter.Close()

				// create a form field writer for name
				formStr, err := multipartWriter.CreateFormField("form")
				if err != nil {
					t.Error(err)
				}

				// write string to the form field writer for form
				formStr.Write(body)
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

func newFalse() *bool {
	b := false
	return &b
}

func newTrue() *bool {
	b := true
	return &b
}

func addProjectSuccess(userId int, collaborateFiles []*multipart.FileHeader, otherFiles []projects.DetailsFiles) (string, error) {
	return "", nil
}
