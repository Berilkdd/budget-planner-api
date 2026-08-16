package fintech

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

	// 2. Calculate the required emergency fund.
	emergencyFund, err := CalculateEmergencyTarget(
		cf.EmploymentStatus,
		cf.Needs,
	)

	if err != nil {
		return assessment
	}

	assessment.EmergencyFundTarget = emergencyFund.TargetAmount
	assessment.EmergencyFundMonths = emergencyFund.MonthsCount

	// 3. Savings below emergency-fund target.
	if cf.CurrentSavings < emergencyFund.TargetAmount {

		// No emergency fund + needs exceed income
		if !needsEqualIncome {
			assessment.Actions = append(
				assessment.Actions,
				ActionSupportAdvised,
			)
		}

		// No emergency fund + debt
		if hasDebt {

			assessment.Pathway = PathwayC

			assessment.Actions = append(
				assessment.Actions,
				ActionDebtAdviceAdvised,
			)

		} else {

			assessment.Pathway = PathwayD

		}

		return assessment
	}

	// 4. Emergency fund is covered.

	// 5. No debt.
	if !hasDebt {

		// Only Needs > Income creates the long-term savings warning.
		assessment.Pathway = PathwayB

		return assessment
	}

	// 6. Debt remains.

	extraSavings := cf.CurrentSavings - emergencyFund.TargetAmount

	// 7. Enough extra savings to clear the debt completely.
	if extraSavings >= cf.UnsettledDebt {

		assessment.Pathway = PathwayA1

		assessment.Actions = append(
			assessment.Actions,
			ActionFullDebtPaymentAdvised,
		)

		return assessment
	}

	// 8. Some extra savings, but not enough to clear the debt.
	if extraSavings > 0 {

		assessment.Pathway = PathwayA2

		assessment.Actions = append(
			assessment.Actions,
			ActionPartialDebtPaymentAdvised,
			ActionDebtAdviceAdvised,
		)

		return assessment
	}

	// 9. Emergency fund is covered exactly, but there is no extra
	// savings available for debt repayment.

	assessment.Pathway = PathwayA3

	assessment.Actions = append(
		assessment.Actions,
		ActionDebtAdviceAdvised,
	)

	return assessment
}

func AssessNeedsPosition(cf CurrentFinances) AssessmentData {

	assessment := AssessmentData{
		EmploymentStatus:              cf.EmploymentStatus,
		NeedsEqualIncome:              cf.Needs == cf.Income,
		CurrentSavings:                cf.CurrentSavings,
		MonthsOfNeedsCoveredBySavings: cf.CurrentSavings / cf.Needs,
		Actions:                       []ActionCode{},
	}

	needsPercentage := (cf.Needs * 100) / cf.Income

	if needsPercentage < 50 {
		assessment.Pathway = PathwayE1
		return assessment
	}

	if needsPercentage == 50 {
		assessment.Pathway = PathwayE2
		return assessment
	}

	if needsPercentage < 60 {
		assessment.Pathway = PathwayE3
		return assessment
	}

	if needsPercentage == 60 {
		assessment.Pathway = PathwayE4
		return assessment
	}

	assessment.Pathway = PathwayE5

	return assessment
}
