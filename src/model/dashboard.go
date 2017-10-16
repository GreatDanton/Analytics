package model

import "github.com/greatdanton/analytics/src/global"

// GetUserWebsites returns array of user websites to be displayed
// on main user dashboard
func GetUserWebsites(userID string) ([]Website, error) {
	rows, err := global.DB.Query(`SELECT id, short_url, website_url, name from website
								  where owner = $1`, userID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var (
		id         string
		shortURL   string
		websiteURL string
		name       string
	)
	websites := []Website{}
	for rows.Next() {
		err = rows.Scan(&id, &shortURL, &websiteURL, &name)
		if err != nil {
			return websites, err
		}
		websites = append(websites, Website{ID: id, Name: name, URL: websiteURL, ShortURL: shortURL})
	}
	err = rows.Err()
	if err != nil {
		return websites, err
	}
	return websites, nil
}
