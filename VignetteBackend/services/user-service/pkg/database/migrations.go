package database

import (
	"database/sql"
	"fmt"
)

// RunMigrations runs database migrations
func RunMigrations(db *sql.DB) error {
	// Create users table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			username VARCHAR(30) NOT NULL UNIQUE,
			email VARCHAR(255) NOT NULL UNIQUE,
			full_name VARCHAR(100) NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			phone_number VARCHAR(20),
			bio VARCHAR(150),
			website VARCHAR(255),
			profile_picture_url TEXT,
			is_private BOOLEAN DEFAULT false,
			is_verified BOOLEAN DEFAULT false,
			is_active BOOLEAN DEFAULT true,
			is_deleted BOOLEAN DEFAULT false,
			followers_count INTEGER DEFAULT 0,
			following_count INTEGER DEFAULT 0,
			posts_count INTEGER DEFAULT 0,
			last_login_at TIMESTAMP,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE is_deleted = false;
		CREATE INDEX IF NOT EXISTS idx_users_username ON users(username) WHERE is_deleted = false;
		CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);
	`)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Create sessions table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			access_token TEXT NOT NULL,
			refresh_token TEXT NOT NULL,
			device_info VARCHAR(255),
			ip_address VARCHAR(45),
			user_agent TEXT,
			expires_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			last_active_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
		CREATE INDEX IF NOT EXISTS idx_sessions_access_token ON sessions(access_token);
		CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);
	`)
	if err != nil {
		return fmt.Errorf("failed to create sessions table: %w", err)
	}

	return nil
}
