package main

import "budget-planner-api/fintech"

func main() {

	if err := fintech.RunPlanner(); err != nil {
		panic(err)
	}
}
