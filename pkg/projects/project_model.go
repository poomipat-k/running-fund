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
	ProjectId            int             `json:"projectId,omitempty"`
	ProjectCode          string          `json:"projectCode,omitempty"`
	ProjectCreatedAt     *time.Time      `json:"projectCreatedAt,omitempty"`
	ProjectName          string          `json:"projectName,omitempty"`
	ReviewId             int             `json:"reviewId,omitempty"`
	ReviewedAt           *time.Time      `json:"reviewedAt,omitempty"`
	IsInterestedPerson   bool            `json:"isInterestedPerson,omitempty"`
	InterestedPersonType string          `json:"interestedPersonType,omitempty"`
	ReviewDetails        []ReviewDetails `json:"reviewDetails,omitempty"`
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
