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
	if payload.General.EventDate.Month == 0 {
		return &MonthRequiredError{}
	}
	if payload.General.EventDate.Day == 0 {
		return &DayRequiredError{}
	}

	// Int values for hour and minute can be zero
	if payload.General.EventDate.FromHour == nil {
		return &FromHourRequiredError{}
	}
	if payload.General.EventDate.FromMinute == nil {
		return &FromMinuteRequiredError{}
	}
	if payload.General.EventDate.ToHour == nil {
		return &ToHourRequiredError{}
	}
	if payload.General.EventDate.ToMinute == nil {
		return &ToMinuteRequiredError{}
	}

	return nil
}
