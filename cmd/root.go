/*
Copyright © 2025 Frank Huettner <frank@huettner.io>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

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
	// Add default db
	viper.SetDefault("DB_URL", "")
	viper.AutomaticEnv()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		configDir, err := os.UserConfigDir()
		if err != nil {
			fmt.Println("Error finding config directory:", err)
			os.Exit(1)
		}

		// Create the gator-cobra directory inside the config directory
		gatorConfigDir := filepath.Join(configDir, "gator-cobra")
		if err := os.MkdirAll(gatorConfigDir, 0755); err != nil {
			fmt.Println("Error creating config directory:", err)
			os.Exit(1)
		}

		viper.AddConfigPath(gatorConfigDir)
		viper.SetConfigName("config") // Changed from .gator-cobra to config
		viper.SetConfigType("yaml")
	}

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			defaultConfig := map[string]interface{}{
				"current_user":    "",
				"current_user_id": "",
			}
			for k, v := range defaultConfig {
				viper.SetDefault(k, v)
			}
			if err := viper.SafeWriteConfig(); err != nil {
				fmt.Printf("Error creating config file: %v\n", err)
			}
		} else {
			fmt.Printf("Error reading config file: %v\n", err)
		}
	}
}
