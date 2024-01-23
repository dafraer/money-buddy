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
	Username    string
	Income      float64
	Expenditure float64
	Category1   Category
	Category2   Category
	Category3   Category
	Category4   Category
	Category5   Category
}

type Category struct {
	Name   string
	Amount float64
}
