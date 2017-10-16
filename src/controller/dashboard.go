package controller

import (
	"fmt"
	"net/http"

	"github.com/greatdanton/analytics/src/model"
	"github.com/greatdanton/analytics/src/sessions"
	"github.com/greatdanton/analytics/src/templates"
)

// dashboardDisplay is used to display all tracked website for
// the currently logged in user
type dashboardDisplay struct {
	Websites []model.Website
	LoggedIn sessions.User
}

// Dashboard takes care of displaying main user dashboard
func Dashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// if user is not logged in redirect him to login page
		user := sessions.LoggedInUser(r)
		// get user websites and display them on dashboard
		dashboard := dashboardDisplay{LoggedIn: user}
		websites, err := model.GetUserWebsites(user.ID)
		if err != nil {
			fmt.Println("GetUserWebsites error:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		dashboard.Websites = websites

		err = templates.Execute(w, "mainDashboard", dashboard)
		if err != nil {
			fmt.Println("template mainDashboard error:", err)
			return
		}
	}
}

type addWebsiteMsg struct {
	Type     string // type of template render: edit or add
	Name     string
	URL      string
	ErrorURL string
	LoggedIn sessions.User
}

// AddWebsite displays template for adding new website to the
// main dashboard
func AddWebsite(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		renderAddWebsite(w, r, addWebsiteMsg{})
	case "POST":
		addWebsite(w, r)
	default:
		renderAddWebsite(w, r, addWebsiteMsg{})
	}
}

// renderAddWebsite renders add website template card
func renderAddWebsite(w http.ResponseWriter, r *http.Request, msg addWebsiteMsg) {
	user := sessions.LoggedInUser(r)
	msg.LoggedIn = user

	err := templates.Execute(w, "addWebsite", msg)
	if err != nil {
		fmt.Println("template addWebsite:", err)
		return
	}
}

// addWebsite adds website to the database
func addWebsite(w http.ResponseWriter, r *http.Request) {
	user := sessions.LoggedInUser(r)

	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error parsing form:", err)
		return
	}

	name := r.Form["name"][0]
	url := r.Form["url"][0]

	// check if url already exist
	exist, err := model.WebsiteURLExist(user.ID, url)
	if err != nil {
		fmt.Println("WebsiteURLExist error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	msg := addWebsiteMsg{LoggedIn: user}
	// website url for this particular user already exist, inform user
	// about that
	if exist {
		msg.ErrorURL = "This url already exist, please choose another"
		renderAddWebsite(w, r, msg)
		return
	}

	// url does not exist in the database add new field into website database
	err = model.TrackNewWebsite(user.ID, name, url)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// New website is added to tracking dashboard
	// redirect to dashboard
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
