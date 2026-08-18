package fintech

import (
	"fmt"
	"time"
)

func PrintCurrentFinances(cf CurrentFinances) {
	fmt.Println()
	fmt.Println()
	fmt.Println("  |-------------------------------------------------------|")
	fmt.Printf(
		"  |          FINANCIAL POSITION — %-20s    |\n",
		cf.CurrentDate.Format("January 2006"),
	)
	fmt.Println("  |-------------------------------------------------------|")

	fmt.Printf(
		"  | %-32s | %-18s |\n",
		"Income",
		fmt.Sprintf("£%.2f", float64(cf.Income)/100),
	)

	fmt.Printf(
		"  | %-32s | %-18s |\n",
		"Needs",
		fmt.Sprintf("£%.2f", float64(cf.Needs)/100),
	)

	fmt.Printf(
		"  | %-32s | %-18s |\n",
		"Current Savings",
		fmt.Sprintf("£%.2f", float64(cf.CurrentSavings)/100),
	)

	fmt.Printf(
		"  | %-32s | %-18s |\n",
		"Unsettled Debt",
		fmt.Sprintf("£%.2f", float64(cf.UnsettledDebt)/100),
	)

	fmt.Printf(
		"  | %-32s | %-18s |\n",
		"Employment Status",
		string(cf.EmploymentStatus),
	)

	fmt.Println("  |----------------------------------|--------------------|")
	fmt.Println()

}

func PrintDebtFreedomStrategies(
	cf CurrentFinances,
	strategies DebtFreedomStrategies,
) {
	// A plan is grey when DMP is required.
	sustainableGrey := strategies.Sustainable.DMPRequired
	moderateGrey := strategies.Moderate.DMPRequired
	aggressiveGrey := strategies.Aggressive.DMPRequired

	// Small local helpers so we don't have to pass the grey status
	// into every single row manually.
	printMoney := func(
		label string,
		sustainable int64,
		moderate int64,
		aggressive int64,
	) {
		fmt.Printf(
			"  | %-30s | %s | %s | %s |\n",
			label,
			formatColumnValue(
				formatMoney(sustainable),
				sustainableGrey,
			),
			formatColumnValue(
				formatMoney(moderate),
				moderateGrey,
			),
			formatColumnValue(
				formatMoney(aggressive),
				aggressiveGrey,
			),
		)
	}

	printInt := func(
		label string,
		sustainable int64,
		moderate int64,
		aggressive int64,
	) {
		fmt.Printf(
			"  | %-30s | %s | %s | %s |\n",
			label,
			formatColumnValue(
				fmt.Sprintf("%d", sustainable),
				sustainableGrey,
			),
			formatColumnValue(
				fmt.Sprintf("%d", moderate),
				moderateGrey,
			),
			formatColumnValue(
				fmt.Sprintf("%d", aggressive),
				aggressiveGrey,
			),
		)
	}

	printDate := func(
		label string,
		sustainableMonths int64,
		moderateMonths int64,
		aggressiveMonths int64,
	) {
		fmt.Printf(
			"  | %-30s | %s | %s | %s |\n",
			label,
			formatColumnValue(
				formatDate(cf.CurrentDate, sustainableMonths),
				sustainableGrey,
			),
			formatColumnValue(
				formatDate(cf.CurrentDate, moderateMonths),
				moderateGrey,
			),
			formatColumnValue(
				formatDate(cf.CurrentDate, aggressiveMonths),
				aggressiveGrey,
			),
		)
	}

	printHeader := func() {
		fmt.Printf(
			"  | %-30s | %s | %s | %s |\n",
			"",
			formatColumnValue("Sustainable", sustainableGrey),
			formatColumnValue("Moderate", moderateGrey),
			formatColumnValue("Aggressive", aggressiveGrey),
		)
	}

	printHeader()

	fmt.Println("  ---------------------------------------------------------------------------------")

	fmt.Println("    PLAN STATUS")

	fmt.Println("  ---------------------------------------------------------------------------------")

	fmt.Printf(
		"  | %-30s | %s | %s | %s |\n",
		"Available",
		formatColumnValue(
			availabilityText(strategies.Sustainable),
			sustainableGrey,
		),
		formatColumnValue(
			availabilityText(strategies.Moderate),
			moderateGrey,
		),
		formatColumnValue(
			availabilityText(strategies.Aggressive),
			aggressiveGrey,
		),
	)

	fmt.Printf(
		"  | %-30s | %s | %s | %s |\n",
		"DMP Required",
		formatColumnValue(
			boolText(strategies.Sustainable.DMPRequired),
			sustainableGrey,
		),
		formatColumnValue(
			boolText(strategies.Moderate.DMPRequired),
			moderateGrey,
		),
		formatColumnValue(
			boolText(strategies.Aggressive.DMPRequired),
			aggressiveGrey,
		),
	)

	fmt.Println("  ---------------------------------------------------------------------------------")

	// -------------------------------------------------------------------
	// PHASE 1 — BUFFER
	// -------------------------------------------------------------------

	fmt.Println("    GENERATE EMERGENCY BUFFER")

	fmt.Println("  ---------------------------------------------------------------------------------")

	printMoney(
		"Starting Balance",
		cf.CurrentSavings,
		cf.CurrentSavings,
		cf.CurrentSavings,
	)

	printMoney(
		"Remaining to Buffer Target",
		bufferRemaining(
			cf.CurrentSavings,
			strategies.Sustainable.BaselineBuffer.TargetAmount,
		),
		bufferRemaining(
			cf.CurrentSavings,
			strategies.Moderate.BaselineBuffer.TargetAmount,
		),
		bufferRemaining(
			cf.CurrentSavings,
			strategies.Aggressive.BaselineBuffer.TargetAmount,
		),
	)

	printMoney(
		"Buffer Target",
		strategies.Sustainable.BaselineBuffer.TargetAmount,
		strategies.Moderate.BaselineBuffer.TargetAmount,
		strategies.Aggressive.BaselineBuffer.TargetAmount,
	)

	printMoney(
		"Monthly Contribution",
		strategies.Sustainable.Allocation.Save,
		strategies.Moderate.Allocation.Save,
		strategies.Aggressive.Allocation.Save,
	)

	printMoney(
		"Interest Gain",
		strategies.Sustainable.BufferForecast.Phase1InterestGain,
		strategies.Moderate.BufferForecast.Phase1InterestGain,
		strategies.Aggressive.BufferForecast.Phase1InterestGain,
	)

	printMoney(
		"Fees",
		strategies.Sustainable.BufferForecast.Phase1Fees,
		strategies.Moderate.BufferForecast.Phase1Fees,
		strategies.Aggressive.BufferForecast.Phase1Fees,
	)

	printInt(
		"Total Months",
		strategies.Sustainable.BufferForecast.Phase1Months,
		strategies.Moderate.BufferForecast.Phase1Months,
		strategies.Aggressive.BufferForecast.Phase1Months,
	)

	printDate(
		"Buffer Target Hit Date",
		strategies.Sustainable.BufferForecast.Phase1Months,
		strategies.Moderate.BufferForecast.Phase1Months,
		strategies.Aggressive.BufferForecast.Phase1Months,
	)

	fmt.Println("  ---------------------------------------------------------------------------------")

	fmt.Println("    PAY DEBT OFF")

	fmt.Println("  ---------------------------------------------------------------------------------")

	printMoney(
		"Starting Debt",
		cf.UnsettledDebt,
		cf.UnsettledDebt,
		cf.UnsettledDebt,
	)

	printMoney(
		"Monthly Contribution",
		strategies.Sustainable.Allocation.Save,
		strategies.Moderate.Allocation.Save,
		strategies.Aggressive.Allocation.Save,
	)

	printMoney(
		"Interest Lost",
		totalDebtInterest(strategies.Sustainable),
		totalDebtInterest(strategies.Moderate),
		totalDebtInterest(strategies.Aggressive),
	)

	printInt(
		"Total Months",
		strategies.Sustainable.DebtForecast.Phase2Months,
		strategies.Moderate.DebtForecast.Phase2Months,
		strategies.Aggressive.DebtForecast.Phase2Months,
	)

	printDate(
		"Debt End Date",
		totalPlanMonths(strategies.Sustainable),
		totalPlanMonths(strategies.Moderate),
		totalPlanMonths(strategies.Aggressive),
	)

	printMoney(
		"Available Surplus",
		strategies.Sustainable.DebtForecast.Phase2Surplus,
		strategies.Moderate.DebtForecast.Phase2Surplus,
		strategies.Aggressive.DebtForecast.Phase2Surplus,
	)

	fmt.Println("  ---------------------------------------------------------------------------------")

	fmt.Println("    PROTECTED BUFFER")

	fmt.Println("  ---------------------------------------------------------------------------------")

	printMoney(
		"Starting Balance",
		strategies.Sustainable.BaselineBuffer.TargetAmount,
		strategies.Moderate.BaselineBuffer.TargetAmount,
		strategies.Aggressive.BaselineBuffer.TargetAmount,
	)

	printMoney(
		"Interest Gain",
		strategies.Sustainable.BufferGrowth.Phase2InterestGain,
		strategies.Moderate.BufferGrowth.Phase2InterestGain,
		strategies.Aggressive.BufferGrowth.Phase2InterestGain,
	)

	printMoney(
		"Fees",
		strategies.Sustainable.BufferGrowth.Phase2Fees,
		strategies.Moderate.BufferGrowth.Phase2Fees,
		strategies.Aggressive.BufferGrowth.Phase2Fees,
	)

	fmt.Println("  ---------------------------------------------------------------------------------")

	fmt.Println()
	fmt.Println()

	// -------------------------------------------------------------------
	// PLAN SUMMARY
	// -------------------------------------------------------------------

	fmt.Println("  =================================================================================")
	fmt.Println("                                    PLAN SUMMARY")
	fmt.Println("  =================================================================================")
	fmt.Println()

	printHeader()

	fmt.Println("  ---------------------------------------------------------------------------------")

	printInt(
		"Total Months",
		totalPlanMonths(strategies.Sustainable),
		totalPlanMonths(strategies.Moderate),
		totalPlanMonths(strategies.Aggressive),
	)

	printDate(
		"End Date",
		totalPlanMonths(strategies.Sustainable),
		totalPlanMonths(strategies.Moderate),
		totalPlanMonths(strategies.Aggressive),
	)

	printMoney(
		"Total Contributions",
		totalContributions(strategies.Sustainable),
		totalContributions(strategies.Moderate),
		totalContributions(strategies.Aggressive),
	)

	printMoney(
		"Total Interest Gain",
		totalInterestGain(strategies.Sustainable),
		totalInterestGain(strategies.Moderate),
		totalInterestGain(strategies.Aggressive),
	)

	printMoney(
		"Total Interest Lost",
		totalInterestLost(strategies.Sustainable),
		totalInterestLost(strategies.Moderate),
		totalInterestLost(strategies.Aggressive),
	)

	printMoney(
		"Total Fees",
		totalPlanFees(strategies.Sustainable),
		totalPlanFees(strategies.Moderate),
		totalPlanFees(strategies.Aggressive),
	)

	printMoney(
		"Monthly Contribution",
		strategies.Sustainable.Allocation.Save,
		strategies.Moderate.Allocation.Save,
		strategies.Aggressive.Allocation.Save,
	)

	printMoney(
		"Monthly Wants",
		strategies.Sustainable.Allocation.Wants,
		strategies.Moderate.Allocation.Wants,
		strategies.Aggressive.Allocation.Wants,
	)

	fmt.Println("  ---------------------------------------------------------------------------------")
	fmt.Println()
}

// CALCULATION HELPERS

func bufferRemaining(startingBalance, target int64) int64 {
	remaining := target - startingBalance

	if remaining < 0 {
		return 0
	}

	return remaining
}

func totalPlanMonths(plan DebtFreedomPlan) int64 {
	return plan.BufferForecast.Phase1Months +
		plan.DebtForecast.Phase2Months
}

func totalContributions(plan DebtFreedomPlan) int64 {
	return totalPlanMonths(plan)*plan.Allocation.Save -
		plan.DebtForecast.Phase2Surplus
}

func totalInterestGain(plan DebtFreedomPlan) int64 {
	return plan.BufferForecast.Phase1InterestGain +
		plan.BufferGrowth.Phase2InterestGain
}

func totalInterestLost(plan DebtFreedomPlan) int64 {
	return plan.DebtForecast.Phase1InterestLost +
		plan.DebtForecast.Phase2InterestLost
}

func totalDebtInterest(plan DebtFreedomPlan) int64 {
	return totalInterestLost(plan)
}

func totalPlanFees(plan DebtFreedomPlan) int64 {
	return plan.BufferForecast.Phase1Fees +
		plan.BufferGrowth.Phase2Fees
}

// FORMATTING HELPERS

const (
	grey  = "\033[90m"
	reset = "\033[0m"
)

func formatMoney(amount int64) string {
	return fmt.Sprintf("£%.2f", float64(amount)/100)
}

func formatDate(
	currentDate time.Time,
	months int64,
) string {
	return currentDate.
		AddDate(0, int(months), 0).
		Format("Jan 2006")
}

func formatColumnValue(value string, greyColumn bool) string {
	padded := fmt.Sprintf("%-13s", value)

	if greyColumn {
		return grey + padded + reset
	}

	return padded
}

func availabilityText(plan DebtFreedomPlan) string {
	if plan.DMPRequired {
		return "Not Available"
	}

	return "Available"
}

func boolText(value bool) string {
	if value {
		return "Yes"
	}

	return "No"
}

func PrintEmergencyFundStrategies(
	cf CurrentFinances,
	strategies EmergencyFundStrategies,
) {
	printValue := func(value string, available bool) string {
		if !available {
			return "Not Available"
		}

		return value
	}

	printMoney := func(
		label string,
		sustainable int64,
		moderate int64,
		aggressive int64,
	) {
		fmt.Printf(
			"  | %-30s | %s | %s | %s |\n",
			label,
			formatColumnValue(
				printValue(
					formatMoney(sustainable),
					strategies.Sustainable.Available,
				),
				!strategies.Sustainable.Available,
			),
			formatColumnValue(
				printValue(
					formatMoney(moderate),
					strategies.Moderate.Available,
				),
				!strategies.Moderate.Available,
			),
			formatColumnValue(
				printValue(
					formatMoney(aggressive),
					strategies.Aggressive.Available,
				),
				!strategies.Aggressive.Available,
			),
		)
	}

	printInt := func(
		label string,
		sustainable int64,
		moderate int64,
		aggressive int64,
	) {
		fmt.Printf(
			"  | %-30s | %s | %s | %s |\n",
			label,
			formatColumnValue(
				printValue(
					fmt.Sprintf("%d", sustainable),
					strategies.Sustainable.Available,
				),
				!strategies.Sustainable.Available,
			),
			formatColumnValue(
				printValue(
					fmt.Sprintf("%d", moderate),
					strategies.Moderate.Available,
				),
				!strategies.Moderate.Available,
			),
			formatColumnValue(
				printValue(
					fmt.Sprintf("%d", aggressive),
					strategies.Aggressive.Available,
				),
				!strategies.Aggressive.Available,
			),
		)
	}

	printDate := func(
		label string,
		sustainableMonths int64,
		moderateMonths int64,
		aggressiveMonths int64,
	) {
		fmt.Printf(
			"  | %-30s | %s | %s | %s |\n",
			label,
			formatColumnValue(
				printValue(
					formatDate(cf.CurrentDate, sustainableMonths),
					strategies.Sustainable.Available,
				),
				!strategies.Sustainable.Available,
			),
			formatColumnValue(
				printValue(
					formatDate(cf.CurrentDate, moderateMonths),
					strategies.Moderate.Available,
				),
				!strategies.Moderate.Available,
			),
			formatColumnValue(
				printValue(
					formatDate(cf.CurrentDate, aggressiveMonths),
					strategies.Aggressive.Available,
				),
				!strategies.Aggressive.Available,
			),
		)
	}

	printHeader := func() {
		fmt.Printf(
			"  | %-30s | %s | %s | %s |\n",
			"",
			formatColumnValue(
				"Sustainable",
				!strategies.Sustainable.Available,
			),
			formatColumnValue(
				"Moderate",
				!strategies.Moderate.Available,
			),
			formatColumnValue(
				"Aggressive",
				!strategies.Aggressive.Available,
			),
		)
	}

	printHeader()

	fmt.Println("  ---------------------------------------------------------------------------------")
	fmt.Println("    GENERATE EMERGENCY FUND")
	fmt.Println("  ---------------------------------------------------------------------------------")

	printMoney(
		"Starting Savings",
		cf.CurrentSavings,
		cf.CurrentSavings,
		cf.CurrentSavings,
	)

	printMoney(
		"Target Amount",
		strategies.Sustainable.TargetAmount,
		strategies.Moderate.TargetAmount,
		strategies.Aggressive.TargetAmount,
	)

	printInt(
		"Total Months",
		strategies.Sustainable.Forecast.Phase3Months,
		strategies.Moderate.Forecast.Phase3Months,
		strategies.Aggressive.Forecast.Phase3Months,
	)

	printDate(
		"EF Target Hit Date",
		strategies.Sustainable.Forecast.Phase3Months,
		strategies.Moderate.Forecast.Phase3Months,
		strategies.Aggressive.Forecast.Phase3Months,
	)

	printMoney(
		"Total Contributions",
		strategies.Sustainable.Forecast.Phase3Months*
			strategies.Sustainable.Allocation.Save,
		strategies.Moderate.Forecast.Phase3Months*
			strategies.Moderate.Allocation.Save,
		strategies.Aggressive.Forecast.Phase3Months*
			strategies.Aggressive.Allocation.Save,
	)

	printMoney(
		"Available Surplus",
		strategies.Sustainable.Forecast.Phase3Surplus,
		strategies.Moderate.Forecast.Phase3Surplus,
		strategies.Aggressive.Forecast.Phase3Surplus,
	)

	printMoney(
		"Total Interest Gain",
		strategies.Sustainable.Forecast.Phase3InterestGain,
		strategies.Moderate.Forecast.Phase3InterestGain,
		strategies.Aggressive.Forecast.Phase3InterestGain,
	)

	printMoney(
		"Total Fees",
		strategies.Sustainable.Forecast.Phase3Fees,
		strategies.Moderate.Forecast.Phase3Fees,
		strategies.Aggressive.Forecast.Phase3Fees,
	)

	printMoney(
		"Monthly Contribution",
		strategies.Sustainable.Allocation.Save,
		strategies.Moderate.Allocation.Save,
		strategies.Aggressive.Allocation.Save,
	)

	printMoney(
		"Monthly Wants",
		strategies.Sustainable.Allocation.Wants,
		strategies.Moderate.Allocation.Wants,
		strategies.Aggressive.Allocation.Wants,
	)

	fmt.Println("  ---------------------------------------------------------------------------------")
	fmt.Println()
}

func printTierOptimisation(
	cf CurrentFinances,
	breakpoints []TierBreakpoint,
) {
	if len(breakpoints) == 0 {
		return
	}

	// The first breakpoint is always Month 0,
	// recommended starting account.
	startingTier := breakpoints[0]

	fmt.Println("     Recommended starting account:")
	fmt.Println()

	fmt.Printf("       %s\n", startingTier.Tier.Name)
	fmt.Printf(
		"       AER: %.2f%%\n",
		float64(startingTier.Tier.AER)/100,
	)
	fmt.Printf(
		"       Fee: £%.2f\n",
		float64(startingTier.Tier.Fee)/100,
	)

	fmt.Println()

	// Everything after the first breakpoint represents
	// a change to a better tier.
	if len(breakpoints) == 1 {
		fmt.Println("     No breakpoints reached.")
		fmt.Println(
			"     Savings remain in the starting account throughout the forecast.",
		)
		return
	}

	fmt.Println("     Breakpoints:")
	fmt.Println()

	for _, breakpoint := range breakpoints[1:] {
		fmt.Printf(
			"       %s\n",
			formatDate(cf.CurrentDate, breakpoint.MonthOffset),
		)

		fmt.Printf(
			"       Balance: £%.2f\n",
			float64(breakpoint.Balance)/100,
		)

		fmt.Printf(
			"       → %s becomes more beneficial\n",
			breakpoint.Tier.Name,
		)

		fmt.Println()
	}
}

type SelectablePlan struct {
	Name         string
	Allocation   Allocation
	WantsPercent int64
	Description  string
	Breakpoints  []TierBreakpoint
}

func BuildDebtFreedomSelectablePlans(
	strategies DebtFreedomStrategies,
) []SelectablePlan {

	plans := []SelectablePlan{}

	debtPlans := []struct {
		name         string
		plan         DebtFreedomPlan
		wantsPercent int64
		description  string
	}{
		{
			name:         "Sustainable Plan",
			plan:         strategies.Sustainable,
			wantsPercent: 30,
			description:  "leaving more flexibility for lifestyle spending while prioritising steady debt repayment and financial progress.",
		},
		{
			name:         "Moderate Plan",
			plan:         strategies.Moderate,
			wantsPercent: 25,
			description:  "balancing lifestyle spending with a stronger contribution towards debt repayment and financial progress.",
		},
		{
			name:         "Aggressive Plan",
			plan:         strategies.Aggressive,
			wantsPercent: 20,
			description:  "prioritising faster debt repayment by directing more of the remaining income towards financial progress.",
		},
	}

	for _, item := range debtPlans {

		// The plan cannot be selected if it is unavailable
		// or requires a Debt Management Plan.
		if !item.plan.Available || item.plan.DMPRequired {
			continue
		}

		var breakpoints []TierBreakpoint

		// Phase 1: generating buffer
		breakpoints = append(
			breakpoints,
			item.plan.BufferForecast.TierBreakpoints...,
		)

		// Phase 2: debt repayment timeline.
		// BufferGrowth also includes the interest earned during this phase.
		// Phase 2 starts from Month 0, so convert its local offset
		// into the overall Debt Freedom timeline.
		for _, breakpoint := range item.plan.BufferGrowth.TierBreakpoints {
			breakpoint.MonthOffset += item.plan.BufferForecast.Phase1Months
			breakpoints = append(breakpoints, breakpoint)
		}

		plans = append(plans, SelectablePlan{
			Name:         item.name,
			Allocation:   item.plan.Allocation,
			WantsPercent: item.wantsPercent,
			Description:  item.description,
			Breakpoints:  breakpoints,
		})
	}

	return plans
}

func BuildEmergencyFundSelectablePlans(
	strategies EmergencyFundStrategies,
) []SelectablePlan {

	plans := []SelectablePlan{}

	emergencyPlans := []struct {
		name         string
		plan         EmergencyFundPlan
		wantsPercent int64
		description  string
	}{
		{
			name:         "Sustainable Plan",
			plan:         strategies.Sustainable,
			wantsPercent: 30,
			description:  "leaving more flexibility while steadily building your financial safety net.",
		},
		{
			name:         "Moderate Plan",
			plan:         strategies.Moderate,
			wantsPercent: 25,
			description:  "balancing lifestyle spending with a stronger contribution towards your financial safety net.",
		},
		{
			name:         "Aggressive Plan",
			plan:         strategies.Aggressive,
			wantsPercent: 20,
			description:  "prioritising financial safety by directing more of the remaining income towards your emergency fund.",
		},
	}

	for _, item := range emergencyPlans {

		if !item.plan.Available {
			continue
		}

		breakpoints := item.plan.Forecast.TierBreakpoints

		plans = append(plans, SelectablePlan{
			Name:         item.name,
			Allocation:   item.plan.Allocation,
			WantsPercent: item.wantsPercent,
			Description:  item.description,
			Breakpoints:  breakpoints,
		})
	}

	return plans
}

func PrintAvailablePlans(
	cf CurrentFinances,
	plans []SelectablePlan,
) {
	fmt.Println()

	fmt.Println(
		"  The figures in the table below include the savings optimisation",
	)
	fmt.Println(
		"  calculations described here.",
	)
	fmt.Println()

	fmt.Println(
		"  To optimise the protected buffer, we compare available",
	)
	fmt.Println(
		"  instant-access savings tiers each month as the balance changes.",
	)
	fmt.Println()

	fmt.Println("  The aim is not to increase your costs.")
	fmt.Println()

	fmt.Println(
		"  We only consider an alternative account when the",
	)
	fmt.Println(
		"  additional interest is expected to outweigh its fees.",
	)
	fmt.Println()

	planNumber := 1

	for _, item := range plans {
		fmt.Printf(
			"  %d. %s\n",
			planNumber,
			item.Name,
		)

		fmt.Printf(
			"     %d%% of income is allocated to wants, %s\n",
			item.WantsPercent,
			item.Description,
		)

		fmt.Println()
		fmt.Println("     Summary:")

		fmt.Printf(
			"       Needs:                 £%.2f\n",
			float64(item.Allocation.Needs)/100,
		)

		fmt.Printf(
			"       Wants:                 £%.2f\n",
			float64(item.Allocation.Wants)/100,
		)

		fmt.Printf(
			"       Monthly contribution:  £%.2f\n",
			float64(item.Allocation.Save)/100,
		)

		fmt.Println()

		printTierOptimisation(
			cf,
			item.Breakpoints,
		)

		fmt.Println()

		planNumber++
	}

	if planNumber == 1 {
		fmt.Println("  No default plans are currently available.")
		fmt.Println()
		fmt.Println("  Please tell us how much you can contribute each month")
		fmt.Println("  to create a Custom Plan.")
		fmt.Println()
	}
}
