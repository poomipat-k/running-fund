package cms

import "time"

type S3UploadResponse struct {
	ObjectKey string `json:"objectKey,omitempty"`
	FullPath  string `json:"fullPath,omitempty"`
}

type GetAdminDashboardDateConfigPreviewRequest struct {
	FromYear  int `json:"fromYear,omitempty"`
	FromMonth int `json:"fromMonth,omitempty"`
	FromDay   int `json:"fromDay,omitempty"`
	ToYear    int `json:"toYear,omitempty"`
	ToMonth   int `json:"toMonth,omitempty"`
	ToDay     int `json:"toDay,omitempty"`
	PageNo    int `json:"pageNo,omitempty"`
	PageSize  int `json:"pageSize,omitempty"`
}

type ReviewPeriod struct {
	Id       int        `json:"id,omitempty"`
	FromDate *time.Time `json:"fromDate,omitempty"`
	ToDate   *time.Time `json:"toDate,omitempty"`
}

type Image struct {
	Id              int     `json:"id,omitempty"`
	FullPath        string  `json:"fullPath,omitempty"`
	ObjectKey       string  `json:"objectKey,omitempty"`
	LinkTo          *string `json:"linkTo,omitempty"`
	WebsiteConfigId *int    `json:"websiteConfigId,omitempty"`
}

type AdminUpdateWebsiteConfigRequest struct {
	Landing     LandingConfig   `json:"landing,omitempty"`
	Dashboard   DashboardConfig `json:"dashboard,omitempty"`
	Faq         []FAQ           `json:"faq,omitempty"`
	HowToCreate []HowToCreate   `json:"howToCreate,omitempty"`
	Footer      FooterConfig    `json:"footer,omitempty"`
}

type LandingConfig struct {
	WebsiteConfigId int     `json:"websiteConfigId,omitempty"`
	Banner          []Image `json:"banner,omitempty"`
	Content         string  `json:"content,omitempty"`
}

type DashboardConfig struct {
	FromYear  int `json:"fromYear,omitempty"`
	FromMonth int `json:"fromMonth,omitempty"`
	FromDay   int `json:"fromDay,omitempty"`
	ToYear    int `json:"toYear,omitempty"`
	ToMonth   int `json:"toMonth,omitempty"`
	ToDay     int `json:"toDay,omitempty"`
}

type FAQ struct {
	Id       int    `json:"id,omitempty"`
	Question string `json:"question,omitempty"`
	Answer   string `json:"answer,omitempty"`
}

type HowToCreate struct {
	Id      int    `json:"id,omitempty"`
	Header  string `json:"header,omitempty"`
	Content string `json:"content,omitempty"`
}

type FooterConfig struct {
	Logo    []Image       `json:"logo,omitempty"`
	Contact FooterContact `json:"contact,omitempty"`
}

type FooterContact struct {
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
	FromHour    string `json:"fromHour,omitempty"`
	FromMinute  string `json:"fromMinute,omitempty"`
	ToHour      string `json:"toHour,omitempty"`
	ToMinute    string `json:"toMinute,omitempty"`
}

type FooterContactResponse struct {
	Email         string `json:"email,omitempty"`
	PhoneNumber   string `json:"phoneNumber,omitempty"`
	OperatingHour string `json:"operatingHour,omitempty"`
}

type FooterResponse struct {
	Id      int                   `json:"id,omitempty"`
	Logo    []Image               `json:"logo,omitempty"`
	Contact FooterContactResponse `json:"contact,omitempty"`
}

type CommonSuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type AdminDateConfigPreviewRow struct {
	ProjectCode      string    `json:"projectCode,omitempty"`
	ProjectCreatedAt time.Time `json:"projectCreatedAt,omitempty"`
	ProjectName      string    `json:"projectName,omitempty"`
	ProjectStatus    string    `json:"projectStatus,omitempty"`
	Count            int       `json:"count,omitempty"`
}

type UploadFileRequest struct {
	Name       string `json:"name,omitempty"`
	PathPrefix string `json:"pathPrefix,omitempty"`
}
