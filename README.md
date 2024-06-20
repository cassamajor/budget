# Busy People Budget

This app queries the You Need A Budget (YNAB) API to quickly generate a snapshot of current finances.

## Install
```shell
go get github.com/cassamajor/budget
```

## Usage
<details>
<summary>Simple Summary: Income v Expenses</summary>

```go
package main

import "github.com/cassamajor/budget"

func main() {
	budget.DefaultBudget()
}
```

Set your [Personal Access Token](https://api.ynab.com/#personal-access-tokens) as an environment variable, and execute the program:
```shell
export YNAB_PAT=personal-access-token
go run main.go
```

Output:
```text
Ready to Assign: $529.15
Assigned: $3,703.49
Underfunded: $5,390.30
Income: $3,812.50
Expenses: $1,737.59
```
</details>


<details>
<summary>Advanced Summary: Detailed and Personalized Insights</summary>

The example can be found here: [examples/advanced/main.go](./examples/advanced/main.go)

Set your [Personal Access Token](https://api.ynab.com/#personal-access-tokens) as an environment variable, and execute the program:
```shell
export YNAB_PAT=personal-access-token
go run examples/advanced/main.go
```

Output:
```
Our Net Worth is $272,005.98.

We have $19,881.21 cash in our bank accounts.

We have $2,079.27 saved towards medical expenses.
We have $25,874.76 saved towards retirement.
We have $159,581.51 equity in our home.

We need an additional $5,390.30 to fund this month. At the end of the month, we will have roughly -$1,682.89 remaining. Unplanned purchases are not accounted for in this estimate.

We have $170,918.49 remaining on our mortgage.
We have $33,067.28 in student loans.
We have $18,315.82 in auto loans.
```
</details>