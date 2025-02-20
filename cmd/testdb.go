// cmd/testdb.go
package cmd

import (
	"fmt"
	"os"

	"github.com/frankhuettner/gator-cobra/internal/database"
	"github.com/spf13/cobra"
)

// connectDB is a helper function that handles database connection
func connectDB() (*database.DB, error) {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("database URL not configured in environment")
	}
	db, err := database.Connect(dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	return db, nil
}

var testdbCmd = &cobra.Command{
	Use:   "testdb",
	Short: "Test database connection",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := connectDB()
		if err != nil {
			fmt.Printf("Database connection test failed: %v\n", err)
			return
		}
		defer db.Close()

		fmt.Println("Successfully connected to database!")
	},
}

func init() {
	rootCmd.AddCommand(testdbCmd)
}
