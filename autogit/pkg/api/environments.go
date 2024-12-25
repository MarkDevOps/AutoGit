package api

import (
	"github.com/MarkDevOps/AutoGit/autogit/pkg/types"
)

func MapEnvironmentData(deployments []types.Deployment, env string) types.EnvData {
	var latestDeployment *types.Deployment
	for _, dep := range deployments {
		if dep.Environment == env {
			if latestDeployment == nil || dep.CreatedAt > latestDeployment.CreatedAt {
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
		CreatedAt:     latestDeployment.CreatedAt,
		DeploymentURL: latestDeployment.StatusesURL,
	}
}
