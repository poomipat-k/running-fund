package projects_test

import (
	"net/http"

	"github.com/poomipat-k/running-fund/pkg/mock"
	"github.com/poomipat-k/running-fund/pkg/projects"
)

var ContactTestCases = []TestCase{
	// contact.projectHead
	{
		name: "should error when contact.projectHead.prefix is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
		},
		store: &mock.MockProjectStore{
			AddProjectFunc: addProjectSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectHeadPrefixRequiredError{},
	},
	{
		name: "should error when contact.projectHead.firstName is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: projects.ProjectHead{
					Prefix: "Mr",
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc: addProjectSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectHeadFirstNameRequiredError{},
	},
	{
		name: "should error when contact.projectHead.lastName is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: projects.ProjectHead{
					Prefix:    "Mr",
					FirstName: "Poomipat",
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc: addProjectSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectHeadLastNameRequiredError{},
	},
	{
		name: "should error when contact.projectHead.organizationPosition is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: projects.ProjectHead{
					Prefix:    "Mr",
					FirstName: "Poomipat",
					LastName:  "Khamai",
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc: addProjectSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectHeadOrganizationPositionRequiredError{},
	},
	{
		name: "should error when contact.projectHead.eventPosition is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: projects.ProjectHead{
					Prefix:               "Mr",
					FirstName:            "Poomipat",
					LastName:             "Khamai",
					OrganizationPosition: "Software Engineer",
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc: addProjectSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectHeadEventPositionRequiredError{},
	},
	// contact.projectManager
	{
		name: "should error when contact.projectManager.prefix is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: projects.ProjectHead{
					Prefix:               "Mr",
					FirstName:            "Poomipat",
					LastName:             "Khamai",
					OrganizationPosition: "Software Engineer",
					EventPosition:        "Head",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc: addProjectSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectManagerPrefixRequiredError{},
	},
	{
		name: "should error when contact.projectManager.firstName is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: projects.ProjectHead{
					Prefix:               "Mr",
					FirstName:            "Poomipat",
					LastName:             "Khamai",
					OrganizationPosition: "Software Engineer",
					EventPosition:        "Head",
				},
				ProjectManager: projects.ProjectManager{
					Prefix: "Mr",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc: addProjectSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectManagerFirstNameRequiredError{},
	},
	{
		name: "should error when contact.projectManager.lastName is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: projects.ProjectHead{
					Prefix:               "Mr",
					FirstName:            "Poomipat",
					LastName:             "Khamai",
					OrganizationPosition: "Software Engineer",
					EventPosition:        "Head",
				},
				ProjectManager: projects.ProjectManager{
					Prefix:    "Mr",
					FirstName: "AA",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc: addProjectSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectManagerLastNameRequiredError{},
	},
	{
		name: "should error when contact.projectManager.organizationPosition is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: projects.ProjectHead{
					Prefix:               "Mr",
					FirstName:            "Poomipat",
					LastName:             "Khamai",
					OrganizationPosition: "Software Engineer",
					EventPosition:        "Head",
				},
				ProjectManager: projects.ProjectManager{
					Prefix:    "Mr",
					FirstName: "AA",
					LastName:  "BB",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc: addProjectSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectManagerOrganizationPositionRequiredError{},
	},
	{
		name: "should error when contact.projectManager.eventPosition is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: projects.ProjectHead{
					Prefix:               "Mr",
					FirstName:            "Poomipat",
					LastName:             "Khamai",
					OrganizationPosition: "Software Engineer",
					EventPosition:        "Head",
				},
				ProjectManager: projects.ProjectManager{
					Prefix:               "Mr",
					FirstName:            "AA",
					LastName:             "BB",
					OrganizationPosition: "COO",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc: addProjectSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectManagerEventPositionRequiredError{},
	},
	// contact.projectCoordinator
	{
		name: "should error when contact.projectCoordinator.prefix is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: projects.ProjectHead{
					Prefix:               "Mr",
					FirstName:            "Poomipat",
					LastName:             "Khamai",
					OrganizationPosition: "Software Engineer",
					EventPosition:        "Head",
				},
				ProjectManager: projects.ProjectManager{
					Prefix:               "Mr",
					FirstName:            "AA",
					LastName:             "BB",
					OrganizationPosition: "COO",
					EventPosition:        "Y",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc: addProjectSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectCoordinatorPrefixRequiredError{},
	},
	{
		name: "should error when contact.projectCoordinator.firstName is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: projects.ProjectHead{
					Prefix:               "Mr",
					FirstName:            "Poomipat",
					LastName:             "Khamai",
					OrganizationPosition: "Software Engineer",
					EventPosition:        "Head",
				},
				ProjectManager: projects.ProjectManager{
					Prefix:               "Mr",
					FirstName:            "AA",
					LastName:             "BB",
					OrganizationPosition: "COO",
					EventPosition:        "Y",
				},
				ProjectCoordinator: projects.ProjectCoordinator{
					Prefix: "Mr",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc: addProjectSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectCoordinatorFirstNameRequiredError{},
	},
	{
		name: "should error when contact.projectCoordinator.lastName is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: projects.ProjectHead{
					Prefix:               "Mr",
					FirstName:            "Poomipat",
					LastName:             "Khamai",
					OrganizationPosition: "Software Engineer",
					EventPosition:        "Head",
				},
				ProjectManager: projects.ProjectManager{
					Prefix:               "Mr",
					FirstName:            "AA",
					LastName:             "BB",
					OrganizationPosition: "COO",
					EventPosition:        "Y",
				},
				ProjectCoordinator: projects.ProjectCoordinator{
					Prefix:    "Mr",
					FirstName: "A",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc: addProjectSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectCoordinatorLastNameRequiredError{},
	},
	{
		name: "should error when contact.projectCoordinator.organizationPosition is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: projects.ProjectHead{
					Prefix:               "Mr",
					FirstName:            "Poomipat",
					LastName:             "Khamai",
					OrganizationPosition: "Software Engineer",
					EventPosition:        "Head",
				},
				ProjectManager: projects.ProjectManager{
					Prefix:               "Mr",
					FirstName:            "AA",
					LastName:             "BB",
					OrganizationPosition: "COO",
					EventPosition:        "Y",
				},
				ProjectCoordinator: projects.ProjectCoordinator{
					Prefix:    "Mr",
					FirstName: "A",
					LastName:  "B",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc: addProjectSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectCoordinatorOrganizationPositionRequiredError{},
	},
	{
		name: "should error when contact.projectCoordinator.eventPosition is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: projects.ProjectHead{
					Prefix:               "Mr",
					FirstName:            "Poomipat",
					LastName:             "Khamai",
					OrganizationPosition: "Software Engineer",
					EventPosition:        "Head",
				},
				ProjectManager: projects.ProjectManager{
					Prefix:               "Mr",
					FirstName:            "AA",
					LastName:             "BB",
					OrganizationPosition: "COO",
					EventPosition:        "Y",
				},
				ProjectCoordinator: projects.ProjectCoordinator{
					Prefix:               "Mr",
					FirstName:            "A",
					LastName:             "B",
					OrganizationPosition: "X",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc: addProjectSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectCoordinatorEventPositionRequiredError{},
	},
	// contact.projectCoordinator.address
	{
		name: "should error when contact.projectCoordinator.address.address is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: projects.ProjectHead{
					Prefix:               "Mr",
					FirstName:            "Poomipat",
					LastName:             "Khamai",
					OrganizationPosition: "Software Engineer",
					EventPosition:        "Head",
				},
				ProjectManager: projects.ProjectManager{
					Prefix:               "Mr",
					FirstName:            "AA",
					LastName:             "BB",
					OrganizationPosition: "COO",
					EventPosition:        "Y",
				},
				ProjectCoordinator: projects.ProjectCoordinator{
					Prefix:               "Mr",
					FirstName:            "A",
					LastName:             "B",
					OrganizationPosition: "X",
					EventPosition:        "Y",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc: addProjectSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectCoordinatorAddressRequiredError{},
	},
	{
		name: "should error when contact.projectCoordinator.address.provinceId is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: projects.ProjectHead{
					Prefix:               "Mr",
					FirstName:            "Poomipat",
					LastName:             "Khamai",
					OrganizationPosition: "Software Engineer",
					EventPosition:        "Head",
				},
				ProjectManager: projects.ProjectManager{
					Prefix:               "Mr",
					FirstName:            "AA",
					LastName:             "BB",
					OrganizationPosition: "COO",
					EventPosition:        "Y",
				},
				ProjectCoordinator: projects.ProjectCoordinator{
					Prefix:               "Mr",
					FirstName:            "A",
					LastName:             "B",
					OrganizationPosition: "X",
					EventPosition:        "Y",
					Address: projects.Address{
						Address: "Test",
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc: addProjectSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectCoordinatorProvinceIdRequiredError{},
	},
	{
		name: "should error when contact.projectCoordinator.address.districtId is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: projects.ProjectHead{
					Prefix:               "Mr",
					FirstName:            "Poomipat",
					LastName:             "Khamai",
					OrganizationPosition: "Software Engineer",
					EventPosition:        "Head",
				},
				ProjectManager: projects.ProjectManager{
					Prefix:               "Mr",
					FirstName:            "AA",
					LastName:             "BB",
					OrganizationPosition: "COO",
					EventPosition:        "Y",
				},
				ProjectCoordinator: projects.ProjectCoordinator{
					Prefix:               "Mr",
					FirstName:            "A",
					LastName:             "B",
					OrganizationPosition: "X",
					EventPosition:        "Y",
					Address: projects.Address{
						Address:    "Test",
						ProvinceId: 1,
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc: addProjectSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectCoordinatorDistrictIdRequiredError{},
	},
	{
		name: "should error when contact.projectCoordinator.address.subdistrictId is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: projects.ProjectHead{
					Prefix:               "Mr",
					FirstName:            "Poomipat",
					LastName:             "Khamai",
					OrganizationPosition: "Software Engineer",
					EventPosition:        "Head",
				},
				ProjectManager: projects.ProjectManager{
					Prefix:               "Mr",
					FirstName:            "AA",
					LastName:             "BB",
					OrganizationPosition: "COO",
					EventPosition:        "Y",
				},
				ProjectCoordinator: projects.ProjectCoordinator{
					Prefix:               "Mr",
					FirstName:            "A",
					LastName:             "B",
					OrganizationPosition: "X",
					EventPosition:        "Y",
					Address: projects.Address{
						Address:    "Test",
						ProvinceId: 1,
						DistrictId: 1,
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc: addProjectSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectCoordinatorSubdistrictIdRequiredError{},
	},
	{
		name: "should error when contact.projectCoordinator.address.postcodeId is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: projects.ProjectHead{
					Prefix:               "Mr",
					FirstName:            "Poomipat",
					LastName:             "Khamai",
					OrganizationPosition: "Software Engineer",
					EventPosition:        "Head",
				},
				ProjectManager: projects.ProjectManager{
					Prefix:               "Mr",
					FirstName:            "AA",
					LastName:             "BB",
					OrganizationPosition: "COO",
					EventPosition:        "Y",
				},
				ProjectCoordinator: projects.ProjectCoordinator{
					Prefix:               "Mr",
					FirstName:            "A",
					LastName:             "B",
					OrganizationPosition: "X",
					EventPosition:        "Y",
					Address: projects.Address{
						Address:       "Test",
						ProvinceId:    1,
						DistrictId:    1,
						SubdistrictId: 1,
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc: addProjectSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectCoordinatorPostcodeIdRequiredError{},
	},
	{
		name: "should error when contact.projectCoordinator.email is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: projects.ProjectHead{
					Prefix:               "Mr",
					FirstName:            "Poomipat",
					LastName:             "Khamai",
					OrganizationPosition: "Software Engineer",
					EventPosition:        "Head",
				},
				ProjectManager: projects.ProjectManager{
					Prefix:               "Mr",
					FirstName:            "AA",
					LastName:             "BB",
					OrganizationPosition: "COO",
					EventPosition:        "Y",
				},
				ProjectCoordinator: projects.ProjectCoordinator{
					Prefix:               "Mr",
					FirstName:            "A",
					LastName:             "B",
					OrganizationPosition: "X",
					EventPosition:        "Y",
					Address: projects.Address{
						Address:       "Test",
						ProvinceId:    1,
						DistrictId:    1,
						SubdistrictId: 1,
						PostcodeId:    2,
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc: addProjectSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectCoordinatorEmailRequiredError{},
	},
	{
		name: "should error when contact.projectCoordinator.lineId is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: projects.ProjectHead{
					Prefix:               "Mr",
					FirstName:            "Poomipat",
					LastName:             "Khamai",
					OrganizationPosition: "Software Engineer",
					EventPosition:        "Head",
				},
				ProjectManager: projects.ProjectManager{
					Prefix:               "Mr",
					FirstName:            "AA",
					LastName:             "BB",
					OrganizationPosition: "COO",
					EventPosition:        "Y",
				},
				ProjectCoordinator: projects.ProjectCoordinator{
					Prefix:               "Mr",
					FirstName:            "A",
					LastName:             "B",
					OrganizationPosition: "X",
					EventPosition:        "Y",
					Address: projects.Address{
						Address:       "Test",
						ProvinceId:    1,
						DistrictId:    1,
						SubdistrictId: 1,
						PostcodeId:    2,
					},
					Email: "abc", // can be free text
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc: addProjectSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectCoordinatorLineIdRequiredError{},
	},
	{
		name: "should error when contact.projectCoordinator.phoneNumber is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: projects.ProjectHead{
					Prefix:               "Mr",
					FirstName:            "Poomipat",
					LastName:             "Khamai",
					OrganizationPosition: "Software Engineer",
					EventPosition:        "Head",
				},
				ProjectManager: projects.ProjectManager{
					Prefix:               "Mr",
					FirstName:            "AA",
					LastName:             "BB",
					OrganizationPosition: "COO",
					EventPosition:        "Y",
				},
				ProjectCoordinator: projects.ProjectCoordinator{
					Prefix:               "Mr",
					FirstName:            "A",
					LastName:             "B",
					OrganizationPosition: "X",
					EventPosition:        "Y",
					Address: projects.Address{
						Address:       "Test",
						ProvinceId:    1,
						DistrictId:    1,
						SubdistrictId: 1,
						PostcodeId:    2,
					},
					Email:  "abc", // can be free text
					LineId: "abcd",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc: addProjectSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectCoordinatorPhoneNumberRequiredError{},
	},
	{
		name: "should error when contact.projectCoordinator.phoneNumber is shorter than 9 numbers",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: projects.ProjectHead{
					Prefix:               "Mr",
					FirstName:            "Poomipat",
					LastName:             "Khamai",
					OrganizationPosition: "Software Engineer",
					EventPosition:        "Head",
				},
				ProjectManager: projects.ProjectManager{
					Prefix:               "Mr",
					FirstName:            "AA",
					LastName:             "BB",
					OrganizationPosition: "COO",
					EventPosition:        "Y",
				},
				ProjectCoordinator: projects.ProjectCoordinator{
					Prefix:               "Mr",
					FirstName:            "A",
					LastName:             "B",
					OrganizationPosition: "X",
					EventPosition:        "Y",
					Address: projects.Address{
						Address:       "Test",
						ProvinceId:    1,
						DistrictId:    1,
						SubdistrictId: 1,
						PostcodeId:    2,
					},
					Email:       "abc", // can be free text
					LineId:      "abcd",
					PhoneNumber: "09912345",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc: addProjectSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectCoordinatorPhoneNumberLengthError{},
	},
	{
		name: "should error when contact.projectCoordinator.phoneNumber is invalid",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: projects.ProjectHead{
					Prefix:               "Mr",
					FirstName:            "Poomipat",
					LastName:             "Khamai",
					OrganizationPosition: "Software Engineer",
					EventPosition:        "Head",
				},
				ProjectManager: projects.ProjectManager{
					Prefix:               "Mr",
					FirstName:            "AA",
					LastName:             "BB",
					OrganizationPosition: "COO",
					EventPosition:        "Y",
				},
				ProjectCoordinator: projects.ProjectCoordinator{
					Prefix:               "Mr",
					FirstName:            "A",
					LastName:             "B",
					OrganizationPosition: "X",
					EventPosition:        "Y",
					Address: projects.Address{
						Address:       "Test",
						ProvinceId:    1,
						DistrictId:    1,
						SubdistrictId: 1,
						PostcodeId:    2,
					},
					Email:       "abc", // can be free text
					LineId:      "abcd",
					PhoneNumber: "099-2131234", // Only numbers allowed
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc: addProjectSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectCoordinatorPhoneNumberInvalidError{},
	},
}
