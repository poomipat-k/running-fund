package users

import (
	"database/sql"
	"errors"
	"log/slog"
	"math/rand"
	"net/http"
	"net/mail"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/poomipat-k/running-fund/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

const alphaNumericBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type UserStore interface {
	GetReviewers() ([]User, error)
	GetReviewerById(id int) (User, error)
	GetUserByEmail(email string) (User, error)
	GetUserById(id int) (User, error)
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
		utils.ErrorJSON(w, err, http.StatusBadRequest)
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

func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUserIdFromRequestHeader(r)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, http.StatusForbidden)
		return
	}
	userRole := r.Header.Get("userRole")
	if userRole == "" {
		err = errors.New("user role is invalid")
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, http.StatusForbidden)
		return
	}
	user, err := h.store.GetUserById(userId)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	// Validate if user exists
	var payload SignUpRequest
	err := utils.ReadJSON(w, r, &payload)
	if err != nil {
		fail(w, err)
		return
	}

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
	err := utils.ReadJSON(w, r, &payload)
	if err != nil {
		fail(w, err)
		return
	}
	err = validateSignInRequest(payload)
	if err != nil {
		fail(w, err)
		return
	}

	user, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		fail(w, &InvalidLoginCredentialError{})
		return
	}

	if !user.Activated {
		fail(w, &UserNotActivatedError{})
		return
	}

	err = comparePassword(payload.Password, user.Password)
	if err != nil {
		fail(w, &InvalidLoginCredentialError{}, http.StatusUnauthorized)
		return
	}

	// Return jwt token when log in successfully
	token, err := generateJwtToken(user)
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
	}
	tokenCookie := http.Cookie{
		Name:     "authToken",
		Value:    token,
		HttpOnly: true,
		// SameSite: http.SameSiteNoneMode,
		Secure: true,
		Path:   "/api",
	}

	http.SetCookie(w, &tokenCookie)
	utils.WriteJSON(w, http.StatusOK, SignInResponse{Success: true})
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

func generateJwtToken(user User) (string, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   user.Id,
		"userRole": user.UserRole,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(5 * time.Hour).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := t.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func comparePassword(inputPassword string, userPassword string) error {
	splitStr := strings.Split(userPassword, "_")
	if len(splitStr) != 2 {
		return errors.New("user password is invalid")
	}
	hash := splitStr[0]
	salt := splitStr[1]
	// user provided password + salt compare to hashed
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(strings.Join([]string{inputPassword, salt}, "")))
	if err != nil {
		return err
	}
	return nil
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

	err = validateFirstName(payload.FirstName)
	if err != nil {
		return 0, err
	}

	err = validateLastName(payload.LastName)
	if err != nil {
		return 0, err
	}

	toBeDeletedUserId, err := isDuplicatedEmail(payload.Email, store)
	if err != nil {
		return 0, err
	}
	return toBeDeletedUserId, nil
}

func validateSignInRequest(payload SignInRequest) error {
	err := validateEmail(payload.Email)
	if err != nil {
		return err
	}
	err = validatePassword(payload.Password)
	if err != nil {
		return err
	}
	return nil
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
		return &PasswordRequiredError{}
	}
	if len(password) < 8 {
		return &PasswordTooShortError{}
	}
	if len(password) > 60 {
		return &PasswordTooLongError{}
	}
	return nil
}

func validateFirstName(firstName string) error {
	if firstName == "" {
		return &FirstNameRequiredError{}
	}
	if len(firstName) > 255 {
		return &FirstNameTooLongError{}
	}
	return nil
}

func validateLastName(lastName string) error {
	if lastName == "" {
		return &LastNameRequiredError{}
	}
	if len(lastName) > 255 {
		return &LastNameTooLongError{}
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
