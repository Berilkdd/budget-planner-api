package fintech

import (
	"errors"
)

var ErrZeroSavingAllocation = errors.New("monthly savings allocation must be greater than zero to forecast")

// BufferForecast holds the timing result (months) and the leftover money from the final month.
type BufferForecast struct {
	MonthsToBuffer      int64 
	SurplusContribution int64 // Leftover cash from the final month 
}

// CalculateBufferTimeline simulates building the core £1,000 baseline cushion with compounding interest.
func CalculateBufferTimeline(cf CurrentFinances, monthlySave int64) (BufferForecast, error) {
	if monthlySave <= 0 {
		return BufferForecast{}, ErrZeroSavingAllocation
	}

	// 1. If user already hits or exceeds the buffer baseline, Phase 1 takes 0 months
	if cf.CurrentSavings >= BaselineBuffer {
		return BufferForecast{
			MonthsToBuffer:      0,
			SurplusContribution: 0,
		}, nil
	}

	// Base rate 2.75% AER (represented as 275 basis points)
	const baseAER int64 = 275

	months := int64(0)
	runningSavings := cf.CurrentSavings

	// 2. Step forward month by month until the baseline cushion is fully collected
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
				MonthsToBuffer:      months,
				SurplusContribution: surplus, 
			}, nil
		}

		// Otherwise, use full monthly save allocation and keep looping
		runningSavings += monthlySave
	}

	return BufferForecast{
		MonthsToBuffer:      months,
		SurplusContribution: 0,
	}, nil
}
