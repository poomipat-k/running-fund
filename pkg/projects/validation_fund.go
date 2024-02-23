package projects

func validateFund(payload AddProjectRequest) error {
	fund := payload.Fund
	if fund.Budget.Total == 0 {
		return &TotalBudgetRequiredError{}
	}
	if fund.Budget.SupportOrganization == "" {
		return &TotalBudgetRequiredError{}
	}
	return nil
}
