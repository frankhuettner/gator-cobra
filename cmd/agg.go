package cmd

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/frankhuettner/gator-cobra/internal/database"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Description string    `xml:"description"`
		Items       []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

var aggCmd = &cobra.Command{
	Use:   "agg [time_between_reqs]",
	Short: "Aggregate RSS feeds periodically",
	Long: `Periodically fetch and aggregate RSS feeds from all subscriptions.
Example:
  gator-cobra agg 1m     # Aggregate every minute
  gator-cobra agg 30s    # Aggregate every 30 seconds`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("usage: %s <time_between_reqs>", cmd.Name())
		}

		timeBetweenRequests, err := time.ParseDuration(args[0])
		if err != nil {
			return fmt.Errorf("error parsing time duration: %v", err)
		}

		fmt.Printf("Collecting feeds every %v:\n", timeBetweenRequests)
		ticker := time.NewTicker(timeBetweenRequests)
		defer ticker.Stop()

		// Run immediately once before starting the ticker
		if err := scrapeFeeds(); err != nil {
			fmt.Printf("Error scraping feeds: %v\n", err)
		}

		// Then run on ticker
		for range ticker.C {
			if err := scrapeFeeds(); err != nil {
				fmt.Printf("Error scraping feeds: %v\n", err)
			}
		}

		return nil
	},
}

func scrapeFeeds() error {
	db, err := connectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	feed, err := db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error getting next feed: %v", err)
	}

	fmt.Printf("Fetching feed %s\n", feed.Url)

	// Fetch the feed
	resp, err := http.Get(feed.Url)
	if err != nil {
		return fmt.Errorf("error fetching feed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}

	var rssFeed RSSFeed
	if err := xml.Unmarshal(body, &rssFeed); err != nil {
		return fmt.Errorf("error parsing XML: %v", err)
	}

	// Store each item in the database
	for _, item := range rssFeed.Channel.Items {
		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			// Try alternative date format
			pubAt, err = time.Parse(time.RFC1123, item.PubDate)
			if err != nil {
				pubAt = time.Now() // Use current time if parsing fails
			}
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			PublishedAt: pubAt,
			FeedID:      feed.ID,
		})
		if err != nil {
			fmt.Printf("Error creating post %s: %v\n", item.Title, err)
			continue
		}
	}

	// Mark the feed as fetched
	err = db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("error marking feed as fetched: %v", err)
	}

	fmt.Printf("Collected %d posts from %s\n", len(rssFeed.Channel.Items), feed.Url)
	return nil
}

func init() {
	rootCmd.AddCommand(aggCmd)
}
