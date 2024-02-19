package projects

import "fmt"

type InvalidError struct {
	Name string
}

func (e *InvalidError) Error() string {
	return fmt.Sprintf("%s is invalid", e.Name)
}

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

type YearInvalidError struct{}

func (e *YearInvalidError) Error() string {
	return "year must greater than 1971"
}

type MonthRequiredError struct{}

func (e *MonthRequiredError) Error() string {
	return "month is required"
}

type MonthOutOfBoundError struct{}

func (e *MonthOutOfBoundError) Error() string {
	return "month must greater than > 0 and <= 12"
}

type DayRequiredError struct{}

func (e *DayRequiredError) Error() string {
	return "day is required"
}

type DayOutOfBoundError struct{}

func (e *DayOutOfBoundError) Error() string {
	return "day is not valid"
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

type AddressRequiredError struct{}

func (e *AddressRequiredError) Error() string {
	return "address is required"
}

type ProvinceRequiredError struct{}

func (e *ProvinceRequiredError) Error() string {
	return "provinceId is required"
}

type DistrictIdRequiredError struct{}

func (e *DistrictIdRequiredError) Error() string {
	return "districtId is required"
}

type SubdistrictIdRequiredError struct{}

func (e *SubdistrictIdRequiredError) Error() string {
	return "subdistrictIdId is required"
}

type PostcodeIdRequiredError struct{}

func (e *PostcodeIdRequiredError) Error() string {
	return "postcodeId is required"
}

type StartPointRequiredError struct{}

func (e *StartPointRequiredError) Error() string {
	return "startPoint is required"
}

type FinishPointRequiredError struct{}

func (e *FinishPointRequiredError) Error() string {
	return "finishPoint is required"
}
