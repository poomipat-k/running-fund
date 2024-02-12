package projects

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/poomipat-k/running-fund/pkg/upload"
	"github.com/poomipat-k/running-fund/pkg/users"
	"github.com/poomipat-k/running-fund/pkg/utils"
)

const MAX_UPLOAD_SIZE = 32 * 1024 * 1024 // 32MB

type projectStore interface {
	GetReviewerDashboard(userId int, from time.Time, to time.Time) ([]ReviewDashboardRow, error)
	GetReviewPeriod() (ReviewPeriod, error)
	GetReviewerProjectDetails(userId int, projectCode string) (ProjectReviewDetails, error)
	GetProjectCriteria(criteriaVersion int) ([]ProjectReviewCriteria, error)
	GetApplicantCriteria(version int) ([]ApplicantSelfScoreCriteria, error)
	AddProject() (string, error)
}

type ProjectHandler struct {
	store         projectStore
	uStore        users.UserStore
	uploadService upload.S3Service
}

func NewProjectHandler(s projectStore, uStore users.UserStore, us upload.S3Service) *ProjectHandler {
	return &ProjectHandler{
		store:         s,
		uStore:        uStore,
		uploadService: us,
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
	userId, err := utils.GetUserIdFromRequestHeader(r)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "userId")
		return
	}
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "")
		return
	}

	projectCode := chi.URLParam(r, "projectCode")
	if len(projectCode) == 0 {
		slog.Error("Please provide a project code.")
		utils.ErrorJSON(w, err, "projectCode")
		return
	}
	projectDetails, err := h.store.GetReviewerProjectDetails(userId, projectCode)
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

	err = json.Unmarshal([]byte(formJsonString), &payload)
	if err != nil {
		utils.ErrorJSON(w, err, "")
		return
	}

	log.Println("======PAYLOAD======")
	log.Println(payload)
	// get a reference to the fileHeaders

	projectCode, err := h.store.AddProject()
	if err != nil {
		utils.ErrorJSON(w, err, "")
		return
	}

	collaborateFiles := r.MultipartForm.File["collaborationFiles"]
	basePrefix := fmt.Sprintf("applicant/user_%d/%s", userId, projectCode)

	err = h.uploadService.UploadToS3(collaborateFiles, fmt.Sprintf("%s/collaboration", basePrefix))
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "collaboration", http.StatusInternalServerError)
		return
	}
	detailsFiles := r.MultipartForm.File["files"]
	err = h.uploadService.UploadToS3(detailsFiles, basePrefix)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "details", http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, "OK")
}

func (h *ProjectHandler) ListObjectsV2(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		BucketName string `json:"bucketName"`
		Prefix     string `json:"prefix"`
	}

	utils.ReadJSON(w, r, &payload)
	objects, err := h.uploadService.ListObjects(payload.BucketName, payload.Prefix)
	if err != nil {
		log.Println("err: ", err)
		utils.ErrorJSON(w, err, "")
		return
	}
	log.Println(objects)
	utils.WriteJSON(w, 200, objects)
}
