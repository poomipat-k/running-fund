package users

import (
	"database/sql"
	"log"
	"log/slog"
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
	rows, err := s.db.Query(getReviewersSQL, "reviewer")
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

func (s *store) GetReviewerById(userId int) (User, error) {
	var user User
	row := s.db.QueryRow(getReviewerByIdSQL, userId)
	err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email)
	switch err {
	case sql.ErrNoRows:
		slog.Error("GetReviewerById() no row were returned!")
		return User{}, err
	case nil:
		return user, nil
	default:
		slog.Error(err.Error())
		return User{}, err
	}
}
