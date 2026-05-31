package main

import (
	"net/url"
	"strings"
)

func normalizeURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(strings.ToLower(rawURL))
	if err != nil {
		return "", err
	}
	fullPath := parsedURL.Hostname() + parsedURL.Path
	output := strings.TrimSuffix(fullPath, "/")
	return output, nil
}
