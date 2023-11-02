package users

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

func (s *store) GetReviewers() ([]User, error) {
	rows, err := s.db.Query(`
	SELECT id, first_name, last_name, email, user_role, created_at
	FROM users WHERE user_role = $1
	`, "reviewer")
	if err != nil {
		log.Println("Error on Query: ", err)
		return nil, err
	}
	defer rows.Close()

	var data []User
	for rows.Next() {
		var row User
		err = rows.Scan(&row.Id, &row.FirstName, &row.LastName, &row.Email, &row.UserRole, &row.CreatedAt)
		if err != nil {
			log.Println("Error on Scan: ", err)
			return nil, err
		}
		data = append(data, row)
	}
	err = rows.Err()
	if err != nil {
		log.Println("Error on rows.Err: ", err)
		return nil, err
	}
	return data, nil
}
