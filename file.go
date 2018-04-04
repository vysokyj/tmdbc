package main

import (
	"os"
)

//ExistFile checks if file exists
func ExistFile(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}
