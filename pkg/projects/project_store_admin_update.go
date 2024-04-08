package projects

import (
	"database/sql"
	"fmt"
	"log/slog"
)

func (s *store) GetProjectStatusByProjectCode(projectCode string) (string, error) {
	var projectStatus string
	row := s.db.QueryRow(GetProjectStatusByProjectCodeSQL, projectCode)
	err := row.Scan(&projectStatus)
	switch err {
	case sql.ErrNoRows:
		slog.Error("GetReviewPeriod(): no row were returned!")
		return "", err
	case nil:
		return projectStatus, nil
	default:
		slog.Error(err.Error())
		return "", fmt.Errorf("GetReviewPeriod() unknown error")
	}
}
