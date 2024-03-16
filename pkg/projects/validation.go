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
	marketingFiles, routeFiles, eventMapFiles, eventDetailsFiles, screenshotFiles []*multipart.FileHeader) error {
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
	if err := validateAttachment(marketingFiles, routeFiles, eventMapFiles, eventDetailsFiles, screenshotFiles); err != nil {
		return err
	}

	return nil
}
