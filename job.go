package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"

	"github.com/disintegration/imaging"
	tmdb "github.com/ryanbradynd05/go-tmdb"
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
	Movie             *tmdb.Movie
}

// NewJob creates new job
func NewJob(file string) *Job {
	j := new(Job)
	j.File = file
	return j
}

func (j *Job) processMovie() {
	fmt.Printf("ID: %d\n", j.Movie.ID)
	fmt.Printf("Title: %s\n", j.Movie.Title)
	fmt.Printf("Original title: %s\n", j.Movie.OriginalTitle)
	fmt.Printf("Release date: %s\n", j.Movie.ReleaseDate)
	images, err := TMDb.GetMovieImages(j.Movie.ID, getOptions())
	check(err)
	//for index, poster := range images.Posters {
	//	fmt.Printf("%d: %s [%dx%d]\n", index+1, poster.FilePath, poster.Width, poster.Height)
	//}
	poster := images.Posters[0]
	//url := "http://image.tmdb.org/t/p/w600" + poster.FilePath
	url := "http://image.tmdb.org/t/p/original" + poster.FilePath
	file := path.Join(os.TempDir(), "original"+path.Ext(poster.FilePath))
	out, err := os.Create(file)
	//fmt.Printf("%s\n", poster.FilePath)
	//fmt.Printf("%s\n", poster.Iso639_1)
	check(err)
	defer out.Close()
	resp, err := http.Get(url)
	check(err)
	defer resp.Body.Close()
	c, err := io.Copy(out, resp.Body)
	check(err)
	fmt.Printf("%s -> %s (%d bytes)\n", url, file, c)

	coverFile := path.Join(os.TempDir(), "cover"+path.Ext(poster.FilePath))
	coverSmallFile := path.Join(os.TempDir(), "cover_small"+path.Ext(poster.FilePath))

	originalWidth := poster.Width
	originalHeight := poster.Height
	coverWidht := 600 // MKV limit
	coverHeight := originalHeight * coverWidht / originalWidth
	coverSmallWidht := 120 // MKV limit
	coverSmallHeight := originalHeight * coverSmallWidht / originalWidth
	originalImage, err := imaging.Open(file)
	check(err)
	coverImage := imaging.Fit(originalImage, coverWidht, coverHeight, imaging.Lanczos)
	imaging.Save(coverImage, coverFile)
	coverSmallImage := imaging.Fit(originalImage, coverSmallWidht, coverSmallHeight, imaging.Lanczos)
	imaging.Save(coverSmallImage, coverSmallFile)

	cmd := exec.Command("mkvpropedit",
		j.File,
		"--edit", "info",
		"--set", "title="+j.Movie.Title,
		"--add-attachment", coverFile,
		"--add-attachment", coverSmallFile,
	)
	var buffer bytes.Buffer
	cmd.Stdout = &buffer
	cmd.Stderr = &buffer
	err = cmd.Run()
	fmt.Println(buffer.String())
	check(err)
}

// SearchMovie search job movie by given string
func (j *Job) SearchMovie(name string) {
	j.SearchString = name
	fmt.Printf("File: %s\n", j.File)
	fmt.Printf("Searching: %s\n", name)
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
	//fmt.Printf("Selected number: %d\n", i)
	//TODO: Check index
	movieShort := movieSearchResults.Results[i-1]
	movie, err2 := TMDb.GetMovieInfo(movieShort.ID, getOptions())
	check(err2)
	j.Movie = movie

	//fmt.Printf("%+v\n", j)
	j.processMovie()
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
