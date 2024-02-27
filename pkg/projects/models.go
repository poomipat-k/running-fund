package projects

import "time"

type ApplicantDashboardItem struct {
	ProjectCode      string    `json:"project_code,omitempty"`
	ProjectCreatedAt time.Time `json:"project_created_at,omitempty"`
	ProjectName      string    `json:"project_name,omitempty"`
	ProjectStatus    string    `json:"project_status,omitempty"`
	ProjectUpdatedAt time.Time `json:"project_updated_at,omitempty"`
	AdminComment     string    `json:"admin_comment,omitempty"`
	ReviewerComment  string    `json:"reviewer_comment,omitempty"`
}
