package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/greatdanton/analytics/src/global"
	"github.com/greatdanton/analytics/src/utilities"
)

// Visitors struct for holding data about visitors that are clicking
// on our tracked websites
type Visitors struct {
	TopVisitor string
	Visitors   []Visitor
}

// Visitor holds data about each visitor
type Visitor struct {
	Country   string // to be implemented with microservice
	IP        string // visitor ip
	LastVisit string // last day when the visitor was present
	Amount    string // How many times the same visitor visited in past timeframe
}

// TopVisitors returns last amount of visitors
func (w Website) TopVisitors(timeStart, timeEnd time.Time, amount int) (Visitors, error) {
	v := Visitors{}
	start, end := utilities.FormatTimeFrame(timeStart, timeEnd)

	rows, err := global.DB.Query(`SELECT ip, count(*) as visitedNum from land
								WHERE time >= $1
								AND time <= $2
								GROUP BY ip
								ORDER BY visitedNum desc
								LIMIT $3`, start, end, amount)
	defer rows.Close()
	if err != nil {
		return v, err
	}
	var (
		ip         string
		visitedNum string
	)
	i := 0
	for rows.Next() {
		i++
		err := rows.Scan(&ip, &visitedNum)
		if err != nil {
			return v, err
		}
		if i == 1 {
			v.TopVisitor = ip
		}

		// Country to be implemented -> via separate app?
		// browserID to be implemented
		visitor := Visitor{IP: ip, Amount: visitedNum}
		v.Visitors = append(v.Visitors, visitor)
	}
	err = rows.Err()
	if err != nil {
		return v, err
	}
	return v, nil
}

// TopVisitorsJSON returns data about last visitors in json format
func (w Website) TopVisitorsJSON(timeStart, timeEnd time.Time, amount int) (string, error) {
	data, err := w.TopVisitors(timeStart, timeEnd, amount)
	if err != nil {
		return "", err
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		e := fmt.Errorf("TopVisitorsJSON: json.Marshal error: %v", err)
		return "", e
	}
	return string(bytes), err
}

// LastVisitors returns last visitors that landed on your page with land
// land number right next to them
func (w Website) LastVisitors(timeStart, timeEnd time.Time, amount int) (Visitors, error) {
	v := Visitors{}
	start, end := utilities.FormatTimeFrame(timeStart, timeEnd)

	rows, err := global.DB.Query(`SELECT ip, to_char(date_trunc('minute', MAX(time)), 'YYYY-MM-DD HH24:MI') as day,
								count(*) as visitedNum from land
								WHERE time >= $1
								AND time <= $2
								GROUP BY ip
								ORDER BY MAX(time) desc
								LIMIT $3`, start, end, amount)
	defer rows.Close()
	if err != nil {
		return v, err
	}

	var (
		ip         string
		lastVisit  string
		visitedNum string
	)
	for rows.Next() {
		err := rows.Scan(&ip, &lastVisit, &visitedNum)
		if err != nil {
			return v, err
		}

		visitor := Visitor{IP: ip, LastVisit: lastVisit, Amount: visitedNum}
		v.Visitors = append(v.Visitors, visitor)
	}

	err = rows.Err()
	if err != nil {
		return v, err
	}
	return v, nil
}

// MostClicked holds all clicks with most clicked url
type MostClicked struct {
	MostClickedURL string
	Clicks         []Clicks
}

// Clicks holds amount of clicks and url
type Clicks struct {
	URL       string `json:"Url"`
	ClicksNum int64
}

// MostClicked returns top amount of most clicked links in tracked website
func (w Website) MostClicked(timeStart, timeEnd time.Time, amount int) (MostClicked, error) {
	start, end := utilities.FormatTimeFrame(timeStart, timeEnd)
	mostClicked := MostClicked{}
	rows, err := global.DB.Query(`SELECT url_clicked, count(*) as num from click
								WHERE website_id = $1
								AND time >= $2
								AND time <= $3
								GROUP BY url_clicked
								ORDER BY num desc
								LIMIT $4`, w.ID, start, end, amount)
	defer rows.Close()
	if err != nil {
		return mostClicked, err
	}
	var (
		url    string
		number int64
	)

	i := 0
	for rows.Next() {
		i++
		err = rows.Scan(&url, &number)
		if err != nil {
			return mostClicked, err
		}
		if i == 1 {
			mostClicked.MostClickedURL = url
		}
		c := Clicks{URL: url, ClicksNum: number}
		mostClicked.Clicks = append(mostClicked.Clicks, c)
	}
	err = rows.Err()
	if err != nil {
		return mostClicked, err
	}
	return mostClicked, nil
}
