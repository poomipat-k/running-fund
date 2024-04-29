package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/poomipat-k/running-fund/pkg/database"
	"github.com/poomipat-k/running-fund/pkg/server"
	"github.com/pressly/goose/v3"
)

const webPort = "8080"

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	db := database.ConnectToDB()
	if db == nil {
		log.Panic("Can't connect to Postgres!")
	}
	defer db.Close()

	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Println("error to SetDialect")
		panic(err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		log.Println("error to goose.Up")
		panic(err)
	}

	app := server.Server{}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.Routes(db),
	}

	log.Println("Ready on port ", webPort)

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
