//core/identity/resolver.go

package identity

type CredentialProvider interface {
    Identify() (*IdentityProfile, error)
    IsAvailable(caps map[core.Capability]bool) bool
}

func ResolveUser(ctx core.RuntimeContext) (*IdentityProfile, error) {
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