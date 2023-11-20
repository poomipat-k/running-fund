package projects

import (
	"time"
)

type ReviewDashboardRow struct {
	ProjectId        int        `json:"projectId,omitempty"`
	ProjectCode      string     `json:"projectCode,omitempty"`
	ProjectCreatedAt *time.Time `json:"projectCreatedAt,omitempty"`
	ProjectName      string     `json:"projectName,omitempty"`
	ReviewId         int        `json:"reviewId,omitempty"`
	ReviewedAt       *time.Time `json:"reviewedAt,omitempty"`
	DownloadLink     string     `json:"downloadLink,omitempty"`
}

type ReviewPeriod struct {
	Id       int        `json:"id,omitempty"`
	FromDate *time.Time `json:"fromDate,omitempty"`
	ToDate   *time.Time `json:"toDate,omitempty"`
}

type ProjectReviewDetails struct {
	ProjectId            int                `json:"projectId,omitempty"`
	ProjectCode          string             `json:"projectCode,omitempty"`
	ProjectCreatedAt     *time.Time         `json:"projectCreatedAt,omitempty"`
	ProjectName          string             `json:"projectName,omitempty"`
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

type ProjectReviewCriteriaMinimal struct {
	CriteriaVersion int `json:"criteriaVersion,omitempty"`
	OrderNumber     int `json:"orderNumber,omitempty"`
}

type ReviewDetails struct {
	ReviewDetailsId     int `json:"reviewDetailsId,omitempty"`
	CriteriaVersion     int `json:"criteriaVersion,omitempty"`
	CriteriaOrderNumber int `json:"criteriaOrderNumber,omitempty"`
	Score               int `json:"score,omitempty"`
}
