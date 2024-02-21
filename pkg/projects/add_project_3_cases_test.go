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
	// details.safety
	{
		name: "should error when none of details.safety.ready is checked",
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
					"q_1_2": 3,
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.SafetyReadyRequiredOneError{},
	},
	{
		name: "should error when details.safety.ready.aed is checked and details.safety.aedCount < 1",
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
					"q_1_2": 3,
				},
				Safety: projects.Safety{
					Ready: projects.SafetyReady{
						RunnerInformation: true,
						AED:               true,
					},
					AEDCount: 0,
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.AEDCountInvalidError{},
	},
	{
		name: "should error when details.safety.ready.other is checked and details.safety.addition is empty",
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
					"q_1_2": 3,
				},
				Safety: projects.Safety{
					Ready: projects.SafetyReady{
						RunnerInformation: true,
						AED:               true,
						Other:             true,
					},
					AEDCount: 5,
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.SafetyAdditionRequiredError{},
	},
	// details.route
	{
		name: "should error when none of details.route.measurement is checked",
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
					"q_1_2": 3,
				},
				Safety: projects.Safety{
					Ready: projects.SafetyReady{
						RunnerInformation: true,
						AED:               true,
						Other:             true,
					},
					AEDCount: 5,
					Addition: "X",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.RouteMeasurementRequiredOneError{},
	},
	{
		name: "should error when details.route.measurement.selfMeasurement is checked and details.route.tool is empty",
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
					"q_1_2": 3,
				},
				Safety: projects.Safety{
					Ready: projects.SafetyReady{
						RunnerInformation: true,
						AED:               true,
						Other:             true,
					},
					AEDCount: 5,
					Addition: "X",
				},
				Route: projects.Route{
					Measurement: projects.RouteMeasurement{
						SelfMeasurement: true,
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.RouteToolRequiredError{},
	},
	{
		name: "should error when none of details.route.trafficManagement is checked",
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
					"q_1_2": 3,
				},
				Safety: projects.Safety{
					Ready: projects.SafetyReady{
						RunnerInformation: true,
						AED:               true,
						Other:             true,
					},
					AEDCount: 5,
					Addition: "X",
				},
				Route: projects.Route{
					Measurement: projects.RouteMeasurement{
						SelfMeasurement: true,
					},
					Tool: "UU",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.RouteTrafficManagementRequiredOneError{},
	},
	// details.judge
	{
		name: "should error when details.judge.type is empty",
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
					"q_1_2": 3,
				},
				Safety: projects.Safety{
					Ready: projects.SafetyReady{
						RunnerInformation: true,
						AED:               true,
						Other:             true,
					},
					AEDCount: 5,
					Addition: "X",
				},
				Route: projects.Route{
					Measurement: projects.RouteMeasurement{
						SelfMeasurement: true,
					},
					Tool: "UU",
					TrafficManagement: projects.TrafficManagement{
						AskPermission: true,
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.JudgeTypeRequiredError{},
	},
	{
		name: "should error when details.judge.type value is invalid",
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
					"q_1_2": 3,
				},
				Safety: projects.Safety{
					Ready: projects.SafetyReady{
						RunnerInformation: true,
						AED:               true,
						Other:             true,
					},
					AEDCount: 5,
					Addition: "X",
				},
				Route: projects.Route{
					Measurement: projects.RouteMeasurement{
						SelfMeasurement: true,
					},
					Tool: "UU",
					TrafficManagement: projects.TrafficManagement{
						AskPermission: true,
					},
				},
				Judge: projects.Judge{
					Type: "NOT_EXIST_KEY",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.JudgeTypeInvalidError{},
	},
	{
		name: "should error when details.judge.type is other and details.judge.otherType is empty",
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
					"q_1_2": 3,
				},
				Safety: projects.Safety{
					Ready: projects.SafetyReady{
						RunnerInformation: true,
						AED:               true,
						Other:             true,
					},
					AEDCount: 5,
					Addition: "X",
				},
				Route: projects.Route{
					Measurement: projects.RouteMeasurement{
						SelfMeasurement: true,
					},
					Tool: "UU",
					TrafficManagement: projects.TrafficManagement{
						AskPermission: true,
					},
				},
				Judge: projects.Judge{
					Type: "other",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.JudgeOtherTypeRequiredError{},
	},
	// details.support
	{
		name: "should error when none of details.support.organization is checked",
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
					"q_1_2": 3,
				},
				Safety: projects.Safety{
					Ready: projects.SafetyReady{
						RunnerInformation: true,
						AED:               true,
						Other:             true,
					},
					AEDCount: 5,
					Addition: "X",
				},
				Route: projects.Route{
					Measurement: projects.RouteMeasurement{
						SelfMeasurement: true,
					},
					Tool: "UU",
					TrafficManagement: projects.TrafficManagement{
						AskPermission: true,
					},
				},
				Judge: projects.Judge{
					Type:      "other",
					OtherType: "Any",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.SupportOrganizationRequiredOneError{},
	},
	{
		name: "should error when details.support.organization.other is checked and details.support.addition is empty",
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
					"q_1_2": 3,
				},
				Safety: projects.Safety{
					Ready: projects.SafetyReady{
						RunnerInformation: true,
						AED:               true,
						Other:             true,
					},
					AEDCount: 5,
					Addition: "X",
				},
				Route: projects.Route{
					Measurement: projects.RouteMeasurement{
						SelfMeasurement: true,
					},
					Tool: "UU",
					TrafficManagement: projects.TrafficManagement{
						AskPermission: true,
					},
				},
				Judge: projects.Judge{
					Type:      "other",
					OtherType: "Any",
				},
				Support: projects.Support{
					Organization: projects.Organization{
						Safety: true,
						Other:  true,
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.SupportAdditionRequiredError{},
	},
	// details.feedback
	{
		name: "should error when details.feedback is empty",
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
					"q_1_2": 3,
				},
				Safety: projects.Safety{
					Ready: projects.SafetyReady{
						RunnerInformation: true,
						AED:               true,
						Other:             true,
					},
					AEDCount: 5,
					Addition: "X",
				},
				Route: projects.Route{
					Measurement: projects.RouteMeasurement{
						SelfMeasurement: true,
					},
					Tool: "UU",
					TrafficManagement: projects.TrafficManagement{
						AskPermission: true,
					},
				},
				Judge: projects.Judge{
					Type:      "other",
					OtherType: "Any",
				},
				Support: projects.Support{
					Organization: projects.Organization{
						Safety: true,
						Other:  true,
					},
					Addition: "X",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.FeedbackRequiredError{},
	},
}
