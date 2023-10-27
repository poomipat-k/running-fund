package projects

import "time"

type GetReviewerDashboardRequest struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}
