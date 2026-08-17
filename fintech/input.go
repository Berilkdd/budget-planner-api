package fintech

import (
	"errors"
	"fmt"
	"time"
)

func CollectCurrentFinances() CurrentFinances {
	var cf CurrentFinances

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("            BUDGET PLANNER")
	fmt.Println("========================================")
	fmt.Println()

	// Income
	fmt.Print("Monthly income: £")
	fmt.Scan(&cf.Income)

	// Employment status
	fmt.Println()
	fmt.Println("Employment status:")
	fmt.Println("1. Employee")
	fmt.Println("2. Self-employed")

	var employmentChoice int
	fmt.Print("Enter your choice: ")
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
		"Monthly essential expenses only (excluding wants and lifestyle spending): £",
	)
	fmt.Scan(&cf.Needs)

	// Current savings
	fmt.Println()
	fmt.Print("Do you have current savings? (y/n): ")

	var savingsAnswer string
	fmt.Scan(&savingsAnswer)

	if savingsAnswer == "y" || savingsAnswer == "Y" {
		fmt.Print("Current savings: £")
		fmt.Scan(&cf.CurrentSavings)
	} else {
		cf.CurrentSavings = 0
	}

	// Debt
	fmt.Println()
	fmt.Print("Do you have unsettled debt? (y/n): ")

	var debtAnswer string
	fmt.Scan(&debtAnswer)

	if debtAnswer == "y" || debtAnswer == "Y" {
		cf.HasDebt = true

		fmt.Print("Unsettled debt: £")
		fmt.Scan(&cf.UnsettledDebt)

		fmt.Print("Debt interest rate (if known): %")

		var interestRate float64
		fmt.Scan(&interestRate)

		cf.DebtInterestRate = int64(interestRate * 100)

		fmt.Println(
			"If unknown, we will use the default UK consumer credit rate for the forecast.",
		)
	} else {
		cf.HasDebt = false
		cf.UnsettledDebt = 0
		cf.DebtInterestRate = 0
	}

	// Current date is provided by the system.
	cf.CurrentDate = time.Now()

	return cf
}

func SelectDebtFreedomPlan(
	cf CurrentFinances,
	strategies DebtFreedomStrategies,
	baselineBuffer BaselineBuffer,
) (DebtFreedomPlan, error) {

	customCreated := false

	for {
		fmt.Println()
		fmt.Println("Choose a Debt Freedom Plan:")
		fmt.Println("1. Sustainable")
		fmt.Println("2. Moderate")
		fmt.Println("3. Aggressive")

		if !customCreated {
			fmt.Println("4. Create a Custom Plan")
		} else {
			fmt.Println("4. Custom")
			fmt.Println("5. Try another custom contribution")
		}

		var choice int
		fmt.Print("Enter your choice: ")
		fmt.Scan(&choice)

		switch choice {

		case 1:
			return strategies.Sustainable, nil

		case 2:
			return strategies.Moderate, nil

		case 3:
			return strategies.Aggressive, nil

		case 4:
			if !customCreated {
				var contribution int64

				fmt.Print("How much can you contribute monthly? £")
				fmt.Scan(&contribution)

				customPlan, err := GenerateCustomDebtFreedomPlan(
					cf,
					baselineBuffer,
					contribution,
				)
				if err != nil {
					return DebtFreedomPlan{}, err
				}

				strategies.Custom = customPlan
				customCreated = true

				fmt.Println()
				fmt.Println("Custom plan generated.")
				continue
			}

			return strategies.Custom, nil

		case 5:
			if customCreated {
				var contribution int64

				fmt.Print("How much can you contribute monthly? £")
				fmt.Scan(&contribution)

				customPlan, err := GenerateCustomDebtFreedomPlan(
					cf,
					baselineBuffer,
					contribution,
				)
				if err != nil {
					return DebtFreedomPlan{}, err
				}

				strategies.Custom = customPlan

				fmt.Println()
				fmt.Println("Custom plan updated.")
				continue
			}

			return DebtFreedomPlan{}, errors.New(
				"invalid debt freedom plan selection",
			)

		default:
			return DebtFreedomPlan{}, errors.New(
				"invalid debt freedom plan selection",
			)
		}
	}
}

func SelectEmergencyFundPlan(
	cf CurrentFinances,
	strategies EmergencyFundStrategies,
	targetAmount int64,
) (EmergencyFundPlan, error) {

	customCreated := false

	for {
		fmt.Println()
		fmt.Println("Choose an Emergency Fund Plan:")
		fmt.Println("1. Sustainable")
		fmt.Println("2. Moderate")
		fmt.Println("3. Aggressive")

		if !customCreated {
			fmt.Println("4. Create a Custom Plan")
		} else {
			fmt.Println("4. Custom")
			fmt.Println("5. Try another custom contribution")
		}

		var choice int
		fmt.Print("Enter your choice: ")
		fmt.Scan(&choice)

		switch choice {

		case 1:
			return strategies.Sustainable, nil

		case 2:
			return strategies.Moderate, nil

		case 3:
			return strategies.Aggressive, nil

		case 4:
			if !customCreated {
				var contribution int64

				fmt.Print("How much can you contribute monthly? £")
				fmt.Scan(&contribution)

				customPlan, err := GenerateCustomEmergencyFundPlan(
					cf,
					targetAmount,
					contribution,
				)
				if err != nil {
					return EmergencyFundPlan{}, err
				}

				strategies.Custom = customPlan
				customCreated = true

				fmt.Println()
				fmt.Println("Custom plan generated.")
				continue
			}

			return strategies.Custom, nil

		case 5:
			if customCreated {
				var contribution int64

				fmt.Print("How much can you contribute monthly? £")
				fmt.Scan(&contribution)

				customPlan, err := GenerateCustomEmergencyFundPlan(
					cf,
					targetAmount,
					contribution,
				)
				if err != nil {
					return EmergencyFundPlan{}, err
				}

				strategies.Custom = customPlan

				fmt.Println()
				fmt.Println("Custom plan updated.")
				continue
			}

			return EmergencyFundPlan{}, errors.New(
				"invalid emergency fund plan selection",
			)

		default:
			return EmergencyFundPlan{}, errors.New(
				"invalid emergency fund plan selection",
			)
		}
	}
}
