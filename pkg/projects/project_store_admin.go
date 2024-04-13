package projects

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"mime/multipart"
	"time"
)

func (s *store) GetProjectStatusByProjectCode(projectCode string) (AdminUpdateParam, error) {
	var payload AdminUpdateParam
	row := s.db.QueryRow(getProjectForAdminUpdateByProjectCodeSQL, projectCode)
	err := row.Scan(
		&payload.CreatedBy,
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

func (s *store) UpdateProjectByAdmin(payload AdminUpdateParam, userId int, projectCode string, additionFiles []*multipart.FileHeader) error {
	// start transaction
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var id int
	err = tx.QueryRowContext(
		ctx,
		updateProjectByAdminSQL,
		payload.ProjectHistoryId,
		payload.ProjectStatus,
		payload.AdminScore,
		payload.FundApprovedAmount,
		payload.AdminComment,
		payload.AdminApprovedAt,
		payload.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return err
	}

	// upload additionFiles
	objectPrefix := fmt.Sprintf("applicant/user_%d/%s/addition", userId, projectCode)
	err = s.awsS3Service.UploadFilesToS3(additionFiles, objectPrefix)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	// committed
	slog.Info("success update a project", "projectHistoryId", payload.ProjectHistoryId)

	return nil
}

func (s *store) GetAdminRequestDashboard(fromDate, toDate time.Time, orderBy string, limit, offset int) ([]AdminRequestDashboardRow, error) {
	orderByStmt := orderBy
	rows, err := s.db.Query(getAdminRequestDashboardSQL, fromDate, toDate, orderByStmt, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var data []AdminRequestDashboardRow
	for rows.Next() {
		var row AdminRequestDashboardRow
		err := rows.Scan(
			&row.ProjectCode,
			&row.ProjectCreatedAt,
			&row.ProjectName,
			&row.ProjectStatus,
			&row.ProjectUpdatedAt,
			&row.AdminComment,
			&row.AvgScore,
		)
		if err != nil {
			return nil, err
		}

		data = append(data, row)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return data, nil
}
