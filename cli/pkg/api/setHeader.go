package api

import (
	"fmt"
	"os"
)

func setHeader() string {
	// checking GITHUB_TOKEN environment variable exists
	GITHUB_TOKEN := os.Getenv("GITHUB_TOKEN")
	if GITHUB_TOKEN == "" {
		fmt.Errorf("GITHUB_TOKEN environment variable not set")
	}
	return GITHUB_TOKEN
}
