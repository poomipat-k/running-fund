package appMiddleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/poomipat-k/running-fund/pkg/utils"
)

func MyFirstMiddleWare(f http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Running before handler")
		f(w, r)
		log.Println("Running after handler")
	})
}

func IsReviewer(f http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Running before IsReviewer")
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")

		if len(authHeader) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
		} else {
			tokenString := authHeader[1]
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// validate the alg is what you expect:
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(os.Getenv("JWT_SECRET_KEY")), nil
			})
			if err != nil {
				log.Println("===ERROR 1 err:", err)
				utils.ErrorJSON(w, err, http.StatusForbidden)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if ok {
				fmt.Println(claims["userRole"], claims["userId"], claims["exp"], claims["iat"])
				f(w, r)
				log.Println("OK Running after IsReviewer")
			} else {
				utils.ErrorJSON(w, errors.New("invalid token2"), http.StatusForbidden)
				return
			}
		}
	})
}
