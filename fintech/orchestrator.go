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

	var userDecisions UserDecisions

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

				"  Your available savings can cover the required safety buffer,",
			)

			fmt.Println()
			fmt.Println("  Baseline buffer: 1 month of essential expenses")

			fmt.Println(
				"  so paying down high-interest debt is prioritised before forecasting.",
			)

			fmt.Println()
			fmt.Println("  Reviewing available savings for immediate debt repayment...")

			if err := cf.ApplyImmediateDebtPayoff(baselineBuffer); err != nil {
				return err
			}
		}

		if cf.UnsettledDebt == 0 {
			goto debtFreedomDone
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

				fmt.Println()
			}

			if strategies.Moderate.DMPRequired {
				fmt.Println("  Moderate Plan: DMP Required")

				for _, reason := range strategies.Moderate.DMPReasons {
					fmt.Println("    -", reason)
				}

				fmt.Println()
			}

			if strategies.Aggressive.DMPRequired {
				fmt.Println("  Aggressive Plan: DMP Required")

				for _, reason := range strategies.Aggressive.DMPReasons {
					fmt.Println("    -", reason)
				}

				fmt.Println()
			}
		}

		// 4.3. Select Available Plans

		debtPlans := BuildDebtFreedomSelectablePlans(strategies)

		PrintAvailablePlans(
			cf,
			debtPlans,
		)

		var selectedDebtPlan DebtFreedomPlan

		selectedPlan, err := SelectPlan(
			&cf,
			debtPlans,

			func(contribution int64) (SelectablePlan, bool, error) {

				contribution *= 100

				cf.CustomContribution = contribution

				customPlan, err := GenerateCustomDebtFreedomPlan(
					cf,
					baselineBuffer,
					contribution,
				)
				if err != nil {
					return SelectablePlan{}, false, err
				}

				selectedDebtPlan = customPlan
				strategies.Custom = customPlan

				var breakpoints []TierBreakpoint

				breakpoints = append(
					breakpoints,
					customPlan.BufferForecast.TierBreakpoints...,
				)

				for _, breakpoint := range customPlan.BufferGrowth.TierBreakpoints {
					breakpoint.MonthOffset += customPlan.BufferForecast.Phase1Months
					breakpoints = append(breakpoints, breakpoint)
				}

				return SelectablePlan{
					Name:        "Custom Plan",
					Allocation:  customPlan.Allocation,
					Description: "a personalised contribution based on the amount you can afford each month.",
					Breakpoints: breakpoints,
					TotalMonths: totalPlanMonths(customPlan),
				}, !customPlan.DMPRequired, nil
			},

			// Custom plan has just been generated.
			func(customPlan SelectablePlan) {
				PrintDebtFreedomCustomComparison(
					cf,
					strategies,
					selectedDebtPlan,
				)

				if strategies.Moderate.DMPRequired {
					fmt.Println("  Moderate Plan: DMP Required")

					for _, reason := range strategies.Moderate.DMPReasons {
						fmt.Println("    -", reason)
					}

					fmt.Println()
				}

				if strategies.Aggressive.DMPRequired {
					fmt.Println("  Aggressive Plan: DMP Required")

					for _, reason := range strategies.Aggressive.DMPReasons {
						fmt.Println("    -", reason)
					}

					fmt.Println()
				}

				if selectedDebtPlan.DMPRequired {
					fmt.Println("  Custom Plan: DMP Required")

					for _, reason := range selectedDebtPlan.DMPReasons {
						fmt.Println("    -", reason)
					}

					fmt.Println()
				}

			},
		)

		userDecisions.DebtFreedomStrategy =
			TranslateDebtFreedomSelection(selectedPlan)

		// 4.4. Apply Debt Freedom Plan

		switch userDecisions.DebtFreedomStrategy {

		case DebtFreedomSustainable:
			selectedDebtPlan = strategies.Sustainable

		case DebtFreedomModerate:
			selectedDebtPlan = strategies.Moderate

		case DebtFreedomAggressive:
			selectedDebtPlan = strategies.Aggressive

		case DebtFreedomCustom:
			// selectedDebtPlan was generated during custom selection.

		default:
			return fmt.Errorf("invalid debt freedom strategy")
		}

		if err := ApplyDebtFreedomPlan(&cf, selectedDebtPlan); err != nil {
			return err
		}
	}

debtFreedomDone:

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

		excessSavings, err := CalculateExcessSavings(
			cf,
			emergencyFund,
		)
		if err != nil {
			return err
		}

		PrintExcessSavings(
			cf,
			emergencyFund,
			excessSavings,
		)
	}

	// 7. Emergency Fund Forecast

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

		PrintEmergencyFundStrategies(
			cf,
			strategies,
		)

		// 7.2. Select Emergency Fund Plan

		emergencyPlans := BuildEmergencyFundSelectablePlans(strategies)

		PrintAvailablePlans(
			cf,
			emergencyPlans,
		)

		var selectedEmergencyPlan EmergencyFundPlan

		selectedPlan, err := SelectPlan(
			&cf,
			emergencyPlans,

			func(contribution int64) (SelectablePlan, bool, error) {

				contribution *= 100

				cf.CustomContribution = contribution

				customPlan, err := GenerateCustomEmergencyFundPlan(
					cf,
					emergencyFund.TargetAmount,
					contribution,
				)
				if err != nil {
					return SelectablePlan{}, false, err
				}

				selectedEmergencyPlan = customPlan

				return SelectablePlan{
					Name:        "Custom Plan",
					Allocation:  customPlan.Allocation,
					Description: "a personalised contribution based on the amount you can afford each month.",
					Breakpoints: customPlan.Forecast.TierBreakpoints,
				}, true, nil
			},

			// Custom plan generated.
			func(customPlan SelectablePlan) {
				PrintEmergencyFundCustomComparison(
					cf,
					strategies,
					selectedEmergencyPlan,
				)
			},
		)

		userDecisions.EmergencyFundStrategy =
			TranslateEmergencyFundSelection(selectedPlan)

		// 7.3. Apply Emergency Fund Plan

		switch userDecisions.EmergencyFundStrategy {

		case EmergencyFundSustainable:
			selectedEmergencyPlan = strategies.Sustainable

		case EmergencyFundModerate:
			selectedEmergencyPlan = strategies.Moderate

		case EmergencyFundAggressive:
			selectedEmergencyPlan = strategies.Aggressive

		case EmergencyFundCustom:
			// selectedEmergencyPlan was generated during custom selection.

		default:
			return fmt.Errorf("invalid emergency fund strategy")
		}

		if err := ApplyEmergencyFundPlan(&cf, selectedEmergencyPlan); err != nil {
			return err
		}

		return nil
	}

	return nil
}
