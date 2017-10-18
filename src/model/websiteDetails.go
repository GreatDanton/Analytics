package model

import (
	"encoding/json"
	"fmt"
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

// GO OOP test

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

// GetLands returns number of lands for each day between the timeFrame
func (w *Website) GetLands(startTime, endTime time.Time) (WebsiteTraffic, error) {
	// get start and end time in correct db format for times ('2006-10-01 15:20:10')
	start := utilities.FormatTime(startTime)
	end := utilities.FormatTime(endTime)
	traffic := WebsiteTraffic{}

	rows, err := global.DB.Query(`SELECT to_char(date_trunc('day', time), 'YYYY-MM-DD') as day, count(*) as lands
								FROM land
								WHERE website_id = $1
								AND time >= $2
								AND time <= $3
								GROUP BY day`, w.ID, start, end)
	defer rows.Close()
	if err != nil {
		return traffic, nil
	}

	var (
		date  string
		count int64
	)
	for rows.Next() {
		err := rows.Scan(&date, &count)
		if err != nil {
			return traffic, err
		}
		ms, err := utilities.ToMiliSecond(date)
		if err != nil {
			return traffic, err
		}

		day := DailyLand{Date: ms, LandNumber: count}
		traffic.Lands = append(traffic.Lands, day)
		traffic.NumOfLands += count
	}
	err = rows.Err()
	if err != nil {
		return traffic, err
	}
	// everything is allright return traffic type
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

// GetLandsJSON returns lands traffic in json string
func (w *Website) GetLandsJSON(timeStart, timeEnd time.Time) (string, error) {
	landsTraffic, err := w.GetLands(timeStart, timeEnd)
	if err != nil {
		return "", fmt.Errorf("GetLandsJSON: GetLands: %v", err)
	}
	bytes, err := json.Marshal(landsTraffic)
	if err != nil {
		return "", fmt.Errorf("GetLandsJSON: cannot unmarshal type: %v", err)
	}
	return string(bytes), err
}

// GetClicks returns number of clicks in the given timeframe
func (w *Website) GetClicks(timeStart, timeEnd time.Time) (WebsiteClicks, error) {
	start := utilities.FormatTime(timeStart)
	end := utilities.FormatTime(timeEnd)

	clicks := WebsiteClicks{}
	rows, err := global.DB.Query(`SELECT to_char(date_trunc('day', time), 'YYYY-MM-DD') as day, count(*) as num from click
								 WHERE website_id = $1
								 AND time >= $2
								 AND time <= $3
								 GROUP BY day`, w.ID, start, end)
	defer rows.Close()
	if err != nil {
		return clicks, err
	}

	var (
		date  string
		count int64
	)

	for rows.Next() {
		err = rows.Scan(&date, &count)
		if err != nil {
			return clicks, err
		}

		ms, err := utilities.ToMiliSecond(date)
		if err != nil {
			return clicks, err
		}
		clicks.Clicks = append(clicks.Clicks, DailyClicks{Date: ms, ClicksNum: count})
		clicks.NumOfClicks += count
	}
	err = rows.Err()
	if err != nil {
		return clicks, err
	}

	return clicks, nil
}

// GetClicksJSON returns json string of clicks data from database in the
// chosen timeframe
func (w *Website) GetClicksJSON(timeStart, timeEnd time.Time) (string, error) {
	clicksTraffic, err := w.GetClicks(timeStart, timeEnd)
	if err != nil {
		return "", fmt.Errorf("GetClicksJSON: GetClicks: %v", err)
	}
	bytes, err := json.Marshal(clicksTraffic)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// GetWebsiteDetail returns website data for website
// with name = website name and owner = userID
func GetWebsiteDetail(websiteID, userID string) (Website, error) {
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
		return Website{}, err
	}
	website := Website{ID: id, Name: name, URL: websiteURL, ShortURL: shortURL}
	return website, nil
}
