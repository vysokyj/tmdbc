package main

import (
	"testing"
)

func TestGetMovieNameAndYear(t *testing.T) {
	s := "Some Movie (1979)"
	name, year := getMovieNameAndYear(s)
	if name != "Some Movie" || year != "1979" {
		t.Fail()
	}
}


func TestFixSearchText1(t *testing.T) {
	input := "Some.Movie.Name"
	output := fixSearchText(input)
	if output != "Some Movie Name" {
		t.Fail()
	}
}

func TestFixSearchText2(t *testing.T) {
	input := "Some_Movie_Name"
	output := fixSearchText(input)
	if output != "Some Movie Name" {
		t.Fail()
	}
}

func TestFixSearchText3(t *testing.T) {
	input := "SomeMovieName"
	output := fixSearchText(input)
	if output != "Some Movie Name" {
		t.Fail()
	}
}

