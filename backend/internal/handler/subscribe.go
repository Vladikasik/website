package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"aynshteyn.dev/backend/internal/model"
)

// Secret key for HMAC verification (in production, this should be in env vars)
var secretKey = []byte("aynshteyn-secret-key-change-in-production")

// SubscribeRequest represents the JSON payload for subscribing
type SubscribeRequest struct {
	Email         string `json:"email"`
	ClientToken   string `json:"clientToken"`
	ClientNonce   string `json:"clientNonce"`
	BrowserData   string `json:"browserData"`
	Timestamp     int64  `json:"timestamp"`
	ClientHMAC    string `json:"clientHMAC"`
}

// SubscribeHandler handles email subscription requests
func (app *Application) SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	// Limit request body size to prevent attacks
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024) // 1MB max

	// Parse JSON request
	var req SubscribeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		if errors.Is(err, io.EOF) {
			http.Error(w, "Request body empty", http.StatusBadRequest)
			return
		}
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.Email == "" || !strings.Contains(req.Email, "@") {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	// Validate timestamp is within a reasonable timeframe (5 minutes)
	now := time.Now().Unix()
	if now-req.Timestamp > 300 || req.Timestamp > now {
		http.Error(w, "Request expired or from the future", http.StatusBadRequest)
		return
	}

	// Validate client HMAC if provided (security feature)
	if req.ClientHMAC != "" && req.ClientToken != "" && req.ClientNonce != "" {
		if !validateClientHMAC(req) {
			app.Logger.Printf("HMAC validation failed for %s", req.Email)
			http.Error(w, "Invalid request signature", http.StatusBadRequest)
			return
		}
	}

	// Get client IP address
	ipAddr := r.Header.Get("X-Forwarded-For")
	if ipAddr == "" {
		ipAddr = r.RemoteAddr
	}

	// Create subscriber record
	subscriber := &model.Subscriber{
		Email:        req.Email,
		EmailHash:    model.HashEmail(req.Email),
		UserAgent:    r.UserAgent(),
		IPAddress:    ipAddr,
		BrowserInfo:  req.BrowserData,
		ChallengeKey: model.CreateChallengeKey(req.Email, r.UserAgent(), time.Now()),
	}

	// Save to database
	err = app.Store.SaveSubscriber(subscriber)
	if err != nil {
		if err.Error() == "email already exists" {
			// Don't reveal if email exists to prevent enumeration
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"message": "Thank you for subscribing",
			})
			return
		}
		
		app.Logger.Printf("Error saving subscriber: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Return success response with safe data
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Thank you for subscribing",
		"data":    subscriber.SafeResponse(),
	})
}

// validateClientHMAC validates the HMAC from the client
func validateClientHMAC(req SubscribeRequest) bool {
	// Create the message that was signed
	message := fmt.Sprintf("%s:%s:%d", req.Email, req.ClientNonce, req.Timestamp)
	
	// Create HMAC
	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(message))
	expectedMAC := hex.EncodeToString(mac.Sum(nil))
	
	// Compare HMACs (constant time compare to prevent timing attacks)
	return hmac.Equal([]byte(req.ClientHMAC), []byte(expectedMAC))
}

// GetAllSubscribersHandler returns all subscribers (admin only)
func (app *Application) GetAllSubscribersHandler(w http.ResponseWriter, r *http.Request) {
	subscribers, err := app.Store.GetAllSubscribers()
	if err != nil {
		app.Logger.Printf("Error getting subscribers: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Convert to safe responses
	safeResponses := make([]map[string]interface{}, 0, len(subscribers))
	for _, sub := range subscribers {
		resp := sub.SafeResponse()
		// Add email for admin view
		resp["email"] = sub.Email
		resp["user_agent"] = sub.UserAgent
		resp["ip_address"] = sub.IPAddress

		safeResponses = append(safeResponses, resp)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"subscribers": safeResponses,
		"count":       len(safeResponses),
	})
}

// AdminAuthMiddleware provides basic authentication for admin routes
func (app *Application) AdminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		// In production, use proper auth with environment variables
		// This is just a simple example
		if !ok || user != "admin" || pass != "aynshteyn-secure-password" {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
} 