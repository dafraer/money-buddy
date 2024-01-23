package main

import (
	"MoneyBuddy/handler"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("super-secret"))

func main() {
	handler.HandleRequest()
}
