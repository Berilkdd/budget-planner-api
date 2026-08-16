package fintech

import (
	"errors"
	"fmt"
)

// BaselineBuffer is one month of needs baseline safety net while user paying the debt off.
type BaselineBuffer struct {
	TargetAmount int64
}

func CalculateBaselineBuffer(cf CurrentFinances) (BaselineBuffer, error) {

	if cf.Needs <= 0 {
		return BaselineBuffer{}, ErrZeroNeeds
	}

	fmt.Printf("[CALCULATED] Baseline buffer: £%.2f\n", float64(cf.Needs)/100)

	return BaselineBuffer{
		TargetAmount: cf.Needs,
	}, nil
}

// CalculateEmergencyTarget determines target cash reserves based on unemployment risk.
type EmergencyFund struct {
	TargetAmount int64
	MonthsCount  int64
}

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
