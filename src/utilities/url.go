package utilities

import (
	"net/http"
	"strings"
)

// GetURLSuffix returns last part of the url
// url /website/myWebsite => returns myWebsite
func GetURLSuffix(r *http.Request) string {
	url := strings.TrimRight(r.URL.Path, "/")
	path := strings.Split(url, "/")
	suffix := path[len(path)-1]
	return suffix
}
