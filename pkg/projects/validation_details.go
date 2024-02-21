package projects

import (
	"fmt"
	"log"
)

func validateDetails(payload AddProjectRequest, criteria []ApplicantSelfScoreCriteria) error {
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
	// score
	log.Println("==criteria")
	log.Println(criteria)
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

func isScoreValid(score int) bool {
	return score >= 1 && score <= 5
}
