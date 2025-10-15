package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"socialink/meta-user-service/internal/model"
	"socialink/meta-user-service/pkg/cache"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

var (
	ErrMetaUserNotFound      = errors.New("meta user not found")
	ErrMetaUserAlreadyExists = errors.New("meta user already exists")
	ErrDuplicateMetaID       = errors.New("duplicate meta ID")
)

// MetaUserRepository provides advanced data access patterns with event sourcing
type MetaUserRepository struct {
	db    *sql.DB
	cache cache.Cache
}

func NewMetaUserRepository(db *sql.DB, cache cache.Cache) *MetaUserRepository {
	return &MetaUserRepository{
		db:    db,
		cache: cache,
	}
}

// Create creates a new meta user with event sourcing
func (r *MetaUserRepository) Create(ctx context.Context, user *model.MetaUser) error {
	query := `
		INSERT INTO meta_users (
			id, meta_id, email, email_verified, phone_number, phone_verified,
			password_hash, status, trust_score, risk_level, account_tier,
			platform_links, security_profile, privacy_settings, data_rights,
			device_fingerprints, biometric_tokens, federated_identities,
			session_management, anomaly_detection, cross_platform_activity,
			compliance_data, metadata, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23, $24, $25
		)
	`

	platformLinksJSON, _ := json.Marshal(user.PlatformLinks)
	securityProfileJSON, _ := json.Marshal(user.SecurityProfile)
	privacySettingsJSON, _ := json.Marshal(user.PrivacySettings)
	dataRightsJSON, _ := json.Marshal(user.DataRights)
	deviceFingerprintsJSON, _ := json.Marshal(user.DeviceFingerprints)
	biometricTokensJSON, _ := json.Marshal(user.BiometricTokens)
	federatedIdentitiesJSON, _ := json.Marshal(user.FederatedIdentities)
	sessionManagementJSON, _ := json.Marshal(user.SessionManagement)
	anomalyDetectionJSON, _ := json.Marshal(user.AnomalyDetection)
	crossPlatformActivityJSON, _ := json.Marshal(user.CrossPlatformActivity)
	complianceDataJSON, _ := json.Marshal(user.ComplianceData)
	metadataJSON, _ := json.Marshal(user.Metadata)

	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.MetaID, user.Email, user.EmailVerified, user.PhoneNumber,
		user.PhoneVerified, user.PasswordHash, user.Status, user.TrustScore,
		user.RiskLevel, user.AccountTier, platformLinksJSON, securityProfileJSON,
		privacySettingsJSON, dataRightsJSON, deviceFingerprintsJSON, biometricTokensJSON,
		federatedIdentitiesJSON, sessionManagementJSON, anomalyDetectionJSON,
		crossPlatformActivityJSON, complianceDataJSON, metadataJSON,
		user.CreatedAt, user.UpdatedAt,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrMetaUserAlreadyExists
		}
		return fmt.Errorf("failed to create meta user: %w", err)
	}

	// Create event sourcing entry
	r.createEvent(ctx, user.ID, "user.created", map[string]interface{}{
		"meta_id": user.MetaID,
		"email":   user.Email,
	})

	// Invalidate cache
	r.cache.Delete(ctx, fmt.Sprintf("meta_user:%s", user.ID))
	r.cache.Delete(ctx, fmt.Sprintf("meta_user:meta_id:%s", user.MetaID))

	return nil
}

// GetByID retrieves a meta user by ID with multi-level caching
func (r *MetaUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.MetaUser, error) {
	// Check L1 cache (Redis)
	cacheKey := fmt.Sprintf("meta_user:%s", id)
	var user model.MetaUser
	
	if err := r.cache.Get(ctx, cacheKey, &user); err == nil {
		return &user, nil
	}

	// Query database
	query := `
		SELECT 
			id, meta_id, email, email_verified, phone_number, phone_verified,
			password_hash, status, trust_score, risk_level, account_tier,
			platform_links, security_profile, privacy_settings, data_rights,
			device_fingerprints, biometric_tokens, federated_identities,
			session_management, anomaly_detection, cross_platform_activity,
			compliance_data, metadata, created_at, updated_at, last_seen_at,
			deleted_at
		FROM meta_users
		WHERE id = $1 AND deleted_at IS NULL
	`

	var platformLinksJSON, securityProfileJSON, privacySettingsJSON, dataRightsJSON []byte
	var deviceFingerprintsJSON, biometricTokensJSON, federatedIdentitiesJSON []byte
	var sessionManagementJSON, anomalyDetectionJSON, crossPlatformActivityJSON []byte
	var complianceDataJSON, metadataJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.MetaID, &user.Email, &user.EmailVerified, &user.PhoneNumber,
		&user.PhoneVerified, &user.PasswordHash, &user.Status, &user.TrustScore,
		&user.RiskLevel, &user.AccountTier, &platformLinksJSON, &securityProfileJSON,
		&privacySettingsJSON, &dataRightsJSON, &deviceFingerprintsJSON,
		&biometricTokensJSON, &federatedIdentitiesJSON, &sessionManagementJSON,
		&anomalyDetectionJSON, &crossPlatformActivityJSON, &complianceDataJSON,
		&metadataJSON, &user.CreatedAt, &user.UpdatedAt, &user.LastSeenAt,
		&user.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrMetaUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get meta user: %w", err)
	}

	// Unmarshal JSON fields
	json.Unmarshal(platformLinksJSON, &user.PlatformLinks)
	json.Unmarshal(securityProfileJSON, &user.SecurityProfile)
	json.Unmarshal(privacySettingsJSON, &user.PrivacySettings)
	json.Unmarshal(dataRightsJSON, &user.DataRights)
	json.Unmarshal(deviceFingerprintsJSON, &user.DeviceFingerprints)
	json.Unmarshal(biometricTokensJSON, &user.BiometricTokens)
	json.Unmarshal(federatedIdentitiesJSON, &user.FederatedIdentities)
	json.Unmarshal(sessionManagementJSON, &user.SessionManagement)
	json.Unmarshal(anomalyDetectionJSON, &user.AnomalyDetection)
	json.Unmarshal(crossPlatformActivityJSON, &user.CrossPlatformActivity)
	json.Unmarshal(complianceDataJSON, &user.ComplianceData)
	json.Unmarshal(metadataJSON, &user.Metadata)

	// Cache the result
	r.cache.Set(ctx, cacheKey, &user, 15*time.Minute)

	return &user, nil
}

// GetByMetaID retrieves a meta user by Meta ID
func (r *MetaUserRepository) GetByMetaID(ctx context.Context, metaID string) (*model.MetaUser, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("meta_user:meta_id:%s", metaID)
	var userID uuid.UUID
	
	if err := r.cache.Get(ctx, cacheKey, &userID); err == nil {
		return r.GetByID(ctx, userID)
	}

	query := `SELECT id FROM meta_users WHERE meta_id = $1 AND deleted_at IS NULL`
	err := r.db.QueryRowContext(ctx, query, metaID).Scan(&userID)
	
	if err == sql.ErrNoRows {
		return nil, ErrMetaUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get meta user by meta_id: %w", err)
	}

	// Cache the meta_id -> id mapping
	r.cache.Set(ctx, cacheKey, &userID, 15*time.Minute)

	return r.GetByID(ctx, userID)
}

// GetByEmail retrieves a meta user by email
func (r *MetaUserRepository) GetByEmail(ctx context.Context, email string) (*model.MetaUser, error) {
	query := `SELECT id FROM meta_users WHERE email = $1 AND deleted_at IS NULL`
	var userID uuid.UUID
	err := r.db.QueryRowContext(ctx, query, email).Scan(&userID)
	
	if err == sql.ErrNoRows {
		return nil, ErrMetaUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get meta user by email: %w", err)
	}

	return r.GetByID(ctx, userID)
}

// Update updates a meta user with optimistic locking
func (r *MetaUserRepository) Update(ctx context.Context, user *model.MetaUser) error {
	user.UpdatedAt = time.Now()

	query := `
		UPDATE meta_users SET
			email = $2, email_verified = $3, phone_number = $4, phone_verified = $5,
			password_hash = $6, status = $7, trust_score = $8, risk_level = $9,
			account_tier = $10, platform_links = $11, security_profile = $12,
			privacy_settings = $13, data_rights = $14, device_fingerprints = $15,
			biometric_tokens = $16, federated_identities = $17, session_management = $18,
			anomaly_detection = $19, cross_platform_activity = $20, compliance_data = $21,
			metadata = $22, updated_at = $23, last_seen_at = $24
		WHERE id = $1 AND deleted_at IS NULL
	`

	platformLinksJSON, _ := json.Marshal(user.PlatformLinks)
	securityProfileJSON, _ := json.Marshal(user.SecurityProfile)
	privacySettingsJSON, _ := json.Marshal(user.PrivacySettings)
	dataRightsJSON, _ := json.Marshal(user.DataRights)
	deviceFingerprintsJSON, _ := json.Marshal(user.DeviceFingerprints)
	biometricTokensJSON, _ := json.Marshal(user.BiometricTokens)
	federatedIdentitiesJSON, _ := json.Marshal(user.FederatedIdentities)
	sessionManagementJSON, _ := json.Marshal(user.SessionManagement)
	anomalyDetectionJSON, _ := json.Marshal(user.AnomalyDetection)
	crossPlatformActivityJSON, _ := json.Marshal(user.CrossPlatformActivity)
	complianceDataJSON, _ := json.Marshal(user.ComplianceData)
	metadataJSON, _ := json.Marshal(user.Metadata)

	result, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.EmailVerified, user.PhoneNumber, user.PhoneVerified,
		user.PasswordHash, user.Status, user.TrustScore, user.RiskLevel, user.AccountTier,
		platformLinksJSON, securityProfileJSON, privacySettingsJSON, dataRightsJSON,
		deviceFingerprintsJSON, biometricTokensJSON, federatedIdentitiesJSON,
		sessionManagementJSON, anomalyDetectionJSON, crossPlatformActivityJSON,
		complianceDataJSON, metadataJSON, user.UpdatedAt, user.LastSeenAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update meta user: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrMetaUserNotFound
	}

	// Create event sourcing entry
	r.createEvent(ctx, user.ID, "user.updated", map[string]interface{}{
		"updated_at": user.UpdatedAt,
	})

	// Invalidate cache
	r.cache.Delete(ctx, fmt.Sprintf("meta_user:%s", user.ID))
	r.cache.Delete(ctx, fmt.Sprintf("meta_user:meta_id:%s", user.MetaID))

	return nil
}

// UpdateTrustScore updates the user's trust score
func (r *MetaUserRepository) UpdateTrustScore(ctx context.Context, userID uuid.UUID, trustScore float64, riskLevel model.RiskLevel) error {
	query := `
		UPDATE meta_users 
		SET trust_score = $2, risk_level = $3, updated_at = $4
		WHERE id = $1 AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query, userID, trustScore, riskLevel, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update trust score: %w", err)
	}

	// Invalidate cache
	r.cache.Delete(ctx, fmt.Sprintf("meta_user:%s", userID))

	return nil
}

// LinkPlatformAccount links a platform-specific account to the meta user
func (r *MetaUserRepository) LinkPlatformAccount(ctx context.Context, metaUserID uuid.UUID, platform string, platformUserID uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get current user
	user, err := r.GetByID(ctx, metaUserID)
	if err != nil {
		return err
	}

	// Update platform links
	if platform == "socialink" {
		user.PlatformLinks.SocialinkUserID = &platformUserID
	} else if platform == "vignette" {
		user.PlatformLinks.VignetteUserID = &platformUserID
	}
	user.PlatformLinks.LinkedAt = time.Now()
	user.PlatformLinks.LinkStatus = model.LinkStatusLinked

	platformLinksJSON, _ := json.Marshal(user.PlatformLinks)

	query := `
		UPDATE meta_users 
		SET platform_links = $2, updated_at = $3
		WHERE id = $1 AND deleted_at IS NULL
	`

	_, err = tx.ExecContext(ctx, query, metaUserID, platformLinksJSON, time.Now())
	if err != nil {
		return fmt.Errorf("failed to link platform account: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Create event
	r.createEvent(ctx, metaUserID, "platform.linked", map[string]interface{}{
		"platform":         platform,
		"platform_user_id": platformUserID,
	})

	// Invalidate cache
	r.cache.Delete(ctx, fmt.Sprintf("meta_user:%s", metaUserID))

	return nil
}

// SoftDelete performs a soft delete on the user
func (r *MetaUserRepository) SoftDelete(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	query := `
		UPDATE meta_users 
		SET deleted_at = $2, status = $3, updated_at = $4
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, userID, now, model.UserStatusDeleted, now)
	if err != nil {
		return fmt.Errorf("failed to soft delete user: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrMetaUserNotFound
	}

	// Create event
	r.createEvent(ctx, userID, "user.deleted", map[string]interface{}{
		"deleted_at": now,
	})

	// Invalidate cache
	r.cache.Delete(ctx, fmt.Sprintf("meta_user:%s", userID))

	return nil
}

// GetHighRiskUsers retrieves users with high risk levels
func (r *MetaUserRepository) GetHighRiskUsers(ctx context.Context, limit int) ([]*model.MetaUser, error) {
	query := `
		SELECT id 
		FROM meta_users 
		WHERE risk_level IN ('high', 'critical') AND deleted_at IS NULL
		ORDER BY trust_score ASC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get high risk users: %w", err)
	}
	defer rows.Close()

	var users []*model.MetaUser
	for rows.Next() {
		var userID uuid.UUID
		if err := rows.Scan(&userID); err != nil {
			continue
		}
		user, err := r.GetByID(ctx, userID)
		if err == nil {
			users = append(users, user)
		}
	}

	return users, nil
}

// createEvent creates an event sourcing entry
func (r *MetaUserRepository) createEvent(ctx context.Context, userID uuid.UUID, eventType string, data map[string]interface{}) error {
	dataJSON, _ := json.Marshal(data)
	query := `
		INSERT INTO meta_user_events (id, user_id, event_type, event_data, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query, uuid.New(), userID, eventType, dataJSON, time.Now())
	return err
}

// GetEventHistory retrieves the event history for a user (Event Sourcing)
func (r *MetaUserRepository) GetEventHistory(ctx context.Context, userID uuid.UUID, limit int) ([]map[string]interface{}, error) {
	query := `
		SELECT id, event_type, event_data, created_at
		FROM meta_user_events
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get event history: %w", err)
	}
	defer rows.Close()

	var events []map[string]interface{}
	for rows.Next() {
		var id uuid.UUID
		var eventType string
		var eventDataJSON []byte
		var createdAt time.Time

		if err := rows.Scan(&id, &eventType, &eventDataJSON, &createdAt); err != nil {
			continue
		}

		var eventData map[string]interface{}
		json.Unmarshal(eventDataJSON, &eventData)

		event := map[string]interface{}{
			"id":         id,
			"event_type": eventType,
			"data":       eventData,
			"created_at": createdAt,
		}
		events = append(events, event)
	}

	return events, nil
}
