package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/poomipat-k/running-fund/pkg/database"
	"github.com/poomipat-k/running-fund/pkg/server"
)

const webPort = "8080"

func main() {
	conn := database.ConnectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}
	defer conn.Close()

	app := server.Server{
		DB: conn,
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.Routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
