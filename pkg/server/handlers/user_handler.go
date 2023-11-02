package server

import (
	"encoding/json"
	"net/http"

	"github.com/poomipat-k/running-fund/pkg/users"
)

type userStore interface {
	GetReviewers() ([]users.User, error)
}

type UserHandler struct {
	store userStore
}

func NewUserHandler(s userStore) *UserHandler {
	return &UserHandler{
		store: s,
	}
}

func (h *UserHandler) GetReviewers(w http.ResponseWriter, r *http.Request) {
	reviewers, err := h.store.GetReviewers()
	if err != nil {
		panic(err)
	}
	jsonBytes, err := json.Marshal(reviewers)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
