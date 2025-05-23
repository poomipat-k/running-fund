package projects_test

import (
	"net/http"

	"github.com/poomipat-k/running-fund/pkg/mock"
	"github.com/poomipat-k/running-fund/pkg/projects"
)

var GeneralDetailsOkPayload = projects.AddProjectGeneralDetails{
	ProjectName: "A",
	EventDate: projects.EventDate{
		Year:       2024,
		Month:      2,
		Day:        20,
		FromHour:   newInt(0),
		FromMinute: newInt(25),
		ToHour:     newInt(10),
		ToMinute:   newInt(20),
	},
	Address: projects.Address{
		Address:       "A",
		ProvinceId:    1,
		DistrictId:    2,
		SubdistrictId: 3,
		PostcodeId:    4,
	},
	StartPoint:  "X",
	FinishPoint: "Y",
	EventDetails: projects.EventDetails{
		Category: projects.Category{
			Available: projects.Available{
				Other:        true,
				RoadRace:     false,
				TrailRunning: true,
			},
			OtherType: "Tri",
		},
		DistanceAndFee: []projects.DistanceAndFee{
			{Checked: true, Type: "half", Fee: newFloat64(330), Dynamic: newFalse()},
			{Checked: true, Type: "iron man", Fee: newFloat64(1000), Dynamic: newTrue()},
		},
		VIP: newFalse(),
	},
	ExpectedParticipants: "3501-4500",
	HasOrganizer:         newTrue(),
	OrganizerName:        "ABC",
}

var GeneralAndCollaboratedTestCases = []TestCase{
	/*
		STEP 0 START (collaborated)
	*/
	{
		name:    "should error collaborated required",
		payload: projects.AddProjectRequest{},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.CollaboratedRequiredError{},
	},
	{
		name: "should not error when collaborated and collaborateFiles is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newTrue()},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectNameRequiredError{},
	},
	// STEP 0 END
	// 1 START - general
	{
		name: "should error when general.projectName is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse()},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectNameRequiredError{},
	},
	// general.eventDate
	{
		name: "should error when general.eventDate.year is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.YearRequiredError{},
	},
	{
		name: "should error when general.eventDate.year is less than 1971",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year: 1969,
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.YearInvalidError{},
	},
	{
		name: "should error when general.eventDate.month is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year: 2024,
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.MonthRequiredError{},
	},
	{
		name: "should error when general.eventDate.month is less than 1 or > 12",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:  2023,
					Month: -1,
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.MonthOutOfBoundError{},
	},
	{
		name: "should error when general.eventDate.day is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:  2024,
					Month: 2,
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.DayRequiredError{},
	},
	{
		name: "should error when general.eventDate.day is 32",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:  2023,
					Month: 1,
					Day:   32,
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.DayOutOfBoundError{},
	},
	{
		name: "should error when general.eventDate.day is 29 Feb on non-leap-year",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:  2023,
					Month: 2,
					Day:   29,
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.DayOutOfBoundError{},
	},
	{
		name: "should error when general.eventDate.day is 31 November",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:  2023,
					Month: 11,
					Day:   31,
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.DayOutOfBoundError{},
	},
	{
		name: "should error when general.eventDate.day is 30 Feb",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:  2023,
					Month: 2,
					Day:   30,
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.DayOutOfBoundError{},
	},
	{
		name: "should error when general.eventDate.fromHour is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:  2024,
					Month: 2,
					Day:   20,
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.FromHourRequiredError{},
	},
	{
		name: "should error when general.eventDate.fromHour is < 0 or > 23",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:     2024,
					Month:    2,
					Day:      20,
					FromHour: newInt(24),
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.InvalidError{Name: "fromHour"},
	},
	{
		name: "should error when general.eventDate.fromMinute is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:     2024,
					Month:    2,
					Day:      20,
					FromHour: newInt(0),
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.FromMinuteRequiredError{},
	},
	{
		name: "should error when general.eventDate.fromMinute < 0 or > 59",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:       2024,
					Month:      2,
					Day:        20,
					FromHour:   newInt(0),
					FromMinute: newInt(60),
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.InvalidError{Name: "fromMinute"},
	},
	{
		name: "should error when general.eventDate.toHour is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:       2024,
					Month:      2,
					Day:        20,
					FromHour:   newInt(0),
					FromMinute: newInt(0),
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ToHourRequiredError{},
	},
	{
		name: "should error when general.eventDate.ToHour < 0 or > 23",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:       2024,
					Month:      2,
					Day:        20,
					FromHour:   newInt(0),
					FromMinute: newInt(25),
					ToHour:     newInt(25),
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.InvalidError{Name: "toHour"},
	},
	{
		name: "should error when general.eventDate.toMinute is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:       2024,
					Month:      2,
					Day:        20,
					FromHour:   newInt(0),
					FromMinute: newInt(25),
					ToHour:     newInt(10),
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ToMinuteRequiredError{},
	},
	{
		name: "should error when general.eventDate.toMinute < 0 or > 59",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:       2024,
					Month:      2,
					Day:        20,
					FromHour:   newInt(0),
					FromMinute: newInt(25),
					ToHour:     newInt(10),
					ToMinute:   newInt(70),
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.InvalidError{Name: "toMinute"},
	},

	// general.address
	{
		name: "should error when general.address.address is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:       2024,
					Month:      2,
					Day:        20,
					FromHour:   newInt(0),
					FromMinute: newInt(25),
					ToHour:     newInt(10),
					ToMinute:   newInt(20),
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.GeneralAddressRequiredError{},
	},
	{
		name: "should error when general.address.provinceId is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:       2024,
					Month:      2,
					Day:        20,
					FromHour:   newInt(0),
					FromMinute: newInt(25),
					ToHour:     newInt(10),
					ToMinute:   newInt(20),
				},
				Address: projects.Address{
					Address: "A",
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.GeneralProvinceRequiredError{},
	},
	{
		name: "should error when general.address.districtId is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:       2024,
					Month:      2,
					Day:        20,
					FromHour:   newInt(0),
					FromMinute: newInt(25),
					ToHour:     newInt(10),
					ToMinute:   newInt(20),
				},
				Address: projects.Address{
					Address:    "A",
					ProvinceId: 1,
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.GeneralDistrictIdRequiredError{},
	},
	{
		name: "should error when general.address.subdistrictId is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:       2024,
					Month:      2,
					Day:        20,
					FromHour:   newInt(0),
					FromMinute: newInt(25),
					ToHour:     newInt(10),
					ToMinute:   newInt(20),
				},
				Address: projects.Address{
					Address:    "A",
					ProvinceId: 1,
					DistrictId: 2,
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.GeneralSubdistrictIdRequiredError{},
	},
	{
		name: "should error when general.address.postcodeId is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:       2024,
					Month:      2,
					Day:        20,
					FromHour:   newInt(0),
					FromMinute: newInt(25),
					ToHour:     newInt(10),
					ToMinute:   newInt(20),
				},
				Address: projects.Address{
					Address:       "A",
					ProvinceId:    1,
					DistrictId:    2,
					SubdistrictId: 3,
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.GeneralPostcodeIdRequiredError{},
	},
	// general.startPoint and general.finishPoint
	{
		name: "should error when general.startPoint is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:       2024,
					Month:      2,
					Day:        20,
					FromHour:   newInt(0),
					FromMinute: newInt(25),
					ToHour:     newInt(10),
					ToMinute:   newInt(20),
				},
				Address: projects.Address{
					Address:       "A",
					ProvinceId:    1,
					DistrictId:    2,
					SubdistrictId: 3,
					PostcodeId:    4,
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.StartPointRequiredError{},
	},
	{
		name: "should error when general.finishPoint is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:       2024,
					Month:      2,
					Day:        20,
					FromHour:   newInt(0),
					FromMinute: newInt(25),
					ToHour:     newInt(10),
					ToMinute:   newInt(20),
				},
				Address: projects.Address{
					Address:       "A",
					ProvinceId:    1,
					DistrictId:    2,
					SubdistrictId: 3,
					PostcodeId:    4,
				},
				StartPoint: "X",
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.FinishPointRequiredError{},
	},
	// general.eventDetails.category
	{
		name: "should error when none of general.category.available is selected",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:       2024,
					Month:      2,
					Day:        20,
					FromHour:   newInt(0),
					FromMinute: newInt(25),
					ToHour:     newInt(10),
					ToMinute:   newInt(20),
				},
				Address: projects.Address{
					Address:       "A",
					ProvinceId:    1,
					DistrictId:    2,
					SubdistrictId: 3,
					PostcodeId:    4,
				},
				StartPoint:  "X",
				FinishPoint: "Y",
				EventDetails: projects.EventDetails{
					Category: projects.Category{
						Available: projects.Available{
							Other:        false,
							RoadRace:     false,
							TrailRunning: false,
						},
					},
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.CategoryAvailableRequiredOneError{},
	},
	{
		name: "should error when general.eventDetails.category.available.other is checked but general.eventDetails.category.otherType is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:       2024,
					Month:      2,
					Day:        20,
					FromHour:   newInt(0),
					FromMinute: newInt(25),
					ToHour:     newInt(10),
					ToMinute:   newInt(20),
				},
				Address: projects.Address{
					Address:       "A",
					ProvinceId:    1,
					DistrictId:    2,
					SubdistrictId: 3,
					PostcodeId:    4,
				},
				StartPoint:  "X",
				FinishPoint: "Y",
				EventDetails: projects.EventDetails{
					Category: projects.Category{
						Available: projects.Available{
							Other:        true,
							RoadRace:     false,
							TrailRunning: false,
						},
					},
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.OtherEventTypeRequiredError{},
	},
	// general.eventDetails.distanceAndFee
	{
		name: "should error when general.eventDetails.distanceAndFee len is 0",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:       2024,
					Month:      2,
					Day:        20,
					FromHour:   newInt(0),
					FromMinute: newInt(25),
					ToHour:     newInt(10),
					ToMinute:   newInt(20),
				},
				Address: projects.Address{
					Address:       "A",
					ProvinceId:    1,
					DistrictId:    2,
					SubdistrictId: 3,
					PostcodeId:    4,
				},
				StartPoint:  "X",
				FinishPoint: "Y",
				EventDetails: projects.EventDetails{
					Category: projects.Category{
						Available: projects.Available{
							Other:        false,
							RoadRace:     false,
							TrailRunning: true,
						},
					},
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.DistanceRequiredOneError{},
	},
	{
		name: "should error when general.eventDetails.distanceAndFee has 0 checked",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:       2024,
					Month:      2,
					Day:        20,
					FromHour:   newInt(0),
					FromMinute: newInt(25),
					ToHour:     newInt(10),
					ToMinute:   newInt(20),
				},
				Address: projects.Address{
					Address:       "A",
					ProvinceId:    1,
					DistrictId:    2,
					SubdistrictId: 3,
					PostcodeId:    4,
				},
				StartPoint:  "X",
				FinishPoint: "Y",
				EventDetails: projects.EventDetails{
					Category: projects.Category{
						Available: projects.Available{
							Other:        false,
							RoadRace:     false,
							TrailRunning: true,
						},
					},
					DistanceAndFee: []projects.DistanceAndFee{
						{Checked: false, Type: "fun", Fee: newFloat64(222.75), Dynamic: newFalse()},
					},
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.DistanceRequiredOneError{},
	},
	{
		name: "should error when general.eventDetails.distanceAndFee.type is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:       2024,
					Month:      2,
					Day:        20,
					FromHour:   newInt(0),
					FromMinute: newInt(25),
					ToHour:     newInt(10),
					ToMinute:   newInt(20),
				},
				Address: projects.Address{
					Address:       "A",
					ProvinceId:    1,
					DistrictId:    2,
					SubdistrictId: 3,
					PostcodeId:    4,
				},
				StartPoint:  "X",
				FinishPoint: "Y",
				EventDetails: projects.EventDetails{
					Category: projects.Category{
						Available: projects.Available{
							Other:        false,
							RoadRace:     false,
							TrailRunning: true,
						},
					},
					DistanceAndFee: []projects.DistanceAndFee{
						{Checked: true, Type: "", Fee: newFloat64(222.50), Dynamic: newFalse()},
					},
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.DistanceTypeRequiredError{},
	},
	{
		name: "should error when general.eventDetails.distanceAndFee.fee is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:       2024,
					Month:      2,
					Day:        20,
					FromHour:   newInt(0),
					FromMinute: newInt(25),
					ToHour:     newInt(10),
					ToMinute:   newInt(20),
				},
				Address: projects.Address{
					Address:       "A",
					ProvinceId:    1,
					DistrictId:    2,
					SubdistrictId: 3,
					PostcodeId:    4,
				},
				StartPoint:  "X",
				FinishPoint: "Y",
				EventDetails: projects.EventDetails{
					Category: projects.Category{
						Available: projects.Available{
							Other:        false,
							RoadRace:     false,
							TrailRunning: true,
						},
					},
					DistanceAndFee: []projects.DistanceAndFee{
						{Checked: true, Type: "half", Dynamic: newFalse()},
					},
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.DistanceFeeRequiredError{},
	},
	{
		name: "should error when general.eventDetails.distanceAndFee.fee is negative",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:       2024,
					Month:      2,
					Day:        20,
					FromHour:   newInt(0),
					FromMinute: newInt(25),
					ToHour:     newInt(10),
					ToMinute:   newInt(20),
				},
				Address: projects.Address{
					Address:       "A",
					ProvinceId:    1,
					DistrictId:    2,
					SubdistrictId: 3,
					PostcodeId:    4,
				},
				StartPoint:  "X",
				FinishPoint: "Y",
				EventDetails: projects.EventDetails{
					Category: projects.Category{
						Available: projects.Available{
							Other:        false,
							RoadRace:     false,
							TrailRunning: true,
						},
					},
					DistanceAndFee: []projects.DistanceAndFee{
						{Checked: true, Type: "half", Fee: newFloat64(-30), Dynamic: newFalse()},
					},
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ValueNegativeError{},
	},
	{
		name: "should error when general.eventDetails.distanceAndFee.dynamic is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:       2024,
					Month:      2,
					Day:        20,
					FromHour:   newInt(0),
					FromMinute: newInt(25),
					ToHour:     newInt(10),
					ToMinute:   newInt(20),
				},
				Address: projects.Address{
					Address:       "A",
					ProvinceId:    1,
					DistrictId:    2,
					SubdistrictId: 3,
					PostcodeId:    4,
				},
				StartPoint:  "X",
				FinishPoint: "Y",
				EventDetails: projects.EventDetails{
					Category: projects.Category{
						Available: projects.Available{
							Other:        false,
							RoadRace:     false,
							TrailRunning: true,
						},
					},
					DistanceAndFee: []projects.DistanceAndFee{
						{Checked: true, Type: "half", Fee: newFloat64(330)},
					},
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.DistanceAndFeeDynamicRequiredError{},
	},
	// general.eventDetails.vip
	{
		name: "should error when general.eventDetails.vip is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:       2024,
					Month:      2,
					Day:        20,
					FromHour:   newInt(0),
					FromMinute: newInt(25),
					ToHour:     newInt(10),
					ToMinute:   newInt(20),
				},
				Address: projects.Address{
					Address:       "A",
					ProvinceId:    1,
					DistrictId:    2,
					SubdistrictId: 3,
					PostcodeId:    4,
				},
				StartPoint:  "X",
				FinishPoint: "Y",
				EventDetails: projects.EventDetails{
					Category: projects.Category{
						Available: projects.Available{
							Other:        false,
							RoadRace:     false,
							TrailRunning: true,
						},
					},
					DistanceAndFee: []projects.DistanceAndFee{
						{Checked: true, Type: "half", Fee: newFloat64(330), Dynamic: newTrue()},
					},
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.VIPRequiredError{},
	},
	{
		name: "should error when general.eventDetails.vip is true and general.eventDetails.vipFee is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:       2024,
					Month:      2,
					Day:        20,
					FromHour:   newInt(0),
					FromMinute: newInt(25),
					ToHour:     newInt(10),
					ToMinute:   newInt(20),
				},
				Address: projects.Address{
					Address:       "A",
					ProvinceId:    1,
					DistrictId:    2,
					SubdistrictId: 3,
					PostcodeId:    4,
				},
				StartPoint:  "X",
				FinishPoint: "Y",
				EventDetails: projects.EventDetails{
					Category: projects.Category{
						Available: projects.Available{
							Other:        false,
							RoadRace:     false,
							TrailRunning: true,
						},
					},
					DistanceAndFee: []projects.DistanceAndFee{
						{Checked: true, Type: "half", Fee: newFloat64(330), Dynamic: newTrue()},
					},
					VIP: newTrue(),
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.VIPFeeRequiredError{},
	},
	{
		name: "should error when general.eventDetails.vip is true and general.eventDetails.vipFee is negative",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:       2024,
					Month:      2,
					Day:        20,
					FromHour:   newInt(0),
					FromMinute: newInt(25),
					ToHour:     newInt(10),
					ToMinute:   newInt(20),
				},
				Address: projects.Address{
					Address:       "A",
					ProvinceId:    1,
					DistrictId:    2,
					SubdistrictId: 3,
					PostcodeId:    4,
				},
				StartPoint:  "X",
				FinishPoint: "Y",
				EventDetails: projects.EventDetails{
					Category: projects.Category{
						Available: projects.Available{
							Other:        false,
							RoadRace:     false,
							TrailRunning: true,
						},
					},
					DistanceAndFee: []projects.DistanceAndFee{
						{Checked: true, Type: "half", Fee: newFloat64(330), Dynamic: newTrue()},
					},
					VIP:    newTrue(),
					VIPFee: newFloat64(-1.234),
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.VIPFeeNegativeError{},
	},
	// general.expectedParticipants
	{
		name: "should error when general.expectedParticipants is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:       2024,
					Month:      2,
					Day:        20,
					FromHour:   newInt(0),
					FromMinute: newInt(25),
					ToHour:     newInt(10),
					ToMinute:   newInt(20),
				},
				Address: projects.Address{
					Address:       "A",
					ProvinceId:    1,
					DistrictId:    2,
					SubdistrictId: 3,
					PostcodeId:    4,
				},
				StartPoint:  "X",
				FinishPoint: "Y",
				EventDetails: projects.EventDetails{
					Category: projects.Category{
						Available: projects.Available{
							Other:        false,
							RoadRace:     false,
							TrailRunning: true,
						},
					},
					DistanceAndFee: []projects.DistanceAndFee{
						{Checked: true, Type: "half", Fee: newFloat64(330), Dynamic: newTrue()},
					},
					VIP: newFalse(),
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ExpectedParticipantsRequiredError{},
	},
	{
		name: "should error when general.expectedParticipants is not a valid value",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:       2024,
					Month:      2,
					Day:        20,
					FromHour:   newInt(0),
					FromMinute: newInt(25),
					ToHour:     newInt(10),
					ToMinute:   newInt(20),
				},
				Address: projects.Address{
					Address:       "A",
					ProvinceId:    1,
					DistrictId:    2,
					SubdistrictId: 3,
					PostcodeId:    4,
				},
				StartPoint:  "X",
				FinishPoint: "Y",
				EventDetails: projects.EventDetails{
					Category: projects.Category{
						Available: projects.Available{
							Other:        false,
							RoadRace:     false,
							TrailRunning: true,
						},
					},
					DistanceAndFee: []projects.DistanceAndFee{
						{Checked: true, Type: "half", Fee: newFloat64(330), Dynamic: newTrue()},
					},
					VIP: newFalse(),
				},
				ExpectedParticipants: ">100000",
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ExpectedParticipantsInvalidError{},
	},
	// general.hasOrganizer
	{
		name: "should error when general.hasOrganizer is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:       2024,
					Month:      2,
					Day:        20,
					FromHour:   newInt(0),
					FromMinute: newInt(25),
					ToHour:     newInt(10),
					ToMinute:   newInt(20),
				},
				Address: projects.Address{
					Address:       "A",
					ProvinceId:    1,
					DistrictId:    2,
					SubdistrictId: 3,
					PostcodeId:    4,
				},
				StartPoint:  "X",
				FinishPoint: "Y",
				EventDetails: projects.EventDetails{
					Category: projects.Category{
						Available: projects.Available{
							Other:        false,
							RoadRace:     false,
							TrailRunning: true,
						},
					},
					DistanceAndFee: []projects.DistanceAndFee{
						{Checked: true, Type: "half", Fee: newFloat64(330), Dynamic: newTrue()},
					},
					VIP: newFalse(),
				},
				ExpectedParticipants: "3501-4500",
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.HasOrganizerRequiredError{},
	},
	{
		name: "should error when general.hasOrganizer is true and organizerName is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General: projects.AddProjectGeneralDetails{
				ProjectName: "A",
				EventDate: projects.EventDate{
					Year:       2024,
					Month:      2,
					Day:        20,
					FromHour:   newInt(0),
					FromMinute: newInt(25),
					ToHour:     newInt(10),
					ToMinute:   newInt(20),
				},
				Address: projects.Address{
					Address:       "A",
					ProvinceId:    1,
					DistrictId:    2,
					SubdistrictId: 3,
					PostcodeId:    4,
				},
				StartPoint:  "X",
				FinishPoint: "Y",
				EventDetails: projects.EventDetails{
					Category: projects.Category{
						Available: projects.Available{
							Other:        false,
							RoadRace:     false,
							TrailRunning: true,
						},
					},
					DistanceAndFee: []projects.DistanceAndFee{
						{Checked: true, Type: "half", Fee: newFloat64(330), Dynamic: newTrue()},
					},
					VIP: newFalse(),
				},
				ExpectedParticipants: "3501-4500",
				HasOrganizer:         newTrue(),
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.OrganizerNameRequiredError{},
	},

	// 1 END - general
}
