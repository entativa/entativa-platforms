package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"messaging-service/internal/logger"
	"messaging-service/internal/repository"
	"messaging-service/internal/util"
)

type KeyHandler struct {
	prekeyRepo *repository.PrekeyRepository
	logger     *logger.Logger
}

func NewKeyHandler(prekeyRepo *repository.PrekeyRepository, logger *logger.Logger) *KeyHandler {
	return &KeyHandler{
		prekeyRepo: prekeyRepo,
		logger:     logger,
	}
}

// UploadPrekeys handles uploading prekeys for E2EE
func (h *KeyHandler) UploadPrekeys(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*repository.User)
	if !ok {
		util.RespondWithUnauthorized(w, "")
		return
	}

	var req struct {
		Prekeys []struct {
			KeyID     int    `json:"key_id"`
			PublicKey string `json:"public_key"`
			Signature string `json:"signature"`
		} `json:"prekeys"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(req.Prekeys) == 0 {
		util.RespondWithError(w, http.StatusBadRequest, "At least one prekey required")
		return
	}

	// Store prekeys
	if err := h.prekeyRepo.StorePrekeys(r.Context(), user.ID, req.Prekeys); err != nil {
		h.logger.Error("Failed to store prekeys", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to store prekeys")
		return
	}

	util.RespondWithCreated(w, "Prekeys uploaded", map[string]interface{}{
		"count": len(req.Prekeys),
	})
}

// GetPrekeyBundle retrieves a user's prekey bundle
func (h *KeyHandler) GetPrekeyBundle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetUserID := vars["userID"]

	// Get identity key
	identityKey, err := h.prekeyRepo.GetIdentityKey(r.Context(), targetUserID)
	if err != nil {
		util.RespondWithNotFound(w, "User keys not found")
		return
	}

	// Get a prekey
	prekey, err := h.prekeyRepo.GetUnusedPrekey(r.Context(), targetUserID)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "No prekeys available")
		return
	}

	bundle := map[string]interface{}{
		"identity_key":     identityKey,
		"signed_prekey":    prekey.PublicKey,
		"prekey_signature": prekey.Signature,
		"prekey_id":        prekey.KeyID,
	}

	util.RespondWithSuccess(w, "", bundle)
}

// UploadIdentityKey uploads a user's identity key
func (h *KeyHandler) UploadIdentityKey(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*repository.User)
	if !ok {
		util.RespondWithUnauthorized(w, "")
		return
	}

	var req struct {
		PublicKey string `json:"public_key"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.prekeyRepo.StoreIdentityKey(r.Context(), user.ID, req.PublicKey); err != nil {
		h.logger.Error("Failed to store identity key", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Failed to store identity key")
		return
	}

	util.RespondWithCreated(w, "Identity key uploaded", nil)
}

// GetIdentityKey retrieves a user's identity key
func (h *KeyHandler) GetIdentityKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetUserID := vars["userID"]

	identityKey, err := h.prekeyRepo.GetIdentityKey(r.Context(), targetUserID)
	if err != nil {
		util.RespondWithNotFound(w, "Identity key not found")
		return
	}

	util.RespondWithSuccess(w, "", map[string]interface{}{
		"public_key": identityKey,
	})
}
