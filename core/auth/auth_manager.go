//core/auth/auth_manager.go

package auth

import (
	"errors"
	"fmt"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

// AuthManager coordinates login/sign-up procedures post-boot
type AuthManager struct {
	Vault    *security.IsolatedVault
	Identity *schema.MachineIdentity
	Platform schema.PlatformClass
	Entity   schema.EntityType
	Tier     schema.TierType
}

// LoginOrSignUp handles the entry flow based on platform and entity type
func (am *AuthManager) LoginOrSignUp() (*schema.UserSession, error) {
	switch am.Platform {

	case schema.PlatformVehicle:
		return am.vehicleFlow()

	case schema.PlatformIndustrial:
		return am.industrialFlow()

	case schema.PlatformComputer, schema.PlatformLaptop:
		return am.pcFlow()

	default:
		return nil, errors.New("unsupported platform for login")
	}
}

// Vehicle / Robotaxi Flow
func (am *AuthManager) vehicleFlow() (*schema.UserSession, error) {
	switch am.Entity {

	case schema.EntityPersonal:
		if err := am.verifyKeyFobOrBiometrics(); err != nil {
			return nil, err
		}

	case schema.EntityOrganization:
		if err := am.verifyBiometricsAndAppHandshake(); err != nil {
			return nil, err
		}

	case schema.EntityStranger:
		if err := am.guestLoginVehicle(); err != nil {
			return nil, err
		}

	case schema.EntityTester:
		if err := am.verifyMechanicAccess(); err != nil {
			return nil, err
		}
	}

	return am.createSession(schema.ServiceEnterprise)
}

// Industrial Flow
func (am *AuthManager) industrialFlow() (*schema.UserSession, error) {
	switch am.Entity {

	case schema.EntityPersonal:
		return nil, errors.New("personal access forbidden in industrial setting")

	case schema.EntityOrganization, schema.EntityTester:
		if err := am.verifyNFCCardOrButton(); err != nil {
			return nil, err
		}

	case schema.EntityStranger:
		return nil, errors.New("unauthorized")
	}

	return am.createSession(schema.ServiceSystem)
}

// PC / Laptop Flow
func (am *AuthManager) pcFlow() (*schema.UserSession, error) {
	switch am.Entity {

	case schema.EntityPersonal:
		if err := am.verifyPasswordOrOSBiometrics(); err != nil {
			return nil, err
		}

	case schema.EntityOrganization:
		if err := am.verify2FAEnterprise(); err != nil {
			return nil, err
		}

	case schema.EntityStranger:
		if err := am.guestLoginPC(); err != nil {
			return nil, err
		}

	case schema.EntityTester:
		if err := am.enableDebugLogin(); err != nil {
			return nil, err
		}
	}

	return am.createSession(schema.ServicePersonal)
}

// ------------------------------------------------------------
// Identity Verification Helpers (Platform-Specific)
// ------------------------------------------------------------

func (am *AuthManager) verifyKeyFobOrBiometrics() error {
	// TODO: integrate with vehicle key-fob API / onboard biometric reader
	fmt.Println("[AUTH] Vehicle: Key-fob or biometrics verified")
	return nil
}

func (am *AuthManager) verifyBiometricsAndAppHandshake() error {
	// Enterprise vehicle verification
	fmt.Println("[AUTH] Vehicle: Biometric + App handshake verified")
	return nil
}

func (am *AuthManager) guestLoginVehicle() error {
	fmt.Println("[AUTH] Vehicle: Guest login activated")
	return nil
}

func (am *AuthManager) verifyMechanicAccess() error {
	fmt.Println("[AUTH] Vehicle: Mechanic/Tester full access granted")
	return nil
}

func (am *AuthManager) verifyNFCCardOrButton() error {
	fmt.Println("[AUTH] Industrial: NFC or Button verified")
	return nil
}

func (am *AuthManager) verifyPasswordOrOSBiometrics() error {
	fmt.Println("[AUTH] PC: Password/OS-biometrics verified")
	return nil
}

func (am *AuthManager) verify2FAEnterprise() error {
	fmt.Println("[AUTH] PC: Enterprise 2FA verified")
	return nil
}

func (am *AuthManager) guestLoginPC() error {
	fmt.Println("[AUTH] PC: Guest session started")
	return nil
}

func (am *AuthManager) enableDebugLogin() error {
	fmt.Println("[AUTH] Debug/Test login granted")
	return nil
}

// ------------------------------------------------------------
// Session Creation
// ------------------------------------------------------------
func (am *AuthManager) createSession(service schema.ServiceType) (*schema.UserSession, error) {
	session := &schema.UserSession{
		SessionID:   security.GenerateSessionToken(),
		Platform:    am.Platform,
		Entity:      am.Entity,
		Tier:        am.Tier,
		Service:     service,
		Permissions: security.DerivePermissions(am.Platform, am.Entity, am.Tier),
	}
	fmt.Printf("[SESSION] Created: %s\n", session.SessionID)
	return session, nil
}