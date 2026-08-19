package fintech

type InstantAccessTier struct {
	Name string
	Fee  int64
	AER  int64
}

// CalculateBestInstantAccessTier compares all instant access tiers and returns the best option for the given balance.
func CalculateBestInstantAccessTier(balance int64) InstantAccessTier {
	tiers := []InstantAccessTier{
		{Name: "Standard Tier (Instant Access)", Fee: 0, AER: 275},
		{Name: "Extra Tier (Instant Access)", Fee: 300, AER: 300},
		{Name: "Perks Tier (Instant Access)", Fee: 700, AER: 325},
		{Name: "Max Tier (Instant Access)", Fee: 1700, AER: 350},
	}

	bestTier := tiers[0]
	bestNetGrowth := balance

	for _, tier := range tiers {
		balanceAfterFee := balance - tier.Fee

		if balanceAfterFee <= 0 {
			continue
		}

		monthlyInterest := (balanceAfterFee * tier.AER) / 10000 / 12
		finalBalance := balanceAfterFee + monthlyInterest

		if finalBalance > bestNetGrowth {
			bestNetGrowth = finalBalance
			bestTier = tier
		}
	}

	return bestTier
}

type TierBreakpoint struct {
	MonthOffset int64
	Balance     int64
	Tier        InstantAccessTier
}

func recordTierBreakpoint(
	breakpoints *[]TierBreakpoint,
	monthOffset int64,
	balance int64,
	previousTier *InstantAccessTier,
	currentTier InstantAccessTier,
) {
	if previousTier != nil && previousTier.Name != currentTier.Name {
		*breakpoints = append(
			*breakpoints,
			TierBreakpoint{
				MonthOffset: monthOffset,
				Balance:     balance,
				Tier:        currentTier,
			},
		)
	}
}

// CalculateExcessSavings identifies surplus savings when the user has no debt and has fully funded their emergency fund.
type ExcessSavingsForecast struct {
	EmergencyFundAmount int64
	InvestmentAmount    int64
	RecommendedTier     InstantAccessTier
}

func CalculateExcessSavings(cf CurrentFinances, emergencyFund EmergencyFund) (ExcessSavingsForecast, error) {
	if cf.UnsettledDebt == 0 && cf.CurrentSavings > emergencyFund.TargetAmount {
		emergencyFundAmount := emergencyFund.TargetAmount
		investmentAmount := cf.CurrentSavings - emergencyFundAmount
		bestTier := CalculateBestInstantAccessTier(emergencyFundAmount)

		return ExcessSavingsForecast{
			EmergencyFundAmount: emergencyFundAmount,
			InvestmentAmount:    investmentAmount,
			RecommendedTier:     bestTier,
		}, nil
	}

	return ExcessSavingsForecast{}, nil
}
