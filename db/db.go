package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// UpdateUserData updates user data in the database
func (u *User) UpdateUserData() {

	//Opening the database
	database, err := sql.Open("sqlite3", "./users.db")
	defer database.Close()
	if err != nil {
		log.Println(err.Error())
	}

	//Updating user's data in the database
	database.Exec(fmt.Sprintf("UPDATE users SET balance = %v, username = '%s', currency = '%s' WHERE user_id = %d", u.Balance, u.Username, u.Currency, u.UserId))
	database.Exec(fmt.Sprintf("UPDATE piggy_bank SET balance = %v, target_amount = %v, target_date = '%s' WHERE user_id = %d", u.PiggyBank.Balance, u.PiggyBank.TargetAmount, u.PiggyBank.TargetDate, u.UserId))
}

// Dec adds expense transaction
func (u *User) Dec(t Transaction) {
	t.Amount *= -1
	u.Add(&t)
}

// Add method adds transactions to the account
func (u *User) Add(t *Transaction) {

	//Opening database
	database, err := sql.Open("sqlite3", "./users.db")
	defer database.Close()
	if err != nil {
		log.Println(err.Error())
	}

	//Adding transaction to the main balance
	u.Balance += t.Amount

	//Adding transaction to the transactions table in the database
	//Transactions are added to the beginning of the slice in order for the added transaction to be displayed at the top
	temp := make([]Transaction, 1)
	temp[0] = *t
	u.Transactions = append(temp, u.Transactions...)
	database.Exec(fmt.Sprintf("INSERT INTO transactions (user_id, transaction_time, amount, category) VALUES('%d', '%s', %v, '%s')", u.UserId, t.TransactionTime.Format(time.DateTime), t.Amount, t.Category))
	u.Analytics = getAnalyticsData(u.Username)
}

// CreateNewUser adds new user to the database
func CreateNewUser(username string, password string) {

	//Opening database
	database, err := sql.Open("sqlite3", "./users.db")
	defer database.Close()
	if err != nil {
		log.Println(err.Error())
	}

	//Adding new user to the users table in the database
	database.Exec(fmt.Sprintf("INSERT INTO users (username, password) VALUES('%s', '%s')", username, password))

	//Getting user_id from the database
	rows, err := database.Query(fmt.Sprintf("SELECT user_id FROM users WHERE username = '%s'", username))
	if err != nil {
		log.Println(err.Error())
	}
	var userId int
	for rows.Next() {
		rows.Scan(&userId)
	}

	//Creating new PiggyBank in piggy_bank table
	database.Exec(fmt.Sprintf("INSERT INTO piggy_bank (user_id) VALUES(%d)", userId))
}

// Exists checks if user exists in the database
func Exists(username string) bool {

	//Opening the database
	database, err := sql.Open("sqlite3", "./users.db")
	defer database.Close()
	if err != nil {
		log.Println(err.Error())
	}

	//Getting all userenames from the database
	rows, err := database.Query("SELECT username FROM users")
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
	}

	//Checking if the same username already exists
	var usernameDB string
	for rows.Next() {
		rows.Scan(&usernameDB)
		if usernameDB == username {
			return true
		}
	}
	return false
}

// Authentication checks if username and password are valid
func Authentication(username string, password string) bool {

	//Opening the database
	database, err := sql.Open("sqlite3", "./users.db")
	defer database.Close()
	if err != nil {
		log.Println(err.Error())
	}

	//Selecting usernames and password from users table in the database
	rows, err := database.Query("SELECT username, password FROM users")
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
	}

	//Checking validity of the username and the password
	var usernameDB, passwordDB string
	for rows.Next() {
		rows.Scan(&usernameDB, &passwordDB)
		if usernameDB == username {
			//Since passwords are not stored directly password hash is compared
			err := bcrypt.CompareHashAndPassword([]byte(passwordDB), []byte(password))
			if err == nil {
				return true
			}
			break
		}
	}
	return false
}

// GetUserData returns user data from the database
func GetUserData(username string) User {

	//Initializing new user
	var u User

	//Opening the database
	database, err := sql.Open("sqlite3", "./users.db")
	defer database.Close()
	if err != nil {
		log.Println(err.Error())
	}

	//Get user data from the database
	rows, err := database.Query(fmt.Sprintf("SELECT users.user_id, username, target_amount, target_date, currency, users.balance, piggy_bank.balance FROM users  JOIN piggy_bank ON users.user_id = piggy_bank.user_id WHERE username = '%s'", username))
	if err != nil {
		log.Println(err.Error())
	}
	for rows.Next() {
		var t string
		rows.Scan(&u.UserId, &u.Username, &u.PiggyBank.TargetAmount, &t, &u.Currency, &u.Balance, &u.PiggyBank.Balance)
		if t == "" {
			u.PiggyBank.TargetDate = time.Now().Format(time.DateOnly)
		} else {
			u.PiggyBank.TargetDate = t
		}

	}

	//Getting transactions from database
	rows, err = database.Query(fmt.Sprintf("SELECT transaction_id, transaction_time, amount, category FROM transactions WHERE user_id = %d ORDER BY transaction_time DESC", u.UserId))
	u.Transactions = make([]Transaction, 0)
	var transactionId int
	var transactionTime, category string
	var amount float64
	for rows.Next() {
		rows.Scan(&transactionId, &transactionTime, &amount, &category)
		t, err := time.Parse(time.DateTime, transactionTime)
		if err != nil {
			log.Println(err.Error())
		}
		u.Transactions = append(u.Transactions, Transaction{TransactionId: transactionId, TransactionTime: t, Amount: amount, Category: category})
	}
	u.Analytics = getAnalyticsData(u.Username)
	return u
}

// GetAnalyticsData returns analysed user data
func getAnalyticsData(username string) Analytics {

	//Initializing analytics variable
	var a Analytics
	a.Username = username

	//Opening the database
	database, err := sql.Open("sqlite3", "./users.db")
	defer database.Close()
	if err != nil {
		log.Println(err.Error())
	}

	//Getting total income and expense
	rows, err := database.Query(fmt.Sprintf("SELECT amount FROM transactions JOIN users ON users.user_id = transactions.user_id WHERE username = '%s'", username))
	defer rows.Close()
	var t float64
	for rows.Next() {
		rows.Scan(&t)
		if t > 0 {
			a.Income += t
		} else {
			a.Expenditure += t
		}
	}

	//Making a.Expenditure equal to 0 in case it is -0
	if a.Expenditure != 0 {
		a.Expenditure *= -1
	}

	//Calculating by category
	rows, err = database.Query(fmt.Sprintf("SELECT category, SUM(amount) AS total_amount FROM users  LEFT JOIN transactions ON users.user_id = transactions.user_id  WHERE username = '%s' AND amount < 0  GROUP BY category  ORDER BY total_amount;", username))
	if err != nil {
		log.Println(err.Error())
	}
	i := 0
	for rows.Next() {
		//Adding first 4 categories into piechart
		//The rest categories are combined and displayed as "Other"
		if i <= 4 {
			rows.Scan(&a.Categories[i].Name, &a.Categories[i].Amount)
			a.Categories[i].Amount *= -1
			i++
		} else {
			var tempString string
			var tempFloat float64
			rows.Scan(&tempString, &tempFloat)
			a.Categories[4].Amount += tempFloat * -1
		}
	}
	return a
}
