package projects_test

import (
	"net/http"

	"github.com/poomipat-k/running-fund/pkg/mock"
	"github.com/poomipat-k/running-fund/pkg/projects"
)

var FundOkPayload = projects.Fund{
	Budget: projects.Budget{
		Total:               50000,
		SupportOrganization: "ABC",
	},
	Request: projects.FundRequest{
		Type: projects.FundRequestType{
			Fund:    true,
			BIB:     true,
			Seminar: true,
			Other:   true,
		},
		Details: projects.FundRequestDetails{
			FundAmount: 50000,
			BibAmount:  500,
			Seminar:    "How to run fast",
			Other:      "Test",
		},
	},
}

var Fund = []TestCase{
	// budget
	{
		name: "should error when fund.budget.total is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details:      DetailsOkPayload,
			Experience:   ExperienceOkPayload,
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.TotalBudgetRequiredError{},
	},
	{
		name: "should error when fund.budget.supportOrganization is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details:      DetailsOkPayload,
			Experience:   ExperienceOkPayload,
			Fund: projects.Fund{
				Budget: projects.Budget{
					Total: 40000,
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.BudgetSupportOrganizationRequiredError{},
	},
	// fund.request
	{
		name: "should error when none of fund.request.type is checked",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details:      DetailsOkPayload,
			Experience:   ExperienceOkPayload,
			Fund: projects.Fund{
				Budget: projects.Budget{
					Total:               50000,
					SupportOrganization: "ABC",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.FundRequestTypeRequiredOneError{},
	},
	{
		name: "should error when fund.request.type.fund is checked and fund.request.details.fundAmount is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details:      DetailsOkPayload,
			Experience:   ExperienceOkPayload,
			Fund: projects.Fund{
				Budget: projects.Budget{
					Total:               50000,
					SupportOrganization: "ABC",
				},
				Request: projects.FundRequest{
					Type: projects.FundRequestType{
						Fund: true,
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.FundRequestAmountRequiredError{},
	},
	{
		name: "should error when fund.request.type.fund is checked and fund.request.details.fundAmount is invalid",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details:      DetailsOkPayload,
			Experience:   ExperienceOkPayload,
			Fund: projects.Fund{
				Budget: projects.Budget{
					Total:               50000,
					SupportOrganization: "ABC",
				},
				Request: projects.FundRequest{
					Type: projects.FundRequestType{
						Fund: true,
					},
					Details: projects.FundRequestDetails{
						FundAmount: 2500,
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.FundRequestAmountInvalidError{},
	},
	{
		name: "should error when fund.request.type.bib is checked and fund.request.details.bibAmount is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details:      DetailsOkPayload,
			Experience:   ExperienceOkPayload,
			Fund: projects.Fund{
				Budget: projects.Budget{
					Total:               50000,
					SupportOrganization: "ABC",
				},
				Request: projects.FundRequest{
					Type: projects.FundRequestType{
						Fund: true,
						BIB:  true,
					},
					Details: projects.FundRequestDetails{
						FundAmount: 50000,
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.BibRequestAmountRequiredError{},
	},
	{
		name: "should error when fund.request.type.bib is checked and fund.request.details.bibAmount is invalid",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details:      DetailsOkPayload,
			Experience:   ExperienceOkPayload,
			Fund: projects.Fund{
				Budget: projects.Budget{
					Total:               50000,
					SupportOrganization: "ABC",
				},
				Request: projects.FundRequest{
					Type: projects.FundRequestType{
						Fund: true,
						BIB:  true,
					},
					Details: projects.FundRequestDetails{
						FundAmount: 50000,
						BibAmount:  -50,
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.BibRequestAmountInvalidError{},
	},
	{
		name: "should error when fund.request.type.seminar is checked and fund.request.details.seminar is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details:      DetailsOkPayload,
			Experience:   ExperienceOkPayload,
			Fund: projects.Fund{
				Budget: projects.Budget{
					Total:               50000,
					SupportOrganization: "ABC",
				},
				Request: projects.FundRequest{
					Type: projects.FundRequestType{
						Fund:    true,
						BIB:     true,
						Seminar: true,
					},
					Details: projects.FundRequestDetails{
						FundAmount: 50000,
						BibAmount:  500,
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.SeminarRequiredError{},
	},
	{
		name: "should error when fund.request.type.other is checked and fund.request.details.other is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact:      ContactOkPayload,
			Details:      DetailsOkPayload,
			Experience:   ExperienceOkPayload,
			Fund: projects.Fund{
				Budget: projects.Budget{
					Total:               50000,
					SupportOrganization: "ABC",
				},
				Request: projects.FundRequest{
					Type: projects.FundRequestType{
						Fund:    true,
						BIB:     true,
						Seminar: true,
						Other:   true,
					},
					Details: projects.FundRequestDetails{
						FundAmount: 50000,
						BibAmount:  500,
						Seminar:    "How to run fast",
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.OtherRequestRequiredError{},
	},
}
