package fintech

import "time"

CurrentDate: time.Now(),

type EmploymentStatus string

const (
	Employee     EmploymentStatus = "EMPLOYEE"
	SelfEmployed EmploymentStatus = "SELF_EMPLOYED"
)

type CurrentFinances struct {
	Income             int64
	Needs              int64
	CurrentSavings     int64
	HasDebt            bool
	UnsettledDebt      int64
	DebtInterestRate   int64
	CustomContribution int64
	EmploymentStatus   EmploymentStatus
	CurrentDate        time.Now()
	AvailableSurplus   int64
}

type UserDecisions struct {
	DebtFreedomStrategy   string
	EmergencyFundStrategy string
}
