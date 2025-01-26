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
			fmt.Println("Error: --type flag is required. Options: deployment-env, Repository, secrets, variables, secrets-variables")
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
		case "secrets":
			for repoName, environments := range config.Repos {
				fmt.Printf("\nRepository: %s\n", repoName)
				for envName, envOptions := range environments {
					if envOptions.CreateSecrets {
						fmt.Printf("\nAttempting to fetch environment public-key for %s/%s/%s\n", config.Org, repoName, envName)
						if publicKey, err := api.GetGithubPublicKey(config.Org, repoName, envName); err != nil {
							fmt.Printf("Error fetching public key for %s/%s/%s: %v\n", config.Org, repoName, envName, err)
						} else {
							fmt.Printf("Successfully fetched public key for %s/%s/%s\n", config.Org, repoName, envName)
							fmt.Printf("\nPublic-key for %s/%s/%s: %s\n", config.Org, repoName, envName, publicKey.(map[string]interface{})["key"].(string)) // Display the public key in terminal (Debug reasons only.)
							for secretName, secretValue := range envOptions.Secrets {
								fmt.Printf("\nAttempting to create/update secret '%s':'%s' within %s/%s/%s\n", secretName, secretValue, config.Org, repoName, envName)
								if err := api.CreateUpdateSecret(config.Org, repoName, envName, secretName, secretValue, publicKey.(map[string]interface{})["key"].(string), publicKey.(map[string]interface{})["key_id"].(string)); err != nil {
									fmt.Printf("Error creating/updating secret %s within %s/%s/%s: %v\n", secretName, config.Org, repoName, envName, err)
								} else {
									fmt.Printf("Successfully created/updated secret %s within %s/%s/%s\n", secretName, config.Org, repoName, envName)
								}
							}
						}
					}
				}
			}
		case "variables":
			for repoName, environments := range config.Repos {
				fmt.Printf("\n  Repository: %s\n", repoName)
				fmt.Println("|-----------------------------|")
				for envName, envOptions := range environments {
					if envOptions.CreateVariables {
						fmt.Printf("\n  Attempting to create/update variables within %s/%s/%s\n\n", config.Org, repoName, envName)
						for variableName, variableValue := range envOptions.Variables {
							fmt.Println("---------------------------------------------------------------------")
							fmt.Printf("\nAttempting to create/update variable '%s':'%s' within %s/%s/%s\n", variableName, variableValue, config.Org, repoName, envName)
							if err := api.CreateUpdateVariable(config.Org, repoName, envName, variableName, variableValue); err != nil {
								fmt.Printf("Error creating/updating variable %s within %s/%s/%s: %s\n", variableName, config.Org, repoName, envName, err)
							}
						}
					}
				}
			}
		case "secrets-variables":

		default:
			fmt.Printf("Error: Unsupported --type value '%s'. Options: deployment-env, Repository, secrets, variables, secrets-variables\n", typeFlag)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringP("type", "t", "", "Type of resource to create. Options repository, secret, variable")

}
