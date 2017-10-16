package model

import (
	"time"

	"github.com/greatdanton/analytics/src/global"
)

// Website will hold the data of user websites
type Website struct {
	ID       string // website id number
	Name     string // website name
	URL      string // website full url
	shortURL string // website short url (identifier for collecting data)
}

// WebsiteTraffic is used to parse dailyLands and number of lands
// from the chosen timeframe
type WebsiteTraffic struct {
	NumOfLands string      // displays whole number of lands in timeframe
	Lands      []DailyLand // all lands fetched by day
}

// DailyLand holds traffic data for chosen day
type DailyLand struct {
	Date       string // which day lands occur
	LandNumber string // number of lands that day
}

// GetWebsiteLands returns number of lands in the chosen timeframe (timeBeginsWith)
func GetWebsiteLands(websiteID string, timeStart time.Time, timeEnd time.Time) (WebsiteTraffic, error) {
	traffic := WebsiteTraffic{}
	// returns: day | count
	// 2017-15-03   | 5

	tStart := timeStart.Format("2006-01-02")   // yyyy-mm-dd
	tEnd := timeEnd.Format("2006-01-02 15:04") // yyyy-mm-dd hh:mm

	rows, err := global.DB.Query(`SELECT to_char(date_trunc('day', time), 'YYYY-MM-DD') as day, count(*) as lands
								  FROM land
								  WHERE website_id = $1
								  AND time >= $2
								  AND time <= $3
								  GROUP BY day`, websiteID, tStart, tEnd)
	defer rows.Close()
	if err != nil {
		return traffic, err
	}

	var (
		date string
		num  string
	)

	for rows.Next() {
		err = rows.Scan(&date, &num)
		if err != nil {
			return traffic, err
		}
		daily := DailyLand{Date: date, LandNumber: num}
		traffic.Lands = append(traffic.Lands, daily)
		traffic.NumOfLands += num
	}
	err = rows.Err()
	if err != nil {
		return traffic, err
	}

	// everything is allright return traffic data
	return traffic, nil
}

// GetWebsiteDetail returns website data for website
// with name = website name and owner = userID
func GetWebsiteDetail(websiteID, userID string) (Website, error) {
	website := Website{}
	var (
		id         string
		websiteURL string
		shortURL   string
		name       string
	)
	err := global.DB.QueryRow(`SELECT id, name, website_url, short_url from website
								WHERE owner = $1
								and id = $2`, userID, websiteID).Scan(&id, &name, &websiteURL, &shortURL)
	if err != nil {
		return website, err
	}

	website.ID = id
	website.Name = name
	website.URL = websiteURL
	website.shortURL = shortURL
	return website, nil
}
