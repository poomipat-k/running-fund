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
	Details      Details                  `json:"details,omitempty"`
	Experience   Experience               `json:"experience,omitempty"`
	Fund         Fund                     `json:"fund,omitempty"`
}

type AddProjectFilesRequest struct {
	ProjectCode string `json:"projectCode,omitempty"`
	UserId      int    `json:"userId,omitempty"`
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
	VIPFee         *float64         `json:"vipFee,omitempty"`
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
	ProjectHead        ProjectHead         `json:"projectHead,omitempty"`
	ProjectManager     ProjectManager      `json:"projectManager,omitempty"`
	ProjectCoordinator ProjectCoordinator  `json:"projectCoordinator,omitempty"`
	RaceDirector       RaceDirector        `json:"raceDirector,omitempty"`
	Organization       ContactOrganization `json:"organization,omitempty"`
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

type ProjectCoordinator struct {
	Prefix               string  `json:"prefix,omitempty"`
	FirstName            string  `json:"firstName,omitempty"`
	LastName             string  `json:"lastName,omitempty"`
	OrganizationPosition string  `json:"organizationPosition,omitempty"`
	EventPosition        string  `json:"eventPosition,omitempty"`
	Address              Address `json:"address,omitempty"`
	Email                string  `json:"email,omitempty"`
	LineId               string  `json:"lineId,omitempty"`
	PhoneNumber          string  `json:"phoneNumber,omitempty"`
}

type RaceDirector struct {
	Who         string                  `json:"who,omitempty"`
	Alternative RaceDirectorAlternative `json:"alternative,omitempty"`
}

type RaceDirectorAlternative struct {
	Prefix    string `json:"prefix,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

type ContactOrganization struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

type Details struct {
	Background string         `json:"background,omitempty"`
	Objective  string         `json:"objective,omitempty"`
	Marketing  Marketing      `json:"marketing,omitempty"`
	Score      map[string]int `json:"score,omitempty"`
	Safety     Safety         `json:"Safety,omitempty"`
	Route      Route          `json:"route,omitempty"`
	Judge      Judge          `json:"judge,omitempty"`
	Support    Support        `json:"support,omitempty"`
	Feedback   string         `json:"feedback,omitempty"`
}

type Marketing struct {
	Online  Online  `json:"online,omitempty"`
	Offline Offline `json:"offline,omitempty"`
}

type Online struct {
	Available OnlineAvailable `json:"available,omitempty"`
	HowTo     OnlineHowTo     `json:"howTo,omitempty"`
}

type OnlineAvailable struct {
	Facebook   bool `json:"facebook,omitempty"`
	Website    bool `json:"website,omitempty"`
	OnlinePage bool `json:"onlinePage,omitempty"`
	Other      bool `json:"other,omitempty"`
}

type OnlineHowTo struct {
	Facebook   string `json:"facebook,omitempty"`
	Website    string `json:"website,omitempty"`
	OnlinePage string `json:"onlinePage,omitempty"`
	Other      string `json:"other,omitempty"`
}

type Offline struct {
	Available OfflineAvailable `json:"available,omitempty"`
	Addition  string           `json:"addition,omitempty"`
}

type OfflineAvailable struct {
	PR            bool `json:"pr,omitempty"`
	LocalOfficial bool `json:"localOfficial,omitempty"`
	Booth         bool `json:"booth,omitempty"`
	Billboard     bool `json:"billboard,omitempty"`
	TV            bool `json:"tv,omitempty"`
	Other         bool `json:"other,omitempty"`
}

type Safety struct {
	Ready    SafetyReady `json:"ready,omitempty"`
	AEDCount int         `json:"aedCount,omitempty"`
	Addition string      `json:"addition,omitempty"`
}

type SafetyReady struct {
	RunnerInformation bool `json:"runnerInformation,omitempty"`
	HealthDecider     bool `json:"healthDecider,omitempty"`
	Ambulance         bool `json:"ambulance,omitempty"`
	FirstAid          bool `json:"firstAid,omitempty"`
	AED               bool `json:"aed,omitempty"`
	Insurance         bool `json:"insurance,omitempty"`
	Other             bool `json:"other,omitempty"`
}

type Route struct {
	Measurement       RouteMeasurement  `json:"measurement,omitempty"`
	Tool              string            `json:"tool,omitempty"`
	TrafficManagement TrafficManagement `json:"trafficManagement,omitempty"`
}

type RouteMeasurement struct {
	AthleticsAssociation bool `json:"athleticsAssociation,omitempty"`
	CalibratedBicycle    bool `json:"calibratedBicycle,omitempty"`
	SelfMeasurement      bool `json:"selfMeasurement,omitempty"`
}

type TrafficManagement struct {
	AskPermission bool `json:"askPermission,omitempty"`
	HasSupporter  bool `json:"hasSupporter,omitempty"`
	RoadClosure   bool `json:"roadClosure,omitempty"`
	Signs         bool `json:"signs,omitempty"`
	Lighting      bool `json:"lighting,omitempty"`
}

type Judge struct {
	Type      string `json:"type,omitempty"`
	OtherType string `json:"otherType,omitempty"`
}

type Support struct {
	Addition     string       `json:"addition,omitempty"`
	Organization Organization `json:"organization,omitempty"`
}

type Organization struct {
	ProvincialAdministration bool `json:"provincialAdministration,omitempty"`
	Safety                   bool `json:"safety,omitempty"`
	Health                   bool `json:"health,omitempty"`
	Volunteer                bool `json:"volunteer,omitempty"`
	Community                bool `json:"community,omitempty"`
	Other                    bool `json:"other,omitempty"`
}

type Experience struct {
	ThisSeries  ThisSeries  `json:"thisSeries,omitempty"`
	OtherSeries OtherSeries `json:"otherSeries,omitempty"`
}

type ThisSeries struct {
	FirstTime *bool             `json:"firstTime,omitempty"`
	History   ThisSeriesHistory `json:"history,omitempty"`
}

type ThisSeriesHistory struct {
	Completed1    HistoryCompleted `json:"completed1,omitempty"`
	Completed2    HistoryCompleted `json:"completed2,omitempty"`
	Completed3    HistoryCompleted `json:"completed3,omitempty"`
	OrdinalNumber int              `json:"ordinalNumber,omitempty"`
	Year          int              `json:"year,omitempty"`
	Month         int              `json:"month,omitempty"`
	Day           int              `json:"day,omitempty"`
}

type HistoryCompleted struct {
	Year        int    `json:"year,omitempty"`
	Name        string `json:"name,omitempty"`
	Participant int    `json:"participant,omitempty"`
}

type OtherSeries struct {
	DoneBefore *bool              `json:"doneBefore,omitempty"`
	History    OtherSeriesHistory `json:"history,omitempty"`
}

type OtherSeriesHistory struct {
	Completed1 HistoryCompleted `json:"completed1,omitempty"`
	Completed2 HistoryCompleted `json:"completed2,omitempty"`
	Completed3 HistoryCompleted `json:"completed3,omitempty"`
}

type Fund struct {
	Budget  Budget      `json:"budget,omitempty"`
	Request FundRequest `json:"request,omitempty"`
}

type Budget struct {
	Total               int    `json:"total,omitempty"`
	SupportOrganization string `json:"supportOrganization,omitempty"`
}

type FundRequest struct {
	Type    FundRequestType    `json:"type,omitempty"`
	Details FundRequestDetails `json:"details,omitempty"`
}

type FundRequestType struct {
	Fund    bool `json:"fund,omitempty"`
	BIB     bool `json:"bib,omitempty"`
	Pr      bool `json:"pr,omitempty"`
	Seminar bool `json:"seminar,omitempty"`
	Other   bool `json:"other,omitempty"`
}

type FundRequestDetails struct {
	FundAmount int    `json:"fundAmount,omitempty"`
	BibAmount  int    `json:"bibAmount,omitempty"`
	Seminar    string `json:"seminar,omitempty"`
	Other      string `json:"other,omitempty"`
}

// End sub-types for AddProjectRequest

type ListFilesRequest struct {
	Prefix    string `json:"prefix,omitempty"`
	CreatedBy int    `json:"createdBy,omitempty"`
}

type ProjectReviewer struct {
	ReviewerId int `json:"reviewerId,omitempty"`
}

type AdminUpdateProjectRequest struct {
	ProjectStatusPrimary   string     `json:"projectStatusPrimary,omitempty"`
	ProjectStatusSecondary string     `json:"projectStatusSecondary,omitempty"`
	AdminScore             *int       `json:"adminScore,omitempty"`
	FundApprovedAmount     *int64     `json:"fundApprovedAmount,omitempty"`
	AdminComment           *string    `json:"adminComment,omitempty"`
	AdminApprovedAt        *time.Time `json:"adminApprovedAt,omitempty"`
}
