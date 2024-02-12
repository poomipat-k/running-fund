package main

import (
	"fmt"
	"log"
	"net/http"

	_ "time/tzdata"

	"github.com/poomipat-k/running-fund/pkg/database"
	"github.com/poomipat-k/running-fund/pkg/server"
)

const webPort = "8080"

func main() {

	db := database.ConnectToDB()
	if db == nil {
		log.Panic("Can't connect to Postgres!")
	}
	defer db.Close()

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
