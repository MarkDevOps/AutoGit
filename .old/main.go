package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// config file mapping.
type Config struct {
	Org   string              `yaml:"org"`
	Repos map[string][]string `yaml:"repos"`
}
type OutputData struct {
	Organization string                        `yaml:"organization"`
	Repositories map[string]map[string]EnvData `yaml:"repositories"`
}
type Release struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	HTMLURL string `json:"html_url"`
}

// workflow responses
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

// fetch workflowAction responses
type WorkflowRepsonse struct {
	TotalCount int        `json:"total_count"`
	Workflows  []Workflow `json:"workflow_runs"`
}

// environment responses
type Environment struct {
	Name             string `json:"name"`
	HTMLURL          string `json:"html_url"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
	LatestDeployment struct {
		ID        int    `json:"id"`
		State     string `json:"state"`
		CreatedAt string `json:"created_at"`
	} `json:"latest_deployment"`
}

type EnvData struct {
	DeploymentID  int    `yaml:"deployment_id,omitempty"`
	Ref           string `yaml:"ref,omitempty"`
	Description   string `yaml:"description,omitempty"`
	CreatedAt     string `yaml:"created_at,omitempty"`
	Status        string `yaml:"status,omitempty"`
	StatusTime    string `yaml:"status_time,omitempty"`
	DeploymentURL string `yaml:"deployment_url,omitempty"`
}

type Deployment struct {
	ID          int    `json:"id"`
	Ref         string `json:"ref"`
	Environment string `json:"environment"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	StatusesURL string `json:"statuses_url"`
}

type DeploymentStatus struct {
	ID        int    `json:"id"`
	State     string `json:"state"`
	CreatedAt string `json:"created_at"`
}

// fetchDeployments fetches deployments for a given repositories
func fetchDeployments(org, repo string) ([]Deployment, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/deployments", org, repo)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch deployments for %s/%s: %w", org, repo, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 response for %s/%s deployments: %s", org, repo, resp.Status)
	}

	var deployments []Deployment
	if err := json.NewDecoder(resp.Body).Decode(&deployments); err != nil {
		return nil, fmt.Errorf("failed to decode deployments for %s/%s: %w", org, repo, err)
	}

	return deployments, nil
}

// fetchLatestDeploymentStatus fetches the latest status for a given deployment
func fetchLatestDeploymentStatus(statusesURL string) (*DeploymentStatus, error) {
	resp, err := http.Get(statusesURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch deployment statuses: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 response for deployment statuses: %s", resp.Status)
	}

	var statuses []DeploymentStatus
	if err := json.NewDecoder(resp.Body).Decode(&statuses); err != nil {
		return nil, fmt.Errorf("failed to decode deployment statuses: %w", err)
	}

	if len(statuses) > 0 {
		return &statuses[0], nil // Return the most recent status
	}
	return nil, nil
}

// fetchEnvironments fetches all environments for a given repository
func fetchEnvironments(org, repo string) ([]Environment, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/environments", org, repo)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch environments for %s/%s: %w", org, repo, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch environments for %s/%s: %w", org, repo, err)
	}

	var response struct {
		Environments []Environment `json:"environments"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode environments for %s/%s: %w", org, repo, err)
	}

	return response.Environments, nil
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

// fetch workflow runs using tag name
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
	outputFilePath := "output.yaml"

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

	output := OutputData{
		Organization: config.Org,
		Repositories: make(map[string]map[string]EnvData),
	}
	// Display organization name
	fmt.Printf("Organization: %s\n\n", config.Org)

	// Iterate over repos and fetch deployments
	for repo, environments := range config.Repos {
		fmt.Printf("Fetching deployments for repo: %s/%s\n", config.Org, repo)
		repoData := make(map[string]EnvData)

		deployments, err := fetchDeployments(config.Org, repo)
		if err != nil {
			fmt.Printf("Error fetching deployments for repo %s: %v\n", repo, err)
			continue
		}

		// Map deployments to environments
		for _, env := range environments {
			fmt.Printf("Environment: %s\n", env)
			var latestDeployment *Deployment
			for _, dep := range deployments {
				if dep.Environment == env {
					if latestDeployment == nil || dep.CreatedAt > latestDeployment.CreatedAt {
						latestDeployment = &dep
					}
				}
			}

			envData := EnvData{}
			if latestDeployment != nil {
				// mapping envData values
				envData.DeploymentID = latestDeployment.ID
				envData.Ref = latestDeployment.Ref

				fmt.Printf("  - Deployment ID: %d\n", latestDeployment.ID)
				fmt.Printf("  - Ref: %s\n", latestDeployment.Ref)

				// Fetch latest deployment status
				status, err := fetchLatestDeploymentStatus(latestDeployment.StatusesURL)
				if err != nil {
					fmt.Printf("  Error fetching deployment status: %v\n", err)
				} else if status != nil {
					fmt.Printf("  - Status: %s\n", status.State)
					fmt.Printf("  - Created At: %s\n", latestDeployment.CreatedAt)
					fmt.Printf("  - Status Created At: %s\n", status.CreatedAt)
					envData.Status = status.State
					envData.StatusTime = status.CreatedAt
				} else {
					fmt.Printf("  - Created At: %s\n", latestDeployment.CreatedAt)
					fmt.Printf("  - No statuses found for this deployment\n")
					envData.CreatedAt = latestDeployment.CreatedAt
				}
				if latestDeployment.Description == "" {
					fmt.Println("  - Description: NOT_SET")
					envData.Description = "NOT_SET"
				} else {
					fmt.Printf("  - Description: %s\n", latestDeployment.Description)
					envData.Description = latestDeployment.Description
				}
			} else {
				fmt.Printf("  - No deployments found for this environment\n")
			}
			repoData[env] = envData
		}
		fmt.Println()

		output.Repositories[repo] = repoData
	}
	yamlOutput, err := yaml.Marshal(&output)
	if err != nil {
		fmt.Printf("Error marshalling output to YAML: %v\n", err)
		return
	}

	if err := os.WriteFile(outputFilePath, yamlOutput, 0644); err != nil {
		fmt.Printf("Error writing to file %s: %v\n", outputFilePath, err)
		return
	}
}
