package projects

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	myCsv "github.com/poomipat-k/running-fund/pkg/csv-app"
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

func (s *store) GetAdminRequestDashboard(
	fromDate, toDate time.Time,
	orderBy string,
	limit, offset int,
	projectCode, projectName, projectStatus *string,
) ([]AdminRequestDashboardRow, error) {
	queryStmt, values := prepareAdminDashboardQuery("request", fromDate, toDate, orderBy, limit, offset, projectCode, projectName, projectStatus)
	rows, err := s.db.Query(queryStmt, values...)
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
			&row.Count,
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

func (s *store) GetAdminStartedDashboard(
	fromDate, toDate time.Time,
	orderBy string,
	limit, offset int,
	projectCode, projectName, projectStatus *string,
) ([]AdminRequestDashboardRow, error) {
	queryStmt, values := prepareAdminDashboardQuery("started", fromDate, toDate, orderBy, limit, offset, projectCode, projectName, projectStatus)
	rows, err := s.db.Query(queryStmt, values...)
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
			&row.Count,
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

func (s *store) GetAdminSummary(fromDate, toDate time.Time) ([]AdminSummaryData, error) {
	rows, err := s.db.Query(getAdminSummarySQL, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []AdminSummaryData
	for rows.Next() {
		var row AdminSummaryData
		err := rows.Scan(
			&row.Status,
			&row.Count,
			&row.FundSum,
		)
		if err != nil {
			return nil, err
		}

		if row.FundSum == nil {
			row.FundSum = newInt64(0)
		}

		data = append(data, row)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *store) GenerateAdminReport() error {
	items := [][]string{
		{"id", "name", "price"},
		{"1", "Book", "18"},
		{"2", "Pencil", "3"},
		{"3", "Rubber", "2"},
	}

	buffer, err := myCsv.GenCsvBuffer(items)
	if err != nil {
		panic(err)
	}
	err = s.awsS3Service.DoUploadFileToS3(buffer, fmt.Sprintf("%s/%s", "csv", "myTestCsv.csv"))
	if err != nil {
		return err
	}
	return nil
}

func newInt64(val int64) *int64 {
	v := val
	return &v
}

func prepareAdminDashboardQuery(
	dashboardType string,
	fromDate, toDate time.Time,
	orderBy string,
	limit, offset int,
	projectCode, projectName, projectStatus *string,
) (string, []any) {
	curPlaceholder := 3
	var where []string
	if dashboardType == "request" {
		where = []string{"project.created_at >= $1 AND project.created_at <= $2 AND (project_history.status != 'Start' AND project_history.status != 'Completed')"}
	} else {
		where = []string{"project.created_at >= $1 AND project.created_at <= $2 AND (project_history.status = 'Start' OR project_history.status = 'Completed')"}
	}
	values := []any{fromDate, toDate}
	if projectCode != nil {
		where = append(where, fmt.Sprintf("AND project_history.project_code = $%d", curPlaceholder))
		values = append(values, *projectCode)
		curPlaceholder++
	}
	if projectName != nil {
		strContainStmt := "AND project_history.project_name LIKE '%' || $" + strconv.Itoa(curPlaceholder) + " || '%'"
		where = append(where, strContainStmt)
		values = append(values, *projectName)
		curPlaceholder++
	}
	if projectStatus != nil {
		where = append(where, fmt.Sprintf("AND project_history.status = $%d", curPlaceholder))
		values = append(values, *projectStatus)
		curPlaceholder++
	}
	whereStmt := strings.Join(where, " ")
	// orderBy must be safe string
	orderLimitOffsetStmt := fmt.Sprintf("ORDER BY %s LIMIT $%d OFFSET $%d", orderBy, curPlaceholder, curPlaceholder+1)
	values = append(values, limit, offset)
	countStmt := fmt.Sprintf(`
	(
		SELECT COUNT(*) FROM project 
		INNER JOIN project_history ON project.project_history_id = project_history.id
		WHERE %s
	) as count`, whereStmt)

	getAdminRequestDashboardSQL := fmt.Sprintf(`
	SELECT
project.project_code as project_code,
project.created_at as created_at,
project_history.project_name as project_name,
project_history.status as project_status,
project_history.updated_at as updated_at,
project_history.admin_comment,
(
SELECT ROUND(AVG(sum_score), 2)
	FROM (
		SELECT
		review.project_history_id as project_history_id,
		review.id as review_id,
		SUM(review_details.score) as sum_score
		FROM review
		INNER JOIN review_details ON review.id = review_details.review_id
		WHERE project_history_id = project.project_history_id
		GROUP BY  review.project_history_id, review.id
		)
) as avg_score,
%s
FROM project 
INNER JOIN project_history ON project.project_history_id = project_history.id
WHERE `, countStmt)

	queryStmt := strings.Join([]string{getAdminRequestDashboardSQL, whereStmt, orderLimitOffsetStmt}, " ") + ";"
	return queryStmt, values
}
