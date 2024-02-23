package projects

var requestFundAmountOptions = map[int]bool{
	20000:  true,
	30000:  true,
	50000:  true,
	100000: true,
}

func validateFund(payload AddProjectRequest) error {
	// budget
	fund := payload.Fund
	if fund.Budget.Total == 0 {
		return &TotalBudgetRequiredError{}
	}
	if fund.Budget.SupportOrganization == "" {
		return &BudgetSupportOrganizationRequiredError{}
	}
	// request
	if !fund.Request.Type.Fund &&
		!fund.Request.Type.BIB &&
		!fund.Request.Type.Pr &&
		!fund.Request.Type.Other {
		return &FundRequestTypeRequiredOneError{}
	}
	if fund.Request.Type.Fund {
		if fund.Request.Details.FundAmount == 0 {
			return &FundRequestAmountRequiredError{}
		}
		_, found := requestFundAmountOptions[fund.Request.Details.FundAmount]
		if !found {
			return &FundRequestAmountInvalidError{}
		}
	}
	if fund.Request.Type.BIB {
		if fund.Request.Details.BibAmount == 0 {
			return &BibRequestAmountRequiredError{}
		}
		if fund.Request.Details.BibAmount < 0 {
			return &BibRequestAmountInvalidError{}
		}
	}
	if fund.Request.Type.Seminar && fund.Request.Details.Seminar == "" {
		return &SeminarRequiredError{}
	}
	if fund.Request.Type.Other && fund.Request.Details.Other == "" {
		return &OtherRequestRequiredError{}
	}
	return nil
}
