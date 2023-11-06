package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/poomipat-k/running-fund/pkg/projects"
)

type projectStore interface {
	GetReviewerDashboard(userId int, from time.Time, to time.Time) ([]projects.ReviewDashboardRow, error)
	GetReviewPeriod() (projects.ReviewPeriod, error)
	GetReviewerProejctDetails(userId int, projectCode string) (projects.ProjectReviewDetails, error)
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
	authHeader := r.Header.Get("Authorization")
	splits := strings.Split(authHeader, " ")
	var token string
	if len(splits) > 1 {
		token = splits[1]
	}
	if token == "" {
		panic("Token not found")
	}

	decoder := json.NewDecoder(r.Body)
	var payload projects.GetReviewerDashboardRequest
	err := decoder.Decode(&payload)
	if err != nil {
		panic(err)
	}

	// To check if the user exists in the db
	userId, err := strconv.Atoi(token)
	if err != nil {
		panic(err)
	}
	projects, err := h.store.GetReviewerDashboard(userId, payload.FromDate, payload.ToDate)
	if err != nil {
		log.Panic(err)
	}

	ResposeJson(w, projects, http.StatusAccepted)
}

func (h *ProjectHandler) GetReviewerProejctDetails(w http.ResponseWriter, r *http.Request) {
	// To check if the user exists in the db
	userId, err := getAuthUserId(r)
	if err != nil {
		panic(err)
	}

	projectCode := chi.URLParam(r, "projectCode")
	if len(projectCode) == 0 {
		panic("Please provide a project code.")
	}
	projectDetails, err := h.store.GetReviewerProejctDetails(userId, projectCode)
	if err != nil {
		panic(err)
	}
	ResposeJson(w, projectDetails, http.StatusOK)
}

func (h *ProjectHandler) GetReviewPeriod(w http.ResponseWriter, r *http.Request) {
	period, err := h.store.GetReviewPeriod()
	if err != nil {
		panic(err)
	}
	ResposeJson(w, period, http.StatusOK)
}
