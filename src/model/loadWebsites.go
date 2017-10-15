package model

import "github.com/greatdanton/analytics/src/global"

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

// SiteURLExistInMemory checks if short url exists in memory
// short url ensures the website is tracked
func SiteURLExistInMemory(shortURL string) bool {
	w := global.Websites
	_, exist := w[shortURL]
	// shortURL exist in memory (and database)
	if exist {
		return true
	}
	// shortURL does not exist in memory
	return false
}

// AddWebsiteToMemory adds website parameters to memory
func AddWebsiteToMemory(shortURL string, id string, websiteURL string) {
	w := global.Websites
	w[shortURL] = global.Website{ID: id, WebsiteURL: websiteURL}
}
