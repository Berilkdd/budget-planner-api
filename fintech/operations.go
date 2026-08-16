package fintech

import (
	"errors"
	"fmt"
)

// ApplyImmediateDebtPayoff uses surplus savings above the baseline buffer to clear active debt.
func (cf *CurrentFinances) ApplyImmediateDebtPayoff(baselineBuffer BaselineBuffer) error {

	// Guard: Error for negative debt amount
	if cf.UnsettledDebt < 0 {
		err := errors.New("unsettled debt cannot be negative.")
		fmt.Printf("[ERROR] %s\n", err)
		return err
	}

	// Verify a real liability exists to pay off
	if cf.UnsettledDebt == 0 {
		fmt.Println("[STATUS] No active liabilities found. Directing to Emergency Fund goals.")
		return nil
	}

	if cf.CurrentSavings <= baselineBuffer.TargetAmount {
		fmt.Println("[STATUS] Balance is below the recommended safety cushion.")
		return nil
	}

	// Calculate exactly how much extra cash we are allowed to use
	availableCash := cf.CurrentSavings - baselineBuffer.TargetAmount

	if availableCash >= cf.UnsettledDebt {
		// Cash completely wipes out debt, leaving the remainder in savings
		cf.CurrentSavings -= cf.UnsettledDebt
		cf.UnsettledDebt = 0

		fmt.Printf(
			"[OPERATION] Debt cleared: £%.2f | Remaining debt: £%.2f | Remaining savings: £%.2f\n",
			float64(availableCash)/100,
			float64(cf.UnsettledDebt)/100,
			float64(cf.CurrentSavings)/100,
		)

		return nil
	}

	// Cash reduces debt partially, safely locking savings exactly at the BaselineBuffer
	cf.UnsettledDebt -= availableCash
	cf.CurrentSavings = baselineBuffer.TargetAmount

	fmt.Printf(
		"[OPERATION] Partial debt payment: £%.2f | Remaining debt: £%.2f | Remaining savings: £%.2f\n",
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
