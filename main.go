package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"
	"unicode"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

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
}

var store = sessions.NewCookieStore([]byte("super-secret"))
var current User

func main() {
	http.HandleFunc("/main", homePage)
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.HandleFunc("/registerauth", registerAuth)
	http.HandleFunc("/goals", financialGoals)
	http.HandleFunc("/expenses", expenseTracking)
	http.HandleFunc("/loginauth", loginAuth)
	http.HandleFunc("/analytics", expenseAnalytics)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/goalssave", goalsSave)
	http.ListenAndServe(":8000", context.ClearHandler(http.DefaultServeMux))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	_, ok := session.Values["username"]
	if !ok {
		tmpl, err := template.ParseFiles("templates/homepage.html")
		if err != nil {
			fmt.Println(err.Error())
		}
		tmpl.Execute(w, nil)
		return
	}
	tmpl, err := template.ParseFiles("templates/homepageacc.html")
	if err != nil {
		fmt.Println(err.Error())
	}
	tmpl.Execute(w, current)
}

func login(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		fmt.Println(err.Error())
	}
	tmpl.Execute(w, nil)
}

func loginAuth(w http.ResponseWriter, r *http.Request) {
	database, err := sql.Open("sqlite3", "./users.db")
	defer database.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	//Check the validity of username and password
	rows, err := database.Query("SELECT username, password FROM users")
	defer rows.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	var usernameDB, passwordDB string
	correct := false
	for rows.Next() {
		rows.Scan(&usernameDB, &passwordDB)
		if usernameDB == username {
			err := bcrypt.CompareHashAndPassword([]byte(passwordDB), []byte(password))
			if err == nil {
				correct = true
			}
			break
		}
	}
	if correct == true {
		//Opening current user data
		rows, err := database.Query(fmt.Sprintf("SELECT users.user_id, username, target_amount, target_date, currency, users.balance, piggy_bank.balance FROM users  JOIN piggy_bank ON users.user_id = piggy_bank.user_id WHERE username = '%s'", username))
		if err != nil {
			fmt.Println(err.Error())
		}
		for rows.Next() {
			var t string
			rows.Scan(&current.UserId, &current.Username, &current.PiggyBank.TargetAmount, &t, &current.Currency, &current.Balance, &current.PiggyBank.Balance)
			if t == "" {
				current.PiggyBank.TargetDate = convertToStringDate(time.Now())
			} else {
				current.PiggyBank.TargetDate = t
			}

		}
		rows, err = database.Query(fmt.Sprintf("SELECT transaction_id, transaction_time, amount, category FROM transactions WHERE user_id = %d", current.UserId))
		current.Transactions = make([]Transaction, 0)
		var transaction_id int
		var transaction_time, category string
		var amount float64
		for rows.Next() {
			rows.Scan(&transaction_id, &transaction_time, &amount, &category)
			current.Transactions = append(current.Transactions, Transaction{transaction_id, convertToTime(transaction_time), amount, category})
		}
		//Creating login session
		session, err := store.Get(r, "session")
		if err != nil {
			fmt.Println(err.Error())
		}
		session.Values["username"] = username
		session.Save(r, w)
		tmpl, err := template.ParseFiles("templates/loginsuccess.html")
		if err != nil {
			fmt.Println(err.Error())
		}
		tmpl.Execute(w, nil)
	} else {
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			fmt.Println(err.Error())
		}
		tmpl.Execute(w, "Incorrect username or password. Please try again.")
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/registration.html")
	if err != nil {
		fmt.Println(err.Error())
	}
	tmpl.Execute(w, nil)
}

func registerAuth(w http.ResponseWriter, r *http.Request) {
	database, err := sql.Open("sqlite3", "./users.db")
	defer database.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	//Username has to  have no spaces and only ASCII characters
	var spacesInUsername, spacesInPassword, notASCIIPassowrd, notASCIIUsername bool
	for _, c := range username {
		if c == ' ' {
			spacesInUsername = true
		}
		if c > unicode.MaxASCII {
			notASCIIUsername = true
		}
	}
	//Password has to have no spaces and only ASCII characters
	for _, c := range password {
		if c == ' ' {
			spacesInPassword = true
		}
		if c > unicode.MaxASCII {
			notASCIIPassowrd = true
		}
	}
	//Checking if the password is at least 8 characters long and fits password requierments
	//Checking if username is from 1 to 20 characters long and fits username requiermants
	if len(username) > 20 || len(username) < 1 || len(password) < 8 || spacesInUsername || spacesInPassword || notASCIIPassowrd || notASCIIUsername {
		tmpl, err := template.ParseFiles("templates/registration.html")
		if err != nil {
			fmt.Println(err.Error())
		}
		tmpl.Execute(w, "Password or username does not meet the requirements.")
		return
	}
	//Checking if user already exists
	rows, err := database.Query("SELECT username FROM users")
	defer rows.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	var usernameDB string
	exists := false
	for rows.Next() {
		rows.Scan(&usernameDB)
		if usernameDB == username {
			exists = true
		}
	}
	if exists == true {
		tmpl, err := template.ParseFiles("templates/registration.html")
		if err != nil {
			fmt.Println(err.Error())
		}
		tmpl.Execute(w, "Account already exists")
		return
	}
	//Hashing the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err.Error())
		tmpl, err := template.ParseFiles("templates/registration.html")
		if err != nil {
			fmt.Println(err.Error())
		}
		tmpl.Execute(w, "There was a problem registering new user")
		return
	}
	//Creating a user
	database.Exec(fmt.Sprintf("INSERT INTO users (username, password) VALUES('%s', '%s')", username, string(hash)))
	rows, err = database.Query(fmt.Sprintf("SELECT user_id FROM users WHERE username = '%s'", username))
	if err != nil {
		fmt.Println(err.Error())
	}
	var user_id int
	for rows.Next() {
		rows.Scan(&user_id)
	}
	database.Exec(fmt.Sprintf("INSERT INTO piggy_bank (user_id) VALUES(%d)", user_id))
	tmpl, err := template.ParseFiles("templates/registrationsuccess.html")
	if err != nil {
		fmt.Println(err.Error())
	}
	tmpl.Execute(w, nil)
}

func financialGoals(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	_, ok := session.Values["username"]
	if !ok {
		tmpl, err := template.ParseFiles("templates/goals.html")
		if err != nil {
			fmt.Println(err.Error())
		}
		tmpl.Execute(w, nil)
		return
	}
	tmpl, err := template.ParseFiles("templates/goalsacc.html")
	if err != nil {
		fmt.Println(err.Error())
	}
	tmpl.Execute(w, current)
}

func goalsSave(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if amount != 0 && err == nil {
		current.PiggyBank.Balance += amount
		var t Transaction
		t.TransactionTime = time.Now()
		t.Amount = amount
		t.Category = "Savings"
		t.Add(-1)
	}
	newAmount, err := strconv.ParseFloat(r.FormValue("newAmount"), 64)
	if newAmount != 0 && err == nil {
		current.PiggyBank.TargetAmount = newAmount
	}
	newDateString := r.FormValue("newDate")
	_, err = time.Parse(time.DateOnly, r.FormValue("newDate"))
	if err == nil {
		current.PiggyBank.TargetDate = newDateString
	}
	current.updateUserData()
	http.Redirect(w, r, "/goals", http.StatusSeeOther)
}

func expenseTracking(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username, ok := session.Values["username"]
	if !ok {
		tmpl, err := template.ParseFiles("templates/expenses.html")
		if err != nil {
			fmt.Println(err.Error())
		}
		tmpl.Execute(w, nil)
		return
	}
	tmpl, err := template.ParseFiles("templates/expensesacc.html")
	if err != nil {
		fmt.Println(err.Error())
	}
	tmpl.Execute(w, username)
}

func expenseAnalytics(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	_, ok := session.Values["username"]
	if !ok {
		tmpl, err := template.ParseFiles("templates/analytics.html")
		if err != nil {
			fmt.Println(err.Error())
		}
		tmpl.Execute(w, nil)
		return
	}
	tmpl, err := template.ParseFiles("templates/analyticsacc.html")
	if err != nil {
		fmt.Println(err.Error())
	}
	var temp Analytics
	temp.Username = current.Username
	for _, v := range current.Transactions {
		if v.Amount > 0 {
			temp.Income++
		} else {
			temp.Expenditure++
		}
	}
	tmpl.Execute(w, temp)
}

func logout(w http.ResponseWriter, r *http.Request) {
	//Saving user data
	current.updateUserData()
	session, _ := store.Get(r, "session")
	//Deleting session
	delete(session.Values, "username")
	session.Save(r, w)
	tmpl, err := template.ParseFiles("templates/logout.html")
	if err != nil {
		fmt.Println(err.Error())
	}
	tmpl.Execute(w, nil)
}

func (u *User) updateUserData() {
	database, err := sql.Open("sqlite3", "./users.db")
	defer database.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	database.Exec(fmt.Sprintf("UPDATE users SET balance = %v, username = '%s', currency = '%s' WHERE user_id = %d", u.Balance, u.Username, u.Currency, u.UserId))
	database.Exec(fmt.Sprintf("UPDATE piggy_bank SET balance = %v, target_amount = %v, target_date = '%s' WHERE user_id = %d", u.PiggyBank.Balance, u.PiggyBank.TargetAmount, u.PiggyBank.TargetDate, u.UserId))
}

// Add method adds transactions to the user account. Variable c has to be 1 when money is recieved and -1 when money is lost
func (t *Transaction) Add(c float64) {
	database, err := sql.Open("sqlite3", "./users.db")
	defer database.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	current.Balance += (t.Amount * c)
	current.Transactions = append(current.Transactions, *t)
	database.Exec(fmt.Sprintf("INSERT INTO transactions (transaction_time, amount, category) VALUES('%s', %v, '%s')", convertToStringTime(t.TransactionTime), t.Amount, t.Category))
}

func (t *Transaction) Remove() {
	database, err := sql.Open("sqlite3", "./users.db")
	defer database.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	current.Transactions = append(current.Transactions[:t.TransactionId-1], current.Transactions[t.TransactionId:]...)
	database.Query(fmt.Sprintf("DELETE FROM transactions WHERE transaction_id = %d", t.TransactionId))
}

func convertToStringTime(t time.Time) string {
	s := t.Format(time.DateTime)
	return s
}

func convertToStringDate(t time.Time) string {
	s := t.Format(time.DateOnly)
	return s
}

func convertToTime(s string) time.Time {
	t, err := time.Parse(time.DateTime, s)
	if err != nil {
		fmt.Println(err.Error())
	}
	return t
}
