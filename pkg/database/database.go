package database

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var counter int64

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func ConnectToDB() *sql.DB {
	dsn := os.Getenv("DSN")
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready...")
			counter++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		if counter > 20 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for half a second...")
		time.Sleep(500 * time.Millisecond)
	}
}
