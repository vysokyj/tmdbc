package main

import (
	"encoding/json"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"fmt"
	"bufio"
	"os"
)

type configuration struct {
	APIKey   string `json:"apiKey"`
	Language string `json:"language"`
}

func getConfigurationFile() string {
	usr, err := user.Current()
	check(err)
	return filepath.Join(usr.HomeDir, ".tmdbc")
}

func loadConfiguration() *configuration {
	var c *configuration
	file := getConfigurationFile()
	if !existFile(file) {
		c = new(configuration)
		c.askApiKey()
		c.askLanguage()
		c.save(file)
	} else {
		dat, err := ioutil.ReadFile(file)
		check(err)
		err = json.Unmarshal(dat, &c)
		check(err)
	}
	return c
}

func (c *configuration) askApiKey() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("API key: ")
	bytes, _, err := reader.ReadLine()
	check(err)
	//TODO: Check with regular expression
	s := string(bytes)
	fmt.Println(s)
	c.APIKey = s
}

func (c *configuration) askLanguage() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Language: ")
	bytes, _, err := reader.ReadLine()
	check(err)
	//TODO: Check with regular expression
	s := string(bytes)
	fmt.Println(s)
	c.Language = s
}

func (c *configuration) save(file string) {
	fmt.Println(c)
	bytes, err := json.Marshal(c)
	check(err)
	err = ioutil.WriteFile(file, bytes, 0644)
	check(err)
}
