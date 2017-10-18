package model

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/greatdanton/analytics/src/global"
	"github.com/greatdanton/analytics/src/memory"
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
	err = memory.Memory.AddWebsite(id, userID, shortURL, websiteURL)
	if err != nil {
		return fmt.Errorf("Memory.AddWebsite: error: %v", err)
	}
	return nil
}

//ErrorShortURLExist is used when the website with short url already exist in memory
// (and therefore in database)
var ErrorShortURLExist = errors.New("Website with this short url already exists")

// EditWebsite handles updating website row
func EditWebsite(userID string, oldWebsite, newWebsite Website) error {
	exist := memory.Memory.ShortURLExist(newWebsite.ShortURL)
	// if shortUrl exist in memory return error
	owner, err := memory.Memory.GetOwner(newWebsite.ShortURL)
	if err != nil {
		return err
	}
	// IF shortURL already exist and the owner of the website is not
	// the same person trying to edit the website => return error
	if exist && owner != userID {
		return ErrorShortURLExist
	}
	//  shortURL does not exist we can update the database
	_, err = global.DB.Exec(`UPDATE website
							 SET name = $1, website_url = $2, short_url = $3
							 where owner = $4
							 and id = $5`, newWebsite.Name, newWebsite.URL, newWebsite.ShortURL, userID, oldWebsite.ID)
	if err != nil {
		return err
	}

	err = memory.Memory.EditWebsite(oldWebsite.ShortURL, newWebsite.ShortURL, newWebsite.URL)
	if err != nil {
		return err
	}

	return nil
}

// DeleteWebsite handles website deletion from the db
func DeleteWebsite(userID string, website Website) error {
	_, err := global.DB.Exec(`DELETE from website
							  where owner = $1
							  and id = $2`, userID, website.ID)
	if err != nil {
		return err
	}
	// Delete website from memory
	memory.Memory.DeleteWebsite(website.ShortURL)
	return nil
}
