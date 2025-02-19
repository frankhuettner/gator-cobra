package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var feedsCmd = &cobra.Command{
	Use:   "feeds",
	Short: "List all feeds in the system",
	Long: `Display a list of all RSS feeds in the system.
Example:
  gator-cobra feeds`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Connect to database
		db, err := connectDB()
		if err != nil {
			return err
		}
		defer db.Close()

		// Get all feeds
		feeds, err := db.ListFeedsWithUsers(context.Background())
		if err != nil {
			return fmt.Errorf("failed to fetch feeds: %v", err)
		}

		// Check if there are any feeds
		if len(feeds) == 0 {
			fmt.Println("No feeds found in the system")
			return nil
		}

		// Display feeds
		fmt.Println("Available feeds:")
		for _, feed := range feeds {
			fmt.Printf("* %s\n  URL: %s\n  Created by: %s\n", feed.Name, feed.Url, feed.UserName)
			fmt.Println()
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(feedsCmd)
}
