package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ryanbradynd05/go-tmdb"
)

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var conf *Configuration

// TMDb access The Movie Database API
var TMDb *tmdb.TMDb

func getOptions() map[string]string {
	var options = make(map[string]string)
	options["language"] = conf.Language
	return options
}

func getYear(date string) string {
	parts := strings.Split(date, "-")
	return parts[0]
}

func main() {
	conf = LoadConfiguration()
	fmt.Printf("Language: %s\n", conf.Language)

	if len(os.Args) < 2 {
		fmt.Println("Missing MKV movie path!")
		os.Exit(1)
	}

	TMDb = tmdb.Init(conf.APIKey)

	for i := 1; i < len(os.Args); i++ {
		job := Job{File: os.Args[i]}
		job.SearchByFilename()
	}
}
