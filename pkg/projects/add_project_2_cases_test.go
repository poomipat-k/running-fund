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
}
