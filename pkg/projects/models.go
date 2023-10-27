package projects

import "time"

type Project struct {
	Id             int       `json:"id"`
	ProjectCode    string    `json:"project_code"`
	ProjectName    string    `json:"project_name"`
	ProjectVersion int       `json:"project_version"`
	CreatedAt      time.Time `json:"created_at"`
}
