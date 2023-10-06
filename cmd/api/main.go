package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/poomipat-k/running-fund/pkg/database"

	_ "github.com/lib/pq"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	db, err := database.GetDbConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome"))
	})

	r.Route("/api", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("API landing page"))
		})

		r.Post("/signup", func(w http.ResponseWriter, r *http.Request) {
			sqlStatement := `
			INSERT INTO users (age, email, first_name, last_name)
			VALUES ($1, $2, $3, $4)
			RETURNING id`
			id := 0
			err = db.QueryRow(sqlStatement, 30, "jon@calhoun.io", "Jonathan", "Calhoun").Scan(&id)
			if err != nil {
				panic(err)
			}
			fmt.Println("New record ID is:", id)
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("From Signup"))
		})
	})

	http.ListenAndServe(":4000", r)
}
