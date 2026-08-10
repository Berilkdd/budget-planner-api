package fintech

import 	(
	"errors"	
	"fmt"	
)
// This error triggers if a user's fixed costs eat up more than 60% of their income.

var ErrNeedsTooHigh = errors.New("monthly needs exceed maximum allowed 60% of income")

type Allocation struct {
	Needs int64
	Wants int64
	Save  int64
}

type BudgetStrategy struct {
	Name        string
	Allocations Allocation
}

// Creates custom plan for user based on their income and needs.
func GenerateAllocation(cf CurrentFinances) (BudgetStrategy, error) {

	if cf.Income <= 0 {
		return BudgetStrategy{}, errors.New("income must be greater than zero")
	}

	// Calculate wants at a flat 30% using precise integer math
	wantsAlloc := (cf.Income * 30) / 100

	// Calculate leftover balance dynamically for savings
	saveAlloc := cf.Income - cf.Needs - wantsAlloc

	// Situation 1: Equal or below 50%
	if cf.Needs*100 <= cf.Income*50 {
		return BudgetStrategy{
			Name: "Ideal Framework (Needs <= 50%)",
			Allocations: Allocation{
				Needs: cf.Needs,
				Wants: wantsAlloc,
				Save:  saveAlloc,
			},
		}, nil
	}

	// Situation 2: Greater than 50% but less than or equal to 60%
	if cf.Needs*100 <= cf.Income*60 {
		return BudgetStrategy{
			Name: "Not Ideal But Balanced (50% < Needs <= 60%)",
			Allocations: Allocation{
				Needs: cf.Needs,
				Wants: wantsAlloc,
				Save:  saveAlloc,
			},
		}, nil
	}

	// Situation 3: More than 60%
	return BudgetStrategy{}, ErrNeedsTooHigh
}


// BaselineBuffer represents strict £1,000.00 baseline safety net.
const BaselineBuffer int64 = 100000

// ApplyImmediateDebtPayoff uses surplus savings above the baseline buffer to clear active debt.
func (cf *CurrentFinances) ApplyImmediateDebtPayoff() (string, error) {

	// Guard 1: Defend the core baseline cushion first
	if cf.CurrentSavings <= BaselineBuffer {
		return "No action. Balance is below the £1000.00 safety cushion.", nil
	}

	// Guard 2: Error for negative debt amount
	if cf.UnsettledDebt < 0 {
    return "", errors.New("unsettled debt cannot be negative.")
	}

	// Guard 3: Verify a real liability exists to pay off
	if cf.UnsettledDebt == 0 {
		return "No action. No active liabilities found. Directing to Emergency Fund goals.", nil
	}

	// Calculate exactly how much extra cash we are allowed to use
	availableCash := cf.CurrentSavings - BaselineBuffer

	if availableCash >= cf.UnsettledDebt {
		// Cash completely wipes out debt, leaving the remainder in savings
		cf.CurrentSavings -= cf.UnsettledDebt
		cf.UnsettledDebt = 0
		return fmt.Sprintf("Debt completely cleared using extra savings. New savings: £%.2f.", float64(cf.CurrentSavings)/100), nil
	} else {
		// Cash reduces debt partially, safely locking savings exactly at the BaselineBuffer
		cf.UnsettledDebt -= availableCash
		cf.CurrentSavings = BaselineBuffer
		return fmt.Sprintf("Extra savings applied to debt. Remaining savings locked at £%.2f. Remaining debt: £%.2f.", float64(cf.CurrentSavings)/100, float64(cf.UnsettledDebt)/100), nil
	}
}



