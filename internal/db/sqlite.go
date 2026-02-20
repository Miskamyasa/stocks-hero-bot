package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS users (
    chat_id     INTEGER PRIMARY KEY,
    username    TEXT,
    state       TEXT DEFAULT 'idle',
    state_data  TEXT DEFAULT '',
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS holdings (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    chat_id     INTEGER NOT NULL REFERENCES users(chat_id),
    symbol      TEXT NOT NULL,
    name        TEXT NOT NULL,
    shares      REAL NOT NULL,
    added_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(chat_id, symbol)
);

CREATE INDEX IF NOT EXISTS idx_holdings_chat ON holdings(chat_id);
`

// DB wraps a sql.DB with SQLite-specific setup.
type DB struct {
	*sql.DB
}

// New opens (or creates) the SQLite database at path and runs migrations.
func New(path string) (*DB, error) {
	sqlDB, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	// SQLite performs best with a single writer connection.
	sqlDB.SetMaxOpenConns(1)

	if _, err := sqlDB.Exec(schema); err != nil {
		return nil, fmt.Errorf("migrate schema: %w", err)
	}

	return &DB{sqlDB}, nil
}
