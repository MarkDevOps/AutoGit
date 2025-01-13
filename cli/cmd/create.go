/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/MarkDevOps/AutoGit/cli/pkg/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A create command for various resources",
	Long: `A create command for various resources. For example:
		- Adding Repositories
		- Adding Secrets
		- Adding Variables`,
	Run: func(cmd *cobra.Command, args []string) {

		var typeFlag string
		if err := viper.Unmarshal(&typeFlag); err != nil {
			fmt.Printf("Error parsing type: %v\n", err)
			return
		}

		if typeFlag == "secret" {
			var config types.ConfigSecrets
			// var secrets types.Secrets

			if err := viper.Unmarshal(&config); err != nil {
				fmt.Printf("Error parsing config: %v\n", err)
				return
			}
			fmt.Println("Creating secrets")
			// Call the CreateConfigSecrets function
			for repo, env := range config.Repos {
				for envName, secrets := range env {
					fmt.Printf("Creating secret for %s in %s/%s\n", secrets, repo, envName)
					fmt.Printf("Name: %s Value: %s\n", secrets, secrets)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	createCmd.Flags().StringP("type", "t", "", "Type of resource to create. Options repository, secret, variable")

}
