//boot/identity/resolver.go

package boot

import (
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/boot"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

type CredentialProvider interface {
	Identify() (*IdentityProfile, error)
	IsAvailable(caps map[schema.Capability]bool) bool
}

func ResolveUser(ctx boot.RuntimeContext) (*IdentityProfile, error) {
	// 1. If we have a Screen/Keyboard (Mobile/PC) -> Use UI Provider
	if ctx.Platform.Capabilities[core.CapKeyboard] {
		return UIProvider.Login()
	}

	// 2. If we have a Camera/NFC but NO Keyboard (Robot/Station) -> Use Passive Provider
	if ctx.Platform.Capabilities[core.CapCamera] || ctx.Platform.Capabilities[core.CapNFC] {
		return BiometricProvider.Scan()
	}

	// 3. Fallback: Trust-based or Remote-Auth (Cloud/Phone-as-Key)
	return RemoteProvider.WaitForKey()
}
