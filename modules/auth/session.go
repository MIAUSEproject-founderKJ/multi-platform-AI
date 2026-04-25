//modules/auth/session.go

package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
"golang.org/x/crypto/bcrypt"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/auth"
	security_persistence "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security/persistence"
)

// UserProfile defines the human operator
type UserProfile struct {
    Username     string
    PasswordHash string
    Entity       schema_system.EntityType
    Tier         schema_identity.TierType
    CreatedAt    time.Time
}

// Session represents an active login
type Session struct {
	User      UserProfile
	Token     string
	ExpiresAt time.Time
}

func NewAuthManager(v security_persistence.VaultStore) *auth.AuthManager {
	return &auth.AuthManager{Vault: v}
}

type MyAuthManager struct {
	*auth.AuthManager
}

// Login verifies credentials and returns a session
func (am *MyAuthManager) Login(username, password string) (*Session, error) {
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
func (am *AuthManager) Signup(
    username, password string,
    entity schema_system.EntityType,
    tier schema_identity.TierType,
) error {

    exists, err := am.Vault.Exists("users", username)
    if err != nil {
        return fmt.Errorf("vault check failed: %w", err)
    }
    if exists {
        return errors.New("username already taken")
    }

    hash, err := hashPassword(password)
    if err != nil {
        return err
    }

    newUser := UserProfile{
        Username:     username,
        PasswordHash: hash,
        Entity:       entity,
        Tier:         tier,
        CreatedAt:    time.Now(),
    }

    if err := am.Vault.Write("users", username, newUser); err != nil {
        return fmt.Errorf("vault write failed: %w", err)
    }

    return nil
}

// Helpers
func hashPassword(pw string) (string, error) {
    b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
    return string(b), err
}

func generateToken() string {
	return fmt.Sprintf("sess_%d", time.Now().UnixNano())
}
