package projects

const MIN_YEAR = 1957

func validateExperience(payload AddProjectRequest) error {
	experience := payload.Experience
	localYear, _, _ := getLocalYearMonthDay()
	// thisSeries
	if experience.ThisSeries.FirstTime == nil {
		return &ThisSeriesFirstTimeRequiredError{}
	}
	if !(*experience.ThisSeries.FirstTime) {
		if experience.ThisSeries.History.OrdinalNumber < 2 {
			return &HistoryOrdinalNumberInvalidError{}
		}
		if experience.ThisSeries.History.Year == 0 {
			return &HistoryYearRequiredError{}
		}
		if experience.ThisSeries.History.Year < MIN_YEAR || experience.ThisSeries.History.Year > localYear {
			return &HistoryYearOutOfBoundError{}
		}
		if experience.ThisSeries.History.Month == 0 {
			return &HistoryMonthRequiredError{}
		}
		if experience.ThisSeries.History.Month <= 0 || experience.ThisSeries.History.Month > 12 {
			return &HistoryMonthOutOfBoundError{}
		}
		if experience.ThisSeries.History.Day == 0 {
			return &HistoryDayRequiredError{}
		}
		if !isValidDay(experience.ThisSeries.History.Year, experience.ThisSeries.History.Month, experience.ThisSeries.History.Day) {
			return &HistoryDayOutOfBoundError{}
		}
		if experience.ThisSeries.History.Completed1.Year == 0 {
			return &CompletedYearRequiredError{}
		}
		if experience.ThisSeries.History.Completed1.Year < MIN_YEAR || experience.ThisSeries.History.Completed1.Year > localYear {
			return &CompletedYearOutOfBoundError{}
		}
		if experience.ThisSeries.History.Completed1.Participant == 0 {
			return &CompletedParticipantRequiredError{}
		}
		if experience.ThisSeries.History.Completed1.Participant < 0 {
			return &CompletedParticipantInvalidError{}
		}
		// thisSeries.completed2
		if experience.ThisSeries.History.Completed2.Year != 0 ||
			experience.ThisSeries.History.Completed2.Participant != 0 {
			err := thisHistoryCompletedIsValid(
				experience.ThisSeries.History.Completed2.Year,
				experience.ThisSeries.History.Completed2.Participant,
			)
			if err != nil {
				return err
			}
		}
		// thisSeries.completed3
		if experience.ThisSeries.History.Completed3.Year != 0 ||
			payload.Experience.ThisSeries.History.Completed3.Participant != 0 {
			err := thisHistoryCompletedIsValid(
				payload.Experience.ThisSeries.History.Completed3.Year,
				payload.Experience.ThisSeries.History.Completed3.Participant,
			)
			if err != nil {
				return err
			}
		}
	}
	// otherSeries
	// otherSeries.history.completed1
	if experience.OtherSeries.DoneBefore == nil {
		return &DoneBeforeRequiredError{}
	}
	if *experience.OtherSeries.DoneBefore {
		if experience.OtherSeries.History.Completed1.Year == 0 {
			return &CompletedYearRequiredError{}
		}
		if experience.OtherSeries.History.Completed1.Year < MIN_YEAR || experience.OtherSeries.History.Completed1.Year > localYear {
			return &CompletedYearOutOfBoundError{}
		}
		if experience.OtherSeries.History.Completed1.Name == "" {
			return &CompletedNameRequiredError{}
		}
		if experience.OtherSeries.History.Completed1.Participant == 0 {
			return &CompletedParticipantRequiredError{}
		}
		if experience.OtherSeries.History.Completed1.Participant < 0 {
			return &CompletedParticipantInvalidError{}
		}
		// otherSeries.history.completed2
		if experience.OtherSeries.History.Completed2.Year != 0 ||
			experience.OtherSeries.History.Completed2.Name != "" ||
			experience.OtherSeries.History.Completed2.Participant != 0 {
			err := otherHistoryCompletedIsValid(
				experience.OtherSeries.History.Completed2.Year,
				experience.OtherSeries.History.Completed2.Name,
				experience.OtherSeries.History.Completed2.Participant,
			)
			if err != nil {
				return err
			}
		}
		// otherSeries.history.completed3
		if experience.OtherSeries.History.Completed3.Year != 0 ||
			experience.OtherSeries.History.Completed3.Name != "" ||
			experience.OtherSeries.History.Completed3.Participant != 0 {
			err := otherHistoryCompletedIsValid(
				experience.OtherSeries.History.Completed3.Year,
				experience.OtherSeries.History.Completed3.Name,
				experience.OtherSeries.History.Completed3.Participant,
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func thisHistoryCompletedIsValid(year int, participant int) error {
	localYear, _, _ := getLocalYearMonthDay()
	if year == 0 {
		return &CompletedYearRequiredError{}
	}
	if year < MIN_YEAR || year > localYear {
		return &CompletedYearOutOfBoundError{}
	}
	if participant == 0 {
		return &CompletedParticipantRequiredError{}
	}
	if participant < 0 {
		return &CompletedParticipantInvalidError{}
	}

	return nil
}

func otherHistoryCompletedIsValid(year int, name string, participant int) error {
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
