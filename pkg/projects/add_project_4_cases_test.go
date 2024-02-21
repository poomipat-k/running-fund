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
}
