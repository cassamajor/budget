package budget

import (
	"encoding/json"
	"fmt"
)

// Month contains the overall financial picture for a specific month.
type Month struct {
	ReadyToAssign Balance
	Assigned      Balance
	Underfunded   Balance
	Income        Balance
	Expenses      Balance
}

// UnmarshalJSON unmarshals the JSON data retrieved from YNAB's API into the Month struct.
func (m *Month) UnmarshalJSON(data []byte) error {
	var format MonthData
	var underfunded int

	err := json.Unmarshal(data, &format)
	if err != nil {
		return err
	}

	for _, category := range format.Data.Month.Categories {
		if category.GoalUnderFunded != nil {
			underfunded += *category.GoalUnderFunded
		}
	}

	// Assigning the values into NewRat avoids floating point errors
	m.ReadyToAssign = NewRat(format.Data.Month.ToBeBudgeted)
	m.Assigned = NewRat(format.Data.Month.Budgeted)
	m.Underfunded = NewRat(underfunded)
	m.Income = NewRat(format.Data.Month.Income)
	m.Expenses = NewRat(-format.Data.Month.Activity)

	// Calculate payments made to the Credit Card Payments category as an expense
	for _, category := range format.Data.Month.Categories {
		if category.CategoryGroupName == "Credit Card Payments" {
			m.Expenses = m.Expenses + NewRat(category.Activity)
		}
	}

	return nil
}

// Report prints the current state of finances for a Month.
func (m Month) Report() {
	fmt.Printf("Ready to Assign: %v\n", m.ReadyToAssign)
	fmt.Printf("Assigned: %v\n", m.Assigned)
	fmt.Printf("Underfunded: %v\n", m.Underfunded)
	fmt.Printf("Income: %v\n", m.Income)
	fmt.Printf("Expenses: %v\n", m.Expenses)
}

// MonthData https://api.ynab.com/v1/#/Months/getBudgetMonth
type MonthData struct {
	Data struct {
		Month struct {
			Month        string `json:"month"`
			Note         string `json:"note"`
			Income       int    `json:"income"`
			Budgeted     int    `json:"budgeted"`
			Activity     int    `json:"activity"`
			ToBeBudgeted int    `json:"to_be_budgeted"`
			AgeOfMoney   int    `json:"age_of_money"`
			Deleted      bool   `json:"deleted"`
			Categories   []struct {
				Id                      string      `json:"id"`
				CategoryGroupId         string      `json:"category_group_id"`
				CategoryGroupName       string      `json:"category_group_name"`
				Name                    string      `json:"name"`
				Hidden                  bool        `json:"hidden"`
				OriginalCategoryGroupId interface{} `json:"original_category_group_id"`
				Note                    *string     `json:"note"`
				Budgeted                int         `json:"budgeted"`
				Activity                int         `json:"activity"`
				Balance                 int         `json:"balance"`
				GoalType                *string     `json:"goal_type"`
				GoalDay                 *int        `json:"goal_day"`
				GoalCadence             *int        `json:"goal_cadence"`
				GoalCadenceFrequency    *int        `json:"goal_cadence_frequency"`
				GoalCreationMonth       *string     `json:"goal_creation_month"`
				GoalTarget              int         `json:"goal_target"`
				GoalTargetMonth         *string     `json:"goal_target_month"`
				GoalPercentageComplete  *int        `json:"goal_percentage_complete"`
				GoalMonthsToBudget      *int        `json:"goal_months_to_budget"`
				GoalUnderFunded         *int        `json:"goal_under_funded"`
				GoalOverallFunded       *int        `json:"goal_overall_funded"`
				GoalOverallLeft         *int        `json:"goal_overall_left"`
				Deleted                 bool        `json:"deleted"`
			} `json:"categories"`
		} `json:"month"`
	} `json:"data"`
}
