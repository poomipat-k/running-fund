package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/poomipat-k/running-fund/pkg/database"
)

const webPort = "8080"

type Config struct {
}

func main() {
	conn := database.ConnectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}
	defer conn.Close()

	app := Config{}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

	// r := chi.NewRouter()
	// r.Use(middleware.Logger)
	// r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Welcome"))
	// })

	// r.Route("/api", func(r chi.Router) {
	// 	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 		w.Write([]byte("API landing page"))
	// 	})

	// 	r.Post("/signup", func(w http.ResponseWriter, r *http.Request) {
	// 		sqlStatement := `
	// 		INSERT INTO users (age, email, first_name, last_name)
	// 		VALUES ($1, $2, $3, $4)
	// 		RETURNING id`
	// 		id := 0
	// 		err = db.QueryRow(sqlStatement, 30, "jon@calhoun.io", "Jonathan", "Calhoun").Scan(&id)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		fmt.Println("New record ID is:", id)
	// 		w.WriteHeader(http.StatusCreated)
	// 		w.Write([]byte("From Signup"))
	// 	})
	// })

	// err := http.ListenAndServe(":8080", r)

	// if err != nil {
	// 	log.Println("There was an error listening on port :8080", err)
	// }
}
