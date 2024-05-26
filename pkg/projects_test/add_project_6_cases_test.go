package projects_test

import (
	"net/http"

	"github.com/poomipat-k/running-fund/pkg/mock"
	"github.com/poomipat-k/running-fund/pkg/projects"
)

var Attachment = []TestCase{
	{
		name: "should error when marketingFiles is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newTrue(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details:      DetailsOkPayload,
			Experience:   ExperienceOkPayload,
			Fund:         FundOkPayload,
		},
		collaborationFilesPath: "test.png",
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.MarketingFilesRequiredError{},
	},
	{
		name: "should error when routesFiles is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newTrue(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details:      DetailsOkPayload,
			Experience:   ExperienceOkPayload,
			Fund:         FundOkPayload,
		},
		collaborationFilesPath: "test.png",
		marketingFilesPath:     "test.png",
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.RouteFilesRequiredError{},
	},
	{
		name: "should error when eventMapFiles is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newTrue(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details:      DetailsOkPayload,
			Experience:   ExperienceOkPayload,
			Fund:         FundOkPayload,
		},
		collaborationFilesPath: "test.png",
		marketingFilesPath:     "test.png",
		routeFilesPath:         "test.png",
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.EventMapFilesRequiredError{},
	},
	{
		name: "should error when eventDetailsFiles is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newTrue(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details:      DetailsOkPayload,
			Experience:   ExperienceOkPayload,
			Fund:         FundOkPayload,
		},
		collaborationFilesPath: "test.png",
		marketingFilesPath:     "test.png",
		routeFilesPath:         "test.png",
		eventMapFilesPath:      "test.png",
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.EventDetailsFilesRequiredError{},
	},
}
