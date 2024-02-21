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

func validateAddProjectPayload(payload AddProjectRequest, collaborateFiles []*multipart.FileHeader, criteria []ApplicantSelfScoreCriteria) error {
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
	if err := validateDetails(payload); err != nil {
		return err
	}

	return nil
}
