package cmd

import (
	"context"
	"fmt"

	"github.com/frankhuettner/gator-cobra/internal/database"
	"github.com/spf13/cobra"
)

var followingCmd = createProtectedCommand(
	"following",
	"List all feeds you're following",
	`Display a list of all RSS feeds that you're currently following.
Example:
  gator-cobra following`,
	func(cmd *cobra.Command, args []string, user database.User) error {
		if len(args) != 0 {
			return fmt.Errorf("usage: %s (no arguments needed)", cmd.Name())
		}

		// Connect to database
		db, err := connectDB()
		if err != nil {
			return err
		}
		defer db.Close()

		// Get all feed follows for the user
		feedFollows, err := db.GetFeedFollowsForUser(context.Background(), user.ID)
		if err != nil {
			return fmt.Errorf("couldn't get feeds followed by user %s: %v", user.Name, err)
		}

		if len(feedFollows) == 0 {
			fmt.Printf("User %s is not following any feeds\n", user.Name)
			return nil
		}

		fmt.Printf("Feeds followed by %s:\n", user.Name)
		for _, follow := range feedFollows {
			fmt.Printf("* %s\n", follow.FeedName)
		}

		return nil
	},
)

func init() {
	rootCmd.AddCommand(followingCmd)
}
