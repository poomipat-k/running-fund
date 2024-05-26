package projects

import (
	"regexp"
)

func validateContact(payload AddProjectRequest) error {
	// projectHead
	if err := validateProjectHead(payload); err != nil {
		return err
	}
	// projectManager
	if err := validateProjectManager(payload); err != nil {
		return err
	}
	// projectCoordinator
	if err := validateProjectCoordinator(payload); err != nil {
		return err
	}
	// raceDirector
	if payload.Contact.RaceDirector.Who == "" {
		return &RaceDirectorWhoRequiredError{}
	}
	if payload.Contact.RaceDirector.Who == "other" {
		if payload.Contact.RaceDirector.Alternative.Prefix == "" {
			return &RaceDirectorAlternativePrefixRequiredError{}
		}
		if payload.Contact.RaceDirector.Alternative.FirstName == "" {
			return &RaceDirectorAlternativeFirstNameRequiredError{}
		}
		if payload.Contact.RaceDirector.Alternative.LastName == "" {
			return &RaceDirectorAlternativeLastNameRequiredError{}
		}
	}
	// organization
	if payload.Contact.Organization.Name == "" {
		return &ContactOrganizationNameRequiredError{}
	}
	if payload.Contact.Organization.Type == "" {
		return &ContactOrganizationTypeRequiredError{}
	}
	return nil
}

func validateProjectHead(payload AddProjectRequest) error {
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
	if payload.Contact.ProjectHead.Address.Address == "" {
		return &ProjectHeadAddressRequiredError{}
	}
	if payload.Contact.ProjectHead.Address.ProvinceId <= 0 {
		return &ProjectHeadProvinceIdRequiredError{}
	}
	if payload.Contact.ProjectHead.Address.DistrictId <= 0 {
		return &ProjectHeadDistrictIdRequiredError{}
	}
	if payload.Contact.ProjectHead.Address.SubdistrictId <= 0 {
		return &ProjectHeadSubdistrictIdRequiredError{}
	}
	if payload.Contact.ProjectHead.Address.PostcodeId <= 0 {
		return &ProjectHeadPostcodeIdRequiredError{}
	}
	if payload.Contact.ProjectHead.Email == "" {
		return &ProjectHeadEmailRequiredError{}
	}
	if payload.Contact.ProjectHead.LineId == "" {
		return &ProjectHeadLineIdRequiredError{}
	}
	if payload.Contact.ProjectHead.PhoneNumber == "" {
		return &ProjectHeadPhoneNumberRequiredError{}
	}
	if len(payload.Contact.ProjectHead.PhoneNumber) < 9 {
		return &ProjectHeadPhoneNumberLengthError{}
	}
	if !IsValidPhoneNumber(payload.Contact.ProjectHead.PhoneNumber) {
		return &ProjectHeadPhoneNumberInvalidError{}
	}
	return nil
}
func validateProjectManager(payload AddProjectRequest) error {
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
	if payload.Contact.ProjectManager.Address.Address == "" {
		return &ProjectManagerAddressRequiredError{}
	}
	if payload.Contact.ProjectManager.Address.ProvinceId <= 0 {
		return &ProjectManagerProvinceIdRequiredError{}
	}
	if payload.Contact.ProjectManager.Address.DistrictId <= 0 {
		return &ProjectManagerDistrictIdRequiredError{}
	}
	if payload.Contact.ProjectManager.Address.SubdistrictId <= 0 {
		return &ProjectManagerSubdistrictIdRequiredError{}
	}
	if payload.Contact.ProjectManager.Address.PostcodeId <= 0 {
		return &ProjectManagerPostcodeIdRequiredError{}
	}
	if payload.Contact.ProjectManager.Email == "" {
		return &ProjectManagerEmailRequiredError{}
	}
	if payload.Contact.ProjectManager.LineId == "" {
		return &ProjectManagerLineIdRequiredError{}
	}
	if payload.Contact.ProjectManager.PhoneNumber == "" {
		return &ProjectManagerPhoneNumberRequiredError{}
	}
	if len(payload.Contact.ProjectManager.PhoneNumber) < 9 {
		return &ProjectManagerPhoneNumberLengthError{}
	}
	if !IsValidPhoneNumber(payload.Contact.ProjectManager.PhoneNumber) {
		return &ProjectManagerPhoneNumberInvalidError{}
	}
	return nil
}
func validateProjectCoordinator(payload AddProjectRequest) error {
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
