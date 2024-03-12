package projects

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-chi/chi"

	s3Service "github.com/poomipat-k/running-fund/pkg/s3-service"
	"github.com/poomipat-k/running-fund/pkg/users"
	"github.com/poomipat-k/running-fund/pkg/utils"
)

const MAX_UPLOAD_SIZE = 25 * 1024 * 1024 // 25MB

type projectStore interface {
	GetReviewerDashboard(userId int, from time.Time, to time.Time) ([]ReviewDashboardRow, error)
	GetReviewPeriod() (ReviewPeriod, error)
	GetReviewerProjectDetails(reviewerId int, projectCode string) (ProjectReviewDetailsResponse, error)
	GetProjectCriteria(criteriaVersion int) ([]ProjectReviewCriteria, error)
	GetApplicantCriteria(version int) ([]ApplicantSelfScoreCriteria, error)
	AddProject(addProject AddProjectRequest, userId int, criteria []ApplicantSelfScoreCriteria, attachments []Attachments) (int, error)
	GetAllProjectDashboardByApplicantId(applicantId int) ([]ApplicantDashboardItem, error)
	GetApplicantProjectDetails(userId int, projectCode string) ([]ApplicantDetailsData, error)
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
	reviewerId, err := utils.GetUserIdFromRequestHeader(r)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "userId")
		return
	}

	projectCode := chi.URLParam(r, "projectCode")
	if len(projectCode) == 0 {
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
	utils.WriteJSON(w, http.StatusOK, projectDetails)
}

func (h *ProjectHandler) GetReviewPeriod(w http.ResponseWriter, r *http.Request) {
	period, err := h.store.GetReviewPeriod()
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "")
		return
	}

	utils.WriteJSON(w, http.StatusOK, period)
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
	slog.Info("AddProject payload", "userId", userId, "payload", r.Form)

	err = json.Unmarshal([]byte(formJsonString), &payload)
	if err != nil {
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
	screenshotFiles := r.MultipartForm.File["screenshotFiles"]
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
			DirName:         formStr,
			ZipName:         formStr,
			InZipFilePrefix: formStr,
			Files:           screenshotFiles,
		},
	}

	v := os.Getenv("APPLICANT_CRITERIA_VERSION")
	criteriaVersion, err := strconv.Atoi(v)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "APPLICANT_CRITERIA_VERSION", http.StatusBadRequest)
		return
	}
	criteria, err := h.store.GetApplicantCriteria(criteriaVersion)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "", http.StatusInternalServerError)
		return
	}

	err = validateAddProjectPayload(payload, collaborateFiles, criteria, marketingFiles, routeFiles, eventMapFiles, eventDetailsFiles, screenshotFiles)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "", http.StatusBadRequest)
		return
	}

	projectId, err := h.store.AddProject(payload, userId, criteria, attachments)
	if err != nil {
		utils.ErrorJSON(w, err, "", http.StatusBadRequest)
		return
	}

	utils.WriteJSON(w, http.StatusOK, projectId)
}

func (h *ProjectHandler) Download(w http.ResponseWriter, r *http.Request) {
	client := h.awsS3Service.S3Client
	manager := manager.NewDownloader(client)

	Prefix := "applicant/user_3/FEB67_2801/"
	Bucket := os.Getenv("AWS_S3_STORE_BUCKET_NAME")

	paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
		Bucket: &Bucket,
		Prefix: &Prefix,
	})
	LocalDirectory := "home/download"

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			utils.ErrorJSON(w, err, "")
		}
		for _, obj := range page.Contents {
			if err := downloadToFile(manager, LocalDirectory, Bucket, aws.ToString(obj.Key)); err != nil {
				utils.ErrorJSON(w, err, "")
			}
		}
	}
}

func (h *ProjectHandler) GetApplicantProjectDetails(w http.ResponseWriter, r *http.Request) {
	projectCode := chi.URLParam(r, "projectCode")

	userId, err := utils.GetUserIdFromRequestHeader(r)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "userId", http.StatusForbidden)
		return
	}
	projectDetails, err := h.store.GetApplicantProjectDetails(userId, projectCode)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "userId + projectCode", http.StatusNotFound)
		return
	}

	utils.WriteJSON(w, http.StatusOK, projectDetails)

}

func downloadToFile(downloader *manager.Downloader, targetDirectory, bucket, key string) error {
	// Create the directories in the path
	file := filepath.Join(targetDirectory, key)
	if err := os.MkdirAll(filepath.Dir(file), 0775); err != nil {
		return err
	}

	// Set up the local file
	fd, err := os.Create(file)
	if err != nil {
		return err
	}
	defer fd.Close()

	// Download the file using the AWS SDK for Go
	fmt.Printf("Downloading s3://%s/%s to %s...\n", bucket, key, file)
	_, err = downloader.Download(context.TODO(), fd, &s3.GetObjectInput{Bucket: &bucket, Key: &key})

	return err
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
