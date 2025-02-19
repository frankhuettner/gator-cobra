package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "List all users",
	Long:  `Display a list of all registered users in the system`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Connect to database
		db, err := connectDB()
		if err != nil {
			return err
		}
		defer db.Close()

		users, err := db.ListUsers(context.Background())
		if err != nil {
			return fmt.Errorf("failed to fetch users: %v", err)
		}

		currentUser := viper.GetString("current_user")
		for _, u := range users {
			if u.Name == currentUser {
				fmt.Printf("* %s (current)\n", u.Name)
			} else {
				fmt.Printf("* %s\n", u.Name)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(usersCmd)
}
