package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/frankhuettner/gator-cobra/internal/database"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var addfeedCmd = createProtectedCommand(
	"addfeed [feedName] [url]",
	"Add a new RSS feed to your subscriptions",
	`Add a new RSS feed to your subscriptions by providing a feed name and its URL.
Example:
  gator-cobra addfeed "Tech News" https://example.com/feed.xml`,
	func(cmd *cobra.Command, args []string, user database.User) error {
		if len(args) != 2 {
			return fmt.Errorf("usage: %s <feedName> <url>", cmd.Name())
		}

		feedName := args[0]
		feedURL := args[1]

		// Connect to database
		db, err := connectDB()
		if err != nil {
			return err
		}
		defer db.Close()

		// Create the feed subscription
		feed, err := db.CreateFeed(context.Background(), database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      feedName,
			Url:       feedURL,
			UserID:    user.ID,
		})
		if err != nil {
			return fmt.Errorf("failed to add feed: %v", err)
		}

		db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
			UserID: user.ID,
			FeedID: feed.ID,
		})

		fmt.Printf("Successfully added feed '%s' with URL: %s \nfollowed by %s\n", feedName, feedURL, user.Name)
		return nil
	},
)

func init() {
	rootCmd.AddCommand(addfeedCmd)
}
