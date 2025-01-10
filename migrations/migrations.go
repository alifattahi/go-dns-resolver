package migrations

import (
	"database/sql"
)

// CreateTableIfNotExists ensures the dns_cache table exists
func CreateTableIfNotExists(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS dns_cache (
		domain TEXT PRIMARY KEY,
		ip TEXT NOT NULL,
		dns_provider TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL
	)`
	_, err := db.Exec(query)
	return err
}
