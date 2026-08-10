package fintech

import (
	"errors"
)

var ErrZeroSavingAllocation = errors.New("monthly savings allocation must be greater than zero to forecast")

// BufferForecast holds the timing result (months) and the leftover money from the final month.
type BufferForecast struct {
	Phase1Months  int64 
	Phase1Surplus int64 // Leftover cash from the final month 
}

// CalculateBufferTimeline simulates building the core £1,000 baseline cushion with compounding interest.
func CalculateBufferTimeline(cf CurrentFinances, monthlySave int64) (BufferForecast, error) {
	if monthlySave <= 0 {
		return BufferForecast{}, ErrZeroSavingAllocation
	}

	// 1. If user has no active debt, bypass Phase 1 immediately to start collecting emergency fund.
	if !cf.HasDebt || cf.UnsettledDebt <= 0 {
		return BufferForecast{
			Phase1Months:  0,
			Phase1Surplus: 0,
		}, nil
	}

	// 2. If user already hits or exceeds the buffer baseline, Phase 1 takes 0 months
	if cf.CurrentSavings >= BaselineBuffer {
		return BufferForecast{
			Phase1Months:  0,
			Phase1Surplus: 0,
		}, nil
	}

	// Base rate 2.75% AER (represented as 275 basis points)
	const baseAER int64 = 275

	months := int64(0)
	runningSavings := cf.CurrentSavings

	// 3. Step forward month by month until the baseline cushion is fully collected
	for runningSavings < BaselineBuffer {
		months++

		// Calculate passive compounding monthly interest: (Balance * AER) / 10000 / 12 months
		monthlyInterest := (runningSavings * baseAER) / 10000 / 12
		runningSavings += monthlyInterest

		// Check if adding the full monthly deposit exceeds our baseline cushion target
		if runningSavings+monthlySave >= BaselineBuffer {
			// Calculate exactly how much was actually needed to hit exactly £1000
			neededToFill := BaselineBuffer - runningSavings
			surplus := monthlySave - neededToFill

			return BufferForecast{
				Phase1Months:  months,
				Phase1Surplus: surplus, 
			}, nil
		}

		// Otherwise, use full monthly save allocation and keep looping
		runningSavings += monthlySave
	}

	return BufferForecast{
		Phase1Months:  months,
		Phase1Surplus:0,
	}, nil
}

// DebtForecast holds the final timeline results and leftover cash after clearing liabilities.
type DebtForecast struct {
	TotalMonths         int64 
	Phase2Surplus       int64 // Leftover cash from the final payoff month 
	Phase2Months        int64 // Months spent inside Phase 2 active payoff loop
}

// CalculateDebtTimeline simulates dynamic, compounding debt payoff 
func CalculateDebtTimeline(cf CurrentFinances, monthlySave, phase1Months, initialSurplus int64) (DebtForecast, error) {
	// If user has no active debt, bypass Phase 2 immediately to start collecting emergency fund.
	if cf.UnsettledDebt <= 0 {
		return DebtForecast{
			TotalMonths:   phase1Months, // Total time matches whatever Phase 1 took
			Phase2Surplus: initialSurplus, // Passes the full cash straight to Phase 3
			Phase2Months:  0,
		}, nil
	}

	// 2. Resolve Interest Rate: Use user input or apply the 24% AER national average fallback 
	interestAER := cf.DebtInterestRate
	if interestAER == 0 {
		interestAER = 2400
	}

	runningDebt := cf.UnsettledDebt
	phase2Months := int64(0)

	// 3. The Phase 1 Time: Compounding debt growth while user was building the buffer
	for i := int64(0); i < phase1Months; i++ {
		monthlyInterest := (runningDebt * interestAER) / 10000 / 12
		runningDebt += monthlyInterest
	}

	// 4. Subtract the Phase 1 leftover surplus immediately before the main loop starts
	if initialSurplus >= runningDebt {
		// If surplus cash from Phase 1 completely wipes out the entire debt 
		leftoverSurplus := initialSurplus - runningDebt
		return DebtForecast{
			TotalMonths:   phase1Months, // Stacks Phase 2 onto Phase 1 timeline
			Phase2Surplus: leftoverSurplus, // Leftover cash to pass to Phase 3
			Phase2Months:  0,
		}, nil
	}
	runningDebt -= initialSurplus

	// 5. Active Payoff Loop: Month by month compounding and reduction
	for runningDebt > 0 {
		phase2Months++

		//Debt grows from monthly interest
		monthlyInterest := (runningDebt * interestAER) / 10000 / 12
		runningDebt += monthlyInterest

		// Check if our monthly savings payment completely pays off remaining debt balance
		if monthlySave >= runningDebt {
			leftoverSurplus := monthlySave - runningDebt
			return DebtForecast{
				TotalMonths: phase1Months + phase2Months, // Stacks Phase 2 onto Phase 1 timeline
				Phase2Surplus: leftoverSurplus, // Leftover cash to pass to Phase 3
				Phase2Months: phase2Months,
			}, nil
		}

		// Otherwise, make the full monthly payment and continue looping
		runningDebt -= monthlySave
	}

	return DebtForecast{
		TotalMonths:   phase1Months + phase2Months,
		Phase2Surplus: 0,
		Phase2Months:  phase2Months,
	}, nil
}
