//core/auth/auth_manager.go

package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

type AuthManager struct {
	Vault    *security.IsolatedVault
	Identity *schema.MachineIdentity
	Platform schema.PlatformClass
	Entity   schema.EntityType
	Tier     schema.TierType
}

// LoginOrSignUp automatically runs the correct verification flow
// based on platform, entity type, and tier
func (am *AuthManager) LoginOrSignUp() (*schema.UserSession, error) {
	var err error
	fmt.Printf("[DEBUG] Platform: %v | Entity: %v | Tier: %v\n",
		am.Platform, am.Entity, am.Tier)
	switch am.Platform {

	// ------------------------------
	// Vehicles / Autonomous Mobility
	// ------------------------------
	case schema.PlatformVehicle, schema.PlatformDrone, schema.PlatformRobot:
		switch am.Entity {
		case schema.EntityPersonal:
			err = am.verifyKeyFobOrBiometrics()
		case schema.EntityOrganization:
			err = am.verifyBiometricsAndAppHandshake()
		case schema.EntityStranger:
			err = am.guestLoginVehicle()
		case schema.EntityTester:
			err = am.verifyMechanicAccess()
		default:
			err = fmt.Errorf("unknown vehicle entity type")
		}

	// ------------------------------
	// Industrial / Embedded / Factory
	// ------------------------------
	case schema.PlatformIndustrial, schema.PlatformEmbedded:
		err = am.verifyNFCCardOrButton()

	// ------------------------------
	// PCs / Laptops / Productivity
	// ------------------------------
	case schema.PlatformComputer, schema.PlatformLaptop, schema.PlatformMobile:
		switch am.Entity {
		case schema.EntityPersonal:
			err = am.verifyPasswordOrOSBiometrics()
		case schema.EntityOrganization:
			err = am.verifyPasswordOrOSBiometrics()
			if err == nil {
				err = am.verify2FAEnterprise()
			}
		case schema.EntityStranger:
			err = am.guestLoginPC()
		case schema.EntityTester:
			err = am.enableDebugLogin()
		default:
			err = fmt.Errorf("unknown PC entity type")
		}

	// ------------------------------
	// Fallback / Unknown
	// ------------------------------
	default:
		err = fmt.Errorf("unsupported platform: %s", am.Platform)
	}

	if err != nil {
		return nil, err
	}

	// Automatically create a session after successful login
	service := schema.ServiceUnknown
	switch am.Platform {
	case schema.PlatformVehicle, schema.PlatformDrone, schema.PlatformRobot:
		service = schema.ServiceEnterprise
	case schema.PlatformIndustrial, schema.PlatformEmbedded:
		service = schema.ServiceSystem
	case schema.PlatformComputer, schema.PlatformLaptop, schema.PlatformMobile:
		service = schema.ServicePersonal
	}

	return am.createSession(service)
}

// ------------------------------------------------------------
// Vehicle / Autonomous Mobility Auth
// ------------------------------------------------------------

func (am *AuthManager) verifyKeyFobOrBiometrics() error {
	// Simulate querying vehicle key-fob API or biometric reader
	time.Sleep(100 * time.Millisecond) // simulate latency
	verified := true                   // replace with actual hardware API
	if !verified {
		return errors.New("key-fob / biometrics verification failed")
	}
	fmt.Println("[AUTH] Vehicle: Key-fob or biometric verified")
	return nil
}

func (am *AuthManager) verifyBiometricsAndAppHandshake() error {
	// Enterprise vehicle: requires both biometrics + companion app handshake
	time.Sleep(150 * time.Millisecond)
	success := true // replace with real verification logic
	if !success {
		return errors.New("biometric + app handshake failed")
	}
	fmt.Println("[AUTH] Vehicle: Biometric + App handshake verified")
	return nil
}

func (am *AuthManager) guestLoginVehicle() error {
	// Minimal verification for passengers
	fmt.Println("[AUTH] Vehicle: Guest login activated (limited privileges)")
	return nil
}

func (am *AuthManager) verifyMechanicAccess() error {
	// Mechanic / Tester: full system access in sandbox
	time.Sleep(50 * time.Millisecond)
	fmt.Println("[AUTH] Vehicle: Mechanic/Tester full access granted")
	return nil
}

// ------------------------------------------------------------
// Industrial / Embedded / Factory Auth
// ------------------------------------------------------------

func (am *AuthManager) verifyNFCCardOrButton() error {
	// NFC card or physical pairing button
	time.Sleep(80 * time.Millisecond)
	valid := true
	if !valid {
		return errors.New("NFC/button verification failed")
	}
	fmt.Println("[AUTH] Industrial: NFC card or pairing button verified")
	return nil
}

// ------------------------------------------------------------
// Personal / Enterprise PC Auth
// ------------------------------------------------------------

func (am *AuthManager) verifyPasswordOrOSBiometrics() error {
	// Standard login: password + optional OS-biometrics
	time.Sleep(50 * time.Millisecond)
	fmt.Println("[AUTH] PC: Password/OS-biometrics verified")
	return nil
}

func (am *AuthManager) verify2FAEnterprise() error {
	// Enterprise 2FA: email/code, token, or app approval
	time.Sleep(100 * time.Millisecond)
	fmt.Println("[AUTH] PC: Enterprise 2FA verified")
	return nil
}

func (am *AuthManager) guestLoginPC() error {
	// Local-only guest session
	fmt.Println("[AUTH] PC: Guest session started (restricted access)")
	return nil
}

// ------------------------------------------------------------
// Debug / Tester Auth
// ------------------------------------------------------------

func (am *AuthManager) enableDebugLogin() error {
	fmt.Println("[AUTH] Debug/Test login granted (sandbox environment)")
	return nil
}

// ------------------------------------------------------------
// Session Creation
// ------------------------------------------------------------

func (am *AuthManager) createSession(service schema.ServiceType) (*schema.UserSession, error) {
	// Derive permissions dynamically from platform, entity, and tier
	permList := security.DerivePermissions(
		am.Platform,
		am.Entity,
		am.Tier,
	)

	// Use map[schema.Permission]bool to match UserSession type
	permMap := make(map[schema.Permission]bool)
	for _, p := range permList {
		permMap[p] = true
	}

	// Ensure the session has at least basic user permission
	permMap[schema.PermUser] = true

	// Build the UserSession struct dynamically
	session := &schema.UserSession{
		SessionID:   fmt.Sprintf("%d", time.Now().UnixNano()), // unique session ID
		Platform:    am.Platform,
		Entity:      am.Entity,
		Tier:        am.Tier,
		Service:     service,
		Permissions: permMap,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(24 * time.Hour), // default 24h expiration
	}

	fmt.Printf("[SESSION] Created: %s | Platform: %s | Entity: %x | Tier: %s | Service: %s\n",
		session.SessionID, session.Platform, session.Entity, session.Tier, session.Service)

	return session, nil
}
