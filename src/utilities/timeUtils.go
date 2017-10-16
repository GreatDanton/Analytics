package utilities

import "time"

// DBtime represents type for easier time conversion
// when communicating with database
type DBtime struct {
	Date time.Time
}
// FormatTime formats time.Time into database time format (2006-01-02 15:04:05)
func (d DBtime) FormatTime() string {
	time := d.Date.Format("2006-01-02 15:04:05") // yyyy-mm-dd hh:mm:ss
	return time
}

// ToMiliSecond returns date from database (2006-01-02)
// into miliseconds that could be displayed in chart
func (d DBtime) ToMiliSecond(t string) (int64, error) {
	u, err := time.Parse(t, "2006-01-02")
	if err != nil {
		return 0, err
	}
	// for some reason chart js needs miliseconds
	micro := u.Unix() * 1000
	return micro, nil
}
