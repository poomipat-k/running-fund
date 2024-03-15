package projects

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/patrickmn/go-cache"
	s3Service "github.com/poomipat-k/running-fund/pkg/s3-service"
)

const applicantCriteriaCachePrefix = "applicant_criteria"
const reviewerCriteriaCachePrefix = "reviewer_criteria"

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

func (s *store) HasPermissionToAddAdditionalFiles(userId int, projectCode string) bool {
	var projectId int
	var projectStatus string
	row := s.db.QueryRow(hasRightToAddAdditionalFilesSQL, userId, projectCode)
	err := row.Scan(&projectId, &projectStatus)
	switch err {
	case sql.ErrNoRows:
		slog.Error("GetReviewPeriod(): no row were returned!")
		return false
	case nil:
		return projectStatus == "RevisedRequired"
	default:
		slog.Error(err.Error())
		return false
	}
}

// Get project [fromDate, toDate)
func (s *store) GetReviewerDashboard(reviewerId int, fromDate, toDate time.Time) ([]ReviewDashboardRow, error) {
	rows, err := s.db.Query(getReviewerDashboardSQL, reviewerId, fromDate, toDate)
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
		err = rows.Scan(&row.UserId, &row.ProjectId, &row.ProjectCode, &row.ProjectCreatedAt, &row.ProjectName, &reviewId, &reviewedAt)
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

		data = append(data, row)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *store) GetApplicantProjectDetails(userId int, projectCode string) ([]ApplicantDetailsData, error) {
	rows, err := s.db.Query(getApplicantProjectDetailsSQL, userId, projectCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []ApplicantDetailsData
	for rows.Next() {
		var row ApplicantDetailsData
		err = rows.Scan(&row.ProjectCode, &row.UserId, &row.ProjectName, &row.ProjectStatus, &row.AdminScore,
			&row.FundApprovedAmount, &row.AdminComment, &row.ReviewId, &row.ReviewerId, &row.ReviewedAt, &row.SumScore)
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

func (s *store) GetReviewerProjectDetails(userId int, projectCode string) (ProjectReviewDetailsResponse, error) {
	rows, err := s.db.Query(getReviewerProjectDetailsSQL, userId, projectCode)
	if err != nil {
		return ProjectReviewDetailsResponse{}, err
	}
	defer rows.Close()

	var data []ProjectReviewDetailsRow
	var distances []DistanceAndFee
	for rows.Next() {
		var row ProjectReviewDetailsRow
		// Nullable columns
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
		err = rows.Scan(
			&row.UserId,
			&row.ProjectId,
			&row.ProjectHistoryId,
			&row.ProjectCode,
			&row.ProjectCreatedAt,
			&row.ProjectName,
			&row.ProjectHeadPrefix,
			&row.ProjectHeadFirstName,
			&row.ProjectHeadLastName,
			&row.FromDate,
			&row.ToDate,
			&row.Address,
			&row.ProvinceName,
			&row.DistrictName,
			&row.SubdistrictName,
			&row.DistanceType,
			&row.DistanceDynamic,
			&row.ExpectedParticipants,
			&row.Collaborated,
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
		if err != nil {
			slog.Error(err.Error())
			return ProjectReviewDetailsResponse{}, err
		}

		if reviewId.Valid {
			row.ReviewId = int(reviewId.Int64)
		}
		if reviewedAt.Valid {
			row.ReviewedAt = &reviewedAt.Time
		}
		if isInterestedPerson.Valid {
			row.IsInterestedPerson = &isInterestedPerson.Bool
		}
		if interestedPersonType.Valid {
			row.InterestedPersonType = interestedPersonType.String
		}
		if reviewerComment.Valid {
			row.ReviewerComment = reviewerComment.String
		}

		if row.ReviewId > 0 {
			rd, err := s.GetReviewDetailsByReviewId(row.ReviewId)
			if err != nil {
				slog.Error(err.Error())
				return ProjectReviewDetailsResponse{}, err
			}
			row.ReviewDetails = rd
		}
		if reviewSummary.Valid {
			row.ReviewSummary = reviewSummary.String
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
			row.ReviewImprovement = imp
		}

		distance := DistanceAndFee{
			Type:    row.DistanceType,
			Dynamic: &row.DistanceDynamic,
		}
		distances = append(distances, distance)
		data = append(data, row)
	}
	err = rows.Err()
	if err != nil {
		return ProjectReviewDetailsResponse{}, err
	}
	var body ProjectReviewDetailsResponse
	if len(data) > 0 {
		f := data[0]
		body.UserId = f.UserId
		body.ProjectId = f.ProjectId
		body.ProjectHistoryId = f.ProjectHistoryId
		body.ProjectCode = f.ProjectCode
		body.ProjectCreatedAt = f.ProjectCreatedAt
		body.ProjectName = f.ProjectName
		body.ProjectHeadPrefix = f.ProjectHeadPrefix
		body.ProjectHeadFirstName = f.ProjectHeadFirstName
		body.ProjectHeadLastName = f.ProjectHeadLastName
		body.FromDate = f.FromDate
		body.ToDate = f.ToDate
		body.Address = f.Address
		body.ProvinceName = f.ProvinceName
		body.DistrictName = f.DistrictName
		body.SubdistrictName = f.SubdistrictName
		body.Distances = distances
		body.ExpectedParticipants = f.ExpectedParticipants
		body.Collaborated = f.Collaborated
		body.ReviewId = f.ReviewId
		body.ReviewedAt = f.ReviewedAt
		body.IsInterestedPerson = f.IsInterestedPerson
		body.InterestedPersonType = f.InterestedPersonType
		body.ReviewDetails = f.ReviewDetails
		body.ReviewSummary = f.ReviewSummary
		body.ReviewerComment = f.ReviewerComment
		body.ReviewImprovement = f.ReviewImprovement
	}

	return body, nil
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
	// check cache
	cacheKey := fmt.Sprintf("%s_%d", applicantCriteriaCachePrefix, criteriaVersion)
	raw, found := s.c.Get(cacheKey)
	if found {
		cachedData, ok := raw.([]ApplicantSelfScoreCriteria)
		if ok {
			return cachedData, nil
		}
	}

	// Fetch data from the db
	rows, err := s.db.Query(getApplicantCriteriaSQL, criteriaVersion)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []ApplicantSelfScoreCriteria
	for rows.Next() {
		var row ApplicantSelfScoreCriteria

		err := rows.Scan(&row.Id, &row.CriteriaVersion, &row.OrderNumber, &row.Display)
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
	if len(data) > 0 {
		s.c.Set(cacheKey, data, cache.NoExpiration)
	}
	return data, nil
}

func (s *store) GetProjectCriteria(criteriaVersion int) ([]ProjectReviewCriteria, error) {
	if criteriaVersion == 0 {
		criteriaVersion = 1
	}
	// check cache
	cacheKey := fmt.Sprintf("%s_%d", reviewerCriteriaCachePrefix, criteriaVersion)
	raw, found := s.c.Get(cacheKey)
	if found {
		cachedData, ok := raw.([]ProjectReviewCriteria)
		if ok {
			return cachedData, nil
		}
	}

	// Check db
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
	if len(data) > 0 {
		s.c.Set(cacheKey, data, cache.NoExpiration)
	}
	return data, nil
}

func (s *store) GetAllProjectDashboardByApplicantId(applicantId int) ([]ApplicantDashboardItem, error) {
	// getAllProjectDashboardByApplicantIdSQL
	rows, err := s.db.Query(getAllProjectDashboardByApplicantIdSQL, applicantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []ApplicantDashboardItem
	for rows.Next() {
		var row ApplicantDashboardItem
		// // Nullable columns
		var adminComment sql.NullString
		err = rows.Scan(&row.ProjectId, &row.ProjectCode, &row.ProjectCreatedAt, &row.ProjectName, &row.ProjectStatus, &row.ProjectUpdatedAt, &adminComment)
		if err != nil {
			return nil, err
		}
		// Check Nullable columns
		if adminComment.Valid {
			row.AdminComment = adminComment.String
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
