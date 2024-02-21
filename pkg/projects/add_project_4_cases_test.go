package projects_test

import (
	"net/http"

	"github.com/poomipat-k/running-fund/pkg/mock"
	"github.com/poomipat-k/running-fund/pkg/projects"
)

var Experience = []TestCase{
	{
		name: "should error when experience.thisSeries.firstTime is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details:      DetailsOkPayload,
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ThisSeriesFirstTimeRequiredError{},
	},
	{
		name: "should error when experience.thisSeries.firstTime false and ordinalNumber is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details:      DetailsOkPayload,
			Experience: projects.Experience{
				ThisSeries: projects.ThisSeries{
					FirstTime: newFalse(),
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.HistoryOrdinalNumberInvalidError{},
	},
	{
		name: "should error when experience.thisSeries.firstTime is false and experience.thisSeries.history.ordinalNumber < 2",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details:      DetailsOkPayload,
			Experience: projects.Experience{
				ThisSeries: projects.ThisSeries{
					FirstTime: newFalse(),
					History: projects.ThisSeriesHistory{
						OrdinalNumber: 1,
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.HistoryOrdinalNumberInvalidError{},
	},
	{
		name: "should error when experience.thisSeries.firstTime is false and experience.thisSeries.history.year is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details:      DetailsOkPayload,
			Experience: projects.Experience{
				ThisSeries: projects.ThisSeries{
					FirstTime: newFalse(),
					History: projects.ThisSeriesHistory{
						OrdinalNumber: 2,
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.HistoryYearRequiredError{},
	},
	{
		name: "should error when experience.thisSeries.firstTime is false and (experience.thisSeries.history.year < 2018 or experience.thisSeries.history.year > currentYear)",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details:      DetailsOkPayload,
			Experience: projects.Experience{
				ThisSeries: projects.ThisSeries{
					FirstTime: newFalse(),
					History: projects.ThisSeriesHistory{
						OrdinalNumber: 2,
						Year:          2016, // 2018 <= validYear <= currentYear
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.HistoryYearOutOfBoundError{},
	},
	{
		name: "should error when experience.thisSeries.firstTime is false and experience.thisSeries.history.month is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details:      DetailsOkPayload,
			Experience: projects.Experience{
				ThisSeries: projects.ThisSeries{
					FirstTime: newFalse(),
					History: projects.ThisSeriesHistory{
						OrdinalNumber: 2,
						Year:          2020,
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.HistoryMonthRequiredError{},
	},
	{
		name: "should error when experience.thisSeries.firstTime is false and experience.thisSeries.history.month is invalid",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details:      DetailsOkPayload,
			Experience: projects.Experience{
				ThisSeries: projects.ThisSeries{
					FirstTime: newFalse(),
					History: projects.ThisSeriesHistory{
						OrdinalNumber: 2,
						Year:          2020,
						Month:         13,
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.HistoryMonthOutOfBoundError{},
	},
	{
		name: "should error when experience.thisSeries.firstTime is false and experience.thisSeries.history.day is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details:      DetailsOkPayload,
			Experience: projects.Experience{
				ThisSeries: projects.ThisSeries{
					FirstTime: newFalse(),
					History: projects.ThisSeriesHistory{
						OrdinalNumber: 2,
						Year:          2020,
						Month:         2,
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.HistoryDayRequiredError{},
	},
	{
		name: "should error when experience.thisSeries.firstTime is false and experience.thisSeries.history.day is 32 Jan",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details:      DetailsOkPayload,
			Experience: projects.Experience{
				ThisSeries: projects.ThisSeries{
					FirstTime: newFalse(),
					History: projects.ThisSeriesHistory{
						OrdinalNumber: 2,
						Year:          2020,
						Month:         1,
						Day:           32,
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.HistoryDayOutOfBoundError{},
	},
	{
		name: "should error when experience.thisSeries.firstTime is false and experience.thisSeries.history.day is 31 November",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details:      DetailsOkPayload,
			Experience: projects.Experience{
				ThisSeries: projects.ThisSeries{
					FirstTime: newFalse(),
					History: projects.ThisSeriesHistory{
						OrdinalNumber: 2,
						Year:          2020,
						Month:         11,
						Day:           31,
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.HistoryDayOutOfBoundError{},
	},
	{
		name: "should error when experience.thisSeries.firstTime is false and experience.thisSeries.history.day is 30 Feb",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details:      DetailsOkPayload,
			Experience: projects.Experience{
				ThisSeries: projects.ThisSeries{
					FirstTime: newFalse(),
					History: projects.ThisSeriesHistory{
						OrdinalNumber: 2,
						Year:          2020,
						Month:         2,
						Day:           30,
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.HistoryDayOutOfBoundError{},
	},
	{
		name: "should error when experience.thisSeries.firstTime is false and experience.thisSeries.history.day is 29 Feb on none-leap-year",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details:      DetailsOkPayload,
			Experience: projects.Experience{
				ThisSeries: projects.ThisSeries{
					FirstTime: newFalse(),
					History: projects.ThisSeriesHistory{
						OrdinalNumber: 2,
						Year:          2023,
						Month:         2,
						Day:           29,
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.HistoryDayOutOfBoundError{},
	},
}
