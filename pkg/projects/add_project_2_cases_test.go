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
					OrganizationPosition: "Senior Manager",
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc: addProjectSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectHeadEventPositionRequiredError{},
	},
}
