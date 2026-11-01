package lib

import (
	"errors"
	"net/url"
	"strings"
	"unicode"
)

func ValidateURL(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Hostname() == "" || strings.ContainsFunc(rawURL, unicode.IsSpace) {
		return errors.New("invalid URL")
	}
	return nil
}
