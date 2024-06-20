package budget_test

import (
	"github.com/cassamajor/budget"
	"testing"
)

func TestNewRat(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected budget.Balance
	}{
		{"Convert an int into a Balance", 1000000, 1000},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := budget.NewRat(tc.input)
			if got != tc.expected {
				t.Errorf("Expected: %v, Got: %v", tc.expected, got)
			}
		})
	}
}

//func TestNewBudget(t *testing.T) {
//	t.Parallel()
//
//	tests := []struct {
//		name string
//		args budget.Budget
//		want budget.Summary
//	}{
//		{
//			name: "No Token is set as an environment variable",
//			args: budget.Budget{
//				Month: "2024-06-01",
//			},
//			want: errors.New("API token is required"),
//		},
//		{
//			name: "WithToken is passed with an empty string",
//			args: budget.Budget{
//				Month:    "2024-06-01",
//				APIToken: "",
//			},
//			want: errors.New("token cannot be empty"),
//		},
//		{
//			name: "Token is set as an environment variable",
//			args: budget.Budget{
//				Month:    "2024-06-01",
//				APIToken: "",
//			},
//			want: budget.Summary{},
//		},
//		{
//			name: "WithMonth is passed with an empty string",
//			args: budget.Budget{
//				Month: "",
//			},
//			want: errors.New("token cannot be empty"),
//		},
//		{
//			name: "No Month Set",
//			args: budget.Budget{
//				APIToken: "gibberish",
//			},
//			want: budget.Summary{},
//		},
//		{
//			name: "Month is set with invalid value",
//			args: budget.Budget{
//				APIToken: "gibberish",
//				Month:    "01-02-2006",
//			},
//			want: budget.Summary{},
//		},
//		{
//			name: "Token and Month Set",
//			args: budget.Budget{
//				APIToken: "gibberish",
//				Month:    "2024-06-01",
//			},
//			want: budget.Summary{},
//		},
//		{
//			name: "Token and Month Unset",
//			args: budget.Budget{},
//			want: budget.Summary{},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			//t.Parallel()
//
//			token := budget.WithToken(tt.args.APIToken)
//			m := budget.WithMonth(tt.args.Month)
//			b, err := budget.NewBudget(token, m)
//
//			if err != nil {
//				t.Fatal(err)
//			}
//
//			if got := b.Budget(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("got =\n %v, want =\n %v", got, tt.want)
//			}
//		})
//	}
//}
