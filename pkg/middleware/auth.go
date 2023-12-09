package mw

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/poomipat-k/running-fund/pkg/utils"
)

const accessExpireDurationMinute = 5

func IsLoggedIn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("===[IsLoggedIn]")

		at, err := getAccessToken(r)
		if err != nil {
			log.Println("===IsLoggedIn err 1")
			next(w, r)
			return
		}

		claims, ok := at.Claims.(jwt.MapClaims)
		if ok {
			userId := fmt.Sprintf("%v", claims["userId"])
			userRole := fmt.Sprintf("%v", claims["userRole"])
			r.Header.Set("userId", userId)
			r.Header.Set("userRole", userRole)
			next(w, r)
		} else {
			log.Println("===IsLoggedIn err 2")
			utils.ErrorJSON(w, errors.New("corrupt token"), http.StatusForbidden)
			return

		}
	})
}

// func TryRefreshAccessToken(next http.HandlerFunc) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Already has a valid access token
// 		if r.Header.Get("userId") != "" {
// 			next(w, r)
// 			return
// 		}
// 		// Try to refresh access token
// 		log.Println("====Attempt to refresh token ====")
// 		token, err := getRefreshToken(r)
// 		if err != nil {
// 			log.Println("===Refresh err 1")
// 			utils.ErrorJSON(w, err, http.StatusForbidden)
// 			return
// 		}

// 		claims, ok := token.Claims.(jwt.MapClaims)
// 		if ok {
// 			userId := fmt.Sprintf("%v", claims["userId"])
// 			userRole := fmt.Sprintf("%v", claims["userRole"])
// 			r.Header.Set("userId", userId)
// 			r.Header.Set("userRole", userRole)

// 			accessExpiredAtUnix := time.Now().Add(accessExpireDurationMinute * time.Second).Unix()
// 			uid, err := strconv.Atoi(userId)
// 			if err != nil {
// 				log.Println("===Refresh err 1.3")
// 				utils.ErrorJSON(w, err, http.StatusForbidden)
// 				return
// 			}
// 			accessToken, err := generateAccessToken(uid, userRole, accessExpiredAtUnix)
// 			if err != nil {
// 				log.Println("===Refresh err 1.2")
// 				utils.ErrorJSON(w, err, http.StatusForbidden)
// 				return
// 			}
// 			newAccessTokenCookie := http.Cookie{
// 				Name:     "authToken",
// 				Value:    accessToken,
// 				HttpOnly: true,
// 				Secure:   true,
// 				Path:     "/api",
// 				Expires:  time.Unix(accessExpiredAtUnix, 0),
// 			}
// 			http.SetCookie(w, &newAccessTokenCookie)

// 			log.Println("===[Token Refreshed] !! ==")
// 			next(w, r)
// 		} else {
// 			utils.ErrorJSON(w, errors.New("corrupt token"), http.StatusForbidden)
// 			return
// 		}
// 	})
// }

func IsReviewer(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getAccessToken(r)
		if err != nil {
			log.Println("===err 1")
			next(w, r)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if ok {
			userId := fmt.Sprintf("%v", claims["userId"])
			userRole := fmt.Sprintf("%v", claims["userRole"])

			if userRole != "reviewer" {
				utils.ErrorJSON(w, errors.New("permission denied"), http.StatusForbidden)
				return
			}
			r.Header.Set("userId", userId)
			r.Header.Set("userRole", userRole)

			next(w, r)
			log.Println("OK Running after IsReviewer")
		} else {
			utils.ErrorJSON(w, errors.New("corrupt token"), http.StatusForbidden)
			return

		}
	})
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

func getAccessToken(r *http.Request) (*jwt.Token, error) {
	// Cookie
	cookie, err := r.Cookie("authToken")
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		// validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_ACCESS_TOKEN_SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
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
