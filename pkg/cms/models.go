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

type Banner struct {
	Id        int    `json:"id,omitempty"`
	LinkTo    string `json:"linkTo,omitempty"`
	ObjectKey string `json:"objectKey,omitempty"`
	FullPath  string `json:"fullPath,omitempty"`
}

type AdminUpdateWebsiteConfigRequest struct {
	Landing   LandingConfig   `json:"landing,omitempty"`
	Dashboard DashboardConfig `json:"dashboard,omitempty"`
}

type LandingConfig struct {
	Banner  []Banner `json:"banner,omitempty"`
	Content string   `json:"content,omitempty"`
}

type DashboardConfig struct {
	FromYear  int `json:"fromYear,omitempty"`
	FromMonth int `json:"fromMonth,omitempty"`
	FromDay   int `json:"fromDay,omitempty"`
	ToYear    int `json:"toYear,omitempty"`
	ToMonth   int `json:"toMonth,omitempty"`
	ToDay     int `json:"toDay,omitempty"`
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
