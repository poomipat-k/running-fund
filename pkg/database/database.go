package database

import (
	"database/sql"
	"fmt"
)

const (
	host     = "compose_postgres"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "running_fund"
)

func GetDbConnection() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
