package projects

func validateDetails(payload AddProjectRequest) error {
	if payload.Details.Background == "" {
		return &BackgroundRequiredError{}
	}
	if payload.Details.Objective == "" {
		return &ObjectiveRequiredError{}
	}
	// marketing
	// marketing.online
	if !payload.Details.Marketing.Online.Available.Facebook &&
		!payload.Details.Marketing.Online.Available.Website &&
		!payload.Details.Marketing.Online.Available.OnlinePage &&
		!payload.Details.Marketing.Online.Available.Other {
		return &OnlineAvailableRequiredOne{}
	}
	if payload.Details.Marketing.Online.Available.Facebook && payload.Details.Marketing.Online.HowTo.Facebook == "" {
		return &FacebookHowToIsRequired{}
	}
	if payload.Details.Marketing.Online.Available.Website && payload.Details.Marketing.Online.HowTo.Website == "" {
		return &WebsiteHowToIsRequired{}
	}
	if payload.Details.Marketing.Online.Available.OnlinePage && payload.Details.Marketing.Online.HowTo.OnlinePage == "" {
		return &OnlinePageHowToIsRequired{}
	}
	if payload.Details.Marketing.Online.Available.Other && payload.Details.Marketing.Online.HowTo.Other == "" {
		return &OtherHowToIsRequired{}
	}
	return nil
}
