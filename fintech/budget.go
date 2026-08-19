package fintech

type Allocation struct {
	Needs int64
	Wants int64
	Save  int64
}

type BudgetStrategy struct {
	Name        string
	Allocations Allocation
	Available   bool
}

type AllocationOptions struct {
	Sustainable BudgetStrategy
	Moderate    BudgetStrategy
	Aggressive  BudgetStrategy
	Custom      BudgetStrategy
}

// Calculates the allocation for a default strategy.
func calculateAllocation(
	cf CurrentFinances,
	wantsPercent int64,
) Allocation {

	wants := (cf.Income * wantsPercent) / 100
	save := cf.Income - cf.Needs - wants

	return Allocation{
		Needs: cf.Needs,
		Wants: wants,
		Save:  save,
	}
}

// Creates the default budget strategies
func GenerateAllocation(cf CurrentFinances) (AllocationOptions, error) {

	if cf.Income <= 0 {
		return AllocationOptions{}, ErrZeroIncome
	}

	sustainable := calculateAllocation(cf, 30)
	moderate := calculateAllocation(cf, 25)
	aggressive := calculateAllocation(cf, 20)

	sustainableStrategy := BudgetStrategy{
		Name:        "Sustainable Plan",
		Allocations: sustainable,
		Available:   sustainable.Save > 0,
	}

	moderateStrategy := BudgetStrategy{
		Name:        "Moderate Plan",
		Allocations: moderate,
		Available:   moderate.Save > 0,
	}

	aggressiveStrategy := BudgetStrategy{
		Name:        "Aggressive Plan",
		Allocations: aggressive,
		Available:   aggressive.Save > 0,
	}

	return AllocationOptions{
		Sustainable: sustainableStrategy,
		Moderate:    moderateStrategy,
		Aggressive:  aggressiveStrategy,
	}, nil
}

// Creates a custom budget plan from the user's chosen monthly contribution.
func GenerateCustomAllocation(cf CurrentFinances, customContribution int64) (BudgetStrategy, error) {
	if customContribution <= 0 {
		return BudgetStrategy{}, ErrZeroSavingAllocation
	}

	availableAfterNeeds := cf.Income - cf.Needs

	if customContribution > availableAfterNeeds {
		return BudgetStrategy{}, ErrContributionExceedsAvailableIncome
	}

	return BudgetStrategy{
		Name:      "Custom Plan",
		Available: true,
		Allocations: Allocation{
			Needs: cf.Needs,
			Wants: cf.Income - cf.Needs - customContribution,
			Save:  customContribution,
		},
	}, nil
}
