package projects_test

import (
	"net/http"

	"github.com/poomipat-k/running-fund/pkg/mock"
	"github.com/poomipat-k/running-fund/pkg/projects"
)

var Details = []TestCase{
	{
		name: "should error when details.background is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.BackgroundRequiredError{},
	},
	{
		name: "should error when details.objective is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details: projects.Details{
				Background: "Some background",
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ObjectiveRequiredError{},
	},
	// marketing
	// marketing.online
	{
		name: "should error when none of details.marketing.online.available is selected",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details: projects.Details{
				Background: "Some background",
				Objective:  "Some objective",
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.OnlineAvailableRequiredOne{},
	},
	{
		name: "should error when  details.marketing.online.available.facebook is checked and details.marketing.online.howTo.facebook is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details: projects.Details{
				Background: "Some background",
				Objective:  "Some objective",
				Marketing: projects.Marketing{
					Online: projects.Online{
						Available: projects.OnlineAvailable{
							Facebook: true,
						},
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.FacebookHowToIsRequired{},
	},
	{
		name: "should error when  details.marketing.online.available.website is checked and details.marketing.online.howTo.website is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details: projects.Details{
				Background: "Some background",
				Objective:  "Some objective",
				Marketing: projects.Marketing{
					Online: projects.Online{
						Available: projects.OnlineAvailable{
							Facebook: true,
							Website:  true,
						},
						HowTo: projects.OnlineHowTo{
							Facebook: "facebook.com/abc",
						},
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.WebsiteHowToIsRequired{},
	},
	{
		name: "should error when  details.marketing.online.available.onlinePage is checked and details.marketing.online.howTo.onlinePage is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details: projects.Details{
				Background: "Some background",
				Objective:  "Some objective",
				Marketing: projects.Marketing{
					Online: projects.Online{
						Available: projects.OnlineAvailable{
							Facebook:   true,
							Website:    true,
							OnlinePage: true,
						},
						HowTo: projects.OnlineHowTo{
							Facebook: "facebook.com/abc",
							Website:  "test.com",
						},
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.OnlinePageHowToIsRequired{},
	},
	{
		name: "should error when  details.marketing.online.available.other is checked and details.marketing.online.howTo.other is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details: projects.Details{
				Background: "Some background",
				Objective:  "Some objective",
				Marketing: projects.Marketing{
					Online: projects.Online{
						Available: projects.OnlineAvailable{
							Facebook:   true,
							Website:    true,
							OnlinePage: true,
							Other:      true,
						},
						HowTo: projects.OnlineHowTo{
							Facebook:   "facebook.com/abc",
							Website:    "test.com",
							OnlinePage: "abc",
						},
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.OtherHowToIsRequired{},
	},
	// marketing offline
	{
		name: "should error when none of details.marketing.offline.available is checked",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details: projects.Details{
				Background: "Some background",
				Objective:  "Some objective",
				Marketing: projects.Marketing{
					Online: projects.Online{
						Available: projects.OnlineAvailable{
							Facebook:   true,
							Website:    true,
							OnlinePage: true,
							Other:      false,
						},
						HowTo: projects.OnlineHowTo{
							Facebook:   "facebook.com/abc",
							Website:    "test.com",
							OnlinePage: "abc",
						},
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.OfflineAvailableRequiredOne{},
	},
	{
		name: "should error when  details.marketing.offline.available is 'other' and details.marketing.offline.addition is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details: projects.Details{
				Background: "Some background",
				Objective:  "Some objective",
				Marketing: projects.Marketing{
					Online: projects.Online{
						Available: projects.OnlineAvailable{
							Facebook:   true,
							Website:    true,
							OnlinePage: true,
							Other:      false,
						},
						HowTo: projects.OnlineHowTo{
							Facebook:   "facebook.com/abc",
							Website:    "test.com",
							OnlinePage: "abc",
						},
					},
					Offline: projects.Offline{
						Available: projects.OfflineAvailable{
							Other: true,
						},
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.OfflineAdditionRequiredError{},
	},
	// details.score
	{
		name: "should error when applicant criteria len = 0",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details: projects.Details{
				Background: "Some background",
				Objective:  "Some objective",
				Marketing: projects.Marketing{
					Online: projects.Online{
						Available: projects.OnlineAvailable{
							Facebook:   true,
							Website:    true,
							OnlinePage: true,
							Other:      false,
						},
						HowTo: projects.OnlineHowTo{
							Facebook:   "facebook.com/abc",
							Website:    "test.com",
							OnlinePage: "abc",
						},
					},
					Offline: projects.Offline{
						Available: projects.OfflineAvailable{
							Other: true,
						},
						Addition: "Test",
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaNotFound,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ApplicantCriteriaNotFoundError{},
	},
	{
		name: "should error when some score.q_{version}_{order} is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details: projects.Details{
				Background: "Some background",
				Objective:  "Some objective",
				Marketing: projects.Marketing{
					Online: projects.Online{
						Available: projects.OnlineAvailable{
							Facebook:   true,
							Website:    true,
							OnlinePage: true,
							Other:      false,
						},
						HowTo: projects.OnlineHowTo{
							Facebook:   "facebook.com/abc",
							Website:    "test.com",
							OnlinePage: "abc",
						},
					},
					Offline: projects.Offline{
						Available: projects.OfflineAvailable{
							Other: true,
						},
						Addition: "Test",
					},
				},
				Score: map[string]int{
					"q_1_1": 4,
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ScoreRequiredError{Name: "q_1_2"},
	},
	{
		name: "should error when first score.q_1_1 is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details: projects.Details{
				Background: "Some background",
				Objective:  "Some objective",
				Marketing: projects.Marketing{
					Online: projects.Online{
						Available: projects.OnlineAvailable{
							Facebook:   true,
							Website:    true,
							OnlinePage: true,
							Other:      false,
						},
						HowTo: projects.OnlineHowTo{
							Facebook:   "facebook.com/abc",
							Website:    "test.com",
							OnlinePage: "abc",
						},
					},
					Offline: projects.Offline{
						Available: projects.OfflineAvailable{
							Other: true,
						},
						Addition: "Test",
					},
				},
				Score: map[string]int{
					"q_1_2": 4,
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ScoreRequiredError{Name: "q_1_1"},
	},
	{
		name: "should error when not all details.score.q_[criteriaVersion]_[orderNumber] is valid",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details: projects.Details{
				Background: "Some background",
				Objective:  "Some objective",
				Marketing: projects.Marketing{
					Online: projects.Online{
						Available: projects.OnlineAvailable{
							Facebook:   true,
							Website:    true,
							OnlinePage: true,
							Other:      false,
						},
						HowTo: projects.OnlineHowTo{
							Facebook:   "facebook.com/abc",
							Website:    "test.com",
							OnlinePage: "abc",
						},
					},
					Offline: projects.Offline{
						Available: projects.OfflineAvailable{
							Other: true,
						},
						Addition: "Test",
					},
				},
				Score: map[string]int{
					"q_1_1": 4,
					"q_1_2": 24,
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ScoreInvalidError{Name: "q_1_2"},
	},
}
