package fintech

import (
	"errors"
	"fmt"
)

func SelectDebtFreedomPlan(
	cf CurrentFinances,
	strategies DebtFreedomStrategies,
	baselineBuffer BaselineBuffer,
) (DebtFreedomPlan, error) {

	for {
		fmt.Println()
		fmt.Println("Choose a Debt Freedom Plan:")
		fmt.Println("1. Sustainable")
		fmt.Println("2. Moderate")
		fmt.Println("3. Aggressive")
		fmt.Println("4. Custom")

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
			var contribution int64

			fmt.Print("Enter your monthly contribution: ")
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
			fmt.Println("Custom plan generated.")
			fmt.Println()
			fmt.Println("Choose a Debt Freedom Plan:")
			fmt.Println("1. Sustainable")
			fmt.Println("2. Moderate")
			fmt.Println("3. Aggressive")
			fmt.Println("4. Custom")

			var customChoice int
			fmt.Print("Enter your choice: ")
			fmt.Scan(&customChoice)

			switch customChoice {
			case 1:
				return strategies.Sustainable, nil

			case 2:
				return strategies.Moderate, nil

			case 3:
				return strategies.Aggressive, nil

			case 4:
				return strategies.Custom, nil

			default:
				return DebtFreedomPlan{}, errors.New(
					"invalid debt freedom plan selection",
				)
			}

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

	for {
		fmt.Println()
		fmt.Println("Choose an Emergency Fund Plan:")
		fmt.Println("1. Sustainable")
		fmt.Println("2. Moderate")
		fmt.Println("3. Aggressive")
		fmt.Println("4. Custom")

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
			var contribution int64

			fmt.Print("Enter your monthly contribution: ")
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
			fmt.Println("Custom plan generated.")
			fmt.Println()
			fmt.Println("Choose an Emergency Fund Plan:")
			fmt.Println("1. Sustainable")
			fmt.Println("2. Moderate")
			fmt.Println("3. Aggressive")
			fmt.Println("4. Custom")

			var customChoice int
			fmt.Print("Enter your choice: ")
			fmt.Scan(&customChoice)

			switch customChoice {
			case 1:
				return strategies.Sustainable, nil

			case 2:
				return strategies.Moderate, nil

			case 3:
				return strategies.Aggressive, nil

			case 4:
				return strategies.Custom, nil

			default:
				return EmergencyFundPlan{}, errors.New(
					"invalid emergency fund plan selection",
				)
			}

		default:
			return EmergencyFundPlan{}, errors.New(
				"invalid emergency fund plan selection",
			)
		}
	}
}
