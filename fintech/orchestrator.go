package fintech

import "errors"

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

type DebtFreedomStrategies struct {
	Sustainable DebtFreedomPlan
	Moderate    DebtFreedomPlan
	Aggressive  DebtFreedomPlan
	Custom      DebtFreedomPlan
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

func SelectFinancialStrategy[T any](
	strategy string,
	sustainable T,
	moderate T,
	aggressive T,
) (T, error) {

	switch strategy {
	case "Sustainable":
		return sustainable, nil

	case "Moderate":
		return moderate, nil

	case "Aggressive":
		return aggressive, nil

	default:
		var zero T
		return zero, errors.New("invalid financial strategy")
	}
}

type EmergencyFundStrategies struct {
	Sustainable EmergencyFundPlan
	Moderate    EmergencyFundPlan
	Aggressive  EmergencyFundPlan
	Custom      EmergencyFundPlan
}

type EmergencyFundPlan struct {
	Allocation Allocation
	Forecast   SavingsTierComparison
}

func GenerateEmergencyFundStrategies(
	cf CurrentFinances,
	allocations AllocationOptions,
	selectedPlan DebtFreedomPlan,
	targetAmount int64,
	baselineBuffer int64,
) (EmergencyFundStrategies, error) {

	// Sustainable emergency-fund strategy
	sustainableForecast, err := SimulateEmergencyFundTiers(
		cf,
		allocations.Sustainable.Allocations.Save,
		selectedPlan.DebtForecast.TotalMonths,
		selectedPlan.DebtForecast.Phase2Months,
		selectedPlan.DebtForecast.Phase2Surplus,
		targetAmount,
		baselineBuffer,
	)
	if err != nil {
		return EmergencyFundStrategies{}, err
	}

	// Moderate emergency-fund strategy
	moderateForecast, err := SimulateEmergencyFundTiers(
		cf,
		allocations.Moderate.Allocations.Save,
		selectedPlan.DebtForecast.TotalMonths,
		selectedPlan.DebtForecast.Phase2Months,
		selectedPlan.DebtForecast.Phase2Surplus,
		targetAmount,
		baselineBuffer,
	)
	if err != nil {
		return EmergencyFundStrategies{}, err
	}

	// Aggressive emergency-fund strategy
	aggressiveForecast, err := SimulateEmergencyFundTiers(
		cf,
		allocations.Aggressive.Allocations.Save,
		selectedPlan.DebtForecast.TotalMonths,
		selectedPlan.DebtForecast.Phase2Months,
		selectedPlan.DebtForecast.Phase2Surplus,
		targetAmount,
		baselineBuffer,
	)
	if err != nil {
		return EmergencyFundStrategies{}, err
	}

	return EmergencyFundStrategies{
		Sustainable: EmergencyFundPlan{
			Allocation: allocations.Sustainable.Allocations,
			Forecast:   sustainableForecast,
		},
		Moderate: EmergencyFundPlan{
			Allocation: allocations.Moderate.Allocations,
			Forecast:   moderateForecast,
		},
		Aggressive: EmergencyFundPlan{
			Allocation: allocations.Aggressive.Allocations,
			Forecast:   aggressiveForecast,
		},
	}, nil
}
