package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/poomipat-k/running-fund/pkg/projects"
)

type projectStore interface {
	GetReviewerDashboard(userId int, from time.Time, to time.Time) ([]projects.ReviewDashboardRow, error)
	GetReviewPeriod() (projects.ReviewPeriod, error)
	GetReviewerProjectDetails(userId int, projectCode string) (projects.ProjectReviewDetails, error)
	GetProjectCriteria(criteriaVersion int) ([]projects.ProjectReviewCriteria, error)
	GetProjectCriteriaMinimalDetails(cv int) ([]projects.ProjectReviewCriteriaMinimal, error)
	AddReview(payload projects.AddReviewRequest, userId int, criteriaList []projects.ProjectReviewCriteriaMinimal) (int, error)
}

type ProjectHandler struct {
	store  projectStore
	uStore userStore
}

func NewProjectHandler(s projectStore, uStore userStore) *ProjectHandler {
	return &ProjectHandler{
		store:  s,
		uStore: uStore,
	}
}

func (h *ProjectHandler) GetReviewerDashboard(w http.ResponseWriter, r *http.Request) {
	// To check if the user exists in the db
	userId, err := getAuthUserId(r)
	if err != nil {
		slog.Error(err.Error())
		errorJSON(w, err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var payload projects.GetReviewerDashboardRequest
	err = decoder.Decode(&payload)
	if err != nil {
		slog.Error(err.Error())
		errorJSON(w, err)
		return
	}

	projects, err := h.store.GetReviewerDashboard(userId, payload.FromDate, payload.ToDate)
	if err != nil {
		slog.Error(err.Error())
		errorJSON(w, err)
		return
	}

	writeJSON(w, http.StatusOK, projects)
}

func (h *ProjectHandler) GetReviewerProjectDetails(w http.ResponseWriter, r *http.Request) {
	// To check if the user exists in the db
	userId, err := getAuthUserId(r)
	if err != nil {
		slog.Error(err.Error())
		errorJSON(w, err)
		return
	}

	projectCode := chi.URLParam(r, "projectCode")
	if len(projectCode) == 0 {
		slog.Error("Please provide a project code.")
		errorJSON(w, err)
		return
	}
	projectDetails, err := h.store.GetReviewerProjectDetails(userId, projectCode)
	if err != nil {
		slog.Error(err.Error())
		errorJSON(w, err)
		return
	}
	// ResponseJson(w, projectDetails, http.StatusOK)
	writeJSON(w, http.StatusOK, projectDetails)
}

func (h *ProjectHandler) GetReviewPeriod(w http.ResponseWriter, r *http.Request) {
	period, err := h.store.GetReviewPeriod()
	if err != nil {
		slog.Error(err.Error())
		errorJSON(w, err)
		return
	}

	writeJSON(w, http.StatusOK, period)
}

func (h *ProjectHandler) GetProjectCriteria(w http.ResponseWriter, r *http.Request) {
	criteriaVersion, err := strconv.Atoi(chi.URLParam(r, "criteriaVersion"))
	if err != nil {
		criteriaVersion = 1
	}
	criteria, err := h.store.GetProjectCriteria(criteriaVersion)
	if err != nil {
		slog.Error(err.Error())
		errorJSON(w, err, http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, criteria)
}
