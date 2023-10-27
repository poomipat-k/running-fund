package projects

import "time"

type Project struct {
	Id             int       `json:"id"`
	ProjectCode    string    `json:"project_code"`
	ProjectVersion int       `json:"project_version"`
	CreatedAt      time.Time `json:"created_at"`
}
