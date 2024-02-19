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
	Collaborated *bool                    `json:"collaborated,omitempty"`
	General      AddProjectGeneralDetails `json:"general"`
}

type AddProjectGeneralDetails struct {
	ProjectName string    `json:"projectName,omitempty"`
	EventDate   EventDate `json:"eventDate,omitempty"`
	Address     Address   `json:"address,omitempty"`
}

type EventDate struct {
	Year       int  `json:"year,omitempty"`
	Month      int  `json:"month,omitempty"`
	Day        int  `json:"day,omitempty"`
	FromHour   *int `json:"fromHour,omitempty"`
	FromMinute *int `json:"fromMinute,omitempty"`
	ToHour     *int `json:"toHour,omitempty"`
	ToMinute   *int `json:"toMinute,omitempty"`
}

type Address struct {
	Address       string `json:"address,omitempty"`
	ProvinceId    int    `json:"provinceId,omitempty"`
	DistrictId    int    `json:"districtId,omitempty"`
	SubdistrictId int    `json:"subdistrictId,omitempty"`
	PostcodeId    int    `json:"postcodeId,omitempty"`
}
