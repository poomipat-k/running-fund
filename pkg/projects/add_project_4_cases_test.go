package projects_test

import (
	"net/http"

	"github.com/poomipat-k/running-fund/pkg/mock"
	"github.com/poomipat-k/running-fund/pkg/projects"
)

var Experience = []TestCase{
	// thisSeries
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
		name: "should error when experience.thisSeries.firstTime is false and (experience.thisSeries.history.year < 2010 or experience.thisSeries.history.year > currentYear)",
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
						Year:          2005, // 2010 <= validYear <= currentYear
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
	// Complete1
	{
		name: "should error when experience.thisSeries.firstTime is false and experience.thisSeries.history.completed1.year is empty",
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
						Day:           20,
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.CompletedYearRequiredError{},
	},
	{
		name: "should error when experience.thisSeries.firstTime is false and experience.thisSeries.history.completed1.year is invalid",
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
						Day:           20,
						Completed1: projects.HistoryCompleted{
							Year: 2009,
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
		expectedError:  &projects.CompletedYearOutOfBoundError{},
	},
	{
		name: "should error when experience.thisSeries.firstTime is false and experience.thisSeries.history.completed1.name is empty",
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
						Day:           20,
						Completed1: projects.HistoryCompleted{
							Year: 2024,
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
		expectedError:  &projects.CompletedNameRequiredError{},
	},
	{
		name: "should error when experience.thisSeries.firstTime is false and experience.thisSeries.history.completed1.participant is empty",
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
						Day:           20,
						Completed1: projects.HistoryCompleted{
							Year: 2024,
							Name: "XX",
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
		expectedError:  &projects.CompletedParticipantRequiredError{},
	},
	{
		name: "should error when experience.thisSeries.firstTime is false and experience.thisSeries.history.completed1.participant < 0",
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
						Day:           20,
						Completed1: projects.HistoryCompleted{
							Year:        2024,
							Name:        "XX",
							Participant: -10,
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
		expectedError:  &projects.CompletedParticipantInvalidError{},
	},
	// Completed2
	{
		name: "should error when experience.thisSeries.firstTime is false and only experience.thisSeries.history.completed2.year is empty",
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
						Day:           20,
						Completed1: projects.HistoryCompleted{
							Year:        2024,
							Name:        "XX",
							Participant: 100,
						},
						Completed2: projects.HistoryCompleted{
							Participant: 100,
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
		expectedError:  &projects.CompletedYearRequiredError{},
	},
	{
		name: "should error when experience.thisSeries.firstTime is false and only experience.thisSeries.history.completed2.year is invalid",
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
						Day:           20,
						Completed1: projects.HistoryCompleted{
							Year:        2024,
							Name:        "XX",
							Participant: 100,
						},
						Completed2: projects.HistoryCompleted{
							Year:        2000,
							Name:        "x",
							Participant: 100,
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
		expectedError:  &projects.CompletedYearOutOfBoundError{},
	},
	{
		name: "should error when experience.thisSeries.firstTime is false and only experience.thisSeries.history.completed2.name is empty",
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
						Day:           20,
						Completed1: projects.HistoryCompleted{
							Year:        2024,
							Name:        "XX",
							Participant: 100,
						},
						Completed2: projects.HistoryCompleted{
							Year: 2000,
							// Name:        "x",
							Participant: 100,
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
		expectedError:  &projects.CompletedYearOutOfBoundError{},
	},
	{
		name: "should error when experience.thisSeries.firstTime is false and only experience.thisSeries.history.completed2.participant is empty",
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
						Day:           20,
						Completed1: projects.HistoryCompleted{
							Year:        2024,
							Name:        "XX",
							Participant: 100,
						},
						Completed2: projects.HistoryCompleted{
							Year: 2020,
							Name: "x",
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
		expectedError:  &projects.CompletedParticipantRequiredError{},
	},
	{
		name: "should error when experience.thisSeries.firstTime is false and only experience.thisSeries.history.completed2.participant is invalid",
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
						Day:           20,
						Completed1: projects.HistoryCompleted{
							Year:        2024,
							Name:        "XX",
							Participant: 100,
						},
						Completed2: projects.HistoryCompleted{
							Year:        2020,
							Name:        "x",
							Participant: -100,
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
		expectedError:  &projects.CompletedParticipantInvalidError{},
	},
	// Completed3
	{
		name: "should error when experience.thisSeries.firstTime is false and only experience.thisSeries.history.completed3.year is empty",
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
						Day:           20,
						Completed1: projects.HistoryCompleted{
							Year:        2024,
							Name:        "XX",
							Participant: 100,
						},
						Completed3: projects.HistoryCompleted{
							Participant: 100,
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
		expectedError:  &projects.CompletedYearRequiredError{},
	},
	{
		name: "should error when experience.thisSeries.firstTime is false and only experience.thisSeries.history.completed3.year is invalid",
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
						Day:           20,
						Completed1: projects.HistoryCompleted{
							Year:        2024,
							Name:        "XX",
							Participant: 100,
						},
						Completed3: projects.HistoryCompleted{
							Year:        2000,
							Name:        "x",
							Participant: 100,
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
		expectedError:  &projects.CompletedYearOutOfBoundError{},
	},
	{
		name: "should error when experience.thisSeries.firstTime is false and only experience.thisSeries.history.completed3.name is empty",
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
						Day:           20,
						Completed1: projects.HistoryCompleted{
							Year:        2024,
							Name:        "XX",
							Participant: 100,
						},
						Completed3: projects.HistoryCompleted{
							Year: 2000,
							// Name:        "x",
							Participant: 100,
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
		expectedError:  &projects.CompletedYearOutOfBoundError{},
	},
	{
		name: "should error when experience.thisSeries.firstTime is false and only experience.thisSeries.history.completed3.participant is empty",
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
						Day:           20,
						Completed1: projects.HistoryCompleted{
							Year:        2024,
							Name:        "XX",
							Participant: 100,
						},
						Completed3: projects.HistoryCompleted{
							Year: 2020,
							Name: "x",
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
		expectedError:  &projects.CompletedParticipantRequiredError{},
	},
	{
		name: "should error when experience.thisSeries.firstTime is false and only experience.thisSeries.history.completed3.participant is invalid",
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
						Day:           20,
						Completed1: projects.HistoryCompleted{
							Year:        2024,
							Name:        "XX",
							Participant: 100,
						},
						Completed3: projects.HistoryCompleted{
							Year:        2020,
							Name:        "x",
							Participant: -100,
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
		expectedError:  &projects.CompletedParticipantInvalidError{},
	},
	// otherSeries
	{
		name: "should error when experience.otherSeries.doneBefore is empty",
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
						Day:           20,
						Completed1: projects.HistoryCompleted{
							Year:        2024,
							Name:        "XX",
							Participant: 100,
						},
						Completed3: projects.HistoryCompleted{
							Year:        2020,
							Name:        "x",
							Participant: 1100,
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
		expectedError:  &projects.DoneBeforeRequiredError{},
	},
	// otherSeries.history.completed1
	{
		name: "should error when experience.thisSeries.doneBefore is true and experience.thisSeries.history.completed1.year is empty",
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
						Day:           20,
						Completed1: projects.HistoryCompleted{
							Year:        2024,
							Name:        "XX",
							Participant: 100,
						},
						Completed3: projects.HistoryCompleted{
							Year:        2020,
							Name:        "x",
							Participant: 1100,
						},
					},
				},
				OtherSeries: projects.OtherSeries{
					DoneBefore: newTrue(),
					// History: projects.OtherSeriesHistory{
					// 	Completed1: projects.HistoryCompleted{
					// 		Year: 2020,
					// 	},
					// },
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.CompletedYearRequiredError{},
	},
	{
		name: "should error when experience.otherSeries.doneBefore is true and experience.otherSeries.history.completed1.year is invalid",
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
						Day:           20,
						Completed1: projects.HistoryCompleted{
							Year:        2024,
							Name:        "XX",
							Participant: 100,
						},
						Completed3: projects.HistoryCompleted{
							Year:        2020,
							Name:        "x",
							Participant: 1100,
						},
					},
				},
				OtherSeries: projects.OtherSeries{
					DoneBefore: newTrue(),
					History: projects.OtherSeriesHistory{
						Completed1: projects.HistoryCompleted{
							Year: 1999,
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
		expectedError:  &projects.CompletedYearOutOfBoundError{},
	},
	{
		name: "should error when experience.otherSeries.doneBefore is true and experience.otherSeries.history.completed1.name is empty",
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
						Day:           20,
						Completed1: projects.HistoryCompleted{
							Year:        2024,
							Name:        "XX",
							Participant: 100,
						},
						Completed3: projects.HistoryCompleted{
							Year:        2020,
							Name:        "x",
							Participant: 1100,
						},
					},
				},
				OtherSeries: projects.OtherSeries{
					DoneBefore: newTrue(),
					History: projects.OtherSeriesHistory{
						Completed1: projects.HistoryCompleted{
							Year: 2022,
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
		expectedError:  &projects.CompletedNameRequiredError{},
	},
	{
		name: "should error when experience.otherSeries.doneBefore is true and experience.otherSeries.history.completed1.participant is empty",
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
						Day:           20,
						Completed1: projects.HistoryCompleted{
							Year:        2024,
							Name:        "XX",
							Participant: 100,
						},
						Completed3: projects.HistoryCompleted{
							Year:        2020,
							Name:        "x",
							Participant: 1100,
						},
					},
				},
				OtherSeries: projects.OtherSeries{
					DoneBefore: newTrue(),
					History: projects.OtherSeriesHistory{
						Completed1: projects.HistoryCompleted{
							Year: 2022,
							Name: "ABC",
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
		expectedError:  &projects.CompletedParticipantRequiredError{},
	},
	{
		name: "should error when experience.otherSeries.doneBefore is true and experience.otherSeries.history.completed1.participant < 0",
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
						Day:           20,
						Completed1: projects.HistoryCompleted{
							Year:        2024,
							Name:        "XX",
							Participant: 100,
						},
						Completed3: projects.HistoryCompleted{
							Year:        2020,
							Name:        "x",
							Participant: 1100,
						},
					},
				},
				OtherSeries: projects.OtherSeries{
					DoneBefore: newTrue(),
					History: projects.OtherSeriesHistory{
						Completed1: projects.HistoryCompleted{
							Year:        2022,
							Name:        "ABC",
							Participant: -10,
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
		expectedError:  &projects.CompletedParticipantInvalidError{},
	},
	// otherSeries.history.completed2
	{
		name: "should error when experience.otherSeries.doneBefore is true and only experience.otherSeries.history.completed2.year is empty",
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
						Day:           20,
						Completed1: projects.HistoryCompleted{
							Year:        2024,
							Name:        "XX",
							Participant: 100,
						},
						Completed3: projects.HistoryCompleted{
							Year:        2020,
							Name:        "x",
							Participant: 1100,
						},
					},
				},
				OtherSeries: projects.OtherSeries{
					DoneBefore: newTrue(),
					History: projects.OtherSeriesHistory{
						Completed1: projects.HistoryCompleted{
							Year:        2022,
							Name:        "ABC",
							Participant: 3000,
						},
						Completed2: projects.HistoryCompleted{
							Name:        "XX",
							Participant: 300,
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
		expectedError:  &projects.CompletedYearRequiredError{},
	},
	{
		name: "should error when experience.otherSeries.doneBefore is true and only experience.otherSeries.history.completed2.year is invalid",
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
						Day:           20,
						Completed1: projects.HistoryCompleted{
							Year:        2024,
							Name:        "XX",
							Participant: 100,
						},
						Completed3: projects.HistoryCompleted{
							Year:        2020,
							Name:        "x",
							Participant: 1100,
						},
					},
				},
				OtherSeries: projects.OtherSeries{
					DoneBefore: newTrue(),
					History: projects.OtherSeriesHistory{
						Completed1: projects.HistoryCompleted{
							Year:        2022,
							Name:        "ABC",
							Participant: 3000,
						},
						Completed2: projects.HistoryCompleted{
							Year:        2000,
							Name:        "A",
							Participant: 300,
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
		expectedError:  &projects.CompletedYearOutOfBoundError{},
	},
	{
		name: "should error when experience.otherSeries.doneBefore is true and experience.otherSeries.history.completed2.name is empty",
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
						Day:           20,
						Completed1: projects.HistoryCompleted{
							Year:        2024,
							Name:        "XX",
							Participant: 100,
						},
						Completed3: projects.HistoryCompleted{
							Year:        2020,
							Name:        "x",
							Participant: 1100,
						},
					},
				},
				OtherSeries: projects.OtherSeries{
					DoneBefore: newTrue(),
					History: projects.OtherSeriesHistory{
						Completed1: projects.HistoryCompleted{
							Year:        2022,
							Name:        "ABC",
							Participant: 3000,
						},
						Completed2: projects.HistoryCompleted{
							Year: 2020,
							// Name:        "A",
							Participant: 300,
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
		expectedError:  &projects.CompletedNameRequiredError{},
	},
	{
		name: "should error when experience.otherSeries.doneBefore is true and experience.otherSeries.history.completed2.participant is empty",
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
						Day:           20,
						Completed1: projects.HistoryCompleted{
							Year:        2024,
							Name:        "XX",
							Participant: 100,
						},
						Completed3: projects.HistoryCompleted{
							Year:        2020,
							Name:        "x",
							Participant: 1100,
						},
					},
				},
				OtherSeries: projects.OtherSeries{
					DoneBefore: newTrue(),
					History: projects.OtherSeriesHistory{
						Completed1: projects.HistoryCompleted{
							Year:        2022,
							Name:        "ABC",
							Participant: 3000,
						},
						Completed2: projects.HistoryCompleted{
							Year: 2020,
							Name: "A",
							// Participant: 300,
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
		expectedError:  &projects.CompletedParticipantRequiredError{},
	},
	{
		name: "should error when experience.otherSeries.doneBefore is true and experience.otherSeries.history.completed2.participant is invalid",
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
						Day:           20,
						Completed1: projects.HistoryCompleted{
							Year:        2024,
							Name:        "XX",
							Participant: 100,
						},
						Completed3: projects.HistoryCompleted{
							Year:        2020,
							Name:        "x",
							Participant: 1100,
						},
					},
				},
				OtherSeries: projects.OtherSeries{
					DoneBefore: newTrue(),
					History: projects.OtherSeriesHistory{
						Completed1: projects.HistoryCompleted{
							Year:        2022,
							Name:        "ABC",
							Participant: 3000,
						},
						Completed2: projects.HistoryCompleted{
							Year:        2020,
							Name:        "A",
							Participant: -300,
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
		expectedError:  &projects.CompletedParticipantInvalidError{},
	},
	// otherSeries.history.completed3
	{
		name: "should error when experience.otherSeries.doneBefore is true and only experience.otherSeries.history.completed3.year is empty",
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
						Day:           20,
						Completed1: projects.HistoryCompleted{
							Year:        2024,
							Name:        "XX",
							Participant: 100,
						},
						Completed3: projects.HistoryCompleted{
							Year:        2020,
							Name:        "x",
							Participant: 1100,
						},
					},
				},
				OtherSeries: projects.OtherSeries{
					DoneBefore: newTrue(),
					History: projects.OtherSeriesHistory{
						Completed1: projects.HistoryCompleted{
							Year:        2022,
							Name:        "ABC",
							Participant: 3000,
						},
						Completed3: projects.HistoryCompleted{
							Name:        "XX",
							Participant: 300,
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
		expectedError:  &projects.CompletedYearRequiredError{},
	},
	{
		name: "should error when experience.otherSeries.doneBefore is true and only experience.otherSeries.history.completed3.year is invalid",
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
						Day:           20,
						Completed1: projects.HistoryCompleted{
							Year:        2024,
							Name:        "XX",
							Participant: 100,
						},
						Completed3: projects.HistoryCompleted{
							Year:        2020,
							Name:        "x",
							Participant: 1100,
						},
					},
				},
				OtherSeries: projects.OtherSeries{
					DoneBefore: newTrue(),
					History: projects.OtherSeriesHistory{
						Completed1: projects.HistoryCompleted{
							Year:        2022,
							Name:        "ABC",
							Participant: 3000,
						},
						Completed3: projects.HistoryCompleted{
							Year:        2000,
							Name:        "A",
							Participant: 300,
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
		expectedError:  &projects.CompletedYearOutOfBoundError{},
	},
	{
		name: "should error when experience.otherSeries.doneBefore is true and experience.otherSeries.history.completed3.name is empty",
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
						Day:           20,
						Completed1: projects.HistoryCompleted{
							Year:        2024,
							Name:        "XX",
							Participant: 100,
						},
						Completed3: projects.HistoryCompleted{
							Year:        2020,
							Name:        "x",
							Participant: 1100,
						},
					},
				},
				OtherSeries: projects.OtherSeries{
					DoneBefore: newTrue(),
					History: projects.OtherSeriesHistory{
						Completed1: projects.HistoryCompleted{
							Year:        2022,
							Name:        "ABC",
							Participant: 3000,
						},
						Completed3: projects.HistoryCompleted{
							Year: 2020,
							// Name:        "A",
							Participant: 300,
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
		expectedError:  &projects.CompletedNameRequiredError{},
	},
	{
		name: "should error when experience.otherSeries.doneBefore is true and experience.otherSeries.history.completed3.participant is empty",
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
						Day:           20,
						Completed1: projects.HistoryCompleted{
							Year:        2024,
							Name:        "XX",
							Participant: 100,
						},
						Completed3: projects.HistoryCompleted{
							Year:        2020,
							Name:        "x",
							Participant: 1100,
						},
					},
				},
				OtherSeries: projects.OtherSeries{
					DoneBefore: newTrue(),
					History: projects.OtherSeriesHistory{
						Completed1: projects.HistoryCompleted{
							Year:        2022,
							Name:        "ABC",
							Participant: 3000,
						},
						Completed3: projects.HistoryCompleted{
							Year: 2020,
							Name: "A",
							// Participant: 300,
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
		expectedError:  &projects.CompletedParticipantRequiredError{},
	},
	{
		name: "should error when experience.otherSeries.doneBefore is true and experience.otherSeries.history.completed3.participant is invalid",
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
						Day:           20,
						Completed1: projects.HistoryCompleted{
							Year:        2024,
							Name:        "XX",
							Participant: 100,
						},
						Completed3: projects.HistoryCompleted{
							Year:        2020,
							Name:        "x",
							Participant: 1100,
						},
					},
				},
				OtherSeries: projects.OtherSeries{
					DoneBefore: newTrue(),
					History: projects.OtherSeriesHistory{
						Completed1: projects.HistoryCompleted{
							Year:        2022,
							Name:        "ABC",
							Participant: 3000,
						},
						Completed3: projects.HistoryCompleted{
							Year:        2020,
							Name:        "A",
							Participant: -300,
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
		expectedError:  &projects.CompletedParticipantInvalidError{},
	},
}
