//runtime/session/session_manager.go

package runtime_session

import (
	"errors"
	"fmt"
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/auth"
	verification_persistence "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security/persistence"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
	"golang.org/x/crypto/bcrypt"
)

// UserProfile defines the human operator
type UserProfile struct {
	Username     string
	PasswordHash string
	Entity       internal_environment.EntityKind
	Tier         user_setting.TierType
	CreatedAt    time.Time
}

// Session represents an active login
type Session struct {
	User      UserProfile
	Token     string
	ExpiresAt time.Time
}

func NewAuthManager(v verification_persistence.VaultStore) *auth.AuthManager {
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
	entity internal_environment.EntityKind,
	tier user_setting.TierType,
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
