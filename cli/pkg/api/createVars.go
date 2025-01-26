package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func ShowVariables(org, repo, env, variable string) (interface{}, error) {
	uri := fmt.Sprintf("https://api.github.com/repos/%s/%s/environments/%s/variables/%s", org, repo, env, variable)
	// Create a new request using http.NewRequest() and set the Authorization header
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send GET to variables Api: %w", err)
	}

	// Set the Authorization header using req.Header.Add()
	req.Header.Add("Authorization", "bearer "+setHeader())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to variables API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get variable: %s", resp.Status)
	}
	// fetching existing variable
	var existingVariable map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&existingVariable)
	if err != nil {
		return nil, fmt.Errorf("failed to decode existing variable: %w", err)
	}

	// fmt.Printf("\nShowing variable for %s in %s\n", variable, repo)
	return existingVariable, nil
}

func PatchVariable(org, repo, env, variable, value string) error {
	uri := fmt.Sprintf("https://api.github.com/repos/%s/%s/environments/%s/variables/%s", org, repo, env, variable)
	// Create a new request using http.NewRequest() and set the Authorization header
	req, err := http.NewRequest("PATCH", uri, nil)
	if err != nil {
		return fmt.Errorf("failed to send PATCH to variables Api: %w", err)
	}

	// set the Authorization header using req.Header.Add()
	req.Header.Add("Authorization", "bearer "+setHeader())
	req.Body = io.NopCloser(
		strings.NewReader(fmt.Sprintf(`
			{
				"value":"%s",
				"name":"%s"
			}
		`, value, variable)),
	)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to variables API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to patch variable: %s", resp.Status)
	}
	fmt.Printf("Updating variable for %s : '%s' in %s/%s\n", variable, value, repo, env)
	return nil

}

func CreateUpdateVariable(org, repo, env, variable, value string) (string, error) {
	var status string
	uri := fmt.Sprintf("https://api.github.com/repos/%s/%s/environments/%s/variables", org, repo, env)
	// Create a new request using http.NewRequest() and set the Authorization header
	req, err := http.NewRequest("POST", uri, nil)
	if err != nil {
		status = "error"
		return "error", fmt.Errorf("failed to send POST to variables Api: %w", err)
	}

	// Set the Authorization header using req.Header.Add()
	req.Header.Add("Authorization", "bearer "+setHeader())
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-GitHub-Api-version", "2022-11-28")
	req.Body = io.NopCloser(strings.NewReader(fmt.Sprintf(`
			{
				"name":"%s",
				"value":"%s"
			}
		`, variable, value)),
	)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		status = "error"
		return "error", fmt.Errorf("failed to send request to variables API: %w", err)
	}
	defer resp.Body.Close()

	status = "Created"

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusConflict {
		status = "error"
		return "error", fmt.Errorf("failed to create variable: %s", resp.Status)
	}
	if resp.StatusCode == http.StatusConflict {

		fmt.Printf("variable already exists: %s", variable)
		fmt.Printf("\nfetching variable: %s", variable)
		existingVariable, err := ShowVariables(org, repo, env, variable)
		fmt.Printf("\nexisting variable: '%s': '%v'\n", variable, existingVariable.(map[string]interface{})["value"].(string))
		if err != nil {
			status = "error"
			return "error", fmt.Errorf("failed to show existing variable: %w", err)
		}
		if existingVariable.(map[string]interface{})["value"].(string) == value {
			fmt.Printf("variable already exists with same value: %s\n\n", variable)
			status = "Unchanged"
		} else {
			fmt.Printf("\nvariable already exists with different value: %s\n", variable)
			fmt.Println("updating variable")
			status = "Changed"
			if err := PatchVariable(org, repo, env, variable, value); err != nil {
				return "error", fmt.Errorf("failed to patch variable: %w", err)
			}
		}
	}
	fmt.Printf("Creating variable for `%s` in %s/%s/%s\n", variable, org, repo, env)
	return status, nil
}
