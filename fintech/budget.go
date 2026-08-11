package fintech

import (
	"errors"
	"fmt"
)
// This error triggers if a user's fixed costs eat up more than 60% of their income.

var ErrNeedsTooHigh = errors.New("monthly needs exceed maximum allowed 60% of income")
var ErrInvalidStatus = errors.New("invalid employment status provided")

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

// BaselineBuffer represents strict one month of needs baseline safety net.
type BaselineBuffer struct {
	TargetAmount int64
}

// Guard 1: Defend the core baseline cushion first
func CalculateBaselineBuffer(cf CurrentFinances) (BaselineBuffer, error) {

	if cf.Needs <= 0 {
		return BaselineBuffer{}, errors.New("monthly needs must be greater than zero")
	}

		return BaselineBuffer{
		TargetAmount: cf.Needs,
	}, nil
}

// ApplyImmediateDebtPayoff uses surplus savings above the baseline buffer to clear active debt.
func (cf *CurrentFinances) ApplyImmediateDebtPayoff() (string, error) {

	baselineBuffer, err := CalculateBaselineBuffer(*cf)

	if err != nil {
		return "", err
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
	availableCash := cf.CurrentSavings - baselineBuffer.TargetAmount

	if availableCash >= cf.UnsettledDebt {
		// Cash completely wipes out debt, leaving the remainder in savings
		cf.CurrentSavings -= cf.UnsettledDebt
		cf.UnsettledDebt = 0
		return fmt.Sprintf("Debt completely cleared using extra savings. New savings: £%.2f.", float64(cf.CurrentSavings)/100), nil
	} else {
		// Cash reduces debt partially, safely locking savings exactly at the BaselineBuffer
		cf.UnsettledDebt -= availableCash
		cf.CurrentSavings = baselineBuffer.TargetAmount
		return fmt.Sprintf("Extra savings applied to debt. Remaining savings locked at £%.2f. Remaining debt: £%.2f.", float64(cf.CurrentSavings)/100, float64(cf.UnsettledDebt)/100), nil
	}
}

// EmergencyFund represents the long-term cash safety cushion target for the user.
type EmergencyFund struct {
	TargetAmount int64 
	MonthsCount  int64 
}

// CalculateEmergencyTarget determines target cash reserves based on unemployment risk.
func CalculateEmergencyTarget(status EmploymentStatus, monthlyNeeds int64) (EmergencyFund, error) {
	if monthlyNeeds <= 0 {
		return EmergencyFund{}, errors.New("monthly needs must be greater than zero")
	}

	switch status {
	case Employee:
		return EmergencyFund{
			TargetAmount: monthlyNeeds * 3,
			MonthsCount:  3,
		}, nil
	case SelfEmployed:
		return EmergencyFund{
			TargetAmount: monthlyNeeds * 6,
			MonthsCount:  6,
		}, nil
	default:
		return EmergencyFund{}, ErrInvalidStatus
	}
}

