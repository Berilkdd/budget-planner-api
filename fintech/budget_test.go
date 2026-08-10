package fintech

import (
	"errors"
	"testing"
)

func TestGenerateAllocation(t *testing.T) {
	tests := []struct {
		name         string
		input        CurrentFinances
		expectedName string
		expectedErr  error
	}{
		{
			name: "Ideal Framework <= 50%",
			input: CurrentFinances{
				Income: 300000,
				Needs:  120000,
			},
			expectedName: "Ideal Framework (Needs <= 50%)",
			expectedErr:  nil,
		},
		{
			name: "Balanced Framework between 50% and 60%",
			input: CurrentFinances{
				Income: 300000,
				Needs:  165000,
			},
			expectedName: "Not Ideal But Balanced (50% < Needs <= 60%)",
			expectedErr:  nil,
		},
		{
			name: "Budget Deficit Error Triggered",
			input: CurrentFinances{
				Income: 300000,
				Needs:  200000,
			},
			expectedName: "",
			expectedErr:  ErrNeedsTooHigh,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Pass the entire CurrentFinances object directly into the function
			strat, err := GenerateAllocation(tt.input)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if strat.Name != tt.expectedName {
				t.Errorf("expected strategy %s, got %s", tt.expectedName, strat.Name)
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
			name: "Savings below baseline cushion",
			input: CurrentFinances{
				CurrentSavings: 80000, 
				UnsettledDebt:  50000, 
			},
			expectedMsg:     "No action. Balance is below the £1000.00 safety cushion.",
			expectedSavings: 80000,
			expectedDebt:    50000,
			expectedErr:     "",
		},
		{
			name: "No active debt liability",
			input: CurrentFinances{
				CurrentSavings: 150000, 
				UnsettledDebt:  0,
			},
			expectedMsg:     "No action. No active liabilities found. Directing to Emergency Fund goals.",
			expectedSavings: 150000,
			expectedDebt:    0,
			expectedErr:     "",
		},
		{
			name: "Negative debt error boundary",
			input: CurrentFinances{
				CurrentSavings: 150000, 
				UnsettledDebt:  -50000, 
			},
			expectedMsg:     "",
			expectedSavings: 150000,
			expectedDebt:    -50000,
			expectedErr:     "unsettled debt cannot be negative.",
		},
		{
			name: "Surplus cash clears total debt",
			input: CurrentFinances{
				CurrentSavings: 250000, // (Surplus: £1,500)
				UnsettledDebt:  50000,  
			},
			expectedMsg:     "Debt completely cleared using extra savings. New savings: £2000.00.",
			expectedSavings: 200000, 
			expectedDebt:    0,      
			expectedErr:     "",
		},
		{
			name: "Surplus cash pays partial debt",
			input: CurrentFinances{
				CurrentSavings: 150000, // (Surplus: £500)
				UnsettledDebt:  120000, 
			},
			expectedMsg:     "Extra savings applied to debt. Remaining savings locked at £1000.00. Remaining debt: £700.00.",
			expectedSavings: 100000, 
			expectedDebt:    70000,  
			expectedErr:     "",
		},
	}


	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Copy instance to test local mutations safely
			cf := tt.input
			msg, err := cf.ApplyImmediateDebtPayoff()

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
