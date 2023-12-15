package users

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jordan-wright/email"
	"github.com/poomipat-k/running-fund/pkg/utils"
)

const accessExpireDurationMinute = 30
const refreshExpireDurationHour = 4320 // 180 days

type UserStore interface {
	GetUserByEmail(email string) (User, error)
	GetUserById(id int) (User, error)
	AddUser(user User, toBeDeletedUserId int) (int, error)
	ActivateUser(activateCode string) (int64, error)
	ForgotPasswordAction(resetPasswordCode string, email string, resetPasswordLink string) (int64, error)
	ResetPassword(resetPasswordCode string, newPassword string) (int64, error)
}

type EmailService interface {
	SendEmail(email email.Email) error
	BuildSignUpConfirmationEmail(email, activateLink string) email.Email
	BuildResetPasswordEmail(to, resetPasswordLink string) email.Email
}

type UserHandler struct {
	store        UserStore
	emailService EmailService
}

func NewUserHandler(s UserStore, es EmailService) *UserHandler {
	return &UserHandler{
		store:        s,
		emailService: es,
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
	activateCode := utils.RandAlphaNum(24)
	newUser := User{
		Email:        payload.Email,
		Password:     passwordToStore,
		FirstName:    payload.FirstName,
		LastName:     payload.LastName,
		UserRole:     "applicant",
		Activated:    false,
		ActivateCode: activateCode,
	}
	// Create a new user and save it
	userId, err := h.store.AddUser(newUser, toBeDeletedUserId)
	if err != nil {
		fail(w, err)
		return
	}

	activateLink := fmt.Sprintf("%s/user/activate-email?activateCode=%s", os.Getenv("BACKEND_URL"), activateCode)
	mail := h.emailService.BuildSignUpConfirmationEmail(newUser.Email, activateLink)
	err = h.emailService.SendEmail(mail)
	if err != nil {
		fail(w, err)
		return
	}
	slog.Info("Sign up confirmation email sent to", "email", newUser.Email)
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

func (h *UserHandler) ActivateUser(w http.ResponseWriter, r *http.Request) {
	activateCode := r.URL.Query().Get("activateCode")
	err := validateActivateCode(activateCode)
	if err != nil {
		fail(w, err, http.StatusBadRequest)
		return
	}

	rowEffected, err := h.store.ActivateUser(activateCode)
	if err != nil {
		fail(w, err, http.StatusNotFound)
		return
	}
	if rowEffected == 0 {
		fail(w, &UserToActivateNotFoundError{}, http.StatusNotFound)
		return
	}
	utils.WriteJSON(w, http.StatusOK, rowEffected)
}

func (h *UserHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var payload ForgotPasswordRequest
	err := utils.ReadJSON(w, r, &payload)
	if err != nil {
		fail(w, err)
		return
	}
	err = validateEmail(payload.Email)
	if err != nil {
		fail(w, err)
		return
	}

	user, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		fail(w, err)
		return
	}
	if !user.Activated {
		fail(w, &UserIsNotActivatedError{}, http.StatusBadRequest)
		return
	}

	resetPasswordCode := utils.RandAlphaNum(24)
	resetPasswordLink := fmt.Sprintf("%s/user/reset-password?resetPasswordCode=%s", os.Getenv("BACKEND_URL"), resetPasswordCode)
	rowEffected, err := h.store.ForgotPasswordAction(resetPasswordCode, user.Email, resetPasswordLink)
	if err != nil {
		fail(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, rowEffected)
}

func (h *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var payload ResetPasswordRequest
	err := utils.ReadJSON(w, r, &payload)
	if err != nil {
		fail(w, err)
		return
	}
	if payload.Password != payload.ConfirmPassword {
		fail(w, &PasswordAndConfirmPasswordNotMatchError{})
		return
	}
	err = validatePassword(payload.Password)
	if err != nil {
		fail(w, err)
		return
	}
	if len(payload.ResetPasswordCode) != 24 {
		fail(w, &ResetPasswordCodeNotValidError{})
		return
	}

	passwordToStore, err := generateHashedAndSaltedPassword(payload.Password, 8, "_")
	if err != nil {
		fail(w, err)
		return
	}

	rowEffected, err := h.store.ResetPassword(payload.ResetPasswordCode, passwordToStore)
	if err != nil {
		fail(w, err)
		return
	}
	if rowEffected == 0 {
		fail(w, &ResetPasswordCodeNotFound{}, http.StatusNotFound)
		return
	}
	utils.WriteJSON(w, http.StatusOK, rowEffected)
}
