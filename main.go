package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"
	"unicode"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username     string
	Balance      float64
	Transactions []Transaction
}

type Transaction struct {
	Amount   float64
	Time     time.Time
	Category string
}

var store = sessions.NewCookieStore([]byte("super-secret"))
var current User

func main() {
	handleRequest()
}

func handleRequest() {
	http.HandleFunc("/main", homePage)
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.HandleFunc("/registerauth", registerAuth)
	http.HandleFunc("/goals", financialGoals)
	http.HandleFunc("/expenses", expenseTracking)
	http.HandleFunc("/loginauth", loginAuth)
	http.HandleFunc("/analytics", expenseAnalytics)
	http.HandleFunc("/logout", logout)
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
		file, err := os.Open(fmt.Sprintf("data/%s.json", username))
		defer file.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
		content, err := io.ReadAll(file)
		if err != nil {
			fmt.Println(err.Error())
		}
		json.Unmarshal(content, &current)
		if err != nil {
			fmt.Println(err.Error())
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
	statement, err := database.Prepare("INSERT INTO users (username, password) VALUES(?, ?)")
	defer statement.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	statement.Exec(username, string(hash))
	tmpl, err := template.ParseFiles("templates/registrationsuccess.html")
	if err != nil {
		fmt.Println(err.Error())
	}
	//Creating and saving json file for user transactions data
	newUser := User{}
	newUser.Username = username
	j, err := json.Marshal(newUser)
	if err != nil {
		fmt.Println(err.Error())
	}
	file, err := os.Create(fmt.Sprintf("data/%s.json", username))
	if err != nil {
		fmt.Println(err.Error())
	}
	defer file.Close()
	file.Write(j)
	tmpl.Execute(w, nil)
}

func financialGoals(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username, ok := session.Values["username"]
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
	tmpl.Execute(w, username)
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
	username, ok := session.Values["username"]
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
	tmpl.Execute(w, username)
}

func logout(w http.ResponseWriter, r *http.Request) {
	//Saving user data
	j, err := json.Marshal(current)
	if err != nil {
		fmt.Println(err.Error())
	}
	file, err := os.Open(fmt.Sprintf("data/%s.json", current.Username))
	if err != nil {
		fmt.Println(err.Error())
	}
	defer file.Close()
	file.Write(j)
	//Deleting session
	session, _ := store.Get(r, "session")
	delete(session.Values, "username")
	session.Save(r, w)
	tmpl, err := template.ParseFiles("templates/logout.html")
	if err != nil {
		fmt.Println(err.Error())
	}
	tmpl.Execute(w, nil)
}
