package users

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"
)

const dbTimeout = time.Second * 5

type store struct {
	db           *sql.DB
	emailService EmailService
}

func NewStore(db *sql.DB, es EmailService) *store {
	return &store{
		db:           db,
		emailService: es,
	}
}

func (s *store) GetUserById(userId int) (User, error) {
	var user User
	row := s.db.QueryRow(getUserByIdSQL, userId)
	err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.UserRole, &user.Activated)
	switch err {
	case sql.ErrNoRows:
		slog.Error("GetUserById() no row were returned!")
		return User{}, err
	case nil:
		return user, nil
	default:
		slog.Error(err.Error())
		return User{}, err
	}
}

func (s *store) GetUserFullNameById(userId int) (UserFullName, error) {
	var user UserFullName
	row := s.db.QueryRow(getUserFullNameByIdSQL, userId)
	err := row.Scan(&user.Id, &user.FirstName, &user.LastName)
	switch err {
	case sql.ErrNoRows:
		slog.Error("GetUserFullNameById() no row were returned!")
		return UserFullName{}, err
	case nil:
		return user, nil
	default:
		slog.Error(err.Error())
		return UserFullName{}, err
	}
}

func (s *store) GetUserByEmail(email string) (User, error) {
	var user User
	row := s.db.QueryRow(getUserByEmailSQL, email)
	err := row.Scan(&user.Id, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.UserRole, &user.Activated, &user.ActivatedBefore, &user.CreatedAt)
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

func (s *store) AddUser(user User, toBeDeletedUserId int) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return failAddUser(err, "")
	}
	defer tx.Rollback()

	// Email already used but we want to replace
	if toBeDeletedUserId > 0 {
		_, _, err := s.DeleteUserById(toBeDeletedUserId, ctx, tx)
		if err != nil {
			return failAddUser(err, "")
		}
	}

	var userId int
	err = tx.QueryRowContext(ctx, addUserSQL, user.Email, user.Password, user.FirstName, user.LastName, "applicant", false, user.ActivateCode).Scan(&userId)
	if err != nil {
		return failAddUser(err, "dbQuery")
	}

	activateLink := fmt.Sprintf("http://%s/signup/activate/%s", os.Getenv("UI_URL"), user.ActivateCode)
	mail := s.emailService.BuildSignUpConfirmationEmail(user.Email, activateLink)
	err = s.emailService.SendEmail(mail)
	if err != nil {
		slog.Error("Signup: failed to send account activation email", "error", err.Error())
		return failAddUser(fmt.Errorf("ไม่สามารถส่งอีเมลไปยังที่อยู่อีเมลนี้ได้ โปรดตรวจสอบที่อยู่อีเมล"), "email")
	}
	slog.Info("Sign up confirmation email sent to", "email", user.Email)

	err = tx.Commit()
	if err != nil {
		return failAddUser(err, "commit")
	}
	return userId, "", nil
}

func (s *store) DeleteUserById(id int, ctx context.Context, tx *sql.Tx) (int, string, error) {
	var deletedId int
	if tx != nil {
		err := tx.QueryRowContext(ctx, DeleteUserByIdSQL, id).Scan(&deletedId)
		if err != nil {
			return failAddUser(err, "")
		}
		return deletedId, "", nil
	}

	err := s.db.QueryRowContext(ctx, DeleteUserByIdSQL, id).Scan(&deletedId)
	if err != nil {
		return failAddUser(err, "")
	}
	return deletedId, "", nil
}

func (s *store) ActivateUser(activateCode string) (int64, error) {
	result, err := s.db.Exec(activateEmailSQL, activateCode)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *store) ForgotPasswordAction(resetPasswordCode string, email string, resetPasswordLink string) (int64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return failForgotPasswordAction(err)
	}
	defer tx.Rollback()

	result, err := tx.Exec(forgotPasswordSQL, resetPasswordCode, email)
	if err != nil {
		return 0, err
	}
	// Send email
	mail := s.emailService.BuildResetPasswordEmail(email, resetPasswordLink)
	err = s.emailService.SendEmail(mail)
	if err != nil {
		return 0, err
	}
	slog.Info("reset password email sent to", "email", email)

	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *store) ResetPassword(resetPasswordCode string, newPassword string) (int64, error) {
	result, err := s.db.Exec(resetPasswordSQL, resetPasswordCode, newPassword)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()

}

func failAddUser(err error, name string) (int, string, error) {
	return 0, name, err
}

func failForgotPasswordAction(err error) (int64, error) {
	return 0, fmt.Errorf("forgotPasswordAction: %w", err)
}
