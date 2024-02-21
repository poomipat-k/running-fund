package projects

func validateExperience(payload AddProjectRequest) error {
	if payload.Experience.ThisSeries.FirstTime == nil {
		return &ThisSeriesFirstTimeRequiredError{}
	}
	if payload.Experience.ThisSeries.History.OrdinalNumber < 2 {
		return &HistoryOrdinalNumberInvalidError{}
	}
	if payload.Experience.ThisSeries.History.Year == 0 {
		return &HistoryYearRequiredError{}
	}
	localYear, _, _ := getLocalYearMonthDay()
	if payload.Experience.ThisSeries.History.Year < 2018 || payload.Experience.ThisSeries.History.Year > localYear {
		return &HistoryYearOutOfBoundError{}
	}
	if payload.Experience.ThisSeries.History.Month == 0 {
		return &HistoryMonthRequiredError{}
	}
	if payload.Experience.ThisSeries.History.Month <= 0 || payload.Experience.ThisSeries.History.Month > 12 {
		return &HistoryMonthOutOfBoundError{}
	}
	if payload.Experience.ThisSeries.History.Day == 0 {
		return &HistoryDayRequiredError{}
	}
	if !isValidDay(payload.Experience.ThisSeries.History.Year, payload.Experience.ThisSeries.History.Month, payload.Experience.ThisSeries.History.Day) {
		return &HistoryDayOutOfBoundError{}
	}
	return nil
}
