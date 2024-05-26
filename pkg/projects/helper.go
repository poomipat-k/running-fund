package projects

func hasPrimaryStatusChanged(currentStatus, primaryStatus string) bool {
	currentVal := PROJECT_STATUS[currentStatus]
	if currentVal <= 3 {
		return primaryStatus != "CurrentBeforeApprove"
	}
	if currentVal == 4 {
		return primaryStatus != "NotApproved"
	}
	if currentVal > 4 {
		return primaryStatus != "Approved"
	}
	return false
}
