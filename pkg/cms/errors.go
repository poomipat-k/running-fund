package cms

type AllowNewProjectRequiredError struct{}

func (e *AllowNewProjectRequiredError) Error() string {
	return "allowNewProject is required"
}

type FromYearRequiredError struct{}

func (e *FromYearRequiredError) Error() string {
	return "fromYear is required"
}

type ToYearRequiredError struct{}

func (e *ToYearRequiredError) Error() string {
	return "toYear is required"
}

type PageNoInvalidError struct{}

func (e *PageNoInvalidError) Error() string {
	return "pageNo is invalid"
}

type PageSizeInvalidError struct{}

func (e *PageSizeInvalidError) Error() string {
	return "pageSize is invalid"
}

type MonthRequiredError struct{}

func (e *MonthRequiredError) Error() string {
	return "month is required"
}

type MonthOutOfBoundError struct{}

func (e *MonthOutOfBoundError) Error() string {
	return "month must greater than > 0 and <= 12"
}

type FromDateExceedToDateError struct{}

func (e *FromDateExceedToDateError) Error() string {
	return "fromDate is later than toDate"
}
