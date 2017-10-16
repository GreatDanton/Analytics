package controller

import (
	"fmt"
	"net/http"

	"github.com/greatdanton/analytics/src/client"
	"github.com/greatdanton/analytics/src/global"
	"github.com/greatdanton/analytics/src/model"
)

// LogTraffic handles "/l/" of the application and logs user lands
// or user clicked links from website
func LogTraffic(w http.ResponseWriter, r *http.Request) {
	fmt.Println("")
	fmt.Println("Origin:", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Origin", "*") // allow access from all origins
	if r.Method == "POST" {
		r.ParseForm()
		// check if 'link' exist in form => user clicked on link
		if len(r.Form) > 0 {
			userClicks(w, r)
			return
		}
		// user landed on main page
		userLands(w, r)
	}
}

// userLands logs user website lands
func userLands(w http.ResponseWriter, r *http.Request) {
	c := client.GetInfo(r)

	request := c.Request
	data, ok := global.Websites[request]
	if !ok { // this key does not exist
		fmt.Printf("This website shortURL %v does not exist in db: %v\n", request, ok)
		return
	}

	// print client info if request is valid
	fmt.Println("User Landed")
	printClientInfo(c)

	clientIP := c.IP
	websiteID := data.ID
	err := model.LogClientLand(clientIP, websiteID)
	if err != nil {
		fmt.Println("Error LogClientLand:", err)
		return
	}
}

// userClicks log user link clicks that occurs on website
func userClicks(w http.ResponseWriter, r *http.Request) {
	c := client.GetInfo(r)

	request := c.Request
	data, ok := global.Websites[request]
	if !ok {
		fmt.Printf("This website shortURL %v does not exist in db: %v", request, ok)
		return
	}

	// display client info if request is valid
	fmt.Println("UserClicked")
	printClientInfo(c)

	clickedLink := r.Form["link"][0]
	if len(clickedLink) > 7 { // http://
		model.LogClientRequest(c.IP, clickedLink, data.ID)
	}
}

// prints client info into console
func printClientInfo(c client.Client) {
	fmt.Println("IP:", c.IP)
	fmt.Println("Browser:", c.Browser)
	fmt.Println("Proxies:", c.Proxies)
	fmt.Println("Request:", c.Request)
	fmt.Println("Referer", c.Referer)
}
