package controller

import "net/http"

// MainHandler handles "/" of the application
// if user is logged in it redirects him to /dashboard
// otherwise middleware redirects him to login page
func MainHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
