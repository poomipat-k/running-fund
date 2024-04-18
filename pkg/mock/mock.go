package mock

import (
	"mime/multipart"
	"time"

	"github.com/poomipat-k/running-fund/pkg/projects"
	"github.com/poomipat-k/running-fund/pkg/users"
)

// users
type MockUserStore struct {
	Users                    map[int]users.User
	UsersMapByEmail          map[string]users.User
	GetUserByEmailFunc       func(email string) (users.User, error)
	AddUserFunc              func(user users.User, toBeDeletedId int) (int, string, error)
	GetUserByIdFunc          func(id int) (users.User, error)
	ActivateUserFunc         func(activateCode string) (int64, error)
	ForgotPasswordActionFunc func(resetPasswordCode string, email string, resetPasswordLink string) (int64, error)
	ResetPasswordFunc        func(resetPasswordCode string, newPassword string) (int64, error)
	GetUserFullNameByIdFunc  func(userId int) (users.UserFullName, error)
}

func (m *MockUserStore) GetUserById(id int) (users.User, error) {
	return m.GetUserByIdFunc(id)
}

func (m *MockUserStore) GetUserByEmail(email string) (users.User, error) {
	return m.GetUserByEmailFunc(email)
}

func (m *MockUserStore) AddUser(user users.User, toBeDeletedId int) (int, string, error) {
	return m.AddUserFunc(user, toBeDeletedId)
}

func (m *MockUserStore) ActivateUser(activateCode string) (int64, error) {
	return m.ActivateUserFunc(activateCode)
}

func (m *MockUserStore) ForgotPasswordAction(resetPasswordCode string, email string, resetPasswordLink string) (int64, error) {
	return m.ForgotPasswordActionFunc(resetPasswordCode, email, resetPasswordLink)
}

func (m *MockUserStore) ResetPassword(resetPasswordCode string, newPassword string) (int64, error) {
	return m.ResetPasswordFunc(resetPasswordCode, newPassword)
}

func (m *MockUserStore) GetUserFullNameById(userId int) (users.UserFullName, error) {
	return m.GetUserFullNameByIdFunc(userId)
}

// projects
type MockProjectStore struct {
	AdminUpdateData projects.AdminUpdateParam

	GetReviewerDashboardFunc                func(userId int, from time.Time, to time.Time) ([]projects.ReviewDashboardRow, error)
	GetReviewPeriodFunc                     func() (projects.ReviewPeriod, error)
	GetReviewerProjectDetailsFunc           func(userId int, projectCode string) (projects.ProjectReviewDetailsResponse, error)
	GetProjectCriteriaFunc                  func(criteriaVersion int) ([]projects.ProjectReviewCriteria, error)
	AddProjectFunc                          func(addProject projects.AddProjectRequest, userId int, criteria []projects.ApplicantSelfScoreCriteria, attachments []projects.Attachments) (int, error)
	GetApplicantCriteriaFunc                func(version int) ([]projects.ApplicantSelfScoreCriteria, error)
	GetAllProjectDashboardByApplicantIdFunc func(applicantId int) ([]projects.ApplicantDashboardItem, error)
	GetApplicantProjectDetailsFunc          func(isAdmin bool, projectCode string, userId int) ([]projects.ApplicantDetailsData, error)
	HasPermissionToAddAdditionalFilesFunc   func(userId int, projectCode string) bool
	GetProjectStatusByProjectCodeFunc       func(projectCode string) (projects.AdminUpdateParam, error)
	GetAdminRequestDashboardFunc            func(fromDate, toDate time.Time, orderBy string, limit, offset int, projectCode, projectName, projectStatus *string) ([]projects.AdminRequestDashboardRow, error)
	GetAdminSummaryFunc                     func(fromDate, toDate time.Time) ([]projects.AdminSummaryData, error)
}

func (m *MockProjectStore) GetReviewerDashboard(userId int, from time.Time, to time.Time) ([]projects.ReviewDashboardRow, error) {
	return m.GetReviewerDashboardFunc(userId, from, to)
}

func (m *MockProjectStore) GetReviewPeriod() (projects.ReviewPeriod, error) {
	return m.GetReviewPeriodFunc()
}

func (m *MockProjectStore) GetReviewerProjectDetails(userId int, projectCode string) (projects.ProjectReviewDetailsResponse, error) {
	return m.GetReviewerProjectDetailsFunc(userId, projectCode)
}

func (m *MockProjectStore) GetProjectCriteria(criteriaVersion int) ([]projects.ProjectReviewCriteria, error) {
	return m.GetProjectCriteriaFunc(criteriaVersion)
}

func (m *MockProjectStore) GetApplicantCriteria(version int) ([]projects.ApplicantSelfScoreCriteria, error) {
	return m.GetApplicantCriteriaFunc(version)
}

func (m *MockProjectStore) AddProject(addProject projects.AddProjectRequest, userId int, criteria []projects.ApplicantSelfScoreCriteria, attachments []projects.Attachments) (int, error) {
	return m.AddProjectFunc(addProject, userId, criteria, attachments)
}

func (m *MockProjectStore) GetAllProjectDashboardByApplicantId(applicantId int) ([]projects.ApplicantDashboardItem, error) {
	return m.GetAllProjectDashboardByApplicantIdFunc(applicantId)
}

func (m *MockProjectStore) GetApplicantProjectDetails(isAdmin bool, projectCode string, userId int) ([]projects.ApplicantDetailsData, error) {
	return m.GetApplicantProjectDetailsFunc(isAdmin, projectCode, userId)
}

func (m *MockProjectStore) HasPermissionToAddAdditionalFiles(userId int, projectCode string) bool {
	return m.HasPermissionToAddAdditionalFilesFunc(userId, projectCode)
}

func (m *MockProjectStore) GetProjectStatusByProjectCode(projectCode string) (projects.AdminUpdateParam, error) {
	return m.GetProjectStatusByProjectCodeFunc(projectCode)
}

func (m *MockProjectStore) UpdateProjectByAdmin(payload projects.AdminUpdateParam, userId int, projectCode string, additionFiles []*multipart.FileHeader) error {
	m.AdminUpdateData = payload
	return nil
}

func (m *MockProjectStore) GetAdminRequestDashboard(
	fromDate,
	toDate time.Time,
	orderBy string,
	limit, offset int,
	projectCode, projectName, projectStatus *string,
) ([]projects.AdminRequestDashboardRow, error) {
	return m.GetAdminRequestDashboardFunc(fromDate, toDate, orderBy, limit, offset, projectCode, projectName, projectStatus)
}

func (m *MockProjectStore) GetAdminSummary(fromDate, toDate time.Time) ([]projects.AdminSummaryData, error) {
	return m.GetAdminSummaryFunc(fromDate, toDate)
}
