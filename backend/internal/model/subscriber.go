package model

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// Subscriber represents someone who submitted their email
type Subscriber struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	EmailHash    string    `json:"email_hash"`
	UserAgent    string    `json:"user_agent"`
	IPAddress    string    `json:"ip_address"`
	BrowserInfo  string    `json:"browser_info"`
	CreatedAt    time.Time `json:"created_at"`
	VerifyToken  string    `json:"-"`
	VerifiedAt   time.Time `json:"verified_at,omitempty"`
	ChallengeKey string    `json:"-"` // For future anti-bot measures
}

// HashEmail creates a hash of the email for storage
func HashEmail(email string) string {
	hash := sha256.Sum256([]byte(email))
	return hex.EncodeToString(hash[:])
}

// CreateChallengeKey creates a unique key for this submission
// This can be used for future anti-bot verification 
func CreateChallengeKey(email, userAgent string, timestamp time.Time) string {
	data := email + userAgent + timestamp.String()
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:8]) // Only use first 8 bytes (16 chars)
}

// SafeResponse returns a sanitized version of the subscriber for response
// This prevents leaking sensitive information
func (s *Subscriber) SafeResponse() map[string]interface{} {
	return map[string]interface{}{
		"email_hash":   s.EmailHash,
		"created_at":   s.CreatedAt,
		"is_verified":  !s.VerifiedAt.IsZero(),
		"challenge_id": s.ChallengeKey,
	}
} 