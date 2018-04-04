package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strconv"
)

// Job store movie job phases
type Job struct {
	File              string
	Filename          string
	Extension         string
	SearchString      string
	CoverFullSizeFile string
	CoverFile         string
	CoverSmallFile    string
}

// NewJob creates new job
func NewJob(file string) *Job {
	j := new(Job)
	j.File = file
	return j
}

// SearchMovie search job movie by given string
func (j *Job) SearchMovie(name string) {
	j.SearchString = name
	fmt.Printf("Searching: '%s'\n", name)
	movieSearchResults, err := TMDb.SearchMovie(name, getOptions())
	reader := bufio.NewReader(os.Stdin)
	check(err)
	if movieSearchResults.TotalResults > 1 {
		fmt.Println("Candidates:")
	}
	for index, movieShort := range movieSearchResults.Results {
		fmt.Printf("%d: %s (%s)\n", index+1, movieShort.Title, getYear(movieShort.ReleaseDate))
	}
	fmt.Println("s: Search")
	fmt.Println("q: Quit")
	a1, _, _ := reader.ReadLine()
	s1 := string(a1)

	if s1 == "s" {
		fmt.Print("Search: ")
		a2, _, _ := reader.ReadLine()
		j.SearchMovie(string(a2[:]))
		return
	}

	if s1 == "q" {
		os.Exit(0)
	}

	i, _ := strconv.Atoi(s1)
	fmt.Printf("Selected number: %d\n", i)
}

// SearchByFilename search movie by filename
func (j *Job) SearchByFilename() {
	j.Filename = path.Base(j.File)
	j.Extension = path.Ext(j.Filename)
	name := j.Filename[0 : len(j.Filename)-len(j.Extension)]
	if j.Extension != ".mkv" {
		fmt.Printf("Unsupported movie extension %s\n", j.Extension)
		os.Exit(1)
	}
	j.SearchMovie(name)
}
