package projects

import (
	"mime/multipart"
)

// []int{4, 6, 9, 11}
var thirtyDaysMonth = map[int]int{
	4:  30,
	6:  30,
	9:  30,
	11: 30,
}

func validateAddProjectPayload(payload AddProjectRequest, collaborateFiles []*multipart.FileHeader) error {
	if payload.Collaborated == nil {
		return &CollaboratedRequiredError{}
	}
	if *payload.Collaborated && len(collaborateFiles) == 0 {
		return &CollaboratedFilesRequiredError{}
	}

	// general.event
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

	if !isValidDay(payload.General.EventDate.Year, payload.General.EventDate.Month, payload.General.EventDate.Day) {
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

	// general.address
	if payload.General.Address.Address == "" {
		return &AddressRequiredError{}
	}
	if payload.General.Address.ProvinceId <= 0 {
		return &ProvinceRequiredError{}
	}
	if payload.General.Address.DistrictId <= 0 {
		return &DistrictIdRequiredError{}
	}
	if payload.General.Address.SubdistrictId <= 0 {
		return &SubdistrictIdRequiredError{}
	}
	if payload.General.Address.PostcodeId <= 0 {
		return &PostcodeIdRequiredError{}
	}

	return nil
}

// Assume we have a valid year and month already
func isValidDay(year, month, day int) bool {
	if day < 1 || day > 31 {
		return false
	}
	if month == 2 {
		leapYear := isLeapYear(year)
		if leapYear && day > 29 {
			return false
		}
		if !leapYear && day > 28 {
			return false
		}
	}
	_, isThirtyDayMonth := thirtyDaysMonth[month]
	if isThirtyDayMonth && day > 30 {
		return false
	}
	return true
}

func isLeapYear(year int) bool {
	return year%4 == 0 && year%100 != 0 || year%400 == 0
}
