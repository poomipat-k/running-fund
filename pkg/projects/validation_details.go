package projects

import (
	"fmt"
)

var judgeTypeOptions = map[string]bool{
	"manual": true,
	"auto":   true,
	"other":  true,
}

func validateDetails(payload AddProjectRequest, criteria []ApplicantSelfScoreCriteria) error {
	if payload.Details.Background == "" {
		return &BackgroundRequiredError{}
	}
	if payload.Details.Objective == "" {
		return &ObjectiveRequiredError{}
	}
	// marketing
	if err := validateMarketing(payload); err != nil {
		return err
	}

	// score
	if err := validateScore(payload, criteria); err != nil {
		return err
	}
	// safety
	if err := validateSafety(payload); err != nil {
		return err
	}
	// route
	if err := validateRoute(payload); err != nil {
		return err
	}
	// judge
	if err := validateJudge(payload); err != nil {
		return err
	}
	// support
	if err := validateSupport(payload); err != nil {
		return err
	}
	// feedback
	if payload.Details.Feedback == "" {
		return &FeedbackRequiredError{}
	}

	return nil
}

func validateMarketing(payload AddProjectRequest) error {
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
	// marketing.offline
	if !payload.Details.Marketing.Offline.Available.PR &&
		!payload.Details.Marketing.Offline.Available.LocalOfficial &&
		!payload.Details.Marketing.Offline.Available.Booth &&
		!payload.Details.Marketing.Offline.Available.Billboard &&
		!payload.Details.Marketing.Offline.Available.TV &&
		!payload.Details.Marketing.Offline.Available.Other {
		return &OfflineAvailableRequiredOne{}
	}
	if payload.Details.Marketing.Offline.Available.Other && payload.Details.Marketing.Offline.Addition == "" {
		return &OfflineAdditionRequiredError{}
	}
	return nil
}

func validateScore(payload AddProjectRequest, criteria []ApplicantSelfScoreCriteria) error {
	criteriaCount := len(criteria)
	if criteriaCount == 0 {
		return &ApplicantCriteriaNotFoundError{}
	}
	for _, c := range criteria {
		key := fmt.Sprintf("q_%d_%d", c.CriteriaVersion, c.OrderNumber)
		input, exists := payload.Details.Score[key]
		if !exists {
			return &ScoreRequiredError{Name: key}
		}
		if !isScoreValid(input) {
			return &ScoreInvalidError{Name: key}
		}
	}
	return nil
}

func validateSafety(payload AddProjectRequest) error {
	if !payload.Details.Safety.Ready.RunnerInformation &&
		!payload.Details.Safety.Ready.HealthDecider &&
		!payload.Details.Safety.Ready.Ambulance &&
		!payload.Details.Safety.Ready.FirstAid &&
		!payload.Details.Safety.Ready.AED &&
		!payload.Details.Safety.Ready.Insurance &&
		!payload.Details.Safety.Ready.Other {
		return &SafetyReadyRequiredOneError{}
	}
	if payload.Details.Safety.Ready.AED && payload.Details.Safety.AEDCount < 1 {
		return &AEDCountInvalidError{}
	}
	if payload.Details.Safety.Ready.Other && payload.Details.Safety.Addition == "" {
		return &SafetyAdditionRequiredError{}
	}
	return nil
}

func validateRoute(payload AddProjectRequest) error {
	if !payload.Details.Route.Measurement.AthleticsAssociation &&
		!payload.Details.Route.Measurement.CalibratedBicycle &&
		!payload.Details.Route.Measurement.SelfMeasurement {
		return &RouteMeasurementRequiredOneError{}
	}
	if payload.Details.Route.Measurement.SelfMeasurement && payload.Details.Route.Tool == "" {
		return &RouteToolRequiredError{}
	}
	if !payload.Details.Route.TrafficManagement.AskPermission &&
		!payload.Details.Route.TrafficManagement.HasSupporter &&
		!payload.Details.Route.TrafficManagement.RoadClosure &&
		!payload.Details.Route.TrafficManagement.Signs &&
		!payload.Details.Route.TrafficManagement.Lighting {
		return &RouteTrafficManagementRequiredOneError{}
	}
	return nil
}

func validateJudge(payload AddProjectRequest) error {
	if payload.Details.Judge.Type == "" {
		return &JudgeTypeRequiredError{}
	}
	if _, found := judgeTypeOptions[payload.Details.Judge.Type]; !found {
		return &JudgeTypeInvalidError{}
	}
	if payload.Details.Judge.Type == "other" && payload.Details.Judge.OtherType == "" {
		return &JudgeOtherTypeRequiredError{}
	}
	return nil
}

func validateSupport(payload AddProjectRequest) error {
	if !payload.Details.Support.Organization.ProvincialAdministration &&
		!payload.Details.Support.Organization.Safety &&
		!payload.Details.Support.Organization.Health &&
		!payload.Details.Support.Organization.Volunteer &&
		!payload.Details.Support.Organization.Community &&
		!payload.Details.Support.Organization.Other {
		return &SupportOrganizationRequiredOneError{}
	}
	if payload.Details.Support.Organization.Other && payload.Details.Support.Addition == "" {
		return &SupportAdditionRequiredError{}
	}
	return nil
}

func isScoreValid(score int) bool {
	return score >= 1 && score <= 5
}
