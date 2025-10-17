package service

import (
	"context"
	"fmt"
	"time"

	"github.com/entativa/socialink/community-service/internal/model"
	"github.com/entativa/socialink/community-service/internal/repository"
	"github.com/google/uuid"
)

type CommunityService struct {
	communityRepo *repository.CommunityRepository
	memberRepo    *repository.MemberRepository
}

func NewCommunityService(
	communityRepo *repository.CommunityRepository,
	memberRepo    *repository.MemberRepository,
) *CommunityService {
	return &CommunityService{
		communityRepo: communityRepo,
		memberRepo:    memberRepo,
	}
}

// ============================================
// COMMUNITY MANAGEMENT
// ============================================

func (s *CommunityService) CreateCommunity(ctx context.Context, req *model.CreateCommunityRequest, creatorID uuid.UUID) (*model.Community, error) {
	// Validate request
	if req.Name == "" || req.Category == "" {
		return nil, fmt.Errorf("name and category are required")
	}

	// Create community
	community, err := s.communityRepo.Create(ctx, req, creatorID)
	if err != nil {
		return nil, err
	}

	return community, nil
}

func (s *CommunityService) GetCommunity(ctx context.Context, communityID uuid.UUID, userID *uuid.UUID) (*model.Community, *model.CommunityMember, error) {
	// Get community
	community, err := s.communityRepo.GetByID(ctx, communityID)
	if err != nil {
		return nil, nil, err
	}

	// Get user's membership if authenticated
	var member *model.CommunityMember
	if userID != nil {
		member, _ = s.memberRepo.GetMember(ctx, communityID, *userID)
	}

	// Check visibility
	if community.Privacy == model.PrivacyHidden {
		if member == nil {
			return nil, nil, fmt.Errorf("community not found")
		}
	}

	return community, member, nil
}

func (s *CommunityService) UpdateCommunity(ctx context.Context, communityID, userID uuid.UUID, req *model.UpdateCommunityRequest) error {
	// Check permissions
	member, err := s.memberRepo.GetMember(ctx, communityID, userID)
	if err != nil {
		return fmt.Errorf("not a member")
	}

	if !member.Permissions.CanEditCommunity {
		return fmt.Errorf("insufficient permissions")
	}

	// Update community
	err = s.communityRepo.Update(ctx, communityID, req)
	if err != nil {
		return err
	}

	return nil
}

func (s *CommunityService) DeleteCommunity(ctx context.Context, communityID, userID uuid.UUID) error {
	// Check permissions (owner only)
	member, err := s.memberRepo.GetMember(ctx, communityID, userID)
	if err != nil {
		return fmt.Errorf("not a member")
	}

	if !member.Permissions.CanDeleteCommunity {
		return fmt.Errorf("only owner can delete community")
	}

	// Delete community (CASCADE will handle members, rules, etc.)
	err = s.communityRepo.Delete(ctx, communityID)
	if err != nil {
		return err
	}

	return nil
}

// ============================================
// MEMBERSHIP
// ============================================

func (s *CommunityService) JoinCommunity(ctx context.Context, communityID, userID uuid.UUID) error {
	// Get community
	community, err := s.communityRepo.GetByID(ctx, communityID)
	if err != nil {
		return err
	}

	// Check if already banned
	banned, err := s.memberRepo.IsBanned(ctx, communityID, userID)
	if err != nil {
		return err
	}
	if banned {
		return fmt.Errorf("you are banned from this community")
	}

	// Check privacy
	if community.Privacy == model.PrivacyPublic {
		// Public: Auto-join
		_, err = s.memberRepo.AddMember(ctx, communityID, userID, model.RoleMember, nil)
		return err
	} else {
		// Private/Hidden: Create join request
		_, err = s.memberRepo.CreateJoinRequest(ctx, communityID, userID, "")
		return err
	}
}

func (s *CommunityService) LeaveCommunity(ctx context.Context, communityID, userID uuid.UUID) error {
	// Check if owner
	member, err := s.memberRepo.GetMember(ctx, communityID, userID)
	if err != nil {
		return err
	}

	if member.Role == model.RoleOwner {
		return fmt.Errorf("owner cannot leave community (transfer ownership first)")
	}

	// Remove member
	err = s.memberRepo.RemoveMember(ctx, communityID, userID)
	return err
}

func (s *CommunityService) InviteMember(ctx context.Context, communityID, inviterID, invitedUserID uuid.UUID, message *string) error {
	// Check permissions
	member, err := s.memberRepo.GetMember(ctx, communityID, inviterID)
	if err != nil {
		return fmt.Errorf("not a member")
	}

	if !member.Permissions.CanInviteMembers {
		return fmt.Errorf("insufficient permissions")
	}

	// Check if already banned
	banned, err := s.memberRepo.IsBanned(ctx, communityID, invitedUserID)
	if err != nil {
		return err
	}
	if banned {
		return fmt.Errorf("user is banned from this community")
	}

	// Create invite (expires in 7 days)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	_, err = s.memberRepo.CreateInvite(ctx, communityID, invitedUserID, inviterID, message, expiresAt)
	return err
}

func (s *CommunityService) RemoveMember(ctx context.Context, communityID, removerID, targetUserID uuid.UUID) error {
	// Check permissions
	member, err := s.memberRepo.GetMember(ctx, communityID, removerID)
	if err != nil {
		return fmt.Errorf("not a member")
	}

	if !member.Permissions.CanRemoveMembers {
		return fmt.Errorf("insufficient permissions")
	}

	// Get target member
	targetMember, err := s.memberRepo.GetMember(ctx, communityID, targetUserID)
	if err != nil {
		return err
	}

	// Cannot remove owner
	if targetMember.Role == model.RoleOwner {
		return fmt.Errorf("cannot remove owner")
	}

	// Remove member
	err = s.memberRepo.RemoveMember(ctx, communityID, targetUserID)
	return err
}

// ============================================
// ROLES & PERMISSIONS
// ============================================

func (s *CommunityService) UpdateMemberRole(ctx context.Context, communityID, updaterID, targetUserID uuid.UUID, newRole model.CommunityRole) error {
	// Check permissions
	member, err := s.memberRepo.GetMember(ctx, communityID, updaterID)
	if err != nil {
		return fmt.Errorf("not a member")
	}

	if !member.Permissions.CanManageRoles {
		return fmt.Errorf("insufficient permissions")
	}

	// Cannot make someone owner (must transfer ownership)
	if newRole == model.RoleOwner {
		return fmt.Errorf("use transfer ownership to make someone owner")
	}

	// Update role
	err = s.memberRepo.UpdateRole(ctx, communityID, targetUserID, newRole)
	return err
}

func (s *CommunityService) UpdateMemberPermissions(ctx context.Context, communityID, updaterID, targetUserID uuid.UUID, permissions model.CommunityPermissions) error {
	// Check permissions (admins+ only)
	member, err := s.memberRepo.GetMember(ctx, communityID, updaterID)
	if err != nil {
		return fmt.Errorf("not a member")
	}

	if !member.Permissions.CanManageRoles {
		return fmt.Errorf("insufficient permissions")
	}

	// Get target member
	targetMember, err := s.memberRepo.GetMember(ctx, communityID, targetUserID)
	if err != nil {
		return err
	}

	// Cannot modify owner permissions
	if targetMember.Role == model.RoleOwner {
		return fmt.Errorf("cannot modify owner permissions")
	}

	// Update permissions
	err = s.memberRepo.UpdatePermissions(ctx, communityID, targetUserID, permissions)
	return err
}

// ============================================
// MODERATION
// ============================================

func (s *CommunityService) BanMember(ctx context.Context, communityID, moderatorID, targetUserID uuid.UUID, reason string, durationDays int) error {
	// Check permissions
	member, err := s.memberRepo.GetMember(ctx, communityID, moderatorID)
	if err != nil {
		return fmt.Errorf("not a member")
	}

	if !member.Permissions.CanBanMembers {
		return fmt.Errorf("insufficient permissions")
	}

	// Get target member
	targetMember, err := s.memberRepo.GetMember(ctx, communityID, targetUserID)
	if err != nil {
		return err
	}

	// Cannot ban owner/admin
	if targetMember.Role == model.RoleOwner || targetMember.Role == model.RoleAdmin {
		return fmt.Errorf("cannot ban owner or admin")
	}

	// Calculate expiry
	var expiresAt *time.Time
	isPermanent := durationDays == 0
	if !isPermanent {
		exp := time.Now().Add(time.Duration(durationDays) * 24 * time.Hour)
		expiresAt = &exp
	}

	// Ban member
	err = s.memberRepo.BanMember(ctx, communityID, targetUserID, moderatorID, reason, isPermanent, expiresAt)
	return err
}

func (s *CommunityService) UnbanMember(ctx context.Context, communityID, moderatorID, targetUserID uuid.UUID) error {
	// Check permissions
	member, err := s.memberRepo.GetMember(ctx, communityID, moderatorID)
	if err != nil {
		return fmt.Errorf("not a member")
	}

	if !member.Permissions.CanBanMembers {
		return fmt.Errorf("insufficient permissions")
	}

	// Unban member
	err = s.memberRepo.UnbanMember(ctx, communityID, targetUserID)
	return err
}

func (s *CommunityService) MuteMember(ctx context.Context, communityID, moderatorID, targetUserID uuid.UUID, durationHours int) error {
	// Check permissions
	member, err := s.memberRepo.GetMember(ctx, communityID, moderatorID)
	if err != nil {
		return fmt.Errorf("not a member")
	}

	if !member.Permissions.CanMuteMembers {
		return fmt.Errorf("insufficient permissions")
	}

	// Calculate mute duration
	mutedUntil := time.Now().Add(time.Duration(durationHours) * time.Hour)

	// Mute member
	err = s.memberRepo.MuteMember(ctx, communityID, targetUserID, &mutedUntil)
	return err
}

func (s *CommunityService) UnmuteMember(ctx context.Context, communityID, moderatorID, targetUserID uuid.UUID) error {
	// Check permissions
	member, err := s.memberRepo.GetMember(ctx, communityID, moderatorID)
	if err != nil {
		return fmt.Errorf("not a member")
	}

	if !member.Permissions.CanMuteMembers {
		return fmt.Errorf("insufficient permissions")
	}

	// Unmute member
	err = s.memberRepo.UnmuteMember(ctx, communityID, targetUserID)
	return err
}

// ============================================
// UTILITY
// ============================================

func (s *CommunityService) CheckPermission(ctx context.Context, communityID, userID uuid.UUID, permission string) (bool, error) {
	member, err := s.memberRepo.GetMember(ctx, communityID, userID)
	if err != nil {
		return false, err
	}

	// Check specific permission
	perms := member.Permissions
	switch permission {
	case "post":
		return perms.CanPost, nil
	case "comment":
		return perms.CanComment, nil
	case "moderate":
		return perms.CanModerate, nil
	case "ban":
		return perms.CanBanMembers, nil
	case "edit":
		return perms.CanEditCommunity, nil
	default:
		return false, fmt.Errorf("unknown permission")
	}
}
