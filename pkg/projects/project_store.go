package projects

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/patrickmn/go-cache"
	s3Service "github.com/poomipat-k/running-fund/pkg/upload"
)

type store struct {
	db           *sql.DB
	c            *cache.Cache
	awsS3Service s3Service.S3Service
}

func NewStore(db *sql.DB, c *cache.Cache, s3service s3Service.S3Service) *store {
	return &store{
		db:           db,
		c:            c,
		awsS3Service: s3service,
	}
}

var monthMap = map[string]string{
	"January":   "JAN",
	"February":  "FEB",
	"March":     "MAR",
	"April":     "APR",
	"May":       "MAY",
	"June":      "JUN",
	"July":      "JUL",
	"August":    "AUG",
	"September": "SEP",
	"October":   "OCT",
	"November":  "NOV",
	"December":  "DEC",
}

const TIMEZONE = "Asia/Bangkok"

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

func (s *store) AddProject(payload AddProjectRequest, userId int, attachments []DetailsFiles) (int, error) {
	projectCode, err := s.generateProjectCode()
	if err != nil {
		return 0, err
	}
	// baseFilePrefix := getBasePrefix(userId, projectCode)
	_ = getBasePrefix(userId, projectCode)
	// start transaction
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return failAdd(err)
	}

	defer tx.Rollback()

	// now := time.Now()
	// Add address rows
	var addressId int
	err = tx.QueryRowContext(ctx, addAddressSQL, payload.General.Address.Address, payload.General.Address.PostcodeId).Scan(&addressId)
	if err != nil {
		return failAdd(err)
	}
	var projectCoordinatorAddressId int
	err = tx.QueryRowContext(
		ctx,
		addAddressSQL,
		payload.Contact.ProjectCoordinator.Address.Address,
		payload.Contact.ProjectCoordinator.Address.PostcodeId,
	).Scan(&projectCoordinatorAddressId)
	if err != nil {
		return failAdd(err)
	}
	// Add contact rows
	var projectHeadContactId int
	err = tx.QueryRowContext(
		ctx,
		addContactMainSQL,
		payload.Contact.ProjectHead.Prefix,
		payload.Contact.ProjectHead.FirstName,
		payload.Contact.ProjectHead.LastName,
		payload.Contact.ProjectHead.OrganizationPosition,
		payload.Contact.ProjectHead.EventPosition,
	).Scan(&projectHeadContactId)
	if err != nil {
		return failAdd(err)
	}

	var projectManagerContactId int
	err = tx.QueryRowContext(
		ctx,
		addContactMainSQL,
		payload.Contact.ProjectManager.Prefix,
		payload.Contact.ProjectManager.FirstName,
		payload.Contact.ProjectManager.LastName,
		payload.Contact.ProjectManager.OrganizationPosition,
		payload.Contact.ProjectManager.EventPosition,
	).Scan(&projectManagerContactId)
	if err != nil {
		return failAdd(err)
	}

	var projectCoordinatorContactId int
	err = tx.QueryRowContext(
		ctx,
		addContactFullSQL,
		payload.Contact.ProjectCoordinator.Prefix,
		payload.Contact.ProjectCoordinator.FirstName,
		payload.Contact.ProjectCoordinator.LastName,
		payload.Contact.ProjectCoordinator.OrganizationPosition,
		payload.Contact.ProjectCoordinator.EventPosition,
		projectCoordinatorAddressId,
		payload.Contact.ProjectCoordinator.Email,
		payload.Contact.ProjectCoordinator.LineId,
		payload.Contact.ProjectCoordinator.PhoneNumber,
	).Scan(&projectCoordinatorContactId)
	if err != nil {
		return failAdd(err)
	}

	var projectRaceDirectorContactId int
	if payload.Contact.RaceDirector.Who == "other" {
		err = tx.QueryRowContext(
			ctx,
			addContactOnlyRequiredParamSQL,
			payload.Contact.RaceDirector.Alternative.Prefix,
			payload.Contact.RaceDirector.Alternative.FirstName,
			payload.Contact.RaceDirector.Alternative.LastName,
		).Scan(&projectRaceDirectorContactId)
		if err != nil {
			return failAdd(err)
		}
	}

	// // upload files
	// for _, files := range attachments {
	// 	err = s.awsS3Service.UploadToS3(files.Files, fmt.Sprintf("%s/%s", baseFilePrefix, files.DirName))
	// 	if err != nil {
	// 		slog.Error("Failed to upload files to s3", "dirName", files.DirName, "error", err.Error())
	// 		return 0, err
	// 	}
	// }

	err = tx.Commit()
	if err != nil {
		return failAdd(err)
	}
	// commit
	return 1, nil
}

func failAdd(err error) (int, error) {
	return 0, fmt.Errorf("addProject: %w", err)
}

func (s *store) generateProjectCode() (string, error) {
	rawYear, rawMonth, day := getLocalYearMonthDay()
	year2digitsBud := (rawYear + 543) % 100
	month := monthMap[rawMonth.String()]
	cnt, found := s.c.Get(fmt.Sprintf("projectCode_%d_%s_%d", year2digitsBud, month, day))
	newCount := 0
	if found {
		newCount = cnt.(int) + 1
		s.c.Set(fmt.Sprintf("projectCode_%d_%s_%d", year2digitsBud, month, day), newCount, time.Duration(24*time.Hour))
	} else {
		// Query count from database
		row := s.db.QueryRow(countProjectCreatedToday)
		var count int
		err := row.Scan(&count)
		if err != nil {
			return "", err
		}
		newCount = count + 1
		s.c.Set(fmt.Sprintf("projectCode_%d_%s_%d", year2digitsBud, month, day), newCount, cache.NoExpiration)

	}
	projectCode := fmt.Sprintf("%s%d_%02d%02d", month, year2digitsBud, day, newCount)
	return projectCode, nil
}

func getLocalYearMonthDay() (int, time.Month, int) {
	//init the loc
	loc, err := time.LoadLocation(TIMEZONE)
	if err != nil {
		log.Fatal(err)
	}
	// set timezone,
	now := time.Now().In(loc)
	year, month, day := now.Date()
	return year, month, day
}

func getBasePrefix(userId int, projectCode string) string {
	return fmt.Sprintf("applicant/user_%d/%s", userId, projectCode)
}
