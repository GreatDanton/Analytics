package memory

import (
	"fmt"
	"sync"

	"github.com/greatdanton/analytics/src/global"
)

//TODO: add interface that will accept both redis and in memory implementation

// Memory stores MemWebsites struct for manipulating with
// memory outside this package
var Memory MemWebsites

// MemWebsites holds websites in memory
type MemWebsites struct {
	Memory map[string]MemWebsite
}

// MemWebsite holds data for each website
type MemWebsite struct {
	ID         string
	WebsiteURL string
	Owner      string
}

// LoadWebsites takes care of loading websites in memory
func (m *MemWebsites) LoadWebsites() error {
	websites := map[string]MemWebsite{}

	var (
		id         string
		owner      string // users.id
		shortURL   string
		websiteURL string
	)
	rows, err := global.DB.Query(`SELECT id, owner, short_url, website_url from website`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &owner, &shortURL, &websiteURL)
		if err != nil {
			return err
		}
		// add website to map
		websites[shortURL] = MemWebsite{ID: id, WebsiteURL: websiteURL, Owner: owner}
	}
	err = rows.Err()
	if err != nil {
		return err
	}

	m.Memory = websites
	return nil
}

// AddWebsite adds website to memory
func (m *MemWebsites) AddWebsite(websiteID, ownerID, shortURL, websiteURL string) error {
	w := m.Memory
	exist := m.ShortURLExist(shortURL)
	if exist {
		return fmt.Errorf("This shortURL already exists in memory")
	}

	var mu sync.Mutex
	mu.Lock()
	w[shortURL] = MemWebsite{ID: websiteID, WebsiteURL: websiteURL, Owner: ownerID}
	mu.Unlock()
	return nil
}

// EditWebsite safely edits website in memory
func (m *MemWebsites) EditWebsite(oldShortURL, newShortURL, websiteURL string) error {
	exist := m.ShortURLExist(oldShortURL)
	if !exist {
		return fmt.Errorf("EditWebsite: old url: %v does not exist", oldShortURL)
	}
	w := m.Memory[oldShortURL]
	ownerID := w.Owner
	websiteID := w.ID

	m.DeleteWebsite(oldShortURL)
	err := m.AddWebsite(websiteID, ownerID, newShortURL, websiteURL)
	if err != nil {
		return err
	}
	return nil
}

// DeleteWebsite takes care of deleting website in memory
func (m *MemWebsites) DeleteWebsite(shortURL string) {
	var mu sync.Mutex
	w := m.Memory
	mu.Lock()
	delete(w, shortURL)
	mu.Unlock()
}

// ShortURLExist checks if shortURL already exists in memory
func (m MemWebsites) ShortURLExist(shortURL string) bool {
	w := m.Memory
	_, exist := w[shortURL]
	// shortURL exist in object
	if exist {
		return true
	}
	// shortURL does not exist in m object
	return false
}

// GetOwner returns owner of the website with shortURL
func (m MemWebsites) GetOwner(shortURL string) (string, error) {
	website, ok := m.Memory[shortURL]
	if !ok {
		return "", fmt.Errorf("This shortURL does not exist")
	}
	return website.Owner, nil
}

// HandleRequest returns data for requested shortURL or error
// if the shortURL does not exist
func (m MemWebsites) HandleRequest(shortURL string) (MemWebsite, error) {
	data, ok := m.Memory[shortURL]
	if !ok {
		return MemWebsite{}, fmt.Errorf("MemWebsites: HandleRequest shortURL: %v does not exist", shortURL)
	}
	return data, nil
}
