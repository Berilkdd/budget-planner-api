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
