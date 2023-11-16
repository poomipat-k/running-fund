package projects

import (
	"time"
)

type ReviewDashboardRow struct {
	ProjectId        int        `json:"project_id,omitempty"`
	ProjectCode      string     `json:"project_code,omitempty"`
	ProjectCreatedAt *time.Time `json:"project_created_at,omitempty"`
	ProjectName      string     `json:"project_name,omitempty"`
	ReviewId         int        `json:"review_id,omitempty"`
	ReviewedAt       *time.Time `json:"reviewed_at,omitempty"`
	DownloadLink     string     `json:"download_link,omitempty"`
}

type ReviewPeriod struct {
	Id       int        `json:"id,omitempty"`
	FromDate *time.Time `json:"from_date,omitempty"`
	ToDate   *time.Time `json:"to_date,omitempty"`
}

type ProjectReviewDetails struct {
	ProjectId        int        `json:"project_id,omitempty"`
	ProjectCode      string     `json:"project_code,omitempty"`
	ProjectCreatedAt *time.Time `json:"project_created_at,omitempty"`
	ProjectName      string     `json:"project_name,omitempty"`
	ReviewId         int        `json:"review_id,omitempty"`
	ReviewedAt       *time.Time `json:"reviewed_at,omitempty"`
}

type ProjectReviewCriteria struct {
	CriteriaId      int    `json:"review_criteria_id,omitempty"`
	CriteriaVersion int    `json:"criteria_version,omitempty"`
	GroupNumber     int    `json:"group_number,omitempty"`
	InGroupNumber   int    `json:"in_group_number,omitempty"`
	OrderNumber     int    `json:"order_number,omitempty"`
	DisplayText     string `json:"display_text,omitempty"`
}
