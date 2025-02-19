package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/frankhuettner/gator-cobra/internal/database"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var registerCmd = &cobra.Command{
	Use:   "register [username]",
	Short: "Register a new user account",
	Long: `Register a new user account and automatically log in.
Example:
  gator-cobra register john    # Registers and logs in as john`,
	Args: cobra.ExactArgs(1), // Require exactly one argument (username)
	RunE: func(cmd *cobra.Command, args []string) error {
		username := args[0]

		// Connect to database
		db, err := connectDB()
		if err != nil {
			return err
		}
		defer db.Close()

		// Create the user in the database
		userparams := database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      username}
		user, err := db.CreateUser(context.Background(), userparams)
		if err != nil {
			return fmt.Errorf("failed to create user '%s': %v", user.Name, err)
		}

		// Automatically log in the new user
		viper.Set("current_user", user.Name)
		viper.Set("current_user_id", user.ID.String())

		// Save the config
		if err := viper.WriteConfig(); err != nil {
			if err := viper.SafeWriteConfig(); err != nil {
				return fmt.Errorf("error writing config file: %v", err)
			}
		}

		fmt.Printf("Successfully registered and logged in as: %s\n", user.Name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
}
