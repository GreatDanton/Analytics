package model

import "github.com/greatdanton/analytics/src/global"

// Website will hold the data of user websites
type Website struct {
	ID   string
	Name string
	URL  string
}

// GetUserWebsites returns array of user websites to be displayed
// on main user dashboard
func GetUserWebsites(userID string) ([]Website, error) {
	rows, err := global.DB.Query(`SELECT id, website_url, name from website
								  where owner = $1`, userID)
	if err != nil {
		return nil, err
	}

	var (
		id         string
		websiteURL string
		name       string
	)
	websites := []Website{}
	for rows.Next() {
		err = rows.Scan(&id, &websiteURL, &name)
		if err != nil {
			return websites, err
		}
		websites = append(websites, Website{ID: id, Name: name, URL: websiteURL})
	}
	err = rows.Err()
	if err != nil {
		return websites, err
	}
	return websites, nil
}
