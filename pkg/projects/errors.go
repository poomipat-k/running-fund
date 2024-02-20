package projects

import "fmt"

type InvalidError struct {
	Name string
}

func (e *InvalidError) Error() string {
	return fmt.Sprintf("%s is invalid", e.Name)
}

// Collaborated
type CollaboratedRequiredError struct{}

func (e *CollaboratedRequiredError) Error() string {
	return "collaborated is required"
}

type CollaboratedFilesRequiredError struct{}

func (e *CollaboratedFilesRequiredError) Error() string {
	return "collaboratedFiles are required"
}

// General
type ProjectNameRequiredError struct{}

func (e *ProjectNameRequiredError) Error() string {
	return "projectName is required"
}

type YearRequiredError struct{}

func (e *YearRequiredError) Error() string {
	return "year is required"
}

type YearInvalidError struct{}

func (e *YearInvalidError) Error() string {
	return "year must greater than 1971"
}

type MonthRequiredError struct{}

func (e *MonthRequiredError) Error() string {
	return "month is required"
}

type MonthOutOfBoundError struct{}

func (e *MonthOutOfBoundError) Error() string {
	return "month must greater than > 0 and <= 12"
}

type DayRequiredError struct{}

func (e *DayRequiredError) Error() string {
	return "day is required"
}

type DayOutOfBoundError struct{}

func (e *DayOutOfBoundError) Error() string {
	return "day is not valid"
}

type FromHourRequiredError struct{}

func (e *FromHourRequiredError) Error() string {
	return "fromHour is required"
}

type FromMinuteRequiredError struct{}

func (e *FromMinuteRequiredError) Error() string {
	return "fromMinute is required"
}

type ToHourRequiredError struct{}

func (e *ToHourRequiredError) Error() string {
	return "toHour is required"
}

type ToMinuteRequiredError struct{}

func (e *ToMinuteRequiredError) Error() string {
	return "toMinute is required"
}

type AddressRequiredError struct{}

func (e *AddressRequiredError) Error() string {
	return "address is required"
}

type ProvinceRequiredError struct{}

func (e *ProvinceRequiredError) Error() string {
	return "provinceId is required"
}

type DistrictIdRequiredError struct{}

func (e *DistrictIdRequiredError) Error() string {
	return "districtId is required"
}

type SubdistrictIdRequiredError struct{}

func (e *SubdistrictIdRequiredError) Error() string {
	return "subdistrictIdId is required"
}

type PostcodeIdRequiredError struct{}

func (e *PostcodeIdRequiredError) Error() string {
	return "postcodeId is required"
}

type StartPointRequiredError struct{}

func (e *StartPointRequiredError) Error() string {
	return "startPoint is required"
}

type FinishPointRequiredError struct{}

func (e *FinishPointRequiredError) Error() string {
	return "finishPoint is required"
}

type CategoryAvailableRequiredOneError struct{}

func (e *CategoryAvailableRequiredOneError) Error() string {
	return "must select at least one category"
}

type OtherEventTypeRequiredError struct{}

func (e *OtherEventTypeRequiredError) Error() string {
	return "otherType is required"
}

type DistanceRequiredOneError struct{}

func (e *DistanceRequiredOneError) Error() string {
	return "must select at least one distance"
}

type DistanceTypeRequiredError struct{}

func (e *DistanceTypeRequiredError) Error() string {
	return "type is required"
}

type DistanceFeeRequiredError struct{}

func (e *DistanceFeeRequiredError) Error() string {
	return "fee is required"
}

type ValueNegativeError struct{}

func (e *ValueNegativeError) Error() string {
	return "value must >= 0"
}

type DistanceAndFeeDynamicRequiredError struct{}

func (e *DistanceAndFeeDynamicRequiredError) Error() string {
	return "dynamic is required"
}

type VIPRequiredError struct{}

func (e *VIPRequiredError) Error() string {
	return "vip is required"
}

type ExpectedParticipantsRequiredError struct{}

func (e *ExpectedParticipantsRequiredError) Error() string {
	return "expectedParticipants is required"
}

type HasOrganizerRequiredError struct{}

func (e *HasOrganizerRequiredError) Error() string {
	return "hasOrganizer is required"
}

type OrganizerNameRequiredError struct{}

func (e *OrganizerNameRequiredError) Error() string {
	return "organizerName is required"
}

// Contact

// ProjectHead
type ProjectHeadPrefixRequiredError struct{}

func (e *ProjectHeadPrefixRequiredError) Error() string {
	return "projectHead prefix is required"
}

type ProjectHeadFirstNameRequiredError struct{}

func (e *ProjectHeadFirstNameRequiredError) Error() string {
	return "projectHead firstName is required"
}

type ProjectHeadLastNameRequiredError struct{}

func (e *ProjectHeadLastNameRequiredError) Error() string {
	return "projectHead lastName is required"
}

type ProjectHeadOrganizationPositionRequiredError struct{}

func (e *ProjectHeadOrganizationPositionRequiredError) Error() string {
	return "projectHead organizationPosition is required"
}

type ProjectHeadEventPositionRequiredError struct{}

func (e *ProjectHeadEventPositionRequiredError) Error() string {
	return "projectHead eventPosition is required"
}

// ProjectManager
type ProjectManagerPrefixRequiredError struct{}

func (e *ProjectManagerPrefixRequiredError) Error() string {
	return "projectManager prefix is required"
}

type ProjectManagerFirstNameRequiredError struct{}

func (e *ProjectManagerFirstNameRequiredError) Error() string {
	return "projectManager firstName is required"
}

type ProjectManagerLastNameRequiredError struct{}

func (e *ProjectManagerLastNameRequiredError) Error() string {
	return "projectManager lastName is required"
}

type ProjectManagerOrganizationPositionRequiredError struct{}

func (e *ProjectManagerOrganizationPositionRequiredError) Error() string {
	return "projectManager organizationPosition is required"
}

type ProjectManagerEventPositionRequiredError struct{}

func (e *ProjectManagerEventPositionRequiredError) Error() string {
	return "projectManager eventPosition is required"
}

// ProjectCoordinator
type ProjectCoordinatorPrefixRequiredError struct{}

func (e *ProjectCoordinatorPrefixRequiredError) Error() string {
	return "projectCoordinator prefix is required"
}

type ProjectCoordinatorFirstNameRequiredError struct{}

func (e *ProjectCoordinatorFirstNameRequiredError) Error() string {
	return "projectCoordinator firstName is required"
}

type ProjectCoordinatorLastNameRequiredError struct{}

func (e *ProjectCoordinatorLastNameRequiredError) Error() string {
	return "projectCoordinator lastName is required"
}

type ProjectCoordinatorOrganizationPositionRequiredError struct{}

func (e *ProjectCoordinatorOrganizationPositionRequiredError) Error() string {
	return "projectCoordinator organizationPosition is required"
}

type ProjectCoordinatorEventPositionRequiredError struct{}

func (e *ProjectCoordinatorEventPositionRequiredError) Error() string {
	return "projectCoordinator eventPosition is required"
}
