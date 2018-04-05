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

var conf *configuration

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

func checkArgs() {
	if len(os.Args) < 2 {
		fmt.Println("Missing MKV movie path!")
		os.Exit(1)
	}
}

func main() {
	checkArgs()
	conf = loadConfiguration()
	fmt.Println("The Movie Database Client")
	fmt.Println("Copyright \u00A9 Jiří Vysoký, 2018")
	fmt.Printf("API Key: %s\n", conf.APIKey)
	fmt.Printf("Language: %s\n", conf.Language)

	TMDb = tmdb.Init(conf.APIKey)

	for i := 1; i < len(os.Args); i++ {
		job := job{File: os.Args[i]}
		job.searchByFilename()
	}
}
