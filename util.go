package main

import (
	"regexp"
	"strings"
)

var movieRegexp = regexp.MustCompile(".*(\\d{4})")

func getYear(date string) string {
	parts := strings.Split(date, "-")
	return parts[0]
}

func getMovieNameAndYear(s string) (name string, year string) {
	if movieRegexp.MatchString(s) {
		l := len(s)
		i := l - len(" (1980)")
		name := s[0:i]
		year := s[i+2 : l-1]
		return name, year
	}

	return s, ""
}

func fixSearchText(s string) string {
	re1 := regexp.MustCompile(`(\.|_|-)`)
	s = re1.ReplaceAllString(s, " ")
	re2 := regexp.MustCompile(`(\S)([A-Z])`)
	s = re2.ReplaceAllString(s, "$1 $2")
	return s
}
