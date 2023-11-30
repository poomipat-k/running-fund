package users

import (
	"database/sql"
	"errors"
	"log/slog"
	"math/rand"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/poomipat-k/running-fund/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

const alphaNumericBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type UserStore interface {
	GetReviewers() ([]User, error)
	GetReviewerById(id int) (User, error)
	GetUserByEmail(email string) (User, error)
	AddUser(user User, toBeDeletedUserId int) (int, error)
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
		utils.ErrorJSON(w, err, 400)
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

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	// Validate if user exists
	var payload SignUpRequest
	utils.ReadJSON(w, r, &payload)
	toBeDeletedUserId, err := validateSignUpRequest(h.store, payload)
	if err != nil {
		signUpFailed(w, err)
		return
	}
	passwordToStore, err := generateHashedAndSaltedPassword(payload.Password, 8, "_")
	if err != nil {
		signUpFailed(w, err)
		return
	}
	newUser := User{
		Email:     payload.Email,
		Password:  passwordToStore,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		UserRole:  "applicant",
		Activated: false,
	}
	// Create a new user and save it
	userId, err := h.store.AddUser(newUser, toBeDeletedUserId)
	if err != nil {
		signUpFailed(w, err)
		return
	}
	// TODO: send email to activate account

	// return the created user id
	utils.WriteJSON(w, http.StatusCreated, userId)
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

func generateHashedAndSaltedPassword(password string, saltLen int, delim string) (string, error) {
	// Generate a salt
	salt := randString(saltLen)

	// Hash the salt and the password together
	toHash := strings.Join([]string{password, salt}, "")
	hashed, err := hashPassword(toHash)
	if err != nil {
		return "", err
	}
	// Join the hashed result and the salt together
	passwordToStore := strings.Join([]string{hashed, salt}, delim)
	return passwordToStore, nil
}

func validateSignUpRequest(store UserStore, payload SignUpRequest) (int, error) {
	if payload.Email == "" {
		return 0, errors.New("email is required")
	}
	if !isValidEmail(payload.Email) {
		return 0, errors.New("email is invalid")
	}

	if payload.Password == "" {
		return 0, errors.New("password is required")
	}
	if len(payload.Password) < 8 {
		return 0, errors.New("password minimum length are 8 characters")
	}
	if len(payload.Password) > 60 {
		return 0, errors.New("password maximum length are 60 characters")
	}
	if payload.FirstName == "" {
		return 0, errors.New("first name is required")
	}
	if payload.LastName == "" {
		return 0, errors.New("last name is required")
	}

	toBeDeletedUserId, err := isDuplicatedEmail(payload.Email, store)
	if err != nil {
		return 0, err
	}
	return toBeDeletedUserId, nil
}

func isDuplicatedEmail(email string, store UserStore) (int, error) {
	user, err := store.GetUserByEmail(email)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if user.Id > 0 && (!user.Activated && time.Now().After(user.ActivatedBefore)) {
		return user.Id, nil
	}
	return 0, errors.New("email is already exist")
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func signUpFailed(w http.ResponseWriter, err error) {
	slog.Error(err.Error())
	utils.ErrorJSON(w, err, http.StatusBadRequest)
}

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphaNumericBytes[rand.Int63()%int64(len(alphaNumericBytes))]
	}
	return string(b)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}
