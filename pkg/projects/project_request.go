package projects

import "time"

type GetReviewerDashboardRequest struct {
	FromDate time.Time `json:"fromDate"`
	ToDate   time.Time `json:"toDate"`
}

type AddReviewRequest struct {
	ProjectHistoryId int    `json:"projectHistoryId,omitempty"`
	Comment          string `json:"comment"`
	Ip               Ip     `json:"ip"`
	Review           Review `json:"review"`
}

type Ip struct {
	IsInterestedPerson   *bool  `json:"isInterestedPerson,omitempty"`
	InterestedPersonType string `json:"interestedPersonType,omitempty"`
}

type Review struct {
	ReviewSummary string            `json:"reviewSummary,omitempty"`
	Improvement   ReviewImprovement `json:"improvement,omitempty"`
	Scores        map[string]int    `json:"scores,omitempty"`
}

type AddProjectRequest struct {
	Collaborated bool                     `json:"collaborated"`
	General      AddProjectGeneralDetails `json:"general"`
}

type AddProjectGeneralDetails struct {
	ProjectName          string `json:"projectName"`
	StartPoint           string `json:"startPoint"`
	FinishPoint          string `json:"finishPoint"`
	ExpectedParticipants int    `json:"expectedParticipants"`
	HasOrganizer         bool   `json:"hasOrganizer"`
	OrganizerName        string `json:"organizerName,omitempty"`
}
