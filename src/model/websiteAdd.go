package model

import (
	"database/sql"
	"fmt"

	"github.com/greatdanton/analytics/src/global"
)

// WebsiteURLExist checks if url for this particular user
// already exist
func WebsiteURLExist(userID string, url string) (bool, error) {
	var id string
	err := global.DB.QueryRow(`SELECT id from website
								where owner = $1
								AND website_url = $2`, userID, url).Scan(&id)
	if err != nil {
		// if there are no rows website url does not exist in database
		if err == sql.ErrNoRows {
			return false, nil
		}
		// an actual error occured during lookup, return error
		return true, err
	}

	// error did not occur url is present in database
	return true, nil
}

// TrackNewWebsite adds website to the database => the software starts
// tracking records for this website
func TrackNewWebsite(userID string, websiteName string, websiteURL string) error {
	// TODO: create short url that do not exist in the in memory database
	shortURL, err := CreateUniqueShortURL()
	if err != nil {
		return fmt.Errorf("TrackNewWebsite: CreateUniqueShortURL error: %v", err)
	}

	var id string
	err = global.DB.QueryRow(`INSERT into website(owner, name, short_url, website_url)
							   values($1, $2, $3, $4)
							   RETURNING id`, userID, websiteName, shortURL, websiteURL).Scan(&id)
	if err != nil {
		return fmt.Errorf("TrackNewWebsite: error while inserting into website db: %v", err)
	}
	// website is added into database
	// add website into memory
	AddWebsiteToMemory(shortURL, id, websiteURL)
	return nil
}

// TODO: both edit and delete should edit, remove website from
// in memory database

// EditWebsite handles updating website row
func EditWebsite(userID string, websiteID string, websiteName string, websiteURL string) error {
	_, err := global.DB.Exec(`UPDATE website
							 SET name = $1, website_url = $2
							 where owner = $3
							 and id = $4`, websiteName, websiteURL, userID, websiteID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteWebsite handles website deletion from the db
func DeleteWebsite(userID string, websiteID string) error {
	_, err := global.DB.Exec(`DELETE from website
							  where owner = $1
							  and id = $2`, userID, websiteID)
	if err != nil {
		return err
	}
	return nil
}
