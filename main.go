package main

import (
	"fmt"
	"os"

	"github.com/ryanbradynd05/go-tmdb"
)

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var config *configuration
var tmdbOptions map[string]string
var tmdbClient *tmdb.TMDb

func checkArgs() {
	if len(os.Args) < 2 {
		fmt.Println("Missing MKV movie path!")
		os.Exit(1)
	}
}

func printHeader() {
	fmt.Println("The Movie Database Client")
	fmt.Println("Copyright \u00A9 Jiří Vysoký, 2018")
	fmt.Println("License: GNU GENERAL PUBLIC LICENSE Version 3, 29 June 2007")
	fmt.Printf("API Key: %s\n", config.APIKey)
	fmt.Printf("Language: %s\n\n", config.Language)
}

func main() {

	config = loadConfiguration()
	checkArgs()
	printHeader()

	var tmdbConfig = tmdb.Config{
		APIKey:   config.APIKey,
	}
	tmdbOptions = make(map[string]string)
	tmdbOptions["language"] = config.Language
	tmdbClient = tmdb.Init(tmdbConfig)

	for i := 1; i < len(os.Args); i++ {
		job := newJob(os.Args[i])
		job.searchByFilename()
	}
}
