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

// NetWorth contains the combined balances Assets and Liabilities.
type NetWorth struct {
	Assets      Assets
	Liabilities Liabilities
}

// Total calculates the total Net Worth.
func (n NetWorth) Total() Balance {
	return n.Assets.Total() + n.Liabilities.Total()
}

// Assets contains the calculated balances for each Account type.
type Assets struct {
	Cash     Balance
	Checking Balance
	Savings  Balance
	Other    Balance
}

// Total calculates the balance across all Assets.
func (a Assets) Total() Balance {
	return a.Cash + a.Checking + a.Savings + a.Other
}

// Liabilities contains balances for Accounts that have an outstanding balance
type Liabilities struct {
	AutoLoans     Balance
	CreditCards   Balance
	StudentLoans  Balance
	Mortgages     Balance
	LinesOfCredit Balance
	PersonalLoans Balance
	MedicalDebt   Balance
	Other         Balance
}

// Total calculates the balance across all Liabilities.
func (l Liabilities) Total() Balance {
	return l.AutoLoans + l.CreditCards + l.StudentLoans + l.Mortgages + l.LinesOfCredit + l.PersonalLoans + l.MedicalDebt + l.Other
}

// NetWorth calculates the balance for all Account types.
func (a *Accounts) NetWorth() NetWorth {
	var nw NetWorth

	for _, account := range *a {
		switch account.Type {
		// Assets
		case "cash":
			nw.Assets.Cash += account.Balance
		case "checking":
			nw.Assets.Checking += account.Balance
		case "savings":
			nw.Assets.Savings += account.Balance
		case "otherAsset": // IRA/HSA/401k/Brokerage
			nw.Assets.Other += account.Balance

		// Liabilities
		case "autoLoan":
			nw.Liabilities.AutoLoans += account.Balance
		case "creditCard":
			nw.Liabilities.CreditCards += account.Balance
		case "studentLoan":
			nw.Liabilities.StudentLoans += account.Balance
		case "mortgage":
			nw.Liabilities.Mortgages += account.Balance
		case "lineOfCredit":
			nw.Liabilities.LinesOfCredit += account.Balance
		case "personalLoan":
			nw.Liabilities.PersonalLoans += account.Balance
		case "medicalDebt":
			nw.Liabilities.MedicalDebt += account.Balance
		case "otherLiability", "otherDebt":
			nw.Liabilities.Other += account.Balance
		}
	}

	return nw
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
