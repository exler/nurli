package core

import (
	urllib "net/url"
	"strings"
)

func GetDomainFromURL(url string) string {
	parsed_url, err := urllib.Parse(url)
	if err != nil {
		return ""
	}

	return strings.TrimPrefix(parsed_url.Hostname(), "www.")
}
