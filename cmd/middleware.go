package cmd

import (
	"context"
	"fmt"

	"github.com/frankhuettner/gator-cobra/internal/database"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type AuthenticatedCommandFunc func(cmd *cobra.Command, args []string, user database.User) error

func requireLogin(next AuthenticatedCommandFunc) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		currentUser := viper.GetString("current_user")
		if currentUser == "" {
			return fmt.Errorf("you must be logged in to use this command. Use 'gator-cobra login [username]' to log in")
		}

		// Connect to database
		db, err := connectDB()
		if err != nil {
			return err
		}
		defer db.Close()

		// Get the user from database
		user, err := db.GetUser(context.Background(), currentUser)
		if err != nil {
			return fmt.Errorf("error fetching user data: %v", err)
		}

		// Pass the user to the command
		return next(cmd, args, user)
	}
}

func createProtectedCommand(use string, short string, long string, run AuthenticatedCommandFunc) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		RunE:  requireLogin(run),
	}
}
