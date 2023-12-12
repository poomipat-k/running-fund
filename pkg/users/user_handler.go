package users

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"net/mail"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jordan-wright/email"
	appEmail "github.com/poomipat-k/running-fund/pkg/email"
	"github.com/poomipat-k/running-fund/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

const alphaNumericBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

const accessExpireDurationMinute = 30
const refreshExpireDurationHour = 4320 // 180 days

type UserStore interface {
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
	mail := email.Email{
		From:    os.Getenv("EMAIL_SENDER"),
		To:      []string{newUser.Email},
		Subject: fmt.Sprintf("Registration confirmation - %s", newUser.Email),
		Text:    []byte("Text Body is, of course, supported!"),
		HTML:    []byte("<h1>Fancy HTML is supported, too!</h1><br><p>Hi Sis A'Serene</p>"),
	}
	err = appEmail.SendEmail(mail)
	if err != nil {
		fail(w, err)
		return
	}
	slog.Info("Sign up confirmation sent to", "email", newUser.Email)
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

	accessExpiredAtUnix := time.Now().Add(accessExpireDurationMinute * time.Minute).Unix()
	accessToken, err := generateAccessToken(user.Id, user.UserRole, accessExpiredAtUnix)
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
	}
	accessTokenCookie := http.Cookie{
		Name:     "authToken",
		Value:    accessToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/api",
		Expires:  time.Unix(accessExpiredAtUnix, 0),
	}

	refreshExpiredAtUnix := time.Now().Add(refreshExpireDurationHour * time.Hour).Unix()
	refreshToken, err := generateRefreshToken(user, refreshExpiredAtUnix)
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
	}
	refreshTokenCookie := http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/api/v1/auth",
		Expires:  time.Unix(refreshExpiredAtUnix, 0),
	}

	http.SetCookie(w, &accessTokenCookie)
	http.SetCookie(w, &refreshTokenCookie)

	utils.WriteJSON(w, http.StatusOK, CommonSuccessResponse{Success: true, Message: "log in successfully"})
}

func (h *UserHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	accessTokenCookie := http.Cookie{
		Name:     "authToken",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		Path:     "/api",
		Expires:  time.Now(),
	}
	http.SetCookie(w, &accessTokenCookie)
	refreshTokenCookie := http.Cookie{
		Name:     "refreshToken",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		Path:     "/api/v1/auth",
		Expires:  time.Now(),
	}
	http.SetCookie(w, &refreshTokenCookie)

	utils.WriteJSON(w, http.StatusOK, CommonSuccessResponse{Success: true, Message: "log out successfully"})
}

func (h *UserHandler) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := getRefreshToken(r)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusForbidden)
		return
	}
	claims, ok := refreshToken.Claims.(jwt.MapClaims)
	if ok {
		userId := fmt.Sprintf("%v", claims["userId"])
		userRole := fmt.Sprintf("%v", claims["userRole"])

		accessExpiredAtUnix := time.Now().Add(accessExpireDurationMinute * time.Minute).Unix()
		uid, err := strconv.Atoi(userId)
		if err != nil {
			utils.ErrorJSON(w, err, http.StatusForbidden)
			return
		}
		accessToken, err := generateAccessToken(uid, userRole, accessExpiredAtUnix)
		if err != nil {
			utils.ErrorJSON(w, err, http.StatusForbidden)
			return
		}
		newAccessTokenCookie := http.Cookie{
			Name:     "authToken",
			Value:    accessToken,
			HttpOnly: true,
			Secure:   true,
			Path:     "/api",
			Expires:  time.Unix(accessExpiredAtUnix, 0),
		}
		http.SetCookie(w, &newAccessTokenCookie)
		utils.WriteJSON(w, http.StatusOK, CommonSuccessResponse{Success: true, Message: "Access token refresh successfully"})
		return
	}
	utils.ErrorJSON(w, errors.New("corrupt refresh token"), http.StatusForbidden)

}

func getRefreshToken(r *http.Request) (*jwt.Token, error) {
	cookie, err := r.Cookie("refreshToken")
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		// validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_REFRESH_TOKEN_SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func generateRefreshToken(user User, expiredAtUnix int64) (string, error) {
	refreshTokenSecretKey := []byte(os.Getenv("JWT_REFRESH_TOKEN_SECRET_KEY"))
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   user.Id,
		"userRole": user.UserRole,
		"iat":      time.Now().Unix(),
		"exp":      expiredAtUnix,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := t.SignedString(refreshTokenSecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func generateAccessToken(userId int, userRole string, expiredAtUnix int64) (string, error) {
	accessSecretKey := []byte(os.Getenv("JWT_ACCESS_TOKEN_SECRET_KEY"))

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   userId,
		"userRole": userRole,
		"iat":      time.Now().Unix(),
		"exp":      expiredAtUnix,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := t.SignedString(accessSecretKey)
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
