//core\auth\auth_types.go

package auth

import (
	internal_common "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/common"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
)

type VaultStore interface {
	Read(namespace, key string, dst any) (bool, error)
	Write(namespace, key string, value any) error
}

type AuthInterface interface {
	StartAuthFlow(auth *AuthManager) (*user_setting.UserSession, error)
}

type Credentials struct {
	UserID   string
	Password string
}

type AuthManager struct {
	Vault    VaultStore
	Identity *internal_environment.MachineIdentity
	Platform internal_common.PlatformClass
	Entity   internal_common.EntityKind
	Tier     internal_common.TierType
}
