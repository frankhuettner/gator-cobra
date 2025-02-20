/*
Copyright © 2025 Frank Huettner <frank@huettner.io>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gator-cobra",
	Short: "An implementation of boot.dev Gator in Go using Cobra and Viper",
	Long: `An implementation of boot.dev Gator in Go using Cobra and Viper.
	It's different in that gator-cobra doesn't pass around a state. Instead, we store the state in an .env and a .gator-cobra.yaml file and have joint access to the db.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Initialize configuration
	cobra.OnInitialize(initConfig)

	// Add a global flag for the config file
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/Library/Application Support/gator-cobra/config.yaml or equivalent)")
}

func initConfig() {
	// Find config directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Printf("Error finding config directory: %v\n", err)
		os.Exit(1)
	}

	// Ensure gator-cobra config directory exists
	gatorConfigDir := filepath.Join(configDir, "gator-cobra")
	if err := os.MkdirAll(gatorConfigDir, 0755); err != nil {
		fmt.Printf("Error creating config directory: %v\n", err)
		os.Exit(1)
	}

	// Set up Viper for config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(gatorConfigDir)

	// Set defaults for session management
	viper.SetDefault("current_user", "")
	viper.SetDefault("current_user_id", "")

	// Load env file separately using godotenv
	if err := godotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("Error loading .env file: %v\n", err)
		}
	}

	// Try to read the config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; create it with defaults
			if err := viper.SafeWriteConfig(); err != nil {
				fmt.Printf("Error creating config file: %v\n", err)
			}
		} else {
			fmt.Printf("Error reading config file: %v\n", err)
		}
	}
}
