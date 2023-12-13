package users

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/poomipat-k/running-fund/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

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
	salt := utils.RandAlphaNum(saltLen)

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

func validateActivateCode(activateCode string) error {
	if len(activateCode) != 24 {
		return &InvalidActivateCodeError{}
	}
	return nil
}

func fail(w http.ResponseWriter, err error, status ...int) {
	slog.Error(err.Error())
	s := http.StatusBadRequest
	if len(status) > 0 {
		s = status[0]
	}
	utils.ErrorJSON(w, err, s)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}
