package projects

func validateExperience(payload AddProjectRequest) error {
	if payload.Experience.ThisSeries.FirstTime == nil {
		return &ThisSeriesFirstTimeRequiredError{}
	}
	if payload.Experience.ThisSeries.History.OrdinalNumber < 1 {
		return &HistoryOrdinalNumberInvalidError{}
	}
	return nil
}
