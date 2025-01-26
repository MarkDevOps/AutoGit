package api

import (
	"fmt"
	"os"
)

func SetHeader() string {
	// checking GITHUB_TOKEN environment variable exists
	GITHUB_TOKEN := os.Getenv("GITHUB_TOKEN")
	if GITHUB_TOKEN == "" {
		fmt.Println("GITHUB_TOKEN environment variable not set")
		return "Error!!"
	}
	return GITHUB_TOKEN
}
