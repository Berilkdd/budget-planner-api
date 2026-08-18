package fintech

// CalculateBufferGrowthDuringDebt simulates the protected buffer
// earning interest throughout Phase 2 while debt is being paid off.
type BufferGrowthForecast struct {
	FinalBuffer        int64
	Phase2InterestGain int64
	Phase2Fees         int64
	TierBreakpoints    []TierBreakpoint
}

func CalculateBufferGrowthDuringDebt(
	initialBuffer int64,
	phase2Months int64,
) BufferGrowthForecast {

	runningBuffer := initialBuffer
	interestGain := int64(0)
	fees := int64(0)

	var tierBreakpoints []TierBreakpoint
	var previousTier *InstantAccessTier

	for i := int64(0); i < phase2Months; i++ {

		// Calculate the best instant access tier for the protected buffer.
		winningTier := CalculateBestInstantAccessTier(runningBuffer)

		recordTierBreakpoint(
			&tierBreakpoints,
			i,
			runningBuffer,
			previousTier,
			winningTier,
		)

		previousTier = &winningTier

		// Apply the tier fee.
		runningBuffer -= winningTier.Fee
		fees += winningTier.Fee

		// Calculate and record interest gained by the protected buffer.
		monthlyInterest := (runningBuffer * winningTier.AER) / 10000 / 12
		runningBuffer += monthlyInterest
		interestGain += monthlyInterest
	}

	return BufferGrowthForecast{
		FinalBuffer:        runningBuffer,
		Phase2InterestGain: interestGain,
		Phase2Fees:         fees,
		TierBreakpoints:    tierBreakpoints,
	}
}

// CalculateDebtGrowthDuringBuffer simulates debt interest accruing
// while Phase 1 is building the protected buffer.
type DebtGrowthForecast struct {
	FinalDebt             int64
	Phase1InterestLost    int64
	InterestOver50Percent bool
}

func CalculateDebtGrowthDuringBuffer(
	cf CurrentFinances,
	monthlySave,
	phase1Months int64,
) DebtGrowthForecast {

	if cf.UnsettledDebt <= 0 || phase1Months <= 0 {
		return DebtGrowthForecast{
			FinalDebt:             cf.UnsettledDebt,
			Phase1InterestLost:    0,
			InterestOver50Percent: false,
		}
	}

	// Use the user's debt interest rate or the 24% AER fallback.
	interestAER := cf.DebtInterestRate
	if interestAER == 0 {
		interestAER = 2400
	}

	runningDebt := cf.UnsettledDebt
	interestLost := int64(0)
	interestOver50Percent := false

	for i := int64(0); i < phase1Months; i++ {

		// Calculate and apply monthly debt interest.
		monthlyInterest := (runningDebt * interestAER) / 10000 / 12

		runningDebt += monthlyInterest
		interestLost += monthlyInterest

		// Check whether debt interest consumed more than
		// 50% of the monthly saving contribution.
		if monthlyInterest > monthlySave/2 {
			interestOver50Percent = true
		}
	}

	return DebtGrowthForecast{
		FinalDebt:             runningDebt,
		Phase1InterestLost:    interestLost,
		InterestOver50Percent: interestOver50Percent,
	}
}

// CalculateBufferTimeline simulates building the core baseline buffer with compounding interest.
type BufferForecast struct {
	Phase1Months       int64
	Phase1Surplus      int64
	Phase1InterestGain int64
	Phase1Fees         int64
	TierBreakpoints    []TierBreakpoint
}

func CalculateBufferTimeline(cf CurrentFinances, monthlySave, baselineBuffer int64) (BufferForecast, error) {
	if monthlySave <= 0 {
		return BufferForecast{}, ErrZeroSavingAllocation
	}

	// If there is no active debt, Phase 1 is bypassed.
	if !cf.HasDebt || cf.UnsettledDebt <= 0 {
		return BufferForecast{
			Phase1Months:       0,
			Phase1Surplus:      0,
			Phase1InterestGain: 0,
			Phase1Fees:         0,
		}, nil
	}

	// If the buffer is already fully funded, Phase 1 takes 0 months.
	if cf.CurrentSavings >= baselineBuffer {
		return BufferForecast{
			Phase1Months:       0,
			Phase1Surplus:      0,
			Phase1InterestGain: 0,
			Phase1Fees:         0,
		}, nil
	}

	months := int64(0)
	runningSavings := cf.CurrentSavings
	interestGain := int64(0)
	fees := int64(0)

	var tierBreakpoints []TierBreakpoint
	var previousTier *InstantAccessTier

	// Step forward month by month until the baseline buffer is reached.
	for runningSavings < baselineBuffer {

		// Calculate the best instant access tier for the current balance.
		winningTier := CalculateBestInstantAccessTier(runningSavings)

		recordTierBreakpoint(
			&tierBreakpoints,
			months,
			runningSavings,
			previousTier,
			winningTier,
		)

		previousTier = &winningTier

		// Apply the tier fee.
		runningSavings -= winningTier.Fee
		fees += winningTier.Fee

		// Calculate and record interest gained.
		monthlyInterest := (runningSavings * winningTier.AER) / 10000 / 12
		runningSavings += monthlyInterest
		interestGain += monthlyInterest

		// Add the monthly saving allocation.
		runningSavings += monthlySave

		// One month of the forecast has now passed.
		months++

		// Check if the baseline buffer has been reached.
		if runningSavings >= baselineBuffer {
			surplus := runningSavings - baselineBuffer

			return BufferForecast{
				Phase1Months:       months,
				Phase1Surplus:      surplus,
				Phase1InterestGain: interestGain,
				Phase1Fees:         fees,
				TierBreakpoints:    tierBreakpoints,
			}, nil
		}
	}

	return BufferForecast{
		Phase1Months:       months,
		Phase1Surplus:      0,
		Phase1InterestGain: interestGain,
		Phase1Fees:         fees,
		TierBreakpoints:    tierBreakpoints,
	}, nil
}

// CalculateDebtTimeline simulates the debt balance from Phase 1 through
// the completion of Phase 2.
type DebtForecast struct {
	Phase2Months          int64
	Phase2Surplus         int64
	Phase1InterestLost    int64
	Phase2InterestLost    int64
	InterestOver50Percent bool
}

func CalculateDebtTimeline(
	cf CurrentFinances,
	monthlySave,
	phase1Months,
	initialSurplus int64,
) (DebtForecast, error) {

	if monthlySave <= 0 {
		return DebtForecast{}, ErrZeroSavingAllocation
	}

	// If there is no active debt, Phase 2 is bypassed.
	if cf.UnsettledDebt <= 0 {
		return DebtForecast{
			Phase2Months:          0,
			Phase2Surplus:         initialSurplus,
			Phase1InterestLost:    0,
			Phase2InterestLost:    0,
			InterestOver50Percent: false,
		}, nil
	}

	// Calculate the debt growth that occurred during Phase 1.
	debtGrowth := CalculateDebtGrowthDuringBuffer(
		cf,
		monthlySave,
		phase1Months,
	)

	runningDebt := debtGrowth.FinalDebt
	phase1InterestLost := debtGrowth.Phase1InterestLost

	phase2InterestLost := int64(0)
	phase2Months := int64(0)

	// Carry the Phase 1 result into Phase 2.
	interestOver50Percent := debtGrowth.InterestOver50Percent

	// Apply the surplus from the final Phase 1 month to the debt.
	if initialSurplus >= runningDebt {
		leftoverSurplus := initialSurplus - runningDebt

		return DebtForecast{
			Phase2Months:          0,
			Phase2Surplus:         leftoverSurplus,
			Phase1InterestLost:    phase1InterestLost,
			Phase2InterestLost:    0,
			InterestOver50Percent: interestOver50Percent,
		}, nil
	}

	runningDebt -= initialSurplus

	// Use the user's debt interest rate or the 24% AER fallback.
	interestAER := cf.DebtInterestRate
	if interestAER == 0 {
		interestAER = 2400
	}

	// Active Phase 2 debt repayment.
	for runningDebt > 0 {
		phase2Months++

		// Debt grows from monthly interest.
		monthlyInterest := (runningDebt * interestAER) / 10000 / 12

		runningDebt += monthlyInterest
		phase2InterestLost += monthlyInterest

		// Only check Phase 2 if Phase 1 has not already triggered
		// the 50% interest threshold.
		if !interestOver50Percent && monthlyInterest > monthlySave/2 {
			interestOver50Percent = true
		}

		// Check if this month's contribution completely clears the debt.
		if monthlySave >= runningDebt {
			leftoverSurplus := monthlySave - runningDebt

			return DebtForecast{
				Phase2Months:          phase2Months,
				Phase2Surplus:         leftoverSurplus,
				Phase1InterestLost:    phase1InterestLost,
				Phase2InterestLost:    phase2InterestLost,
				InterestOver50Percent: interestOver50Percent,
			}, nil
		}

		// Apply the full monthly contribution.
		runningDebt -= monthlySave
	}

	return DebtForecast{
		Phase2Months:          phase2Months,
		Phase2Surplus:         0,
		Phase1InterestLost:    phase1InterestLost,
		Phase2InterestLost:    phase2InterestLost,
		InterestOver50Percent: interestOver50Percent,
	}, nil
}

// EmergencyFundForecast holds the Phase 3 timing result and financial results.
type EmergencyFundForecast struct {
	Phase3Months       int64
	Phase3Surplus      int64
	Phase3InterestGain int64
	Phase3Fees         int64
	TierBreakpoints    []TierBreakpoint
}

// CalculateEmergencyFundTimeline simulates building the emergency fund
// with compounding interest.
func CalculateEmergencyFundTimeline(
	startingSavings,
	monthlySave,
	targetAmount int64,
) (EmergencyFundForecast, error) {

	if monthlySave <= 0 {
		return EmergencyFundForecast{}, ErrZeroSavingAllocation
	}

	// If the emergency fund is already fully funded, Phase 3 takes 0 months.
	if startingSavings >= targetAmount {
		return EmergencyFundForecast{
			Phase3Months:       0,
			Phase3Surplus:      startingSavings - targetAmount,
			Phase3InterestGain: 0,
			Phase3Fees:         0,
		}, nil
	}

	months := int64(0)
	runningSavings := startingSavings
	interestGain := int64(0)
	fees := int64(0)

	var tierBreakpoints []TierBreakpoint
	var previousTier *InstantAccessTier

	// Step forward month by month until the emergency fund target is reached.
	for runningSavings < targetAmount {

		// Calculate the best instant access tier for the current balance.
		winningTier := CalculateBestInstantAccessTier(runningSavings)

		recordTierBreakpoint(
			&tierBreakpoints,
			months,
			runningSavings,
			previousTier,
			winningTier,
		)

		previousTier = &winningTier

		// Apply the tier fee.
		runningSavings -= winningTier.Fee
		fees += winningTier.Fee

		// Calculate and record interest gained.
		monthlyInterest := (runningSavings * winningTier.AER) / 10000 / 12
		runningSavings += monthlyInterest
		interestGain += monthlyInterest

		// Add the monthly saving allocation.
		runningSavings += monthlySave

		// One month of the forecast has now passed.
		months++

		// Check if the emergency fund target has been reached.
		if runningSavings >= targetAmount {
			surplus := runningSavings - targetAmount

			return EmergencyFundForecast{
				Phase3Months:       months,
				Phase3Surplus:      surplus,
				Phase3InterestGain: interestGain,
				Phase3Fees:         fees,
				TierBreakpoints:    tierBreakpoints,
			}, nil
		}
	}

	return EmergencyFundForecast{
		Phase3Months:       months,
		Phase3Surplus:      0,
		Phase3InterestGain: interestGain,
		Phase3Fees:         fees,
		TierBreakpoints:    tierBreakpoints,
	}, nil
}

// calculateDebtFreedomPlan generates the complete forecast for one debt strategy.
// It runs Phase 1 (buffer) followed by Phase 2 (debt payoff).
type DebtFreedomPlan struct {
	Available      bool
	Allocation     Allocation
	BaselineBuffer BaselineBuffer
	BufferForecast BufferForecast
	DebtForecast   DebtForecast
	BufferGrowth   BufferGrowthForecast
	DMPRequired    bool
	DMPReasons     []string
}

func calculateDebtFreedomPlan(
	cf CurrentFinances,
	strategy BudgetStrategy,
	baselineBuffer BaselineBuffer,
) (DebtFreedomPlan, error) {

	if !strategy.Available {
		return DebtFreedomPlan{
			Available: false,
		}, nil
	}

	bufferForecast, err := CalculateBufferTimeline(
		cf,
		strategy.Allocations.Save,
		baselineBuffer.TargetAmount,
	)

	if err != nil {
		return DebtFreedomPlan{}, err
	}

	debtForecast, err := CalculateDebtTimeline(
		cf,
		strategy.Allocations.Save,
		bufferForecast.Phase1Months,
		bufferForecast.Phase1Surplus,
	)
	if err != nil {
		return DebtFreedomPlan{}, err
	}

	bufferGrowth := CalculateBufferGrowthDuringDebt(
		baselineBuffer.TargetAmount,
		debtForecast.Phase2Months,
	)

	plan := DebtFreedomPlan{
		Available:      true,
		Allocation:     strategy.Allocations,
		BaselineBuffer: baselineBuffer,
		BufferForecast: bufferForecast,
		DebtForecast:   debtForecast,
		BufferGrowth:   bufferGrowth,
	}

	dmpAssessment := AssessDMPNeed(cf, plan)

	plan.DMPRequired = dmpAssessment.DMPRequired
	plan.DMPReasons = dmpAssessment.Reasons

	return plan, nil
}

// calculateEmergencyFundPlan generates the Phase 3 forecast for one strategy.
type EmergencyFundPlan struct {
	Available    bool
	Allocation   Allocation
	TargetAmount int64
	Forecast     EmergencyFundForecast
}

func calculateEmergencyFundPlan(
	cf CurrentFinances,
	strategy BudgetStrategy,
	targetAmount int64,
) (EmergencyFundPlan, error) {

	if !strategy.Available {
		return EmergencyFundPlan{
			Available: false,
		}, nil
	}

	startingSavings := cf.CurrentSavings

	forecast, err := CalculateEmergencyFundTimeline(
		startingSavings,
		strategy.Allocations.Save,
		targetAmount,
	)
	if err != nil {
		return EmergencyFundPlan{}, err
	}

	return EmergencyFundPlan{
		Available:    true,
		Allocation:   strategy.Allocations,
		TargetAmount: targetAmount,
		Forecast:     forecast,
	}, nil
}

//GenerateCustomDebtFreedomPlan generates a debt strategy
// from the user's chosen monthly contribution.
func GenerateCustomDebtFreedomPlan(
	cf CurrentFinances,
	baselineBuffer BaselineBuffer,
	customContribution int64,
) (DebtFreedomPlan, error) {

	customStrategy, err := GenerateCustomAllocation(
		cf,
		customContribution,
	)
	if err != nil {
		return DebtFreedomPlan{}, err
	}

	return calculateDebtFreedomPlan(
		cf,
		customStrategy,
		baselineBuffer,
	)
}

// GenerateCustomEmergencyFundPlan generates an emergency-fund strategy
// from the user's chosen monthly contribution.
func GenerateCustomEmergencyFundPlan(
	cf CurrentFinances,
	targetAmount int64,
	customContribution int64,
) (EmergencyFundPlan, error) {

	customStrategy, err := GenerateCustomAllocation(
		cf,
		customContribution,
	)
	if err != nil {
		return EmergencyFundPlan{}, err
	}

	return calculateEmergencyFundPlan(
		cf,
		customStrategy,
		targetAmount,
	)
}

// GenerateDebtFreedomStrategies generates the default debt strategies
// for comparison. Custom strategy is added separately if requested.
type DebtFreedomStrategies struct {
	Sustainable DebtFreedomPlan
	Moderate    DebtFreedomPlan
	Aggressive  DebtFreedomPlan
	Custom      DebtFreedomPlan
}

func GenerateDebtFreedomStrategies(
	cf CurrentFinances,
	allocations AllocationOptions,
	baselineBuffer BaselineBuffer,
) (DebtFreedomStrategies, error) {

	sustainable, err := calculateDebtFreedomPlan(
		cf,
		allocations.Sustainable,
		baselineBuffer,
	)
	if err != nil {
		return DebtFreedomStrategies{}, err
	}

	moderate, err := calculateDebtFreedomPlan(
		cf,
		allocations.Moderate,
		baselineBuffer,
	)
	if err != nil {
		return DebtFreedomStrategies{}, err
	}

	aggressive, err := calculateDebtFreedomPlan(
		cf,
		allocations.Aggressive,
		baselineBuffer,
	)
	if err != nil {
		return DebtFreedomStrategies{}, err
	}

	strategies := DebtFreedomStrategies{
		Sustainable: sustainable,
		Moderate:    moderate,
		Aggressive:  aggressive,
	}

	PrintDebtFreedomStrategies(cf, strategies)

	return strategies, nil
}

// GenerateEmergencyFundStrategies generates the default emergency-fund
// strategies for comparison. Custom strategy is added separately if requested.
type EmergencyFundStrategies struct {
	Sustainable EmergencyFundPlan
	Moderate    EmergencyFundPlan
	Aggressive  EmergencyFundPlan
	Custom      EmergencyFundPlan
}

func GenerateEmergencyFundStrategies(
	cf CurrentFinances,
	allocations AllocationOptions,
	targetAmount int64,
) (EmergencyFundStrategies, error) {

	var sustainable EmergencyFundPlan

	if allocations.Sustainable.Available {
		var err error

		sustainable, err = calculateEmergencyFundPlan(
			cf,
			allocations.Sustainable,
			targetAmount,
		)
		if err != nil {
			return EmergencyFundStrategies{}, err
		}
	}

	var moderate EmergencyFundPlan

	if allocations.Moderate.Available {
		var err error

		moderate, err = calculateEmergencyFundPlan(
			cf,
			allocations.Moderate,
			targetAmount,
		)
		if err != nil {
			return EmergencyFundStrategies{}, err
		}
	}

	var aggressive EmergencyFundPlan

	if allocations.Aggressive.Available {
		var err error

		aggressive, err = calculateEmergencyFundPlan(
			cf,
			allocations.Aggressive,
			targetAmount,
		)
		if err != nil {
			return EmergencyFundStrategies{}, err
		}
	}

	strategies := EmergencyFundStrategies{
		Sustainable: sustainable,
		Moderate:    moderate,
		Aggressive:  aggressive,
	}

	PrintEmergencyFundStrategies(cf, strategies)

	return strategies, nil
}
