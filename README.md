# Gator-Cobra: RSS Feed Aggregator

A modern command-line RSS feed aggregator built in Go that helps users follow and aggregate content from multiple RSS feeds in one place.

## Features

- 👤 **User Management**: Create and manage user accounts
- 📰 **Feed Management**: Add and track RSS feeds
- ✅ **Feed Following**: Subscribe/unsubscribe to feeds of interest
- 🔄 **Automated Updates**: Regular feed fetching and post aggregation
- 🎯 **Post Aggregation**: View posts from all followed feeds in one place

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/gator-cobra.git

# Navigate to project directory
cd gator-cobra

# Install dependencies
go mod download

# Build the project
go build
```

## Configuration

Create a `.env` file in the project root:

```env
PORT=8080
DB_URL=postgres://username:password@localhost:5432/gator_cobra
```

## Database Setup

The application uses PostgreSQL. Run migrations to set up the database:

```bash
goose -dir sql/schema postgres "postgres://username:password@localhost:5432/gator_cobra" up
```

## Usage

### User Management

```bash
# Create a new user
gator-cobra user create --name john_doe

# List all users
gator-cobra user list
```

### Feed Management

```bash
# Add a new feed
gator-cobra feed create --name "Tech News" --url https://example.com/feed.xml

# List all feeds
gator-cobra feed list

# Follow a feed
gator-cobra feed follow --id <feed_id>

# Unfollow a feed
gator-cobra feed unfollow --id <feed_id>
```

### Post Management

```bash
# List posts from followed feeds
gator-cobra post list --limit 10
```

## Database Schema

The application uses four main tables:

- `users`: Store user information
- `feeds`: Track RSS feed details
- `feed_follows`: Manage feed subscriptions
- `posts`: Store aggregated posts

## Technology Stack

- **Language**: Go
- **CLI Framework**: Cobra
- **Database**: PostgreSQL
- **Migration Tool**: Goose
- **Query Builder**: SQLC

## Project Structure

```
gator-cobra/
├── cmd/            # CLI commands
├── sql/
│   ├── queries/    # SQLC queries
│   └── schema/     # Database migrations
├── internal/       # Internal packages
└── main.go        # Application entry point
```

## Development

### Adding New Commands

```bash
# Add a new command
cobra-cli add [command]

# Add a subcommand
cobra-cli add [subcommand] -p '[parentCmd]Cmd'
```

### Database Migrations

```bash
# Create a new migration
goose create add_new_feature sql

# Run migrations
goose up

# Rollback last migration
goose down
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra)
- Database queries generated using [SQLC](https://github.com/kyleconroy/sqlc)
- Database migrations managed with [Goose](https://github.com/pressly/goose)
