package fintech

import "fmt"

func PrintCurrentFinances(cf CurrentFinances) {
	fmt.Println("=== Current Financial Position ===")
	fmt.Printf("Income: £%.2f\n", float64(cf.Income)/100)
	fmt.Printf("Needs: £%.2f\n", float64(cf.Needs)/100)
	fmt.Printf("Current savings: £%.2f\n", float64(cf.CurrentSavings)/100)
	fmt.Printf("Available surplus: £%.2f\n", float64(cf.AvailableSurplus)/100)
	fmt.Printf("Unsettled debt: £%.2f\n", float64(cf.UnsettledDebt)/100)
	fmt.Printf("Has debt: %t\n", cf.HasDebt)
	fmt.Printf("Employment status: %s\n", cf.EmploymentStatus)
	fmt.Printf("Current date: %s\n", cf.CurrentDate.Format("January 2006"))
	fmt.Println("=================================")
}

// PrintDMPAssessment prints the result of the DMP safety assessment.
func PrintDMPAssessment(assessment DMPAssessment) {
	if !assessment.DMPRequired {
		return
	}

	for _, reason := range assessment.Reasons {
		fmt.Printf("[WARNING] %s\n", reason)
	}

	fmt.Println("[ACTION] DMP support is advised.")
}

// PrintAssessment prints the pathway, warnings, and recommended actions
// produced by the assessment functions.
func PrintAssessment(assessment AssessmentData) {
	switch assessment.Pathway {
	case PathwayA1:
		fmt.Printf(
			"[WARNING] %s\n",
			WarningDefinitions[WarningUnsettledDebt].Description,
		)
		fmt.Printf(
			"[ACTION] %s\n",
			ActionDefinitions[ActionFullDebtPaymentAdvised].Description,
		)

	case PathwayA2:
		fmt.Printf(
			"[WARNING] %s\n",
			WarningDefinitions[WarningUnsettledDebt].Description,
		)
		fmt.Printf(
			"[ACTION] %s\n",
			ActionDefinitions[ActionPartialDebtPaymentAdvised].Description,
		)
		fmt.Printf(
			"[ACTION] %s\n",
			ActionDefinitions[ActionDebtAdviceAdvised].Description,
		)

	case PathwayA3:
		fmt.Printf(
			"[WARNING] %s\n",
			WarningDefinitions[WarningUnsettledDebt].Description,
		)
		fmt.Printf(
			"[ACTION] %s\n",
			ActionDefinitions[ActionDebtAdviceAdvised].Description,
		)

	case PathwayB:
		fmt.Printf(
			"[WARNING] %s\n",
			WarningDefinitions[WarningEmergencyFundCovered].Description,
		)
		fmt.Printf(
			"[WARNING] %s\n",
			WarningDefinitions[WarningNoUnsettledDebt].Description,
		)

	case PathwayC:
		fmt.Printf(
			"[WARNING] %s\n",
			WarningDefinitions[WarningBelowEmergencyFund].Description,
		)
		fmt.Printf(
			"[WARNING] %s\n",
			WarningDefinitions[WarningUnsettledDebt].Description,
		)
		fmt.Printf(
			"[ACTION] %s\n",
			ActionDefinitions[ActionDebtAdviceAdvised].Description,
		)

	case PathwayD:
		fmt.Printf(
			"[WARNING] %s\n",
			WarningDefinitions[WarningBelowEmergencyFund].Description,
		)
		fmt.Printf(
			"[WARNING] %s\n",
			WarningDefinitions[WarningNoUnsettledDebt].Description,
		)

	case PathwayE1:
		fmt.Printf(
			"[WARNING] %s\n",
			WarningDefinitions[WarningNeedsBelow50].Description,
		)

	case PathwayE2:
		fmt.Printf(
			"[WARNING] %s\n",
			WarningDefinitions[WarningNeedsEqual50].Description,
		)

	case PathwayE3:
		fmt.Printf(
			"[WARNING] %s\n",
			WarningDefinitions[WarningNeedsBetween50And60].Description,
		)

	case PathwayE4:
		fmt.Printf(
			"[WARNING] %s\n",
			WarningDefinitions[WarningNeedsEqual60].Description,
		)

	case PathwayE5:
		fmt.Printf(
			"[WARNING] %s\n",
			WarningDefinitions[WarningNeedsAbove60].Description,
		)
	}

	if assessment.Pathway != "" {
		fmt.Printf("[PATHWAY] %s\n", assessment.Pathway)
	}

	for _, action := range assessment.Actions {
		fmt.Printf(
			"[ACTION] %s\n",
			ActionDefinitions[action].Description,
		)
	}
}

// PrintAllocation prints the calculated allocation for a budget strategy.
func PrintAllocation(strategy BudgetStrategy) {
	fmt.Printf(
		"[CALCULATED] %s: Needs £%.2f | Wants £%.2f | Contribution £%.2f\n",
		strategy.Name,
		float64(strategy.Allocations.Needs)/100,
		float64(strategy.Allocations.Wants)/100,
		float64(strategy.Allocations.Save)/100,
	)
}

// PrintStrategyAvailability prints whether a budget strategy can be used.
func PrintStrategyAvailability(strategy BudgetStrategy) {
	if strategy.Available {
		fmt.Printf("[STRATEGY AVAILABLE] %s\n", strategy.Name)
	} else {
		fmt.Printf("[STRATEGY UNAVAILABLE] %s\n", strategy.Name)
	}
}
