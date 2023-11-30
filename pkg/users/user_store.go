package users

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

const dbTimeout = time.Second * 5

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
		slog.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var data []User
	for rows.Next() {
		var row User
		err = rows.Scan(&row.Id, &row.FirstName, &row.LastName, &row.Email, &row.UserRole, &row.CreatedAt)
		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}
		data = append(data, row)
	}
	err = rows.Err()
	if err != nil {
		slog.Error(err.Error())
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

func (s *store) GetUserByEmail(email string) (User, error) {
	var user User
	row := s.db.QueryRow(getUserByEmailSQL, email)
	err := row.Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName, &user.UserRole, &user.Activated, &user.ActivatedBefore, &user.CreatedAt)
	switch err {
	case sql.ErrNoRows:
		slog.Error("GetUserByEmail() no row were returned!")
		return User{}, err
	case nil:
		return user, nil
	default:
		slog.Error(err.Error())
		return User{}, err
	}
}

func (s *store) AddUser(user User, toBeDeletedUserId int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return failAddUser(err)
	}
	defer tx.Rollback()

	// Email already used but we want to replace
	if toBeDeletedUserId > 0 {
		var deletedId int
		err = tx.QueryRowContext(ctx, DeleteUserByIdSQL, toBeDeletedUserId).Scan(&deletedId)
		if err != nil {
			return failAddUser(err)
		}
	}

	var userId int
	err = tx.QueryRowContext(ctx, addUserSQL, user.Email, user.Password, user.FirstName, user.LastName, "applicant", false).Scan(&userId)
	if err != nil {
		return failAddUser(err)
	}
	err = tx.Commit()
	if err != nil {
		return failAddUser(err)
	}
	return userId, nil
}

func failAddUser(err error) (int, error) {
	return 0, fmt.Errorf("addUser: %w", err)
}
