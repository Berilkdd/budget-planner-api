package fintech

import (
	"errors"
	"fmt"
)

type PathwayCode string

type Action struct {
	Code        ActionCode
	Description string
}

const (
	PathwayA1 PathwayCode = "A1" //Emergency fund covered + enough savings to clear debt
	PathwayA2 PathwayCode = "A2" //Emergency fund covered + partial debt repayment
	PathwayA3 PathwayCode = "A3" //Emergency fund covered + debt + no extra savings
	PathwayB  PathwayCode = "B"  //Emergency fund covered + no debt
	PathwayC  PathwayCode = "C"  //Emergency fund below target + debt
	PathwayD  PathwayCode = "D"  //Emergency fund below target + no debt
)

type AssessmentData struct {
	Pathway                       PathwayCode
	EmploymentStatus              EmploymentStatus
	NeedsEqualIncome              bool
	CurrentSavings                int64
	MonthsOfNeedsCoveredBySavings int64
	EmergencyFundTarget           int64
	EmergencyFundMonths           int64
	Actions                       []ActionCode
}

func AssessDeficitPosition(cf *CurrentFinances) AssessmentData {

	assessment := AssessmentData{
	EmploymentStatus:              cf.EmploymentStatus,
	NeedsEqualIncome:              cf.Needs == cf.Income,
	CurrentSavings:                cf.CurrentSavings,
	MonthsOfNeedsCoveredBySavings: cf.CurrentSavings / cf.Needs,
	Actions:                       []ActionCode{},
	}

	// This function is only for cases where needs meet or exceed income.
	if cf.Needs < cf.Income {
		return assessment
	}
	
	needsEqualIncome := cf.Needs == cf.Income
	hasDebt := cf.UnsettledDebt > 0

	// 1. Income position
	if needsEqualIncome {
		fmt.Println("WARNING")
		fmt.Println(WarningDefinitions[WarningNeedsEqualIncome].Description)
	} else {
		fmt.Println("WARNING")
		fmt.Println(WarningDefinitions[WarningNeedsExceedIncome].Description)
	}

	// 2. Calculate the required emergency fund.
	emergencyFund, err := CalculateEmergencyTarget(
		cf.EmploymentStatus,
		cf.Needs,
	)

	if err != nil {
		fmt.Println("ERROR")
		fmt.Println(err)
		return assessment
	}

	assessment.EmergencyFundTarget = emergencyFund.TargetAmount
	assessment.EmergencyFundMonths = emergencyFund.MonthsCount

	// 3. Savings below emergency-fund target.
	if cf.CurrentSavings < emergencyFund.TargetAmount {

		fmt.Println("WARNING")
		fmt.Println(WarningDefinitions[WarningBelowEmergencyFund].Description)

		// No emergency fund + needs exceed income
		if !needsEqualIncome {
			fmt.Println("ACTION")
			fmt.Println(ActionDefinitions[ActionSupportAdvised].Description)

			assessment.Actions = append(
				assessment.Actions,
				ActionSupportAdvised,
			)
		}

		// No emergency fund + debt
		if hasDebt {
			fmt.Println("WARNING")
			fmt.Println(WarningDefinitions[WarningUnsettledDebt].Description)

			fmt.Println("ACTION")
			fmt.Println(ActionDefinitions[ActionDebtAdviceAdvised].Description)

			assessment.Pathway = PathwayC

			assessment.Actions = append(
				assessment.Actions,
				ActionDebtAdviceAdvised,
			)

			fmt.Println("PATHWAY")
			fmt.Println(PathwayC)

			
		} else {
			fmt.Println("WARNING")
			fmt.Println(WarningDefinitions[WarningNoUnsettledDebt].Description)

			assessment.Pathway = PathwayD

			fmt.Println("PATHWAY")
			fmt.Println(PathwayD)		
			
		}

		return assessment
	}

	// 4. Emergency fund is covered.
	fmt.Println("WARNING")
	fmt.Println(WarningDefinitions[WarningEmergencyFundCovered].Description)

	// 5. No debt.
	if !hasDebt {

		fmt.Println("WARNING")
		fmt.Println(WarningDefinitions[WarningNoUnsettledDebt].Description)

		// Only Needs > Income creates the long-term savings warning.
		if !needsEqualIncome {
			fmt.Println("WARNING")
			fmt.Println(WarningDefinitions[WarningSavingsMayBeUsed].Description)
		}

		assessment.Pathway = PathwayB

		fmt.Println("PATHWAY")
		fmt.Println(PathwayB)
			
		return assessment
	}

	// 6. Debt remains.
	fmt.Println("WARNING")
	fmt.Println(WarningDefinitions[WarningUnsettledDebt].Description)

	extraSavings := cf.CurrentSavings - emergencyFund.TargetAmount

	// 7. Enough extra savings to clear the debt completely.
	if extraSavings >= cf.UnsettledDebt {

		fmt.Println("ACTION")
		fmt.Println(ActionDefinitions[ActionFullDebtPaymentAdvised].Description)

		assessment.Pathway = PathwayA1

		assessment.Actions = append(
			assessment.Actions,
			ActionFullDebtPaymentAdvised,
		)

		if !needsEqualIncome {
			fmt.Println("WARNING")
			fmt.Println(WarningDefinitions[WarningSavingsMayBeUsed].Description)
		}

		fmt.Println("PATHWAY")
		fmt.Println(PathwayA1)		

		return assessment
	}

	// 8. Some extra savings, but not enough to clear the debt.
	if extraSavings > 0 {

		fmt.Println("ACTION")
		fmt.Println(ActionDefinitions[ActionPartialDebtPaymentAdvised].Description)

		fmt.Println("ACTION")
		fmt.Println(ActionDefinitions[ActionDebtAdviceAdvised].Description)

		assessment.Pathway = PathwayA2

		assessment.Actions = append(
			assessment.Actions,
			ActionPartialDebtPaymentAdvised,
			ActionDebtAdviceAdvised,
		)

		if !needsEqualIncome {
			fmt.Println("WARNING")
			fmt.Println(WarningDefinitions[WarningSavingsMayBeUsed].Description)
		}		

		fmt.Println("PATHWAY")
		fmt.Println(PathwayA2)

		return assessment
	}

	// 9. Emergency fund is covered exactly, but there is no extra
	// savings available for debt repayment.
	fmt.Println("ACTION")
	fmt.Println(ActionDefinitions[ActionDebtAdviceAdvised].Description)

	assessment.Pathway = PathwayA3

	assessment.Actions = append(
		assessment.Actions,		
		ActionDebtAdviceAdvised,
	)

	fmt.Println("PATHWAY")
	fmt.Println(PathwayA3)

	
	return assessment
}
	
type DebtFreedomStrategies struct {
	Sustainable DebtFreedomPlan
	Moderate    DebtFreedomPlan
	Aggressive  DebtFreedomPlan
	Custom      DebtFreedomPlan
}

type DebtFreedomPlan struct {
	Available      bool
	Allocation     Allocation
	BufferForecast BufferForecast
	DebtForecast   DebtForecast
}

func GenerateDebtFreedomStrategies(
	cf CurrentFinances,
	allocations AllocationOptions,
	baselineBuffer BaselineBuffer,
	customContribution int64,
) (DebtFreedomStrategies, error) {

	var sustainableBuffer BufferForecast
	var sustainableDebt DebtForecast

	var moderateBuffer BufferForecast
	var moderateDebt DebtForecast

	var aggressiveBuffer BufferForecast
	var aggressiveDebt DebtForecast

	var customBuffer BufferForecast
	var customDebt DebtForecast

	var err error

	// Sustainable Plan

	if allocations.Sustainable.Available {

		sustainableBuffer, err = CalculateBufferTimeline(
			cf,
			allocations.Sustainable.Allocations.Save,
			baselineBuffer.TargetAmount,
		)
		if err != nil {
			return DebtFreedomStrategies{}, err
		}

		sustainableDebt, err = CalculateDebtTimeline(
			cf,
			allocations.Sustainable.Allocations.Save,
			sustainableBuffer.Phase1Months,
			sustainableBuffer.Phase1Surplus,
		)
		if err != nil {
			return DebtFreedomStrategies{}, err
		}
	}

	// Moderate Plan

	if allocations.Moderate.Available {

		moderateBuffer, err = CalculateBufferTimeline(
			cf,
			allocations.Moderate.Allocations.Save,
			baselineBuffer.TargetAmount,
		)
		if err != nil {
			return DebtFreedomStrategies{}, err
		}

		moderateDebt, err = CalculateDebtTimeline(
			cf,
			allocations.Moderate.Allocations.Save,
			moderateBuffer.Phase1Months,
			moderateBuffer.Phase1Surplus,
		)
		if err != nil {
			return DebtFreedomStrategies{}, err
		}
	}

	// Aggressive Plan

	if allocations.Aggressive.Available {

		aggressiveBuffer, err = CalculateBufferTimeline(
			cf,
			allocations.Aggressive.Allocations.Save,
			baselineBuffer.TargetAmount,
		)
		if err != nil {
			return DebtFreedomStrategies{}, err
		}

		aggressiveDebt, err = CalculateDebtTimeline(
			cf,
			allocations.Aggressive.Allocations.Save,
			aggressiveBuffer.Phase1Months,
			aggressiveBuffer.Phase1Surplus,
		)
		if err != nil {
			return DebtFreedomStrategies{}, err
		}	
	}	

	// Custom Plan

	customPlan, err := GenerateCustomAllocation(
		cf,
		customContribution,
	)
	if err != nil {
		return DebtFreedomStrategies{}, err
	}

	customBuffer, err = CalculateBufferTimeline(
		cf,
		customPlan.Allocations.Save,
		baselineBuffer.TargetAmount,
	)
	if err != nil {
		return DebtFreedomStrategies{}, err
	}

	customDebt, err = CalculateDebtTimeline(
		cf,
		customPlan.Allocations.Save,
		customBuffer.Phase1Months,
		customBuffer.Phase1Surplus,
	)
	if err != nil {
		return DebtFreedomStrategies{}, err
	}

	return DebtFreedomStrategies{
		Sustainable: DebtFreedomPlan{
			Available:      allocations.Sustainable.Available,
			Allocation:     allocations.Sustainable.Allocations,
			BufferForecast: sustainableBuffer,
			DebtForecast:   sustainableDebt,
		},
		Moderate: DebtFreedomPlan{
			Available:      allocations.Moderate.Available,
			Allocation:     allocations.Moderate.Allocations,
			BufferForecast: moderateBuffer,
			DebtForecast:   moderateDebt,
		},
		Aggressive: DebtFreedomPlan{
			Available:      allocations.Aggressive.Available,
			Allocation:     allocations.Aggressive.Allocations,
			BufferForecast: aggressiveBuffer,
			DebtForecast:   aggressiveDebt,
		},
		Custom: DebtFreedomPlan{
			Available:      true,
			Allocation:     customPlan.Allocations,
			BufferForecast: customBuffer,
			DebtForecast:   customDebt,
		},
	}, nil
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