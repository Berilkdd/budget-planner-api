package fintech

import "time"

type EmploymentStatus string

const (
	Employee     EmploymentStatus = "Employee"
	SelfEmployed EmploymentStatus = "Self Employed"
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
	CurrentDate        time.Time
	AvailableSurplus   int64
}

type DebtFreedomStrategy string

const (
	DebtFreedomSustainable DebtFreedomStrategy = "Sustainable"
	DebtFreedomModerate    DebtFreedomStrategy = "Moderate"
	DebtFreedomAggressive  DebtFreedomStrategy = "Aggressive"
	DebtFreedomCustom      DebtFreedomStrategy = "Custom"
)

type EmergencyFundStrategy string

const (
	EmergencyFundSustainable EmergencyFundStrategy = "Sustainable"
	EmergencyFundModerate    EmergencyFundStrategy = "Moderate"
	EmergencyFundAggressive  EmergencyFundStrategy = "Aggressive"
	EmergencyFundCustom      EmergencyFundStrategy = "Custom"
)

type UserDecisions struct {
	DebtFreedomStrategy   DebtFreedomStrategy
	EmergencyFundStrategy EmergencyFundStrategy
}
