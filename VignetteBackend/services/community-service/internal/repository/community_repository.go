package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/entativa/vignette/community-service/internal/model"
	"github.com/google/uuid"
)

type CommunityRepository struct {
	db *sql.DB
}

func NewCommunityRepository(db *sql.DB) *CommunityRepository {
	return &CommunityRepository{db: db}
}

// ============================================
// COMMUNITY CRUD
// ============================================

func (r *CommunityRepository) Create(ctx context.Context, req *model.CreateCommunityRequest, creatorID uuid.UUID) (*model.Community, error) {
	community := &model.Community{
		ID:              uuid.New(),
		Name:            req.Name,
		Description:     req.Description,
		CoverPhoto:      req.CoverPhoto,
		Category:        req.Category,
		Privacy:         req.Privacy,
		Visibility:      req.Visibility,
		AllowPosts:      req.AllowPosts,
		RequireApproval: req.RequireApproval,
		CreatorID:       creatorID,
		MemberCount:     1, // Creator is first member
		PostCount:       0,
		Tags:            req.Tags,
		Location:        req.Location,
		Website:         req.Website,
		IsVerified:      false,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	query := `
		INSERT INTO communities (
			id, name, description, cover_photo, category,
			privacy, visibility, is_verified, allow_posts, require_approval,
			creator_id, member_count, post_count, tags, location, website,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`

	_, err := r.db.ExecContext(ctx, query,
		community.ID, community.Name, community.Description, community.CoverPhoto, community.Category,
		community.Privacy, community.Visibility, community.IsVerified, community.AllowPosts, community.RequireApproval,
		community.CreatorID, community.MemberCount, community.PostCount, model.StringArray(community.Tags),
		community.Location, community.Website, community.CreatedAt, community.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create community: %w", err)
	}

	// Add creator as owner
	ownerPerms := model.GetDefaultPermissions(model.RoleOwner)
	memberQuery := `
		INSERT INTO community_members (
			community_id, user_id, role, permissions, status, joined_at
		) VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = r.db.ExecContext(ctx, memberQuery,
		community.ID, creatorID, model.RoleOwner, ownerPerms, model.StatusActive, time.Now(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add creator as member: %w", err)
	}

	return community, nil
}

func (r *CommunityRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Community, error) {
	query := `
		SELECT id, name, description, cover_photo, category, privacy, visibility,
		       is_verified, allow_posts, require_approval, creator_id,
		       member_count, post_count, tags, location, website,
		       created_at, updated_at
		FROM communities
		WHERE id = $1
	`

	community := &model.Community{}
	var tags model.StringArray
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&community.ID, &community.Name, &community.Description, &community.CoverPhoto, &community.Category,
		&community.Privacy, &community.Visibility, &community.IsVerified, &community.AllowPosts,
		&community.RequireApproval, &community.CreatorID, &community.MemberCount, &community.PostCount,
		&tags, &community.Location, &community.Website, &community.CreatedAt, &community.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("community not found")
	}
	if err != nil {
		return nil, err
	}

	community.Tags = tags
	return community, nil
}

func (r *CommunityRepository) Update(ctx context.Context, id uuid.UUID, req *model.UpdateCommunityRequest) error {
	// Build dynamic UPDATE query
	updates := []string{}
	args := []interface{}{}
	argIdx := 1

	if req.Name != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", argIdx))
		args = append(args, *req.Name)
		argIdx++
	}
	if req.Description != nil {
		updates = append(updates, fmt.Sprintf("description = $%d", argIdx))
		args = append(args, *req.Description)
		argIdx++
	}
	if req.CoverPhoto != nil {
		updates = append(updates, fmt.Sprintf("cover_photo = $%d", argIdx))
		args = append(args, *req.CoverPhoto)
		argIdx++
	}
	if req.Category != nil {
		updates = append(updates, fmt.Sprintf("category = $%d", argIdx))
		args = append(args, *req.Category)
		argIdx++
	}
	if req.Privacy != nil {
		updates = append(updates, fmt.Sprintf("privacy = $%d", argIdx))
		args = append(args, *req.Privacy)
		argIdx++
	}
	if req.Visibility != nil {
		updates = append(updates, fmt.Sprintf("visibility = $%d", argIdx))
		args = append(args, *req.Visibility)
		argIdx++
	}
	if req.AllowPosts != nil {
		updates = append(updates, fmt.Sprintf("allow_posts = $%d", argIdx))
		args = append(args, *req.AllowPosts)
		argIdx++
	}
	if req.RequireApproval != nil {
		updates = append(updates, fmt.Sprintf("require_approval = $%d", argIdx))
		args = append(args, *req.RequireApproval)
		argIdx++
	}
	if req.Tags != nil {
		updates = append(updates, fmt.Sprintf("tags = $%d", argIdx))
		args = append(args, model.StringArray(req.Tags))
		argIdx++
	}
	if req.Location != nil {
		updates = append(updates, fmt.Sprintf("location = $%d", argIdx))
		args = append(args, *req.Location)
		argIdx++
	}
	if req.Website != nil {
		updates = append(updates, fmt.Sprintf("website = $%d", argIdx))
		args = append(args, *req.Website)
		argIdx++
	}

	if len(updates) == 0 {
		return nil // Nothing to update
	}

	updates = append(updates, "updated_at = NOW()")
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE communities
		SET %s
		WHERE id = $%d
	`, strings.Join(updates, ", "), argIdx)

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update community: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("community not found")
	}

	return nil
}

func (r *CommunityRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM communities WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("community not found")
	}

	return nil
}

func (r *CommunityRepository) List(ctx context.Context, limit, offset int, category, search string) ([]*model.Community, error) {
	query := `
		SELECT id, name, description, cover_photo, category, privacy, visibility,
		       is_verified, allow_posts, require_approval, creator_id,
		       member_count, post_count, tags, location, website,
		       created_at, updated_at
		FROM communities
		WHERE 1=1
	`
	args := []interface{}{}
	argIdx := 1

	if category != "" {
		query += fmt.Sprintf(" AND category = $%d", argIdx)
		args = append(args, category)
		argIdx++
	}

	if search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", argIdx, argIdx)
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern)
		argIdx++
	}

	// Only show public/listed communities for discovery
	query += " AND (privacy = 'public' OR privacy = 'private') AND visibility = 'listed'"

	query += fmt.Sprintf(" ORDER BY member_count DESC, created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	communities := []*model.Community{}
	for rows.Next() {
		community := &model.Community{}
		var tags model.StringArray
		err := rows.Scan(
			&community.ID, &community.Name, &community.Description, &community.CoverPhoto, &community.Category,
			&community.Privacy, &community.Visibility, &community.IsVerified, &community.AllowPosts,
			&community.RequireApproval, &community.CreatorID, &community.MemberCount, &community.PostCount,
			&tags, &community.Location, &community.Website, &community.Created At, &community.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		community.Tags = tags
		communities = append(communities, community)
	}

	return communities, nil
}

func (r *CommunityRepository) GetByCreator(ctx context.Context, creatorID uuid.UUID, limit, offset int) ([]*model.Community, error) {
	query := `
		SELECT id, name, description, cover_photo, category, privacy, visibility,
		       is_verified, allow_posts, require_approval, creator_id,
		       member_count, post_count, tags, location, website,
		       created_at, updated_at
		FROM communities
		WHERE creator_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, creatorID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	communities := []*model.Community{}
	for rows.Next() {
		community := &model.Community{}
		var tags model.StringArray
		err := rows.Scan(
			&community.ID, &community.Name, &community.Description, &community.CoverPhoto, &community.Category,
			&community.Privacy, &community.Visibility, &community.IsVerified, &community.AllowPosts,
			&community.RequireApproval, &community.CreatorID, &community.MemberCount, &community.PostCount,
			&tags, &community.Location, &community.Website, &community.CreatedAt, &community.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		community.Tags = tags
		communities = append(communities, community)
	}

	return communities, nil
}

// IncrementPostCount increments the post count for a community
func (r *CommunityRepository) IncrementPostCount(ctx context.Context, communityID uuid.UUID) error {
	query := `UPDATE communities SET post_count = post_count + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, communityID)
	return err
}

// DecrementPostCount decrements the post count for a community
func (r *CommunityRepository) DecrementPostCount(ctx context.Context, communityID uuid.UUID) error {
	query := `UPDATE communities SET post_count = GREATEST(0, post_count - 1) WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, communityID)
	return err
}
