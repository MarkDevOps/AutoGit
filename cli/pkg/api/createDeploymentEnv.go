package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/MarkDevOps/AutoGit/cli/pkg/types"
)

func CreateDeploymentEnv(org, repo, env string, envOptions types.DeploymentEnvOptions) (string, error) {
	var status string
	uri := fmt.Sprintf("https://api.github.com/repos/%s/%s/environments/%s", org, repo, env)
	// Create a new request using http.NewRequest() and set the Authorization header
	req, err := http.NewRequest("PUT", uri, nil)
	if err != nil {
		status = "error"
		return "error:", fmt.Errorf("failed to create deployment environment: %w", err)
	}

	req.Header.Add("Authorization", "bearer "+SetHeader())
	// req.Body = io.NopCloser(
	// 	strings.NewReader(fmt.Sprintln(`{
	// 			"prevent_self_review": false,
	// 			// "reviewers": [
	// 			// 	"type": "Team",
	// 			// 	"id": 12345
	// 			// 	"type": "User",
	// 			// 	"id": 113201461
	// 			// ]
	// 		}`)),
	// )
	// Send the request using http.DefaultClient.Do() and check the response
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		status = "error"
		return "error:", fmt.Errorf("failed to send PUT request to environment API: %w", err)
	}
	defer resp.Body.Close()
	status = "Created&Updated"

	if resp.StatusCode != http.StatusOK {
		status = "error"
		body, _ := io.ReadAll(resp.Body)
		return "error:", fmt.Errorf("failed to create environment: %s", body)
	}

	fmt.Printf("Creating deployment environment for %s in %s\n", env, repo)
	return status, nil
}

func CheckDeployEnv(org, repo, env string, envOptions types.DeploymentEnvOptions) ([]types.EnvCheck, error) {
	uri := fmt.Sprintf("https://api.github.com/repos/%s/%s/environments/%s", org, repo, env)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create deployment environment: %w", err)
	}

	req.Header.Add("Authorization", "bearer "+SetHeader())
	// Send the request using http.DefaultClient.Do() and check the response
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send GET request to environment API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to GET environment: %s", body)
	}

	fmt.Printf("Getting deployment environment for %s in %s\n", env, repo)

	var envCheck []types.EnvCheck
	if err := json.NewDecoder(resp.Body).Decode(&envCheck); err != nil {
		return nil, fmt.Errorf("failed to decode deployments: %v", err)
	}

	return envCheck, nil
}
