package main

import (
	"fmt"

	"budget-planner-api/fintech"
)

func main() {

	for {

		if err := fintech.RunPlanner(); err != nil {
			panic(err)
		}

		fmt.Println()
		fmt.Println("  =================================================================================")
		fmt.Println("                              FORECAST COMPLETE")
		fmt.Println("  =================================================================================")
		fmt.Println()
		fmt.Println("  1. Run another forecast")
		fmt.Println("  2. Exit")
		fmt.Println()
		fmt.Print("  Enter your choice: ")

		for {
			var choice int

			fmt.Scan(&choice)

			if choice == 1 {
				break
			}

			if choice == 2 {
				return
			}

			fmt.Println()
			fmt.Println("  Invalid choice. Please enter 1 or 2.")
			fmt.Print("  Enter your choice: ")
		}
	}
}
