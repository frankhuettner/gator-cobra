package cmd

import (
	"context"
	"fmt"

	"github.com/frankhuettner/gator-cobra/internal/database"
	"github.com/spf13/cobra"
)

var followCmd = createProtectedCommand(
	"follow [url]",
	"Follow an RSS feed by its URL",
	`Follow an existing RSS feed by providing its URL.
Example:
  gator-cobra follow https://example.com/feed.xml`,
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

		// Create the feed follow
		feedFollow, err := db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
			UserID: user.ID,
			FeedID: feed.ID,
		})
		if err != nil {
			return fmt.Errorf("couldn't create feed follow: %v", err)
		}

		fmt.Printf("Feed '%s' followed successfully by %s\n", feedFollow.FeedName, feedFollow.UserName)
		return nil
	},
)

func init() {
	rootCmd.AddCommand(followCmd)
}
