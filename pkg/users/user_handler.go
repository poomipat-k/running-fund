package users

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/poomipat-k/running-fund/pkg/utils"
)

type UserStore interface {
	GetReviewers() ([]User, error)
	GetReviewerById(id int) (User, error)
	GetUserByEmail(email string) (User, error)
}

type UserHandler struct {
	store UserStore
}

func NewUserHandler(s UserStore) *UserHandler {
	return &UserHandler{
		store: s,
	}
}

func (h *UserHandler) GetReviewers(w http.ResponseWriter, r *http.Request) {
	reviewers, err := h.store.GetReviewers()
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, reviewers)
}

func (h *UserHandler) GetReviewerById(w http.ResponseWriter, r *http.Request) {
	userId, err := GetAuthUserId(r)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err)
		return
	}
	reviewer, err := h.store.GetReviewerById(userId)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, reviewer)
}

func (h *UserHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {

}

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	// Validate if user exists
	var payload SignUpRequest
	utils.ReadJSON(w, r, &payload)
	err := validateSignUpRequest(h.store, payload)
	if err != nil {
		signUpFailed(w, err)
		return
	}

	// Generate a salt

	// Hash the salt and the password together

	// Join the hashed result and the salt together

	// Create a new user and save it

	// return the user
	utils.WriteJSON(w, http.StatusCreated, 1)
}

func validateSignUpRequest(store UserStore, payload SignUpRequest) error {
	if payload.Email == "" {
		return errors.New("email is required")
	}
	if payload.Password == "" {
		return errors.New("password is required")
	}
	if payload.FirstName == "" {
		return errors.New("first name is required")
	}
	if payload.LastName == "" {
		return errors.New("last name is required")
	}
	_, err := store.GetUserByEmail(payload.Email)
	if err == nil {
		return errors.New("email is already exist")
	}
	return nil
}

func signUpFailed(w http.ResponseWriter, err error) {
	slog.Error(err.Error())
	utils.ErrorJSON(w, err, http.StatusBadRequest)
}

func GetAuthUserId(r *http.Request) (int, error) {
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
		return 0, errors.New("no token provided")
	}
}
