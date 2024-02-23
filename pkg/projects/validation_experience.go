package projects

const MIN_YEAR = 2010

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
	if payload.Experience.ThisSeries.History.Year < MIN_YEAR || payload.Experience.ThisSeries.History.Year > localYear {
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
	if payload.Experience.ThisSeries.History.Completed1.Year == 0 {
		return &CompletedYearRequiredError{}
	}
	if payload.Experience.ThisSeries.History.Completed1.Year < MIN_YEAR || payload.Experience.ThisSeries.History.Completed1.Year > localYear {
		return &CompletedYearOutOfBoundError{}
	}
	if payload.Experience.ThisSeries.History.Completed1.Name == "" {
		return &CompletedNameRequiredError{}
	}
	if payload.Experience.ThisSeries.History.Completed1.Participant == 0 {
		return &CompletedParticipantRequiredError{}
	}
	if payload.Experience.ThisSeries.History.Completed1.Participant < 0 {
		return &CompletedParticipantInvalidError{}
	}
	// Validate completed2
	if payload.Experience.ThisSeries.History.Completed2.Year != 0 ||
		payload.Experience.ThisSeries.History.Completed2.Name != "" ||
		payload.Experience.ThisSeries.History.Completed2.Participant != 0 {
		err := historyCompletedIsValid(
			payload.Experience.ThisSeries.History.Completed2.Year,
			payload.Experience.ThisSeries.History.Completed2.Name,
			payload.Experience.ThisSeries.History.Completed2.Participant,
		)
		if err != nil {
			return err
		}
	}
	// Validate completed3
	if payload.Experience.ThisSeries.History.Completed3.Year != 0 ||
		payload.Experience.ThisSeries.History.Completed3.Name != "" ||
		payload.Experience.ThisSeries.History.Completed3.Participant != 0 {
		err := historyCompletedIsValid(
			payload.Experience.ThisSeries.History.Completed3.Year,
			payload.Experience.ThisSeries.History.Completed3.Name,
			payload.Experience.ThisSeries.History.Completed3.Participant,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func historyCompletedIsValid(year int, name string, participant int) error {
	localYear, _, _ := getLocalYearMonthDay()
	if year == 0 {
		return &CompletedYearRequiredError{}
	}
	if year < MIN_YEAR || year > localYear {
		return &CompletedYearOutOfBoundError{}
	}
	if name == "" {
		return &CompletedNameRequiredError{}
	}
	if participant == 0 {
		return &CompletedParticipantRequiredError{}
	}
	if participant < 0 {
		return &CompletedParticipantInvalidError{}
	}

	return nil
}
