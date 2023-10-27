package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/poomipat-k/running-fund/pkg/projects"
)

type projectStore interface {
	GetReviewerDashboard(from time.Time, to time.Time) ([]projects.Project, error)
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
	decoder := json.NewDecoder(r.Body)
	var payload projects.GetReviewerDashboardRequest
	err := decoder.Decode(&payload)
	if err != nil {
		panic(err)
	}

	projects, err := h.store.GetReviewerDashboard(payload.From, payload.To)
	if err != nil {
		log.Panic(err)
	}

	jsonBytes, err := json.Marshal(projects)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
