package review

type AddReviewRequest struct {
	Comment string  `json:"comment"`
	Ip      ip      `json:"ip"`
	Score   scoreV1 `json:"score"`
}

type ip struct {
	IsInterestedPerson   bool   `json:"isInterestedPerson"`
	InterestedPersonType string `json:"interestedPersonType"`
}

type scoreV1 struct {
	ReviewSummary string      `json:"reviewSummary"`
	Improvement   improvement `json:"improvement"`
	Q_1_1         int         `json:"q_1_1"`
	Q_1_2         int         `json:"q_1_2"`
	Q_1_3         int         `json:"q_1_3"`
	Q_1_4         int         `json:"q_1_4"`
	Q_1_5         int         `json:"q_1_5"`
	Q_1_6         int         `json:"q_1_6"`
	Q_1_7         int         `json:"q_1_7"`
	Q_1_8         int         `json:"q_1_8"`
	Q_1_9         int         `json:"q_1_9"`
	Q_1_10        int         `json:"q_1_10"`
	Q_1_11        int         `json:"q_1_11"`
	Q_1_12        int         `json:"q_1_12"`
	Q_1_13        int         `json:"q_1_13"`
	Q_1_14        int         `json:"q_1_14"`
	Q_1_15        int         `json:"q_1_15"`
	Q_1_16        int         `json:"q_1_16"`
	Q_1_17        int         `json:"q_1_17"`
	Q_1_18        int         `json:"q_1_18"`
	Q_1_19        int         `json:"q_1_19"`
	Q_1_20        int         `json:"q_1_20"`
}

type improvement struct {
	ProjectQuality           bool `json:"projectQuality"`
	ProjectStandard          bool `json:"projectStandard"`
	VisionAndImage           bool `json:"visionAndImage"`
	Benefit                  bool `json:"benefit"`
	ExperienceAndReliability bool `json:"experienceAndReliability"`
	FundAndOutput            bool `json:"fundAndOutput"`
}
