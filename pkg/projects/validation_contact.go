package projects

func validateContact(payload AddProjectRequest) error {
	// projectHead
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
	// projectManager
	if payload.Contact.ProjectManager.Prefix == "" {
		return &ProjectManagerPrefixRequiredError{}
	}
	if payload.Contact.ProjectManager.FirstName == "" {
		return &ProjectManagerFirstNameRequiredError{}
	}
	if payload.Contact.ProjectManager.LastName == "" {
		return &ProjectManagerLastNameRequiredError{}
	}
	if payload.Contact.ProjectManager.OrganizationPosition == "" {
		return &ProjectManagerOrganizationPositionRequiredError{}
	}
	if payload.Contact.ProjectManager.EventPosition == "" {
		return &ProjectManagerEventPositionRequiredError{}
	}
	// projectCoordinator
	if payload.Contact.ProjectCoordinator.Prefix == "" {
		return &ProjectCoordinatorPrefixRequiredError{}
	}
	if payload.Contact.ProjectCoordinator.FirstName == "" {
		return &ProjectCoordinatorFirstNameRequiredError{}
	}
	if payload.Contact.ProjectCoordinator.LastName == "" {
		return &ProjectCoordinatorLastNameRequiredError{}
	}
	if payload.Contact.ProjectCoordinator.OrganizationPosition == "" {
		return &ProjectCoordinatorOrganizationPositionRequiredError{}
	}
	if payload.Contact.ProjectCoordinator.EventPosition == "" {
		return &ProjectCoordinatorEventPositionRequiredError{}
	}
	return nil
}
