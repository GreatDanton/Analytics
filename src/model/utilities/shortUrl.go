package utilities

import (
	"crypto/rand"
	"fmt"
	"strings"
)

// CreateRandomShortURL creates short url for website
// (websites are identified with short url)
// creates 8 character string
func CreateRandomShortURL() (string, error) {
	n := 4
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	s := fmt.Sprintf("%X", b)
	return strings.ToLower(s), nil
}
