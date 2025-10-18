package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// User represents a user in the database
type User struct {
	ID                string
	FirstName         string
	LastName          string
	Email             string
	Username          string
	PasswordHash      string
	Birthday          *time.Time
	Gender            *string
	PhoneNumber       *string
	Bio               *string
	ProfilePictureURL *string
	CoverPhotoURL     *string
	IsActive          bool
	IsDeleted         bool
	IsVerified        bool
	LastLoginAt       *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// UserRepository handles database operations for users
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser creates a new user
func (r *UserRepository) CreateUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (
			id, first_name, last_name, email, username, password_hash,
			birthday, gender, phone_number, bio, profile_picture_url,
			cover_photo_url, is_active, is_deleted, is_verified,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
		)
	`
	
	_, err := r.db.ExecContext(
		ctx, query,
		user.ID, user.FirstName, user.LastName, user.Email, user.Username,
		user.PasswordHash, user.Birthday, user.Gender, user.PhoneNumber,
		user.Bio, user.ProfilePictureURL, user.CoverPhotoURL,
		user.IsActive, user.IsDeleted, user.IsVerified,
		user.CreatedAt, user.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	
	return nil
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id string) (*User, error) {
	query := `
		SELECT id, first_name, last_name, email, username, password_hash,
		       birthday, gender, phone_number, bio, profile_picture_url,
		       cover_photo_url, is_active, is_deleted, is_verified,
		       last_login_at, created_at, updated_at
		FROM users
		WHERE id = $1 AND is_deleted = false
	`
	
	user := &User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Email,
		&user.Username, &user.PasswordHash, &user.Birthday, &user.Gender,
		&user.PhoneNumber, &user.Bio, &user.ProfilePictureURL,
		&user.CoverPhotoURL, &user.IsActive, &user.IsDeleted, &user.IsVerified,
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	
	return user, nil
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, first_name, last_name, email, username, password_hash,
		       birthday, gender, phone_number, bio, profile_picture_url,
		       cover_photo_url, is_active, is_deleted, is_verified,
		       last_login_at, created_at, updated_at
		FROM users
		WHERE LOWER(email) = LOWER($1) AND is_deleted = false
	`
	
	user := &User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Email,
		&user.Username, &user.PasswordHash, &user.Birthday, &user.Gender,
		&user.PhoneNumber, &user.Bio, &user.ProfilePictureURL,
		&user.CoverPhotoURL, &user.IsActive, &user.IsDeleted, &user.IsVerified,
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	
	return user, nil
}

// FindByUsername finds a user by username
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	query := `
		SELECT id, first_name, last_name, email, username, password_hash,
		       birthday, gender, phone_number, bio, profile_picture_url,
		       cover_photo_url, is_active, is_deleted, is_verified,
		       last_login_at, created_at, updated_at
		FROM users
		WHERE LOWER(username) = LOWER($1) AND is_deleted = false
	`
	
	user := &User{}
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Email,
		&user.Username, &user.PasswordHash, &user.Birthday, &user.Gender,
		&user.PhoneNumber, &user.Bio, &user.ProfilePictureURL,
		&user.CoverPhotoURL, &user.IsActive, &user.IsDeleted, &user.IsVerified,
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	
	return user, nil
}

// FindByEmailOrUsername finds a user by email or username
func (r *UserRepository) FindByEmailOrUsername(ctx context.Context, emailOrUsername string) (*User, error) {
	query := `
		SELECT id, first_name, last_name, email, username, password_hash,
		       birthday, gender, phone_number, bio, profile_picture_url,
		       cover_photo_url, is_active, is_deleted, is_verified,
		       last_login_at, created_at, updated_at
		FROM users
		WHERE (LOWER(email) = LOWER($1) OR LOWER(username) = LOWER($1))
		  AND is_deleted = false
	`
	
	user := &User{}
	err := r.db.QueryRowContext(ctx, query, emailOrUsername).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Email,
		&user.Username, &user.PasswordHash, &user.Birthday, &user.Gender,
		&user.PhoneNumber, &user.Bio, &user.ProfilePictureURL,
		&user.CoverPhotoURL, &user.IsActive, &user.IsDeleted, &user.IsVerified,
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	
	return user, nil
}

// UpdatePassword updates a user's password
func (r *UserRepository) UpdatePassword(ctx context.Context, userID, hashedPassword string) error {
	query := `
		UPDATE users
		SET password_hash = $1, updated_at = $2
		WHERE id = $3
	`
	
	result, err := r.db.ExecContext(ctx, query, hashedPassword, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	
	return nil
}

// UpdateLastLogin updates the last login timestamp
func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	query := `
		UPDATE users
		SET last_login_at = $1, updated_at = $2
		WHERE id = $3
	`
	
	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, now, now, userID)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}
	
	return nil
}

// UpdateUser updates a user's profile information
func (r *UserRepository) UpdateUser(ctx context.Context, user *User) error {
	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, bio = $3,
		    profile_picture_url = $4, cover_photo_url = $5,
		    phone_number = $6, updated_at = $7
		WHERE id = $8
	`
	
	_, err := r.db.ExecContext(
		ctx, query,
		user.FirstName, user.LastName, user.Bio,
		user.ProfilePictureURL, user.CoverPhotoURL,
		user.PhoneNumber, time.Now(), user.ID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	
	return nil
}

// DeleteUser soft deletes a user
func (r *UserRepository) DeleteUser(ctx context.Context, userID string) error {
	query := `
		UPDATE users
		SET is_deleted = true, updated_at = $1
		WHERE id = $2
	`
	
	_, err := r.db.ExecContext(ctx, query, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	
	return nil
}

// LinkCrossPlatformAccount links a cross-platform account
func (r *UserRepository) LinkCrossPlatformAccount(ctx context.Context, userID, platform, platformUserID string) error {
	// This would typically be stored in a separate table
	// For now, we'll use a simple JSON column or separate table
	query := `
		INSERT INTO cross_platform_links (user_id, platform, platform_user_id, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, platform) DO UPDATE
		SET platform_user_id = $3, updated_at = $4
	`
	
	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, userID, platform, platformUserID, now)
	if err != nil {
		return fmt.Errorf("failed to link cross-platform account: %w", err)
	}
	
	return nil
}

// CheckUsernameExists checks if a username is already taken
func (r *UserRepository) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE LOWER(username) = LOWER($1) AND is_deleted = false`
	
	var count int
	err := r.db.QueryRowContext(ctx, query, username).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check username: %w", err)
	}
	
	return count > 0, nil
}

// CheckEmailExists checks if an email is already taken
func (r *UserRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE LOWER(email) = LOWER($1) AND is_deleted = false`
	
	var count int
	err := r.db.QueryRowContext(ctx, query, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check email: %w", err)
	}
	
	return count > 0, nil
}
