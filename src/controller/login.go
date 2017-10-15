package controller

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/greatdanton/analytics/src/model"
	"github.com/greatdanton/analytics/src/sessions"
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
	Username      string
	Password      string
	ErrorUsername string
	ErrorPassword string
}

// displayLogin takes care of rendering login template and displaying
// error messages when user posted wrong credentials
func displayLogin(w http.ResponseWriter, r *http.Request, msg loginMsg) {
	// if user is already logged in redirect him to the dashboard page
	user := sessions.LoggedInUser(r)
	if user.LoggedIn {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
	// user is not logged in, display login page
	err := templates.Execute(w, "login", msg)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// checks login post (actual Login) request
func checkLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("checkLogin: parseForm:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	username := strings.TrimSpace(r.Form["username"][0])
	password := r.Form["password"][0]

	errMsg := loginMsg{Username: username, Password: password}

	LoginUser, err := model.GetUserLoginData(username, password)
	if err != nil {
		// if user does not exist
		if err == model.ErrUserDoesNotExist {
			errMsg.ErrorUsername = fmt.Sprintf("%v", model.ErrUserDoesNotExist)
			displayLogin(w, r, errMsg)
			return
		}
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	bytePass := []byte(password)
	err = bcrypt.CompareHashAndPassword(LoginUser.PasswordHash, bytePass)
	if err != nil {
		fmt.Println("checkLogin: Wrong password")
		errMsg.ErrorPassword = "Wrong password"
		errMsg.Password = ""
		displayLogin(w, r, errMsg)
		return
	}
	// password is correct, create user session and log in user
	err = sessions.CreateUserSession(w, LoginUser.ID, LoginUser.Username)
	if err != nil {
		fmt.Println("CreateUserSession:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// Logout logs user out
func Logout(w http.ResponseWriter, r *http.Request) {
	user := sessions.LoggedInUser(r)
	if !user.LoggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err := sessions.DestroyUserSession(w, r)
	if err != nil {
		fmt.Println("Logout: DestroyUserSession:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = templates.Execute(w, "logout", nil)
	if err != nil {
		fmt.Println("Logout:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
