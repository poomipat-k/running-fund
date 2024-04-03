package projects

var expectedParticipantsOptions = map[string]bool{
	"<=500":     true,
	"501-1500":  true,
	"1501-2500": true,
	"2501-3500": true,
	"3501-4500": true,
	"4501-5500": true,
	">=5501":    true,
}

func validateGeneral(payload AddProjectRequest) error {

	if payload.General.ProjectName == "" {
		return &ProjectNameRequiredError{}
	}

	// general.eventDate
	if err := validateEventDate(payload); err != nil {
		return err
	}

	// general.address
	if err := validateGeneralAddress(payload); err != nil {
		return err
	}

	// general.startPoint and general.finishPoint
	if payload.General.StartPoint == "" {
		return &StartPointRequiredError{}
	}
	if payload.General.FinishPoint == "" {
		return &FinishPointRequiredError{}
	}
	// general.eventDetails
	if err := validateGeneralEventDetails(payload); err != nil {
		return err
	}

	// general.expectedParticipants
	if payload.General.ExpectedParticipants == "" {
		return &ExpectedParticipantsRequiredError{}
	}
	if _, found := expectedParticipantsOptions[payload.General.ExpectedParticipants]; !found {
		return &ExpectedParticipantsInvalidError{}
	}
	// general.hasOrganizer
	if payload.General.HasOrganizer == nil {
		return &HasOrganizerRequiredError{}
	}
	if *payload.General.HasOrganizer && payload.General.OrganizerName == "" {
		return &OrganizerNameRequiredError{}
	}
	return nil
}

func validateEventDate(payload AddProjectRequest) error {
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

func validateGeneralAddress(payload AddProjectRequest) error {
	if payload.General.Address.Address == "" {
		return &GeneralAddressRequiredError{}
	}
	if payload.General.Address.ProvinceId <= 0 {
		return &GeneralProvinceRequiredError{}
	}
	if payload.General.Address.DistrictId <= 0 {
		return &GeneralDistrictIdRequiredError{}
	}
	if payload.General.Address.SubdistrictId <= 0 {
		return &GeneralSubdistrictIdRequiredError{}
	}
	if payload.General.Address.PostcodeId <= 0 {
		return &GeneralPostcodeIdRequiredError{}
	}
	return nil
}

func validateGeneralEventDetails(payload AddProjectRequest) error {
	// general.eventDetails.category
	if !payload.General.EventDetails.Category.Available.RoadRace &&
		!payload.General.EventDetails.Category.Available.TrailRunning &&
		!payload.General.EventDetails.Category.Available.Other {
		return &CategoryAvailableRequiredOneError{}
	}
	if payload.General.EventDetails.Category.Available.Other && payload.General.EventDetails.Category.OtherType == "" {
		return &OtherEventTypeRequiredError{}
	}
	// general.eventDetails.distanceAndFee
	if len(payload.General.EventDetails.DistanceAndFee) == 0 {
		return &DistanceRequiredOneError{}
	}
	dfCount := 0
	for _, df := range payload.General.EventDetails.DistanceAndFee {
		if df.Checked {
			dfCount++
			if df.Type == "" {
				return &DistanceTypeRequiredError{}
			}
			if df.Fee == nil {
				return &DistanceFeeRequiredError{}
			}
			if *df.Fee < 0 {
				return &ValueNegativeError{}
			}
			if df.Dynamic == nil {
				return &DistanceAndFeeDynamicRequiredError{}
			}
		}
	}
	if dfCount == 0 {
		return &DistanceRequiredOneError{}
	}
	// general.eventDetails.vip
	if payload.General.EventDetails.VIP == nil {
		return &VIPRequiredError{}
	}
	if *payload.General.EventDetails.VIP && payload.General.EventDetails.VIPFree == nil {
		return &VIPFeeRequiredError{}
	}
	if *payload.General.EventDetails.VIP && payload.General.EventDetails.VIPFree != nil && *payload.General.EventDetails.VIPFree < 0 {
		return &VIPFeeNegativeError{}
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
