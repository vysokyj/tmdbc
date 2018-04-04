package main

import (
	"encoding/json"
	"io/ioutil"
	"os/user"
	"path"
)

// Configuration stores persistent application settings
type Configuration struct {
	APIKey   string `json:"apiKey"`
	Language string `json:"language"`
}

// LoadConfiguration configuration from ~/.tmdbc file
func LoadConfiguration() *Configuration {
	var c *Configuration
	usr, err1 := user.Current()
	check(err1)
	file := path.Join(usr.HomeDir, ".tmdbc")
	if !ExistFile(file) {
		c = new(Configuration)
	}
	dat, err3 := ioutil.ReadFile(file)
	check(err3)
	err4 := json.Unmarshal(dat, &c)
	check(err4)
	return c
}

// IsSet check if configuration is set
func (c Configuration) IsSet() bool {
	return true
}
