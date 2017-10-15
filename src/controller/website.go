package controller

import (
	"fmt"
	"net/http"
)

// Website renders traffic data for each user website
func Website(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Displaying website")
}
