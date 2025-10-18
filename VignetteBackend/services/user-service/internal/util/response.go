package util

import (
	"encoding/json"
	"log"
	"net/http"
)

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

// RespondWithJSON sends a JSON response
func RespondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// RespondWithError sends an error response
func RespondWithError(w http.ResponseWriter, statusCode int, message string) {
	RespondWithJSON(w, statusCode, APIResponse{
		Success: false,
		Error:   message,
	})
}

// RespondWithSuccess sends a success response
func RespondWithSuccess(w http.ResponseWriter, message string, data interface{}) {
	RespondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// RespondWithValidationError sends a validation error response
func RespondWithValidationError(w http.ResponseWriter, field, message string) {
	details := map[string][]string{
		field: {message},
	}
	
	RespondWithJSON(w, http.StatusBadRequest, APIResponse{
		Success: false,
		Error:   "Validation error",
		Details: details,
	})
}

// RespondWithUnauthorized sends an unauthorized response
func RespondWithUnauthorized(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Unauthorized"
	}
	RespondWithError(w, http.StatusUnauthorized, message)
}

// RespondWithForbidden sends a forbidden response
func RespondWithForbidden(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Forbidden"
	}
	RespondWithError(w, http.StatusForbidden, message)
}

// RespondWithNotFound sends a not found response
func RespondWithNotFound(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Resource not found"
	}
	RespondWithError(w, http.StatusNotFound, message)
}

// RespondWithInternalError sends an internal server error response
func RespondWithInternalError(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Internal server error"
	}
	RespondWithError(w, http.StatusInternalServerError, message)
}

// RespondWithCreated sends a created response
func RespondWithCreated(w http.ResponseWriter, message string, data interface{}) {
	RespondWithJSON(w, http.StatusCreated, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}
