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
	General      AddProjectGeneralDetails `json:"general,omitempty"`
	Contact      Contact                  `json:"contact,omitempty"`
}

// Sub-types for AddProjectRequest
type AddProjectGeneralDetails struct {
	ProjectName          string       `json:"projectName,omitempty"`
	EventDate            EventDate    `json:"eventDate,omitempty"`
	Address              Address      `json:"address,omitempty"`
	StartPoint           string       `json:"startPoint,omitempty"`
	FinishPoint          string       `json:"finishPoint,omitempty"`
	EventDetails         EventDetails `json:"eventDetails,omitempty"`
	ExpectedParticipants string       `json:"expectedParticipants,omitempty"`
	HasOrganizer         *bool        `json:"hasOrganizer,omitempty"`
	OrganizerName        string       `json:"organizerName,omitempty"`
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

type EventDetails struct {
	Category       Category         `json:"category,omitempty"`
	DistanceAndFee []DistanceAndFee `json:"distanceAndFee,omitempty"`
	VIP            *bool            `json:"vip,omitempty"`
}

type Category struct {
	Available Available `json:"available,omitempty"`
	OtherType string    `json:"otherType,omitempty"`
}

type Available struct {
	Other        bool `json:"other,omitempty"`
	RoadRace     bool `json:"roadRace,omitempty"`
	TrailRunning bool `json:"trailRunning,omitempty"`
}

type DistanceAndFee struct {
	Checked bool     `json:"checked,omitempty"`
	Type    string   `json:"type,omitempty"`
	Fee     *float64 `json:"fee,omitempty"`
	Dynamic *bool    `json:"dynamic,omitempty"`
}

type Contact struct {
	ProjectHead    ProjectHead    `json:"projectHead,omitempty"`
	ProjectManager ProjectManager `json:"projectManager,omitempty"`
}

type ProjectHead struct {
	Prefix               string `json:"prefix,omitempty"`
	FirstName            string `json:"firstName,omitempty"`
	LastName             string `json:"lastName,omitempty"`
	OrganizationPosition string `json:"organizationPosition,omitempty"`
	EventPosition        string `json:"eventPosition,omitempty"`
}

type ProjectManager struct {
	Prefix               string `json:"prefix,omitempty"`
	FirstName            string `json:"firstName,omitempty"`
	LastName             string `json:"lastName,omitempty"`
	OrganizationPosition string `json:"organizationPosition,omitempty"`
	EventPosition        string `json:"eventPosition,omitempty"`
}

// End sub-types for AddProjectRequest
