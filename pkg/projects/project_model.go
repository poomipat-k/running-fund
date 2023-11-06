package projects

import (
	"time"
)

type ReviewDashboardRow struct {
	ProjectId        int        `json:"project_id"`
	ProjectCode      string     `json:"project_code"`
	ProjectCreatedAt *time.Time `json:"project_created_at"`
	ProjectName      string     `json:"project_name"`
	ReviewId         int        `json:"review_id,omitempty"`
	ReviewedAt       *time.Time `json:"reviewed_at,omitempty"`
	DownloadLink     string     `json:"download_link,omitempty"`
}

type ReviewPeriod struct {
	Id       int        `json:"id"`
	FromDate *time.Time `json:"from_date"`
	ToDate   *time.Time `json:"to_date"`
}

type ProjectReviewDetails struct {
	ProjectId   int    `json:"project_id"`
	ProjectCode string `json:"project_code"`
	ProjectName string `json:"project_name"`
}
