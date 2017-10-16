package controller

import (
	"fmt"
	"net/http"

	"github.com/greatdanton/analytics/src/sessions"
	"github.com/greatdanton/analytics/src/templates"
)

type navbar struct {
	LoggedIn sessions.User
}

// Website renders traffic data for each user website
func Website(w http.ResponseWriter, r *http.Request) {
	user := sessions.LoggedInUser(r)

	err := templates.Execute(w, "displayTraffic", navbar{LoggedIn: user})
	if err != nil {
		fmt.Println("Website: websiteTraffic:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// EditWebsite -> renders new website template filled with
// existing data updates database
func EditWebsite(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Editing website")
}

// DeleteWebsite deletes chosen website from the database
// new traffic is no longer logged in db
// TODO: decide what to do with the existing traffic (delete or not?)
func DeleteWebsite(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Deleting website")
}
