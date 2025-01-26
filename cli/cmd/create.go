/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/MarkDevOps/AutoGit/cli/pkg/api"
	"github.com/MarkDevOps/AutoGit/cli/pkg/types"
	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A create command for various resources",
	Long: `A create command for various resources. For example:
		- Adding Deployment Environments
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
			fmt.Println("Error: --type flag is required. Options: deployment-env, ALL, secrets, variables, secrets-variables")
			return
		}

		// Define summary map
		summary := make(map[string]string)

		switch typeFlag {
		case "deployment-env":
			for repoName, environments := range config.Repos {
				fmt.Printf("\nRepository: %s\n", repoName)
				for envName, envOptions := range environments {
					if envOptions.CreateDeploymentEnv {
						fmt.Printf("\nAttempting to create %s/%s/%s\n", config.Org, repoName, envName)
						if status, err := api.CreateDeploymentEnv(config.Org, repoName, envName, envOptions); err != nil {
							fmt.Printf("Error creating deployment environment for %s/%s/%s: %v\n", config.Org, repoName, envName, err)
							summary[fmt.Sprintf(" %s/%s/%s/%s/%s/%s", config.Org, repoName, envName, "N/A", "N/A", "N/A")] = "error"
						} else {
							fmt.Printf("Successfully created deployment environment for %s/%s/%s\n", config.Org, repoName, envName)
							summary[fmt.Sprintf(" %s/%s/%s/%s/%s/%s", config.Org, repoName, envName, "N/A", "N/A", "N/A")] = status
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
								if status, err := api.CreateUpdateSecret(config.Org, repoName, envName, secretName, secretValue, publicKey.(map[string]interface{})["key"].(string), publicKey.(map[string]interface{})["key_id"].(string)); err != nil {
									fmt.Printf("Error creating/updating secret %s within %s/%s/%s: %v\n", secretName, config.Org, repoName, envName, err)
									summary[fmt.Sprintf(" %s/%s/%s/%s/%s/%s", config.Org, repoName, envName, "secret", secretName, secretValue)] = "error"
								} else {
									fmt.Printf("Successfully created/updated secret %s within %s/%s/%s\n", secretName, config.Org, repoName, envName)
									summary[fmt.Sprintf(" %s/%s/%s/%s/%s/%s", config.Org, repoName, envName, "secret", secretName, secretValue)] = status
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
							status, err := api.CreateUpdateVariable(config.Org, repoName, envName, variableName, variableValue)
							if err != nil {
								fmt.Printf("Error creating/updating variable %s within %s/%s/%s: %s\n", variableName, config.Org, repoName, envName, err)
								summary[fmt.Sprintf(" %s/%s/%s/%s/%s/%s", config.Org, repoName, envName, "variable", variableName, variableValue)] = "error"
							} else {
								summary[fmt.Sprintf(" %s/%s/%s/%s/%s/%s", config.Org, repoName, envName, "variable", variableName, variableValue)] = status
							}
						}
					}
				}
			}
		case "secrets-variables":
			for repoName, environments := range config.Repos {
				fmt.Printf("\nRepository: %s\n", repoName)
				for envName, envOptions := range environments {
					if envOptions.CreateSecrets && envOptions.CreateVariables {
						fmt.Printf("\nAttempting to fetch environment public-key for %s/%s/%s\n", config.Org, repoName, envName)
						if publicKey, err := api.GetGithubPublicKey(config.Org, repoName, envName); err != nil {
							fmt.Printf("Error fetching public key for %s/%s/%s: %v\n", config.Org, repoName, envName, err)
						} else {
							fmt.Printf("Successfully fetched public key for %s/%s/%s\n", config.Org, repoName, envName)
							fmt.Printf("\nPublic-key for %s/%s/%s: %s\n", config.Org, repoName, envName, publicKey.(map[string]interface{})["key"].(string)) // Display the public key in terminal (Debug reasons only.)
							for secretName, secretValue := range envOptions.Secrets {
								fmt.Println("---------------------------------------------------------------------")
								fmt.Printf("\nAttempting to create/update secret '%s':'%s' within %s/%s/%s\n", secretName, secretValue, config.Org, repoName, envName)
								if status, err := api.CreateUpdateSecret(config.Org, repoName, envName, secretName, secretValue, publicKey.(map[string]interface{})["key"].(string), publicKey.(map[string]interface{})["key_id"].(string)); err != nil {
									fmt.Printf("Error creating/updating secret %s within %s/%s/%s: %v\n", secretName, config.Org, repoName, envName, err)
									summary[fmt.Sprintf(" %s/%s/%s/%s/%s/%s", config.Org, repoName, envName, "secret", secretName, secretValue)] = "error"
								} else {
									fmt.Printf("Successfully created/updated secret %s within %s/%s/%s\n", secretName, config.Org, repoName, envName)
									summary[fmt.Sprintf(" %s/%s/%s/%s/%s/%s", config.Org, repoName, envName, "secret", secretName, secretValue)] = status
								}
							}
							for variableName, variableValue := range envOptions.Variables {
								fmt.Println("---------------------------------------------------------------------")
								fmt.Printf("\nAttempting to create/update variable '%s':'%s' within %s/%s/%s\n", variableName, variableValue, config.Org, repoName, envName)
								status, err := api.CreateUpdateVariable(config.Org, repoName, envName, variableName, variableValue)
								if err != nil {
									fmt.Printf("Error creating/updating variable %s within %s/%s/%s: %s\n", variableName, config.Org, repoName, envName, err)
									summary[fmt.Sprintf(" %s/%s/%s/%s/%s/%s", config.Org, repoName, envName, "variable", variableName, variableValue)] = "error"
								} else {
									summary[fmt.Sprintf(" %s/%s/%s/%s/%s/%s", config.Org, repoName, envName, "variable", variableName, variableValue)] = status
								}
							}
						}
					}
				}
			}
		case "ALL":
			for repoName, environments := range config.Repos {
				fmt.Printf("\nRepository: %s\n", repoName)
				// Deployment Environments
				for envName, envOptions := range environments {
					if envOptions.CreateDeploymentEnv {
						fmt.Printf("\nAttempting to create %s/%s/%s\n", config.Org, repoName, envName)
						if status, err := api.CreateDeploymentEnv(config.Org, repoName, envName, envOptions); err != nil {
							fmt.Printf("Error creating deployment environment for %s/%s/%s: %v\n", config.Org, repoName, envName, err)
							summary[fmt.Sprintf(" %s/%s/%s/%s/%s/%s", config.Org, repoName, envName, "N.A", "N.A", "N.A")] = "error"
						} else {
							fmt.Printf("Successfully created deployment environment for %s/%s/%s\n", config.Org, repoName, envName)
							summary[fmt.Sprintf(" %s/%s/%s/%s/%s/%s", config.Org, repoName, envName, "N.A", "N.A", "N.A")] = status
						}
					} else {
						fmt.Printf("\nSkipping environment %s in repository %s as 'createDeploymentEnv' is false\n", envName, repoName)
					}
					// Secrets
					if envOptions.CreateSecrets {
						fmt.Printf("\nAttempting to fetch environment public-key for %s/%s/%s\n", config.Org, repoName, envName)
						if publicKey, err := api.GetGithubPublicKey(config.Org, repoName, envName); err != nil {
							fmt.Printf("Error fetching public key for %s/%s/%s: %v\n", config.Org, repoName, envName, err)
						} else {
							fmt.Printf("Successfully fetched public key for %s/%s/%s\n", config.Org, repoName, envName)
							fmt.Printf("\nPublic-key for %s/%s/%s: %s\n", config.Org, repoName, envName, publicKey.(map[string]interface{})["key"].(string)) // Display the public key in terminal (Debug reasons only.)
							for secretName, secretValue := range envOptions.Secrets {
								fmt.Printf("\nAttempting to create/update secret '%s':'%s' within %s/%s/%s\n", secretName, secretValue, config.Org, repoName, envName)
								if status, err := api.CreateUpdateSecret(config.Org, repoName, envName, secretName, secretValue, publicKey.(map[string]interface{})["key"].(string), publicKey.(map[string]interface{})["key_id"].(string)); err != nil {
									fmt.Printf("Error creating/updating secret %s within %s/%s/%s: %v\n", secretName, config.Org, repoName, envName, err)
									summary[fmt.Sprintf(" %s/%s/%s/%s/%s/%s", config.Org, repoName, envName, "secret", secretName, secretValue)] = "error"
								} else {
									fmt.Printf("Successfully created/updated secret %s within %s/%s/%s\n", secretName, config.Org, repoName, envName)
									summary[fmt.Sprintf(" %s/%s/%s/%s/%s/%s", config.Org, repoName, envName, "secret", secretName, secretValue)] = status
								}
							}
						}
					}
					// Variables
					if envOptions.CreateVariables {
						fmt.Printf("\n  Attempting to create/update variables within %s/%s/%s\n\n", config.Org, repoName, envName)
						for variableName, variableValue := range envOptions.Variables {
							fmt.Println("---------------------------------------------------------------------")
							fmt.Printf("\nAttempting to create/update variable '%s':'%s' within %s/%s/%s\n", variableName, variableValue, config.Org, repoName, envName)
							status, err := api.CreateUpdateVariable(config.Org, repoName, envName, variableName, variableValue)
							if err != nil {
								fmt.Printf("Error creating/updating variable %s within %s/%s/%s: %s\n", variableName, config.Org, repoName, envName, err)
								summary[fmt.Sprintf(" %s/%s/%s/%s/%s/%s", config.Org, repoName, envName, "variable", variableName, variableValue)] = "error"
							} else {
								summary[fmt.Sprintf(" %s/%s/%s/%s/%s/%s", config.Org, repoName, envName, "variable", variableName, variableValue)] = status
							}
						}
					}
				}
			}
		}
		// Print summary
		fmt.Printf("\n\n\n")

		// Collect summary entries into a slice for sorting
		var summaryEntries []struct {
			Org         string
			Repo        string
			Env         string
			VarOrSecret string
			Name        string
			Value       string
			Status      string
		}

		for item, status := range summary {
			parts := strings.Split(item, "/")
			if len(parts) >= 5 {
				org := parts[0]
				repo := parts[1]
				env := parts[2]
				varOrSecret := parts[3]
				name := parts[4]
				value := parts[5]

				summaryEntries = append(summaryEntries, struct {
					Org         string
					Repo        string
					Env         string
					VarOrSecret string
					Name        string
					Value       string
					Status      string
				}{
					Org:         org,
					Repo:        repo,
					Env:         env,
					VarOrSecret: varOrSecret,
					Name:        name,
					Value:       value,
					Status:      status,
				})
			}
		}
		// Sort the summary entries by org, repo, env, and var/secret
		sort.Slice(summaryEntries, func(i, j int) bool {
			if summaryEntries[i].Repo != summaryEntries[j].Repo {
				return summaryEntries[i].Repo < summaryEntries[j].Repo
			}
			if summaryEntries[i].Env != summaryEntries[j].Env {
				return summaryEntries[i].Env < summaryEntries[j].Env
			}
			return summaryEntries[i].VarOrSecret < summaryEntries[j].VarOrSecret
		})

		// Print summary using table package for better outputformatting
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Org", "Repo", "Environment", "Var/Secret", "Name", "Value", "Status"})

		for _, entry := range summaryEntries {
			status := entry.Status
			if status == "error" {
				status = fmt.Sprintf("%s ❌", status)
			} else if status == "Unchanged" {
				status = fmt.Sprintf("%s 🆗", status)
			} else if status == "Changed" {
				status = fmt.Sprintf("%s 💣", status)
			} else if status == "Created" {
				status = fmt.Sprintf("%s ✅", status)
			} else if status == "Created&Updated" {
				status = fmt.Sprintf("%s 🔄", status) // This is mainly for Deployment Environments are Create and Update are the same PUT operation on the same API
			}
			t.AppendRow([]interface{}{entry.Org, entry.Repo, entry.Env, entry.VarOrSecret, entry.Name, entry.Value, status})
		}
		t.SetStyle(table.StyleColoredBlackOnYellowWhite)
		t.Render()

	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringP("type", "t", "", "Type of resource to create. Options repository, secret, variable")
}
