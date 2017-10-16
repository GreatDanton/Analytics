package model

import (
	"crypto/rand"
	"fmt"
	"strings"

	"github.com/greatdanton/analytics/src/global"
)

// LoadWebsitesToMemory loads all website data into memory
// TODO - replace this with REDIS db
func LoadWebsitesToMemory() (map[string]global.Website, error) {
	websites := map[string]global.Website{}

	var (
		id         string
		shortURL   string
		websiteURL string
	)
	rows, err := global.DB.Query(`SELECT id, short_url, website_url from website`)
	if err != nil {
		return websites, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &shortURL, &websiteURL)
		if err != nil {
			return websites, err
		}
		// add website to map
		websites[shortURL] = global.Website{ID: id, WebsiteURL: websiteURL}
	}
	err = rows.Err()
	if err != nil {
		return websites, err
	}

	return websites, nil
}

// CreateUniqueShortURL creates unique shortURL that
// can be used to add new website into website database
func CreateUniqueShortURL() (string, error) {
	for {
		// Create unique key
		shortURL, err := createRandomShortURL()
		if err != nil {
			return "", err
		}
		exist := siteURLExistInMemory(shortURL)
		if !exist {
			return shortURL, nil
		}
	}
}

// SiteURLExistInMemory checks if short url exists in memory
// short url ensures the website is tracked
func siteURLExistInMemory(shortURL string) bool {
	w := global.Websites
	_, exist := w[shortURL]
	// shortURL exist in memory (and database)
	if exist {
		return true
	}
	// shortURL does not exist in memory
	return false
}

// CreateRandomShortURL creates short url for website
// (websites are identified with short url)
// creates 8 character string
func createRandomShortURL() (string, error) {
	n := 4
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	s := fmt.Sprintf("%X", b)
	return strings.ToLower(s), nil
}

// AddWebsiteToMemory adds website parameters to memory
func AddWebsiteToMemory(shortURL string, id string, websiteURL string) {
	w := global.Websites
	w[shortURL] = global.Website{ID: id, WebsiteURL: websiteURL}
}
