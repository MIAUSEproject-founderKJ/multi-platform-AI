// bootstrap/boot_context_builder.go
package bootstrap

import verification_persistence "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security/persistence"

func NewBootContext(vault verification_persistence.VaultStore) BootContext {
	return BootContext{vault: vault}
}
