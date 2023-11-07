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

func (s *store) GetReviewPeriod() (ReviewPeriod, error) {
	var period ReviewPeriod
	row := s.db.QueryRow(getReviewPeriodSQL)
	err := row.Scan(&period.Id, &period.FromDate, &period.ToDate)
	switch err {
	case sql.ErrNoRows:
		log.Println("No row were returned!")
		return ReviewPeriod{}, err
	case nil:
		return period, nil
	default:
		panic(err)
	}
}

// Get project [fromDate, toDate)
func (s *store) GetReviewerDashboard(userId int, fromDate, toDate time.Time) ([]ReviewDashboardRow, error) {
	rows, err := s.db.Query(getReviewerDashboardSQL, userId, fromDate, toDate)
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

func (s *store) GetReviewerProejctDetails(userId int, projectCode string) (ProjectReviewDetails, error) {
	var details ProjectReviewDetails
	row := s.db.QueryRow(getReviewerProejctDetailsSQL, userId, projectCode)
	// Nullable
	var reviewId sql.NullInt64
	var reviewedAt sql.NullTime
	var reviewDetailsId sql.NullInt64
	err := row.Scan(&details.ProjectId, &details.ProjectCode, &details.ProjectName, &reviewId, &reviewedAt, &reviewDetailsId)
	if reviewId.Valid {
		details.ReviewId = int(reviewId.Int64)
	}
	if reviewedAt.Valid {
		details.ReviewedAt = &reviewedAt.Time
	}
	if reviewDetailsId.Valid {
		details.ReviewDetailsId = int(reviewDetailsId.Int64)
	}
	switch err {
	case sql.ErrNoRows:
		log.Println("No row were returned!")
		return ProjectReviewDetails{}, err
	case nil:
		return details, nil
	default:
		panic(err)
	}
}
