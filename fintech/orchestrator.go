package fintech

import "fmt"

// RunPlanner runs the complete financial planning process.
func RunPlanner() error {

	fmt.Println()
	fmt.Println("  =================================================================================")
	fmt.Println("                                   BUDGET PLANNER")
	fmt.Println("  =================================================================================")
	fmt.Println()

	// 1. Collect initial user input
	cf := CollectCurrentFinances()

	allocations, err := GenerateAllocation(cf)
	if err != nil {
		return err
	}

	// 2. Assess deficit position
	fmt.Println()
	fmt.Println("  Assessing your current financial position...")

	fmt.Println()
	fmt.Println("  =================================================================================")
	fmt.Println("                                    WHERE WE ARE")
	fmt.Println("  =================================================================================")
	fmt.Println()
	fmt.Println()

	AssessDeficitPosition(cf)

	if cf.Needs >= cf.Income {
		// deficit / non-growth pathway
		return nil
	}

	// 3. Assess needs position
	AssessNeedsPosition(cf)

	// 4. Debt Freedom

	var baselineBuffer BaselineBuffer

	if cf.HasDebt && cf.UnsettledDebt > 0 {

		var err error

		baselineBuffer, err = CalculateBaselineBuffer(cf)
		if err != nil {
			return err
		}

		//4.1. Immidiate Action

		if cf.CurrentSavings > baselineBuffer.TargetAmount {

			fmt.Println()
			fmt.Println("  =================================================================================")
			fmt.Println("                          IMMEDIETE DEBT PAYMENT AVAILABE")
			fmt.Println("  =================================================================================")
			fmt.Println()
			fmt.Println()
			fmt.Println(

				"Your available savings can cover the required safety buffer,",
			)

			fmt.Println()
			fmt.Println("Baseline buffer: 1 month of essential expenses")

			fmt.Println(
				"so paying down high-interest debt is prioritised before forecasting.",
			)

			fmt.Println()
			fmt.Println("Reviewing available savings for immediate debt repayment...")

			if err := cf.ApplyImmediateDebtPayoff(baselineBuffer); err != nil {
				return err
			}
		}

		if cf.UnsettledDebt == 0 {
			return nil
		}

		// 4.2. Generate Debt Freedom Strategies

		strategies, err := GenerateDebtFreedomStrategies(
			cf,
			allocations,
			baselineBuffer,
		)
		if err != nil {
			return err
		}

		// DMP assessment

		hasDMPPlan :=
			strategies.Sustainable.DMPRequired ||
				strategies.Moderate.DMPRequired ||
				strategies.Aggressive.DMPRequired

		if hasDMPPlan {
			fmt.Println()
			fmt.Println("  DMP ASSESSMENT")
			fmt.Println()

			if strategies.Sustainable.DMPRequired {
				fmt.Println("  Sustainable Plan: DMP Required")

				for _, reason := range strategies.Sustainable.DMPReasons {
					fmt.Println("    -", reason)
				}
			}

			if strategies.Moderate.DMPRequired {
				fmt.Println("  Moderate Plan: DMP Required")

				for _, reason := range strategies.Sustainable.DMPReasons {
					fmt.Println("    -", reason)
				}
			}

			if strategies.Aggressive.DMPRequired {
				fmt.Println("  Aggressive Plan: DMP Required")

				for _, reason := range strategies.Sustainable.DMPReasons {
					fmt.Println("    -", reason)
				}
			}
		}

		// 4.3. Select Debt Freedom Plan

		selectedPlan, err := SelectDebtFreedomPlan(
			cf,
			strategies,
			baselineBuffer,
		)
		if err != nil {
			return err
		}

		// 4.4. Apply Debt Freedom Plan

		if err := ApplyDebtFreedomPlan(&cf, selectedPlan); err != nil {
			return err
		}
	}

	// 5. Calculate emergency fund target

	emergencyFund, err := CalculateEmergencyTarget(
		cf.EmploymentStatus,
		cf.Needs,
	)
	if err != nil {
		return err
	}

	// 6. Excess savings

	if cf.UnsettledDebt == 0 &&
		cf.CurrentSavings > emergencyFund.TargetAmount {

		fmt.Println()
		fmt.Println("EXCESS SAVINGS IDENTIFIED")
		fmt.Println()
		fmt.Println(
			"Your emergency fund is fully covered.",
		)
		fmt.Println(
			"We recommend keeping your emergency fund in an instant-access savings account",
		)
		fmt.Println(
			"so it remains available when needed.",
		)
		fmt.Println()
		fmt.Println(
			"Any savings above your emergency fund can be considered for",
		)
		fmt.Println(
			"other investment or longer-term financial goals.",
		)

		excessSavings, err := CalculateExcessSavings(
			cf,
			emergencyFund,
		)
		if err != nil {
			return err
		}

		_ = excessSavings
	}

	// 7. Emergency Fund

	if cf.UnsettledDebt == 0 &&
		cf.CurrentSavings < emergencyFund.TargetAmount {

		// 7.1. Generate Emergency Fund Strategies

		strategies, err := GenerateEmergencyFundStrategies(
			cf,
			allocations,
			emergencyFund.TargetAmount,
		)
		if err != nil {
			return err
		}

		// 7.2. Select Emergency Fund Plan

		selectedPlan, err := SelectEmergencyFundPlan(
			cf,
			strategies,
			emergencyFund.TargetAmount,
		)
		if err != nil {
			return err
		}

		// 7.3. Apply Emergency Fund Plan

		if err := ApplyEmergencyFundPlan(&cf, selectedPlan); err != nil {
			return err
		}
	}

	return nil
}
