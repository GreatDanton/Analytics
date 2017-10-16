package controller

import (
	"fmt"
	"net/http"

	"github.com/greatdanton/analytics/src/sessions"
	"github.com/greatdanton/analytics/src/templates"
)

// Handle403 handles displaying 403 forbidden template
func Handle403(w http.ResponseWriter, r *http.Request) {
	user := sessions.LoggedInUser(r)
	err := templates.Execute(w, "403", navbar{LoggedIn: user})
	if err != nil {
		fmt.Println("Handle403: template:", err)
		return
	}
}
