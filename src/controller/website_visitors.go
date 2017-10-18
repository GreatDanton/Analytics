package controller

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/greatdanton/analytics/src/model"
	"github.com/greatdanton/analytics/src/sessions"
	"github.com/greatdanton/analytics/src/templates"
)

type displayVisitors struct {
	Website     model.Website
	TopVisitors model.Visitors
	MostClicked model.MostClicked
	LoggedIn    sessions.User
}

// WebsiteVisitors displays analyzed visitors data
func WebsiteVisitors(w http.ResponseWriter, r *http.Request) {
	user := sessions.LoggedInUser(r)
	websiteID := strings.Split(r.RequestURI, "/")[2]
	// TODO: Write owner middleware to get rid of this repetition
	website, err := model.GetWebsiteDetail(websiteID, user.ID)
	if err != nil {
		fmt.Println("GetWebsiteDetail:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	timeEnd := time.Now()
	timeStart := timeEnd.AddDate(0, -1, 0)

	visitors := displayVisitors{LoggedIn: user, Website: website}

	// get visitor data
	v, err := website.LastVisitors(timeStart, timeEnd, 10)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	visitors.TopVisitors = v

	visitors.MostClicked, err = website.MostClicked(timeStart, timeEnd, 10)
	if err != nil {
		fmt.Println("Website.MostClicked error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = templates.Execute(w, "displayVisitors", visitors)
	if err != nil {
		return
	}
}
