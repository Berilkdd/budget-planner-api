package fintech

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func readMoney() int64 {
	var amount float64
	fmt.Scan(&amount)
	return int64(amount * 100)
}

func readOptionalPercentage() int64 {
	reader := bufio.NewReader(os.Stdin)

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return 0
	}

	percentage, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0
	}

	return int64(percentage * 100)
}

func CollectCurrentFinances() CurrentFinances {
	var cf CurrentFinances

	// Income
	fmt.Print("  Monthly income: £")
	cf.Income = readMoney()

	// Employment status
	fmt.Println()
	fmt.Println("  Employment status:")
	fmt.Println("  1. Employee")
	fmt.Println("  2. Self-employed")

	var employmentChoice int
	fmt.Print("  Enter your choice: ")
	fmt.Scan(&employmentChoice)

	switch employmentChoice {
	case 1:
		cf.EmploymentStatus = Employee
	case 2:
		cf.EmploymentStatus = SelfEmployed
	default:
		cf.EmploymentStatus = Employee
	}

	// Essential expenses
	fmt.Println()
	fmt.Print(
		"  Monthly essential expenses only (excluding wants and lifestyle spending): £",
	)
	cf.Needs = readMoney()

	// Current savings

	for {
		fmt.Print("  Do you have any savings? (y/n): ")

		var answer string
		fmt.Scan(&answer)

		answer = strings.ToLower(answer)

		if answer == "y" {
			fmt.Print("  Current savings: £")
			cf.CurrentSavings = readMoney()
			break
		}

		if answer == "n" {
			cf.CurrentSavings = 0
			break
		}

		fmt.Println("  Invalid input. Please enter y or n.")
	}

	// Debt

	for {
		fmt.Print("  Do you have unsettled debt? (y/n): ")

		var answer string
		fmt.Scan(&answer)

		answer = strings.ToLower(answer)

		if answer == "y" {
			cf.HasDebt = true

			fmt.Print("  Unsettled debt: £")
			cf.UnsettledDebt = readMoney()

			fmt.Scanln()

			fmt.Print("  Debt interest rate (if known): %")
			cf.DebtInterestRate = readOptionalPercentage()

			fmt.Println(
				"  If unknown, we will use the default UK consumer credit rate for the forecast.",
			)

			break
		}

		if answer == "n" {
			cf.HasDebt = false
			cf.UnsettledDebt = 0
			cf.DebtInterestRate = 0
			break
		}

		fmt.Println("  Invalid input. Please enter y or n.")
	}

	// Current date is provided by the system.
	cf.CurrentDate = time.Now()

	return cf
}

// Translates the user's selected Debt Freedom plan into the corresponding backend strategy
func TranslateDebtFreedomSelection(
	selectedPlan SelectablePlan,
) DebtFreedomStrategy {

	switch selectedPlan.Name {
	case "Sustainable Plan":
		return DebtFreedomSustainable

	case "Moderate Plan":
		return DebtFreedomModerate

	case "Aggressive Plan":
		return DebtFreedomAggressive

	case "Custom Plan":
		return DebtFreedomCustom

	default:
		return ""
	}
}

// Translates the user's selected Emergency Fund plan into the corresponding backend strategy
func TranslateEmergencyFundSelection(
	selectedPlan SelectablePlan,
) EmergencyFundStrategy {

	switch selectedPlan.Name {
	case "Sustainable Plan":
		return EmergencyFundSustainable

	case "Moderate Plan":
		return EmergencyFundModerate

	case "Aggressive Plan":
		return EmergencyFundAggressive

	case "Custom Plan":
		return EmergencyFundCustom

	default:
		return ""
	}
}

// Handles plan selection, custom contribution input, and regeneration of the custom plan forecast
func SelectPlan(
	cf *CurrentFinances,
	plans []SelectablePlan,
	createCustomPlan func(int64) (SelectablePlan, bool, error),
	onCustomCreated func(SelectablePlan),
) (SelectablePlan, error) {

	customCreated := false

	handleCustomError := func(err error) bool {
		if errors.Is(err, ErrContributionExceedsAvailableIncome) {
			fmt.Println()
			fmt.Println("  The contribution you entered is higher than the amount available after your needs.")
			fmt.Println("  Please enter a lower monthly contribution.")
			fmt.Println()
			return true
		}

		if errors.Is(err, ErrZeroSavingAllocation) {
			fmt.Println()
			fmt.Println("  Please enter a contribution greater than £0.")
			fmt.Println()
			return true
		}

		return false
	}

	for {
		fmt.Println()

		// No standard plans are available.
		if len(plans) == 0 && !customCreated {
			fmt.Println("  Your current budget does not support any of the standard plans.")
			fmt.Println()
			fmt.Println("  Please enter the monthly amount you can contribute.")
			fmt.Println()

			var contribution int64

			fmt.Print("  Monthly contribution: £")
			fmt.Scan(&contribution)

			customPlan, selectable, err := createCustomPlan(contribution)
			if err != nil {
				if handleCustomError(err) {
					continue
				}

				return SelectablePlan{}, err
			}

			customCreated = true

			onCustomCreated(customPlan)

			if selectable {
				plans = append(plans, customPlan)
			}

			PrintAvailablePlans(*cf, plans)

			continue
		}

		// Custom plan selected but still no plans are available.

		if len(plans) == 0 && customCreated {
			fmt.Println()
			fmt.Println("  Your current budget does not support any of the available plans.")
			fmt.Println()
			fmt.Println("  Please enter another monthly contribution.")
			fmt.Println()

			var contribution int64

			fmt.Print("  Monthly contribution: £")
			fmt.Scan(&contribution)

			customPlan, selectable, err := createCustomPlan(contribution)
			if err != nil {
				if handleCustomError(err) {
					continue
				}

				return SelectablePlan{}, err
			}

			onCustomCreated(customPlan)

			if selectable {
				plans = append(plans, customPlan)
			}

			PrintAvailablePlans(*cf, plans)

			continue
		}

		fmt.Println("  Choose a Plan:")

		for i, plan := range plans {
			fmt.Printf("  %d. %s\n", i+1, plan.Name)
		}

		customOption := len(plans) + 1

		if !customCreated {
			fmt.Printf("  %d. Create a Custom Plan\n", customOption)
		} else {
			fmt.Printf(
				"  %d. Try another custom contribution\n",
				customOption,
			)
		}

		var choice int

		fmt.Print("  Enter your choice: ")
		fmt.Scan(&choice)

		// Existing selectable plan.
		if choice >= 1 && choice <= len(plans) {
			return plans[choice-1], nil
		}

		// Create the first custom plan.
		if !customCreated && choice == customOption {

			for {
				var contribution int64

				fmt.Print("  How much can you contribute monthly? £")
				fmt.Scan(&contribution)

				customPlan, selectable, err := createCustomPlan(contribution)
				if err != nil {
					if handleCustomError(err) {
						continue
					}

					return SelectablePlan{}, err
				}

				customCreated = true

				onCustomCreated(customPlan)

				if selectable {
					plans = append(plans, customPlan)
				}

				PrintAvailablePlans(*cf, plans)

				break
			}

			continue
		}
		// Generate another custom plan.
		if customCreated && choice == customOption {

			for {
				var contribution int64

				fmt.Print("  How much can you contribute monthly? £")
				fmt.Scan(&contribution)

				customPlan, selectable, err := createCustomPlan(contribution)
				if err != nil {
					if handleCustomError(err) {
						continue
					}

					return SelectablePlan{}, err
				}

				onCustomCreated(customPlan)

				if selectable {

					customIndex := -1

					for i, plan := range plans {
						if plan.Name == "Custom Plan" {
							customIndex = i
							break
						}
					}

					if customIndex >= 0 {
						// Replace the existing Custom Plan.
						plans[customIndex] = customPlan
					} else {
						// No Custom Plan exists yet.
						plans = append(plans, customPlan)
					}
				}

				PrintAvailablePlans(*cf, plans)

				break
			}

			continue
		}

		return SelectablePlan{}, errors.New(
			"invalid plan selection",
		)

	}
}
