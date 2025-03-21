package store

import (
	"database/sql"
	"errors"
	"time"

	"aynshteyn.dev/backend/internal/model"
	_ "github.com/mattn/go-sqlite3"
)

// SQLiteStore is the database layer for storing subscriber data
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore creates a new SQLite database connection
func NewSQLiteStore(dsn string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	// Test database connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &SQLiteStore{db: db}, nil
}

// Close closes the database connection
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

// Migrate creates the necessary database schema if it doesn't exist
func (s *SQLiteStore) Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS subscribers (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        email TEXT NOT NULL UNIQUE,
        email_hash TEXT NOT NULL,
        user_agent TEXT NOT NULL,
        ip_address TEXT NOT NULL,
        browser_info TEXT NOT NULL,
        created_at TIMESTAMP NOT NULL,
        verify_token TEXT,
        verified_at TIMESTAMP,
        challenge_key TEXT
    );
    CREATE INDEX IF NOT EXISTS idx_email_hash ON subscribers(email_hash);
    `

	_, err := s.db.Exec(query)
	return err
}

// SaveSubscriber stores a new subscriber in the database
func (s *SQLiteStore) SaveSubscriber(subscriber *model.Subscriber) error {
	query := `
    INSERT INTO subscribers (
        email, email_hash, user_agent, ip_address, browser_info, 
        created_at, verify_token, challenge_key
    ) 
    VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `

	now := time.Now().UTC()
	subscriber.CreatedAt = now

	_, err := s.db.Exec(
		query,
		subscriber.Email,
		subscriber.EmailHash,
		subscriber.UserAgent,
		subscriber.IPAddress,
		subscriber.BrowserInfo,
		subscriber.CreatedAt,
		subscriber.VerifyToken,
		subscriber.ChallengeKey,
	)

	if err != nil {
		// Check for duplicate email
		if err.Error() == "UNIQUE constraint failed: subscribers.email" {
			return errors.New("email already exists")
		}
		return err
	}

	return nil
}

// GetSubscriberByEmail retrieves a subscriber by email
func (s *SQLiteStore) GetSubscriberByEmail(email string) (*model.Subscriber, error) {
	query := `
    SELECT id, email, email_hash, user_agent, ip_address, browser_info, 
           created_at, verify_token, verified_at, challenge_key
    FROM subscribers
    WHERE email = ?
    `

	var subscriber model.Subscriber
	var verifiedAt sql.NullTime

	err := s.db.QueryRow(query, email).Scan(
		&subscriber.ID,
		&subscriber.Email,
		&subscriber.EmailHash,
		&subscriber.UserAgent,
		&subscriber.IPAddress,
		&subscriber.BrowserInfo,
		&subscriber.CreatedAt,
		&subscriber.VerifyToken,
		&verifiedAt,
		&subscriber.ChallengeKey,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("subscriber not found")
		}
		return nil, err
	}

	if verifiedAt.Valid {
		subscriber.VerifiedAt = verifiedAt.Time
	}

	return &subscriber, nil
}

// GetAllSubscribers retrieves all subscribers
func (s *SQLiteStore) GetAllSubscribers() ([]*model.Subscriber, error) {
	query := `
    SELECT id, email, email_hash, user_agent, ip_address, browser_info, 
           created_at, verify_token, verified_at, challenge_key
    FROM subscribers
    ORDER BY created_at DESC
    `

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscribers []*model.Subscriber

	for rows.Next() {
		var subscriber model.Subscriber
		var verifiedAt sql.NullTime

		err := rows.Scan(
			&subscriber.ID,
			&subscriber.Email,
			&subscriber.EmailHash,
			&subscriber.UserAgent,
			&subscriber.IPAddress,
			&subscriber.BrowserInfo,
			&subscriber.CreatedAt,
			&subscriber.VerifyToken,
			&verifiedAt,
			&subscriber.ChallengeKey,
		)

		if err != nil {
			return nil, err
		}

		if verifiedAt.Valid {
			subscriber.VerifiedAt = verifiedAt.Time
		}

		subscribers = append(subscribers, &subscriber)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return subscribers, nil
} 