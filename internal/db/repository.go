package db

import (
	"database/sql"
	"fmt"
)

// Holding represents a single portfolio position.
type Holding struct {
	ID      int64
	ChatID  int64
	Symbol  string
	Name    string
	Shares  float64
}

// Repository provides CRUD operations for users and holdings.
type Repository struct {
	db *sql.DB
}

// NewRepository creates a Repository backed by the given DB.
func NewRepository(database *DB) *Repository {
	return &Repository{db: database.DB}
}

// UpsertUser inserts or updates a user record.
func (r *Repository) UpsertUser(chatID int64, username string) error {
	_, err := r.db.Exec(`
		INSERT INTO users (chat_id, username)
		VALUES (?, ?)
		ON CONFLICT(chat_id) DO UPDATE SET username = excluded.username`,
		chatID, username,
	)
	return err
}

// SetUserState updates the FSM state and optional JSON payload for a user.
func (r *Repository) SetUserState(chatID int64, state, stateData string) error {
	_, err := r.db.Exec(`
		UPDATE users SET state = ?, state_data = ? WHERE chat_id = ?`,
		state, stateData, chatID,
	)
	return err
}

// GetUserState returns the current FSM state and payload for a user.
func (r *Repository) GetUserState(chatID int64) (state, stateData string, err error) {
	row := r.db.QueryRow(`
		SELECT state, state_data FROM users WHERE chat_id = ?`, chatID)
	err = row.Scan(&state, &stateData)
	if err == sql.ErrNoRows {
		return "idle", "", nil
	}
	return
}

// UpsertHolding inserts or updates a holding (updates shares on conflict).
func (r *Repository) UpsertHolding(chatID int64, symbol, name string, shares float64) error {
	_, err := r.db.Exec(`
		INSERT INTO holdings (chat_id, symbol, name, shares)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(chat_id, symbol) DO UPDATE SET
			name   = excluded.name,
			shares = excluded.shares`,
		chatID, symbol, name, shares,
	)
	return err
}

// GetHoldings returns all holdings for a user.
func (r *Repository) GetHoldings(chatID int64) ([]Holding, error) {
	rows, err := r.db.Query(`
		SELECT id, chat_id, symbol, name, shares
		FROM holdings WHERE chat_id = ?
		ORDER BY symbol`, chatID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var holdings []Holding
	for rows.Next() {
		var h Holding
		if err := rows.Scan(&h.ID, &h.ChatID, &h.Symbol, &h.Name, &h.Shares); err != nil {
			return nil, err
		}
		holdings = append(holdings, h)
	}
	return holdings, rows.Err()
}

// DeleteHolding removes a specific holding for a user.
func (r *Repository) DeleteHolding(chatID int64, symbol string) error {
	_, err := r.db.Exec(`
		DELETE FROM holdings WHERE chat_id = ? AND symbol = ?`, chatID, symbol)
	return err
}

// GetAllActiveUsers returns chat IDs of all users who have at least one holding.
func (r *Repository) GetAllActiveUsers() ([]int64, error) {
	rows, err := r.db.Query(`
		SELECT DISTINCT chat_id FROM holdings`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// SaveReport records a balance report in the history table.
func (r *Repository) SaveReport(chatID int64, totalUSD float64) error {
	_, err := r.db.Exec(`
		INSERT INTO history (chat_id, total_usd) VALUES (?, ?)`,
		chatID, totalUSD,
	)
	return err
}

// GetLastReport returns the most recent historical total for a user.
// Returns 0, nil if no previous report exists.
func (r *Repository) GetLastReport(chatID int64) (float64, error) {
	var total float64
	err := r.db.QueryRow(`
		SELECT total_usd FROM history
		WHERE chat_id = ?
		ORDER BY reported_at DESC
		LIMIT 1`, chatID).Scan(&total)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return total, err
}

// GetDistinctSymbols returns all unique ticker symbols across all users.
func (r *Repository) GetDistinctSymbols() ([]string, error) {
	rows, err := r.db.Query(`
		SELECT DISTINCT symbol FROM holdings ORDER BY symbol`)
	if err != nil {
		return nil, fmt.Errorf("query distinct symbols: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var symbols []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		symbols = append(symbols, s)
	}
	return symbols, rows.Err()
}
