package mw

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/poomipat-k/running-fund/pkg/utils"
)

func IsLoggedIn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("===[IsLoggedIn]")

		at, err := getAccessToken(r)
		if err != nil {
			log.Println("===IsLoggedIn err 1")
			utils.ErrorJSON(w, err, http.StatusForbidden)
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

func IsReviewer(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getAccessToken(r)
		if err != nil {
			log.Println("===err 1")
			utils.ErrorJSON(w, err, http.StatusForbidden)
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
