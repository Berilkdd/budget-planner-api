package fintech
import "errors"

var ErrZeroIncome = errors.New("income must be greater than zero")
var ErrZeroNeeds = errors.New("monthly needs must be greater than zero")
var ErrInvalidStatus = errors.New("invalid employment status provided")
var ErrZeroSavingAllocation = errors.New("monthly savings allocation must be greater than zero to forecast")
var ErrContributionExceedsAvailableIncome = errors.New("custom contribution cannot exceed income remaining after needs")
var ErrSubscriptionFeeDeficit = errors.New("monthly savings allocation is too low to sustain this tier's subscription fee")

// Warning codes
type WarningCode string

const (
	WarningLowWantsBudget                  WarningCode = "LOW_WANTS_BUDGET"
	WarningNeedsExceedIncome               WarningCode = "NEEDS_EXCEED_INCOME"

	WarningNeedsEqualIncome                WarningCode = "NEEDS_EQUAL_INCOME"
	WarningBelowEmergencyFund              WarningCode = "BELOW_EMERGENCY_FUND"
	WarningEmergencyFundCovered            WarningCode = "EMERGENCY_FUND_COVERED"
	WarningSavingsMayBeUsed                WarningCode = "SAVINGS_MAY_BE_USED"
	WarningUnsettledDebt                   WarningCode = "UNSETTLED_DEBT"
	WarningNoUnsettledDebt                 WarningCode = "NO_UNSETTLED_DEBT"	
)

// Action codes
type ActionCode string

const (
	ActionSupportAdvised             ActionCode = "SUPPORT_ADVISED"
	ActionDebtAdviceAdvised          ActionCode = "DEBT_ADVICE_ADVISED"
	ActionFullDebtPaymentAdvised     ActionCode = "FULL_DEBT_PAYMENT_ADVISED"
	ActionPartialDebtPaymentAdvised  ActionCode = "PARTIAL_DEBT_PAYMENT_ADVISED"
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
		Code: WarningLowWantsBudget,
		Description: "Discretionary money is below £100 per month for a debt repayment period of 6 months or longer.",
	},
	
	WarningNeedsExceedIncome: {
		Code: WarningNeedsExceedIncome,
		Description: "Essential needs exceed monthly income.",
	},

	WarningNeedsEqualIncome: {
		Code: WarningNeedsEqualIncome,
		Description: "Essential needs meet monthly income.",
	},

	WarningBelowEmergencyFund: {
		Code: WarningBelowEmergencyFund,
		Description: "Current savings are below the required emergency fund.",
	},

	WarningEmergencyFundCovered: {
		Code: WarningEmergencyFundCovered,
		Description: "Current savings cover the required emergency fund.",
	},

	WarningSavingsMayBeUsed: {
		Code: WarningSavingsMayBeUsed,
		Description: "Savings may be needed to cover the gap between essential expenses and income over time.",
	},

	WarningUnsettledDebt: {
		Code: WarningUnsettledDebt,
		Description: "Outstanding debt remains.",
	},

	WarningNoUnsettledDebt: {
		Code: WarningNoUnsettledDebt,
		Description: "No outstanding debt remains.",
	},	
}