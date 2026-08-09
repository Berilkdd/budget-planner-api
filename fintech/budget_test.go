package fintech

import (
	"errors"
	"testing"
)

// Verifies that the budget allocator correctly matches strategies and catches deficits.
func TestGenerateAllocation(t *testing.T) {
	tests := []struct {
		name         string
		income       int64
		needs        int64
		expectedName string
		expectedErr  error
	}{
		{
			name:         "Ideal Framework 50/30/20", 
			income:       300000, 
			needs:        120000, 
			expectedName: "50/30/20 Strategy (Ideal)", 
			expectedErr:  nil,
		},
		{
			name:         "Fallback Framework 60/30/10", 
			income:       300000, 
			needs:        165000, 
			expectedName: "60/30/10 Strategy (Adjusted)", 
			expectedErr:  nil,
		},
		{
			name:         "Budget Deficit Error Triggered", 
			income:       300000, 
			needs:        200000, 
			expectedName: "", 
			expectedErr:  ErrNeedsTooHigh,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strat, err := GenerateAllocation(tt.income, tt.needs)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if strat.Name != tt.expectedName {
				t.Errorf("expected strategy %s, got %s", tt.expectedName, strat.Name)
			}
		})
	}
} 
