package projects

import (
	"database/sql"
	"fmt"
	"log/slog"
)

func (s *store) GetProjectStatusByProjectCode(projectCode string) (AdminUpdateParam, error) {
	var payload AdminUpdateParam
	row := s.db.QueryRow(GetProjectForAdminUpdateByProjectCodeSQL, projectCode)
	err := row.Scan(
		&payload.ProjectHistoryId,
		&payload.ProjectStatus,
		&payload.AdminScore,
		&payload.FundApprovedAmount,
		&payload.AdminComment,
		&payload.AdminApprovedAt,
		&payload.UpdatedAt,
	)
	switch err {
	case sql.ErrNoRows:
		slog.Error("GetProjectStatusByProjectCode(): no row were returned!")
		return AdminUpdateParam{}, err
	case nil:
		return payload, nil
	default:
		slog.Error(err.Error())
		return AdminUpdateParam{}, fmt.Errorf("GetProjectStatusByProjectCode() unknown error")
	}
}
