# Gator-Cobra: RSS Feed Aggregator CLI

A command-line RSS feed aggregator built in Go that helps users follow and aggregate content from multiple RSS feeds in one place.

## Prerequisites

- Go 1.23.4 or later
- PostgreSQL 12 or later
- [goose](https://github.com/pressly/goose) for database migrations
- [sqlc](https://sqlc.dev/) for generating type-safe Go from SQL

## Installation

### From Source

```bash
# Install the CLI
go install github.com/frankhuettner/gator-cobra@latest

# Verify installation
gator-cobra --help
```

### Build from Source

```bash
# Clone the repository
git clone https://github.com/frankhuettner/gator-cobra.git
cd gator-cobra

# Build the binary
go build -o bin/gator-cobra

# Optional: Move to PATH
sudo mv bin/gator-cobra /usr/local/bin/
```

## Configuration

1. Create a PostgreSQL database:
```sql
CREATE DATABASE gator;
```

2. Set up your database configuration in `.env`:
```env
DB_URL=postgres://username:password@localhost:5432/gator?sslmode=disable
```

3. Run database migrations:
```bash
goose -dir sql/schema postgres "postgres://username:password@localhost:5432/gator" up
```

## Usage

### User Management

```bash
# Register a new user
gator-cobra register john_doe

# Log in as a user
gator-cobra login john_doe

# List all users
gator-cobra users
```

### Feed Management

```bash
# Add a new feed
gator-cobra addfeed "Tech News" https://example.com/feed.xml

# List all feeds
gator-cobra feeds

# Follow an existing feed
gator-cobra follow https://example.com/feed.xml

# List feeds you're following
gator-cobra following

# Unfollow a feed
gator-cobra unfollow https://example.com/feed.xml
```

### Reading Posts

```bash
# Browse latest posts (default: 2 posts)
gator-cobra browse

# Browse more posts
gator-cobra browse 10

# Start feed aggregation (runs continuously)
gator-cobra agg 1m  # Aggregate every minute
```

## Example Workflow

```bash
# 1. Register and log in
gator-cobra register alice
gator-cobra login alice

# 2. Add some feeds
gator-cobra addfeed "Hacker News" https://news.ycombinator.com/rss
gator-cobra addfeed "Go Blog" https://go.dev/blog/feed.atom

# 3. Start aggregation in one terminal
gator-cobra agg 5m

# 4. Browse posts in another terminal
gator-cobra browse 5
```

## Database Schema

The application uses four main tables:
- `users`: Store user information
- `feeds`: Track RSS feed details
- `feed_follows`: Manage feed subscriptions
- `posts`: Store aggregated posts





## Acknowledgments
- Following boot.dev's [Go tutorial](https://www.boot.dev/)
- Built with [Cobra](https://github.com/spf13/cobra) for CLI functionality
- Database queries generated using [SQLC](https://github.com/kyleconroy/sqlc)
- Database migrations managed with [Goose](https://github.com/pressly/goose)
- Configuration management with [Viper](https://github.com/spf13/viper)
```
