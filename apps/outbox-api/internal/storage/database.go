package storage

import (
	"database/sql"
	"fmt"

	"github.com/jared-scarr/portfolio-monorepo/apps/outbox-api/internal/config"
	_ "github.com/lib/pq"
)

// DB wraps a database connection
type DB struct {
	conn *sql.DB
}

// NewDB creates a new database connection
func NewDB(cfg config.DatabaseConfig) (*DB, error) {
	conn, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := &DB{conn: conn}

	// Create tables if they don't exist
	if err := db.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// Conn returns the underlying database connection
func (db *DB) Conn() *sql.DB {
	return db.conn
}

// createTables creates the necessary database tables
func (db *DB) createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS outbox_events (
		id VARCHAR(255) PRIMARY KEY,
		type VARCHAR(255) NOT NULL,
		source VARCHAR(255) NOT NULL,
		data JSONB NOT NULL,
		metadata JSONB,
		status VARCHAR(50) NOT NULL DEFAULT 'pending',
		retry_count INTEGER NOT NULL DEFAULT 0,
		last_error TEXT,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		published_at TIMESTAMP WITH TIME ZONE
	);

	CREATE INDEX IF NOT EXISTS idx_outbox_events_status ON outbox_events(status);
	CREATE INDEX IF NOT EXISTS idx_outbox_events_created_at ON outbox_events(created_at);
	CREATE INDEX IF NOT EXISTS idx_outbox_events_type ON outbox_events(type);
	CREATE INDEX IF NOT EXISTS idx_outbox_events_source ON outbox_events(source);
	`

	_, err := db.conn.Exec(query)
	return err
}
