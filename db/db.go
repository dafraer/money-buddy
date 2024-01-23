package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func (u *User) UpdateUserData() {
	database, err := sql.Open("sqlite3", "./users.db")
	defer database.Close()
	if err != nil {
		log.Println(err.Error())
	}
	database.Exec(fmt.Sprintf("UPDATE users SET balance = %v, username = '%s', currency = '%s' WHERE user_id = %d", u.Balance, u.Username, u.Currency, u.UserId))
	database.Exec(fmt.Sprintf("UPDATE piggy_bank SET balance = %v, target_amount = %v, target_date = '%s' WHERE user_id = %d", u.PiggyBank.Balance, u.PiggyBank.TargetAmount, u.PiggyBank.TargetDate, u.UserId))
}

// Add method adds transactions to the user account. Variable c has to be 1 when money is recieved and -1 when money is lost
func (u *User) Add(t *Transaction, c float64) {
	database, err := sql.Open("sqlite3", "./users.db")
	defer database.Close()
	if err != nil {
		log.Println(err.Error())
	}
	u.Balance += (t.Amount * c)
	u.Transactions = append(u.Transactions, *t)
	database.Exec(fmt.Sprintf("INSERT INTO transactions (transaction_time, amount, category) VALUES('%s', %v, '%s')", t.TransactionTime.Format(time.DateTime), t.Amount, t.Category))
}

func (u *User) Remove(t *Transaction) {
	database, err := sql.Open("sqlite3", "./users.db")
	defer database.Close()
	if err != nil {
		log.Println(err.Error())
	}
	u.Transactions = append(u.Transactions[:t.TransactionId-1], u.Transactions[t.TransactionId:]...)
	database.Query(fmt.Sprintf("DELETE FROM transactions WHERE transaction_id = %d", t.TransactionId))
}

func CreateNewUser(username string, password string) {
	database, err := sql.Open("sqlite3", "./users.db")
	defer database.Close()
	if err != nil {
		log.Println(err.Error())
	}
	database.Exec(fmt.Sprintf("INSERT INTO users (username, password) VALUES('%s', '%s')", username, password))
	rows, err := database.Query(fmt.Sprintf("SELECT user_id FROM users WHERE username = '%s'", username))
	if err != nil {
		log.Println(err.Error())
	}
	var user_id int
	for rows.Next() {
		rows.Scan(&user_id)
	}
	database.Exec(fmt.Sprintf("INSERT INTO piggy_bank (user_id) VALUES(%d)", user_id))
}

func Exists(username string) bool {
	database, err := sql.Open("sqlite3", "./users.db")
	defer database.Close()
	if err != nil {
		log.Println(err.Error())
	}
	rows, err := database.Query("SELECT username FROM users")
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
	}
	var usernameDB string
	for rows.Next() {
		rows.Scan(&usernameDB)
		if usernameDB == username {
			return true
		}
	}
	return false
}

func Authentication(username string, password string) bool {
	database, err := sql.Open("sqlite3", "./users.db")
	defer database.Close()
	if err != nil {
		log.Println(err.Error())
	}
	rows, err := database.Query("SELECT username, password FROM users")
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
	}
	var usernameDB, passwordDB string
	for rows.Next() {
		rows.Scan(&usernameDB, &passwordDB)
		if usernameDB == username {
			err := bcrypt.CompareHashAndPassword([]byte(passwordDB), []byte(password))
			if err == nil {
				return true
			}
			break
		}
	}
	return false
}

func GetUserData(username string) User {
	var u User
	database, err := sql.Open("sqlite3", "./users.db")
	defer database.Close()
	if err != nil {
		log.Println(err.Error())
	}
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
	rows, err = database.Query(fmt.Sprintf("SELECT transaction_id, transaction_time, amount, category FROM transactions WHERE user_id = %d", u.UserId))
	u.Transactions = make([]Transaction, 0)
	var transaction_id int
	var transaction_time, category string
	var amount float64
	for rows.Next() {
		rows.Scan(&transaction_id, &transaction_time, &amount, &category)
		t, err := time.Parse(time.DateTime, transaction_time)
		if err != nil {
			log.Println(err.Error())
		}
		u.Transactions = append(u.Transactions, Transaction{TransactionId: transaction_id, TransactionTime: t, Amount: amount, Category: category})
	}
	return u
}
