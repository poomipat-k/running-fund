package projects

import "time"

type GetReviewerDashboardRequest struct {
	FromDate time.Time `json:"fromDate"`
	ToDate   time.Time `json:"toDate"`
}
