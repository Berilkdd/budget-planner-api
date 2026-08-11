package fintech

import (
	"errors"
	"fmt"
)

var ErrZeroSavingAllocation = errors.New("monthly savings allocation must be greater than zero to forecast")
var ErrSubscriptionFeeDeficit = errors.New("monthly savings allocation is too low to sustain this tier's subscription fee")

// BufferForecast holds the timing result (months) and the leftover money from the final month.
type BufferForecast struct {
	Phase1Months  int64 
	Phase1Surplus int64 // Leftover cash from the final month 
}

// CalculateBufferTimeline simulates building the core £1,000 baseline cushion with compounding interest.
func CalculateBufferTimeline(cf CurrentFinances, monthlySave, baselineBuffer int64) (BufferForecast, error) {
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
	if cf.CurrentSavings >= baselineBuffer {
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
	for runningSavings < baselineBuffer {
		months++

		// Calculate passive compounding monthly interest: (Balance * AER) / 10000 / 12 months
		monthlyInterest := (runningSavings * baseAER) / 10000 / 12
		runningSavings += monthlyInterest

		// Check if adding the full monthly deposit exceeds our baseline cushion target
		if runningSavings+monthlySave >= baselineBuffer {
			// Calculate exactly how much was actually needed to hit exactly £1000
			neededToFill := baselineBuffer - runningSavings
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
	if monthlySave <= 0 {
			return DebtForecast{}, ErrZeroSavingAllocation
		}
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

// EmergencyFundOption isolates the forecasting metrics for a single instant access tier.
type EmergencyFundForecast struct {
	TotalMonths     int64 // Phase 1 + 2 + 3
	Phase3Months    int64 // Months spent inside Phase 3 active saving loop
	TotalFeesPaid   int64 // Accumulated subscription fees over the Phase 3 loop 
	TotalInterest   int64 // Passive compounding interest earned over the Phase 3 loop 
	Phase3Surplus   int64 // Leftover cash from the final target month 
	RecommendedTier string //The final tier the user lands on at the target milestone
}

// SavingsTierComparison collects the results for 4 monzo instant access options.
type SavingsTierComparison struct {
	TrueStartingBalance int64 // The cash in CurrentSavings at the exact moment Phase 3 opened
	Forecast            EmergencyFundForecast
}

// SimulateEmergencyFundTiers evaluates all instant access plans to find the fastest path to the user's target.
func SimulateEmergencyFundTiers(cf CurrentFinances, monthlySave, phase1Months, phase2Months, phase2Surplus, targetAmount, baselineBuffer int64) (SavingsTierComparison, error) {
	// 1. Calculate the True Starting Balance for Phase 3
	var baselineCache int64
	
	if !cf.HasDebt {
		// Bypassed Phase 1 & 2: Start with original savings balance
		baselineCache = cf.CurrentSavings
	} else {
		// Went through Phase 1: User secured the baseline buffer target
		baselineCache = baselineBuffer
	}

	// Calculate interest growth on that buffer while user spent time on Phase 2 for paying off debt
	const freeAER int64 = 275
	for i := int64(0); i < phase2Months; i++ {
		monthlyInterest := (baselineCache * freeAER) / 10000 / 12
		baselineCache += monthlyInterest
	}

	// Add any leftover cash surplus from Phase 2
	trueStartingBalance := baselineCache + phase2Surplus

	// 2. Used Monzo's Instant Access account parameters 
	tiers := []struct {
		name string
		fee  int64 
		aer  int64 
	}{
		{name: "Standard Tier (Instant Access)", fee: 0, aer: 275}, // 2.75% AER
		{name: "Extra Tier (Instant Access)", fee: 300, aer: 300},  // 3.00% AER
		{name: "Perks Tier (Instant Access)", fee: 700, aer: 325},  // 3.25% AER
		{name: "Max Tier (Instant Access)", fee: 1700, aer: 350},   // 3.50% AER
	}

	// 3. One single time-machine loop to simulate the timeline step-by-step
	phase3Months := int64(0)
	runningSavings := trueStartingBalance
	accumulatedFees := int64(0)
	accumulatedInterest := int64(0)
	var activeTierName string
	var previousTierName string

	for runningSavings < targetAmount {
		phase3Months++

		// Monthly Challenge 
		bestPlanIdx := -1
		var bestNetGrowth int64 = -999999999 // Initial baseline floor to be able to compare data with previous one on each loop

		// Challenge all 4 plans using the current balance to find this month's absolute best performer
		for idx, t := range tiers {
			projectedBalance := runningSavings + (monthlySave - t.fee)
			projectedInterest := (projectedBalance * t.aer) / 10000 / 12
			totalNetGrowth := monthlySave + projectedInterest - t.fee

			if totalNetGrowth > bestNetGrowth {
				bestNetGrowth = totalNetGrowth
				bestPlanIdx = idx
			}
		}

		// Prrovide winning tier parameters for messages and this month iteration
		winningTier := tiers[bestPlanIdx]
		activeTierName = winningTier.name
		

		if activeTierName != previousTierName {
			globalMonthTimeline := phase1Months + phase2Months + phase3Months
			
			var logMessage string
			if phase3Months == 1 {
				logMessage = fmt.Sprintf("Month %d: Core emergency phase opens. Start on %s for maximized starting returns.", globalMonthTimeline, activeTierName)
			} else {
				logMessage = fmt.Sprintf("Month %d: Break-even tipping point reached! Growth scales faster here. Plan switched dynamically to %s.", globalMonthTimeline, activeTierName)
			}
			
			var switchLog []string
			switchLog = append(switchLog, logMessage)
			previousTierName = activeTierName
		}		

		// --- APPLYING THE WINNING TIER'S MATH FOR EACH MONTH ---
		
		accumulatedFees += winningTier.fee
		runningSavings += (monthlySave - winningTier.fee)
		monthlyInterest := (runningSavings * winningTier.aer) / 10000 / 12
		accumulatedInterest += monthlyInterest
		runningSavings += monthlyInterest
	}

	// Calculate leftover surplus from the final target month
	var finalSurplus int64
	if runningSavings > targetAmount {
		finalSurplus = runningSavings - targetAmount
	}

	return SavingsTierComparison{
		TrueStartingBalance: trueStartingBalance,
		Forecast: EmergencyFundForecast{
			TotalMonths:     phase1Months + phase2Months + phase3Months, // Stacks Phase 3 onto the calendar clock
			Phase3Months:    phase3Months,
			TotalFeesPaid:   accumulatedFees,
			TotalInterest:   accumulatedInterest,
			Phase3Surplus:   finalSurplus,
			RecommendedTier: activeTierName,
		},
	}, nil
}