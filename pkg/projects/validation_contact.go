package projects

import (
	"regexp"
)

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
	if payload.Contact.ProjectCoordinator.Address.Address == "" {
		return &ProjectCoordinatorAddressRequiredError{}
	}
	if payload.Contact.ProjectCoordinator.Address.ProvinceId <= 0 {
		return &ProjectCoordinatorProvinceIdRequiredError{}
	}
	if payload.Contact.ProjectCoordinator.Address.DistrictId <= 0 {
		return &ProjectCoordinatorDistrictIdRequiredError{}
	}
	if payload.Contact.ProjectCoordinator.Address.SubdistrictId <= 0 {
		return &ProjectCoordinatorSubdistrictIdRequiredError{}
	}
	if payload.Contact.ProjectCoordinator.Address.PostcodeId <= 0 {
		return &ProjectCoordinatorPostcodeIdRequiredError{}
	}
	if payload.Contact.ProjectCoordinator.Email == "" {
		return &ProjectCoordinatorEmailRequiredError{}
	}
	if payload.Contact.ProjectCoordinator.LineId == "" {
		return &ProjectCoordinatorLineIdRequiredError{}
	}
	if payload.Contact.ProjectCoordinator.PhoneNumber == "" {
		return &ProjectCoordinatorPhoneNumberRequiredError{}
	}
	if len(payload.Contact.ProjectCoordinator.PhoneNumber) < 9 {
		return &ProjectCoordinatorPhoneNumberLengthError{}
	}
	if !IsValidPhoneNumber(payload.Contact.ProjectCoordinator.PhoneNumber) {
		return &ProjectCoordinatorPhoneNumberInvalidError{}
	}
	return nil
}

func IsValidPhoneNumber(phoneNumber string) bool {
	regex := `^[0-9]{9,}$`
	re := regexp.MustCompile(regex)
	return re.Find([]byte(phoneNumber)) != nil
}
