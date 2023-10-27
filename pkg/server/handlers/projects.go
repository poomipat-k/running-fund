package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/poomipat-k/running-fund/pkg/projects"
)

type projectStore interface {
	GetAll() ([]projects.Project, error)
}

type ProjectHandler struct {
	store projectStore
}

func NewProjectHandler(s projectStore) *ProjectHandler {
	return &ProjectHandler{
		store: s,
	}
}

func (h *ProjectHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	projects, err := h.store.GetAll()
	if err != nil {
		log.Panic(err)
	}

	jsonBytes, err := json.Marshal(projects)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
