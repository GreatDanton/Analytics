package model

import (
	"crypto/rand"
	"fmt"
	"strings"

	"github.com/greatdanton/analytics/src/memory"
)

// CreateUniqueShortURL creates unique shortURL that
// can be used to add new website into website database
func CreateUniqueShortURL() (string, error) {
	for {
		// Create unique key
		shortURL, err := createRandomShortURL()
		if err != nil {
			return "", err
		}
		exist := memory.Memory.ShortURLExist(shortURL)
		if !exist {
			return shortURL, nil
		}
	}
}

// CreateRandomShortURL creates short url for website
// (websites are identified with short url)
// creates 8 character string
func createRandomShortURL() (string, error) {
	n := 4
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	s := fmt.Sprintf("%X", b)
	return strings.ToLower(s), nil
}
