package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"user-service/internal/handler"
	"user-service/internal/middleware"
)

// SetupRoutes configures all API routes
func SetupRoutes(authHandler *handler.AuthHandler, authMiddleware *middleware.AuthMiddleware) *mux.Router {
	r := mux.NewRouter()
	
	// API version prefix
	api := r.PathPrefix("/api/v1").Subrouter()
	
	// Health check
	api.HandleFunc("/health", healthCheck).Methods("GET")
	
	// Public auth routes
	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/signup", authHandler.HandleSignUp).Methods("POST")
	auth.HandleFunc("/login", authHandler.HandleLogin).Methods("POST")
	auth.HandleFunc("/forgot-password", authHandler.HandleForgotPassword).Methods("POST")
	auth.HandleFunc("/reset-password", authHandler.HandleResetPassword).Methods("POST")
	auth.HandleFunc("/verify-reset-token/{token}", authHandler.HandleVerifyResetToken).Methods("GET")
	
	// Cross-platform authentication
	crossPlatform := auth.PathPrefix("/cross-platform").Subrouter()
	crossPlatform.HandleFunc("/signin", authHandler.HandleCrossPlatformSignIn).Methods("POST")
	crossPlatform.HandleFunc("/check", authHandler.HandleCheckCrossPlatformAccount).Methods("GET")
	
	// Protected auth routes
	authProtected := api.PathPrefix("/auth").Subrouter()
	authProtected.Use(authMiddleware.RequireAuth)
	authProtected.HandleFunc("/me", authHandler.HandleGetCurrentUser).Methods("GET")
	authProtected.HandleFunc("/logout", authHandler.HandleLogout).Methods("POST")
	authProtected.HandleFunc("/refresh", authHandler.HandleRefreshToken).Methods("POST")
	
	// User management routes (protected)
	users := api.PathPrefix("/users").Subrouter()
	users.Use(authMiddleware.RequireAuth)
	users.HandleFunc("/{id}", authHandler.HandleGetUser).Methods("GET")
	users.HandleFunc("/{id}", authHandler.HandleUpdateUser).Methods("PUT")
	users.HandleFunc("/{id}", authHandler.HandleDeleteUser).Methods("DELETE")
	
	// CORS middleware
	r.Use(corsMiddleware)
	
	// Logging middleware
	r.Use(loggingMiddleware)
	
	return r
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy","service":"entativa-user-service","version":"1.0.0"}`))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
