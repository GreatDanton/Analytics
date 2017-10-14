package controller

import (
	"fmt"
	"net/http"

	"github.com/greatdanton/analytics/src/templates"
)

// Login takes care of handling get and post request of the login part
// of the application
func Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		displayLogin(w, r, loginMsg{})
	case "POST":
		checkLogin(w, r)
	default:
		displayLogin(w, r, loginMsg{})
	}
}

type loginMsg struct {
	ErrorPassword string
}

// displayLogin takes care of rendering login template and displaying
// error messages when user posted wrong credentials
func displayLogin(w http.ResponseWriter, r *http.Request, msg loginMsg) {
	err := templates.Execute(w, "login", msg)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// checks login post request
func checkLogin(w http.ResponseWriter, r *http.Request) {

}
