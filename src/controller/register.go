package controller

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/greatdanton/analytics/src/global"
	"github.com/greatdanton/analytics/src/templates"
)

// Register takes care of admin registration
func Register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		displayRegister(w, r, registerMsg{})
	case "POST":
		registerAdmin(w, r)
	default:
		displayRegister(w, r, registerMsg{})
	}
}

type registerMsg struct {
	Username        string
	Email           string
	Password        string
	PasswordConfirm string
	ErrorUsername   string
	ErrorEmail      string
	ErrorPassword   string
}

// displayRegister displays register template
func displayRegister(w http.ResponseWriter, r *http.Request, msg registerMsg) {
	err := templates.Execute(w, "register", msg)
	if err != nil {
		return
	}
}

// registerAdmin is handling post requests on register page
// and admin registration. For now registering users is only possible
// when the application is ran with -env=setup.
//WARNING: -env setup drops the entire database, so use it only
// when you are setting your application for the first time

// NOTE: registration part will change in future
// so this is just a temporary solution, which will be replaced
// with more robust one.
func registerAdmin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("registerAdmin: parseForm:", err)
		return
	}
	username := strings.TrimSpace(r.Form["username"][0])
	email := strings.TrimSpace(r.Form["email"][0])
	password := r.Form["password"][0]
	passwordConfirm := r.Form["password-confirm"][0]

	msg := registerMsg{Username: username, Email: email}

	// if passwords do not match
	if password != passwordConfirm {
		msg.ErrorPassword = "Passwords do not match"
		msg.Password = password
		displayRegister(w, r, msg)
		return
	}
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		fmt.Println("Register error: bycrypt:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// TODO: check if user already exists
	// check if email already exists

	// TODO: check if email is correct

	// Insert new user into database
	//err := model.RegisterUser(username, email, passHash)
	_, err = global.DB.Exec(`INSERT into users(username, email, password, active)
							  values($1, $2, $3, TRUE)`, username, email, passHash)

	if err != nil {
		fmt.Println("Register: can not register new admin:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// new user is inserted, inform user to stop executing software with env=setup flag
	fmt.Fprintf(w, "Admin: %s succesfully registered!\n", username)
	fmt.Fprintf(w, "Please stop the analytics software and start it again without -env flag")
}
