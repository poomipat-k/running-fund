package review

import (
	"database/sql"
	"log"
)

type store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *store {
	return &store{
		db: db,
	}
}

func (s *store) AddReview() {
	log.Println("Review Store")
}
