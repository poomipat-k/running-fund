package server

import (
	"encoding/json"
	"log"
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
}

type ProjectHandler struct {
	store projectStore
}

func NewProjectHandler(s projectStore) *ProjectHandler {
	return &ProjectHandler{
		store: s,
	}
}

func (h *ProjectHandler) GetReviewerDashboard(w http.ResponseWriter, r *http.Request) {
	// To check if the user exists in the db
	userId, err := getAuthUserId(r)
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(r.Body)
	var payload projects.GetReviewerDashboardRequest
	err = decoder.Decode(&payload)
	if err != nil {
		panic(err)
	}

	projects, err := h.store.GetReviewerDashboard(userId, payload.FromDate, payload.ToDate)
	if err != nil {
		log.Panic(err)
	}

	ResponseJson(w, projects, http.StatusAccepted)
}

func (h *ProjectHandler) GetReviewerProjectDetails(w http.ResponseWriter, r *http.Request) {
	// To check if the user exists in the db
	userId, err := getAuthUserId(r)
	if err != nil {
		panic(err)
	}

	projectCode := chi.URLParam(r, "projectCode")
	if len(projectCode) == 0 {
		panic("Please provide a project code.")
	}
	projectDetails, err := h.store.GetReviewerProjectDetails(userId, projectCode)
	if err != nil {
		panic(err)
	}
	ResponseJson(w, projectDetails, http.StatusOK)
}

func (h *ProjectHandler) GetReviewPeriod(w http.ResponseWriter, r *http.Request) {
	period, err := h.store.GetReviewPeriod()
	if err != nil {
		panic(err)
	}
	ResponseJson(w, period, http.StatusOK)
}

func (h *ProjectHandler) GetProjectCriteria(w http.ResponseWriter, r *http.Request) {
	criteriaVersion, err := strconv.Atoi(chi.URLParam(r, "criteriaVersion"))
	if err != nil {
		criteriaVersion = 1
	}
	criteria, err := h.store.GetProjectCriteria(criteriaVersion)
	if err != nil {
		panic(err)
	}
	ResponseJson(w, criteria, http.StatusOK)
}
