package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DB struct {
    *sql.DB
    *Queries
}

func Connect(dbURL string) (*DB, error) {
    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        return nil, fmt.Errorf("error connecting to db: %w", err)
    }

    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("error pinging db: %w", err)
    }

    return &DB{
        DB:      db,
        Queries: New(db),
    }, nil
}
