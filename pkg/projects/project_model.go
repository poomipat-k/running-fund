package projects

import (
	"mime/multipart"
	"time"
)

type ReviewDashboardRow struct {
	UserId           int        `json:"userId,omitempty"`
	ProjectId        int        `json:"projectId,omitempty"`
	ProjectCode      string     `json:"projectCode,omitempty"`
	ProjectCreatedAt *time.Time `json:"projectCreatedAt,omitempty"`
	ProjectName      string     `json:"projectName,omitempty"`
	ReviewId         int        `json:"reviewId,omitempty"`
	ReviewedAt       *time.Time `json:"reviewedAt,omitempty"`
}

type ReviewPeriod struct {
	Id       int        `json:"id,omitempty"`
	FromDate *time.Time `json:"fromDate,omitempty"`
	ToDate   *time.Time `json:"toDate,omitempty"`
}

type ProjectReviewDetailsResponse struct {
	UserId               int                `json:"userId,omitempty"`
	ProjectId            int                `json:"projectId,omitempty"`
	ProjectHistoryId     int                `json:"projectHistoryId,omitempty"`
	ProjectCode          string             `json:"projectCode,omitempty"`
	ProjectCreatedAt     *time.Time         `json:"projectCreatedAt,omitempty"`
	ProjectName          string             `json:"projectName,omitempty"`
	ProjectHeadPrefix    string             `json:"projectHeadPrefix,omitempty"`
	ProjectHeadFirstName string             `json:"projectHeadFirstName,omitempty"`
	ProjectHeadLastName  string             `json:"projectHeadLastName,omitempty"`
	FromDate             time.Time          `json:"fromDate,omitempty"`
	ToDate               time.Time          `json:"toDate,omitempty"`
	Address              string             `json:"address,omitempty"`
	ProvinceName         string             `json:"provinceName,omitempty"`
	DistrictName         string             `json:"districtName,omitempty"`
	SubdistrictName      string             `json:"subdistrictName,omitempty"`
	Distances            []DistanceAndFee   `json:"distances,omitempty"`
	ExpectedParticipants string             `json:"expectedParticipants,omitempty"`
	Collaborated         *bool              `json:"collaborated,omitempty"`
	ReviewId             int                `json:"reviewId,omitempty"`
	ReviewedAt           *time.Time         `json:"reviewedAt,omitempty"`
	IsInterestedPerson   *bool              `json:"isInterestedPerson,omitempty"`
	InterestedPersonType string             `json:"interestedPersonType,omitempty"`
	ReviewDetails        []ReviewDetails    `json:"reviewDetails,omitempty"`
	ReviewSummary        string             `json:"reviewSummary,omitempty"`
	ReviewerComment      string             `json:"reviewerComment,omitempty"`
	ReviewImprovement    *ReviewImprovement `json:"reviewImprovement,omitempty"`
}

type ProjectReviewDetailsRow struct {
	UserId               int                `json:"userId,omitempty"`
	ProjectId            int                `json:"projectId,omitempty"`
	ProjectHistoryId     int                `json:"projectHistoryId,omitempty"`
	ProjectCode          string             `json:"projectCode,omitempty"`
	ProjectCreatedAt     *time.Time         `json:"projectCreatedAt,omitempty"`
	ProjectName          string             `json:"projectName,omitempty"`
	ProjectHeadPrefix    string             `json:"projectHeadPrefix,omitempty"`
	ProjectHeadFirstName string             `json:"projectHeadFirstName,omitempty"`
	ProjectHeadLastName  string             `json:"projectHeadLastName,omitempty"`
	FromDate             time.Time          `json:"fromDate,omitempty"`
	ToDate               time.Time          `json:"toDate,omitempty"`
	Address              string             `json:"address,omitempty"`
	ProvinceName         string             `json:"provinceName,omitempty"`
	DistrictName         string             `json:"districtName,omitempty"`
	SubdistrictName      string             `json:"subdistrictName,omitempty"`
	DistanceType         string             `json:"distanceType,omitempty"`
	DistanceDynamic      bool               `json:"distanceDynamic,omitempty"`
	ExpectedParticipants string             `json:"expectedParticipants,omitempty"`
	Collaborated         *bool              `json:"collaborated,omitempty"`
	ReviewId             int                `json:"reviewId,omitempty"`
	ReviewedAt           *time.Time         `json:"reviewedAt,omitempty"`
	IsInterestedPerson   *bool              `json:"isInterestedPerson,omitempty"`
	InterestedPersonType string             `json:"interestedPersonType,omitempty"`
	ReviewDetails        []ReviewDetails    `json:"reviewDetails,omitempty"`
	ReviewSummary        string             `json:"reviewSummary,omitempty"`
	ReviewerComment      string             `json:"reviewerComment,omitempty"`
	ReviewImprovement    *ReviewImprovement `json:"reviewImprovement,omitempty"`
}

type ReviewImprovement struct {
	Benefit                  *bool `json:"benefit,omitempty"`
	ExperienceAndReliability *bool `json:"experienceAndReliability,omitempty"`
	FundAndOutput            *bool `json:"fundAndOutput,omitempty"`
	ProjectQuality           *bool `json:"projectQuality,omitempty"`
	ProjectStandard          *bool `json:"projectStandard,omitempty"`
	VisionAndImage           *bool `json:"visionAndImage,omitempty"`
}

type ProjectReviewCriteria struct {
	CriteriaId      int    `json:"reviewCriteriaId,omitempty"`
	CriteriaVersion int    `json:"criteriaVersion,omitempty"`
	GroupNumber     int    `json:"groupNumber,omitempty"`
	InGroupNumber   int    `json:"inGroupNumber,omitempty"`
	OrderNumber     int    `json:"orderNumber,omitempty"`
	DisplayText     string `json:"displayText,omitempty"`
}

type ReviewDetails struct {
	ReviewDetailsId     int `json:"reviewDetailsId,omitempty"`
	CriteriaVersion     int `json:"criteriaVersion,omitempty"`
	CriteriaOrderNumber int `json:"criteriaOrderNumber,omitempty"`
	Score               int `json:"score,omitempty"`
}

type ApplicantSelfScoreCriteria struct {
	Id              int    `json:"id,omitempty"`
	CriteriaVersion int    `json:"criteriaVersion,omitempty"`
	OrderNumber     int    `json:"orderNumber,omitempty"`
	Display         string `json:"display,omitempty"`
}

type Attachments struct {
	DirName         string
	ZipName         string
	InZipFilePrefix string
	Files           []*multipart.FileHeader
}

type ApplicantDashboardItem struct {
	ProjectId        int       `json:"projectId,omitempty"`
	ProjectCode      string    `json:"projectCode,omitempty"`
	ProjectCreatedAt time.Time `json:"projectCreatedAt,omitempty"`
	ProjectName      string    `json:"projectName,omitempty"`
	ProjectStatus    string    `json:"projectStatus,omitempty"`
	ProjectUpdatedAt time.Time `json:"projectUpdatedAt,omitempty"`
	AdminComment     string    `json:"adminComment,omitempty"`
}

type ApplicantDetailsData struct {
	ProjectCode   string     `json:"projectCode,omitempty"`
	UserId        int        `json:"userId,omitempty"`
	ProjectName   string     `json:"projectName,omitempty"`
	ProjectStatus string     `json:"projectStatus,omitempty"`
	ReviewId      *int       `json:"reviewId,omitempty"`
	ReviewedAt    *time.Time `json:"reviewedAt,omitempty"`
	SumScore      *int       `json:"sumScore,omitempty"`
}

type S3ObjectDetails struct {
	Key          string    `json:"key,omitempty"`
	LastModified time.Time `json:"lastModified,omitempty"`
}
