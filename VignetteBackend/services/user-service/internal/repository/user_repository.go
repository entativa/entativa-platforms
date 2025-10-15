package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"vignette/user-service/internal/model"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrEmailExists    = errors.New("email already exists")
	ErrUsernameExists = errors.New("username already exists")
	ErrDuplicateUser  = errors.New("user with this email or username already exists")
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *model.User) error {
	query := `
		INSERT INTO users (
			id, username, email, full_name, password_hash, 
			is_private, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		user.ID,
		user.Username,
		user.Email,
		user.FullName,
		user.Password,
		user.IsPrivate,
		user.IsActive,
		time.Now(),
		time.Now(),
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		// Check for duplicate email or username
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			return ErrEmailExists
		}
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_username_key\"" {
			return ErrUsernameExists
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(id uuid.UUID) (*model.User, error) {
	query := `
		SELECT id, username, email, full_name, password_hash, phone_number, 
		       bio, website, profile_picture_url, is_private, is_verified, 
		       is_active, is_deleted, followers_count, following_count, 
		       posts_count, last_login_at, created_at, updated_at
		FROM users
		WHERE id = $1 AND is_deleted = false
	`

	user := &model.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FullName,
		&user.Password,
		&user.PhoneNumber,
		&user.Bio,
		&user.Website,
		&user.ProfilePictureURL,
		&user.IsPrivate,
		&user.IsVerified,
		&user.IsActive,
		&user.IsDeleted,
		&user.FollowersCount,
		&user.FollowingCount,
		&user.PostsCount,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	query := `
		SELECT id, username, email, full_name, password_hash, phone_number, 
		       bio, website, profile_picture_url, is_private, is_verified, 
		       is_active, is_deleted, followers_count, following_count, 
		       posts_count, last_login_at, created_at, updated_at
		FROM users
		WHERE email = $1 AND is_deleted = false
	`

	user := &model.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FullName,
		&user.Password,
		&user.PhoneNumber,
		&user.Bio,
		&user.Website,
		&user.ProfilePictureURL,
		&user.IsPrivate,
		&user.IsVerified,
		&user.IsActive,
		&user.IsDeleted,
		&user.FollowersCount,
		&user.FollowingCount,
		&user.PostsCount,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return user, nil
}

// FindByUsername finds a user by username
func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	query := `
		SELECT id, username, email, full_name, password_hash, phone_number, 
		       bio, website, profile_picture_url, is_private, is_verified, 
		       is_active, is_deleted, followers_count, following_count, 
		       posts_count, last_login_at, created_at, updated_at
		FROM users
		WHERE username = $1 AND is_deleted = false
	`

	user := &model.User{}
	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FullName,
		&user.Password,
		&user.PhoneNumber,
		&user.Bio,
		&user.Website,
		&user.ProfilePictureURL,
		&user.IsPrivate,
		&user.IsVerified,
		&user.IsActive,
		&user.IsDeleted,
		&user.FollowersCount,
		&user.FollowingCount,
		&user.PostsCount,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user by username: %w", err)
	}

	return user, nil
}

// Update updates a user
func (r *UserRepository) Update(user *model.User) error {
	query := `
		UPDATE users
		SET full_name = $1, bio = $2, website = $3, phone_number = $4,
		    profile_picture_url = $5, is_private = $6, updated_at = $7
		WHERE id = $8 AND is_deleted = false
	`

	result, err := r.db.Exec(
		query,
		user.FullName,
		user.Bio,
		user.Website,
		user.PhoneNumber,
		user.ProfilePictureURL,
		user.IsPrivate,
		time.Now(),
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// UpdateLastLogin updates the user's last login timestamp
func (r *UserRepository) UpdateLastLogin(userID uuid.UUID) error {
	query := `UPDATE users SET last_login_at = $1 WHERE id = $2`
	_, err := r.db.Exec(query, time.Now(), userID)
	return err
}

// EmailExists checks if an email already exists
func (r *UserRepository) EmailExists(email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND is_deleted = false)`
	var exists bool
	err := r.db.QueryRow(query, email).Scan(&exists)
	return exists, err
}

// UsernameExists checks if a username already exists
func (r *UserRepository) UsernameExists(username string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 AND is_deleted = false)`
	var exists bool
	err := r.db.QueryRow(query, username).Scan(&exists)
	return exists, err
}
