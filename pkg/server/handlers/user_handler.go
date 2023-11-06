package server

import (
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
	ResposeJson(w, reviewers, http.StatusOK)
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
	ResposeJson(w, reviewer, http.StatusOK)
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
