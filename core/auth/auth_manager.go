// core\auth\auth_service.go
package auth

import (
	verification_persistence "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security/persistence"
	internal_common "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/common"
)

func NewAuthManager(
	vault verification_persistence.VaultStore,
	platform internal_common.PlatformClass,
) *AuthManager {
	return &AuthManager{
		Vault:    vault,
		Platform: platform,
	}
}
