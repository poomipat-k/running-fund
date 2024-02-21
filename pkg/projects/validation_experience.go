package projects

func validateExperience(payload AddProjectRequest) error {
	if payload.Experience.ThisSeries.FirstTime == nil {
		return &ThisSeriesFirstTimeRequiredError{}
	}
	return nil
}
