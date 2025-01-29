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
				AdminScore:             newFloat64(-10),
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
				AdminScore:             newFloat64(110),
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
				AdminScore:             newFloat64(60),
				FundApprovedAmount:     newInt64(-20),
			},
			store:          &mock.MockProjectStore{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &projects.FundApprovedAmountNegativeError{},
		},
		{
			name: "should error when adminComment is longer than 512 characters",
			payload: projects.AdminUpdateProjectRequest{
				ProjectStatusPrimary:   "CurrentBeforeApprove",
				ProjectStatusSecondary: "Reviewing",
				AdminScore:             newFloat64(60),
				FundApprovedAmount:     newInt64(200000),
				AdminComment:           newString("012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012"),
			},
			store:          &mock.MockProjectStore{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &projects.AdminCommentTooLongError{Length: 513},
		},
		{
			name: "should error when adminComment is longer than 512 characters (Thai)",
			payload: projects.AdminUpdateProjectRequest{
				ProjectStatusPrimary:   "CurrentBeforeApprove",
				ProjectStatusSecondary: "Reviewing",
				AdminScore:             newFloat64(60),
				FundApprovedAmount:     newInt64(200000),
				AdminComment:           newString("การเดินทางตลอดหนึ่งปีที่ผ่านมา เราต้องเจอกับเรื่องราวมากมาย เผชิญหน้ากับเหตุการณ์ไม่คาดคิด และรับมือกับหลายความรู้สึกที่เกาะกุมอยู่ในใจ ด้วยเหตุนี้ ยิ่งใกล้ช่วงท้ายปี หลายคนเลยอยากปล่อยให้ ‘ปีเก่า’ เป็นเรื่องราวของ ‘ปีเก่า’ พร้อมทิ้งเรื่องราวเดิมๆ ไว้ข้างหลังและมุ่งหน้าสู่การเดินทางใหม่ที่กำลังจะมาถึงการเดินทางตลอดหนึ่งปีที่ผ่านมา เราต้องเจอกับเรื่องราวมากมาย เผชิญหน้ากับเหตุการณ์ไม่คาดคิด และรับมือกับหลายความรู้สึกที่เกาะกุมอยู่ในใจ ด้วยเหตุนี้ ยิ่งใกล้ช่วงท้ายปี หลายคนเลยอยากปล่อยให้ ‘ปีเก่า’ เป็นเรื่องราวของ ‘ปีเก่า’ พร้อมทิ้งเรื่องราวเดิมๆ ไว้ข้างหลังและมุ่งหน้าสู่การเดินทางใหม่ที่กำลังจะมาถึง"),
			},
			store:          &mock.MockProjectStore{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  &projects.AdminCommentTooLongError{Length: 604},
		},

		// Success cases

		// Primary and Secondary not changed
		{
			name: "should save data correctly when ProjectStatusPrimary and ProjectStatusSecondary haven't changed v1",
			payload: projects.AdminUpdateProjectRequest{
				ProjectStatusPrimary:   "CurrentBeforeApprove",
				ProjectStatusSecondary: "Reviewing",
				AdminScore:             newFloat64(70),
				FundApprovedAmount:     newInt64(200000),
				AdminComment:           newString("Admin comment 1"),
			},
			store: &mock.MockProjectStore{
				GetProjectStatusByProjectCodeFunc: func(projectCode string) (projects.AdminUpdateParam, error) {
					return projects.AdminUpdateParam{
						ProjectHistoryId: 1,
						ProjectStatus:    "Reviewing",
						AdminApprovedAt:  newTime(time.Date(2024, 1, 20, 15, 30, 45, 9, time.UTC)),
					}, nil
				},
			},
			expectedStatus: http.StatusCreated,
			expectedUpdatedData: projects.AdminUpdateParam{
				ProjectHistoryId:   1,
				ProjectStatus:      "Reviewing",
				AdminScore:         newFloat64(70),
				FundApprovedAmount: newInt64(200000),
				AdminComment:       newString("Admin comment 1"),
				AdminApprovedAt:    newTime(time.Date(2024, 1, 20, 15, 30, 45, 9, time.UTC)),
			},
		},
		{
			name: "should save data correctly when ProjectStatusPrimary and ProjectStatusSecondary haven't changed v2",
			payload: projects.AdminUpdateProjectRequest{
				ProjectStatusPrimary:   "Approved",
				ProjectStatusSecondary: "Approved",
				AdminScore:             newFloat64(70),
				FundApprovedAmount:     newInt64(200000),
				AdminComment:           newString("Admin comment 1"),
			},
			store: &mock.MockProjectStore{
				GetProjectStatusByProjectCodeFunc: func(projectCode string) (projects.AdminUpdateParam, error) {
					return projects.AdminUpdateParam{
						ProjectHistoryId: 1,
						ProjectStatus:    "Approved",
						AdminApprovedAt:  newTime(time.Date(2024, 1, 20, 15, 30, 45, 9, time.UTC)),
					}, nil
				},
			},
			expectedStatus: http.StatusCreated,
			expectedUpdatedData: projects.AdminUpdateParam{
				ProjectHistoryId:   1,
				ProjectStatus:      "Approved",
				AdminScore:         newFloat64(70),
				FundApprovedAmount: newInt64(200000),
				AdminComment:       newString("Admin comment 1"),
				AdminApprovedAt:    newTime(time.Date(2024, 1, 20, 15, 30, 45, 9, time.UTC)),
			},
		},
		{
			name: "should save data correctly when ProjectStatusPrimary and ProjectStatusSecondary haven't changed v3",
			payload: projects.AdminUpdateProjectRequest{
				ProjectStatusPrimary:   "NotApproved",
				ProjectStatusSecondary: "NotApproved",
				AdminScore:             newFloat64(70),
				FundApprovedAmount:     newInt64(200000),
				AdminComment:           newString("Admin comment 1"),
			},
			store: &mock.MockProjectStore{
				GetProjectStatusByProjectCodeFunc: func(projectCode string) (projects.AdminUpdateParam, error) {
					return projects.AdminUpdateParam{
						ProjectHistoryId: 1,
						ProjectStatus:    "NotApproved",
						AdminApprovedAt:  newTime(time.Date(2024, 1, 20, 15, 30, 45, 9, time.UTC)),
					}, nil
				},
			},
			expectedStatus: http.StatusCreated,
			expectedUpdatedData: projects.AdminUpdateParam{
				ProjectHistoryId:   1,
				ProjectStatus:      "NotApproved",
				AdminScore:         newFloat64(70),
				FundApprovedAmount: newInt64(200000),
				AdminComment:       newString("Admin comment 1"),
				AdminApprovedAt:    newTime(time.Date(2024, 1, 20, 15, 30, 45, 9, time.UTC)),
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
			expectedStatus: http.StatusCreated,
			expectedUpdatedData: projects.AdminUpdateParam{
				ProjectHistoryId:   1,
				ProjectStatus:      "Reviewing",
				AdminScore:         nil,
				FundApprovedAmount: nil,
				AdminComment:       newString("Admin comment 1"),
				AdminApprovedAt:    nil,
			},
		},
		// Either Primary and Secondary status changed
		{
			name: "should save data correctly when ProjectStatusPrimary changed to Approved and ProjectStatusSecondary haven't changed",
			payload: projects.AdminUpdateProjectRequest{
				ProjectStatusPrimary:   "Approved",
				ProjectStatusSecondary: "Reviewed",
				AdminScore:             newFloat64(70),
				FundApprovedAmount:     newInt64(200000),
				AdminComment:           newString("Admin comment 2"),
				AdminApprovedAt:        newTime(time.Date(2024, 1, 20, 15, 30, 45, 9, time.UTC)),
			},
			store: &mock.MockProjectStore{
				GetProjectStatusByProjectCodeFunc: func(projectCode string) (projects.AdminUpdateParam, error) {
					return projects.AdminUpdateParam{
						ProjectHistoryId: 1,
						ProjectStatus:    "Reviewed",
					}, nil
				},
			},
			expectedStatus: http.StatusCreated,
			expectedUpdatedData: projects.AdminUpdateParam{
				ProjectHistoryId:   1,
				ProjectStatus:      "Approved",
				AdminScore:         newFloat64(70),
				FundApprovedAmount: newInt64(200000),
				AdminComment:       newString("Admin comment 2"),
				AdminApprovedAt:    newTime(time.Date(2024, 1, 20, 15, 30, 45, 9, time.UTC)),
			},
		},
		{
			name: "should save data correctly when ProjectStatusPrimary changed to Approved and ProjectStatusSecondary changed to Reviewing",
			payload: projects.AdminUpdateProjectRequest{
				ProjectStatusPrimary:   "Approved",
				ProjectStatusSecondary: "Reviewing",
				AdminScore:             newFloat64(70),
				FundApprovedAmount:     newInt64(200000),
				AdminComment:           newString("Admin comment 2"),
				AdminApprovedAt:        newTime(time.Date(2024, 1, 20, 15, 30, 45, 9, time.UTC)),
			},
			store: &mock.MockProjectStore{
				GetProjectStatusByProjectCodeFunc: func(projectCode string) (projects.AdminUpdateParam, error) {
					return projects.AdminUpdateParam{
						ProjectHistoryId: 1,
						ProjectStatus:    "Reviewed",
					}, nil
				},
			},
			expectedStatus: http.StatusCreated,
			expectedUpdatedData: projects.AdminUpdateParam{
				ProjectHistoryId:   1,
				ProjectStatus:      "Approved",
				AdminScore:         newFloat64(70),
				FundApprovedAmount: newInt64(200000),
				AdminComment:       newString("Admin comment 2"),
				AdminApprovedAt:    newTime(time.Date(2024, 1, 20, 15, 30, 45, 9, time.UTC)),
			},
		},
		{
			name: "should save data correctly when ProjectStatusPrimary haven't changed and ProjectStatusSecondary changed to Approved",
			payload: projects.AdminUpdateProjectRequest{
				ProjectStatusPrimary:   "CurrentBeforeApprove",
				ProjectStatusSecondary: "Approved",
				AdminScore:             newFloat64(70),
				FundApprovedAmount:     newInt64(200000),
				AdminComment:           newString("Admin comment 2"),
				AdminApprovedAt:        newTime(time.Date(2024, 1, 20, 15, 30, 45, 9, time.UTC)),
			},
			store: &mock.MockProjectStore{
				GetProjectStatusByProjectCodeFunc: func(projectCode string) (projects.AdminUpdateParam, error) {
					return projects.AdminUpdateParam{
						ProjectHistoryId: 1,
						ProjectStatus:    "Reviewing",
					}, nil
				},
			},
			expectedStatus: http.StatusCreated,
			expectedUpdatedData: projects.AdminUpdateParam{
				ProjectHistoryId:   1,
				ProjectStatus:      "Approved",
				AdminScore:         newFloat64(70),
				FundApprovedAmount: newInt64(200000),
				AdminComment:       newString("Admin comment 2"),
				AdminApprovedAt:    newTime(time.Date(2024, 1, 20, 15, 30, 45, 9, time.UTC)),
			},
		},
		{
			name: "should save data correctly when ProjectStatusPrimary haven't changed (NotApproved) and ProjectStatusSecondary changed to Approved",
			payload: projects.AdminUpdateProjectRequest{
				ProjectStatusPrimary:   "NotApproved",
				ProjectStatusSecondary: "Approved",
				AdminScore:             newFloat64(70),
				FundApprovedAmount:     newInt64(200000),
				AdminComment:           newString("Admin comment 2"),
				AdminApprovedAt:        newTime(time.Date(2024, 1, 20, 15, 30, 45, 9, time.UTC)),
			},
			store: &mock.MockProjectStore{
				GetProjectStatusByProjectCodeFunc: func(projectCode string) (projects.AdminUpdateParam, error) {
					return projects.AdminUpdateParam{
						ProjectHistoryId: 1,
						ProjectStatus:    "NotApproved",
					}, nil
				},
			},
			expectedStatus: http.StatusCreated,
			expectedUpdatedData: projects.AdminUpdateParam{
				ProjectHistoryId:   1,
				ProjectStatus:      "Approved",
				AdminScore:         newFloat64(70),
				FundApprovedAmount: newInt64(200000),
				AdminComment:       newString("Admin comment 2"),
				AdminApprovedAt:    newTime(time.Date(2024, 1, 20, 15, 30, 45, 9, time.UTC)),
			},
		},
		{
			name: "should save data correctly when ProjectStatusPrimary changed to NotApproved and ProjectStatusSecondary haven't changed",
			payload: projects.AdminUpdateProjectRequest{
				ProjectStatusPrimary:   "NotApproved",
				ProjectStatusSecondary: "Reviewed",
				AdminScore:             newFloat64(70),
				FundApprovedAmount:     newInt64(200000),
				AdminComment:           newString("Admin comment 2"),
				AdminApprovedAt:        newTime(time.Date(2024, 1, 20, 15, 30, 45, 9, time.UTC)),
			},
			store: &mock.MockProjectStore{
				GetProjectStatusByProjectCodeFunc: func(projectCode string) (projects.AdminUpdateParam, error) {
					return projects.AdminUpdateParam{
						ProjectHistoryId: 1,
						ProjectStatus:    "Reviewed",
						AdminApprovedAt:  newTime(time.Date(2023, 1, 20, 15, 30, 45, 9, time.UTC)),
					}, nil
				},
			},
			expectedStatus: http.StatusCreated,
			expectedUpdatedData: projects.AdminUpdateParam{
				ProjectHistoryId:   1,
				ProjectStatus:      "NotApproved",
				AdminScore:         newFloat64(70),
				FundApprovedAmount: newInt64(200000),
				AdminComment:       newString("Admin comment 2"),
				AdminApprovedAt:    nil,
			},
		},
		{
			name: "should save data correctly when ProjectStatusPrimary changed to NotApproved and ProjectStatusSecondary change to Approved",
			payload: projects.AdminUpdateProjectRequest{
				ProjectStatusPrimary:   "NotApproved",
				ProjectStatusSecondary: "Approved",
				AdminScore:             newFloat64(70),
				FundApprovedAmount:     newInt64(200000),
				AdminComment:           newString("Admin comment 2"),
				AdminApprovedAt:        nil,
			},
			store: &mock.MockProjectStore{
				GetProjectStatusByProjectCodeFunc: func(projectCode string) (projects.AdminUpdateParam, error) {
					return projects.AdminUpdateParam{
						ProjectHistoryId: 1,
						ProjectStatus:    "Reviewed",
						AdminApprovedAt:  newTime(time.Date(2023, 1, 20, 15, 30, 45, 9, time.UTC)),
					}, nil
				},
			},
			expectedStatus: http.StatusCreated,
			expectedUpdatedData: projects.AdminUpdateParam{
				ProjectHistoryId:   1,
				ProjectStatus:      "NotApproved",
				AdminScore:         newFloat64(70),
				FundApprovedAmount: newInt64(200000),
				AdminComment:       newString("Admin comment 2"),
				AdminApprovedAt:    nil,
			},
		},
		{
			name: "should save data correctly when ProjectStatusPrimary changed to NotApproved and currentProject is Approved",
			payload: projects.AdminUpdateProjectRequest{
				ProjectStatusPrimary:   "NotApproved",
				ProjectStatusSecondary: "Approved",
				AdminScore:             newFloat64(70),
				FundApprovedAmount:     newInt64(200000),
				AdminComment:           newString("Admin comment 2"),
				AdminApprovedAt:        nil,
			},
			store: &mock.MockProjectStore{
				GetProjectStatusByProjectCodeFunc: func(projectCode string) (projects.AdminUpdateParam, error) {
					return projects.AdminUpdateParam{
						ProjectHistoryId: 1,
						ProjectStatus:    "Approved",
						AdminApprovedAt:  newTime(time.Date(2023, 1, 20, 15, 30, 45, 9, time.UTC)),
					}, nil
				},
			},
			expectedStatus: http.StatusCreated,
			expectedUpdatedData: projects.AdminUpdateParam{
				ProjectHistoryId:   1,
				ProjectStatus:      "NotApproved",
				AdminScore:         newFloat64(70),
				FundApprovedAmount: newInt64(200000),
				AdminComment:       newString("Admin comment 2"),
				AdminApprovedAt:    nil,
			},
		},
		{
			name: "should save data correctly when ProjectStatusPrimary haven't changed and ProjectStatusSecondary change to NotApproved",
			payload: projects.AdminUpdateProjectRequest{
				ProjectStatusPrimary:   "Approved",
				ProjectStatusSecondary: "NotApproved",
				AdminScore:             newFloat64(70),
				FundApprovedAmount:     newInt64(200000),
				AdminComment:           newString("Admin comment 2"),
				AdminApprovedAt:        nil,
			},
			store: &mock.MockProjectStore{
				GetProjectStatusByProjectCodeFunc: func(projectCode string) (projects.AdminUpdateParam, error) {
					return projects.AdminUpdateParam{
						ProjectHistoryId: 1,
						ProjectStatus:    "Approved",
						AdminApprovedAt:  newTime(time.Date(2023, 1, 20, 15, 30, 45, 9, time.UTC)),
					}, nil
				},
			},
			expectedStatus: http.StatusCreated,
			expectedUpdatedData: projects.AdminUpdateParam{
				ProjectHistoryId:   1,
				ProjectStatus:      "NotApproved",
				AdminScore:         newFloat64(70),
				FundApprovedAmount: newInt64(200000),
				AdminComment:       newString("Admin comment 2"),
				AdminApprovedAt:    nil,
			},
		},
		// Primary or Secondary Changed and newStatus is not Approved/NotApproved
		{
			name: "should save data correctly when ProjectStatusPrimary haven't changed and ProjectStatusSecondary change to Revise",
			payload: projects.AdminUpdateProjectRequest{
				ProjectStatusPrimary:   "CurrentBeforeApprove",
				ProjectStatusSecondary: "Revise",
				AdminScore:             nil,
				FundApprovedAmount:     nil,
				AdminComment:           newString("Revise something"),
				AdminApprovedAt:        nil,
			},
			store: &mock.MockProjectStore{
				GetProjectStatusByProjectCodeFunc: func(projectCode string) (projects.AdminUpdateParam, error) {
					return projects.AdminUpdateParam{
						ProjectHistoryId: 1,
						ProjectStatus:    "Reviewing",
						AdminApprovedAt:  nil,
					}, nil
				},
			},
			expectedStatus: http.StatusCreated,
			expectedUpdatedData: projects.AdminUpdateParam{
				ProjectHistoryId:   1,
				ProjectStatus:      "Revise",
				AdminScore:         nil,
				FundApprovedAmount: nil,
				AdminComment:       newString("Revise something"),
				AdminApprovedAt:    nil,
			},
		},
		{
			name: "should save data correctly when ProjectStatusPrimary haven't changed and ProjectStatusSecondary change to Reviewed",
			payload: projects.AdminUpdateProjectRequest{
				ProjectStatusPrimary:   "CurrentBeforeApprove",
				ProjectStatusSecondary: "Reviewed",
				AdminScore:             nil,
				FundApprovedAmount:     nil,
				AdminComment:           newString("Revise something"),
				AdminApprovedAt:        nil,
			},
			store: &mock.MockProjectStore{
				GetProjectStatusByProjectCodeFunc: func(projectCode string) (projects.AdminUpdateParam, error) {
					return projects.AdminUpdateParam{
						ProjectHistoryId: 1,
						ProjectStatus:    "Revise",
						AdminApprovedAt:  nil,
					}, nil
				},
			},
			expectedStatus: http.StatusCreated,
			expectedUpdatedData: projects.AdminUpdateParam{
				ProjectHistoryId:   1,
				ProjectStatus:      "Reviewed",
				AdminScore:         nil,
				FundApprovedAmount: nil,
				AdminComment:       newString("Revise something"),
				AdminApprovedAt:    nil,
			},
		},
		{
			name: "should save data correctly when ProjectStatusPrimary haven't changed and ProjectStatusSecondary change to Start",
			payload: projects.AdminUpdateProjectRequest{
				ProjectStatusPrimary:   "Approved",
				ProjectStatusSecondary: "Start",
				AdminScore:             newFloat64(60),
				FundApprovedAmount:     newInt64(202020),
				AdminComment:           nil,
				AdminApprovedAt:        nil,
			},
			store: &mock.MockProjectStore{
				GetProjectStatusByProjectCodeFunc: func(projectCode string) (projects.AdminUpdateParam, error) {
					return projects.AdminUpdateParam{
						ProjectHistoryId:   1,
						ProjectStatus:      "Approved",
						AdminScore:         newFloat64(55),
						FundApprovedAmount: newInt64(200),
						AdminApprovedAt:    nil,
					}, nil
				},
			},
			expectedStatus: http.StatusCreated,
			expectedUpdatedData: projects.AdminUpdateParam{
				ProjectHistoryId:   1,
				ProjectStatus:      "Start",
				AdminScore:         newFloat64(60),
				FundApprovedAmount: newInt64(202020),
				AdminComment:       nil,
				AdminApprovedAt:    nil,
			},
		},
		{
			name: "should save data correctly when ProjectStatusPrimary haven't changed and ProjectStatusSecondary change to Completed",
			payload: projects.AdminUpdateProjectRequest{
				ProjectStatusPrimary:   "Approved",
				ProjectStatusSecondary: "Completed",
				AdminScore:             newFloat64(60),
				FundApprovedAmount:     newInt64(202020),
				AdminComment:           nil,
				AdminApprovedAt:        nil,
			},
			store: &mock.MockProjectStore{
				GetProjectStatusByProjectCodeFunc: func(projectCode string) (projects.AdminUpdateParam, error) {
					return projects.AdminUpdateParam{
						ProjectHistoryId:   1,
						ProjectStatus:      "Start",
						AdminScore:         newFloat64(55),
						FundApprovedAmount: newInt64(200),
						AdminApprovedAt:    nil,
					}, nil
				},
			},
			expectedStatus: http.StatusCreated,
			expectedUpdatedData: projects.AdminUpdateParam{
				ProjectHistoryId:   1,
				ProjectStatus:      "Completed",
				AdminScore:         newFloat64(60),
				FundApprovedAmount: newInt64(202020),
				AdminComment:       nil,
				AdminApprovedAt:    nil,
			},
		},
		{
			name: "should save data correctly when ProjectStatusPrimary haven't changed and ProjectStatusSecondary change to Completed v2",
			payload: projects.AdminUpdateProjectRequest{
				ProjectStatusPrimary:   "Approved",
				ProjectStatusSecondary: "Completed",
				AdminScore:             newFloat64(60),
				FundApprovedAmount:     newInt64(202020),
				AdminComment:           nil,
				AdminApprovedAt:        nil,
			},
			store: &mock.MockProjectStore{
				GetProjectStatusByProjectCodeFunc: func(projectCode string) (projects.AdminUpdateParam, error) {
					return projects.AdminUpdateParam{
						ProjectHistoryId:   1,
						ProjectStatus:      "Approved",
						AdminScore:         newFloat64(55),
						FundApprovedAmount: newInt64(200),
						AdminApprovedAt:    nil,
					}, nil
				},
			},
			expectedStatus: http.StatusCreated,
			expectedUpdatedData: projects.AdminUpdateParam{
				ProjectHistoryId:   1,
				ProjectStatus:      "Completed",
				AdminScore:         newFloat64(60),
				FundApprovedAmount: newInt64(202020),
				AdminComment:       nil,
				AdminApprovedAt:    nil,
			},
		},
		{
			name: "should save data correctly when ProjectStatusPrimary haven't changed and ProjectStatusSecondary change to Completed v3",
			payload: projects.AdminUpdateProjectRequest{
				ProjectStatusPrimary:   "NotApproved",
				ProjectStatusSecondary: "Completed",
				AdminScore:             newFloat64(60),
				FundApprovedAmount:     newInt64(202020),
				AdminComment:           nil,
				AdminApprovedAt:        nil,
			},
			store: &mock.MockProjectStore{
				GetProjectStatusByProjectCodeFunc: func(projectCode string) (projects.AdminUpdateParam, error) {
					return projects.AdminUpdateParam{
						ProjectHistoryId:   1,
						ProjectStatus:      "NotApproved",
						AdminScore:         newFloat64(55),
						FundApprovedAmount: newInt64(200),
						AdminApprovedAt:    nil,
					}, nil
				},
			},
			expectedStatus: http.StatusCreated,
			expectedUpdatedData: projects.AdminUpdateParam{
				ProjectHistoryId:   1,
				ProjectStatus:      "Completed",
				AdminScore:         newFloat64(60),
				FundApprovedAmount: newInt64(202020),
				AdminComment:       nil,
				AdminApprovedAt:    nil,
			},
		},
		{
			name: "should save data correctly when ProjectStatusPrimary haven't changed and ProjectStatusSecondary change to Completed v3",
			payload: projects.AdminUpdateProjectRequest{
				ProjectStatusPrimary:   "CurrentBeforeApprove",
				ProjectStatusSecondary: "Start",
				AdminScore:             newFloat64(60),
				FundApprovedAmount:     newInt64(202020),
				AdminComment:           nil,
				AdminApprovedAt:        nil,
			},
			store: &mock.MockProjectStore{
				GetProjectStatusByProjectCodeFunc: func(projectCode string) (projects.AdminUpdateParam, error) {
					return projects.AdminUpdateParam{
						ProjectHistoryId:   1,
						ProjectStatus:      "Reviewing",
						AdminScore:         newFloat64(55),
						FundApprovedAmount: newInt64(200),
						AdminApprovedAt:    nil,
					}, nil
				},
			},
			expectedStatus: http.StatusCreated,
			expectedUpdatedData: projects.AdminUpdateParam{
				ProjectHistoryId:   1,
				ProjectStatus:      "Start",
				AdminScore:         newFloat64(60),
				FundApprovedAmount: newInt64(202020),
				AdminComment:       nil,
				AdminApprovedAt:    nil,
			},
		},
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
