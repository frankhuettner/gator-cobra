package cmd

import (
	"context"
	"fmt"

	"github.com/frankhuettner/gator-cobra/internal/database"
	"github.com/spf13/cobra"
)

var unfollowCmd = createProtectedCommand(
	"unfollow [url]",
	"Unfollow an RSS feed by its URL",
	`Stop following an RSS feed by providing its URL.
Example:
  gator-cobra unfollow https://example.com/feed.xml`,
	func(cmd *cobra.Command, args []string, user database.User) error {
		if len(args) != 1 {
			return fmt.Errorf("usage: %s <url>", cmd.Name())
		}
		feedURL := args[0]

		// Connect to database
		db, err := connectDB()
		if err != nil {
			return err
		}
		defer db.Close()

		// Get the feed by URL
		feed, err := db.GetFeedByURL(context.Background(), feedURL)
		if err != nil {
			return fmt.Errorf("couldn't find a feed with this URL: %v", err)
		}

		// Delete the feed follow
		err = db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
			UserID: user.ID,
			FeedID: feed.ID,
		})
		if err != nil {
			return fmt.Errorf("couldn't unfollow feed: %v", err)
		}

		fmt.Printf("Successfully unfollowed feed: %s\n", feedURL)
		return nil
	},
)

func init() {
	rootCmd.AddCommand(unfollowCmd)
}
