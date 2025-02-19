// cmd/login.go
package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var loginCmd = &cobra.Command{
	Use:   "login [username]",
	Short: "Log in as a specific user",
	Long: `Log in as a specific user for the RSS aggregator or show the current user if no username is provided.
Example:
  gator-cobra login john    # Logs in as john
  gator-cobra login        # Shows the current logged-in user`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			// Show current user
			currentUser := viper.GetString("current_user")
			if currentUser == "" {
				fmt.Println("Not logged in")
				return nil
			}
			fmt.Printf("Logged in as: %s\n", currentUser)
			return nil
		}

		username := args[0]

		db, err := connectDB()
		if err != nil {
			fmt.Printf("Database connection test failed: %v\n", err)
			return err
		}
		defer db.Close()

		user, err := db.GetUser(context.Background(), username)
		if err != nil {
			return fmt.Errorf("user '%s' not found in database: %v", username, err)
		}

		// Set the current user in the config
		viper.Set("current_user", username)
		viper.Set("current_user_id", user.ID.String())

		// Save the config
		if err := viper.WriteConfig(); err != nil {
			if err := viper.SafeWriteConfig(); err != nil {
				return fmt.Errorf("error writing config file: %v", err)
			}
		}

		fmt.Printf("Logged in as: %s\n", username)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Add default values
	viper.SetDefault("current_user", "")
	viper.SetDefault("current_user_id", "")
}
