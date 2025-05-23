package projects

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"

	s3Service "github.com/poomipat-k/running-fund/pkg/s3-service"
	"github.com/poomipat-k/running-fund/pkg/users"
	"github.com/poomipat-k/running-fund/pkg/utils"
)

const MAX_UPLOAD_SIZE = 25 * 1024 * 1024 // 25MB

type projectStore interface {
	GetReviewerDashboard(userId int, from time.Time, to time.Time) ([]ReviewDashboardRow, error)
	GetReviewerProjectDetails(reviewerId int, projectCode string) (ProjectReviewDetailsResponse, error)
	GetProjectCriteria(criteriaVersion int) ([]ProjectReviewCriteria, error)
	GetApplicantCriteria(version int) ([]ApplicantSelfScoreCriteria, error)
	AddProject(addProject AddProjectRequest, userId int, criteria []ApplicantSelfScoreCriteria, attachments []Attachments) (int, error)
	GetAllProjectDashboardByApplicantId(applicantId int) ([]ApplicantDashboardItem, error)
	GetApplicantProjectDetails(isAdmin bool, projectCode string, userId int) ([]ApplicantDetailsData, error)
	HasPermissionToAddAdditionalFiles(userId int, projectCode string) bool
	GetProjectStatusByProjectCode(projectCode string) (AdminUpdateParam, error)
	UpdateProjectByAdmin(payload AdminUpdateParam, userId int, projectCode string, additionFiles []*multipart.FileHeader, etcFiles []*multipart.FileHeader) error
	GetAdminRequestDashboard(fromDate, toDate time.Time, orderBy string, limit, offset int, projectCode, projectName, projectStatus *string) ([]AdminRequestDashboardRow, error)
	GetAdminStartedDashboard(fromDate, toDate time.Time, orderBy string, limit, offset int, projectCode, projectName, projectStatus *string) ([]AdminRequestDashboardRow, error)
	GetAdminSummary(fromDate, toDate time.Time) ([]AdminSummaryData, error)
	GenerateAdminReport(fromDate, toDate time.Time) (*bytes.Buffer, error)
}

type ProjectHandler struct {
	store        projectStore
	uStore       users.UserStore
	awsS3Service s3Service.S3Service
}

func NewProjectHandler(s projectStore, uStore users.UserStore, awsS3Service s3Service.S3Service) *ProjectHandler {
	return &ProjectHandler{
		store:        s,
		uStore:       uStore,
		awsS3Service: awsS3Service,
	}
}

func (h *ProjectHandler) GetReviewerDashboard(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUserIdFromRequestHeader(r)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "userId")
		return
	}

	var payload GetReviewerDashboardRequest
	err = utils.ReadJSON(w, r, &payload)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "")
		return
	}

	projects, err := h.store.GetReviewerDashboard(userId, payload.FromDate, payload.ToDate)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "")
		return
	}

	utils.WriteJSON(w, http.StatusOK, projects)
}

func (h *ProjectHandler) GetReviewerProjectDetails(w http.ResponseWriter, r *http.Request) {
	var reviewerId int
	loggedInUserId, err := utils.GetUserIdFromRequestHeader(r)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "userId")
		return
	}
	userRole := utils.GetUserRoleFromRequestHeader(r)
	if userRole == "applicant" {
		utils.ErrorJSON(w, errors.New("no permission"), "userRole", http.StatusForbidden)
		return
	}
	if userRole == "admin" {
		var payload ProjectReviewer
		utils.ReadJSON(w, r, &payload)
		if payload.ReviewerId == 0 {
			utils.ErrorJSON(w, &ReviewerIdRequiredError{}, "reviewerId")
			return
		}
		reviewerId = payload.ReviewerId

	} else if userRole == "reviewer" {
		reviewerId = loggedInUserId
	}

	projectCode := chi.URLParam(r, "projectCode")
	if projectCode == "" {
		slog.Error("Please provide a project code.")
		utils.ErrorJSON(w, err, "projectCode")
		return
	}
	projectDetails, err := h.store.GetReviewerProjectDetails(reviewerId, projectCode)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "")
		return
	}
	if userRole == "applicant" && projectDetails.UserId != loggedInUserId {
		slog.Error("applicant cannot get other applicant review details")
		utils.ErrorJSON(w, &ProjectNotFoundError{}, "", http.StatusNotFound)
		return
	}
	utils.WriteJSON(w, http.StatusOK, projectDetails)
}

func (h *ProjectHandler) GetProjectCriteria(w http.ResponseWriter, r *http.Request) {
	criteriaVersion, err := strconv.Atoi(chi.URLParam(r, "criteriaVersion"))
	if err != nil {
		criteriaVersion = 1
	}
	criteria, err := h.store.GetProjectCriteria(criteriaVersion)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "", http.StatusBadRequest)
		return
	}

	utils.WriteJSON(w, http.StatusOK, criteria)
}

func (h *ProjectHandler) GetApplicantCriteria(w http.ResponseWriter, r *http.Request) {
	applicantCriteria, err := strconv.Atoi(chi.URLParam(r, "applicantCriteriaVersion"))
	if err != nil {
		applicantCriteria = 1
	}
	criteria, err := h.store.GetApplicantCriteria(applicantCriteria)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "", http.StatusBadRequest)
		return
	}

	utils.WriteJSON(w, http.StatusOK, criteria)
}

func (h *ProjectHandler) GetAllProjectDashboardByApplicantId(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUserIdFromRequestHeader(r)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "userId", http.StatusForbidden)
		return
	}
	data, err := h.store.GetAllProjectDashboardByApplicantId(userId)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "", http.StatusBadRequest)
		return
	}
	utils.WriteJSON(w, http.StatusOK, data)
}

// ADD PROJECT START
func (h *ProjectHandler) AddProject(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUserIdFromRequestHeader(r)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "userId", http.StatusForbidden)
		return
	}

	if err := r.ParseMultipartForm(25 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	formJsonString := r.FormValue("form")
	payload := AddProjectRequest{}
	// slog.Info("AddProject payload", "userId", userId, "payload", r.Form)
	slog.Info("AddProject payload", "userId", userId)

	err = json.Unmarshal([]byte(formJsonString), &payload)
	if err != nil {
		slog.Error(err.Error(), "payload", r.Form)
		utils.ErrorJSON(w, err, "")
		return
	}

	var collaborateFiles []*multipart.FileHeader
	if payload.Collaborated != nil && *payload.Collaborated {
		collaborateFiles = r.MultipartForm.File["collaborationFiles"]
	}
	marketingFiles := r.MultipartForm.File["marketingFiles"]
	routeFiles := r.MultipartForm.File["routeFiles"]
	eventMapFiles := r.MultipartForm.File["eventMapFiles"]
	eventDetailsFiles := r.MultipartForm.File["eventDetailsFiles"]
	etcFiles := r.MultipartForm.File["etcFiles"]
	attachments := []Attachments{
		{
			DirName:         collaborationStr,
			ZipName:         collaborationStr,
			InZipFilePrefix: collaborationStr,
			Files:           collaborateFiles,
		},
		{
			DirName:         fmt.Sprintf("%s/ป้ายประชาสัมพันธ์กิจกรรม", attachmentsStr),
			ZipName:         attachmentsStr,
			InZipFilePrefix: "ป้ายประชาสัมพันธ์กิจกรรม",
			Files:           marketingFiles,
		},
		{
			DirName:         fmt.Sprintf("%s/เส้นทางจุดเริ่มต้นถึงจุดสิ้นสุดและเส้นทางวิ่งในทุกระยะ", attachmentsStr),
			ZipName:         attachmentsStr,
			InZipFilePrefix: "เส้นทางจุดเริ่มต้นถึงจุดสิ้นสุดและเส้นทางวิ่งในทุกระยะ",
			Files:           routeFiles,
		},
		{
			DirName:         fmt.Sprintf("%s/แผนผังบริเวณการจัดงาน", attachmentsStr),
			ZipName:         attachmentsStr,
			InZipFilePrefix: "แผนผังบริเวณการจัดงาน",
			Files:           eventMapFiles,
		},
		{
			DirName:         fmt.Sprintf("%s/กำหนดการการจัดกิจกรรม", attachmentsStr),
			ZipName:         attachmentsStr,
			InZipFilePrefix: "กำหนดการการจัดกิจกรรม",
			Files:           eventDetailsFiles,
		},
		{
			DirName:         fmt.Sprintf("%s/เอกสารอื่นๆ", attachmentsStr),
			ZipName:         attachmentsStr,
			InZipFilePrefix: "เอกสารอื่นๆ",
			Files:           etcFiles,
		},
	}

	v := os.Getenv("APPLICANT_CRITERIA_VERSION")
	criteriaVersion, err := strconv.Atoi(v)
	if err != nil {
		slog.Error(err.Error(), "criteriaVersion", criteriaVersion)
		utils.ErrorJSON(w, err, "APPLICANT_CRITERIA_VERSION", http.StatusBadRequest)
		return
	}
	criteria, err := h.store.GetApplicantCriteria(criteriaVersion)
	if err != nil {
		slog.Error("GetApplicantCriteria error", "error", err.Error())
		utils.ErrorJSON(w, err, "", http.StatusInternalServerError)
		return
	}

	err = validateAddProjectPayload(payload, criteria, marketingFiles, routeFiles, eventMapFiles, eventDetailsFiles)
	if err != nil {
		slog.Error("error validateAddProjectPayload", "error", err.Error(), "payload", payload)
		utils.ErrorJSON(w, err, "", http.StatusBadRequest)
		return
	}

	projectId, err := h.store.AddProject(payload, userId, criteria, attachments)
	if err != nil {
		slog.Error("error add project store", "error", err.Error(), "payload", payload)
		utils.ErrorJSON(w, err, "", http.StatusBadRequest)
		return
	}

	utils.WriteJSON(w, http.StatusOK, projectId)
}

func (h *ProjectHandler) GetApplicantProjectDetails(w http.ResponseWriter, r *http.Request) {
	projectCode := chi.URLParam(r, "projectCode")
	userRole := utils.GetUserRoleFromRequestHeader(r)
	if userRole != "admin" && userRole != "applicant" {
		utils.ErrorJSON(w, errors.New("access denied. No permission"), "userRole", http.StatusForbidden)
		return
	}

	var userId int
	userId, err := utils.GetUserIdFromRequestHeader(r)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "userId", http.StatusForbidden)
		return
	}
	isAdmin := userRole == "admin"
	projectDetails, err := h.store.GetApplicantProjectDetails(isAdmin, projectCode, userId)

	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "userId + projectCode", http.StatusNotFound)
		return
	}

	utils.WriteJSON(w, http.StatusOK, projectDetails)

}

func (h *ProjectHandler) ListApplicantFiles(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUserIdFromRequestHeader(r)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "userId", http.StatusForbidden)
		return
	}
	userRole := utils.GetUserRoleFromRequestHeader(r)
	if userRole == "" {
		msg := "userRole is empty"
		err = errors.New(msg)
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "userRole", http.StatusForbidden)
		return
	}

	var payload ListFilesRequest
	utils.ReadJSON(w, r, &payload)
	var objectKey string
	if userRole == "applicant" {
		objectKey = fmt.Sprintf("applicant/user_%d/%s", userId, payload.Prefix)
	} else {
		objectKey = fmt.Sprintf("applicant/user_%d/%s", payload.CreatedBy, payload.Prefix)
	}
	objects, err := h.awsS3Service.ListObjects(os.Getenv("AWS_S3_STORE_BUCKET_NAME"), objectKey)
	if err != nil {
		utils.ErrorJSON(w, err, "")
		return
	}

	var data []S3ObjectDetails
	for _, obj := range objects {
		data = append(data, S3ObjectDetails{
			Key:          *obj.Key,
			LastModified: *obj.LastModified,
		})
	}

	utils.WriteJSON(w, http.StatusOK, data)
}

// ADD project addition files
func (h *ProjectHandler) AddProjectAdditionFiles(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUserIdFromRequestHeader(r)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "userId", http.StatusForbidden)
		return
	}

	if err := r.ParseMultipartForm(25 << 20); err != nil {
		utils.ErrorJSON(w, err, "", http.StatusBadRequest)
		return
	}

	formJsonString := r.FormValue("form")
	payload := AddProjectFilesRequest{}
	err = json.Unmarshal([]byte(formJsonString), &payload)
	if err != nil {
		utils.ErrorJSON(w, err, "")
		return
	}
	if payload.ProjectCode == "" {
		utils.ErrorJSON(w, &ProjectCodeRequiredError{}, "projectCode")
		return
	}

	// validation
	additionFiles := r.MultipartForm.File["additionFiles"]
	etcFiles := r.MultipartForm.File["etcFiles"]
	if len(additionFiles) == 0 && len(etcFiles) == 0 {
		utils.ErrorJSON(w, &FilesRequiredError{}, "additionFiles or etcFiles", http.StatusBadRequest)
		return
	}
	userRole := utils.GetUserRoleFromRequestHeader(r)
	if userRole != "applicant" && userRole != "admin" {
		utils.ErrorJSON(w, errors.New("permission denied"), "additionFiles  or etcFiles", http.StatusForbidden)
		return
	}
	if userRole == "admin" {
		userId = payload.UserId
	}

	canAddFiles := h.store.HasPermissionToAddAdditionalFiles(userId, payload.ProjectCode)
	if !canAddFiles {
		utils.ErrorJSON(w, &ProjectNotFoundError{}, "userId,ProjectCode", http.StatusNotFound)
		return
	}

	if additionFiles != nil {
		objectPrefix := fmt.Sprintf("applicant/user_%d/%s/addition", userId, payload.ProjectCode)
		bucketName := os.Getenv("AWS_S3_STORE_BUCKET_NAME")
		err = h.awsS3Service.UploadFilesToS3(additionFiles, bucketName, objectPrefix)
		if err != nil {
			slog.Error(err.Error())
			utils.ErrorJSON(w, err, "additionFiles", http.StatusForbidden)
			return
		}
	}
	if etcFiles != nil {
		// Unused update zip file content

		// // update attachment zip file (for reviewer to download)
		// zipFile, err := h.awsS3Service.UpdateAttachmentZipContent(userId, payload.ProjectCode, etcFiles)
		// if err != nil {
		// 	slog.Error(err.Error())
		// 	utils.ErrorJSON(w, err, "updateZip", http.StatusBadRequest)
		// 	return
		// }
		// defer zipFile.Close()

		// // upload the updated zip file to s3 bucket
		// s3ZipFilePath := fmt.Sprintf("applicant/user_%d/%s/zip/%s_เอกสารแนบ.zip", userId, payload.ProjectCode, payload.ProjectCode)
		// err = h.awsS3Service.DoUploadFileToS3(zipFile, bucketName, s3ZipFilePath)
		// if err != nil {
		// 	slog.Error(err.Error())
		// 	utils.ErrorJSON(w, err, "upload zip file", http.StatusBadRequest)
		// 	return
		// }

		// END unused update zip file content
		bucketName := os.Getenv("AWS_S3_STORE_BUCKET_NAME")
		objectPrefix := fmt.Sprintf("applicant/user_%d/%s/เอกสารแนบ/เอกสารอื่นๆ", userId, payload.ProjectCode)
		err = h.awsS3Service.UploadFilesToS3(etcFiles, bucketName, objectPrefix)
		if err != nil {
			slog.Error(err.Error())
			utils.ErrorJSON(w, err, "etcFiles", http.StatusForbidden)
			return
		}
	}

	utils.WriteJSON(w, http.StatusOK, CommonSuccessResponse{Success: true, Message: "upload files successfully"})
}
