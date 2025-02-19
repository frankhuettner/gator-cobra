package cmd

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/frankhuettner/gator-cobra/internal/database"
	"github.com/spf13/cobra"
)

var browseCmd = createProtectedCommand("browse [limit]",
	"List posts from your feed subscriptions",
	`Display posts from all feeds you're subscribed to.
		Optionally specify the maximum number of posts to show (defaults to 2).
		Example:
		gator-cobra browse     # Shows 2 posts
  		gator-cobra browse 5   # Shows 5 posts`,
	func(cmd *cobra.Command, args []string, user database.User) error {
		limit := 2
		if len(args) > 0 {
			parsedLimit, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid limit argument: %v", err)
			}
			if parsedLimit < 1 {
				return fmt.Errorf("limit must be greater than 0")
			}
			limit = parsedLimit
		}

		db, err := connectDB()
		if err != nil {
			return err
		}
		defer db.Close()

		posts, err := db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
			UserID: user.ID,
			Limit:  int32(limit),
		})
		if err != nil {
			return fmt.Errorf("failed to fetch posts: %v", err)
		}

		if len(posts) == 0 {
			fmt.Printf("No posts found for user %v with user ID %v\n", user.Name, user.ID)
			return nil
		}

		fmt.Printf("Latest %d posts:\n\n", limit)
		for i, post := range posts {
			fmt.Printf("%d. %s\n", i+1, post.Title)
			fmt.Printf("   Published: %s\n", post.PublishedAt.Format(time.RFC822))
			fmt.Printf("   URL: %s\n", post.Url)
			if post.Description.Valid {
				fmt.Printf("   Description: %s\n", post.Description.String)
			}
			fmt.Println() // Empty line between browse
		}

		return nil
	})

func init() {
	rootCmd.AddCommand(browseCmd)
}
