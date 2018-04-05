package main

import (
	"encoding/json"
	"io/ioutil"
	"os/user"
	"path/filepath"
)

type configuration struct {
	APIKey   string `json:"apiKey"`
	Language string `json:"language"`
}

func loadConfiguration() *configuration {
	var c *configuration
	usr, err1 := user.Current()
	check(err1)
	file := filepath.Join(usr.HomeDir, ".tmdbc")
	if !existFile(file) {
		c = new(configuration)
	}
	dat, err3 := ioutil.ReadFile(file)
	check(err3)
	err4 := json.Unmarshal(dat, &c)
	check(err4)
	return c
}

func (c configuration) isSet() bool {
	return true
}
