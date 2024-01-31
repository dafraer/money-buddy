package handler

import (
	"MoneyBuddy/db"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"
	"time"
	"unicode"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var store *sessions.CookieStore
var current db.User

func HandleRequest() {
	store = sessions.NewCookieStore([]byte("super-secret"))
	current = db.User{}
	http.HandleFunc("/images/", imageHandler)
	http.HandleFunc("/support/", headerFooterHandler)
	http.HandleFunc("/support/css/", cssHandler)
	http.HandleFunc("/support/js/", jsHandler)
	http.HandleFunc("/main", homePageHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/registerauth", registerAuthHandler)
	http.HandleFunc("/goals", financialGoalsHandler)
	http.HandleFunc("/expenses", expenseTrackingHandler)
	http.HandleFunc("/loginauth", loginAuthHandler)
	http.HandleFunc("/analytics", expenseAnalyticsHandler)
	http.HandleFunc("/getuserdata", getUserDataHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/goalssave", financialGoalsSaveHandler)
	http.HandleFunc("/postuserdata", updateUserDataHandler)
	http.HandleFunc("/addtransaction", addTransactionHandler)
	http.ListenAndServe(":8000", context.ClearHandler(http.DefaultServeMux))
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path
	filePath = path.Base(filePath)
	fullPath := path.Join("templates/images/", filePath)
	http.ServeFile(w, r, fullPath)
}

func jsHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path
	filePath = path.Base(filePath)
	fullPath := path.Join("templates/support/js/", filePath)
	http.ServeFile(w, r, fullPath)
}

func cssHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path
	filePath = path.Base(filePath)
	fullPath := path.Join("templates/support/css/", filePath)
	http.ServeFile(w, r, fullPath)
}

func headerFooterHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path
	filePath = path.Base(filePath)
	fullPath := path.Join("templates/support/", filePath)
	http.ServeFile(w, r, fullPath)
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	_, ok := session.Values["username"]
	if !ok {
		tmpl, err := template.ParseFiles("templates/homepage.html")
		if err != nil {
			log.Println(err.Error())
		}
		tmpl.Execute(w, nil)
		return
	}
	tmpl, err := template.ParseFiles("templates/homepageacc.html")
	if err != nil {
		log.Println(err.Error())
	}
	tmpl.Execute(w, current)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		log.Println(err.Error())
	}
	tmpl.Execute(w, nil)
}

func loginAuthHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	//Check the validity of username and password
	correct := db.Authentication(username, password)
	if correct == true {
		//Opening current user data
		current = db.GetUserData(username)
		current.Analytics = db.GetAnalyticsData(username)
		//Creating login session
		session, err := store.Get(r, "session")
		if err != nil {
			log.Println(err.Error())
		}
		session.Values["username"] = username
		session.Save(r, w)
		tmpl, err := template.ParseFiles("templates/loginsuccess.html")
		if err != nil {
			log.Println(err.Error())
		}
		tmpl.Execute(w, nil)
	} else {
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			log.Println(err.Error())
		}
		tmpl.Execute(w, "Incorrect username or password. Please try again.")
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/registration.html")
	if err != nil {
		log.Println(err.Error())
	}
	tmpl.Execute(w, nil)
}

func registerAuthHandler(w http.ResponseWriter, r *http.Request) {
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
			log.Println(err.Error())
		}
		tmpl.Execute(w, "Password or username does not meet the requirements.")
		return
	}
	//Checking if user already exists
	exists := db.Exists(username)
	if exists {
		tmpl, err := template.ParseFiles("templates/registration.html")
		if err != nil {
			log.Println(err.Error())
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
			log.Println(err.Error())
		}
		tmpl.Execute(w, "There was a problem registering new user")
		return
	}
	//Creating a user
	db.CreateNewUser(username, string(hash))
	tmpl, err := template.ParseFiles("templates/registrationsuccess.html")
	if err != nil {
		log.Println(err.Error())
	}
	tmpl.Execute(w, nil)
}

func financialGoalsHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	_, ok := session.Values["username"]
	if !ok {
		tmpl, err := template.ParseFiles("templates/goals.html")
		if err != nil {
			log.Println(err.Error())
		}
		tmpl.Execute(w, nil)
		return
	}
	tmpl, err := template.ParseFiles("templates/goalsacc.html")
	if err != nil {
		log.Println(err.Error())
	}
	tmpl.Execute(w, current)
}

func financialGoalsSaveHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if amount != 0 && err == nil {
		current.PiggyBank.Balance += amount
		var t db.Transaction
		t.TransactionTime = time.Now()
		t.Amount = amount
		t.Category = "Savings"
		current.Add(&t, -1)
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
	current.Analytics = db.GetAnalyticsData(current.Username)
	current.UpdateUserData()
	http.Redirect(w, r, "/goals", http.StatusSeeOther)
}

func expenseTrackingHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username, ok := session.Values["username"]
	if !ok {
		tmpl, err := template.ParseFiles("templates/expenses.html")
		if err != nil {
			log.Println(err.Error())
		}
		tmpl.Execute(w, nil)
		return
	}
	tmpl, err := template.ParseFiles("templates/expensesacc.html")
	if err != nil {
		log.Println(err.Error())
	}
	tmpl.Execute(w, username)
}

func expenseAnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	_, ok := session.Values["username"]
	if !ok {
		tmpl, err := template.ParseFiles("templates/analytics.html")
		if err != nil {
			log.Println(err.Error())
		}
		tmpl.Execute(w, nil)
		return
	}
	tmpl, err := template.ParseFiles("templates/analyticsacc.html")
	if err != nil {
		log.Println(err.Error())
	}
	tmpl.Execute(w, nil)
}

func getUserDataHandler(w http.ResponseWriter, r *http.Request) {
	jsonData, err := json.Marshal(current)
	if err != nil {
		log.Println(err.Error())
	}
	w.Write(jsonData)
}

func updateUserDataHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&current)
	current.UpdateUserData()
	if err != nil {
		log.Println(err.Error())
	}
}

func addTransactionHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	t := db.Transaction{}
	err := decoder.Decode(&t)
	current.Add(&t, 1)
	current.UpdateUserData()
	if err != nil {
		log.Println(err.Error())
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	//Saving user data
	current.UpdateUserData()
	session, _ := store.Get(r, "session")
	//Deleting session
	delete(session.Values, "username")
	session.Save(r, w)
	tmpl, err := template.ParseFiles("templates/logout.html")
	if err != nil {
		log.Println(err.Error())
	}
	tmpl.Execute(w, nil)
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
		log.Println(err.Error())
	}
	return t
}
