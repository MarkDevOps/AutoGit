package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Org   string              `yaml:"org"`
	Repos map[string][]string `yaml:"repos"`
}
type Release struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	HTMLURL string `json:"html_url"`
}
type Workflow struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	HTMLURL    string `json:"html_url"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
	HeadBranch string `json:"head_branch"`
	Headcommit struct {
		ID string `json:"id"`
	} `json:"head_commit"`
}
type WorkflowRepsonse struct {
	TotalCount int        `json:"total_count"`
	Workflows  []Workflow `json:"workflow_runs"`
}

// fetch latest release for a given repo from config file
func fetchLatestRelease(org, repo string) (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", org, repo)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release for %s: %w", repo, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 response for %s: %s", repo, resp.Status)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("Faield to decode response for %s: %w", resp, err)
	}
	return &release, nil
}

func fetchWorkflowRuns(org, repo, tag string) ([]Workflow, error) {

	// url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/runs?event=release", org, repo) // Filtered for events release specificlly
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/runs", org, repo)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch workflows for %s/%s: %w", org, repo, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 repsonse for %s %s: %w", org, repo, err)
	}

	var workflowRepsonse WorkflowRepsonse
	// DEBUGGING \\
	// body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("!! DBUG !! ")
	// fmt.Printf("Raw API Response: %s\n", string(body))
	// fmt.Println("!! ------------------- !! ")
	// DEBUGGING \\

	if err := json.NewDecoder(resp.Body).Decode(&workflowRepsonse); err != nil {
		return nil, fmt.Errorf("failed to decode workflows for %s/%s: %w", org, repo, err)
	}

	// Filter workflows that match release tag
	var matchingWorkflows []Workflow
	for _, wf := range workflowRepsonse.Workflows {
		if strings.Contains(wf.HeadBranch, tag) {
			matchingWorkflows = append(matchingWorkflows, wf)
		}
	}
	return matchingWorkflows, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <path-to-config-file")
		return
	}

	configFilePath := os.Args[1]

	// Load and parse the Yaml config file
	yamlFile, err := os.ReadFile(configFilePath)
	if err != nil {
		fmt.Printf("Error reading YAML file: %v\n", err)
		return
	}

	var config Config
	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		fmt.Printf("Error parsing YAML file: %v\n", err)
		return
	}

	// Display organization name
	fmt.Printf("Organization: %s\n\n", config.Org)

	// Iterate over repos and their environments
	for repo, environments := range config.Repos { // Removed until a use for environments is needed.
		// for repo := range config.Repos {
		fmt.Printf("Fetching latest release for repo: %s/%s\n", config.Org, repo)

		release, err := fetchLatestRelease(config.Org, repo)
		if err != nil {
			fmt.Printf("Error fetching release for repo: %s, %v\n", repo, err)
			continue
		}
		fmt.Printf("  - Latest Release: %s\n", release.TagName)
		fmt.Printf("  - Release Name: %s\n", release.Name)
		fmt.Printf("  - Release URL: %s\n", release.HTMLURL)

		// fetch and display workflows linked to the release
		workflows, err := fetchWorkflowRuns(config.Org, repo, release.TagName)
		if err != nil {
			fmt.Printf("  Error fetching workflows for repo %s: %v\n", repo, err)
		} else {
			fmt.Printf("  - Linked Workflows:\n")
			for _, wf := range workflows {
				fmt.Printf("      - Name:		%s\n", wf.Name)
				fmt.Printf("      - Status:		%s\n", wf.Status)
				fmt.Printf("      - Conclusion: 	%s\n", wf.Conclusion)
				fmt.Printf("      - Workflow URL: 	%s\n", wf.HTMLURL)
			}
		}

		// Display results for each environments
		for _, env := range environments {
			fmt.Printf("Environment: %s\n", env)
			fmt.Printf("  - Latest Release: %s\n", release.TagName)
			fmt.Printf("  - Release Name: %s\n", release.Name)
			fmt.Printf("  - Release URL: %s\n", release.HTMLURL)
		}
		fmt.Println()
	}
}
