package fintech

import "fmt"

// BaselineBuffer is one month of needs baseline safety net while user paying the debt off.
type BaselineBuffer struct {
	TargetAmount int64
}

func CalculateBaselineBuffer(cf CurrentFinances) (BaselineBuffer, error) {

	if cf.Needs <= 0 {
		return BaselineBuffer{}, ErrZeroNeeds
	}

	return BaselineBuffer{
		TargetAmount: cf.Needs,
	}, nil
}

// CalculateEmergencyTarget determines target cash reserves based on unemployment risk.
type EmergencyFund struct {
	TargetAmount int64
	MonthsCount  int64
}

func CalculateEmergencyTarget(
	status EmploymentStatus,
	monthlyNeeds int64,
) (EmergencyFund, error) {

	if monthlyNeeds <= 0 {
		return EmergencyFund{}, ErrZeroNeeds
	}

	fmt.Println(
		"  Calculating your emergency fund target based on your employment status...",
	)

	var emergencyFund EmergencyFund

	switch status {
	case Employee:
		emergencyFund = EmergencyFund{
			TargetAmount: monthlyNeeds * 3,
			MonthsCount:  3,
		}

	case SelfEmployed:
		emergencyFund = EmergencyFund{
			TargetAmount: monthlyNeeds * 6,
			MonthsCount:  6,
		}

	default:
		return EmergencyFund{}, ErrInvalidStatus
	}

	fmt.Println()
	fmt.Printf(
		"  Emergency fund target calculated: £%.2f (%d months of essential expenses)\n",
		float64(emergencyFund.TargetAmount)/100,
		emergencyFund.MonthsCount,
	)
	fmt.Println()
	fmt.Println()

	return emergencyFund, nil
}
