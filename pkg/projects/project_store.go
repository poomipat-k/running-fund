package projects

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
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
		slog.Error("GetReviewPeriod(): no row were returned!")
		return ReviewPeriod{}, err
	case nil:
		return period, nil
	default:
		slog.Error(err.Error())
		return ReviewPeriod{}, fmt.Errorf("GetReviewPeriod() unknown error")
	}
}

// Get project [fromDate, toDate)
func (s *store) GetReviewerDashboard(userId int, fromDate, toDate time.Time) ([]ReviewDashboardRow, error) {
	rows, err := s.db.Query(getReviewerDashboardSQL, userId, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []ReviewDashboardRow
	for rows.Next() {
		var row ReviewDashboardRow
		// Nullable columns
		var reviewId sql.NullInt64
		var reviewedAt sql.NullTime
		var download sql.NullString
		err = rows.Scan(&row.ProjectId, &row.ProjectCode, &row.ProjectCreatedAt, &row.ProjectName, &reviewId, &reviewedAt, &download)
		if err != nil {
			return nil, err
		}
		// Check Nullable columns
		if reviewId.Valid {
			row.ReviewId = int(reviewId.Int64)
		}
		if reviewedAt.Valid {
			row.ReviewedAt = &reviewedAt.Time
		}
		if download.Valid {
			row.DownloadLink = download.String
		}

		data = append(data, row)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *store) GetReviewerProjectDetails(userId int, projectCode string) (ProjectReviewDetails, error) {
	var details ProjectReviewDetails
	row := s.db.QueryRow(getReviewerProjectDetailsSQL, userId, projectCode)
	// Nullable
	var reviewId sql.NullInt64
	var reviewedAt sql.NullTime
	var isInterestedPerson sql.NullBool
	var interestedPersonType sql.NullString
	var reviewSummary sql.NullString
	var reviewerComment sql.NullString
	var benefit sql.NullBool
	var exp sql.NullBool
	var fund sql.NullBool
	var projectQuality sql.NullBool
	var projectStandard sql.NullBool
	var vision sql.NullBool
	err := row.Scan(
		&details.ProjectId,
		&details.ProjectHistoryId,
		&details.ProjectCode,
		&details.ProjectCreatedAt,
		&details.ProjectName,
		&reviewId,
		&reviewedAt,
		&isInterestedPerson,
		&interestedPersonType,
		&reviewSummary,
		&reviewerComment,
		&benefit,
		&exp,
		&fund,
		&projectQuality,
		&projectStandard,
		&vision,
	)
	if reviewId.Valid {
		details.ReviewId = int(reviewId.Int64)
	}
	if reviewedAt.Valid {
		details.ReviewedAt = &reviewedAt.Time
	}
	if isInterestedPerson.Valid {
		details.IsInterestedPerson = &isInterestedPerson.Bool
	}
	if interestedPersonType.Valid {
		details.InterestedPersonType = interestedPersonType.String
	}
	if reviewerComment.Valid {
		details.ReviewerComment = reviewerComment.String
	}

	if details.ReviewId > 0 {
		rd, err := s.GetReviewDetailsByReviewId(details.ReviewId)
		if err != nil {
			slog.Error(err.Error())
			return ProjectReviewDetails{}, err
		}
		details.ReviewDetails = rd
	}
	if reviewSummary.Valid {
		details.ReviewSummary = reviewSummary.String
	}
	imp := &ReviewImprovement{}
	hasImprovement := false
	if benefit.Valid {
		hasImprovement = true
		imp.Benefit = &benefit.Bool
	}
	if exp.Valid {
		hasImprovement = true
		imp.ExperienceAndReliability = &exp.Bool
	}
	if fund.Valid {
		hasImprovement = true
		imp.FundAndOutput = &fund.Bool
	}
	if projectQuality.Valid {
		hasImprovement = true
		imp.ProjectQuality = &projectQuality.Bool
	}
	if projectStandard.Valid {
		hasImprovement = true
		imp.ProjectStandard = &projectStandard.Bool
	}
	if vision.Valid {
		hasImprovement = true
		imp.VisionAndImage = &vision.Bool
	}
	if hasImprovement {
		details.ReviewImprovement = imp
	}

	switch err {
	case sql.ErrNoRows:
		slog.Error("GetReviewerProjectDetails() no row were returned!")
		return ProjectReviewDetails{}, err
	case nil:
		return details, nil
	default:
		slog.Error(err.Error())
		return ProjectReviewDetails{}, err
	}
}

func (s *store) GetReviewDetailsByReviewId(reviewId int) ([]ReviewDetails, error) {
	rows, err := s.db.Query(getReviewDetailsByReviewIdSQL, reviewId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var data []ReviewDetails
	for rows.Next() {
		var row ReviewDetails
		err := rows.Scan(&row.ReviewDetailsId, &row.CriteriaVersion, &row.CriteriaOrderNumber, &row.Score)
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

func (s *store) GetApplicantCriteria(criteriaVersion int) ([]ApplicantSelfScoreCriteria, error) {
	if criteriaVersion == 0 {
		criteriaVersion = 1
	}
	rows, err := s.db.Query(getApplicantCriteriaSQL, criteriaVersion)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []ApplicantSelfScoreCriteria
	for rows.Next() {
		var row ApplicantSelfScoreCriteria

		err := rows.Scan(&row.CriteriaVersion, &row.OrderNumber, &row.Display)
		if err != nil {
			return nil, err
		}
		data = append(data, row)
	}
	// get any error occur during iteration
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("criteria version not found")
	}
	return data, nil
}

func (s *store) GetProjectCriteria(criteriaVersion int) ([]ProjectReviewCriteria, error) {
	if criteriaVersion == 0 {
		criteriaVersion = 1
	}
	rows, err := s.db.Query(getProjectCriteriaSQL, criteriaVersion)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []ProjectReviewCriteria
	for rows.Next() {
		var row ProjectReviewCriteria
		var groupNumber sql.NullInt64
		var inGroupNumber sql.NullInt64
		var displayText sql.NullString

		err := rows.Scan(&row.CriteriaVersion, &row.OrderNumber, &groupNumber, &inGroupNumber, &displayText)
		if err != nil {
			return nil, err
		}

		// Check Nullable columns
		if groupNumber.Valid {
			row.GroupNumber = int(groupNumber.Int64)
		}
		if inGroupNumber.Valid {
			row.InGroupNumber = int(inGroupNumber.Int64)
		}
		if displayText.Valid {
			row.DisplayText = displayText.String
		}

		data = append(data, row)
	}
	// get any error occur during iteration
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("criteria version not found")
	}
	return data, nil
}
