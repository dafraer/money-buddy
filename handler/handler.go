package handler

import (
	"MoneyBuddy/db"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var store *sessions.CookieStore

func HandleRequest() {

	//Creating new cookies session
	store = sessions.NewCookieStore([]byte("super-secret"))

	//Loading all the pages
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
	http.HandleFunc("/postpiggybank", updatePiggyBankHandler)
	http.HandleFunc("/addtransaction", addTransactionHandler)

	//Get port from the environment if possible
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	port = ":" + port
	http.ListenAndServe(port, context.ClearHandler(http.DefaultServeMux))

}

// imageHandler serves images
func imageHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path
	filePath = path.Base(filePath)
	fullPath := path.Join("templates/images/", filePath)
	http.ServeFile(w, r, fullPath)
}

// jsHandler serves javascript files
func jsHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path
	filePath = path.Base(filePath)
	fullPath := path.Join("templates/support/js/", filePath)
	http.ServeFile(w, r, fullPath)
}

// cssHandler serves css files
func cssHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path
	filePath = path.Base(filePath)
	fullPath := path.Join("templates/support/css/", filePath)
	http.ServeFile(w, r, fullPath)
}

// headerFooter serves footer and header
func headerFooterHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path
	filePath = path.Base(filePath)
	fullPath := path.Join("templates/support/", filePath)
	http.ServeFile(w, r, fullPath)
}

// homePageHandler handles home page
func homePageHandler(w http.ResponseWriter, r *http.Request) {
	//Checking if user is logged in
	session, _ := store.Get(r, "session")
	sessionId, ok := session.Values["sessionId"]

	//If user is not logged in execute regular home page
	if !ok {
		tmpl, err := template.ParseFiles("templates/homepage.html")
		if err != nil {
			log.Println(err.Error())
		}
		tmpl.Execute(w, nil)
		return
	}

	//Else execute template with user account details
	username := db.GetUsername(sessionId.(string))
	tmpl, err := template.ParseFiles("templates/homepageacc.html")
	if err != nil {
		log.Println(err.Error())
	}
	current := db.GetUserData(username)
	tmpl.Execute(w, current)
}

// loginHandler provides login form for the user
func loginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		log.Println(err.Error())
	}
	tmpl.Execute(w, nil)
}

// loginAuthHandler authenticates user
func loginAuthHandler(w http.ResponseWriter, r *http.Request) {

	//Parsing login form
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")

	//Check the validity of username and password
	correct := db.Authentication(username, password)
	//if data is valid proceed
	if correct == true {
		//Creating login session
		session, err := store.Get(r, "session")
		if err != nil {
			log.Println(err.Error())
		}
		token := uuid.New()
		sessionId := token.String()
		session.Values["sessionId"] = sessionId
		db.AddToken(sessionId, username)
		session.Save(r, w)
		http.Redirect(w, r, "/main", http.StatusSeeOther)

		log.Println(fmt.Sprintf("User %s logged in", username))
	} else {
		//Else notify user that data is not valid
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			log.Println(err.Error())
		}
		tmpl.Execute(w, "Incorrect username or password. Please try again.")
	}
}

// registerHandler provides registration from for the user
func registerHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/registration.html")
	if err != nil {
		log.Println(err.Error())
	}
	tmpl.Execute(w, nil)
}

// registerAuthHandler registers new user
func registerAuthHandler(w http.ResponseWriter, r *http.Request) {

	//Parsing data from the form
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")

	//Checking that username has no spaces and only ASCII characters
	var spacesInUsername, spacesInPassword, notASCIIPassowrd, notASCIIUsername bool
	for _, c := range username {
		if c == ' ' {
			spacesInUsername = true
		}
		if c > unicode.MaxASCII {
			notASCIIUsername = true
		}
	}

	//Checking if the password has no spaces and only ASCII characters
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

	log.Println(fmt.Sprintf("Created New User: %s", username))
}

// FinancialGoalsHandler handles PiggyBank page
func financialGoalsHandler(w http.ResponseWriter, r *http.Request) {

	//Checking if user is logged in
	session, _ := store.Get(r, "session")
	//If user is not logged in execute template that tells them to log in
	_, ok := session.Values["sessionId"]
	if !ok {
		tmpl, err := template.ParseFiles("templates/goals.html")
		if err != nil {
			log.Println(err.Error())
		}
		tmpl.Execute(w, nil)
		return
	}

	//Else execute template with user's PiggyBank
	tmpl, err := template.ParseFiles("templates/goalsacc.html")
	if err != nil {
		log.Println(err.Error())
	}
	tmpl.Execute(w, nil)
}

// expenseTrackingHandler handles expense tracking page
func expenseTrackingHandler(w http.ResponseWriter, r *http.Request) {
	//Check if user is logged in
	session, _ := store.Get(r, "session")
	sessionId, ok := session.Values["sessionId"]
	//If user is not logged in execute template that tells them to log in
	if !ok {
		tmpl, err := template.ParseFiles("templates/expenses.html")
		if err != nil {
			log.Println(err.Error())
		}
		tmpl.Execute(w, nil)
		return
	}
	//Else execute template with users transactions
	username := db.GetUsername(sessionId.(string))
	tmpl, err := template.ParseFiles("templates/expensesacc.html")
	if err != nil {
		log.Println(err.Error())
	}
	tmpl.Execute(w, username)

}

// expenseAnalyticsHandler handles expense analytics page
func expenseAnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	//Check if user is logged in
	session, _ := store.Get(r, "session")
	_, ok := session.Values["sessionId"]
	//If user is not logged in execute template that tells them to log in
	if !ok {
		tmpl, err := template.ParseFiles("templates/analytics.html")
		if err != nil {
			log.Println(err.Error())
		}
		tmpl.Execute(w, nil)
		return
	}
	//Else execute template with expense analytics data
	tmpl, err := template.ParseFiles("templates/analyticsacc.html")
	if err != nil {
		log.Println(err.Error())
	}
	tmpl.Execute(w, nil)
}

// getUserDataHandler handles an HTTP request and returns User data encoded in json
func getUserDataHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	sessionId, ok := session.Values["sessionId"]
	if ok {
		username := db.GetUsername(sessionId.(string))
		current := db.GetUserData(username)

		//Encode current user in json
		jsonData, err := json.Marshal(current)
		if err != nil {
			log.Println(err.Error())
		}

		//return json encoded user data
		w.Write(jsonData)

		//log request
		log.Println(fmt.Sprintf("User %s requested account data", current.Username))
	}
	//log request
	log.Println("Unauthorised user tried to request data")
	return
}

// updatePiggyBankHandler handles an HTTP request to update data in the PiggyBank
func updatePiggyBankHandler(w http.ResponseWriter, r *http.Request) {
	//Unmarshall new PiggyBank data into a json
	decoder := json.NewDecoder(r.Body)
	p := db.PiggyBank{}
	if err := decoder.Decode(&p); err != nil {
		log.Println(err.Error())
	}

	//Check if the user is logged in
	session, _ := store.Get(r, "session")
	sessionId, ok := session.Values["sessionId"]
	if ok {
		username := db.GetUsername(sessionId.(string))
		current := db.GetUserData(username)
		//Adding PiggyBank transaction
		if p.Balance > 0 {
			var t db.Transaction
			t.TransactionTime = time.Now().UTC()
			t.Amount = p.Balance
			t.Category = "Savings"
			current.Dec(t)
		}

		//Updating PiggyBank balance
		p.Balance += current.PiggyBank.Balance
		current.PiggyBank = p

		//Save recieved data
		current.UpdateUserData()

		//log request
		log.Println(fmt.Sprintf("User %s updated PiggyBank. Piggybank: %v", current.Username, p))
	}
	//log request
	log.Println("Unauthorised user tried to update PigyBank")
	return
}

// addTransactionHandler handles HTTP request to add new transaction
func addTransactionHandler(w http.ResponseWriter, r *http.Request) {
	//Unmarshall transaction data from json file
	decoder := json.NewDecoder(r.Body)
	t := db.Transaction{}
	err := decoder.Decode(&t)

	//Check if the user is logged in
	session, _ := store.Get(r, "session")
	sessionId, ok := session.Values["sessionId"]

	if ok {
		//Add the transaction
		username := db.GetUsername(sessionId.(string))
		current := db.GetUserData(username)
		current.Add(&t)
		current.UpdateUserData()
		if err != nil {
			log.Println(err.Error())
		}

		//log request
		log.Println(fmt.Sprintf("User %s added transaction. Amount: %v, Category: %s", current.Username, t.Amount, t.Category))
	}
	//log request
	log.Println("Unauthorised user tried to add transaction")
	return
}

// logoutHandler logs out user by ending cookie session
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	//Deleting session
	session, _ := store.Get(r, "session")
	sessionId, ok := session.Values["sessionId"]
	if ok {
		db.DeleteToken(sessionId.(string))
	}
	delete(session.Values, "sessionId")
	session.Save(r, w)
	tmpl, err := template.ParseFiles("templates/logout.html")
	if err != nil {
		log.Println(err.Error())
	}
	tmpl.Execute(w, nil)
}
