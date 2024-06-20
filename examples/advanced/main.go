package main

import (
	"fmt"
	"github.com/cassamajor/budget"
)

func main() {
	b, err := budget.NewBudget()
	if err != nil {
		fmt.Println("Error creating budget:", err)
		return
	}

	// Get the accounts and month data
	summary := b.Budget()
	Report(summary.Month, summary.Accounts)
}

// Report prints a summary of the budget for the month.
func Report(m *budget.Month, a *budget.Accounts) {
	nw := a.CalculateNetWorth()

	m.Income = 7410.90 // Forecasted income
	spread := m.Income - m.Underfunded - m.Assigned

	hsa, _ := budget.AccountMap["HSA"]
	auto, _ := budget.AccountMap["Car Payment 🚗🚘"]
	homeValue := budget.Balance(330500)
	equity := homeValue + nw.Mortgage
	retirement := nw.NonLiquidAssets - hsa.Balance

	fmt.Printf("Our Net Worth is %v.\n\n", nw.Total())
	fmt.Printf("We have %v cash in our bank accounts.\n", nw.LiquidAssets)
	fmt.Printf("We have %v saved towards medical expenses.\n", hsa.Balance)
	fmt.Printf("We have %v saved towards retirement.\n", retirement)
	fmt.Printf("We have %v equity in our home.\n\n", equity)
	fmt.Printf("We need an additional %v to fund this month. At the end of the month, we will have roughly %v remaining. Unplanned purchases are not accounted for in this estimate.\n\n", m.Underfunded, spread)
	fmt.Printf("We have %v remaining on our mortgage.\n", -nw.Mortgage)
	fmt.Printf("We have %v in student loans.\n", -nw.StudentLoans)
	fmt.Printf("We have %v in auto loans.\n", -auto.Balance)
}
