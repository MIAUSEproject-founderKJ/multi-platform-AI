//runtime/session/session_manager.go

package runtime_session

import (
	"errors"
	"fmt"
	"time"

	auth "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/auth"
	verification_persistence "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security/persistence"
	internal_common "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/common"
	"golang.org/x/crypto/bcrypt"
)

// UserProfile defines the human operator
type UserProfile struct {
	Username     string
	PasswordHash string
	Entity       internal_common.EntityKind
	Tier         internal_common.TierType
	CreatedAt    time.Time
}

// AuthSession represents an active login
type AuthSession struct {
	User      UserProfile
	Token     string
	ExpiresAt time.Time
}

func NewAuthManager(
	vault verification_persistence.VaultStore,
	platform internal_common.PlatformClass,
) *auth.AuthManager {

	return &auth.AuthManager{
		Vault:    vault,
		Platform: platform,
	}
}

type MyAuthManager struct {
	*auth.AuthManager
}

// Login verifies credentials and returns a session
func (am *MyAuthManager) Login(username, password string) (*AuthSession, error) {
	// 1. Fetch user from Vault
	var user UserProfile
	found, err := am.Vault.Read("users", username, &user)
	if err != nil || !found {
		return nil, errors.New("user not found")
	}

	// 2. Verify Password (Using SHA256 for demo; use bcrypt in prod)
	hash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}
	if user.PasswordHash != hash {
		return nil, errors.New("invalid credentials")
	}

	// 3. Create Session
	return &AuthSession{
		User:      user,
		Token:     generateToken(),
		ExpiresAt: time.Now().Add(12 * time.Hour),
	}, nil
}

// Signup creates a new user in the Vault
func (am *auth.AuthManager) Signup(
	username, password string,
	entity internal_common.EntityKind,
	tier internal_common.TierType,
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
