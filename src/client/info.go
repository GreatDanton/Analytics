package client

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Client struct holds all client data
type Client struct {
	IP      string   // client ip
	Browser string   // client browser
	Proxies []string // possible proxies
	Request string   // client requested url (website shortURL)
	Referer string   // website the client is coming from
}

// GetInfo parses entire client request information
func GetInfo(r *http.Request) Client {
	c := Client{}
	c.IP = GetIP(r)
	c.Browser = GetBrowserVersion(r)
	c.Request = GetURLRequest(r)
	c.Referer = r.Referer()

	p, err := GetProxies(r)
	if err != ErrorNoProxy {
		fmt.Println(err)
	}
	c.Proxies = p

	return c
}

// GetIP returns ip string parsed from request
func GetIP(r *http.Request) string {
	// using reverse proxy in front of the app
	fullIP := r.Header.Get("x-real-ip")
	//fullIP := r.RemoteAddr
	ipArr := strings.Split(fullIP, ":")
	// check for ipv6 ip
	if len(ipArr) > 2 {
		return fullIP
	}
	// ip address is ipv4
	return ipArr[0]
}

// ErrorNoProxy returns an error when x-forwarded-for returns empty array
var ErrorNoProxy = errors.New("getClientProxies: No proxy is present")

// GetProxies returns array of ips fetched from headers x-forwarded-for
// or error if no proxy is present
func GetProxies(r *http.Request) ([]string, error) {
	forwarded := r.Header.Get("x-forwarded-for")
	if len(forwarded) < 2 {
		return []string{}, ErrorNoProxy
	}
	ips := strings.Split(forwarded, ",")
	return ips, nil
}

// GetBrowserVersion returns a browser version (user-agent string)
func GetBrowserVersion(r *http.Request) string {
	b := r.Header.Get("User-Agent")
	return b
}

// GetURLRequest returns the last part of the url
// that is requested by client ex: urlString/8digitNumber
// returns: 8digitNumber
func GetURLRequest(r *http.Request) string {
	u := r.RequestURI[1:] //remove first slash
	return u
}
