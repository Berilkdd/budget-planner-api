package fintech

import 	"errors"	

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





