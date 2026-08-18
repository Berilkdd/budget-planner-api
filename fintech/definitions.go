package fintech

import "errors"

var ErrZeroIncome = errors.New("income must be greater than zero")
var ErrZeroNeeds = errors.New("monthly needs must be greater than zero")
var ErrNegativeUnsettledDebt = errors.New("unsettled debt cannot be negative")
var ErrInsufficientSavingsForDebtPayoff = errors.New("savings do not exceed the required baseline buffer")
var ErrInvalidBaselineBuffer = errors.New("baseline buffer must be greater than zero")
var ErrInvalidStatus = errors.New("invalid employment status provided")
var ErrZeroSavingAllocation = errors.New("monthly savings allocation must be greater than zero to forecast")
var ErrContributionExceedsAvailableIncome = errors.New("custom contribution cannot exceed income remaining after needs")
var ErrSubscriptionFeeDeficit = errors.New("monthly savings allocation is too low to sustain this tier's subscription fee")

// Warning codes
type WarningCode string

const (
	WarningLowWantsBudget    WarningCode = "LOW_WANTS_BUDGET"
	WarningNeedsExceedIncome WarningCode = "NEEDS_EXCEED_INCOME"

	WarningNeedsEqualIncome     WarningCode = "NEEDS_EQUAL_INCOME"
	WarningBelowEmergencyFund   WarningCode = "BELOW_EMERGENCY_FUND"
	WarningEmergencyFundCovered WarningCode = "EMERGENCY_FUND_COVERED"
	WarningSavingsAtRisk        WarningCode = "SAVINGS_AT_RISK"
	WarningUnsettledDebt        WarningCode = "UNSETTLED_DEBT"
	WarningNoUnsettledDebt      WarningCode = "NO_UNSETTLED_DEBT"

	WarningNeedsBelow50        WarningCode = "NEEDS_BELOW_50"
	WarningNeedsEqual50        WarningCode = "NEEDS_EQUAL_50"
	WarningNeedsBetween50And60 WarningCode = "NEEDS_BELOW_60"
	WarningNeedsEqual60        WarningCode = "NEEDS_EQUAL_60"
	WarningNeedsAbove60        WarningCode = "NEEDS_ABOVE_60"
)

// Action codes
type ActionCode string
type Action struct {
	Code        ActionCode
	Description string
}

const (
	ActionSupportAdvised            ActionCode = "SUPPORT_ADVISED"
	ActionDebtAdviceAdvised         ActionCode = "DEBT_ADVICE_ADVISED"
	ActionFullDebtPaymentAdvised    ActionCode = "FULL_DEBT_PAYMENT_ADVISED"
	ActionPartialDebtPaymentAdvised ActionCode = "PARTIAL_DEBT_PAYMENT_ADVISED"
)

var ActionDefinitions = map[ActionCode]Action{
	ActionSupportAdvised: {
		Code:        ActionSupportAdvised,
		Description: "Check eligibility for available financial support and benefits.",
	},

	ActionDebtAdviceAdvised: {
		Code:        ActionDebtAdviceAdvised,
		Description: "Professional debt advice is advised.",
	},

	ActionFullDebtPaymentAdvised: {
		Code:        ActionFullDebtPaymentAdvised,
		Description: "Full repayment of the outstanding debt is advised while maintaining the required emergency fund.",
	},

	ActionPartialDebtPaymentAdvised: {
		Code:        ActionPartialDebtPaymentAdvised,
		Description: "Partial repayment of the outstanding debt is advised using savings above the required emergency fund.",
	},
}

type Warning struct {
	Code        WarningCode
	Description string
}

// Warning definitions
var WarningDefinitions = map[WarningCode]Warning{
	WarningLowWantsBudget: {
		Code:        WarningLowWantsBudget,
		Description: "Discretionary money is below £100 per month for a debt repayment period of 6 months or longer.",
	},

	WarningNeedsExceedIncome: {
		Code:        WarningNeedsExceedIncome,
		Description: "Essential needs exceed monthly income.",
	},

	WarningNeedsEqualIncome: {
		Code:        WarningNeedsEqualIncome,
		Description: "Essential needs meet monthly income.",
	},

	WarningBelowEmergencyFund: {
		Code:        WarningBelowEmergencyFund,
		Description: "Current savings are below the required emergency fund.",
	},

	WarningEmergencyFundCovered: {
		Code:        WarningEmergencyFundCovered,
		Description: "Current savings cover the required emergency fund.",
	},

	WarningSavingsAtRisk: {
		Code:        WarningSavingsAtRisk,
		Description: "Savings may be needed to cover the gap between essential expenses and income over time.",
	},

	WarningUnsettledDebt: {
		Code:        WarningUnsettledDebt,
		Description: "Outstanding debt remains.",
	},

	WarningNoUnsettledDebt: {
		Code:        WarningNoUnsettledDebt,
		Description: "No outstanding debt remains.",
	},

	WarningNeedsBelow50: {
		Code:        WarningNeedsBelow50,
		Description: "Needs below 50% of income. Healthy.",
	},

	WarningNeedsEqual50: {
		Code:        WarningNeedsEqual50,
		Description: "Needs equal 50% of income.Healthy.",
	},

	WarningNeedsBetween50And60: {
		Code:        WarningNeedsBetween50And60,
		Description: "Needs above 50% but below 60% of income. Reasonable.",
	},

	WarningNeedsEqual60: {
		Code:        WarningNeedsEqual60,
		Description: "Needs equal 60% of income. Reasonable.",
	},

	WarningNeedsAbove60: {
		Code:        WarningNeedsAbove60,
		Description: "Needs above 60% of income. Tight",
	},
}
