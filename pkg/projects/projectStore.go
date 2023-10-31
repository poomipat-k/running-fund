package projects

import (
	"database/sql"
	"log"
	"time"
)

type store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *store {
	return &store{
		db: db,
	}
}

func (s *store) GetReviewerDashboard(userId int, fromDate, toDate time.Time) ([]ReviewDashboardRow, error) {
	rows, err := s.db.Query(`
	SELECT project.id as project_id, project.project_code, 
	project.created_at, project_history.project_name, 
	review.id as review_id, review.created_at as reviewed_at,
	project_history.download_link
	FROM project
	INNER JOIN project_history
	ON project.project_history_id = project_history.id
	LEFT JOIN review
	ON project.project_history_id = review.project_history_id AND user_id = $1
	WHERE project.created_at >= $2
	AND project.created_at <= $3`, userId, fromDate, toDate)
	if err != nil {
		log.Println("Error on Query: ", err)
		return nil, err
	}
	defer rows.Close()

	var data []ReviewDashboardRow
	for rows.Next() {
		var row ReviewDashboardRow
		// Nullable columns
		var reviewId sql.NullInt64
		var reviewedAt sql.NullTime
		var dowloadLink sql.NullString
		err = rows.Scan(&row.ProjectId, &row.ProjectCode, &row.ProjectCreatedAt, &row.ProjectName, &reviewId, &reviewedAt, &dowloadLink)
		if err != nil {
			log.Println("Error on Scan: ", err)
			return nil, err
		}
		// Check Nullable columns
		if reviewId.Valid {
			row.ReviewId = int(reviewId.Int64)
		}
		if reviewedAt.Valid {
			row.ReviewedAt = &reviewedAt.Time
		}
		if dowloadLink.Valid {
			row.DownloadLink = dowloadLink.String
		}

		data = append(data, row)
	}
	// get any error cncountered during iteration
	err = rows.Err()
	if err != nil {
		log.Println("Error on rows.Err: ", err)
		return nil, err
	}
	return data, nil
}
