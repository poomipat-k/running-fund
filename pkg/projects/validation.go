package projects

import (
	"mime/multipart"
)

func validateAddProjectPayload(payload AddProjectRequest, collaborateFiles []*multipart.FileHeader) error {
	if payload.Collaborated == nil {
		return &CollaboratedRequiredError{}
	}
	if *payload.Collaborated && len(collaborateFiles) == 0 {
		return &CollaboratedFilesRequiredError{}
	}

	if payload.General.ProjectName == "" {
		return &ProjectNameRequiredError{}
	}
	if payload.General.EventDate.Year == 0 {
		return &YearRequiredError{}
	}

	return nil
}
