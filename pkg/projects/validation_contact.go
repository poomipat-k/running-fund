package projects

func validateContact(payload AddProjectRequest) error {
	if payload.Contact.ProjectHead.Prefix == "" {
		return &ProjectHeadPrefixRequiredError{}
	}
	if payload.Contact.ProjectHead.FirstName == "" {
		return &ProjectHeadFirstNameRequiredError{}
	}
	if payload.Contact.ProjectHead.LastName == "" {
		return &ProjectHeadLastNameRequiredError{}
	}
	if payload.Contact.ProjectHead.OrganizationPosition == "" {
		return &ProjectHeadOrganizationPositionRequiredError{}
	}
	if payload.Contact.ProjectHead.EventPosition == "" {
		return &ProjectHeadEventPositionRequiredError{}
	}
	return nil
}
