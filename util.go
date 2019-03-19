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

var multipleSpaces = regexp.MustCompile(`\s+`)

func filterNonPrintableCharacters(s string) string {
	var sb strings.Builder
	for _, r := range s {
		ascii := int(r)
		if ascii >= 0 && ascii < 32 {
			continue
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func getSafeFileName(s string) string {
	// 	https://stackoverflow.com/questions/1976007/what-characters-are-forbidden-in-windows-and-linux-directory-names

	// Non-printable characters
	s = filterNonPrintableCharacters(s)

	// The forbidden printable ASCII characters
	s = strings.ReplaceAll(s, ":", " - ") // readable solution
	s = strings.ReplaceAll(s, "<", "")
	s = strings.ReplaceAll(s, ">", "")
	s = strings.ReplaceAll(s, "\\", "")
	s = strings.ReplaceAll(s, "|", "")
	s = strings.ReplaceAll(s, "/", "")
	s = strings.ReplaceAll(s, "?", "")
	s = strings.ReplaceAll(s, "*", "")

	// Reserved filenames omitted

	// Remove double spaces and trim
	s = multipleSpaces.ReplaceAllString(s, " ")
	s = strings.Trim(s, " ")

	// Windows: Filenames cannot end in a space or dot.
	s = strings.TrimRight(s, ".")

	return s
}
