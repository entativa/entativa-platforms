package repository

import (
	"context"
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type SettingsRepository struct {
	db *sql.DB
}

type UserSettings struct {
	// Account
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Bio      string `json:"bio"`
	Website  string `json:"website"`
	Name     string `json:"name"`
	Username string `json:"username"`

	// Privacy
	IsPrivateAccount          bool   `json:"is_private_account"`
	ShowActivityStatus        bool   `json:"show_activity_status"`
	ReadReceipts              bool   `json:"read_receipts"`
	AllowMessageRequests      bool   `json:"allow_message_requests"`
	PostsVisibility           string `json:"posts_visibility"`
	CommentsAllowed           string `json:"comments_allowed"`
	MentionsAllowed           string `json:"mentions_allowed"`
	StorySharing              string `json:"story_sharing"`
	SimilarAccountSuggestions bool   `json:"similar_account_suggestions"`
	IncludeInRecommendations  bool   `json:"include_in_recommendations"`

	// Notifications
	NotifyLikes           bool `json:"notify_likes"`
	NotifyComments        bool `json:"notify_comments"`
	NotifyFollowers       bool `json:"notify_followers"`
	NotifyMessages        bool `json:"notify_messages"`
	NotifyFriendRequests  bool `json:"notify_friend_requests"`
	NotifyVideoViews      bool `json:"notify_video_views"`
	NotifyLiveVideos      bool `json:"notify_live_videos"`
	EmailWeeklySummary    bool `json:"email_weekly_summary"`
	EmailProductUpdates   bool `json:"email_product_updates"`
	EmailTips             bool `json:"email_tips"`
	NotificationSound     bool `json:"notification_sound"`
	NotificationVibration bool `json:"notification_vibration"`
	ShowBadgeCount        bool `json:"show_badge_count"`

	// Data
	UploadQuality    string `json:"upload_quality"`
	AutoplaySettings string `json:"autoplay_settings"`
	DataSaverMode    bool   `json:"data_saver_mode"`
	UseLessData      bool   `json:"use_less_data"`
}

type BlockedUser struct {
	UserID      string    `json:"user_id"`
	Username    string    `json:"username"`
	Name        string    `json:"name"`
	AvatarURL   string    `json:"avatar_url"`
	BlockedAt   time.Time `json:"blocked_at"`
}

type LoginActivity struct {
	ID        string    `json:"id"`
	DeviceName string   `json:"device_name"`
	Location   string   `json:"location"`
	IPAddress  string   `json:"ip_address"`
	LoginAt    time.Time `json:"login_at"`
	IsActive   bool      `json:"is_active"`
}

func NewSettingsRepository(db *sql.DB) *SettingsRepository {
	return &SettingsRepository{db: db}
}

// GetUserSettings retrieves all settings for a user
func (r *SettingsRepository) GetUserSettings(ctx context.Context, userID string) (*UserSettings, error) {
	var settings UserSettings

	query := `
		SELECT 
			u.email, 
			COALESCE(u.phone, '') as phone, 
			COALESCE(u.bio, '') as bio, 
			COALESCE(u.website, '') as website,
			u.name,
			u.username,
			COALESCE(s.is_private_account, false),
			COALESCE(s.show_activity_status, true),
			COALESCE(s.read_receipts, true),
			COALESCE(s.allow_message_requests, true),
			COALESCE(s.posts_visibility, 'everyone'),
			COALESCE(s.comments_allowed, 'everyone'),
			COALESCE(s.mentions_allowed, 'everyone'),
			COALESCE(s.story_sharing, 'everyone'),
			COALESCE(s.similar_account_suggestions, true),
			COALESCE(s.include_in_recommendations, true),
			COALESCE(s.notify_likes, true),
			COALESCE(s.notify_comments, true),
			COALESCE(s.notify_followers, true),
			COALESCE(s.notify_messages, true),
			COALESCE(s.notify_friend_requests, true),
			COALESCE(s.notify_video_views, true),
			COALESCE(s.notify_live_videos, true),
			COALESCE(s.email_weekly_summary, true),
			COALESCE(s.email_product_updates, false),
			COALESCE(s.email_tips, false),
			COALESCE(s.notification_sound, true),
			COALESCE(s.notification_vibration, true),
			COALESCE(s.show_badge_count, true),
			COALESCE(s.upload_quality, 'high'),
			COALESCE(s.autoplay_settings, 'wifi'),
			COALESCE(s.data_saver_mode, false),
			COALESCE(s.use_less_data, false)
		FROM users u
		LEFT JOIN user_settings s ON u.id = s.user_id
		WHERE u.id = $1
	`

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&settings.Email,
		&settings.Phone,
		&settings.Bio,
		&settings.Website,
		&settings.Name,
		&settings.Username,
		&settings.IsPrivateAccount,
		&settings.ShowActivityStatus,
		&settings.ReadReceipts,
		&settings.AllowMessageRequests,
		&settings.PostsVisibility,
		&settings.CommentsAllowed,
		&settings.MentionsAllowed,
		&settings.StorySharing,
		&settings.SimilarAccountSuggestions,
		&settings.IncludeInRecommendations,
		&settings.NotifyLikes,
		&settings.NotifyComments,
		&settings.NotifyFollowers,
		&settings.NotifyMessages,
		&settings.NotifyFriendRequests,
		&settings.NotifyVideoViews,
		&settings.NotifyLiveVideos,
		&settings.EmailWeeklySummary,
		&settings.EmailProductUpdates,
		&settings.EmailTips,
		&settings.NotificationSound,
		&settings.NotificationVibration,
		&settings.ShowBadgeCount,
		&settings.UploadQuality,
		&settings.AutoplaySettings,
		&settings.DataSaverMode,
		&settings.UseLessData,
	)

	if err != nil {
		return nil, err
	}

	return &settings, nil
}

// UsernameExists checks if username is taken by another user
func (r *SettingsRepository) UsernameExists(ctx context.Context, username, excludeUserID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 AND id != $2)`
	err := r.db.QueryRowContext(ctx, query, username, excludeUserID).Scan(&exists)
	return exists, err
}

// UpdateAccountSettings updates account information
func (r *SettingsRepository) UpdateAccountSettings(ctx context.Context, userID string, req interface{}) error {
	// Type assertion to get the request data
	reqMap := req.(*struct {
		Email    string `json:"email,omitempty"`
		Phone    string `json:"phone,omitempty"`
		Bio      string `json:"bio,omitempty"`
		Website  string `json:"website,omitempty"`
		Name     string `json:"name,omitempty"`
		Username string `json:"username,omitempty"`
	})

	query := `
		UPDATE users 
		SET 
			email = COALESCE(NULLIF($2, ''), email),
			phone = COALESCE(NULLIF($3, ''), phone),
			bio = COALESCE(NULLIF($4, ''), bio),
			website = COALESCE(NULLIF($5, ''), website),
			name = COALESCE(NULLIF($6, ''), name),
			username = COALESCE(NULLIF($7, ''), username),
			updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, userID, reqMap.Email, reqMap.Phone, reqMap.Bio, reqMap.Website, reqMap.Name, reqMap.Username)
	return err
}

// UpdatePrivacySettings updates privacy settings
func (r *SettingsRepository) UpdatePrivacySettings(ctx context.Context, userID string, req interface{}) error {
	reqMap := req.(*struct {
		IsPrivateAccount          bool   `json:"is_private_account"`
		ShowActivityStatus        bool   `json:"show_activity_status"`
		ReadReceipts              bool   `json:"read_receipts"`
		AllowMessageRequests      bool   `json:"allow_message_requests"`
		PostsVisibility           string `json:"posts_visibility"`
		CommentsAllowed           string `json:"comments_allowed"`
		MentionsAllowed           string `json:"mentions_allowed"`
		StorySharing              string `json:"story_sharing"`
		SimilarAccountSuggestions bool   `json:"similar_account_suggestions"`
		IncludeInRecommendations  bool   `json:"include_in_recommendations"`
	})

	// Upsert settings
	query := `
		INSERT INTO user_settings (
			user_id, 
			is_private_account, 
			show_activity_status, 
			read_receipts, 
			allow_message_requests,
			posts_visibility,
			comments_allowed,
			mentions_allowed,
			story_sharing,
			similar_account_suggestions,
			include_in_recommendations
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (user_id) 
		DO UPDATE SET
			is_private_account = $2,
			show_activity_status = $3,
			read_receipts = $4,
			allow_message_requests = $5,
			posts_visibility = $6,
			comments_allowed = $7,
			mentions_allowed = $8,
			story_sharing = $9,
			similar_account_suggestions = $10,
			include_in_recommendations = $11
	`

	_, err := r.db.ExecContext(ctx, query,
		userID,
		reqMap.IsPrivateAccount,
		reqMap.ShowActivityStatus,
		reqMap.ReadReceipts,
		reqMap.AllowMessageRequests,
		reqMap.PostsVisibility,
		reqMap.CommentsAllowed,
		reqMap.MentionsAllowed,
		reqMap.StorySharing,
		reqMap.SimilarAccountSuggestions,
		reqMap.IncludeInRecommendations,
	)

	return err
}

// UpdateNotificationSettings updates notification preferences
func (r *SettingsRepository) UpdateNotificationSettings(ctx context.Context, userID string, req interface{}) error {
	reqMap := req.(*struct {
		NotifyLikes           bool `json:"notify_likes"`
		NotifyComments        bool `json:"notify_comments"`
		NotifyFollowers       bool `json:"notify_followers"`
		NotifyMessages        bool `json:"notify_messages"`
		NotifyFriendRequests  bool `json:"notify_friend_requests"`
		NotifyVideoViews      bool `json:"notify_video_views"`
		NotifyLiveVideos      bool `json:"notify_live_videos"`
		EmailWeeklySummary    bool `json:"email_weekly_summary"`
		EmailProductUpdates   bool `json:"email_product_updates"`
		EmailTips             bool `json:"email_tips"`
		NotificationSound     bool `json:"notification_sound"`
		NotificationVibration bool `json:"notification_vibration"`
		ShowBadgeCount        bool `json:"show_badge_count"`
	})

	query := `
		INSERT INTO user_settings (
			user_id,
			notify_likes,
			notify_comments,
			notify_followers,
			notify_messages,
			notify_friend_requests,
			notify_video_views,
			notify_live_videos,
			email_weekly_summary,
			email_product_updates,
			email_tips,
			notification_sound,
			notification_vibration,
			show_badge_count
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		ON CONFLICT (user_id)
		DO UPDATE SET
			notify_likes = $2,
			notify_comments = $3,
			notify_followers = $4,
			notify_messages = $5,
			notify_friend_requests = $6,
			notify_video_views = $7,
			notify_live_videos = $8,
			email_weekly_summary = $9,
			email_product_updates = $10,
			email_tips = $11,
			notification_sound = $12,
			notification_vibration = $13,
			show_badge_count = $14
	`

	_, err := r.db.ExecContext(ctx, query,
		userID,
		reqMap.NotifyLikes,
		reqMap.NotifyComments,
		reqMap.NotifyFollowers,
		reqMap.NotifyMessages,
		reqMap.NotifyFriendRequests,
		reqMap.NotifyVideoViews,
		reqMap.NotifyLiveVideos,
		reqMap.EmailWeeklySummary,
		reqMap.EmailProductUpdates,
		reqMap.EmailTips,
		reqMap.NotificationSound,
		reqMap.NotificationVibration,
		reqMap.ShowBadgeCount,
	)

	return err
}

// UpdateDataSettings updates data usage preferences
func (r *SettingsRepository) UpdateDataSettings(ctx context.Context, userID string, req interface{}) error {
	reqMap := req.(*struct {
		UploadQuality    string `json:"upload_quality"`
		AutoplaySettings string `json:"autoplay_settings"`
		DataSaverMode    bool   `json:"data_saver_mode"`
		UseLessData      bool   `json:"use_less_data"`
	})

	query := `
		INSERT INTO user_settings (
			user_id,
			upload_quality,
			autoplay_settings,
			data_saver_mode,
			use_less_data
		) VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id)
		DO UPDATE SET
			upload_quality = $2,
			autoplay_settings = $3,
			data_saver_mode = $4,
			use_less_data = $5
	`

	_, err := r.db.ExecContext(ctx, query,
		userID,
		reqMap.UploadQuality,
		reqMap.AutoplaySettings,
		reqMap.DataSaverMode,
		reqMap.UseLessData,
	)

	return err
}

// GetBlockedUsers retrieves list of blocked users
func (r *SettingsRepository) GetBlockedUsers(ctx context.Context, userID string) ([]BlockedUser, error) {
	query := `
		SELECT 
			u.id,
			u.username,
			u.name,
			COALESCE(u.avatar_url, '') as avatar_url,
			b.blocked_at
		FROM blocked_users b
		JOIN users u ON b.blocked_user_id = u.id
		WHERE b.user_id = $1
		ORDER BY b.blocked_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blockedUsers []BlockedUser
	for rows.Next() {
		var bu BlockedUser
		if err := rows.Scan(&bu.UserID, &bu.Username, &bu.Name, &bu.AvatarURL, &bu.BlockedAt); err != nil {
			return nil, err
		}
		blockedUsers = append(blockedUsers, bu)
	}

	return blockedUsers, nil
}

// BlockUser blocks a user
func (r *SettingsRepository) BlockUser(ctx context.Context, userID, targetUserID string) error {
	query := `INSERT INTO blocked_users (user_id, blocked_user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, userID, targetUserID)
	return err
}

// UnblockUser unblocks a user
func (r *SettingsRepository) UnblockUser(ctx context.Context, userID, targetUserID string) error {
	query := `DELETE FROM blocked_users WHERE user_id = $1 AND blocked_user_id = $2`
	_, err := r.db.ExecContext(ctx, query, userID, targetUserID)
	return err
}

// VerifyPassword verifies the user's current password
func (r *SettingsRepository) VerifyPassword(ctx context.Context, userID, password string) (bool, error) {
	var hashedPassword string
	query := `SELECT password_hash FROM users WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&hashedPassword)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil, nil
}

// ChangePassword changes the user's password
func (r *SettingsRepository) ChangePassword(ctx context.Context, userID, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `UPDATE users SET password_hash = $2, updated_at = NOW() WHERE id = $1`
	_, err = r.db.ExecContext(ctx, query, userID, string(hashedPassword))
	return err
}

// DeleteAccount soft deletes a user account
func (r *SettingsRepository) DeleteAccount(ctx context.Context, userID, reason string) error {
	query := `
		UPDATE users 
		SET 
			deleted_at = NOW(),
			deletion_reason = $2,
			updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, userID, reason)
	return err
}

// GetLoginActivity retrieves recent login activity
func (r *SettingsRepository) GetLoginActivity(ctx context.Context, userID string, limit int) ([]LoginActivity, error) {
	query := `
		SELECT 
			id,
			COALESCE(device_name, 'Unknown Device') as device_name,
			COALESCE(location, 'Unknown') as location,
			ip_address,
			login_at,
			(expires_at > NOW()) as is_active
		FROM sessions
		WHERE user_id = $1
		ORDER BY login_at DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []LoginActivity
	for rows.Next() {
		var la LoginActivity
		if err := rows.Scan(&la.ID, &la.DeviceName, &la.Location, &la.IPAddress, &la.LoginAt, &la.IsActive); err != nil {
			return nil, err
		}
		activities = append(activities, la)
	}

	return activities, nil
}
