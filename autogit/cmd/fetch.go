package cmd

import (
	"fmt"

	"github.com/MarkDevOps/AutoGit/autogit/pkg/api"
	"github.com/MarkDevOps/AutoGit/autogit/pkg/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var outputFile string

func init() {
	fetchCmd.Flags().StringVarP(&outputFile, "output", "o", "output.yaml", "Path to the output YAML file")
	rootCmd.AddCommand(fetchCmd)
}

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch deployment and release data",
	Long:  "Fetch deployment, release, and workflow data from GitHub repositories specified in the configuration file.",
	Run: func(cmd *cobra.Command, args []string) {
		var config types.Config
		if err := viper.Unmarshal(&config); err != nil {
			fmt.Printf("Error parsing config: %v\n", err)
			return
		}

		output := types.OutputData{
			Organization: config.Org,
			Repositories: make(map[string]map[string]types.EnvData),
		}

		fmt.Printf("Organization: %s\n\n", config.Org)

		for repo, environments := range config.Repos {
			fmt.Printf("Fetching deployments for repo: %s/%s\n", config.Org, repo)
			repoData := make(map[string]types.EnvData)

			deployments, err := api.FetchDeployments(config.Org, repo)
			if err != nil {
				fmt.Printf("Error fetching deployments for repo %s: %v\n", repo, err)
				continue
			}

			for _, env := range environments {
				envData := api.MapEnvironmentData(deployments, env)
				repoData[env] = envData
			}

			output.Repositories[repo] = repoData
		}

		if err := api.WriteOutput(output, outputFile); err != nil {
			fmt.Printf("Error writing output: %v\n", err)
		}
	},
}
