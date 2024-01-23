package db

import "time"

type User struct {
	UserId       int
	Username     string
	Currency     string
	Balance      float64
	Transactions []Transaction
	PiggyBank    PiggyBank
}

type PiggyBank struct {
	TargetAmount float64
	TargetDate   string
	Balance      float64
}

type Transaction struct {
	TransactionId int
	//Storing UserId seems useless, mb delete later
	//UserId          int
	TransactionTime time.Time
	Amount          float64
	Category        string
}

type Analytics struct {
	Username          string
	Income            float64
	Expenditure       float64
	CategoriesNames   [5]string
	CategoriesExpense [5]float64
}
