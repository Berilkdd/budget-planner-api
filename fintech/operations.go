package fintech

import (
	"errors"
	"fmt"
)

// ApplyImmediateDebtPayoff uses savings above the baseline buffer to repay active debt.
func (cf *CurrentFinances) ApplyImmediateDebtPayoff(
	baselineBuffer BaselineBuffer,
) error {

	if cf.UnsettledDebt < 0 {
		return ErrNegativeUnsettledDebt
	}

	if baselineBuffer.TargetAmount <= 0 {
		return ErrInvalidBaselineBuffer
	}

	if cf.CurrentSavings <= baselineBuffer.TargetAmount {
		return ErrInsufficientSavingsForDebtPayoff
	}

	availableCash := cf.CurrentSavings - baselineBuffer.TargetAmount

	if availableCash >= cf.UnsettledDebt {

		cf.CurrentSavings -= cf.UnsettledDebt
		cf.UnsettledDebt = 0

		fmt.Printf(
			"Debt cleared: £%.2f | Remaining debt: £%.2f | Remaining savings: £%.2f\n",
			float64(availableCash)/100,
			float64(cf.UnsettledDebt)/100,
			float64(cf.CurrentSavings)/100,
		)

		PrintCurrentFinances(*cf)

		return nil
	}

	cf.UnsettledDebt -= availableCash
	cf.CurrentSavings = baselineBuffer.TargetAmount

	fmt.Printf(
		"Partial debt payment: £%.2f | Remaining debt: £%.2f | Remaining savings: £%.2f\n",
		float64(availableCash)/100,
		float64(cf.UnsettledDebt)/100,
		float64(cf.CurrentSavings)/100,
	)

	PrintCurrentFinances(*cf)

	return nil
}

func ApplyDebtFreedomPlan(
	cf *CurrentFinances,
	plan DebtFreedomPlan,
) error {

	if !plan.Available {
		return errors.New("debt freedom plan is unavailable")
	}

	totalMonths :=
		plan.BufferForecast.Phase1Months +
			plan.DebtForecast.Phase2Months

	cf.CurrentSavings =
		plan.BaselineBuffer.TargetAmount +
			plan.DebtForecast.Phase2Surplus

	cf.AvailableSurplus = 0
	cf.UnsettledDebt = 0
	cf.HasDebt = false

	cf.CurrentDate = cf.CurrentDate.AddDate(
		0,
		int(totalMonths),
		0,
	)

	PrintCurrentFinances(*cf)

	return nil
}

func ApplyEmergencyFundPlan(
	cf *CurrentFinances,
	plan EmergencyFundPlan,
) error {

	if !plan.Available {
		return errors.New("emergency fund plan is unavailable")
	}

	cf.CurrentSavings = plan.TargetAmount
	cf.AvailableSurplus = plan.Forecast.Phase3Surplus

	cf.CurrentDate = cf.CurrentDate.AddDate(
		0,
		int(plan.Forecast.Phase3Months),
		0,
	)

	PrintCurrentFinances(*cf)

	return nil
}
