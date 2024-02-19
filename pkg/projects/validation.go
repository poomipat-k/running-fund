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
	if payload.General.EventDate.Year < 1971 {
		return &YearInvalidError{}
	}
	if payload.General.EventDate.Month == 0 {
		return &MonthRequiredError{}
	}
	if payload.General.EventDate.Month < 1 || payload.General.EventDate.Month > 12 {
		return &MonthOutOfBoundError{}
	}
	if payload.General.EventDate.Day == 0 {
		return &DayRequiredError{}
	}
	if payload.General.EventDate.Day < 1 || payload.General.EventDate.Day > 31 {
		return &DayOutOfBoundError{}
	}

	// Int values for hour and minute can be zero
	if payload.General.EventDate.FromHour == nil {
		return &FromHourRequiredError{}
	}
	if *payload.General.EventDate.FromHour < 0 || *payload.General.EventDate.FromHour > 23 {
		return &InvalidError{Name: "fromHour"}
	}
	if payload.General.EventDate.FromMinute == nil {
		return &FromMinuteRequiredError{}
	}
	if *payload.General.EventDate.FromMinute < 0 || *payload.General.EventDate.FromMinute > 59 {
		return &InvalidError{Name: "fromMinute"}
	}
	if payload.General.EventDate.ToHour == nil {
		return &ToHourRequiredError{}
	}
	if *payload.General.EventDate.ToHour < 0 || *payload.General.EventDate.ToHour > 23 {
		return &InvalidError{Name: "toHour"}
	}
	if payload.General.EventDate.ToMinute == nil {
		return &ToMinuteRequiredError{}
	}
	if *payload.General.EventDate.ToMinute < 0 || *payload.General.EventDate.ToMinute > 59 {
		return &InvalidError{Name: "toMinute"}
	}

	return nil
}
