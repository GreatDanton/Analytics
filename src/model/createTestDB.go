package model

import (
	"fmt"

	"github.com/greatdanton/analytics/src/global"
)

// CreateTestDB drops all old tables, creates new and inserts initial data
func CreateTestDB() {
	CreateDB() // drop old database and create new tables
	_, err := global.DB.Exec(`INSERT into users(username, email, password)
							  values('user1', 'some@email.com', '12345');`)
	if err != nil {
		fmt.Println("Problem inserting data into users:", err)
		return
	}
	_, err = global.DB.Exec(`INSERT into website(short_url, owner, active, website_url)
							  values('12345678', 1, TRUE, 'http://jan.pribosek.com');`)
	if err != nil {
		fmt.Println("Problem inserting data into website:", err)
		return
	}
}
