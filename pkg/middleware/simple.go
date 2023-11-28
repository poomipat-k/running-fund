package customMiddleware

import (
	"log"
	"net/http"
)

func MyFirstMiddleWare(f http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Running before handler")
		f(w, r)
		log.Println("Running after handler")
	})
}
