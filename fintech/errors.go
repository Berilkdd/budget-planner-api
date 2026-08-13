package fintech
import "errors"

var ErrZeroIncome = errors.New("income must be greater than zero")
var ErrInvalidStatus = errors.New("invalid employment status provided")
var ErrZeroSavingAllocation = errors.New("monthly savings allocation must be greater than zero to forecast")
var ErrContributionExceedsAvailableIncome = errors.New("custom contribution cannot exceed income remaining after needs")
var ErrSubscriptionFeeDeficit = errors.New("monthly savings allocation is too low to sustain this tier's subscription fee")

// Warning codes
type WarningCode string

const (
	WarningLowDiscretionaryBudget WarningCode = "LOW_DISCRETIONARY_BUDGET"
)

type Warning struct {
	Code        WarningCode
	Description string
}

// Warning definitions
var WarningDefinitions = map[WarningCode]Warning{
	WarningLowDiscretionaryBudget: {
		Code:        WarningLowDiscretionaryBudget,
		Description: "Discretionary money is below £100 per month for a debt repayment period of 6 months or longer.",
	},
}
