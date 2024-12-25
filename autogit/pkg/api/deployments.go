package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MarkDevOps/AutoGit/autogit/pkg/types"
)

func FetchDeployments(org, repo string) ([]types.Deployment, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/deployments", org, repo)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch deployments: %w", err)
	}
	defer resp.Body.Close()

	var deployments []types.Deployment
	if err := json.NewDecoder(resp.Body).Decode(&deployments); err != nil {
		return nil, fmt.Errorf("failed to decode deployments: %w", err)
	}

	return deployments, nil
}
