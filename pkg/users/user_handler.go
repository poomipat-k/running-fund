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
		fail(w, err)
		return
	}
	passwordToStore, err := generateHashedAndSaltedPassword(payload.Password, 8, "_")
	if err != nil {
		fail(w, err)
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
		fail(w, err)
		return
	}
	// TODO: send email to activate account

	// return the created user id
	utils.WriteJSON(w, http.StatusCreated, userId)
}

func (h *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	// Find user from email
	var payload SignInRequest
	utils.ReadJSON(w, r, &payload)
	err := validateSignInRequest(payload)
	if err != nil {
		fail(w, err)
		return
	}

	user, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		fail(w, err)
		return
	}

	// get hash and salt from user.password
	splitStr := strings.Split(user.Password, "_")
	if len(splitStr) != 2 {
		fail(w, errors.New("something wrong with user password"))
		return
	}
	hash := splitStr[0]
	salt := splitStr[1]
	// user provided password + salt compare to hashed
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(strings.Join([]string{payload.Password, salt}, "")))
	if err != nil {
		fail(w, err, http.StatusUnauthorized)
		return
	}

	// Return jwt token when log in successfully
	utils.WriteJSON(w, http.StatusOK, "logged in")
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
	salt := randString(saltLen)

	toHash := strings.Join([]string{password, salt}, "")
	hashed, err := hashPassword(toHash)
	if err != nil {
		return "", err
	}

	passwordToStore := strings.Join([]string{hashed, salt}, delim)
	return passwordToStore, nil
}

func validateSignUpRequest(store UserStore, payload SignUpRequest) (int, error) {
	err := validateEmail(payload.Email)
	if err != nil {
		return 0, err
	}

	err = validatePassword(payload.Password)
	if err != nil {
		return 0, err
	}
	if payload.FirstName == "" {
		return 0, &FirstNameRequiredError{}
	}
	if len(payload.FirstName) > 255 {
		return 0, &FirstNameTooLongError{}
	}
	if payload.LastName == "" {
		return 0, &LastNameRequiredError{}
	}
	if len(payload.LastName) > 255 {
		return 0, &LastNameTooLongError{}
	}

	toBeDeletedUserId, err := isDuplicatedEmail(payload.Email, store)
	if err != nil {
		return 0, err
	}
	return toBeDeletedUserId, nil
}

func validateEmail(email string) error {
	if email == "" {
		return &EmailRequiredError{}
	}
	if len(email) > 255 {
		return &EmailTooLongError{}
	}
	if !isValidEmail(email) {
		return &InvalidEmailError{}
	}
	return nil
}

func validatePassword(password string) error {
	if password == "" {
		return errors.New("password is required")
	}
	if len(password) < 8 {
		return &PasswordTooShortError{}
	}
	if len(password) > 60 {
		return &PasswordTooLongError{}
	}
	return nil
}

func validateSignInRequest(payload SignInRequest) error {
	if payload.Email == "" {
		return &EmailRequiredError{}
	}
	return nil
}

func isDuplicatedEmail(email string, store UserStore) (int, error) {
	user, err := store.GetUserByEmail(email)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if user.Id > 0 && (!user.Activated && time.Now().After(user.ActivatedBefore)) {
		return user.Id, nil
	}
	return 0, &DuplicatedEmailError{}
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func fail(w http.ResponseWriter, err error, status ...int) {
	slog.Error(err.Error())
	s := http.StatusBadRequest
	if len(status) > 0 {
		s = status[0]
	}
	utils.ErrorJSON(w, err, s)
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
