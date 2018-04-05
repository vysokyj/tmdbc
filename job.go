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
	"strings"

	"github.com/disintegration/imaging"
	tmdb "github.com/ryanbradynd05/go-tmdb"
)

const mkvCoverLimit = 600
const mkvCoverSmallLimit = 120

// Job store movie job phases
type job struct {
	File           string
	Filename       string
	Extension      string
	SearchString   string
	Poster         *tmdb.MovieImage
	PosterFile     string
	CoverFile      string
	CoverSmallFile string
	Movie          *tmdb.Movie
	CoverIds       []int
}

// NewJob creates new job
func newJob(file string) *job {
	j := new(job)
	j.File = file
	j.CoverIds = make([]int, 0, 10)
	return j
}

func (j *job) downloadPoster() {
	images, err := tmdbClient.GetMovieImages(j.Movie.ID, tmdbOptions)
	check(err)
	//for index, poster := range images.Posters {
	//	fmt.Printf("%d: %s [%dx%d]\n", index+1, poster.FilePath, poster.Width, poster.Height)
	//}
	if len(images.Posters) < 1 {
		fmt.Println("No covers found!")
		os.Exit(1)
	}
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
	fmt.Printf("Cover image: %s %dx%d px - %d bytes\n", file, poster.Width, poster.Height, c)
	j.Poster = &poster
	j.PosterFile = file
}

func (j *job) prepareCovers() {
	ext := path.Ext(j.Poster.FilePath)
	coverFile := path.Join(os.TempDir(), "cover"+ext)
	coverSmallFile := path.Join(os.TempDir(), "cover_small"+ext)
	originalWidth := j.Poster.Width
	originalHeight := j.Poster.Height
	coverWidht := mkvCoverLimit
	coverHeight := originalHeight * coverWidht / originalWidth
	coverSmallWidht := mkvCoverSmallLimit
	coverSmallHeight := originalHeight * coverSmallWidht / originalWidth
	originalImage, err := imaging.Open(j.PosterFile)
	check(err)
	coverImage := imaging.Fit(originalImage, coverWidht, coverHeight, imaging.Lanczos)
	imaging.Save(coverImage, coverFile)
	coverSmallImage := imaging.Fit(originalImage, coverSmallWidht, coverSmallHeight, imaging.Lanczos)
	imaging.Save(coverSmallImage, coverSmallFile)
	j.CoverFile = coverFile
	j.CoverSmallFile = coverSmallFile
}

func (j *job) loadOldMetadata() {
	cmd := exec.Command("mkvmerge", "-i", j.File)
	var buffer bytes.Buffer
	cmd.Stdout = &buffer
	cmd.Stderr = &buffer
	err := cmd.Run()
	check(err)
	lines := strings.Split(buffer.String(), "\n")
	for _, line := range lines {
		if !strings.Contains(line, ":") {
			// skip non indexed lines
			continue
		}
		lineParts := strings.Split(line, ":")
		key := lineParts[0]
		value := lineParts[1]
		var id int
		if strings.Contains(key, "Attachment ID") {
			id, err = strconv.Atoi(strings.TrimPrefix(key, "Attachment ID "))
			check(err)
		}
		if strings.Contains(value, "cover") || strings.Contains(value, "cover_small") {
			j.CoverIds = append(j.CoverIds, id)
		}
	}
}

func (j *job) addNewMetadata() {
	args := make([]string, 0, 20)
	args = append(args, j.File)
	args = append(args, "--edit")
	args = append(args, "info")
	args = append(args, "--set")
	args = append(args, "title="+j.Movie.Title)

	for _, id := range j.CoverIds {
		args = append(args, "--delete-attachment")
		args = append(args, strconv.Itoa(id))
	}

	args = append(args, "--add-attachment")
	args = append(args, j.CoverFile)
	args = append(args, "--add-attachment")
	args = append(args, j.CoverSmallFile)

	cmd := exec.Command("mkvpropedit", args...)

	fmt.Print("mkvpropedit")
	for _, arg := range args {
		fmt.Print(" ")
		if strings.Contains(arg, " ") {
			fmt.Printf("\"%s\"", arg)
		} else {
			fmt.Print(arg)
		}

	}
	fmt.Print("\n")

	var buffer bytes.Buffer
	cmd.Stdout = &buffer
	cmd.Stderr = &buffer
	err := cmd.Run()
	fmt.Println(buffer.String())
	// check(err)
	if err != nil {
		cmd := exec.Command("mkvpropedit",
			j.File,
			"--edit", "info",
			"--set", "title="+j.Movie.Title,
			"--add-attachment", j.CoverFile,
			"--add-attachment", j.CoverSmallFile,
		)
		var buffer bytes.Buffer
		cmd.Stdout = &buffer
		cmd.Stderr = &buffer
		err := cmd.Run()
		fmt.Println(buffer.String())
		check(err)
	}

}

func (j *job) processMovie() {
	fmt.Printf("Processing file: %s\n", j.File)
	fmt.Printf("ID: %d\n", j.Movie.ID)
	fmt.Printf("Title: %s\n", j.Movie.Title)
	fmt.Printf("Original title: %s\n", j.Movie.OriginalTitle)
	fmt.Printf("Release date: %s\n", j.Movie.ReleaseDate)
	j.downloadPoster()
	j.prepareCovers()
	j.loadOldMetadata()
	j.addNewMetadata()
}

// SearchMovie search job movie by given string
func (j *job) searchMovie(name string) {
	j.SearchString = name
	fmt.Printf("File: %s - searching \"%s\"\n", j.File, name)
	movieSearchResults, err := tmdbClient.SearchMovie(name, tmdbOptions)
	reader := bufio.NewReader(os.Stdin)
	check(err)
	if movieSearchResults.TotalResults > 1 {
		fmt.Println("Please select index or action:")
	}
	for index, movieShort := range movieSearchResults.Results {
		fmt.Printf("%d: %s (%s)\n", index+1, movieShort.Title, getYear(movieShort.ReleaseDate))
	}
	fmt.Println("s: Search")
	fmt.Println("q: Quit")
	a1, _, _ := reader.ReadLine()
	s1 := string(a1)

	switch s1 {
	case "s":
		fmt.Print("Search: ")
		a2, _, _ := reader.ReadLine()
		j.searchMovie(string(a2[:]))
		return
	case "q":
		os.Exit(0)
	}

	i, err := strconv.Atoi(s1)
	if err != nil || i < 0 || i > len(movieSearchResults.Results) {
		fmt.Println("This is not valid index!")
		j.searchMovie(name)
	}
	movieShort := movieSearchResults.Results[i-1]
	movie, err2 := tmdbClient.GetMovieInfo(movieShort.ID, tmdbOptions)
	check(err2)
	j.Movie = movie

	//fmt.Printf("%+v\n", j)
	j.processMovie()
}

// SearchByFilename search movie by filename
func (j *job) searchByFilename() {
	j.Filename = path.Base(j.File)
	j.Extension = path.Ext(j.Filename)
	name := j.Filename[0 : len(j.Filename)-len(j.Extension)]
	if j.Extension != ".mkv" {
		fmt.Printf("Unsupported movie extension %s\n", j.Extension)
		os.Exit(1)
	}
	j.searchMovie(name)
}
