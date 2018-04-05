package main

import (
	"fmt"
	"testing"
)

func TestGetMovieNameAndYear(t *testing.T) {
	s := "Some Movie (1979)"
	name, year := getMovieNameAndYear(s)
	fmt.Println(name)
	fmt.Println(year)
	if name != "Some Movie" || year != "1979" {
		t.Fail()
	}
}
