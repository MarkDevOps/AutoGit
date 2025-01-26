package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Workflow represents a GitHub Actions workflow.
type Workflow struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	State     string `json:"state"`
	HTMLURL   string `json:"html_url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// WorkflowRun represents a GitHub Actions workflow run.
type WorkflowRun struct {
	ID           int    `json:"id"`
	Status       string `json:"status"`
	Conclusion   string `json:"conclusion"`
	RunNumber    int    `json:"run_number"`
	HTMLURL      string `json:"html_url"`
	HeadBranch   string `json:"head_branch"`
	HeadSHA      string `json:"head_sha"`
	Event        string `json:"event"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	RunStartedAt string `json:"run_started_at"`
}

// FetchWorkflows retrieves all workflows for a given repository.
func FetchWorkflows(org, repo string) ([]Workflow, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/workflows", org, repo)

	// Create a new request using http.NewRequest() and set the Authorization header
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch workflows: %w", err)
	}
	// Set the Authorization header using req.Header.Set()
	req.Header.Add("Authorization", "bearer "+SetHeader())

	// Send the request using http.DefaultClient.Do() and check the response
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch workflows: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch workflows: HTTP %d", resp.StatusCode)
	}

	var response struct {
		Workflows []Workflow `json:"workflows"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode workflows: %w", err)
	}

	return response.Workflows, nil
}

// FetchWorkflowRuns retrieves all runs for a specific workflow in a repository.
func FetchWorkflowRuns(org, repo string, workflowID int) ([]WorkflowRun, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/workflows/%d/runs", org, repo, workflowID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch workflow runs: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch workflow runs: HTTP %d", resp.StatusCode)
	}

	var response struct {
		WorkflowRuns []WorkflowRun `json:"workflow_runs"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode workflow runs: %w", err)
	}

	return response.WorkflowRuns, nil
}
