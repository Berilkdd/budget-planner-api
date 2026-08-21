package fintech

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

//assesment.go tests

func TestAssessDeficitPosition(t *testing.T) {

	tests := []struct {
		name             string
		finances         CurrentFinances
		expectedEF       int64
		expectedEFMonths int64
	}{
		{
			name: "1 - savings >= emergency fund + debt",
			finances: CurrentFinances{
				Income:           6000,
				Needs:            7000,
				CurrentSavings:   26000,
				UnsettledDebt:    5000,
				EmploymentStatus: Employee,
			},
			expectedEF:       21000,
			expectedEFMonths: 3,
		},

		{
			name: "2 - savings < debt + emergency fund and savings > emergency fund",
			finances: CurrentFinances{
				Income:           6000,
				Needs:            7000,
				CurrentSavings:   23000,
				UnsettledDebt:    5000,
				EmploymentStatus: Employee,
			},
			expectedEF:       21000,
			expectedEFMonths: 3,
		},

		{
			name: "3 - savings < emergency fund and debt > 0",
			finances: CurrentFinances{
				Income:           6000,
				Needs:            7000,
				CurrentSavings:   15000,
				UnsettledDebt:    5000,
				EmploymentStatus: Employee,
			},
			expectedEF:       21000,
			expectedEFMonths: 3,
		},

		{
			name: "4 - savings == emergency fund and debt > 0",
			finances: CurrentFinances{
				Income:           6000,
				Needs:            7000,
				CurrentSavings:   21000,
				UnsettledDebt:    5000,
				EmploymentStatus: Employee,
			},
			expectedEF:       21000,
			expectedEFMonths: 3,
		},

		{
			name: "5 - emergency fund < savings and no debt",
			finances: CurrentFinances{
				Income:           6000,
				Needs:            7000,
				CurrentSavings:   25000,
				UnsettledDebt:    0,
				EmploymentStatus: Employee,
			},
			expectedEF:       21000,
			expectedEFMonths: 3,
		},

		{
			name: "6 - emergency fund > savings and no debt",
			finances: CurrentFinances{
				Income:           6000,
				Needs:            7000,
				CurrentSavings:   15000,
				UnsettledDebt:    0,
				EmploymentStatus: Employee,
			},
			expectedEF:       21000,
			expectedEFMonths: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			AssessDeficitPosition(tt.finances)

			emergencyFund, err := CalculateEmergencyTarget(
				tt.finances.EmploymentStatus,
				tt.finances.Needs,
			)

			if err != nil {
				t.Fatalf("unexpected error calculating emergency fund: %v", err)
			}

			if emergencyFund.TargetAmount != tt.expectedEF {
				t.Errorf(
					"expected emergency fund target £%d, got £%d",
					tt.expectedEF,
					emergencyFund.TargetAmount,
				)
			}

			if emergencyFund.MonthsCount != tt.expectedEFMonths {
				t.Errorf(
					"expected emergency fund months %d, got %d",
					tt.expectedEFMonths,
					emergencyFund.MonthsCount,
				)
			}
		})
	}
}

func TestAssessDMPNeed(t *testing.T) {
	tests := []struct {
		name            string
		cf              CurrentFinances
		aggressive      DebtFreedomPlan
		wantDMP         bool
		wantReasonCount int
	}{
		{
			name: "Interest exceeds 50 percent of monthly contribution",
			cf: CurrentFinances{
				Income:        300000,
				UnsettledDebt: 1000000,
			},
			aggressive: DebtFreedomPlan{
				DebtForecast: DebtForecast{
					InterestOver50Percent: true,
					Phase2Months:          20,
				},
			},
			wantDMP:         true,
			wantReasonCount: 1,
		},
		{
			name: "Debt exceeds 50 percent of annual income",
			cf: CurrentFinances{
				Income:        200000,
				UnsettledDebt: 1500000,
			},
			aggressive: DebtFreedomPlan{
				DebtForecast: DebtForecast{
					InterestOver50Percent: false,
					Phase2Months:          20,
				},
			},
			wantDMP:         true,
			wantReasonCount: 1,
		},
		{
			name: "Aggressive payoff takes more than 36 months",
			cf: CurrentFinances{
				Income:        300000,
				UnsettledDebt: 500000,
			},
			aggressive: DebtFreedomPlan{
				BufferForecast: BufferForecast{
					Phase1Months: 10,
				},
				DebtForecast: DebtForecast{
					InterestOver50Percent: false,
					Phase2Months:          27,
				},
			},
			wantDMP:         true,
			wantReasonCount: 1,
		},
		{
			name: "No DMP conditions triggered",
			cf: CurrentFinances{
				Income:        300000,
				UnsettledDebt: 500000,
			},
			aggressive: DebtFreedomPlan{
				BufferForecast: BufferForecast{
					Phase1Months: 10,
				},
				DebtForecast: DebtForecast{
					InterestOver50Percent: false,
					Phase2Months:          14,
				},
			},
			wantDMP:         false,
			wantReasonCount: 0,
		},
		{
			name: "Multiple DMP conditions triggered",
			cf: CurrentFinances{
				Income:        200000,
				UnsettledDebt: 1500000,
			},
			aggressive: DebtFreedomPlan{
				BufferForecast: BufferForecast{
					Phase1Months: 20,
				},
				DebtForecast: DebtForecast{
					InterestOver50Percent: true,
					Phase2Months:          28,
				},
			},
			wantDMP:         true,
			wantReasonCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AssessDMPNeed(tt.cf, tt.aggressive)

			if got.DMPRequired != tt.wantDMP {
				t.Errorf(
					"expected DMPRequired %v, got %v",
					tt.wantDMP,
					got.DMPRequired,
				)
			}

			if len(got.Reasons) != tt.wantReasonCount {
				t.Errorf(
					"expected %d reasons, got %d",
					tt.wantReasonCount,
					len(got.Reasons),
				)
			}
		})
	}
}

//budget.go tests

func TestGenerateAllocation(t *testing.T) {
	tests := []struct {
		name                string
		input               CurrentFinances
		expectedErr         error
		expectedNeeds       int64
		expectedSustainable Allocation
		expectedModerate    Allocation
		expectedAggressive  Allocation
	}{
		{
			name: "Three allocation plans calculated correctly",
			input: CurrentFinances{
				Income: 300000, // £3,000
				Needs:  120000, // £1,200
			},
			expectedErr:   nil,
			expectedNeeds: 120000,
			expectedSustainable: Allocation{
				Needs: 120000,
				Wants: 90000, // 30%
				Save:  90000,
			},
			expectedModerate: Allocation{
				Needs: 120000,
				Wants: 75000, // 25%
				Save:  105000,
			},
			expectedAggressive: Allocation{
				Needs: 120000,
				Wants: 60000, // 20%
				Save:  120000,
			},
		},
		{
			name: "Different income and needs calculate correctly",
			input: CurrentFinances{
				Income: 240000, // £2,400
				Needs:  150000, // £1,500
			},
			expectedErr:   nil,
			expectedNeeds: 150000,
			expectedSustainable: Allocation{
				Needs: 150000,
				Wants: 72000, // 30%
				Save:  18000,
			},
			expectedModerate: Allocation{
				Needs: 150000,
				Wants: 60000, // 25%
				Save:  30000,
			},
			expectedAggressive: Allocation{
				Needs: 150000,
				Wants: 48000, // 20%
				Save:  42000,
			},
		},
		{
			name: "Needs above 60 percent still calculates plans",
			input: CurrentFinances{
				Income: 300000, // £3,000
				Needs:  200000, // £2,000 = 66.67%
			},
			expectedErr:   nil,
			expectedNeeds: 200000,
			expectedSustainable: Allocation{
				Needs: 200000,
				Wants: 90000, // 30%
				Save:  10000,
			},
			expectedModerate: Allocation{
				Needs: 200000,
				Wants: 75000, // 25%
				Save:  25000,
			},
			expectedAggressive: Allocation{
				Needs: 200000,
				Wants: 60000, // 20%
				Save:  40000,
			},
		},
		{
			name: "Zero income error boundary",
			input: CurrentFinances{
				Income: 0,
				Needs:  100000,
			},
			expectedErr: ErrZeroIncome,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options, err := GenerateAllocation(tt.input)

			if tt.expectedErr != nil {
				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Check strategy names.
			if options.Sustainable.Name != "Sustainable Plan" {
				t.Errorf("expected Sustainable Plan, got %s", options.Sustainable.Name)
			}

			if options.Moderate.Name != "Moderate Plan" {
				t.Errorf("expected Moderate Debt Acceleration, got %s", options.Moderate.Name)
			}

			if options.Aggressive.Name != "Aggressive Plan" {
				t.Errorf("expected Aggressive Debt Acceleration, got %s", options.Aggressive.Name)
			}

			// Check all three calculations.
			if options.Sustainable.Allocations != tt.expectedSustainable {
				t.Errorf("unexpected Sustainable allocation: got %+v, want %+v",
					options.Sustainable.Allocations, tt.expectedSustainable)
			}

			if options.Moderate.Allocations != tt.expectedModerate {
				t.Errorf("unexpected Moderate allocation: got %+v, want %+v",
					options.Moderate.Allocations, tt.expectedModerate)
			}

			if options.Aggressive.Allocations != tt.expectedAggressive {
				t.Errorf("unexpected Aggressive allocation: got %+v, want %+v",
					options.Aggressive.Allocations, tt.expectedAggressive)
			}
		})
	}
}

func TestGenerateCustomAllocation(t *testing.T) {
	tests := []struct {
		name               string
		input              CurrentFinances
		customContribution int64
		expectedErr        error
		expected           Allocation
	}{
		{
			name: "Custom contribution calculates wants correctly",
			input: CurrentFinances{
				Income: 300000, // £3,000
				Needs:  120000, // £1,200
			},
			customContribution: 80000, // £800
			expectedErr:        nil,
			expected: Allocation{
				Needs: 120000,
				Wants: 100000,
				Save:  80000,
			},
		},
		{
			name: "Different income and needs calculate correctly",
			input: CurrentFinances{
				Income: 240000, // £2,400
				Needs:  150000, // £1,500
			},
			customContribution: 50000, // £500
			expectedErr:        nil,
			expected: Allocation{
				Needs: 150000,
				Wants: 40000,
				Save:  50000,
			},
		},
		{
			name: "Zero custom contribution error",
			input: CurrentFinances{
				Income: 300000,
				Needs:  120000,
			},
			customContribution: 0,
			expectedErr:        ErrZeroSavingAllocation,
		},
		{
			name: "Contribution exceeds available income",
			input: CurrentFinances{
				Income: 300000,
				Needs:  120000,
			},
			customContribution: 200000, // More than £1,800 available
			expectedErr:        ErrContributionExceedsAvailableIncome,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan, err := GenerateCustomAllocation(
				tt.input,
				tt.customContribution,
			)

			if tt.expectedErr != nil {
				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if plan.Name != "Custom Plan" {
				t.Errorf("expected Custom Plan, got %s", plan.Name)
			}

			if plan.Allocations != tt.expected {
				t.Errorf(
					"unexpected Custom allocation: got %+v, want %+v",
					plan.Allocations,
					tt.expected,
				)
			}
		})
	}
}

//operations.go tests

func TestApplyImmediateDebtPayoff(t *testing.T) {
	tests := []struct {
		name            string
		input           CurrentFinances
		expectedMsg     string
		expectedSavings int64
		expectedDebt    int64
		expectedErr     string
	}{
		{
			name: "Savings below one month needs buffer",
			input: CurrentFinances{
				Needs:          120000, // £1,200 buffer
				CurrentSavings: 90000,  // £900
				UnsettledDebt:  50000,  // £500
			},
			expectedMsg:     "No action. Balance is below the recommended safety cushion.",
			expectedSavings: 90000,
			expectedDebt:    50000,
			expectedErr:     "savings do not exceed the required baseline buffer",
		},
		{
			name: "No active debt liability",
			input: CurrentFinances{
				Needs:          180000, // £1,800 buffer
				CurrentSavings: 250000, // £2,500
				UnsettledDebt:  0,
			},
			expectedMsg:     "No action. No active liabilities found. Directing to Emergency Fund goals.",
			expectedSavings: 250000,
			expectedDebt:    0,
			expectedErr:     "",
		},
		{
			name: "Negative debt error boundary",
			input: CurrentFinances{
				Needs:          250000, // £2,500 buffer
				CurrentSavings: 300000, // £3,000
				UnsettledDebt:  -50000, // -£500
			},
			expectedMsg:     "",
			expectedSavings: 300000,
			expectedDebt:    -50000,
			expectedErr:     "unsettled debt cannot be negative",
		},
		{
			name: "Surplus cash clears total debt",
			input: CurrentFinances{
				Needs:          120000, // £1,200 buffer
				CurrentSavings: 300000, // £3,000
				UnsettledDebt:  50000,  // £500
			},
			expectedMsg:     "Debt completely cleared using extra savings. New savings: £2500.00.",
			expectedSavings: 250000,
			expectedDebt:    0,
			expectedErr:     "",
		},
		{
			name: "Surplus cash pays partial debt",
			input: CurrentFinances{
				Needs:          150000, // £1,500 buffer
				CurrentSavings: 220000, // £2,200
				UnsettledDebt:  100000, // £1,000
			},
			expectedMsg:     "Extra savings applied to debt. Remaining savings locked at £1500.00. Remaining debt: £300.00.",
			expectedSavings: 150000,
			expectedDebt:    30000,
			expectedErr:     "",
		},
		{
			name: "Buffer is based on one month of needs",
			input: CurrentFinances{
				Needs:          250000, // £2,500 buffer
				CurrentSavings: 350000, // £3,500
				UnsettledDebt:  200000, // £2,000
			},
			expectedMsg:     "Extra savings applied to debt. Remaining savings locked at £2500.00. Remaining debt: £1000.00.",
			expectedSavings: 250000,
			expectedDebt:    100000,
			expectedErr:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cf := tt.input

			baselineBuffer, err := CalculateBaselineBuffer(cf)
			if err != nil {
				t.Fatal(err)
			}

			err = cf.ApplyImmediateDebtPayoff(baselineBuffer)

			if tt.expectedErr != "" {
				if err == nil {
					t.Fatalf(
						"expected error containing %q, got nil",
						tt.expectedErr,
					)
				}

				if err.Error() != tt.expectedErr {
					t.Errorf(
						"expected error string %q, got %q",
						tt.expectedErr,
						err.Error(),
					)
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if cf.CurrentSavings != tt.expectedSavings {
				t.Errorf(
					"expected mutated savings %d, got %d",
					tt.expectedSavings,
					cf.CurrentSavings,
				)
			}

			if cf.UnsettledDebt != tt.expectedDebt {
				t.Errorf(
					"expected mutated debt %d, got %d",
					tt.expectedDebt,
					cf.UnsettledDebt,
				)
			}
		})
	}
}

func TestApplyDebtFreedomPlan(t *testing.T) {
	startDate := time.Date(2026, time.August, 1, 0, 0, 0, 0, time.UTC)

	cf := CurrentFinances{
		Income:           240000,
		Needs:            100000,
		CurrentSavings:   100000,
		HasDebt:          true,
		UnsettledDebt:    500000,
		DebtInterestRate: 2400,
		CurrentDate:      startDate,
	}

	plan := DebtFreedomPlan{
		Available: true,
		BaselineBuffer: BaselineBuffer{
			TargetAmount: 100000,
		},
		BufferForecast: BufferForecast{
			Phase1Months: 1,
		},
		DebtForecast: DebtForecast{
			Phase2Months:  1,
			Phase2Surplus: 25000,
		},
	}

	err := ApplyDebtFreedomPlan(&cf, plan)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cf.UnsettledDebt != 0 {
		t.Errorf("expected debt to be 0, got %d", cf.UnsettledDebt)
	}

	if cf.HasDebt {
		t.Error("expected HasDebt to be false")
	}

	if cf.CurrentSavings != 125000 {
		t.Errorf("expected current savings to be 125000, got %d", cf.CurrentSavings)
	}

	if cf.AvailableSurplus != 0 {
		t.Errorf("expected available surplus to be 0, got %d", cf.AvailableSurplus)
	}

	expectedDate := time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC)

	if !cf.CurrentDate.Equal(expectedDate) {
		t.Errorf(
			"expected date to be %v, got %v",
			expectedDate,
			cf.CurrentDate,
		)
	}
}

func TestApplyEmergencyFundPlan(t *testing.T) {
	startDate := time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC)

	cf := CurrentFinances{
		Income:           240000,
		Needs:            100000,
		CurrentSavings:   100000,
		AvailableSurplus: 25000,
		HasDebt:          false,
		UnsettledDebt:    0,
		DebtInterestRate: 2400,
		CurrentDate:      startDate,
	}

	plan := EmergencyFundPlan{
		Available:    true,
		TargetAmount: 600000,
		Forecast: EmergencyFundForecast{
			Phase3Months:       3,
			Phase3Surplus:      15000,
			Phase3InterestGain: 5000,
			Phase3Fees:         1000,
		},
	}

	err := ApplyEmergencyFundPlan(&cf, plan)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cf.CurrentSavings != 600000 {
		t.Errorf(
			"expected current savings to be 600000, got %d",
			cf.CurrentSavings,
		)
	}

	if cf.AvailableSurplus != 15000 {
		t.Errorf(
			"expected available surplus to be 15000, got %d",
			cf.AvailableSurplus,
		)
	}

	expectedDate := time.Date(
		2027,
		time.January,
		1,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	if !cf.CurrentDate.Equal(expectedDate) {
		t.Errorf(
			"expected date to be %v, got %v",
			expectedDate,
			cf.CurrentDate,
		)
	}

	if cf.UnsettledDebt != 0 {
		t.Errorf(
			"expected debt to remain 0, got %d",
			cf.UnsettledDebt,
		)
	}

	if cf.HasDebt {
		t.Error("expected HasDebt to remain false")
	}
}

//safety.go

func TestCalculateEmergencyTarget(t *testing.T) {
	tests := []struct {
		name           string
		status         EmploymentStatus
		monthlyNeeds   int64
		expectedAmount int64
		expectedMonths int64
		expectErr      bool
	}{
		{
			name:           "Full-time Employee 3 Months Buffer",
			status:         Employee,
			monthlyNeeds:   120000,
			expectedAmount: 360000,
			expectedMonths: 3,
			expectErr:      false,
		},
		{
			name:           "Self-Employed 6 Months Buffer",
			status:         SelfEmployed,
			monthlyNeeds:   150000,
			expectedAmount: 900000,
			expectedMonths: 6,
			expectErr:      false,
		},
		{
			name:           "Zero Monthly Needs Error Boundary",
			status:         Employee,
			monthlyNeeds:   0,
			expectedAmount: 0,
			expectedMonths: 0,
			expectErr:      true,
		},
		{
			name:           "Invalid Employment Status Handling",
			status:         EmploymentStatus("CONTRACTOR"),
			monthlyNeeds:   100000,
			expectedAmount: 0,
			expectedMonths: 0,
			expectErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fund, err := CalculateEmergencyTarget(tt.status, tt.monthlyNeeds)

			if (err != nil) != tt.expectErr {
				t.Fatalf("unexpected error state: got %v", err)
			}

			if fund.TargetAmount != tt.expectedAmount {
				t.Errorf("expected target amount %d, got %d", tt.expectedAmount, fund.TargetAmount)
			}

			if fund.MonthsCount != tt.expectedMonths {
				t.Errorf("expected months count %d, got %d", tt.expectedMonths, fund.MonthsCount)
			}
		})
	}
}

//savings.go tests

func TestCalculateBestInstantAccessTier(t *testing.T) {
	tests := []struct {
		name     string
		balance  int64
		expected string
	}{
		{
			name:     "£1000 balance",
			balance:  100000,
			expected: "Standard Tier (Instant Access)",
		},
		{
			name:     "£5000 balance",
			balance:  500000,
			expected: "Standard Tier (Instant Access)",
		},
		{
			name:     "£10000 balance",
			balance:  1000000,
			expected: "Standard Tier (Instant Access)",
		},
		{
			name:     "£50000 balance",
			balance:  5000000,
			expected: "Max Tier (Instant Access)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateBestInstantAccessTier(tt.balance)

			if result.Name != tt.expected {
				t.Errorf(
					"expected %s, got %s",
					tt.expected,
					result.Name,
				)
			}
		})
	}
}

func TestCalculateExcessSavings(t *testing.T) {
	tests := []struct {
		name              string
		cf                CurrentFinances
		emergencyFund     EmergencyFund
		wantEmergencyFund int64
		wantInvestment    int64
		wantTier          string
	}{
		{
			name: "No debt and savings above emergency fund",
			cf: CurrentFinances{
				CurrentSavings: 2000000, // £20,000
				UnsettledDebt:  0,
			},
			emergencyFund: EmergencyFund{
				TargetAmount: 1200000, // £12,000
			},
			wantEmergencyFund: 1200000,
			wantInvestment:    800000,
			wantTier:          "Standard Tier (Instant Access)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateExcessSavings(tt.cf, tt.emergencyFund)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.EmergencyFundAmount != tt.wantEmergencyFund {
				t.Errorf("expected emergency fund %d, got %d",
					tt.wantEmergencyFund, got.EmergencyFundAmount)
			}

			if got.InvestmentAmount != tt.wantInvestment {
				t.Errorf("expected investment amount %d, got %d",
					tt.wantInvestment, got.InvestmentAmount)
			}

			if got.RecommendedTier.Name != tt.wantTier {
				t.Errorf("expected %s, got %s",
					tt.wantTier, got.RecommendedTier.Name)
			}
		})
	}
}

//timeline.go tests

func TestCalculateBufferTimeline(t *testing.T) {
	tests := []struct {
		name                 string
		input                CurrentFinances
		monthlySave          int64
		baselineBuffer       int64
		expectedMonths       int64
		expectedSurplus      int64
		expectedInterestGain int64
		expectedFees         int64
		expectedErr          error
	}{
		{
			name: "Zero monthly saving allocation",
			input: CurrentFinances{
				CurrentSavings: 50000,
				UnsettledDebt:  100000,
				HasDebt:        true,
			},
			monthlySave:          0,
			baselineBuffer:       100000,
			expectedMonths:       0,
			expectedSurplus:      0,
			expectedInterestGain: 0,
			expectedFees:         0,
			expectedErr:          ErrZeroSavingAllocation,
		},
		{
			name: "No active debt bypasses buffer phase",
			input: CurrentFinances{
				CurrentSavings: 50000,
				UnsettledDebt:  0,
				HasDebt:        false,
			},
			monthlySave:          50000,
			baselineBuffer:       100000,
			expectedMonths:       0,
			expectedSurplus:      0,
			expectedInterestGain: 0,
			expectedFees:         0,
			expectedErr:          nil,
		},
		{
			name: "Savings already reaches buffer",
			input: CurrentFinances{
				CurrentSavings: 150000,
				UnsettledDebt:  100000,
				HasDebt:        true,
			},
			monthlySave:          50000,
			baselineBuffer:       100000,
			expectedMonths:       0,
			expectedSurplus:      0,
			expectedInterestGain: 0,
			expectedFees:         0,
			expectedErr:          nil,
		},
		{
			name: "Buffer requires multiple months",
			input: CurrentFinances{
				CurrentSavings: 50000,
				UnsettledDebt:  200000,
				HasDebt:        true,
			},
			monthlySave:          25000,
			baselineBuffer:       100000,
			expectedMonths:       2,
			expectedSurplus:      286,
			expectedInterestGain: 286,
			expectedFees:         0,
			expectedErr:          nil,
		},
		{
			name: "Final month creates surplus",
			input: CurrentFinances{
				CurrentSavings: 60000,
				UnsettledDebt:  200000,
				HasDebt:        true,
			},
			monthlySave:          50000,
			baselineBuffer:       100000,
			expectedMonths:       1,
			expectedSurplus:      10137,
			expectedInterestGain: 137,
			expectedFees:         0,
			expectedErr:          nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CalculateBufferTimeline(
				tt.input,
				tt.monthlySave,
				tt.baselineBuffer,
			)

			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}

			if result.Phase1Months != tt.expectedMonths {
				t.Errorf(
					"expected Phase1Months %d, got %d",
					tt.expectedMonths,
					result.Phase1Months,
				)
			}

			if result.Phase1Surplus != tt.expectedSurplus {
				t.Errorf(
					"expected Phase1Surplus %d, got %d",
					tt.expectedSurplus,
					result.Phase1Surplus,
				)
			}

			if result.Phase1InterestGain != tt.expectedInterestGain {
				t.Errorf(
					"expected Phase1InterestGain %d, got %d",
					tt.expectedInterestGain,
					result.Phase1InterestGain,
				)
			}

			if result.Phase1Fees != tt.expectedFees {
				t.Errorf(
					"expected Phase1Fees %d, got %d",
					tt.expectedFees,
					result.Phase1Fees,
				)
			}
		})
	}
}

func TestCalculateDebtTimeline(t *testing.T) {
	tests := []struct {
		name                       string
		input                      CurrentFinances
		monthlySave                int64
		phase1Months               int64
		initialSurplus             int64
		expectedPhase2Months       int64
		expectedPhase2Surplus      int64
		expectedPhase1DebtInterest int64
		expectedPhase2DebtInterest int64
		expectedInterestOver50     bool
		expectedErr                error
	}{
		{
			name: "Zero monthly saving allocation",
			input: CurrentFinances{
				UnsettledDebt: 100000,
			},
			monthlySave:                0,
			phase1Months:               0,
			initialSurplus:             0,
			expectedPhase2Months:       0,
			expectedPhase2Surplus:      0,
			expectedPhase1DebtInterest: 0,
			expectedPhase2DebtInterest: 0,
			expectedInterestOver50:     false,
			expectedErr:                ErrZeroSavingAllocation,
		},
		{
			name: "No active debt bypasses Phase 2",
			input: CurrentFinances{
				UnsettledDebt: 0,
			},
			monthlySave:                50000,
			phase1Months:               3,
			initialSurplus:             20000,
			expectedPhase2Months:       0,
			expectedPhase2Surplus:      20000,
			expectedPhase1DebtInterest: 0,
			expectedPhase2DebtInterest: 0,
			expectedInterestOver50:     false,
			expectedErr:                nil,
		},
		{
			name: "Phase 1 surplus completely clears debt",
			input: CurrentFinances{
				UnsettledDebt:    50000,
				DebtInterestRate: 0,
			},
			monthlySave:                50000,
			phase1Months:               2,
			initialSurplus:             100000,
			expectedPhase2Months:       0,
			expectedPhase2Surplus:      47980,
			expectedPhase1DebtInterest: 2020,
			expectedPhase2DebtInterest: 0,
			expectedInterestOver50:     false,
			expectedErr:                nil,
		},
		{
			name: "  Debt cleared during first Phase 2 month",
			input: CurrentFinances{
				UnsettledDebt:    50000,
				DebtInterestRate: 0,
			},
			monthlySave:                100000,
			phase1Months:               0,
			initialSurplus:             0,
			expectedPhase2Months:       1,
			expectedPhase2Surplus:      49000,
			expectedPhase1DebtInterest: 0,
			expectedPhase2DebtInterest: 1000,
			expectedInterestOver50:     false,
			expectedErr:                nil,
		},
		{
			name: "Debt requires multiple Phase 2 months",
			input: CurrentFinances{
				UnsettledDebt:    200000,
				DebtInterestRate: 0,
			},
			monthlySave:                50000,
			phase1Months:               0,
			initialSurplus:             0,
			expectedPhase2Months:       5,
			expectedPhase2Surplus:      39387,
			expectedPhase1DebtInterest: 0,
			expectedPhase2DebtInterest: 10613,
			expectedInterestOver50:     false,
			expectedErr:                nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CalculateDebtTimeline(
				tt.input,
				tt.monthlySave,
				tt.phase1Months,
				tt.initialSurplus,
			)

			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}

			if result.Phase2Months != tt.expectedPhase2Months {
				t.Errorf(
					"expected Phase2Months %d, got %d",
					tt.expectedPhase2Months,
					result.Phase2Months,
				)
			}

			if result.Phase2Surplus != tt.expectedPhase2Surplus {
				t.Errorf(
					"expected Phase2Surplus %d, got %d",
					tt.expectedPhase2Surplus,
					result.Phase2Surplus,
				)
			}

			if result.Phase1InterestLost != tt.expectedPhase1DebtInterest {
				t.Errorf(
					"expected Phase1DebtInterestPaid %d, got %d",
					tt.expectedPhase1DebtInterest,
					result.Phase1InterestLost,
				)
			}

			if result.Phase2InterestLost != tt.expectedPhase2DebtInterest {
				t.Errorf(
					"expected Phase2DebtInterestPaid %d, got %d",
					tt.expectedPhase2DebtInterest,
					result.Phase2InterestLost,
				)
			}

			if result.InterestOver50Percent != tt.expectedInterestOver50 {
				t.Errorf(
					"expected InterestOver50Percent %v, got %v",
					tt.expectedInterestOver50,
					result.InterestOver50Percent,
				)
			}
		})
	}
}

func TestGenerateDebtFreedomStrategies(t *testing.T) {
	input := CurrentFinances{
		Income:           240000,
		CurrentSavings:   50000,
		Needs:            100000,
		UnsettledDebt:    200000,
		HasDebt:          true,
		DebtInterestRate: 2400,
	}

	allocations := AllocationOptions{
		Sustainable: BudgetStrategy{
			Available: true,
			Allocations: Allocation{
				Save: 5000,
			},
		},
		Moderate: BudgetStrategy{
			Available: true,
			Allocations: Allocation{
				Save: 50000,
			},
		},
		Aggressive: BudgetStrategy{
			Available: true,
			Allocations: Allocation{
				Save: 75000,
			},
		},
	}

	baselineBuffer := BaselineBuffer{
		TargetAmount: 100000,
	}

	result, err := GenerateDebtFreedomStrategies(
		input,
		allocations,
		baselineBuffer,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedSustainable, err := calculateDebtFreedomPlan(
		input,
		allocations.Sustainable,
		baselineBuffer,
	)
	if err != nil {
		t.Fatalf("unexpected Sustainable calculation error: %v", err)
	}

	expectedModerate, err := calculateDebtFreedomPlan(
		input,
		allocations.Moderate,
		baselineBuffer,
	)
	if err != nil {
		t.Fatalf("unexpected Moderate calculation error: %v", err)
	}

	expectedAggressive, err := calculateDebtFreedomPlan(
		input,
		allocations.Aggressive,
		baselineBuffer,
	)
	if err != nil {
		t.Fatalf("unexpected Aggressive calculation error: %v", err)
	}

	if !reflect.DeepEqual(result.Sustainable, expectedSustainable) {
		t.Error("Sustainable strategy does not match calculated plan")
	}

	if !reflect.DeepEqual(result.Moderate, expectedModerate) {
		t.Error("Moderate strategy does not match calculated plan")
	}

	if !reflect.DeepEqual(result.Aggressive, expectedAggressive) {
		t.Error("Aggressive strategy does not match calculated plan")
	}

	// DMP assessment
	if !result.Sustainable.DMPRequired {
		t.Error("expected Sustainable plan to require DMP")
	}

	if result.Moderate.DMPRequired {
		t.Error("expected Moderate plan not to require DMP")
	}

	if result.Aggressive.DMPRequired {
		t.Error("expected Aggressive plan not to require DMP")
	}
}

func TestGenerateEmergencyFundStrategies(t *testing.T) {

	startDate := time.Date(
		2026,
		time.August,
		1,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	cf := CurrentFinances{
		Income:           300000,
		Needs:            150000,
		CurrentSavings:   100000, // £1,000
		HasDebt:          true,
		UnsettledDebt:    100,
		DebtInterestRate: 2400,
		EmploymentStatus: Employee,
		CurrentDate:      startDate,
	}

	allocations := AllocationOptions{
		Sustainable: BudgetStrategy{
			Name:      "Sustainable Plan",
			Available: false,
			Allocations: Allocation{
				Needs: 150000,
				Wants: 90000,
				Save:  60000,
			},
		},

		Moderate: BudgetStrategy{
			Name:      "Moderate Plan",
			Available: true,
			Allocations: Allocation{
				Needs: 150000,
				Wants: 75000, // £750
				Save:  75000, // £750
			},
		},

		Aggressive: BudgetStrategy{
			Name:      "Aggressive Plan",
			Available: true,
			Allocations: Allocation{
				Needs: 150000,
				Wants: 60000, // £600
				Save:  90000, // £900
			},
		},
	}

	targetAmount := int64(300000) // £3,000

	result, err := GenerateEmergencyFundStrategies(
		cf,
		allocations,
		targetAmount,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Sustainable should be unavailable.
	if result.Sustainable.Available {
		t.Error("expected Sustainable plan to be unavailable")
	}

	// Moderate should be available.
	if !result.Moderate.Available {
		t.Error("expected Moderate plan to be available")
	}

	// Aggressive should be available.
	if !result.Aggressive.Available {
		t.Error("expected Aggressive plan to be available")
	}

	// Both available plans should use the correct target.
	if result.Moderate.TargetAmount != targetAmount {
		t.Errorf(
			"expected Moderate target to be %d, got %d",
			targetAmount,
			result.Moderate.TargetAmount,
		)
	}

	if result.Aggressive.TargetAmount != targetAmount {
		t.Errorf(
			"expected Aggressive target to be %d, got %d",
			targetAmount,
			result.Aggressive.TargetAmount,
		)
	}

	// Both plans start with the user's current savings.
	if result.Moderate.Forecast.Phase3Months != 3 {
		t.Errorf(
			"expected Moderate forecast to take 3 months, got %d",
			result.Moderate.Forecast.Phase3Months,
		)
	}

	if result.Aggressive.Forecast.Phase3Months != 3 {
		t.Errorf(
			"expected Aggressive forecast to take 3 months, got %d",
			result.Aggressive.Forecast.Phase3Months,
		)
	}

	// The available plans should retain their allocations.
	if result.Moderate.Allocation.Save != 75000 {
		t.Errorf(
			"expected Moderate saving allocation to be 75000, got %d",
			result.Moderate.Allocation.Save,
		)
	}

	if result.Aggressive.Allocation.Save != 90000 {
		t.Errorf(
			"expected Aggressive saving allocation to be 90000, got %d",
			result.Aggressive.Allocation.Save,
		)
	}

	if result.Moderate.Allocation.Wants != 75000 {
		t.Errorf(
			"expected Moderate wants allocation to be 75000, got %d",
			result.Moderate.Allocation.Wants,
		)
	}

	if result.Aggressive.Allocation.Wants != 60000 {
		t.Errorf(
			"expected Aggressive wants allocation to be 60000, got %d",
			result.Aggressive.Allocation.Wants,
		)
	}
}
