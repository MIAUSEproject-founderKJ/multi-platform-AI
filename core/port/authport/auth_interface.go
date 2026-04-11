package authport

import "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"

type AuthInterface interface {
	StartAuthFlow(AuthManager) (*schema.UserSession, error)
}
