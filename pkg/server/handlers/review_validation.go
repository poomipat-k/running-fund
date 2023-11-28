package server

import (
	"fmt"
	"log"

	"github.com/poomipat-k/running-fund/pkg/projects"
)

func validateAddPayload(payload projects.AddReviewRequest, criteriaList []projects.ProjectReviewCriteriaMinimal) error {
	if payload.ProjectHistoryId == 0 {
		return fmt.Errorf("projectHistoryId is required")
	}

	err := validateInterestedPerson(payload.Ip)
	if err != nil {
		return err
	}

	err = validateReview(payload.Review, criteriaList)
	if err != nil {
		return err
	}
	log.Println("Looks good")
	return nil
}

func validateInterestedPerson(ip projects.Ip) error {
	if ip.IsInterestedPerson == nil {
		return fmt.Errorf("ip.isInterestedPerson is required")
	}
	if *ip.IsInterestedPerson && ip.InterestedPersonType == "" {
		return fmt.Errorf("ip.InterestedPersonType is required")
	}
	return nil
}

func validateReview(review projects.Review, criteriaList []projects.ProjectReviewCriteriaMinimal) error {
	err := validateScores(review.Scores, criteriaList)
	if err != nil {
		return err
	}
	err = validateReviewSummary(review)
	if err != nil {
		return err
	}
	return nil
}

func validateScores(scores map[string]int, criteriaList []projects.ProjectReviewCriteriaMinimal) error {
	for _, v := range criteriaList {
		name := fmt.Sprintf("q_%d_%d", v.CriteriaVersion, v.OrderNumber)
		score, valid := scores[name]
		if !valid || score == 0 {
			return fmt.Errorf("review.scores.%s is required", name)
		}
		if score < 1 || score > 5 {
			return fmt.Errorf("review.scores.%s value is out of range [1, 5]", name)
		}
	}
	return nil
}

func validateReviewSummary(review projects.Review) error {
	log.Println(review)
	switch review.ReviewSummary {
	case "":
		return fmt.Errorf("review.reviewSummary is required")
	case "ok":
		return nil
	case "not_ok":
		return nil
	case "to_be_revised":
		err := validateImprovement(review.Improvement)
		if err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("invalid review.reviewSummary")
	}
}

func validateImprovement(im projects.ReviewImprovement) error {
	if im.Benefit == nil {
		return fmt.Errorf("review.improvement.benefit is required")
	}
	if im.ExperienceAndReliability == nil {
		return fmt.Errorf("review.improvement.experienceAndReliability is required")
	}
	if im.FundAndOutput == nil {
		return fmt.Errorf("review.improvement.fundAndOutput is required")
	}
	if im.ProjectQuality == nil {
		return fmt.Errorf("review.improvement.projectQuality is required")
	}
	if im.ProjectStandard == nil {
		return fmt.Errorf("review.improvement.projectStandard is required")
	}
	if im.VisionAndImage == nil {
		return fmt.Errorf("review.improvement.visionAndImage is required")
	}
	return nil
}
