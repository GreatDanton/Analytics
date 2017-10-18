package utilities

import (
	"fmt"
	"time"
)

// FormatTime formats time.Time into database time format
func FormatTime(date time.Time) string {
	// turn time.Time into yyyy-mm-dd hh:mm:ss
	return date.Format("2006-01-02 15:04:05")
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
