package fintech

import (
	"fmt"
)

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
func calculateAllocation(cf CurrentFinances, strategyName string, wantsPercent int64) Allocation {
	wants := (cf.Income * wantsPercent) / 100
	save := cf.Income - cf.Needs - wants

	fmt.Printf(
		"[CALCULATED] %s: Needs £%.2f | Wants £%.2f | Contribution £%.2f\n",
		strategyName,
		float64(cf.Needs)/100,
		float64(wants)/100,
		float64(save)/100,
	)

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

	// Calculate wants for each strategy.
	sustainable := calculateAllocation(cf, "Sustainable", 30)
	moderate := calculateAllocation(cf, "Moderate", 25)
	aggressive := calculateAllocation(cf, "Aggressive", 20)

	// Create budget strategies.
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

	// Print the availability of each strategy.
	strategies := []BudgetStrategy{
		sustainableStrategy,
		moderateStrategy,
		aggressiveStrategy,
	}

	for _, strategy := range strategies {
		if strategy.Available {
			fmt.Printf("[STRATEGY AVAILABLE] %s\n", strategy.Name)
		} else {
			fmt.Printf("[STRATEGY UNAVAILABLE] %s\n", strategy.Name)
		}
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

	wants := cf.Income - cf.Needs - customContribution

	fmt.Printf(
		"[CALCULATED] Custom allocation: Needs £%.2f | Wants £%.2f | Contribution £%.2f\n",
		float64(cf.Needs)/100,
		float64(wants)/100,
		float64(customContribution)/100,
	)

	return BudgetStrategy{
		Name: "Custom Plan",
		Allocations: Allocation{
			Needs: cf.Needs,
			Wants: cf.Income - cf.Needs - customContribution,
			Save:  customContribution,
		},
	}, nil
}
