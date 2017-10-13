package model

import "github.com/greatdanton/analytics/src/global"

// this file will take care of inserting data into db on get and post request
// clients

// LogClientLand logs ip and website on which the client landed
func LogClientLand(ip string, websiteID string) error {
	// we got the correct data, insert user lands into database
	_, err := global.DB.Exec(`INSERT into LAND(website_id, ip)
						   	  values($1, $2)`, websiteID, ip)
	if err != nil {
		return err
	}
	return nil
}

// LogClientRequest adds client click request when the link on website is clicked
func LogClientRequest(userIP string, urlClicked string, websiteID string) error {
	_, err := global.DB.Exec(`INSERT into CLICK(ip, url_clicked, website_id)
							  values($1, $2, $3)`, userIP, urlClicked, websiteID)
	if err != nil {
		return err
	}
	return nil
}
