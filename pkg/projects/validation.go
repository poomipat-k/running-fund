package projects

import (
	"mime/multipart"
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

func validateGetAdminDashboardRequestPayload(payload GetAdminDashboardRequest) (string, error) {
	if payload.FromYear < minDashboardYear {
		return "fromYear", &FromYearRequiredError{}
	}
	if payload.ToYear < minDashboardYear {
		return "toYear", &ToYearRequiredError{}
	}
	if payload.FromYear > payload.ToYear {
		return "fromYear", &FromYearExceedToYearError{}
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
	return "", nil
}

func validateGetAdminSummaryRequestPayload(payload GetAdminSummaryRequest) (string, error) {
	if payload.FromYear < minDashboardYear {
		return "fromYear", &FromYearRequiredError{}
	}
	if payload.ToYear < minDashboardYear {
		return "toYear", &ToYearRequiredError{}
	}
	if payload.FromYear > payload.ToYear {
		return "fromYear", &FromYearExceedToYearError{}
	}
	return "", nil
}
