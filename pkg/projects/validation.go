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
	if payload.ProjectStatusSecondary == "" {
		return "projectStatusSecondary", &ProjectStatusSecondaryRequiredError{}
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
