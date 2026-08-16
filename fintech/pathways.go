package fintech

type PathwayCode string

const (
	PathwayA1 PathwayCode = "A1" // Emergency fund covered + enough savings to clear debt
	PathwayA2 PathwayCode = "A2" // Emergency fund covered + partial debt repayment
	PathwayA3 PathwayCode = "A3" // Emergency fund covered + debt + no extra savings
	PathwayB  PathwayCode = "B"  // Emergency fund covered + no debt
	PathwayC  PathwayCode = "C"  // Emergency fund below target + debt
	PathwayD  PathwayCode = "D"  // Emergency fund below target + no debt

	PathwayE1 PathwayCode = "E1" // Needs below 50% of income
	PathwayE2 PathwayCode = "E2" // Needs equal 50% of income
	PathwayE3 PathwayCode = "E3" // Needs above 50% but below 60%
	PathwayE4 PathwayCode = "E4" // Needs equal 60% of income
	PathwayE5 PathwayCode = "E5" // Needs above 60% of income
)
