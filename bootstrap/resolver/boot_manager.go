//bootstrap/resolver/boot_manager.go

package bootstrap_resolver

import (
	verification_persistence "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security/persistence"
	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/environment"
)

type BootManager struct {
	Vault    verification_persistence.VaultStore
	Identity *internal_environment.MachineIdentity
}
