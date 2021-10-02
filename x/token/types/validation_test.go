package types_test

import (
	"testing"
)

func TestValidateSymbol(t *testing.T) {
	type args struct {
		symbol string
	}

	tests := []struct {
		testCase   string
		args	   args
		expectPass bool
	}{
		{
			testCase:   "right case",
			args:       args{symbol: "gauss"},
			expectPass: false,
		},
		{
			testCase:   "start with a capital letter",
			args:       args{symbol: "Gauss"},
			expectPass: true,
		},
		{
			testCase:   "contain a capital letter",
			args:       args{symbol: "gAuss"},
			expectPass: true,
		},
		{
			testCase:   "less than 3 characters in length",
			args:       args{symbol: "ab"},
			expectPass: true,
		},
		{
			testCase:   "equal 64 characters in length",
			args:       args{symbol: "gauss1234567btc1234567btc1234567btc1234567btc1234567btc1234567bct1"},
			expectPass: false,
		},
		{
			testCase:   "more than 64 characters in length",
			args:       args{symbol: "gauss1234567btc1234567btc1234567btc1234567btc1234567btc1234567bct12"},
			expectPass: true,
		},
		{
			testCase:   "contain peg",
			args:       args{symbol: "aaaaaa"},
			expectPass: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.testCase, func(t *testing.T) {
			if err := ValidateSymbol(tc.args.symbol); (err != nil) != tc.expectPass {
				t.Errorf("ValidateSymbol() error = %v, expectPass = %v", err, tc.expectPass)
			}
		})
	}
}
