package fintech

type DebtFreedomPlan struct {
	Allocation     Allocation
	BufferForecast BufferForecast
	DebtForecast   DebtForecast
}

type EmergencyFundPlan struct {
	Allocation Allocation
	Forecast   SavingsTierComparison
}

type FinancialAssessment struct {
	CurrentFinances CurrentFinances
	BaselineBuffer  BaselineBuffer
	ActionMessage   string
	Allocations     AllocationOptions
}

func AssessCurrentFinancialPosition(cf CurrentFinances) (FinancialAssessment, error) {
	// 1. Calculate baseline buffer
	baselineBuffer, err := CalculateBaselineBuffer(cf)
	if err != nil {
		return FinancialAssessment{}, err
	}

	// 2. Apply immediate debt action
	actionMessage, err := cf.ApplyImmediateDebtPayoff(baselineBuffer)
	if err != nil {
		return FinancialAssessment{}, err
	}

	// 3. Generate allocation options using the updated finances
	allocations, err := GenerateAllocation(cf)
	if err != nil {
		return FinancialAssessment{}, err
	}

	// 4. Return the updated financial position and results
	return FinancialAssessment{
		CurrentFinances: cf,
		BaselineBuffer:  baselineBuffer,
		ActionMessage:   actionMessage,
		Allocations:     allocations,
	}, nil
}