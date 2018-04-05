package main

import (
	"os"
)

//ExistFile checks if file exists
func existFile(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}
