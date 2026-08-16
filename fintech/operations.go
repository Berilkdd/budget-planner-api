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
	return nil
}
