package projects_test

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/poomipat-k/running-fund/pkg/mock"
	"github.com/poomipat-k/running-fund/pkg/projects"
	s3Service "github.com/poomipat-k/running-fund/pkg/s3-service"
)

type UpdateProjectTestCase struct {
	name                string
	payload             projects.AdminUpdateProjectRequest
	store               *mock.MockProjectStore
	expectedStatus      int
	expectedError       error
	additionFilesPath   string
	expectedUpdatedData projects.AdminUpdateParam
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
			name: "should error when projectStatusPrimary is invalid",
			payload: projects.AdminUpdateProjectRequest{
				ProjectStatusPrimary: "Something",
			},
			store:          &mock.MockProjectStore{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &projects.ProjectStatusPrimaryInvalidError{},
		},
		{
			name: "should error when projectStatusSecondary is missing",
			payload: projects.AdminUpdateProjectRequest{
				ProjectStatusPrimary: "CurrentBeforeApprove",
			},
			store:          &mock.MockProjectStore{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &projects.ProjectStatusSecondaryRequiredError{},
		},
		{
			name: "should error when projectStatusSecondary is invalid",
			payload: projects.AdminUpdateProjectRequest{
				ProjectStatusPrimary:   "CurrentBeforeApprove",
				ProjectStatusSecondary: "XXX",
			},
			store:          &mock.MockProjectStore{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &projects.ProjectStatusSecondaryInvalidError{},
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
		{
			name: "should save data correctly when ProjectStatusPrimary and ProjectStatusSecondary haven't changed",
			payload: projects.AdminUpdateProjectRequest{
				ProjectStatusPrimary:   "CurrentBeforeApprove",
				ProjectStatusSecondary: "Reviewing",
				AdminScore:             newInt(70),
				FundApprovedAmount:     newInt64(200000),
				AdminComment:           newString("Admin comment 1"),
			},
			store: &mock.MockProjectStore{
				GetProjectStatusByProjectCodeFunc: func(projectCode string) (projects.AdminUpdateParam, error) {
					return projects.AdminUpdateParam{
						ProjectHistoryId: 1,
						ProjectStatus:    "Reviewing",
					}, nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedUpdatedData: projects.AdminUpdateParam{
				ProjectHistoryId:   1,
				ProjectStatus:      "Reviewing",
				AdminScore:         newInt(70),
				FundApprovedAmount: newInt64(200000),
				AdminComment:       newString("Admin comment 1"),
				AdminApprovedAt:    nil,
			},
		},
		{
			name: "should save data correctly when ProjectStatusPrimary and ProjectStatusSecondary haven't changed and some optional payload is nil",
			payload: projects.AdminUpdateProjectRequest{
				ProjectStatusPrimary:   "CurrentBeforeApprove",
				ProjectStatusSecondary: "Reviewing",
				AdminComment:           newString("Admin comment 1"),
			},
			store: &mock.MockProjectStore{
				GetProjectStatusByProjectCodeFunc: func(projectCode string) (projects.AdminUpdateParam, error) {
					return projects.AdminUpdateParam{
						ProjectHistoryId: 1,
						ProjectStatus:    "Reviewing",
					}, nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedUpdatedData: projects.AdminUpdateParam{
				ProjectHistoryId:   1,
				ProjectStatus:      "Reviewing",
				AdminScore:         nil,
				FundApprovedAmount: nil,
				AdminComment:       newString("Admin comment 1"),
				AdminApprovedAt:    nil,
			},
		},
		// // ok
		// {
		// 	name: "should be okay",
		// 	payload: projects.AdminUpdateProjectRequest{
		// 		ProjectStatusPrimary:   "Approved",
		// 		ProjectStatusSecondary: "Reviewed",
		// 		AdminScore:             newInt(66),
		// 		FundApprovedAmount:     newInt64(200000),
		// 		AdminComment:           newString("Admin comment 1"),
		// 	},
		// 	store: &mock.MockProjectStore{
		// 		GetProjectStatusByProjectCodeFunc: func(projectCode string) (projects.AdminUpdateParam, error) {
		// 			return projects.AdminUpdateParam{
		// 				ProjectStatus: "Reviewing",
		// 			}, nil
		// 		},
		// 	},
		// 	expectedStatus: http.StatusOK,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := tt.store
			userStore := &mock.MockUserStore{}
			handler := projects.NewProjectHandler(store, userStore, s3Service.S3Service{})

			// multipart/form-data set up
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

				if tt.additionFilesPath != "" {
					mk, err := multipartWriter.CreateFormFile("additionFiles", filePath)
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
			req := httptest.NewRequest(http.MethodPost, "/admin/project/APR67_0501", pipeReader)
			// Set content-type to multipart
			req.Header.Add("content-type", multipartWriter.FormDataContentType())

			handler.AdminUpdateProject(res, req)
			assertStatus(t, res.Code, tt.expectedStatus)
			if tt.expectedError != nil {
				errBody := getErrorResponse(t, res)
				assertErrorMessage(t, errBody.Message, tt.expectedError.Error())
			}

			if tt.expectedUpdatedData.ProjectHistoryId != 0 {
				assertUpdatedData(t, tt.store.AdminUpdateData, tt.expectedUpdatedData)
			}
		})
	}
}
