package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "autogit",
	Short: "A CLI tool for fetching GitHub deployment and release data",
	Long: `autogit is a CLI tool for fetching deployment, release, and workflow data
from GitHub repositories specified in a configuration file.`,
}

// Execute initializes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Persistent flag for specifying configuration file
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.yaml", "Configuration file (default is config.yaml)")

	// Bind Viper to config flag
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		os.Exit(1)
	}
}
