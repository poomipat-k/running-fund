package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	_ "github.com/lib/pq"
)

const (
	host     = "compose_postgres"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "running_fund"
)

func main() {
	r := chi.NewRouter()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("===1")
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("===2")
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome"))
	})

	r.Route("/api", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			sqlStatement := `
			INSERT INTO users (age, email, first_name, last_name)
			VALUES (30, 'a@gmail.com', 'P', 'K')`

			_, err = db.Exec(sqlStatement)
			if err != nil {
				panic(err)
			}

			w.Write([]byte("Hello API Home"))
		})

		r.Post("/signup", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("From Signup"))
		})
	})

	http.ListenAndServe(":4000", r)
}
