package projects_test

import (
	"net/http"

	"github.com/poomipat-k/running-fund/pkg/mock"
	"github.com/poomipat-k/running-fund/pkg/projects"
)

var ContactOkPayload = projects.Contact{
	ProjectHead: projects.ContactPerson{
		Prefix:               "Mr",
		FirstName:            "Poomipat",
		LastName:             "Khamai",
		OrganizationPosition: "Software Engineer",
		EventPosition:        "Head",
		Address: projects.Address{
			Address:       "Address 1",
			ProvinceId:    1,
			DistrictId:    1,
			SubdistrictId: 1,
			PostcodeId:    1,
		},
		Email:       "a@test.com",
		LineId:      "c",
		PhoneNumber: "123456789",
	},
	ProjectManager: projects.ContactPerson{
		Prefix:               "Mr",
		FirstName:            "AA",
		LastName:             "BB",
		OrganizationPosition: "COO",
		EventPosition:        "Y",
		Address: projects.Address{
			Address:       "Address 1",
			ProvinceId:    1,
			DistrictId:    1,
			SubdistrictId: 1,
			PostcodeId:    1,
		},
		Email:       "a@test.com",
		LineId:      "c",
		PhoneNumber: "123456789",
	},
	ProjectCoordinator: projects.ContactPerson{
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
		PhoneNumber: "0992131234", // Only numbers allowed
	},
	RaceDirector: projects.RaceDirector{
		Who: "other",
		Alternative: projects.RaceDirectorAlternative{
			Prefix:    "Mr",
			FirstName: "A",
			LastName:  "B",
		},
	},
	Organization: projects.ContactOrganization{
		Name: "government",
		Type: "XX",
	},
}

var ContactTestCases = []TestCase{
	// contact.projectHead
	{
		name: "should error when contact.projectHead.prefix is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
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
				ProjectHead: projects.ContactPerson{
					Prefix: "Mr",
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
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
				ProjectHead: projects.ContactPerson{
					Prefix:    "Mr",
					FirstName: "Poomipat",
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
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
				ProjectHead: projects.ContactPerson{
					Prefix:    "Mr",
					FirstName: "Poomipat",
					LastName:  "Khamai",
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
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
				ProjectHead: projects.ContactPerson{
					Prefix:               "Mr",
					FirstName:            "Poomipat",
					LastName:             "Khamai",
					OrganizationPosition: "Software Engineer",
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectHeadEventPositionRequiredError{},
	},
	{
		name: "should error when contact.projectHead.address.address is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: projects.ContactPerson{
					Prefix:               "Mr",
					FirstName:            "Poomipat",
					LastName:             "Khamai",
					OrganizationPosition: "Software Engineer",
					EventPosition:        "ABC",
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectHeadAddressRequiredError{},
	},
	{
		name: "should error when contact.projectHead.address.provinceId is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: projects.ContactPerson{
					Prefix:               "Mr",
					FirstName:            "Poomipat",
					LastName:             "Khamai",
					OrganizationPosition: "Software Engineer",
					EventPosition:        "ABC",
					Address: projects.Address{
						Address: "Address 1",
					},
				},
			}},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectHeadProvinceIdRequiredError{},
	},
	{
		name: "should error when contact.ProjectHead.address.districtId is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: projects.ContactPerson{
					Prefix:               "Mr",
					FirstName:            "Poomipat",
					LastName:             "Khamai",
					OrganizationPosition: "Software Engineer",
					EventPosition:        "Head",
					Address: projects.Address{
						Address:    "Address 1",
						ProvinceId: 1,
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectHeadDistrictIdRequiredError{},
	},
	{
		name: "should error when contact.ProjectHead.address.subdistrictId is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: projects.ContactPerson{
					Prefix:               "Mr",
					FirstName:            "Poomipat",
					LastName:             "Khamai",
					OrganizationPosition: "Software Engineer",
					EventPosition:        "Head",
					Address: projects.Address{
						Address:    "Address 1",
						ProvinceId: 1,
						DistrictId: 1,
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectHeadSubdistrictIdRequiredError{},
	},
	{
		name: "should error when contact.ProjectHead.address.postcodeId is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: projects.ContactPerson{
					Prefix:               "Mr",
					FirstName:            "Poomipat",
					LastName:             "Khamai",
					OrganizationPosition: "Software Engineer",
					EventPosition:        "Head",
					Address: projects.Address{
						Address:       "Address 1",
						ProvinceId:    1,
						DistrictId:    1,
						SubdistrictId: 1,
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectHeadPostcodeIdRequiredError{},
	},
	{
		name: "should error when contact.ProjectHead.email is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: projects.ContactPerson{
					Prefix:               "Mr",
					FirstName:            "Poomipat",
					LastName:             "Khamai",
					OrganizationPosition: "Software Engineer",
					EventPosition:        "Head",
					Address: projects.Address{
						Address:       "Address 1",
						ProvinceId:    1,
						DistrictId:    1,
						SubdistrictId: 1,
						PostcodeId:    1,
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectHeadEmailRequiredError{},
	},
	{
		name: "should error when contact.ProjectHead.lineId is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: projects.ContactPerson{
					Prefix:               "Mr",
					FirstName:            "Poomipat",
					LastName:             "Khamai",
					OrganizationPosition: "Software Engineer",
					EventPosition:        "Head",
					Address: projects.Address{
						Address:       "Address 1",
						ProvinceId:    1,
						DistrictId:    1,
						SubdistrictId: 1,
						PostcodeId:    1,
					},
					Email: "a@test.com",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectHeadLineIdRequiredError{},
	},
	{
		name: "should error when contact.ProjectHead.phoneNumber is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: projects.ContactPerson{
					Prefix:               "Mr",
					FirstName:            "Poomipat",
					LastName:             "Khamai",
					OrganizationPosition: "Software Engineer",
					EventPosition:        "Head",
					Address: projects.Address{
						Address:       "Address 1",
						ProvinceId:    1,
						DistrictId:    1,
						SubdistrictId: 1,
						PostcodeId:    1,
					},
					Email:  "a@test.com",
					LineId: "c",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectHeadPhoneNumberRequiredError{},
	},
	{
		name: "should error when contact.ProjectHead.phoneNumber is shorter than 9 numbers",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: projects.ContactPerson{
					Prefix:               "Mr",
					FirstName:            "Poomipat",
					LastName:             "Khamai",
					OrganizationPosition: "Software Engineer",
					EventPosition:        "Head",
					Address: projects.Address{
						Address:       "Address 1",
						ProvinceId:    1,
						DistrictId:    1,
						SubdistrictId: 1,
						PostcodeId:    1,
					},
					Email:       "a@test.com",
					LineId:      "c",
					PhoneNumber: "12345678",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectHeadPhoneNumberLengthError{},
	},
	{
		name: "should error when contact.ProjectHead.phoneNumber is invalid",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: projects.ContactPerson{
					Prefix:               "Mr",
					FirstName:            "Poomipat",
					LastName:             "Khamai",
					OrganizationPosition: "Software Engineer",
					EventPosition:        "Head",
					Address: projects.Address{
						Address:       "Address 1",
						ProvinceId:    1,
						DistrictId:    1,
						SubdistrictId: 1,
						PostcodeId:    1,
					},
					Email:       "a@test.com",
					LineId:      "12345678",
					PhoneNumber: "1234120ab2",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectHeadPhoneNumberInvalidError{},
	},
	// contact.projectManager
	{
		name: "should error when contact.projectManager.prefix is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: ContactOkPayload.ProjectHead,
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
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
				ProjectHead: ContactOkPayload.ProjectHead,
				ProjectManager: projects.ContactPerson{
					Prefix: "Mr",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
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
				ProjectHead: ContactOkPayload.ProjectHead,
				ProjectManager: projects.ContactPerson{
					Prefix:    "Mr",
					FirstName: "AA",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
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
				ProjectHead: ContactOkPayload.ProjectHead,
				ProjectManager: projects.ContactPerson{
					Prefix:    "Mr",
					FirstName: "AA",
					LastName:  "BB",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
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
				ProjectHead: ContactOkPayload.ProjectHead,
				ProjectManager: projects.ContactPerson{
					Prefix:               "Mr",
					FirstName:            "AA",
					LastName:             "BB",
					OrganizationPosition: "COO",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectManagerEventPositionRequiredError{},
	},
	{
		name: "should error when contact.ProjectManager.address.address is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: ContactOkPayload.ProjectHead,
				ProjectManager: projects.ContactPerson{
					Prefix:               "Mr",
					FirstName:            "AA",
					LastName:             "BB",
					OrganizationPosition: "COO",
					EventPosition:        "ABC",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectManagerAddressRequiredError{},
	},
	{
		name: "should error when contact.ProjectManager.address.provinceId is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: ContactOkPayload.ProjectHead,
				ProjectManager: projects.ContactPerson{
					Prefix:               "Mr",
					FirstName:            "AA",
					LastName:             "BB",
					OrganizationPosition: "COO",
					EventPosition:        "ABC",
					Address: projects.Address{
						Address: "address a",
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectManagerProvinceIdRequiredError{},
	},
	{
		name: "should error when contact.ProjectManager.address.districtId is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: ContactOkPayload.ProjectHead,
				ProjectManager: projects.ContactPerson{
					Prefix:               "Mr",
					FirstName:            "AA",
					LastName:             "BB",
					OrganizationPosition: "COO",
					EventPosition:        "ABC",
					Address: projects.Address{
						Address:    "address a",
						ProvinceId: 1,
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectManagerDistrictIdRequiredError{},
	},
	{
		name: "should error when contact.ProjectManager.address.subdistrictId is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: ContactOkPayload.ProjectHead,
				ProjectManager: projects.ContactPerson{
					Prefix:               "Mr",
					FirstName:            "AA",
					LastName:             "BB",
					OrganizationPosition: "COO",
					EventPosition:        "ABC",
					Address: projects.Address{
						Address:    "address a",
						ProvinceId: 1,
						DistrictId: 1,
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectManagerSubdistrictIdRequiredError{},
	},
	{
		name: "should error when contact.ProjectManager.address.postcodeId is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: ContactOkPayload.ProjectHead,
				ProjectManager: projects.ContactPerson{
					Prefix:               "Mr",
					FirstName:            "AA",
					LastName:             "BB",
					OrganizationPosition: "COO",
					EventPosition:        "ABC",
					Address: projects.Address{
						Address:       "address a",
						ProvinceId:    1,
						DistrictId:    1,
						SubdistrictId: 1,
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectManagerPostcodeIdRequiredError{},
	},
	{
		name: "should error when contact.ProjectManager.email is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: ContactOkPayload.ProjectHead,
				ProjectManager: projects.ContactPerson{
					Prefix:               "Mr",
					FirstName:            "AA",
					LastName:             "BB",
					OrganizationPosition: "COO",
					EventPosition:        "ABC",
					Address: projects.Address{
						Address:       "address a",
						ProvinceId:    1,
						DistrictId:    1,
						SubdistrictId: 1,
						PostcodeId:    1,
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectManagerEmailRequiredError{},
	},
	{
		name: "should error when contact.ProjectManager.lineId is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: ContactOkPayload.ProjectHead,
				ProjectManager: projects.ContactPerson{
					Prefix:               "Mr",
					FirstName:            "AA",
					LastName:             "BB",
					OrganizationPosition: "COO",
					EventPosition:        "ABC",
					Address: projects.Address{
						Address:       "address a",
						ProvinceId:    1,
						DistrictId:    1,
						SubdistrictId: 1,
						PostcodeId:    1,
					},
					Email: "abc@test.com",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectManagerLineIdRequiredError{},
	},
	{
		name: "should error when contact.ProjectManager.phoneNumber is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: ContactOkPayload.ProjectHead,
				ProjectManager: projects.ContactPerson{
					Prefix:               "Mr",
					FirstName:            "AA",
					LastName:             "BB",
					OrganizationPosition: "COO",
					EventPosition:        "ABC",
					Address: projects.Address{
						Address:       "address a",
						ProvinceId:    1,
						DistrictId:    1,
						SubdistrictId: 1,
						PostcodeId:    1,
					},
					Email:  "abc@test.com",
					LineId: "@test",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectManagerPhoneNumberRequiredError{},
	},
	{
		name: "should error when contact.ProjectManager.phoneNumber is shorter than 9 numbers",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: ContactOkPayload.ProjectHead,
				ProjectManager: projects.ContactPerson{
					Prefix:               "Mr",
					FirstName:            "AA",
					LastName:             "BB",
					OrganizationPosition: "COO",
					EventPosition:        "ABC",
					Address: projects.Address{
						Address:       "address a",
						ProvinceId:    1,
						DistrictId:    1,
						SubdistrictId: 1,
						PostcodeId:    1,
					},
					Email:       "abc@test.com",
					LineId:      "@test",
					PhoneNumber: "12345",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectManagerPhoneNumberLengthError{},
	},
	{
		name: "should error when contact.ProjectManager.phoneNumber is invalid",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead: ContactOkPayload.ProjectHead,
				ProjectManager: projects.ContactPerson{
					Prefix:               "Mr",
					FirstName:            "AA",
					LastName:             "BB",
					OrganizationPosition: "COO",
					EventPosition:        "ABC",
					Address: projects.Address{
						Address:       "address a",
						ProvinceId:    1,
						DistrictId:    1,
						SubdistrictId: 1,
						PostcodeId:    1,
					},
					Email:       "abc@test.com",
					LineId:      "@test",
					PhoneNumber: "12345abadwqw123",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectManagerPhoneNumberInvalidError{},
	},
	// contact.projectCoordinator
	{
		name: "should error when contact.projectCoordinator.prefix is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead:    ContactOkPayload.ProjectHead,
				ProjectManager: ContactOkPayload.ProjectManager,
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
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
				ProjectHead:    ContactOkPayload.ProjectHead,
				ProjectManager: ContactOkPayload.ProjectManager,
				ProjectCoordinator: projects.ContactPerson{
					Prefix: "Mr",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
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
				ProjectHead:    ContactOkPayload.ProjectHead,
				ProjectManager: ContactOkPayload.ProjectManager,
				ProjectCoordinator: projects.ContactPerson{
					Prefix:    "Mr",
					FirstName: "A",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
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
				ProjectHead:    ContactOkPayload.ProjectHead,
				ProjectManager: ContactOkPayload.ProjectManager,
				ProjectCoordinator: projects.ContactPerson{
					Prefix:    "Mr",
					FirstName: "A",
					LastName:  "B",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
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
				ProjectHead:    ContactOkPayload.ProjectHead,
				ProjectManager: ContactOkPayload.ProjectManager,
				ProjectCoordinator: projects.ContactPerson{
					Prefix:               "Mr",
					FirstName:            "A",
					LastName:             "B",
					OrganizationPosition: "X",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
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
				ProjectHead:    ContactOkPayload.ProjectHead,
				ProjectManager: ContactOkPayload.ProjectManager,
				ProjectCoordinator: projects.ContactPerson{
					Prefix:               "Mr",
					FirstName:            "A",
					LastName:             "B",
					OrganizationPosition: "X",
					EventPosition:        "Y",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
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
				ProjectHead:    ContactOkPayload.ProjectHead,
				ProjectManager: ContactOkPayload.ProjectManager,
				ProjectCoordinator: projects.ContactPerson{
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
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
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
				ProjectHead:    ContactOkPayload.ProjectHead,
				ProjectManager: ContactOkPayload.ProjectManager,
				ProjectCoordinator: projects.ContactPerson{
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
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
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
				ProjectHead:    ContactOkPayload.ProjectHead,
				ProjectManager: ContactOkPayload.ProjectManager,
				ProjectCoordinator: projects.ContactPerson{
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
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
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
				ProjectHead:    ContactOkPayload.ProjectHead,
				ProjectManager: ContactOkPayload.ProjectManager,
				ProjectCoordinator: projects.ContactPerson{
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
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
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
				ProjectHead:    ContactOkPayload.ProjectHead,
				ProjectManager: ContactOkPayload.ProjectManager,
				ProjectCoordinator: projects.ContactPerson{
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
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
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
				ProjectHead:    ContactOkPayload.ProjectHead,
				ProjectManager: ContactOkPayload.ProjectManager,
				ProjectCoordinator: projects.ContactPerson{
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
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
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
				ProjectHead:    ContactOkPayload.ProjectHead,
				ProjectManager: ContactOkPayload.ProjectManager,
				ProjectCoordinator: projects.ContactPerson{
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
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
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
				ProjectHead:    ContactOkPayload.ProjectHead,
				ProjectManager: ContactOkPayload.ProjectManager,
				ProjectCoordinator: projects.ContactPerson{
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
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
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
				ProjectHead:    ContactOkPayload.ProjectHead,
				ProjectManager: ContactOkPayload.ProjectManager,
				ProjectCoordinator: projects.ContactPerson{
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
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ProjectCoordinatorPhoneNumberInvalidError{},
	},
	// project.raceDirector
	{
		name: "should error when contact.raceDirector.who is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead:    ContactOkPayload.ProjectHead,
				ProjectManager: ContactOkPayload.ProjectManager,
				ProjectCoordinator: projects.ContactPerson{
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
					PhoneNumber: "0992131234", // Only numbers allowed
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.RaceDirectorWhoRequiredError{},
	},
	{
		name: "should error when contact.raceDirector.who is other and raceDirector.alternative.prefix is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead:    ContactOkPayload.ProjectHead,
				ProjectManager: ContactOkPayload.ProjectManager,
				ProjectCoordinator: projects.ContactPerson{
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
					PhoneNumber: "0992131234", // Only numbers allowed
				},
				RaceDirector: projects.RaceDirector{
					Who: "other",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.RaceDirectorAlternativePrefixRequiredError{},
	},
	{
		name: "should error when contact.raceDirector.who is other and raceDirector.alternative.firstName is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead:    ContactOkPayload.ProjectHead,
				ProjectManager: ContactOkPayload.ProjectManager,
				ProjectCoordinator: projects.ContactPerson{
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
					PhoneNumber: "0992131234", // Only numbers allowed
				},
				RaceDirector: projects.RaceDirector{
					Who: "other",
					Alternative: projects.RaceDirectorAlternative{
						Prefix: "Mr",
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.RaceDirectorAlternativeFirstNameRequiredError{},
	},
	{
		name: "should error when contact.raceDirector.who is other and raceDirector.alternative.lastName is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead:    ContactOkPayload.ProjectHead,
				ProjectManager: ContactOkPayload.ProjectManager,
				ProjectCoordinator: projects.ContactPerson{
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
					PhoneNumber: "0992131234", // Only numbers allowed
				},
				RaceDirector: projects.RaceDirector{
					Who: "other",
					Alternative: projects.RaceDirectorAlternative{
						Prefix:    "Mr",
						FirstName: "A",
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.RaceDirectorAlternativeLastNameRequiredError{},
	},
	// organization
	{
		name: "should error when contact.organization.name is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead:    ContactOkPayload.ProjectHead,
				ProjectManager: ContactOkPayload.ProjectManager,
				ProjectCoordinator: projects.ContactPerson{
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
					PhoneNumber: "0992131234", // Only numbers allowed
				},
				RaceDirector: projects.RaceDirector{
					Who: "other",
					Alternative: projects.RaceDirectorAlternative{
						Prefix:    "Mr",
						FirstName: "A",
						LastName:  "B",
					},
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ContactOrganizationNameRequiredError{},
	},
	{
		name: "should error when contact.organization.type is empty",
		payload: projects.AddProjectRequest{
			Collaborated: newFalse(),
			General:      GeneralDetailsOkPayload,
			Contact: projects.Contact{
				ProjectHead:    ContactOkPayload.ProjectHead,
				ProjectManager: ContactOkPayload.ProjectManager,
				ProjectCoordinator: projects.ContactPerson{
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
					PhoneNumber: "0992131234", // Only numbers allowed
				},
				RaceDirector: projects.RaceDirector{
					Who: "other",
					Alternative: projects.RaceDirectorAlternative{
						Prefix:    "Mr",
						FirstName: "A",
						LastName:  "B",
					},
				},
				Organization: projects.ContactOrganization{
					Name: "government",
				},
			},
		},
		store: &mock.MockProjectStore{
			AddProjectFunc:           addProjectSuccess,
			GetApplicantCriteriaFunc: getApplicantCriteriaSuccess,
		},
		expectedStatus: http.StatusBadRequest,
		expectedError:  &projects.ContactOrganizationTypeRequiredError{},
	},
}
