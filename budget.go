package budget

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
)

// Summary contains financial data for each account and the specified month.
type Summary struct {
	Accounts *Accounts
	Month    *Month
}

// Budget contains the API token and month for the budget.
type Budget struct {
	APIToken string
	Month    string
}

// Budget retrieves the account and month data for the budget.
func (b *Budget) Budget() Summary {
	a, _ := b.GetAccounts()
	m, _ := b.GetMonth()

	return Summary{Accounts: AccessData[Accounts](a), Month: AccessData[Month](m)}
}

// GetAccounts retrieves account data from the YNAB API.
func (b *Budget) GetAccounts() ([]byte, error) {
	url := "https://api.youneedabudget.com/v1/budgets/last-used/accounts"
	a, err := GetURL(b.APIToken, url)
	if err != nil {
		fmt.Println("Error getting account data:", err)
		return nil, err
	}

	return a, nil
}

// GetMonth retrieves the month data from the YNAB API.
func (b *Budget) GetMonth() ([]byte, error) {
	// Gather month-specific budget data
	url := fmt.Sprintf("https://api.youneedabudget.com/v1/budgets/last-used/months/%s", b.Month)
	m, err := GetURL(b.APIToken, url)
	if err != nil {
		fmt.Println("Error getting month data:", err)
		return nil, err
	}

	return m, nil
}

// option is a function that sets a value on the Budget.
type option func(*Budget) error

// WithToken sets the API token for the Budget.
// If this is not set, NewBudget will attempt to use the `YNAB_PAT` environment variable.
func WithToken(t string) option {
	return func(b *Budget) error {
		if t == "" {
			return errors.New("token cannot be empty")
		}
		b.APIToken = t
		return nil
	}
}

// WithMonth sets the month for the Budget. Expects a string in the format `2024-12-01`, or `current`.
// If this is not set, NewBudget will use the current month.
func WithMonth(t string) option {
	return func(b *Budget) error {
		pattern := `^\d{4}-\d{2}-01$|^current$`
		re, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("error compiling regex: %w", err)
		}
		if !re.MatchString(t) {
			return errors.New("month must be `current` or match the `YYYY-MM-01` format")
		}
		b.Month = t
		return nil
	}
}

// NewBudget creates a new Budget with the provided options.
func NewBudget(opts ...option) (*Budget, error) {
	b := &Budget{
		Month:    "current",
		APIToken: os.Getenv("YNAB_PAT"),
	}

	for _, opt := range opts {
		err := opt(b)
		if err != nil {
			return nil, err
		}
	}

	if b.APIToken == "" {
		return nil, errors.New("API token is required. Set the `YNAB_PAT` environment variable or specify the `WithToken` option")
	}

	return b, nil
}

// DefaultBudget creates a new Budget with the default options and prints a summary for the Month.
func DefaultBudget() int {
	b, err := NewBudget()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	summary := b.Budget()
	summary.Month.Report()
	return 0
}

// AccessData reads the JSON data and unmarshals it into the provided struct.
func AccessData[T Month | Accounts](jsonData []byte) *T {
	var data T
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return nil
	}

	return &data
}

// AccessDataFiles reads the JSON data from the specified file and unmarshals it into the provided struct.
func AccessDataFiles[T Month | Accounts](fileLocation string) *T {
	// Read the JSON data from the file
	jsonData, err := os.ReadFile(fileLocation)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	return AccessData[T](jsonData)
}
