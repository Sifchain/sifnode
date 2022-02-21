package utils

import "fmt"

var (
	verbose = false
)

func SetVerbose() {
	verbose = true
}

func Print(msg string) {
	if verbose {
		fmt.Println(msg)
	}
}
