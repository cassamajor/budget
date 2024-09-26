package budget

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"
)

// Budget contains financial data for each Account and the specified Month.
type Budget struct {
	Accounts *Accounts
	Month    *Month
}

// Session contains the API token and month for the budget.
type Session struct {
	APIToken string
	Month    string
}

// Summary retrieves the Account and Month data for the budget.
func (s *Session) Summary() Budget {
	a, _ := s.GetAccounts()
	m, _ := s.GetMonth()

	return Budget{Accounts: AccessData[Accounts](a), Month: AccessData[Month](m)}
}

// GetAccounts retrieves Account data from the YNAB API.
func (s *Session) GetAccounts() ([]byte, error) {
	url := "https://api.youneedabudget.com/v1/budgets/last-used/accounts"
	a, err := GetURL(s.APIToken, url)
	if err != nil {
		fmt.Println("Error getting account data:", err)
		return nil, err
	}

	return a, nil
}

// GetMonth retrieves the Month data from the YNAB API.
func (s *Session) GetMonth() ([]byte, error) {
	// Gather month-specific budget data
	url := fmt.Sprintf("https://api.youneedabudget.com/v1/budgets/last-used/months/%s", s.Month)
	m, err := GetURL(s.APIToken, url)
	if err != nil {
		fmt.Println("Error getting month data:", err)
		return nil, err
	}

	return m, nil
}

// option is a function that sets a value on the Session.
type option func(*Session) error

// WithToken sets the API token for the Session.
// If this is not set, NewSession will attempt to use the `YNAB_PAT` environment variable.
func WithToken(t string) option {
	return func(s *Session) error {
		if t == "" {
			return errors.New("token cannot be empty")
		}
		s.APIToken = t
		return nil
	}
}

// WithMonth sets the month for the Session. Expects a string in the format `2024-12-01`, or `current`.
// If this is not set, NewSession will use the current month.
func WithMonth(t string) option {
	return func(s *Session) error {
		pattern := `^\d{4}-\d{2}-01$|^current$`
		re, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("error compiling regex: %w", err)
		}
		if !re.MatchString(t) {
			return errors.New("month must be `current` or match the `YYYY-MM-01` format")
		}
		s.Month = t
		return nil
	}
}

// NewSession creates a new Session with the provided options.
func NewSession(opts ...option) (*Session, error) {
	s := &Session{
		Month:    "current",
		APIToken: os.Getenv("YNAB_PAT"),
	}

	for _, opt := range opts {
		err := opt(s)
		if err != nil {
			return nil, err
		}
	}

	if s.APIToken == "" {
		return nil, errors.New("API token is required. Set the `YNAB_PAT` environment variable or specify the `WithToken` option")
	}

	return s, nil
}

// DefaultSession creates a new Session with the default options and prints a summary for the Month.
func DefaultSession() int {
	s, err := NewSession()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	b := s.Summary()
	b.Month.Report()
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

func (s *Session) Rollover() (Balance, error) {
	t := s.Date()

	// Subtract one month to get the prior month. Also set the date to the first day of the month, as expected by the YNAB API.
	date := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())

	// Re-format the date as a YYYY-MM-01 string
	f := date.Format("2006-01-02")

	// Get the accounts and month data
	month := WithMonth(f)
	session, err := NewSession(month)

	if err != nil {
		token := WithToken(s.APIToken)
		session, err = NewSession(month, token)

		if err != nil {
			return 0, err
		}
	}

	m := session.Summary().Month
	rollover := m.Income - m.Expenses

	if rollover != 0 {
		return rollover, nil
	}

	return 0, nil
}

// Date converts a string into a time.Time
func (s *Session) Date() time.Time {
	var t time.Time

	switch s.Month {
	case "current":
		t = time.Now().Local()
	default:
		// Parse the month into a time.Time object
		t, _ = time.Parse("2006-01-02", s.Month)
	}

	return t
}
