package fintech

type FinancialAssessment struct {
	CurrentFinances CurrentFinances
	BaselineBuffer  BaselineBuffer
	ActionMessage   string
	Allocations     AllocationOptions
}

type EmergencyFundPlan struct {
	Allocation Allocation
	Forecast   SavingsTierComparison
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

type DebtFreedomStrategies struct {
	Sustainable DebtFreedomPlan
	Moderate    DebtFreedomPlan
	Aggressive  DebtFreedomPlan
}

type DebtFreedomPlan struct {
	Allocation     Allocation
	BufferForecast BufferForecast
	DebtForecast   DebtForecast
}


func GenerateDebtFreedomStrategies(
	cf CurrentFinances,
	allocations AllocationOptions,
	baselineBuffer BaselineBuffer,
	) (DebtFreedomStrategies, error) {

	// Sustainable Plan
	sustainableBuffer, err := CalculateBufferTimeline(
		cf,
		allocations.Sustainable.Allocations.Save,
		baselineBuffer.TargetAmount,
	)
	if err != nil {
		return DebtFreedomStrategies{}, err
	}

	sustainableDebt, err := CalculateDebtTimeline(
		cf,
		allocations.Sustainable.Allocations.Save,
		sustainableBuffer.Phase1Months,
		sustainableBuffer.Phase1Surplus,
	)
	if err != nil {
		return DebtFreedomStrategies{}, err
	}

	// Moderate Plan
	moderateBuffer, err := CalculateBufferTimeline(
		cf,
		allocations.Moderate.Allocations.Save,
		baselineBuffer.TargetAmount,
	)
	if err != nil {
		return DebtFreedomStrategies{}, err
	}

	moderateDebt, err := CalculateDebtTimeline(
		cf,
		allocations.Moderate.Allocations.Save,
		moderateBuffer.Phase1Months,
		moderateBuffer.Phase1Surplus,
	)
	if err != nil {
		return DebtFreedomStrategies{}, err
	}

	// Aggressive Plan
	aggressiveBuffer, err := CalculateBufferTimeline(
		cf,
		allocations.Aggressive.Allocations.Save,
		baselineBuffer.TargetAmount,
	)
	if err != nil {
		return DebtFreedomStrategies{}, err
	}

	aggressiveDebt, err := CalculateDebtTimeline(
		cf,
		allocations.Aggressive.Allocations.Save,
		aggressiveBuffer.Phase1Months,
		aggressiveBuffer.Phase1Surplus,
	)
	if err != nil {
		return DebtFreedomStrategies{}, err
	}

	return DebtFreedomStrategies{
		Sustainable: DebtFreedomPlan{
			Allocation:     allocations.Sustainable.Allocations,
			BufferForecast: sustainableBuffer,
			DebtForecast:   sustainableDebt,
		},
		Moderate: DebtFreedomPlan{
			Allocation:     allocations.Moderate.Allocations,
			BufferForecast: moderateBuffer,
			DebtForecast:   moderateDebt,
		},
		Aggressive: DebtFreedomPlan{
			Allocation:     allocations.Aggressive.Allocations,
			BufferForecast: aggressiveBuffer,
			DebtForecast:   aggressiveDebt,
		},
	}, nil
}