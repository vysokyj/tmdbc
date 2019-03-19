package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/ryanbradynd05/go-tmdb"
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

	p, err := filepath.Abs(j.File)
	check(err)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		fmt.Println("File \"" + p + "\" not found!")
		os.Exit(1)
	}
	return j
}

func (j *job) getPosterInDefaultLanguage() tmdb.MovieImage {
	var images *tmdb.MovieImages
	var err error
	ops := make(map[string]string)
	for key, value := range tmdbOptions {
		ops[key] = value
	}
	ops["language"] = "en"
	images, err = tmdbClient.GetMovieImages(j.Movie.ID, ops)
	check(err)
	if len(images.Posters) < 1 {
		fmt.Println("No covers found!")
		os.Exit(1)
	}
	return images.Posters[0]
}

func (j *job) getPoster() tmdb.MovieImage {
	var images *tmdb.MovieImages
	var err error
	images, err = tmdbClient.GetMovieImages(j.Movie.ID, tmdbOptions)
	check(err)
	if len(images.Posters) < 1 {
		fmt.Println("No covers found - trying default language!")
		return j.getPosterInDefaultLanguage()
	}
	return images.Posters[0]
}

func (j *job) downloadPoster() {
	poster := j.getPoster()
	url := "http://image.tmdb.org/t/p/original" + poster.FilePath
	file := filepath.Join(os.TempDir(), "original"+filepath.Ext(poster.FilePath))
	out, err := os.Create(file)
	check(err)
	resp, err := http.Get(url)
	check(err)
	c, err := io.Copy(out, resp.Body)
	check(err)
	fmt.Printf("Poster: %d px x %d px, %d bytes\n", poster.Width, poster.Height, c)
	j.Poster = &poster
	j.PosterFile = file
	check(resp.Body.Close())
	check(out.Close())
}

func (j *job) prepareCovers() {
	ext := filepath.Ext(j.Poster.FilePath)
	coverFile := filepath.Join(os.TempDir(), "cover"+ext)
	coverSmallFile := filepath.Join(os.TempDir(), "cover_small"+ext)
	originalWidth := j.Poster.Width
	originalHeight := j.Poster.Height
	coverWidth := mkvCoverLimit
	coverHeight := originalHeight * coverWidth / originalWidth
	coverSmallWidth := mkvCoverSmallLimit
	coverSmallHeight := originalHeight * coverSmallWidth / originalWidth
	originalImage, err := imaging.Open(j.PosterFile)
	check(err)
	coverImage := imaging.Fit(originalImage, coverWidth, coverHeight, imaging.Lanczos)
	check(imaging.Save(coverImage, coverFile))
	coverSmallImage := imaging.Fit(originalImage, coverSmallWidth, coverSmallHeight, imaging.Lanczos)
	check(imaging.Save(coverSmallImage, coverSmallFile))
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

	/*
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
	*/

	var buffer bytes.Buffer
	cmd.Stdout = &buffer
	cmd.Stderr = &buffer
	err := cmd.Run()
	fmt.Println(buffer.String())
	check(err)
}

func (j *job) getFilename() string {
	var sb strings.Builder
	sb.WriteString(getSafeFileName(j.Movie.Title))
	sb.WriteString(" (")
	sb.WriteString(getYear(j.Movie.ReleaseDate))
	sb.WriteString(")")
	sb.WriteString(j.Extension)
	return sb.String()
}

func (j *job) renameFile() {
	source, err := filepath.Abs(j.File)
	check(err)
	target := filepath.Join(filepath.Dir(source), j.getFilename())
	check(os.Rename(source, target))
	fmt.Println(source + " -> " + target)
}

func (j *job) processMovie() {
	/*
		fmt.Printf("Processing file: %s\n", j.File)
		fmt.Printf("ID: %d\n", j.Movie.ID)
		fmt.Printf("Title: %s\n", j.Movie.Title)
		fmt.Printf("Original title: %s\n", j.Movie.OriginalTitle)
		fmt.Printf("Release date: %s\n", j.Movie.ReleaseDate)
	*/
	fmt.Printf("Found:  %s (%s)\n", j.Movie.Title, getYear(j.Movie.ReleaseDate))
	j.downloadPoster()
	j.prepareCovers()
	j.loadOldMetadata()
	j.addNewMetadata()
	j.renameFile()
}

func (j *job) getMovie(movieShort *tmdb.MovieShort) {
	movie, err := tmdbClient.GetMovieInfo(movieShort.ID, tmdbOptions)
	check(err)
	j.Movie = movie
	j.processMovie()
}

func (j *job) promptSelectAction(name string, year string, results *tmdb.MovieSearchResults) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please select index or action:")
	for index, movieShort := range results.Results {
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
		newName, newYear := getMovieNameAndYear(string(a2[:]))
		j.searchMovie(newName, newYear)
		return
	case "q":
		os.Exit(0)
	}

	i, err := strconv.Atoi(s1)
	if err != nil || i < 0 || i > len(results.Results) {
		fmt.Println("This is not valid index!")
		j.searchMovie(name, year)
	}
	movieShort := results.Results[i-1]
	j.getMovie(&movieShort)
}

func (j *job) searchMovie(name string, year string) {
	var movieShort tmdb.MovieShort
	j.SearchString = name
	fmt.Printf("File:   %s\n", filepath.Base(j.File))
	if year == "" {
		fmt.Printf("Search: %s (????)\n", name)
	} else {
		fmt.Printf("Search: %s (%s)\n", name, year)
	}
	results, err := tmdbClient.SearchMovie(name, tmdbOptions)
	check(err)
	if results.TotalResults == 1 && (year == "" || year == getYear(results.Results[0].ReleaseDate)) {
		movieShort = results.Results[0]
		j.getMovie(&movieShort)
	} else if year != "" {
		count := 0
		for _, m := range results.Results {
			if year == getYear(m.ReleaseDate) {
				count++
				movieShort = m
			}
		}
		if count == 1 {
			j.getMovie(&movieShort)
		} else {
			j.promptSelectAction(name, year, results)
		}
	} else {
		j.promptSelectAction(name, year, results)
	}
}


// searchByFilename search movie by filename
func (j *job) searchByFilename() {
	j.Filename = filepath.Base(j.File)
	j.Extension = filepath.Ext(j.Filename)
	fullName := j.Filename[0 : len(j.Filename)-len(j.Extension)]
	if j.Extension != ".mkv" {
		fmt.Printf("Unsupported movie extension %s\n", j.Extension)
		os.Exit(1)
	}
	name, year := getMovieNameAndYear(fixSearchText(fullName))
	j.searchMovie(name, year)
}
