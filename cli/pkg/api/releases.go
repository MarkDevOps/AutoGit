package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Release represents a GitHub release.
type Release struct {
	ID          int    `json:"id"`
	TagName     string `json:"tag_name"`
	Name        string `json:"name"`
	Body        string `json:"body"`
	PublishedAt string `json:"published_at"`
	HTMLURL     string `json:"html_url"`
}

// FetchReleases retrieves all releases for a given repository.
func FetchReleases(org, repo string) ([]Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", org, repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch releases: %w", err)
	}
	// Set the Authorization header using req.Header.Set()
	req.Header.Set("Authorization", "bearer "+SetHeader())

	// Send the request using http.DefaultClient.Do() and check the response
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch releases: HTTP: %w", err)
	}
	// Defer closing the response body
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch releases: status: %s", resp.Status)
	}

	// Decode the response body into a slice of Release
	var releases []Release
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, fmt.Errorf("failed to decode releases: %w", err)
	}

	return releases, nil
}
