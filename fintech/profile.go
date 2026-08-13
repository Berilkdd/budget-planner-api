package fintech

type EmploymentStatus string

const (
	Employee     EmploymentStatus = "EMPLOYEE"
	SelfEmployed EmploymentStatus = "SELF_EMPLOYED"
)

type CurrentFinances struct {
	Income           int64
	Needs            int64
	CurrentSavings   int64
	HasDebt          bool
	UnsettledDebt    int64
	DebtInterestRate int64
	CustomContribution int64
	EmploymentStatus EmploymentStatus
}

type UserDecisions struct {
	DebtFreedomStrategy string
	EmergencyFundStrategy string
}