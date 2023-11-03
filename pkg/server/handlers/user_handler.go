package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/poomipat-k/running-fund/pkg/users"
)

type userStore interface {
	GetReviewers() ([]users.User, error)
	GetReviewerById(id int) (users.User, error)
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

func (h *UserHandler) GetReviewerById(w http.ResponseWriter, r *http.Request) {
	userId, err := getAuthUserId(r)
	if err != nil {
		panic(err)
	}
	reviewer, err := h.store.GetReviewerById(userId)
	if err != nil {
		panic(err)
	}

	jsonBytes, err := json.Marshal(reviewer)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func getAuthUserId(r *http.Request) (int, error) {
	authHeader := r.Header.Get("Authorization")
	splits := strings.Split(authHeader, " ")
	var token string
	if len(splits) > 1 {
		token = splits[1]
		userId, err := strconv.Atoi(token)
		if err != nil {
			return 0, errors.New("invalid token")
		}
		return userId, nil
	} else {
		return 0, errors.New("invalid token")
	}
}
