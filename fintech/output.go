package fintech

import (
	"fmt"
	"time"
)

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
	if plan.Available {
		return "Available"
	}

	return "UNAVAILABLE"
}

func availabilityTextEF(plan EmergencyFundPlan) string {
	if plan.Available {
		return "Available"
	}

	return "UNAVAILABLE"
}

func boolText(value bool) string {
	if value {
		return "Yes"
	}

	return "No"
}

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

func PrintEmergencyFundTarget(emergencyFund EmergencyFund) {

	fmt.Println(
		"  Calculating your emergency fund target based on your employment status...",
	)

	fmt.Println()

	fmt.Printf(
		"  Emergency fund target calculated: £%.2f (%d months of essential expenses)\n",
		float64(emergencyFund.TargetAmount)/100,
		emergencyFund.MonthsCount,
	)

	fmt.Println()
	fmt.Println()
}

func PrintDebtFreedomStrategies(
	cf CurrentFinances,
	strategies DebtFreedomStrategies,
) {
	// A plan is grey when DMP is required.
	// A plan is grey when it is unavailable or requires a DMP.
	sustainableGrey := !strategies.Sustainable.Available ||
		strategies.Sustainable.DMPRequired

	moderateGrey := !strategies.Moderate.Available ||
		strategies.Moderate.DMPRequired

	aggressiveGrey := !strategies.Aggressive.Available ||
		strategies.Aggressive.DMPRequired

	formatPlanValue := func(
		value string,
		unavailable bool,
		grey bool,
	) string {
		if unavailable {
			return formatColumnValue("-", grey)
		}

		return formatColumnValue(value, grey)
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
			formatPlanValue(
				formatMoney(sustainable),
				!strategies.Sustainable.Available,
				sustainableGrey,
			),
			formatPlanValue(
				formatMoney(moderate),
				!strategies.Moderate.Available,
				moderateGrey,
			),
			formatPlanValue(
				formatMoney(aggressive),
				!strategies.Aggressive.Available,
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
			formatPlanValue(
				fmt.Sprintf("%d", sustainable),
				!strategies.Sustainable.Available,
				sustainableGrey,
			),
			formatPlanValue(
				fmt.Sprintf("%d", moderate),
				!strategies.Moderate.Available,
				moderateGrey,
			),
			formatPlanValue(
				fmt.Sprintf("%d", aggressive),
				!strategies.Aggressive.Available,
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
			formatPlanValue(
				formatDate(cf.CurrentDate, sustainableMonths),
				!strategies.Sustainable.Available,
				sustainableGrey,
			),
			formatPlanValue(
				formatDate(cf.CurrentDate, moderateMonths),
				!strategies.Moderate.Available,
				moderateGrey,
			),
			formatPlanValue(
				formatDate(cf.CurrentDate, aggressiveMonths),
				!strategies.Aggressive.Available,
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

	fmt.Println()
	fmt.Println("  =================================================================================")
	fmt.Println("                            DEBT FREEDOM PLAN FORECAST")
	fmt.Println("  =================================================================================")
	fmt.Println()

	printHeader()

	fmt.Println("  ---------------------------------------------------------------------------------")

	fmt.Println("    PLAN STATUS")

	fmt.Println("  ---------------------------------------------------------------------------------")

	fmt.Printf(
		"  | %-30s | %s | %s | %s |\n",
		"Availability",
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

	if cf.CurrentSavings < cf.Needs {
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
	}

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

func PrintDebtFreedomCustomComparison(
	cf CurrentFinances,
	strategies DebtFreedomStrategies,
	customPlan DebtFreedomPlan,
) {
	// A plan is grey when it is unavailable or requires a DMP.
	sustainableGrey := !strategies.Sustainable.Available ||
		strategies.Sustainable.DMPRequired

	moderateGrey := !strategies.Moderate.Available ||
		strategies.Moderate.DMPRequired

	aggressiveGrey := !strategies.Aggressive.Available ||
		strategies.Aggressive.DMPRequired

	customGrey := !customPlan.Available ||
		customPlan.DMPRequired

	// A plan shows "-" for calculated values only when unavailable.
	formatPlanValue := func(
		value string,
		unavailable bool,
		grey bool,
	) string {
		if unavailable {
			return formatColumnValue("-", grey)
		}

		return formatColumnValue(value, grey)
	}

	printMoney := func(
		label string,
		sustainable int64,
		moderate int64,
		aggressive int64,
		custom int64,
	) {
		fmt.Printf(
			"  | %-26s | %s | %s | %s | %s |\n",
			label,
			formatPlanValue(
				fmt.Sprintf("%13s", formatMoney(sustainable)),
				!strategies.Sustainable.Available,
				sustainableGrey,
			),
			formatPlanValue(
				fmt.Sprintf("%13s", formatMoney(moderate)),
				!strategies.Moderate.Available,
				moderateGrey,
			),
			formatPlanValue(
				fmt.Sprintf("%13s", formatMoney(aggressive)),
				!strategies.Aggressive.Available,
				aggressiveGrey,
			),
			formatPlanValue(
				fmt.Sprintf("%13s", formatMoney(custom)),
				!customPlan.Available,
				customGrey,
			),
		)
	}

	printInt := func(
		label string,
		sustainable int64,
		moderate int64,
		aggressive int64,
		custom int64,
	) {
		fmt.Printf(
			"  | %-26s | %s | %s | %s | %s |\n",
			label,
			formatPlanValue(
				fmt.Sprintf("%13d", sustainable),
				!strategies.Sustainable.Available,
				sustainableGrey,
			),
			formatPlanValue(
				fmt.Sprintf("%13d", moderate),
				!strategies.Moderate.Available,
				moderateGrey,
			),
			formatPlanValue(
				fmt.Sprintf("%13d", aggressive),
				!strategies.Aggressive.Available,
				aggressiveGrey,
			),
			formatPlanValue(
				fmt.Sprintf("%13d", custom),
				!customPlan.Available,
				customGrey,
			),
		)
	}

	printDate := func(
		label string,
		sustainableMonths int64,
		moderateMonths int64,
		aggressiveMonths int64,
		customMonths int64,
	) {
		fmt.Printf(
			"  | %-26s | %s | %s | %s | %s |\n",
			label,
			formatPlanValue(
				fmt.Sprintf(
					"%13s",
					formatDate(cf.CurrentDate, sustainableMonths),
				),
				!strategies.Sustainable.Available,
				sustainableGrey,
			),
			formatPlanValue(
				fmt.Sprintf(
					"%13s",
					formatDate(cf.CurrentDate, moderateMonths),
				),
				!strategies.Moderate.Available,
				moderateGrey,
			),
			formatPlanValue(
				fmt.Sprintf(
					"%13s",
					formatDate(cf.CurrentDate, aggressiveMonths),
				),
				!strategies.Aggressive.Available,
				aggressiveGrey,
			),
			formatPlanValue(
				fmt.Sprintf(
					"%13s",
					formatDate(cf.CurrentDate, customMonths),
				),
				!customPlan.Available,
				customGrey,
			),
		)
	}

	printHeader := func() {
		fmt.Printf(
			"  | %-26s | %s | %s | %s | %s |\n",
			"",
			formatColumnValue(
				fmt.Sprintf("%13s", "Sustainable"),
				sustainableGrey,
			),
			formatColumnValue(
				fmt.Sprintf("%13s", "Moderate"),
				moderateGrey,
			),
			formatColumnValue(
				fmt.Sprintf("%13s", "Aggressive"),
				aggressiveGrey,
			),
			formatColumnValue(
				fmt.Sprintf("%13s", "Custom"),
				customGrey,
			),
		)
	}

	fmt.Println()
	fmt.Println("  ==============================================================================================")
	fmt.Println("                                DEBT FREEDOM PLAN FORECAST")
	fmt.Println("  ==============================================================================================")
	fmt.Println()

	printHeader()

	fmt.Println("  ----------------------------------------------------------------------------------------------")

	fmt.Println("    PLAN STATUS")

	fmt.Println("  ----------------------------------------------------------------------------------------------")

	fmt.Printf(
		"  | %-26s | %s | %s | %s | %s |\n",
		"Availability",
		formatColumnValue(
			fmt.Sprintf("%13s", availabilityText(strategies.Sustainable)),
			sustainableGrey,
		),
		formatColumnValue(
			fmt.Sprintf("%13s", availabilityText(strategies.Moderate)),
			moderateGrey,
		),
		formatColumnValue(
			fmt.Sprintf("%13s", availabilityText(strategies.Aggressive)),
			aggressiveGrey,
		),
		formatColumnValue(
			fmt.Sprintf("%13s", availabilityText(customPlan)),
			customGrey,
		),
	)

	fmt.Printf(
		"  | %-26s | %s | %s | %s | %s |\n",
		"DMP Required",
		formatColumnValue(
			fmt.Sprintf(
				"%13s",
				boolText(strategies.Sustainable.DMPRequired),
			),
			sustainableGrey,
		),
		formatColumnValue(
			fmt.Sprintf(
				"%13s",
				boolText(strategies.Moderate.DMPRequired),
			),
			moderateGrey,
		),
		formatColumnValue(
			fmt.Sprintf(
				"%13s",
				boolText(strategies.Aggressive.DMPRequired),
			),
			aggressiveGrey,
		),
		formatColumnValue(
			fmt.Sprintf(
				"%13s",
				boolText(customPlan.DMPRequired),
			),
			customGrey,
		),
	)

	if cf.CurrentSavings < cf.Needs {

		fmt.Println("  ----------------------------------------------------------------------------------------------")

		// -------------------------------------------------------------------
		// PHASE 1 — BUFFER
		// -------------------------------------------------------------------

		fmt.Println("    GENERATE EMERGENCY BUFFER")

		fmt.Println("  ----------------------------------------------------------------------------------------------")

		printMoney(
			"Starting Balance",
			cf.CurrentSavings,
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
			bufferRemaining(
				cf.CurrentSavings,
				customPlan.BaselineBuffer.TargetAmount,
			),
		)

		printMoney(
			"Buffer Target",
			strategies.Sustainable.BaselineBuffer.TargetAmount,
			strategies.Moderate.BaselineBuffer.TargetAmount,
			strategies.Aggressive.BaselineBuffer.TargetAmount,
			customPlan.BaselineBuffer.TargetAmount,
		)

		printMoney(
			"Monthly Contribution",
			strategies.Sustainable.Allocation.Save,
			strategies.Moderate.Allocation.Save,
			strategies.Aggressive.Allocation.Save,
			customPlan.Allocation.Save,
		)

		printMoney(
			"Interest Gain",
			strategies.Sustainable.BufferForecast.Phase1InterestGain,
			strategies.Moderate.BufferForecast.Phase1InterestGain,
			strategies.Aggressive.BufferForecast.Phase1InterestGain,
			customPlan.BufferForecast.Phase1InterestGain,
		)

		printMoney(
			"Fees",
			strategies.Sustainable.BufferForecast.Phase1Fees,
			strategies.Moderate.BufferForecast.Phase1Fees,
			strategies.Aggressive.BufferForecast.Phase1Fees,
			customPlan.BufferForecast.Phase1Fees,
		)

		printInt(
			"Total Months",
			strategies.Sustainable.BufferForecast.Phase1Months,
			strategies.Moderate.BufferForecast.Phase1Months,
			strategies.Aggressive.BufferForecast.Phase1Months,
			customPlan.BufferForecast.Phase1Months,
		)

		printDate(
			"Buffer Target Hit Date",
			strategies.Sustainable.BufferForecast.Phase1Months,
			strategies.Moderate.BufferForecast.Phase1Months,
			strategies.Aggressive.BufferForecast.Phase1Months,
			customPlan.BufferForecast.Phase1Months,
		)
	}

	fmt.Println("  ----------------------------------------------------------------------------------------------")

	fmt.Println("    PAY DEBT OFF")

	fmt.Println("  ----------------------------------------------------------------------------------------------")

	printMoney(
		"Starting Debt",
		cf.UnsettledDebt,
		cf.UnsettledDebt,
		cf.UnsettledDebt,
		cf.UnsettledDebt,
	)

	printMoney(
		"Monthly Contribution",
		strategies.Sustainable.Allocation.Save,
		strategies.Moderate.Allocation.Save,
		strategies.Aggressive.Allocation.Save,
		customPlan.Allocation.Save,
	)

	printMoney(
		"Interest Lost",
		totalDebtInterest(strategies.Sustainable),
		totalDebtInterest(strategies.Moderate),
		totalDebtInterest(strategies.Aggressive),
		totalDebtInterest(customPlan),
	)

	printInt(
		"Total Months",
		strategies.Sustainable.DebtForecast.Phase2Months,
		strategies.Moderate.DebtForecast.Phase2Months,
		strategies.Aggressive.DebtForecast.Phase2Months,
		customPlan.DebtForecast.Phase2Months,
	)

	printDate(
		"Debt End Date",
		totalPlanMonths(strategies.Sustainable),
		totalPlanMonths(strategies.Moderate),
		totalPlanMonths(strategies.Aggressive),
		totalPlanMonths(customPlan),
	)

	printMoney(
		"Available Surplus",
		strategies.Sustainable.DebtForecast.Phase2Surplus,
		strategies.Moderate.DebtForecast.Phase2Surplus,
		strategies.Aggressive.DebtForecast.Phase2Surplus,
		customPlan.DebtForecast.Phase2Surplus,
	)

	fmt.Println("  ----------------------------------------------------------------------------------------------")

	fmt.Println("    PROTECTED BUFFER")

	fmt.Println("  ----------------------------------------------------------------------------------------------")

	printMoney(
		"Starting Balance",
		strategies.Sustainable.BaselineBuffer.TargetAmount,
		strategies.Moderate.BaselineBuffer.TargetAmount,
		strategies.Aggressive.BaselineBuffer.TargetAmount,
		customPlan.BaselineBuffer.TargetAmount,
	)

	printMoney(
		"Interest Gain",
		strategies.Sustainable.BufferGrowth.Phase2InterestGain,
		strategies.Moderate.BufferGrowth.Phase2InterestGain,
		strategies.Aggressive.BufferGrowth.Phase2InterestGain,
		customPlan.BufferGrowth.Phase2InterestGain,
	)

	printMoney(
		"Fees",
		strategies.Sustainable.BufferGrowth.Phase2Fees,
		strategies.Moderate.BufferGrowth.Phase2Fees,
		strategies.Aggressive.BufferGrowth.Phase2Fees,
		customPlan.BufferGrowth.Phase2Fees,
	)

	fmt.Println("  ----------------------------------------------------------------------------------------------")

	fmt.Println()
	fmt.Println()

	// -------------------------------------------------------------------
	// PLAN SUMMARY
	// -------------------------------------------------------------------

	fmt.Println("  ==============================================================================================")
	fmt.Println("                                       PLAN SUMMARY")
	fmt.Println("  ==============================================================================================")
	fmt.Println()

	printHeader()

	fmt.Println("  ----------------------------------------------------------------------------------------------")

	printInt(
		"Total Months",
		totalPlanMonths(strategies.Sustainable),
		totalPlanMonths(strategies.Moderate),
		totalPlanMonths(strategies.Aggressive),
		totalPlanMonths(customPlan),
	)

	printDate(
		"End Date",
		totalPlanMonths(strategies.Sustainable),
		totalPlanMonths(strategies.Moderate),
		totalPlanMonths(strategies.Aggressive),
		totalPlanMonths(customPlan),
	)

	printMoney(
		"Total Contributions",
		totalContributions(strategies.Sustainable),
		totalContributions(strategies.Moderate),
		totalContributions(strategies.Aggressive),
		totalContributions(customPlan),
	)

	printMoney(
		"Total Interest Gain",
		totalInterestGain(strategies.Sustainable),
		totalInterestGain(strategies.Moderate),
		totalInterestGain(strategies.Aggressive),
		totalInterestGain(customPlan),
	)

	printMoney(
		"Total Interest Lost",
		totalInterestLost(strategies.Sustainable),
		totalInterestLost(strategies.Moderate),
		totalInterestLost(strategies.Aggressive),
		totalInterestLost(customPlan),
	)

	printMoney(
		"Total Fees",
		totalPlanFees(strategies.Sustainable),
		totalPlanFees(strategies.Moderate),
		totalPlanFees(strategies.Aggressive),
		totalPlanFees(customPlan),
	)

	printMoney(
		"Monthly Contribution",
		strategies.Sustainable.Allocation.Save,
		strategies.Moderate.Allocation.Save,
		strategies.Aggressive.Allocation.Save,
		customPlan.Allocation.Save,
	)

	printMoney(
		"Monthly Wants",
		strategies.Sustainable.Allocation.Wants,
		strategies.Moderate.Allocation.Wants,
		strategies.Aggressive.Allocation.Wants,
		customPlan.Allocation.Wants,
	)

	fmt.Println("  ----------------------------------------------------------------------------------------------")
	fmt.Println()
}

func PrintEmergencyFundStrategies(
	cf CurrentFinances,
	strategies EmergencyFundStrategies,
) {

	// A plan is grey when it is unavailable.
	sustainableGrey := !strategies.Sustainable.Available
	moderateGrey := !strategies.Moderate.Available
	aggressiveGrey := !strategies.Aggressive.Available

	formatPlanValue := func(
		value string,
		unavailable bool,
		grey bool,
	) string {
		if unavailable {
			return formatColumnValue("-", grey)
		}
		return formatColumnValue(value, grey)
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
			formatPlanValue(
				formatMoney(sustainable),
				!strategies.Sustainable.Available,
				sustainableGrey,
			),
			formatPlanValue(
				formatMoney(moderate),
				!strategies.Moderate.Available,
				moderateGrey,
			),
			formatPlanValue(
				formatMoney(aggressive),
				!strategies.Aggressive.Available,
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
			formatPlanValue(
				fmt.Sprintf("%d", sustainable),
				!strategies.Sustainable.Available,
				sustainableGrey,
			),
			formatPlanValue(
				fmt.Sprintf("%d", moderate),
				!strategies.Moderate.Available,
				moderateGrey,
			),
			formatPlanValue(
				fmt.Sprintf("%d", aggressive),
				!strategies.Aggressive.Available,
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
			formatPlanValue(
				formatDate(cf.CurrentDate, sustainableMonths),
				!strategies.Sustainable.Available,
				sustainableGrey,
			),
			formatPlanValue(
				formatDate(cf.CurrentDate, moderateMonths),
				!strategies.Moderate.Available,
				moderateGrey,
			),
			formatPlanValue(
				formatDate(cf.CurrentDate, aggressiveMonths),
				!strategies.Aggressive.Available,
				aggressiveGrey,
			),
		)
	}

	printHeader := func() {
		fmt.Printf(
			"  | %-30s | %s | %s | %s |\n",
			"",
			formatColumnValue(
				"Sustainable",
				sustainableGrey,
			),
			formatColumnValue(
				"Moderate",
				moderateGrey,
			),
			formatColumnValue(
				"Aggressive",
				aggressiveGrey,
			),
		)
	}

	fmt.Println()
	fmt.Println("  =================================================================================")
	fmt.Println("                            EMERGENCY FUND PLAN FORECAST")
	fmt.Println("  =================================================================================")
	fmt.Println()

	printHeader()

	fmt.Println("  ----------------------------------------------------------------------------------")

	fmt.Printf(
		"  | %-30s | %s | %s | %s |\n",
		"PLAN STATUS",
		formatColumnValue(
			availabilityTextEF(strategies.Sustainable),
			sustainableGrey,
		),
		formatColumnValue(
			availabilityTextEF(strategies.Moderate),
			moderateGrey,
		),
		formatColumnValue(
			availabilityTextEF(strategies.Aggressive),
			aggressiveGrey,
		),
	)

	fmt.Println("  ----------------------------------------------------------------------------------")

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

func PrintEmergencyFundCustomComparison(
	cf CurrentFinances,
	strategies EmergencyFundStrategies,
	customPlan EmergencyFundPlan,
) {

	// A plan is grey when it is unavailable.
	sustainableGrey := !strategies.Sustainable.Available
	moderateGrey := !strategies.Moderate.Available
	aggressiveGrey := !strategies.Aggressive.Available
	customGrey := !customPlan.Available

	formatPlanValue := func(
		value string,
		unavailable bool,
		grey bool,
	) string {
		if unavailable {
			return formatColumnValue("-", grey)
		}

		return formatColumnValue(value, grey)
	}

	printMoney := func(
		label string,
		sustainable int64,
		moderate int64,
		aggressive int64,
		custom int64,
	) {
		fmt.Printf(
			"  | %-26s | %13s | %13s | %13s | %13s |\n",
			label,
			formatPlanValue(
				formatMoney(sustainable),
				!strategies.Sustainable.Available,
				sustainableGrey,
			),
			formatPlanValue(
				formatMoney(moderate),
				!strategies.Moderate.Available,
				moderateGrey,
			),
			formatPlanValue(
				formatMoney(aggressive),
				!strategies.Aggressive.Available,
				aggressiveGrey,
			),
			formatPlanValue(
				formatMoney(custom),
				!customPlan.Available,
				customGrey,
			),
		)
	}

	printInt := func(
		label string,
		sustainable int64,
		moderate int64,
		aggressive int64,
		custom int64,
	) {
		fmt.Printf(
			"  | %-26s | %13s | %13s | %13s | %13s |\n",
			label,
			formatPlanValue(
				fmt.Sprintf("%d", sustainable),
				!strategies.Sustainable.Available,
				sustainableGrey,
			),
			formatPlanValue(
				fmt.Sprintf("%d", moderate),
				!strategies.Moderate.Available,
				moderateGrey,
			),
			formatPlanValue(
				fmt.Sprintf("%d", aggressive),
				!strategies.Aggressive.Available,
				aggressiveGrey,
			),
			formatPlanValue(
				fmt.Sprintf("%d", custom),
				!customPlan.Available,
				customGrey,
			),
		)
	}

	printDate := func(
		label string,
		sustainableMonths int64,
		moderateMonths int64,
		aggressiveMonths int64,
		customMonths int64,
	) {
		fmt.Printf(
			"  | %-26s | %13s | %13s | %13s | %13s |\n",
			label,
			formatPlanValue(
				formatDate(cf.CurrentDate, sustainableMonths),
				!strategies.Sustainable.Available,
				sustainableGrey,
			),
			formatPlanValue(
				formatDate(cf.CurrentDate, moderateMonths),
				!strategies.Moderate.Available,
				moderateGrey,
			),
			formatPlanValue(
				formatDate(cf.CurrentDate, aggressiveMonths),
				!strategies.Aggressive.Available,
				aggressiveGrey,
			),
			formatPlanValue(
				formatDate(cf.CurrentDate, customMonths),
				!customPlan.Available,
				customGrey,
			),
		)
	}

	fmt.Println()
	fmt.Println("  ================================================================================================")
	fmt.Println("                                  EMERGENCY FUND PLAN FORECAST")
	fmt.Println("  ================================================================================================")
	fmt.Println()

	fmt.Printf(
		"  | %-26s | %13s | %13s | %13s | %13s |\n",
		"",
		formatColumnValue(
			"Sustainable",
			sustainableGrey,
		),
		formatColumnValue(
			"Moderate",
			moderateGrey,
		),
		formatColumnValue(
			"Aggressive",
			aggressiveGrey,
		),
		formatColumnValue(
			"Custom",
			customGrey,
		),
	)

	fmt.Println("  -----------------------------------------------------------------------------------------------")

	fmt.Printf(
		"  | %-26s | %13s | %13s | %13s | %13s |\n",
		"PLAN STATUS",
		formatColumnValue(
			availabilityTextEF(strategies.Sustainable),
			sustainableGrey,
		),
		formatColumnValue(
			availabilityTextEF(strategies.Moderate),
			moderateGrey,
		),
		formatColumnValue(
			availabilityTextEF(strategies.Aggressive),
			aggressiveGrey,
		),
		formatColumnValue(
			availabilityTextEF(customPlan),
			customGrey,
		),
	)

	fmt.Println("  -----------------------------------------------------------------------------------------------")

	printMoney(
		"Starting Savings",
		cf.CurrentSavings,
		cf.CurrentSavings,
		cf.CurrentSavings,
		cf.CurrentSavings,
	)

	printMoney(
		"Target Amount",
		strategies.Sustainable.TargetAmount,
		strategies.Moderate.TargetAmount,
		strategies.Aggressive.TargetAmount,
		customPlan.TargetAmount,
	)

	printInt(
		"Total Months",
		strategies.Sustainable.Forecast.Phase3Months,
		strategies.Moderate.Forecast.Phase3Months,
		strategies.Aggressive.Forecast.Phase3Months,
		customPlan.Forecast.Phase3Months,
	)

	printDate(
		"EF Target Hit Date",
		strategies.Sustainable.Forecast.Phase3Months,
		strategies.Moderate.Forecast.Phase3Months,
		strategies.Aggressive.Forecast.Phase3Months,
		customPlan.Forecast.Phase3Months,
	)

	printMoney(
		"Total Contributions",
		strategies.Sustainable.Forecast.Phase3Months*
			strategies.Sustainable.Allocation.Save,
		strategies.Moderate.Forecast.Phase3Months*
			strategies.Moderate.Allocation.Save,
		strategies.Aggressive.Forecast.Phase3Months*
			strategies.Aggressive.Allocation.Save,
		customPlan.Forecast.Phase3Months*
			customPlan.Allocation.Save,
	)

	printMoney(
		"Available Surplus",
		strategies.Sustainable.Forecast.Phase3Surplus,
		strategies.Moderate.Forecast.Phase3Surplus,
		strategies.Aggressive.Forecast.Phase3Surplus,
		customPlan.Forecast.Phase3Surplus,
	)

	printMoney(
		"Total Interest Gain",
		strategies.Sustainable.Forecast.Phase3InterestGain,
		strategies.Moderate.Forecast.Phase3InterestGain,
		strategies.Aggressive.Forecast.Phase3InterestGain,
		customPlan.Forecast.Phase3InterestGain,
	)

	printMoney(
		"Total Fees",
		strategies.Sustainable.Forecast.Phase3Fees,
		strategies.Moderate.Forecast.Phase3Fees,
		strategies.Aggressive.Forecast.Phase3Fees,
		customPlan.Forecast.Phase3Fees,
	)

	printMoney(
		"Monthly Contribution",
		strategies.Sustainable.Allocation.Save,
		strategies.Moderate.Allocation.Save,
		strategies.Aggressive.Allocation.Save,
		customPlan.Allocation.Save,
	)

	printMoney(
		"Monthly Wants",
		strategies.Sustainable.Allocation.Wants,
		strategies.Moderate.Allocation.Wants,
		strategies.Aggressive.Allocation.Wants,
		customPlan.Allocation.Wants,
	)

	fmt.Println("  -----------------------------------------------------------------------------------------------")
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
	TotalMonths  int64
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
			TotalMonths:  totalPlanMonths(item.plan),
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

func PrintSavingsOptimisationInfo() {
	fmt.Println()
	fmt.Println("  =================================================================================")
	fmt.Println("                           SAVINGS ACCOUNT OPTIMISATION")
	fmt.Println("  =================================================================================")
	fmt.Println()

	fmt.Println(
		"  We compare Monzo's Instant Access Savings options",
	)
	fmt.Println(
		"  across Free, Extra, Perks and Max plans as your balance changes.",
	)
	fmt.Println()

	fmt.Println(
		"  We review these options each month throughout your forecast",
	)
	fmt.Println(
		"  and automatically use the most beneficial option in our calculations",
	)
	fmt.Println(
		"  as your balance changes.",
	)
	fmt.Println()

	fmt.Println(
		"  This means the savings optimisation is built into your forecast",
	)
	fmt.Println(
		"  and your plan is calculated using the most suitable option",
	)
	fmt.Println(
		"  at each stage of your forecast.",
	)
	fmt.Println()

	fmt.Println(
		"  The aim is not to increase your costs.",
	)
	fmt.Println()

	fmt.Println(
		"  We only consider an alternative account when the",
	)
	fmt.Println(
		"  additional interest is expected to outweigh its fees.",
	)
	fmt.Println()

	fmt.Println(
		"  We also show the balance breakpoints where a different",
	)
	fmt.Println(
		"  option becomes more beneficial, so you can see when",
	)
	fmt.Println(
		"  these changes occur throughout your forecast.",
	)
	fmt.Println()

	fmt.Println(
		"  You can learn more about Monzo Instant Access Savings here:",
	)
	fmt.Println(
		"  https://monzo.com/savings-isas/instant-access",
	)
	fmt.Println()
}

func PrintAvailablePlans(
	cf CurrentFinances,
	plans []SelectablePlan,
) {

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

		if item.Allocation.Save < 10000 && item.TotalMonths >= 6 {
			fmt.Println()
			fmt.Println("     Note:")
			fmt.Println("       Your discretionary money is below £100 per month and")
			fmt.Println("       your debt repayment period is 6 months or longer.")
			fmt.Println("       You may wish to consider whether a Debt Management")
			fmt.Println("       Plan could be appropriate for you.")
		}

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

func PrintExcessSavings(
	cf CurrentFinances,
	emergencyFund EmergencyFund,
	excessSavings ExcessSavingsForecast,
) {

	fmt.Println()
	fmt.Println("  =================================================================================")
	fmt.Println("                            EXCESS SAVINGS IDENTIFIED")
	fmt.Println("  =================================================================================")
	fmt.Println()

	PrintEmergencyFundTarget(emergencyFund)

	fmt.Println(
		"  Your emergency fund is fully covered.",
	)

	fmt.Println()
	fmt.Printf(
		"  You currently have £%.2f in savings.\n",
		float64(cf.CurrentSavings)/100,
	)

	fmt.Printf(
		"  £%.2f is needed to fully cover your emergency fund.\n",
		float64(excessSavings.EmergencyFundAmount)/100,
	)

	fmt.Printf(
		"  This leaves £%.2f available above your emergency fund target.\n",
		float64(excessSavings.InvestmentAmount)/100,
	)

	fmt.Println()
	fmt.Println(
		"  We recommend keeping your emergency fund in an instant-access savings account",
	)
	fmt.Println(
		"  so it remains available when needed.",
	)

	fmt.Println()
	fmt.Printf(
		"  Based on your emergency fund balance, the most advantageous available tier is %s.\n",
		excessSavings.RecommendedTier.Name,
	)

	fmt.Printf(
		"  AER: %.2f%%\n",
		float64(excessSavings.RecommendedTier.AER)/100,
	)

	fmt.Printf(
		"  Monthly fee: £%.2f\n",
		float64(excessSavings.RecommendedTier.Fee)/100,
	)

	fmt.Println()
	fmt.Println(
		"  The remaining excess could be considered for investments",
	)
	fmt.Println(
		"  or other longer-term financial goals, depending on your circumstances.",
	)
}
