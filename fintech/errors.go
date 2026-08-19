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
