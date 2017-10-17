package model

import (
	"time"

	"github.com/greatdanton/analytics/src/global"
	"github.com/greatdanton/analytics/src/utilities"
)

// Website will hold the data of user websites
type Website struct {
	ID       string // website id number
	Name     string // website name
	URL      string // website full url
	ShortURL string // website short url (identifier for collecting data)
}

// WebsiteTraffic is used to parse dailyLands and number of lands
// from the chosen timeframe
type WebsiteTraffic struct {
	NumOfLands int64       // displays whole number of lands in timeframe
	Lands      []DailyLand // all lands fetched by day
}

// DailyLand holds traffic data for chosen day
type DailyLand struct {
	Date       int64 // which day lands occur
	LandNumber int64 // number of lands that day
}

// GetWebsiteLands returns number of lands in the chosen timeframe (timeBeginsWith)
func GetWebsiteLands(websiteID string, timeStart time.Time, timeEnd time.Time) (WebsiteTraffic, error) {
	traffic := WebsiteTraffic{}
	// returns: day | count
	// 2017-15-03  | 5

	tStart := timeStart.Format("2006-01-02")      // yyyy-mm-dd
	tEnd := timeEnd.Format("2006-01-02 15:04:05") // yyyy-mm-dd hh:mm:ss

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
		num  int64
	)

	dbtime := utilities.DBtime{}
	for rows.Next() {
		err = rows.Scan(&date, &num)
		if err != nil {
			return traffic, err
		}

		ms, err := dbtime.ToMiliSecond(date)
		if err != nil {
			return traffic, err
		}
		daily := DailyLand{Date: ms, LandNumber: num}
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

// WebsiteClicks is used to hold clicks data
type WebsiteClicks struct {
	NumOfClicks int64
	Clicks      []DailyClicks
}

// DailyClicks holds clicks per day data
type DailyClicks struct {
	Date      int64
	ClicksNum int64
}

// GetNumOfClicks returns number of clicks in the given timeframe for chosen websiteID
func GetNumOfClicks(websiteID string, timeStart time.Time, timeEnd time.Time) (WebsiteClicks, error) {
	timeS := timeStart.Format("2006-01-02")
	timeE := timeEnd.Format("2006-01-02 15:04:05")
	clicks := WebsiteClicks{}
	rows, err := global.DB.Query(`SELECT to_char(date_trunc('day', time), 'YYYY-MM-DD') as day, count(*) as num from click
								 WHERE website_id = $1
								 AND time >= $2
								 AND time <= $3
								 GROUP BY day`, websiteID, timeS, timeE)
	defer rows.Close()
	if err != nil {
		return clicks, err
	}

	var (
		date string
		num  int64
	)

	dbtime := utilities.DBtime{}
	for rows.Next() {
		err = rows.Scan(&date, &num)
		if err != nil {
			return clicks, err
		}

		ms, err := dbtime.ToMiliSecond(date)
		if err != nil {
			return clicks, err
		}

		clicks.Clicks = append(clicks.Clicks, DailyClicks{Date: ms, ClicksNum: num})
		clicks.NumOfClicks += num
	}
	err = rows.Err()
	if err != nil {
		return clicks, err
	}

	return clicks, nil
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
	website.ShortURL = shortURL
	return website, nil
}
