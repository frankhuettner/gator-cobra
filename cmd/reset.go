package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the user database",
	Long: `Reset the user database by deleting all users and resetting the ID counter.
This action cannot be undone!`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// // Prompt for confirmation
		// fmt.Print("This will delete ALL users from the database. Are you sure? (y/N): ")
		// var response string
		// fmt.Scanln(&response)

		// if response != "y" && response != "Y" {
		// 	fmt.Println("Operation cancelled")
		// 	return nil
		// }

		db, err := connectDB()
		if err != nil {
			return err
		}
		defer db.Close()

		if err := db.DeleteAllUsers(context.Background()); err != nil {
			return fmt.Errorf("failed to delete users: %v", err)
		}

		viper.Set("current_user", "")
		viper.Set("current_user_id", "")
		if err := viper.WriteConfig(); err != nil {
			return fmt.Errorf("error writing config file: %v", err)
		}

		fmt.Println("Database reset successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
