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

type GeneralAddressRequiredError struct{}

func (e *GeneralAddressRequiredError) Error() string {
	return "general address is required"
}

type GeneralProvinceRequiredError struct{}

func (e *GeneralProvinceRequiredError) Error() string {
	return "general provinceId is required"
}

type GeneralDistrictIdRequiredError struct{}

func (e *GeneralDistrictIdRequiredError) Error() string {
	return "general districtId is required"
}

type GeneralSubdistrictIdRequiredError struct{}

func (e *GeneralSubdistrictIdRequiredError) Error() string {
	return "general subdistrictIdId is required"
}

type GeneralPostcodeIdRequiredError struct{}

func (e *GeneralPostcodeIdRequiredError) Error() string {
	return "general postcodeId is required"
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

type VIPFeeRequiredError struct{}

func (e *VIPFeeRequiredError) Error() string {
	return "vipFee is required"
}

type VIPFeeNegativeError struct{}

func (e *VIPFeeNegativeError) Error() string {
	return "vipFee must greater equal to zero"
}

type ExpectedParticipantsRequiredError struct{}

func (e *ExpectedParticipantsRequiredError) Error() string {
	return "expectedParticipants is required"
}

type ExpectedParticipantsInvalidError struct{}

func (e *ExpectedParticipantsInvalidError) Error() string {
	return "expectedParticipants is invalid"
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

// Head address
type ProjectHeadAddressRequiredError struct{}

func (e *ProjectHeadAddressRequiredError) Error() string {
	return "ProjectHead address is required"
}

type ProjectHeadProvinceIdRequiredError struct{}

func (e *ProjectHeadProvinceIdRequiredError) Error() string {
	return "ProjectHead provinceId is required"
}

type ProjectHeadDistrictIdRequiredError struct{}

func (e *ProjectHeadDistrictIdRequiredError) Error() string {
	return "ProjectHead districtId is required"
}

type ProjectHeadSubdistrictIdRequiredError struct{}

func (e *ProjectHeadSubdistrictIdRequiredError) Error() string {
	return "ProjectHead subdistrictId is required"
}

type ProjectHeadPostcodeIdRequiredError struct{}

func (e *ProjectHeadPostcodeIdRequiredError) Error() string {
	return "ProjectHead postcodeId is required"
}

type ProjectHeadEmailRequiredError struct{}

func (e *ProjectHeadEmailRequiredError) Error() string {
	return "ProjectHead email is required"
}

type ProjectHeadLineIdRequiredError struct{}

func (e *ProjectHeadLineIdRequiredError) Error() string {
	return "ProjectHead lineId is required"
}

type ProjectHeadPhoneNumberRequiredError struct{}

func (e *ProjectHeadPhoneNumberRequiredError) Error() string {
	return "ProjectHead phoneNumber is required"
}

type ProjectHeadPhoneNumberLengthError struct{}

func (e *ProjectHeadPhoneNumberLengthError) Error() string {
	return "ProjectHead phoneNumber is shorter than 9 numbers"
}

type ProjectHeadPhoneNumberInvalidError struct{}

func (e *ProjectHeadPhoneNumberInvalidError) Error() string {
	return "ProjectHead phoneNumber is invalid"
}

// Manager address
type ProjectManagerAddressRequiredError struct{}

func (e *ProjectManagerAddressRequiredError) Error() string {
	return "ProjectManager address is required"
}

type ProjectManagerProvinceIdRequiredError struct{}

func (e *ProjectManagerProvinceIdRequiredError) Error() string {
	return "ProjectManager provinceId is required"
}

type ProjectManagerDistrictIdRequiredError struct{}

func (e *ProjectManagerDistrictIdRequiredError) Error() string {
	return "ProjectManager districtId is required"
}

type ProjectManagerSubdistrictIdRequiredError struct{}

func (e *ProjectManagerSubdistrictIdRequiredError) Error() string {
	return "ProjectManager subdistrictId is required"
}

type ProjectManagerPostcodeIdRequiredError struct{}

func (e *ProjectManagerPostcodeIdRequiredError) Error() string {
	return "ProjectManager postcodeId is required"
}

type ProjectManagerEmailRequiredError struct{}

func (e *ProjectManagerEmailRequiredError) Error() string {
	return "ProjectManager email is required"
}

type ProjectManagerLineIdRequiredError struct{}

func (e *ProjectManagerLineIdRequiredError) Error() string {
	return "ProjectManager lineId is required"
}

type ProjectManagerPhoneNumberRequiredError struct{}

func (e *ProjectManagerPhoneNumberRequiredError) Error() string {
	return "ProjectManager phoneNumber is required"
}

type ProjectManagerPhoneNumberLengthError struct{}

func (e *ProjectManagerPhoneNumberLengthError) Error() string {
	return "ProjectManager phoneNumber is shorter than 9 numbers"
}

type ProjectManagerPhoneNumberInvalidError struct{}

func (e *ProjectManagerPhoneNumberInvalidError) Error() string {
	return "ProjectManager phoneNumber is invalid"
}

// Coordinator address
type ProjectCoordinatorAddressRequiredError struct{}

func (e *ProjectCoordinatorAddressRequiredError) Error() string {
	return "projectCoordinator address is required"
}

type ProjectCoordinatorProvinceIdRequiredError struct{}

func (e *ProjectCoordinatorProvinceIdRequiredError) Error() string {
	return "projectCoordinator provinceId is required"
}

type ProjectCoordinatorDistrictIdRequiredError struct{}

func (e *ProjectCoordinatorDistrictIdRequiredError) Error() string {
	return "projectCoordinator districtId is required"
}

type ProjectCoordinatorSubdistrictIdRequiredError struct{}

func (e *ProjectCoordinatorSubdistrictIdRequiredError) Error() string {
	return "projectCoordinator subdistrictId is required"
}

type ProjectCoordinatorPostcodeIdRequiredError struct{}

func (e *ProjectCoordinatorPostcodeIdRequiredError) Error() string {
	return "projectCoordinator postcodeId is required"
}

type ProjectCoordinatorEmailRequiredError struct{}

func (e *ProjectCoordinatorEmailRequiredError) Error() string {
	return "projectCoordinator email is required"
}

type ProjectCoordinatorLineIdRequiredError struct{}

func (e *ProjectCoordinatorLineIdRequiredError) Error() string {
	return "projectCoordinator lineId is required"
}

type ProjectCoordinatorPhoneNumberRequiredError struct{}

func (e *ProjectCoordinatorPhoneNumberRequiredError) Error() string {
	return "projectCoordinator phoneNumber is required"
}

type ProjectCoordinatorPhoneNumberLengthError struct{}

func (e *ProjectCoordinatorPhoneNumberLengthError) Error() string {
	return "projectCoordinator phoneNumber is shorter than 9 numbers"
}

type ProjectCoordinatorPhoneNumberInvalidError struct{}

func (e *ProjectCoordinatorPhoneNumberInvalidError) Error() string {
	return "projectCoordinator phoneNumber is invalid"
}

type RaceDirectorWhoRequiredError struct{}

func (e *RaceDirectorWhoRequiredError) Error() string {
	return "raceDirector who is required"
}

type RaceDirectorAlternativePrefixRequiredError struct{}

func (e *RaceDirectorAlternativePrefixRequiredError) Error() string {
	return "raceDirector alternative prefix is required"
}

type RaceDirectorAlternativeFirstNameRequiredError struct{}

func (e *RaceDirectorAlternativeFirstNameRequiredError) Error() string {
	return "raceDirector alternative firstName is required"
}

type RaceDirectorAlternativeLastNameRequiredError struct{}

func (e *RaceDirectorAlternativeLastNameRequiredError) Error() string {
	return "raceDirector alternative lastName is required"
}

type ContactOrganizationNameRequiredError struct{}

func (e *ContactOrganizationNameRequiredError) Error() string {
	return "organization name is required"
}

type ContactOrganizationTypeRequiredError struct{}

func (e *ContactOrganizationTypeRequiredError) Error() string {
	return "organization type is required"
}

type BackgroundRequiredError struct{}

func (e *BackgroundRequiredError) Error() string {
	return "background is required"
}

type ObjectiveRequiredError struct{}

func (e *ObjectiveRequiredError) Error() string {
	return "objective is required"
}

type OnlineAvailableRequiredOne struct{}

func (e *OnlineAvailableRequiredOne) Error() string {
	return "online marketing available required one"
}

type FacebookHowToIsRequired struct{}

func (e *FacebookHowToIsRequired) Error() string {
	return "facebook link is required"
}

type WebsiteHowToIsRequired struct{}

func (e *WebsiteHowToIsRequired) Error() string {
	return "website link is required"
}

type OnlinePageHowToIsRequired struct{}

func (e *OnlinePageHowToIsRequired) Error() string {
	return "onlinePage link is required"
}

type OtherHowToIsRequired struct{}

func (e *OtherHowToIsRequired) Error() string {
	return "other link is required"
}

type OfflineAvailableRequiredOne struct{}

func (e *OfflineAvailableRequiredOne) Error() string {
	return "offline marketing available required one"
}

type OfflineAdditionRequiredError struct{}

func (e *OfflineAdditionRequiredError) Error() string {
	return "offline addition is required"
}

type ApplicantCriteriaNotFoundError struct{}

func (e *ApplicantCriteriaNotFoundError) Error() string {
	return "applicant criteria not found"
}

type ScoreRequiredError struct {
	Name string
}

func (e *ScoreRequiredError) Error() string {
	return fmt.Sprintf("score %s is required", e.Name)
}

type ScoreInvalidError struct {
	Name string
}

func (e *ScoreInvalidError) Error() string {
	return fmt.Sprintf("score %s is invalid. 1 <= score <= 5", e.Name)
}

type SafetyReadyRequiredOneError struct{}

func (e *SafetyReadyRequiredOneError) Error() string {
	return "safety ready required one"
}

type AEDCountInvalidError struct{}

func (e *AEDCountInvalidError) Error() string {
	return "safety aedCount is invalid. aedCount must >= 1"
}

type SafetyAdditionRequiredError struct{}

func (e *SafetyAdditionRequiredError) Error() string {
	return "safety addition is required"
}

type RouteMeasurementRequiredOneError struct{}

func (e *RouteMeasurementRequiredOneError) Error() string {
	return "route measurement required one"
}

type RouteToolRequiredError struct{}

func (e *RouteToolRequiredError) Error() string {
	return "route tool is required"
}

type RouteTrafficManagementRequiredOneError struct{}

func (e *RouteTrafficManagementRequiredOneError) Error() string {
	return "route trafficManagement required one"
}

type JudgeTypeRequiredError struct{}

func (e *JudgeTypeRequiredError) Error() string {
	return "judge type is required"
}

type JudgeTypeInvalidError struct{}

func (e *JudgeTypeInvalidError) Error() string {
	return "judge type is invalid"
}

type JudgeOtherTypeRequiredError struct{}

func (e *JudgeOtherTypeRequiredError) Error() string {
	return "judge otherType is required"
}

type SupportOrganizationRequiredOneError struct{}

func (e *SupportOrganizationRequiredOneError) Error() string {
	return "support organization required one"
}

type SupportAdditionRequiredError struct{}

func (e *SupportAdditionRequiredError) Error() string {
	return "support addition is required"
}

type FeedbackRequiredError struct{}

func (e *FeedbackRequiredError) Error() string {
	return "feedback is required"
}

type ThisSeriesFirstTimeRequiredError struct{}

func (e *ThisSeriesFirstTimeRequiredError) Error() string {
	return "firstTime is required"
}

type HistoryOrdinalNumberInvalidError struct{}

func (e *HistoryOrdinalNumberInvalidError) Error() string {
	return "history ordinalNumber is invalid"
}

type HistoryYearRequiredError struct{}

func (e *HistoryYearRequiredError) Error() string {
	return "history year is required"
}

type HistoryYearOutOfBoundError struct{}

func (e *HistoryYearOutOfBoundError) Error() string {
	return "history year is out of bound"
}

type HistoryMonthRequiredError struct{}

func (e *HistoryMonthRequiredError) Error() string {
	return "history month is required"
}

type HistoryMonthOutOfBoundError struct{}

func (e *HistoryMonthOutOfBoundError) Error() string {
	return "history month is out of bound"
}

type HistoryDayRequiredError struct{}

func (e *HistoryDayRequiredError) Error() string {
	return "history day is required"
}

type HistoryDayOutOfBoundError struct{}

func (e *HistoryDayOutOfBoundError) Error() string {
	return "history day is out of bound"
}

type CompletedYearRequiredError struct{}

func (e *CompletedYearRequiredError) Error() string {
	return "history completed year is required"
}

type CompletedYearOutOfBoundError struct{}

func (e *CompletedYearOutOfBoundError) Error() string {
	return "history completed year is out of bound"
}

type CompletedNameRequiredError struct{}

func (e *CompletedNameRequiredError) Error() string {
	return "history completed name is required"
}

type CompletedParticipantRequiredError struct{}

func (e *CompletedParticipantRequiredError) Error() string {
	return "history completed participant is required"
}

type CompletedParticipantInvalidError struct{}

func (e *CompletedParticipantInvalidError) Error() string {
	return "history completed participant is invalid"
}

type DoneBeforeRequiredError struct{}

func (e *DoneBeforeRequiredError) Error() string {
	return "otherSeries doneBefore is required"
}

type TotalBudgetRequiredError struct{}

func (e *TotalBudgetRequiredError) Error() string {
	return "budget total is required"
}

type BudgetSupportOrganizationRequiredError struct{}

func (e *BudgetSupportOrganizationRequiredError) Error() string {
	return "budget supportOrganization is required"
}

type NoAlcoholSponsorError struct{}

func (e *NoAlcoholSponsorError) Error() string {
	return "budget noAlcohol sponsor must be checked"
}

type FundRequestTypeRequiredOneError struct{}

func (e *FundRequestTypeRequiredOneError) Error() string {
	return "fund request type is required"
}

type FundRequestAmountRequiredError struct{}

func (e *FundRequestAmountRequiredError) Error() string {
	return "request fundAmount is required"
}

type FundRequestAmountInvalidError struct{}

func (e *FundRequestAmountInvalidError) Error() string {
	return "request fundAmount is invalid"
}

type BibRequestAmountRequiredError struct{}

func (e *BibRequestAmountRequiredError) Error() string {
	return "request bibAmount is required"
}

type BibRequestAmountInvalidError struct{}

func (e *BibRequestAmountInvalidError) Error() string {
	return "request bibAmount is invalid"
}

type SeminarRequiredError struct{}

func (e *SeminarRequiredError) Error() string {
	return "request details seminar is required"
}

type OtherRequestRequiredError struct{}

func (e *OtherRequestRequiredError) Error() string {
	return "request details otherRequest is required"
}

type MarketingFilesRequiredError struct{}

func (e *MarketingFilesRequiredError) Error() string {
	return "marketingFiles are required"
}

type RouteFilesRequiredError struct{}

func (e *RouteFilesRequiredError) Error() string {
	return "routeFiles are required"
}

type EventMapFilesRequiredError struct{}

func (e *EventMapFilesRequiredError) Error() string {
	return "EventMapFiles are required"
}

type EventDetailsFilesRequiredError struct{}

func (e *EventDetailsFilesRequiredError) Error() string {
	return "eventDetailsFiles are required"
}

type FilesRequiredError struct{}

func (e *FilesRequiredError) Error() string {
	return "additionFiles or etcFiles are required"
}

type ProjectNotFoundError struct{}

func (e *ProjectNotFoundError) Error() string {
	return "project is not found"
}

type ProjectCodeRequiredError struct{}

func (e *ProjectCodeRequiredError) Error() string {
	return "projectCode is required"
}

type ReviewerIdRequiredError struct{}

func (e *ReviewerIdRequiredError) Error() string {
	return "reviewerId is required"
}

type ProjectStatusPrimaryRequiredError struct{}

func (e *ProjectStatusPrimaryRequiredError) Error() string {
	return "projectStatusPrimary is required"
}

type ProjectStatusPrimaryInvalidError struct{}

func (e *ProjectStatusPrimaryInvalidError) Error() string {
	return "projectStatusPrimary is invalid"
}

type ProjectStatusSecondaryRequiredError struct{}

func (e *ProjectStatusSecondaryRequiredError) Error() string {
	return "projectStatusSecondary is required"
}

type ProjectStatusSecondaryInvalidError struct{}

func (e *ProjectStatusSecondaryInvalidError) Error() string {
	return "projectStatusSecondary is invalid"
}

type AdminScoreOutOfRangeError struct{}

func (e *AdminScoreOutOfRangeError) Error() string {
	return "adminScore is out of range 0-100"
}

type FundApprovedAmountNegativeError struct{}

func (e *FundApprovedAmountNegativeError) Error() string {
	return "fundApprovedAmount is negative"
}

type AdminCommentTooLongError struct {
	Length int
}

func (e *AdminCommentTooLongError) Error() string {
	return fmt.Sprintf("adminComment length is over 512 characters, got %d", e.Length)
}

type FromYearRequiredError struct{}

func (e *FromYearRequiredError) Error() string {
	return "fromYear is required"
}

type ToYearRequiredError struct{}

func (e *ToYearRequiredError) Error() string {
	return "toYear is required"
}

type PageNoInvalidError struct{}

func (e *PageNoInvalidError) Error() string {
	return "pageNo is invalid"
}

type PageSizeInvalidError struct{}

func (e *PageSizeInvalidError) Error() string {
	return "pageSize is invalid"
}

type SortByRequiredError struct{}

func (e *SortByRequiredError) Error() string {
	return "sortBy is required"
}

type SortByInvalidError struct{}

func (e *SortByInvalidError) Error() string {
	return "sortBy is invalid"
}

type FromDateExceedToDateError struct{}

func (e *FromDateExceedToDateError) Error() string {
	return "fromDate is later than toDate"
}
