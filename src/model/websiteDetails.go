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

// WebsiteTraffic is used to parse dailyLands and number of lands
// from the chosen timeframe
type WebsiteTraffic struct {
	NumOfLands int64   // displays whole number of lands in timeframe
	Lands      []Daily // all lands fetched by day
}

// Daily holds date and number of occurrences for each day
// currently we are using it to hold data for
// - number of lands per day
// - number of clicks per day
type Daily struct {
	Date  int64 // which day it occurs
	Count int64 // number of occurrences
}

// GetLands returns number of lands for each day between the timeFrame
func (w Website) GetLands(startTime, endTime time.Time) (WebsiteTraffic, error) {
	// get start and end time in correct db format for times ('2006-10-01 15:20:10')
	start, end := utilities.FormatTimeFrame(startTime, endTime)
	traffic := WebsiteTraffic{}

	rows, err := global.DB.Query(`SELECT to_char(date_trunc('day', time), 'YYYY-MM-DD') as day, count(*) as lands
								FROM land
								WHERE website_id = $1
								AND time >= $2
								AND time <= $3
								GROUP BY day
								ORDER BY day`, w.ID, start, end)
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
		day := Daily{Date: ms, Count: count}
		traffic.Lands = append(traffic.Lands, day)
		traffic.NumOfLands += count
	}
	err = rows.Err()
	if err != nil {
		return traffic, err
	}

	// fix missing dates
	traffic.Lands = addMissingDates(&traffic.Lands, startTime, endTime)

	// everything is allright return traffic type
	return traffic, nil
}

// addMissingDates adds dates with count = 0 between the two dates
// parsed from db that are more than one day apart
func addMissingDates(dateArr *[]Daily, startTime time.Time, endTime time.Time) []Daily {
	fixedArr := []Daily{}
	const dayMS = 24 * 60 * 60 * 1000 // num of miliseconds in one day

	for i := 0; i < len(*dateArr); i++ {
		item := (*dateArr)[i]

		// add missing dates before the first item in dateArr
		// this ensures graph is always displayed for the whole range
		// between startTime and endTime
		if i == 0 {
			start := startTime.Unix()*1000 + dayMS // in ms
			d := item.Date
			if start < d {
				for start < d {
					fixedArr = append(fixedArr, Daily{Date: start, Count: 0})
					start += dayMS
				}
			}
		}

		// add missing dates after the last item in dateArr
		if i+1 == len(*dateArr) {
			fixedArr = append(fixedArr, item)
			// if there are dates missing add them as zero here
			end := endTime.Unix() * 1000 // in ms
			d := item.Date + dayMS
			if d < end {
				for d < end {
					fixedArr = append(fixedArr, Daily{Date: d, Count: 0})
					d += dayMS
				}
			}
			break
		}

		nextItem := (*dateArr)[i+1]
		tmp := Daily{}
		tmp.Date = item.Date
		tmp.Count = item.Count
		fixedArr = append(fixedArr, tmp)

		// add missing dates in between start and end where necessary
		if item.Date+dayMS < nextItem.Date {
			d := item.Date + dayMS
			for d < nextItem.Date {
				fixedArr = append(fixedArr, Daily{Date: d, Count: 0})
				d += dayMS
			}
		}
	}
	return fixedArr
}

// WebsiteClicks is used to hold clicks data
type WebsiteClicks struct {
	NumOfClicks int64 // all clicks counted
	Clicks      []Daily
}

// GetLandsJSON returns lands traffic in json string
func (w Website) GetLandsJSON(timeStart, timeEnd time.Time) (string, error) {
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
func (w Website) GetClicks(timeStart, timeEnd time.Time) (WebsiteClicks, error) {
	start, end := utilities.FormatTimeFrame(timeStart, timeEnd)

	clicks := WebsiteClicks{}
	rows, err := global.DB.Query(`SELECT to_char(date_trunc('day', time), 'YYYY-MM-DD') as day, count(*) as num from click
								 WHERE website_id = $1
								 AND time >= $2
								 AND time <= $3
								 GROUP BY day
								 ORDER BY day`, w.ID, start, end)
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
		clicks.Clicks = append(clicks.Clicks, Daily{Date: ms, Count: count})
		clicks.NumOfClicks += count
	}
	err = rows.Err()
	if err != nil {
		return clicks, err
	}
	clicks.Clicks = addMissingDates(&clicks.Clicks, timeStart, timeEnd)

	return clicks, nil
}

// GetClicksJSON returns json string of clicks data from database in the
// chosen timeframe
func (w Website) GetClicksJSON(timeStart, timeEnd time.Time) (string, error) {
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
