package fintech

import "errors"

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

func GenerateAllocation(income, needs int64) (BudgetStrategy, error) {
	if income <= 0 {
		return BudgetStrategy{}, errors.New("income must be greater than zero")
	}

	// Rule: ratio <= 0.50. Checked cleanly using integer cross-multiplication: (needs * 100 <= income * 50)
	if needs*100 <= income*50 {
		return BudgetStrategy{
			Name: "50/30/20 Strategy (Ideal)",
			Allocations: Allocation{
				Needs: (income * 50) / 100,
				Wants: (income * 30) / 100,
				Save:  (income * 20) / 100,
			},
		}, nil
	} else if needs*100 <= income*60 {
		return BudgetStrategy{
			Name: "60/30/10 Strategy (Adjusted)",
			Allocations: Allocation{
				Needs: (income * 60) / 100,
				Wants: (income * 30) / 100,
				Save:  (income * 10) / 100,
			},
		}, nil
	}

	return BudgetStrategy{}, ErrNeedsTooHigh
}
