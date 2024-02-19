package projects

type CollaboratedRequiredError struct{}

func (e *CollaboratedRequiredError) Error() string {
	return "collaborated is required"
}

type CollaboratedFilesRequiredError struct{}

func (e *CollaboratedFilesRequiredError) Error() string {
	return "collaboratedFiles are required"
}

type ProjectNameRequiredError struct{}

func (e *ProjectNameRequiredError) Error() string {
	return "projectName is required"
}

type YearRequiredError struct{}

func (e *YearRequiredError) Error() string {
	return "year is required"
}

type MonthRequiredError struct{}

func (e *MonthRequiredError) Error() string {
	return "month is required"
}

type DayRequiredError struct{}

func (e *DayRequiredError) Error() string {
	return "day is required"
}

type FromHourRequiredError struct{}

func (e *FromHourRequiredError) Error() string {
	return "fromHour is required"
}

type FromMinuteRequiredError struct{}

func (e *FromMinuteRequiredError) Error() string {
	return "fromMinute is required"
}

type ToHourRequiredError struct{}

func (e *ToHourRequiredError) Error() string {
	return "toHour is required"
}

type ToMinuteRequiredError struct{}

func (e *ToMinuteRequiredError) Error() string {
	return "toMinute is required"
}
