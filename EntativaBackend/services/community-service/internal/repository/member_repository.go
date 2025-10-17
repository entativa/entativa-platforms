package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/entativa/socialink/community-service/internal/model"
	"github.com/google/uuid"
)

type MemberRepository struct {
	db *sql.DB
}

func NewMemberRepository(db *sql.DB) *MemberRepository {
	return &MemberRepository{db: db}
}

// ============================================
// MEMBERSHIP MANAGEMENT
// ============================================

func (r *MemberRepository) AddMember(ctx context.Context, communityID, userID uuid.UUID, role model.CommunityRole, invitedBy *uuid.UUID) (*model.CommunityMember, error) {
	member := &model.CommunityMember{
		ID:           uuid.New(),
		CommunityID:  communityID,
		UserID:       userID,
		Role:         role,
		Permissions:  model.GetDefaultPermissions(role),
		Status:       model.StatusActive,
		IsMuted:      false,
		PostCount:    0,
		CommentCount: 0,
		InvitedBy:    invitedBy,
		JoinedAt:     time.Now(),
		UpdatedAt:    time.Now(),
	}

	query := `
		INSERT INTO community_members (
			id, community_id, user_id, role, permissions, status, is_muted,
			post_count, comment_count, invited_by, joined_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := r.db.ExecContext(ctx, query,
		member.ID, member.CommunityID, member.UserID, member.Role, member.Permissions,
		member.Status, member.IsMuted, member.PostCount, member.CommentCount,
		member.InvitedBy, member.JoinedAt, member.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add member: %w", err)
	}

	return member, nil
}

func (r *MemberRepository) GetMember(ctx context.Context, communityID, userID uuid.UUID) (*model.CommunityMember, error) {
	query := `
		SELECT id, community_id, user_id, role, permissions, status, is_muted,
		       muted_until, post_count, comment_count, invited_by, joined_at, updated_at
		FROM community_members
		WHERE community_id = $1 AND user_id = $2
	`

	member := &model.CommunityMember{}
	err := r.db.QueryRowContext(ctx, query, communityID, userID).Scan(
		&member.ID, &member.CommunityID, &member.UserID, &member.Role, &member.Permissions,
		&member.Status, &member.IsMuted, &member.MutedUntil, &member.PostCount, &member.CommentCount,
		&member.InvitedBy, &member.JoinedAt, &member.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("member not found")
	}
	if err != nil {
		return nil, err
	}

	return member, nil
}

func (r *MemberRepository) RemoveMember(ctx context.Context, communityID, userID uuid.UUID) error {
	query := `DELETE FROM community_members WHERE community_id = $1 AND user_id = $2`
	result, err := r.db.ExecContext(ctx, query, communityID, userID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("member not found")
	}

	return nil
}

func (r *MemberRepository) UpdateRole(ctx context.Context, communityID, userID uuid.UUID, role model.CommunityRole) error {
	perms := model.GetDefaultPermissions(role)
	query := `
		UPDATE community_members
		SET role = $1, permissions = $2, updated_at = NOW()
		WHERE community_id = $3 AND user_id = $4
	`

	result, err := r.db.ExecContext(ctx, query, role, perms, communityID, userID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("member not found")
	}

	return nil
}

func (r *MemberRepository) UpdatePermissions(ctx context.Context, communityID, userID uuid.UUID, permissions model.CommunityPermissions) error {
	query := `
		UPDATE community_members
		SET permissions = $1, updated_at = NOW()
		WHERE community_id = $2 AND user_id = $3
	`

	result, err := r.db.ExecContext(ctx, query, permissions, communityID, userID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("member not found")
	}

	return nil
}

func (r *MemberRepository) ListMembers(ctx context.Context, communityID uuid.UUID, limit, offset int, role string, status string) ([]*model.CommunityMember, error) {
	query := `
		SELECT id, community_id, user_id, role, permissions, status, is_muted,
		       muted_until, post_count, comment_count, invited_by, joined_at, updated_at
		FROM community_members
		WHERE community_id = $1
	`
	args := []interface{}{communityID}
	argIdx := 2

	if role != "" {
		query += fmt.Sprintf(" AND role = $%d", argIdx)
		args = append(args, role)
		argIdx++
	}

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY joined_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := []*model.CommunityMember{}
	for rows.Next() {
		member := &model.CommunityMember{}
		err := rows.Scan(
			&member.ID, &member.CommunityID, &member.UserID, &member.Role, &member.Permissions,
			&member.Status, &member.IsMuted, &member.MutedUntil, &member.PostCount, &member.CommentCount,
			&member.InvitedBy, &member.JoinedAt, &member.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	return members, nil
}

func (r *MemberRepository) MuteMember(ctx context.Context, communityID, userID uuid.UUID, duration *time.Time) error {
	query := `
		UPDATE community_members
		SET is_muted = TRUE, muted_until = $1, updated_at = NOW()
		WHERE community_id = $2 AND user_id = $3
	`

	_, err := r.db.ExecContext(ctx, query, duration, communityID, userID)
	return err
}

func (r *MemberRepository) UnmuteMember(ctx context.Context, communityID, userID uuid.UUID) error {
	query := `
		UPDATE community_members
		SET is_muted = FALSE, muted_until = NULL, updated_at = NOW()
		WHERE community_id = $2 AND user_id = $3
	`

	_, err := r.db.ExecContext(ctx, query, communityID, userID)
	return err
}

// ============================================
// JOIN REQUESTS
// ============================================

func (r *MemberRepository) CreateJoinRequest(ctx context.Context, communityID, userID uuid.UUID, message string) (*model.JoinRequest, error) {
	req := &model.JoinRequest{
		ID:          uuid.New(),
		CommunityID: communityID,
		UserID:      userID,
		Message:     message,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	query := `
		INSERT INTO join_requests (id, community_id, user_id, message, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query,
		req.ID, req.CommunityID, req.UserID, req.Message, req.Status, req.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (r *MemberRepository) ApproveJoinRequest(ctx context.Context, requestID, reviewerID uuid.UUID) error {
	query := `
		UPDATE join_requests
		SET status = 'approved', reviewed_by = $1, reviewed_at = NOW()
		WHERE id = $2 AND status = 'pending'
	`

	result, err := r.db.ExecContext(ctx, query, reviewerID, requestID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("join request not found or already processed")
	}

	return nil
}

func (r *MemberRepository) RejectJoinRequest(ctx context.Context, requestID, reviewerID uuid.UUID) error {
	query := `
		UPDATE join_requests
		SET status = 'rejected', reviewed_by = $1, reviewed_at = NOW()
		WHERE id = $2 AND status = 'pending'
	`

	result, err := r.db.ExecContext(ctx, query, reviewerID, requestID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("join request not found or already processed")
	}

	return nil
}

func (r *MemberRepository) ListJoinRequests(ctx context.Context, communityID uuid.UUID, status string, limit, offset int) ([]*model.JoinRequest, error) {
	query := `
		SELECT id, community_id, user_id, message, status, reviewed_by, created_at, reviewed_at
		FROM join_requests
		WHERE community_id = $1
	`
	args := []interface{}{communityID}
	argIdx := 2

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requests := []*model.JoinRequest{}
	for rows.Next() {
		req := &model.JoinRequest{}
		err := rows.Scan(
			&req.ID, &req.CommunityID, &req.UserID, &req.Message,
			&req.Status, &req.ReviewedBy, &req.CreatedAt, &req.ReviewedAt,
		)
		if err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}

	return requests, nil
}

// ============================================
// INVITES
// ============================================

func (r *MemberRepository) CreateInvite(ctx context.Context, communityID, invitedUserID, invitedBy uuid.UUID, message *string, expiresAt time.Time) (*model.MemberInvite, error) {
	invite := &model.MemberInvite{
		ID:            uuid.New(),
		CommunityID:   communityID,
		InvitedUserID: invitedUserID,
		InvitedBy:     invitedBy,
		Status:        "pending",
		Message:       message,
		ExpiresAt:     expiresAt,
		CreatedAt:     time.Now(),
	}

	query := `
		INSERT INTO member_invites (id, community_id, invited_user_id, invited_by, status, message, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		invite.ID, invite.CommunityID, invite.InvitedUserID, invite.InvitedBy,
		invite.Status, invite.Message, invite.ExpiresAt, invite.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return invite, nil
}

// ============================================
// BANS
// ============================================

func (r *MemberRepository) BanMember(ctx context.Context, communityID, userID, bannedBy uuid.UUID, reason string, isPermanent bool, expiresAt *time.Time) error {
	// First remove member if exists
	_ = r.RemoveMember(ctx, communityID, userID)

	// Add to banned list
	query := `
		INSERT INTO banned_members (community_id, user_id, banned_by, reason, is_permanent, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (community_id, user_id) DO UPDATE
		SET banned_by = $3, reason = $4, is_permanent = $5, expires_at = $6, created_at = NOW()
	`

	_, err := r.db.ExecContext(ctx, query, communityID, userID, bannedBy, reason, isPermanent, expiresAt)
	return err
}

func (r *MemberRepository) UnbanMember(ctx context.Context, communityID, userID uuid.UUID) error {
	query := `DELETE FROM banned_members WHERE community_id = $1 AND user_id = $2`
	_, err := r.db.ExecContext(ctx, query, communityID, userID)
	return err
}

func (r *MemberRepository) IsBanned(ctx context.Context, communityID, userID uuid.UUID) (bool, error) {
	query := `
		SELECT COUNT(*) FROM banned_members
		WHERE community_id = $1 AND user_id = $2
		  AND (is_permanent OR expires_at > NOW())
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, communityID, userID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
