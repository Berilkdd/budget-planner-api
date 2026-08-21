package fintech

import "fmt"

type DMPAssessment struct {
	DMPRequired bool
	Reasons     []string
}

// AssessDMPNeed checks the aggressive debt forecast against the DMP safety thresholds.
func AssessDMPNeed(
	cf CurrentFinances,
	plan DebtFreedomPlan,
) DMPAssessment {
	var reasons []string

	if plan.DebtForecast.InterestOver50Percent {
		reasons = append(
			reasons,
			"More than 50% of the monthly debt contribution is being consumed by interest.",
		)
	}

	if cf.UnsettledDebt > (cf.Income * 12 / 2) {
		reasons = append(
			reasons,
			"Debt exceeds 50% of annual income.",
		)
	}

	totalDebtMonths :=
		plan.BufferForecast.Phase1Months +
			plan.DebtForecast.Phase2Months

	if totalDebtMonths > 36 {
		reasons = append(
			reasons,
			"Debt payoff takes more than 36 months.",
		)
	}

	return DMPAssessment{
		DMPRequired: len(reasons) > 0,
		Reasons:     reasons,
	}
}

func AssessDeficitPosition(cf CurrentFinances) {

	if cf.Needs < cf.Income {
		return
	}

	fmt.Println(
		"  Your essential expenses currently meet or exceed your monthly income.\n\n" +
			"  This isn't necessarily a sign that something is wrong.\n" +
			"  Your current situation may be temporary, you could be between jobs,\n" +
			"  starting a new career, studying, building a new business,\n" +
			"  or going through a period of higher essential costs.\n\n" +
			"  Because your income currently doesn't fully cover your essential expenses,\n" +
			"  we look for a stronger financial safety net to help protect you \n" +
			"  while your situation is changing.",
	)
	fmt.Println()

	// Calculate the required emergency fund.
	emergencyFund, err := CalculateEmergencyTarget(
		cf.EmploymentStatus,
		cf.Needs,
	)

	PrintEmergencyFundTarget(emergencyFund)

	if err != nil {
		return
	}

	// 1. Savings >= Emergency Fund + Debt
	if cf.UnsettledDebt > 0 &&
		cf.CurrentSavings >= emergencyFund.TargetAmount+cf.UnsettledDebt {

		fmt.Println(
			"  Good news! Your savings are enough to cover your recommended emergency fund and \n" +
				"  clear your outstanding debt.\n\n" +
				"  If your debt carries a high interest rate,\n" +
				"  prioritising repayment could help reduce the amount of interest you pay over time.\n\n" +
				fmt.Sprintf(
					"  You could use £%.2f of your savings to clear the debt,\n"+
						"  leaving approximately £%.2f in savings\n"+
						"  while keeping your recommended emergency fund fully covered.",
					float64(cf.UnsettledDebt)/100,
					float64(cf.CurrentSavings-cf.UnsettledDebt)/100,
				),
		)
		fmt.Println()
		fmt.Println()
	}

	// 2. Savings < Debt + EF && Savings > EF
	if cf.CurrentSavings < emergencyFund.TargetAmount+cf.UnsettledDebt &&
		cf.CurrentSavings > emergencyFund.TargetAmount {

		availableForDebt := cf.CurrentSavings - emergencyFund.TargetAmount
		remainingDebt := cf.UnsettledDebt - availableForDebt

		fmt.Println(
			"  Your savings cover your recommended emergency fund,\n" +
				"  but they are not enough to clear your outstanding debt in full\n" +
				"  while keeping this safety buffer protected.\n\n" +
				"  If your debt carries a high interest rate,\n" +
				"  you may want to prioritise a portion of your available savings towards repayment.\n\n" +
				fmt.Sprintf(
					"  You could use £%.2f towards your debt while keeping your emergency fund protected.\n"+
						"  This would leave approximately £%.2f of debt remaining.",
					float64(availableForDebt)/100,
					float64(remainingDebt)/100,
				),
		)
		fmt.Println()
		fmt.Println()
	}

	// 3. Savings < EF && Debt > 0
	if cf.CurrentSavings < emergencyFund.TargetAmount &&
		cf.UnsettledDebt > 0 {

		fmt.Println(
			"  Your current savings do not yet cover your recommended emergency fund,\n" +
				"  and you also have outstanding debt.\n\n" +
				"  Because your essential expenses currently meet or exceed your income,\n" +
				"  you may need to use your savings to cover the gap.\n" +
				"  Continuing to use your savings in this way could reduce your financial safety net over time.\n\n" +
				"  At the same time, interest on your debt may continue to increase the amount you owe.\n" +
				"  we recommend seeking debt advice before making any debt repayments from your savings,\n" +
				"  rather than using your limited savings to repay the debt and \n" +
				"  leaving yourself without a financial safety net.\n\n" +
				"  You can also seek support to explore ways to reduce your essential expenses or increase your income,\n" +
				"  so that you can move away from relying on your savings to cover your monthly needs.",
		)
		fmt.Println()
		fmt.Println()

	}

	// 4. Savings == EF && Debt > 0
	if cf.CurrentSavings == emergencyFund.TargetAmount &&
		cf.UnsettledDebt > 0 {

		fmt.Println(
			"  Your savings currently cover your recommended emergency fund,\n" +
				"  but you do not have additional savings available to repay your outstanding debt.\n\n" +
				"  Because your essential expenses currently meet or exceed your income,\n" +
				"  this safety buffer may be needed to cover unexpected costs \n" +
				"  or a continued gap between your income and essential expenses.\n\n" +
				"  Because you are currently at risk of running a deficit,\n" +
				"  we recommend seeking debt advice before making any debt repayments from your savings.\n\n" +
				"  If your debt carries a high interest rate,\n" +
				"  seeking debt support can help you explore ways of managing the debt\n" +
				"  without putting your financial safety net at risk.\n\n" +
				"  You can also look at ways to reduce your essential expenses or increase your income,\n" +
				"  so that you are no longer relying on your savings to cover your monthly needs.",
		)
		fmt.Println()
		fmt.Println()
	}

	// 5. EF < CurrentSavings && Debt == 0
	if emergencyFund.TargetAmount < cf.CurrentSavings &&
		cf.UnsettledDebt == 0 {

		fmt.Println(
			"  Good news! You currently have no outstanding debt,\n" +
				"  and your savings cover your recommended emergency fund.\n\n" +
				fmt.Sprintf(
					"  This gives you a safety net while your essential expenses currently meet or exceed your income.\n"+
						"  You have approximately £%.2f in savings above your recommended emergency fund.",
					float64(cf.CurrentSavings-emergencyFund.TargetAmount)/100,
				) +
				"\n\n" +
				"  We recommend looking at ways to reduce your essential expenses or increase your income,\n" +
				"  so that you are no longer relying on your savings to cover your monthly needs.",
		)
		fmt.Println()
		fmt.Println()
	}

	// 6. EF > CurrentSavings && Debt == 0
	if emergencyFund.TargetAmount > cf.CurrentSavings &&
		cf.UnsettledDebt == 0 {

		fmt.Println(
			"  You currently have no outstanding debt,\n" +
				"  but your savings do not yet cover your recommended emergency fund.\n\n" +
				fmt.Sprintf(
					"  Because your essential expenses currently meet or exceed your income,\n"+
						"  you may need to use your savings to cover the gap.\n"+
						"  You currently have £%.2f in savings against a recommended emergency fund of £%.2f.",
					float64(cf.CurrentSavings)/100,
					float64(emergencyFund.TargetAmount)/100,
				) +
				"\n\n" +
				"  If this continues, your savings could gradually be used up.\n\n" +
				"  We recommend looking at ways to reduce your essential expenses or increase your income,\n" +
				"  and seeking support if you need help identifying what could change in your current situation.",
		)
		fmt.Println()
		fmt.Println()
	}
}

func AssessNeedsPosition(cf CurrentFinances) {

	needsPercentage := (float64(cf.Needs) * 100) / float64(cf.Income)

	if needsPercentage < 50 {
		fmt.Println(
			"  Great news!\n" +
				"  Your must-have expenses are below 50% of your income.\n" +
				"  This gives you a healthy amount of room to put money towards your goals.\n" +
				"  Try to keep your must-have expenses at 50% or below as your circumstances change.",
		)
	} else if needsPercentage == 50 {
		fmt.Println(
			"  Great news!\n" +
				"  Your must-have expenses are 50% of your income.\n" +
				"  This gives you a healthy amount of room to put money towards your goals.\n" +
				"  Try to keep your must-have expenses at 50% or below as your circumstances change.",
		)
	} else if needsPercentage < 60 {
		fmt.Println(
			"  Good news!\n" +
				"  Your must-have expenses are below 60% of your income.\n" +
				"  We recommend aiming for 50% or below, but you're still in a reasonable starting position.\n" +
				"  As your situation changes, look for opportunities to bring this closer to 50%.",
		)
	} else if needsPercentage == 60 {
		fmt.Println(
			"  Good news!\n" +
				"  Your must-have expenses are 60% of your income.\n" +
				"  We recommend aiming for 50% or below, but you're still in a reasonable starting position.\n" +
				"  As your situation changes, look for opportunities to bring this closer to 50%.",
		)
	} else {
		fmt.Println(
			"  Your must-have expenses currently use more than 60% of your income.\n" +
				"  That's okay, everyone's situation is different,\n" +
				"  and this is simply where you're starting from.\n" +
				"  We usually recommend keeping must-have expenses at 50% or below,\n" +
				"  but we'll work with what you have available and help you find a plan that feels\n" +
				"  realistic for you.",
		)
	}

	PrintCurrentFinances(cf)
}
