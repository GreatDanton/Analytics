package utilities

import (
	"fmt"
	"time"
)

// FormatTime formats time.Time into database time format
func formatTime(date time.Time) string {
	// turn time.Time into yyyy-mm-dd hh:mm:ss
	return date.Format("2006-01-02 15:04:05")
}

// FormatTimeFrame returns starting, ending time formatted
// in db format
func FormatTimeFrame(start, end time.Time) (string, string) {
	s := formatTime(start)
	e := formatTime(end)
	return s, e
}

// FormatTimeMany returns array of strings of time arguments
// formatted for db timeframe usage
func FormatTimeMany(args ...interface{}) ([]string, error) {
	timeArr := make([]string, 0, len(args))
	for _, v := range args {
		// type asertion
		t, ok := v.(time.Time)
		if !ok {
			err := fmt.Errorf("FormatTimeMany: provided argument %v is not type of time.Time", v)
			return timeArr, err
		}
		timeArr = append(timeArr, formatTime(t))
	}
	return timeArr, nil
}

// ToMiliSecond returns date from databse(2006-01-02)
// into miliseconds that could be displayed in chart
func ToMiliSecond(date string) (int64, error) {
	u, err := time.Parse("2006-01-02", date)
	if err != nil {
		return 0, fmt.Errorf("ToMiliSecond error:%v", err)
	}
	// for some reason chart js needs miliseconds
	ms := u.Unix() * 1000
	return ms, nil
}
