package review

type AddReviewRequest struct {
	ProjectHistoryId int    `json:"projectHistoryId,omitempty"`
	Comment          string `json:"comment"`
	Ip               ip     `json:"ip"`
	Review           review `json:"review"`
}

type ip struct {
	IsInterestedPerson   *bool  `json:"isInterestedPerson,omitempty"`
	InterestedPersonType string `json:"interestedPersonType,omitempty"`
}

type review struct {
	ReviewSummary string            `json:"reviewSummary,omitempty"`
	Improvement   reviewImprovement `json:"improvement,omitempty"`
	Scores        map[string]int    `json:"scores,omitempty"`
}

type reviewImprovement struct {
	Benefit                  *bool `json:"benefit,omitempty"`
	ExperienceAndReliability *bool `json:"experienceAndReliability,omitempty"`
	FundAndOutput            *bool `json:"fundAndOutput,omitempty"`
	ProjectQuality           *bool `json:"projectQuality,omitempty"`
	ProjectStandard          *bool `json:"projectStandard,omitempty"`
	VisionAndImage           *bool `json:"visionAndImage,omitempty"`
}

type ProjectReviewCriteriaMinimal struct {
	CriteriaId      int `json:"reviewCriteriaId,omitempty"`
	CriteriaVersion int `json:"criteriaVersion,omitempty"`
	OrderNumber     int `json:"orderNumber,omitempty"`
}
