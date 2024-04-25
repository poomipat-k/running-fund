package projects

import (
	"mime/multipart"
	"time"
)

var thirtyDaysMonth = map[int]int{
	4:  30,
	6:  30,
	9:  30,
	11: 30,
}

var PROJECT_STATUS = map[string]int{
	"Reviewing":   1,
	"Reviewed":    2,
	"Revise":      3,
	"NotApproved": 4,
	"Approved":    5,
	"Start":       6,
	"Completed":   7,
}

var sortByWhiteList = map[string]bool{
	"project_history.project_name": true,
	"project_history.created_at":   true,
	"project_history.updated_at":   true,
	"project_history.status":       true,
	"POSITION(project_history.status::text IN 'Reviewing,Reviewed,Revise,NotApproved,Approved')": true,
	"POSITION(project_history.status::text IN 'Start,Completed')":                                true,
}

var PRIMARY_STATUS = map[string]bool{
	"CurrentBeforeApprove": true,
	"Approved":             true,
	"NotApproved":          true,
}

func validateAddProjectPayload(
	payload AddProjectRequest,
	collaborateFiles []*multipart.FileHeader,
	criteria []ApplicantSelfScoreCriteria,
	marketingFiles, routeFiles, eventMapFiles, eventDetailsFiles []*multipart.FileHeader) error {
	if payload.Collaborated == nil {
		return &CollaboratedRequiredError{}
	}
	if *payload.Collaborated && len(collaborateFiles) == 0 {
		return &CollaboratedFilesRequiredError{}
	}

	if err := validateGeneral(payload); err != nil {
		return err
	}
	if err := validateContact(payload); err != nil {
		return err
	}
	if err := validateDetails(payload, criteria); err != nil {
		return err
	}
	if err := validateExperience(payload); err != nil {
		return err
	}
	if err := validateFund(payload); err != nil {
		return err
	}
	if err := validateAttachment(marketingFiles, routeFiles, eventMapFiles, eventDetailsFiles); err != nil {
		return err
	}

	return nil
}

func validateAdminUpdateProjectPayload(payload AdminUpdateProjectRequest) (string, error) {
	if payload.ProjectStatusPrimary == "" {
		return "projectStatusPrimary", &ProjectStatusPrimaryRequiredError{}
	}
	if !PRIMARY_STATUS[payload.ProjectStatusPrimary] {
		return "projectStatusPrimary", &ProjectStatusPrimaryInvalidError{}
	}
	if payload.ProjectStatusSecondary == "" {
		return "projectStatusSecondary", &ProjectStatusSecondaryRequiredError{}
	}
	if PROJECT_STATUS[payload.ProjectStatusSecondary] == 0 {
		return "projectStatusSecondary", &ProjectStatusSecondaryInvalidError{}
	}
	if payload.AdminScore != nil {
		if *payload.AdminScore < 0 || *payload.AdminScore > 100 {
			return "adminScore", &AdminScoreOutOfRangeError{}
		}
	}
	if payload.FundApprovedAmount != nil && *payload.FundApprovedAmount < 0 {
		return "fundApprovedAmount", &FundApprovedAmountNegativeError{}
	}
	return "", nil
}

func validateGetAdminDashboardPayload(payload GetAdminDashboardRequest) (string, error) {
	if payload.FromYear < minDashboardYear {
		return "fromYear", &FromYearRequiredError{}
	}
	if payload.FromMonth == 0 {
		return "fromMonth", &MonthRequiredError{}
	}
	if payload.FromMonth < 1 || payload.FromMonth > 12 {
		return "fromMonth", &MonthOutOfBoundError{}
	}

	if payload.ToYear < minDashboardYear {
		return "toYear", &ToYearRequiredError{}
	}
	if payload.ToMonth == 0 {
		return "toMonth", &MonthRequiredError{}
	}
	if payload.ToMonth < 1 || payload.ToMonth > 12 {
		return "toMonth", &MonthOutOfBoundError{}
	}
	loc, err := getTimeLocation()
	if err != nil {
		return "timeLocation", nil
	}
	fromDate := time.Date(payload.FromYear, time.Month(payload.FromMonth), payload.FromDay, 0, 0, 0, 0, loc)
	toDate := time.Date(payload.ToYear, time.Month(payload.ToMonth), payload.ToDay, 23, 59, 59, 999999999, loc)
	if fromDate.After(toDate) {
		return "fromDate", &FromDateExceedToDateError{}
	}
	if payload.PageNo <= 0 {
		return "pageNo", &PageNoInvalidError{}
	}
	if payload.PageSize < 1 {
		return "pageSize", &PageSizeInvalidError{}
	}
	if len(payload.SortBy) == 0 {
		return "sortBy", &SortByRequiredError{}
	}

	for i := 0; i < len(payload.SortBy); i++ {
		if !sortByWhiteList[payload.SortBy[i]] {
			return "sortBy", &SortByInvalidError{}
		}
	}
	return "", nil
}

func validateGetAdminSummaryRequestPayload(payload GetAdminSummaryRequest) (string, error) {
	if payload.FromYear < minDashboardYear {
		return "fromYear", &FromYearRequiredError{}
	}
	if payload.FromMonth == 0 {
		return "fromMonth", &MonthRequiredError{}
	}
	if payload.FromMonth < 1 || payload.FromMonth > 12 {
		return "fromMonth", &MonthOutOfBoundError{}
	}

	if payload.ToYear < minDashboardYear {
		return "toYear", &ToYearRequiredError{}
	}
	if payload.ToMonth == 0 {
		return "toMonth", &MonthRequiredError{}
	}
	if payload.ToMonth < 1 || payload.ToMonth > 12 {
		return "toMonth", &MonthOutOfBoundError{}
	}
	loc, err := getTimeLocation()
	if err != nil {
		return "timeLocation", nil
	}
	fromDate := time.Date(payload.FromYear, time.Month(payload.FromMonth), payload.FromDay, 0, 0, 0, 0, loc)
	toDate := time.Date(payload.ToYear, time.Month(payload.ToMonth), payload.ToDay, 23, 59, 59, 999999999, loc)
	if fromDate.After(toDate) {
		return "fromDate", &FromDateExceedToDateError{}
	}

	return "", nil
}
