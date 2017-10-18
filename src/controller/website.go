package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
	Clicks   string // json string about clicks
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

	t, err := website.GetLands(timeStart, timeEnd)
	if err != nil {
		fmt.Println("Website.GetLands error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(t)
	if err != nil {
		fmt.Println("Cannot marshal:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	details.Traffic = string(b)

	clicks, err := website.GetClicks(timeStart, timeEnd)
	if err != nil {
		fmt.Println("GetNumOfClicks error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	c, err := json.Marshal(clicks)
	if err != nil {
		fmt.Println("Cannot marshal:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	details.Clicks = string(c)

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

	switch r.Method {
	case "GET":
		user := sessions.LoggedInUser(r)
		id := strings.Split(r.RequestURI, "/")[2] // grab website id

		website, err := model.GetWebsiteDetail(id, user.ID)
		if err != nil {
			// if there are no rows in database user is not the owner
			if err == sql.ErrNoRows {
				http.Redirect(w, r, "/403", http.StatusForbidden)
				return
			}
			fmt.Println("EditWebsite: GetWebsiteDetail:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		msg := addWebsiteMsg{Name: website.Name, URL: website.URL,
			Type: "Edit", ShortURL: website.ShortURL}
		renderAddWebsite(w, r, msg)

	// for post request on /website/id/edit page
	case "POST":
		user := sessions.LoggedInUser(r)
		id := strings.Split(r.RequestURI, "/")[2] // grab website id

		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		website, err := model.GetWebsiteDetail(id, user.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Redirect(w, r, "/403", http.StatusForbidden)
				return
			}
			fmt.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		name := r.Form["name"][0]
		url := r.Form["url"][0]
		newShortURL := r.Form["shortURL"][0]
		msg := addWebsiteMsg{Name: name, URL: url, ShortURL: newShortURL, Type: "Edit"}

		if len(newShortURL) > 8 {
			msg.ErrorShortURL = "Short URL should be max 8 characters long"
			renderAddWebsite(w, r, msg)
			return
		}

		// TODO: error check if logged in user is really owner of this url

		exist, err := model.WebsiteShortURLExist(newShortURL)
		if err != nil {
			fmt.Println("WebsiteShortURLExist:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if exist {
			msg.ErrorShortURL = "This short url already exist in db, please pick another one"
			renderAddWebsite(w, r, msg)
			return
		}

		err = model.EditWebsite(user.ID, website.ID, name, url, newShortURL)
		if err != nil {
			fmt.Println("Edit Website:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		err = model.EditWebsiteInMemory(website.ShortURL, newShortURL, website.ID, url)
		if err != nil {
			fmt.Println("EditWebsite handler error:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// successfully updated, redirect to dashboard
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

// DeleteWebsite deletes chosen website from the database
// new traffic is no longer logged in db
// TODO: decide what to do with the existing traffic (delete or not?)
func DeleteWebsite(w http.ResponseWriter, r *http.Request) {
	user := sessions.LoggedInUser(r)
	websiteID := strings.Split(r.RequestURI, "/")[2] // grab website id

	website, err := model.GetWebsiteDetail(websiteID, user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Redirect(w, r, "/403", http.StatusForbidden)
		}
		fmt.Println("Delete website: GetWebsiteDetail:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = model.DeleteWebsite(user.ID, website.ID)
	if err != nil {
		fmt.Println("DeleteWebsite:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	model.DeleteWebsiteInMemory(website.ShortURL)

	// successful delete
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
