package projects

import "time"

type ApplicantDashboardItem struct {
	ProjectId        int       `json:"projectId,omitempty"`
	ProjectCode      string    `json:"projectCode,omitempty"`
	ProjectCreatedAt time.Time `json:"projectCreatedAt,omitempty"`
	ProjectName      string    `json:"projectName,omitempty"`
	ProjectStatus    string    `json:"projectStatus,omitempty"`
	ProjectUpdatedAt time.Time `json:"projectUpdatedAt,omitempty"`
	AdminComment     string    `json:"adminComment,omitempty"`
}
