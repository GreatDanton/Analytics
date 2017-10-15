package controller

import (
	"fmt"
	"net/http"

	"github.com/greatdanton/analytics/src/sessions"
	"github.com/greatdanton/analytics/src/templates"
)

// Dashboard takes care of displaying main user dashboard
func Dashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// if user is not logged in redirect him to login page
		user := sessions.LoggedInUser(r)
		if !user.LoggedIn {
			fmt.Println("User is not logged in")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		err := templates.Execute(w, "mainDashboard", nil)
		if err != nil {
			fmt.Println(err)
			return
		}

	}
}
