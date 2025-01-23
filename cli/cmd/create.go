/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/MarkDevOps/AutoGit/cli/pkg/api"
	"github.com/MarkDevOps/AutoGit/cli/pkg/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A create command for various resources",
	Long: `A create command for various resources. For example:
		- Adding Deployment Environments
		- Adding Repositories
		- Adding Secrets
		- Adding Variables`,
	Run: func(cmd *cobra.Command, args []string) {
		var config types.Config
		if err := viper.Unmarshal(&config); err != nil {
			fmt.Printf("Error parsing config: %v\n", err)
			return
		}
		// Get the type flag
		typeFlag, _ := cmd.Flags().GetString("type")
		if typeFlag == "" {
			fmt.Println("Error: --type flag is required. Options: deployment-env, Repository, Secret, Variable")
			return
		}

		switch typeFlag {
		case "deployment-env":
			for repoName, environments := range config.Repos {
				fmt.Printf("\nRepository: %s\n", repoName)
				for envName, envOptions := range environments {
					if envOptions.CreateDeploymentEnv {
						fmt.Printf("\nAttempting to create %s/%s/%s\n", config.Org, repoName, envName)
						if err := api.CreateDeploymentEnv(config.Org, repoName, envName, envOptions); err != nil {
							fmt.Printf("Error creating deployment environment for %s/%s/%s: %v\n", config.Org, repoName, envName, err)
						} else {
							fmt.Printf("Successfully created deployment environment for %s/%s/%s\n", config.Org, repoName, envName)
						}
					} else {
						fmt.Printf("\nSkipping environment %s in repository %s as 'createDeploymentEnv' is false\n", envName, repoName)
					}
				}
			}
		default:
			fmt.Printf("Error: Unsupported --type value '%s'. Options: DeploymentEnv, Repository, Secret, Variable\n", typeFlag)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringP("type", "t", "", "Type of resource to create. Options repository, secret, variable")

}
