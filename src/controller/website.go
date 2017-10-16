package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/greatdanton/analytics/src/model"
	"github.com/greatdanton/analytics/src/sessions"
	"github.com/greatdanton/analytics/src/templates"
	"github.com/greatdanton/analytics/src/utilities"
)

type navbar struct {
	LoggedIn sessions.User
}

type websiteTraffic struct {
	Website  model.Website
	LoggedIn sessions.User
	Traffic  string // json string about traffic
}

// WebsiteTraffic renders traffic data for each user website
func WebsiteTraffic(w http.ResponseWriter, r *http.Request) {
	user := sessions.LoggedInUser(r)
	websiteID := utilities.GetURLSuffix(r)

	details := websiteTraffic{LoggedIn: user}
	timeEnd := time.Now()                  // now
	timeStart := timeEnd.AddDate(0, -1, 0) // 1 month before now

	// check if user is the owner of the website
	website, err := model.GetWebsiteDetail(websiteID, user.ID)
	if err != nil {
		// if there are no rows returned, user is trying to access
		// website traffic data via url (without being owner)
		if err == sql.ErrNoRows {
			fmt.Println("Logged in user is not the owner of the website")
			err = templates.Execute(w, "403", navbar{LoggedIn: user})
			if err != nil {
				fmt.Println("Website traffic: 403 error:", err)
			}
			return
		}
		// an actual error occured
		fmt.Println("GetWebsiteDetail error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	t, err := model.GetWebsiteLands(website.ID, timeStart, timeEnd)
	if err != nil {
		fmt.Println("GetWebsiteLands error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	details.Website = website

	b, err := json.Marshal(t)
	if err != nil {
		fmt.Println("Cannot marshal:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	details.Traffic = string(b)

	err = templates.Execute(w, "displayTraffic", details)
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
