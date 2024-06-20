package budget

import (
	"encoding/json"
	"fmt"
	"time"
)

// AccountMap allows an Account to be queried by name.
var AccountMap = make(map[string]Account)

// Account contains the details of a specific account.
type Account struct {
	ID      string
	Name    string
	Type    string
	Balance Balance
}

// String returns the Name, Balance, and Type as a string. The ID has been omitted.
func (a Account) String() string {
	return fmt.Sprintf("Account Name: %s\nBalance: %v\nType: %s\n", a.Name, a.Balance, a.Type)
}

// Accounts is a collection of Account structs.
type Accounts []Account

// UnmarshalJSON unmarshals the JSON data retrieved from YNAB's API into the Accounts struct.
func (a *Accounts) UnmarshalJSON(data []byte) error {
	var format AccountData
	err := json.Unmarshal(data, &format)
	if err != nil {
		return err
	}

	for _, account := range format.Data.Accounts {
		if !account.Closed {
			acc := Account{
				ID:      account.Id,
				Name:    account.Name,
				Balance: NewRat(account.Balance),
				Type:    account.Type,
			}
			*a = append(*a, acc)

			AccountMap[account.Name] = acc
		}
	}

	return nil
}

// Report prints the account balance for each account.
func (a *Accounts) Report() {
	// Print the account names and balances
	for _, account := range *a {
		fmt.Println(account)
	}
}

// NetWorth contains the calculated balances for each Account type.
type NetWorth struct {
	LiquidAssets    Balance
	NonLiquidAssets Balance
	StudentLoans    Balance
	ConsumerDebt    Balance
	Mortgage        Balance
}

// Total returns the total Net Worth.
func (s NetWorth) Total() Balance {
	return (s.LiquidAssets + s.NonLiquidAssets) - (s.StudentLoans + s.ConsumerDebt + s.Mortgage)
}

// CalculateNetWorth calculates the Net Worth of the Accounts.
func (a *Accounts) CalculateNetWorth() NetWorth {
	var total NetWorth

	for _, account := range *a {
		switch account.Type {
		case "otherAsset": // IRA/HSA/401k
			total.NonLiquidAssets += account.Balance
		case "studentLoan":
			total.StudentLoans += account.Balance
		case "autoLoan", "creditCard":
			total.ConsumerDebt += account.Balance
		case "mortgage":
			total.Mortgage += account.Balance
		case "savings", "checking":
			total.LiquidAssets += account.Balance
		}
	}

	return total
}

// AccountData https://api.ynab.com/v1/#/Accounts/getAccounts
type AccountData struct {
	Data struct {
		Accounts []struct {
			Id                  string     `json:"id"`
			Name                string     `json:"name"`
			Type                string     `json:"type"`
			OnBudget            bool       `json:"on_budget"`
			Closed              bool       `json:"closed"`
			Note                *string    `json:"note"`
			Balance             int        `json:"balance"`
			ClearedBalance      int        `json:"cleared_balance"`
			UnclearedBalance    int        `json:"uncleared_balance"`
			TransferPayeeId     string     `json:"transfer_payee_id"`
			DirectImportLinked  bool       `json:"direct_import_linked"`
			DirectImportInError bool       `json:"direct_import_in_error"`
			LastReconciledAt    *time.Time `json:"last_reconciled_at"`
			DebtOriginalBalance *int       `json:"debt_original_balance"`
			DebtInterestRates   struct {
				Field1 int `json:"2021-10-01,omitempty"`
				Field2 int `json:"2022-06-01,omitempty"`
				Field3 int `json:"2022-11-01,omitempty"`
				Field4 int `json:"2023-09-01,omitempty"`
			} `json:"debt_interest_rates"`
			DebtMinimumPayments struct {
				Field1 int `json:"2021-10-01,omitempty"`
				Field2 int `json:"2022-06-01,omitempty"`
				Field3 int `json:"2022-11-01,omitempty"`
				Field4 int `json:"2023-09-01,omitempty"`
			} `json:"debt_minimum_payments"`
			DebtEscrowAmounts struct {
				Field1 int `json:"2022-08-01,omitempty"`
				Field2 int `json:"2021-10-01,omitempty"`
				Field3 int `json:"2022-11-01,omitempty"`
				Field4 int `json:"2023-09-01,omitempty"`
			} `json:"debt_escrow_amounts"`
			Deleted bool `json:"deleted"`
		} `json:"accounts"`
		ServerKnowledge int `json:"server_knowledge"`
	} `json:"data"`
}
