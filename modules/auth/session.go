//modules/auth/session.go

package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security"
)

// UserProfile defines the human operator
type UserProfile struct {
	Username     string    `json:"username"`
	PasswordHash string    `json:"password_hash"`
	Role         string    `json:"role"` // "OWNER", "OPERATOR", "GUEST"
	CreatedAt    time.Time `json:"created_at"`
}

// Session represents an active login
type Session struct {
	User      UserProfile
	Token     string
	ExpiresAt time.Time
}

type AuthManager struct {
	Vault *security.IsolatedVault
}

func NewAuthManager(v *security.IsolatedVault) *AuthManager {
	return &AuthManager{Vault: v}
}

// Login verifies credentials and returns a session
func (am *AuthManager) Login(username, password string) (*Session, error) {
	// 1. Fetch user from Vault
	var user UserProfile
	found, err := am.Vault.Read("users", username, &user)
	if err != nil || !found {
		return nil, errors.New("user not found")
	}

	// 2. Verify Password (Using SHA256 for demo; use bcrypt in prod)
	hash := hashPassword(password)
	if user.PasswordHash != hash {
		return nil, errors.New("invalid credentials")
	}

	// 3. Create Session
	return &Session{
		User:      user,
		Token:     generateToken(),
		ExpiresAt: time.Now().Add(12 * time.Hour),
	}, nil
}

// Signup creates a new user in the Vault
func (am *AuthManager) Signup(username, password string, role string) error {
	// Check if user exists
	exists, _ := am.Vault.Exists("users", username)
	if exists {
		return errors.New("username already taken")
	}

	newUser := UserProfile{
		Username:     username,
		PasswordHash: hashPassword(password),
		Role:         role,
		CreatedAt:    time.Now(),
	}

	// Save to Vault
	return am.Vault.Write("users", username, newUser)
}

// Helpers
func hashPassword(pw string) string {
	h := sha256.Sum256([]byte(pw))
	return hex.EncodeToString(h[:])
}

func generateToken() string {
	return fmt.Sprintf("sess_%d", time.Now().UnixNano())
}