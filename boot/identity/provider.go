//boot/identity/provider.go
package identity

import "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core"

type ProviderType string

const (
	ProviderUI      ProviderType = "UI_INTERACTIVE" // Keyboard/Phone screen
	ProviderPassive ProviderType = "PASSIVE_SENSE"  // NFC/Biometric
	ProviderRemote  ProviderType = "REMOTE_RELAY"   // Approve via another device
)

type IdentityProvider interface {
	Type() ProviderType
	// Authenticate returns a valid IdentityProfile or an error
	Authenticate(ctx runtime.RuntimeContext) (*core.IdentityProfile, error)
}
