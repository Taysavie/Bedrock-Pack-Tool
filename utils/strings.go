package utils

import (
	"net/url"
	"regexp"
	"strings"
)

var trimMultiSpace = regexp.MustCompile(`\s+`)

func GetURLHost(link string) string {
	url, err := url.Parse(link)
	if err != nil {
		return ""
	}
	return strings.TrimPrefix(url.Host, "www.")
}

func RemoveExtension(s string) string {
	n := strings.LastIndexByte(s, '.')
	if n == -1 {
		return s
	}
	return s[:n]
}

func TrimMultiSpace(s string) string {
	return trimMultiSpace.ReplaceAllString(s, " ")
}
