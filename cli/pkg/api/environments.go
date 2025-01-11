package api

import (
	"time"

	"github.com/MarkDevOps/AutoGit/cli/pkg/types"
)

// func convertStringToInt(s string) int {
// 	i, err := strconv.Atoi(s)
// 	if err != nil {
// 		return 0
// 	}
// 	return i
// }

func MapEnvironmentData(deployments []types.Deployment, env string) types.EnvData {
	var latestDeployment *types.Deployment
	for _, dep := range deployments {
		if dep.Environment == env {
			if latestDeployment == nil || dep.Timestamp > latestDeployment.Timestamp {
				latestDeployment = &dep
			}
		}
	}

	if latestDeployment == nil {
		return types.EnvData{}
	}

	return types.EnvData{
		DeploymentID:  latestDeployment.ID,
		Ref:           latestDeployment.Ref,
		Description:   latestDeployment.Description,
		CreatedAt:     latestDeployment.CreatedAt.Format(time.RFC3339),
		DeploymentURL: latestDeployment.StatusesURL,
	}
}
