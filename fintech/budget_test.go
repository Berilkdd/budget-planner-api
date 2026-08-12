package fintech

import (
	"errors"
	"testing"
)

func TestGenerateAllocation(t *testing.T) {
	tests := []struct {
		name         string
		input        CurrentFinances
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

			if options.Moderate.Name != "Moderate Debt Acceleration" {
				t.Errorf("expected Moderate Debt Acceleration, got %s", options.Moderate.Name)
			}

			if options.Aggressive.Name != "Aggressive Debt Acceleration" {
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
			expectedErr:     "",
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
			expectedErr:     "unsettled debt cannot be negative.",
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
			// Copy instance to test local mutations safely
			cf := tt.input

			baselineBuffer, err := CalculateBaselineBuffer(cf)
			if err != nil {
				t.Fatal(err)
			}
			msg, err := cf.ApplyImmediateDebtPayoff(baselineBuffer)

			// Check if error matches our expected text
			if tt.expectedErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.expectedErr)
				}
				if err.Error() != tt.expectedErr && err.Error() != tt.expectedErr {
					t.Errorf("expected error string %q, got %q", tt.expectedErr, err.Error())
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Validate string output message
			if msg != tt.expectedMsg {
				t.Errorf("expected message:\n%q\ngot:\n%q", tt.expectedMsg, msg)
			}

			// Validate pointer state mutations inside CurrentFinances
			if cf.CurrentSavings != tt.expectedSavings {
				t.Errorf("expected mutated savings %d, got %d", tt.expectedSavings, cf.CurrentSavings)
			}
			if cf.UnsettledDebt != tt.expectedDebt {
				t.Errorf("expected mutated debt %d, got %d", tt.expectedDebt, cf.UnsettledDebt)
			}
		})
	}
}

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

func TestCalculateBufferTimeline(t *testing.T) {
	tests := []struct {
		name           string
		input          CurrentFinances
		monthlySave    int64
		baselineBuffer int64
		expectedMonths int64
		expectedSurplus int64
		expectedErr    error
	}{
		{
			name: "Zero monthly saving allocation",
			input: CurrentFinances{
				CurrentSavings: 50000,
				UnsettledDebt:  100000,
				HasDebt:        true,
			},
			monthlySave:     0,
			baselineBuffer: 100000,
			expectedMonths:  0,
			expectedSurplus: 0,
			expectedErr:     ErrZeroSavingAllocation,
		},
		{
			name: "No active debt bypasses buffer phase",
			input: CurrentFinances{
				CurrentSavings: 50000,
				UnsettledDebt:  0,
				HasDebt:        false,
			},
			monthlySave:     50000,
			baselineBuffer: 100000,
			expectedMonths:  0,
			expectedSurplus: 0,
			expectedErr:     nil,
		},
		{
			name: "Savings already reaches buffer",
			input: CurrentFinances{
				CurrentSavings: 150000,
				UnsettledDebt:  100000,
				HasDebt:        true,
			},
			monthlySave:     50000,
			baselineBuffer: 100000,
			expectedMonths:  0,
			expectedSurplus: 0,
			expectedErr:     nil,
		},
		{
			name: "Buffer requires multiple months",
			input: CurrentFinances{
				CurrentSavings: 50000,
				UnsettledDebt:  200000,
				HasDebt:        true,
			},
			monthlySave:     25000,
			baselineBuffer: 100000,
			expectedMonths:  2,
			expectedSurplus: 286,
			expectedErr:     nil,
		},
		{
			name: "Final month creates surplus",
			input: CurrentFinances{
				CurrentSavings: 60000,
				UnsettledDebt:  200000,
				HasDebt:        true,
			},
			monthlySave:     50000,
			baselineBuffer: 100000,
			expectedMonths:  1,
			expectedSurplus: 10137,
			expectedErr:     nil,
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
		})
	}
}

func TestCalculateDebtTimeline(t *testing.T) {
	tests := []struct {
		name             string
		input            CurrentFinances
		monthlySave      int64
		phase1Months     int64
		initialSurplus   int64
		expectedTotal    int64
		expectedPhase2   int64
		expectedSurplus  int64
		expectedErr      error
	}{
		{
			name: "Zero monthly saving allocation",
			input: CurrentFinances{
				UnsettledDebt: 100000,
			},
			monthlySave:     0,
			phase1Months:    0,
			initialSurplus:  0,
			expectedTotal:   0,
			expectedPhase2:  0,
			expectedSurplus: 0,
			expectedErr:     ErrZeroSavingAllocation,
		},
		{
			name: "No active debt bypasses Phase 2",
			input: CurrentFinances{
				UnsettledDebt: 0,
			},
			monthlySave:     50000,
			phase1Months:    3,
			initialSurplus: 20000,
			expectedTotal:   3,
			expectedPhase2:  0,
			expectedSurplus: 20000,
			expectedErr:     nil,
		},
		{
			name: "Phase 1 surplus completely clears debt",
			input: CurrentFinances{
				UnsettledDebt: 50000,
				DebtInterestRate: 0,
			},
			monthlySave:     50000,
			phase1Months:    2,
			initialSurplus: 100000,
			expectedTotal:   2,
			expectedPhase2:  0,
			expectedSurplus: 47980,
			expectedErr:     nil,
		},
		{
			name: "Debt cleared during first Phase 2 month",
			input: CurrentFinances{
				UnsettledDebt:     50000,
				DebtInterestRate: 0,
			},
			monthlySave:     100000,
			phase1Months:    0,
			initialSurplus:  0,
			expectedTotal:   1,
			expectedPhase2:  1,
			expectedSurplus: 49000,
			expectedErr:      nil,
		},
		{
			name: "Debt requires multiple Phase 2 months",
			input: CurrentFinances{
				UnsettledDebt:     200000,
				DebtInterestRate: 0,
			},
			monthlySave:     50000,
			phase1Months:    0,
			initialSurplus:  0,
			expectedTotal:   5,
			expectedPhase2:  5,
			expectedSurplus: 39387,
			expectedErr:      nil,
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

			if result.TotalMonths != tt.expectedTotal {
				t.Errorf(
					"expected TotalMonths %d, got %d",
					tt.expectedTotal,
					result.TotalMonths,
				)
			}

			if result.Phase2Months != tt.expectedPhase2 {
				t.Errorf(
					"expected Phase2Months %d, got %d",
					tt.expectedPhase2,
					result.Phase2Months,
				)
			}

			if result.Phase2Surplus != tt.expectedSurplus {
				t.Errorf(
					"expected Phase2Surplus %d, got %d",
					tt.expectedSurplus,
					result.Phase2Surplus,
				)
			}
		})
	}
}
func TestSimulateEmergencyFundTiers(t *testing.T) {
	tests := []struct {
		name                string
		input               CurrentFinances
		monthlySave         int64
		phase1Months        int64
		phase2Months        int64
		phase2Surplus       int64
		targetAmount        int64
		baselineBuffer      int64
		expectedErr         error
		expectedStart       int64
		expectedPhase3      int64
		expectedTotalMonths int64
	}{
		{
			name: "Zero monthly saving allocation",
			input: CurrentFinances{
				CurrentSavings: 100000,
				HasDebt:        false,
			},
			monthlySave:         0,
			phase1Months:        0,
			phase2Months:        0,
			phase2Surplus:       0,
			targetAmount:        200000,
			baselineBuffer:      100000,
			expectedErr:         ErrZeroSavingAllocation,
		},
		{
			name: "No debt starts Phase 3 from current savings",
			input: CurrentFinances{
				CurrentSavings: 100000,
				HasDebt:        false,
			},
			monthlySave:         50000,
			phase1Months:        0,
			phase2Months:        0,
			phase2Surplus:       0,
			targetAmount:        100000,
			baselineBuffer:      100000,
			expectedErr:         nil,
			expectedStart:       100000,
			expectedPhase3:      0,
			expectedTotalMonths: 0,
		},
		{
			name: "Debt path starts from buffer plus Phase 2 surplus",
			input: CurrentFinances{
				CurrentSavings: 50000,
				HasDebt:        true,
			},
			monthlySave:         50000,
			phase1Months:        2,
			phase2Months:        0,
			phase2Surplus:       25000,
			targetAmount:        125000,
			baselineBuffer:      100000,
			expectedErr:         nil,
			expectedStart:       125000,
			expectedPhase3:      0,
			expectedTotalMonths: 2,
		},
		{
			name: "Emergency target requires Phase 3 saving",
			input: CurrentFinances{
				CurrentSavings: 100000,
				HasDebt:        false,
			},
			monthlySave:         50000,
			phase1Months:        0,
			phase2Months:        0,
			phase2Surplus:       0,
			targetAmount:        200000,
			baselineBuffer:      100000,
			expectedErr:         nil,
			expectedStart:       100000,
			expectedPhase3:      2,
			expectedTotalMonths: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SimulateEmergencyFundTiers(
				tt.input,
				tt.monthlySave,
				tt.phase1Months,
				tt.phase2Months,
				tt.phase2Surplus,
				tt.targetAmount,
				tt.baselineBuffer,
			)

			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
			}

			if tt.expectedErr != nil {
				return
			}

			if result.TrueStartingBalance != tt.expectedStart {
				t.Errorf(
					"expected TrueStartingBalance %d, got %d",
					tt.expectedStart,
					result.TrueStartingBalance,
				)
			}

			if result.Forecast.Phase3Months != tt.expectedPhase3 {
				t.Errorf(
					"expected Phase3Months %d, got %d",
					tt.expectedPhase3,
					result.Forecast.Phase3Months,
				)
			}

			if result.Forecast.TotalMonths != tt.expectedTotalMonths {
				t.Errorf(
					"expected TotalMonths %d, got %d",
					tt.expectedTotalMonths,
					result.Forecast.TotalMonths,
				)
			}
		})
	}
}