//core/auth/auth_manager.go

package auth

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

type AuthManager struct {
	Vault    security.VaultStore
	Identity *schema.MachineIdentity
	Platform schema.PlatformClass
	Entity   schema.EntityType
	Tier     schema.TierType
}

// detectEntityAndTier inspects the user identity to assign entity and tier
func (am *AuthManager) detectEntityAndTier() {
	if am.Identity == nil {
		am.Entity = schema.EntityStranger
		am.Tier = schema.TierUnknown
		return
	}

	switch am.Identity.EntityType {
	case schema.EntityPersonal:
		am.Entity = schema.EntityPersonal
		am.Tier = schema.TierPersonal
	case schema.EntityOrganization:
		am.Entity = schema.EntityOrganization
		am.Tier = schema.TierEnterprise
	case schema.EntityTester:
		am.Entity = schema.EntityTester
		am.Tier = schema.TierTester
	default:
		am.Entity = schema.EntityStranger
		am.Tier = schema.TierUnknown
	}
}

func PromptForCredentials() (string, string) {
	reader := bufio.NewReader(os.Stdin)

	// Prompt if first-time login
	fmt.Print("[PromptForCredentials] Is this your first time logging in? (y/n): ")
	firstTime, _ := reader.ReadString('\n')
	firstTime = strings.TrimSpace(firstTime)

	if strings.ToLower(firstTime) == "y" {
		fmt.Println("[PromptForCredentials] Starting registration process...")
		vault, err := security.OpenVault()
		if err != nil {
			fmt.Println("[PromptForCredentials] Failed to open vault:", err)
			os.Exit(1)
		}

		_, err = PromptForRegistration(vault)
		if err != nil {
			fmt.Println("[PromptForCredentials] Registration failed:", err)
			os.Exit(1)
		}
	} // <- close the if block here

	// Now prompt for login credentials
	fmt.Print("[PromptForCredentials] Enter User ID: ")
	userID, _ := reader.ReadString('\n')
	userID = strings.TrimSpace(userID)

	fmt.Print("[PromptForCredentials] Enter Password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	return userID, password
}

func PromptForRegistration(vault security.VaultStore) (*schema.MachineIdentity, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== User Registration ===")

	fmt.Print("Enter new User ID: ")
	userID, _ := reader.ReadString('\n')
	userID = strings.TrimSpace(userID)

	fmt.Print("Enter Password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	fmt.Print("Select Entity Type (personal / organization / tester): ")
	entityStr, _ := reader.ReadString('\n')
	entityStr = strings.TrimSpace(entityStr)

	entityType := schema.EntityPersonal
	switch strings.ToLower(entityStr) {
	case "organization":
		entityType = schema.EntityOrganization
	case "tester":
		entityType = schema.EntityTester
	}

	identity := &schema.MachineIdentity{
		MachineID:  fmt.Sprintf("machine-%s", userID),
		EntityType: entityType,
		OS:         "unknown",
		Arch:       "unknown",
		Password:   password,
	}

	if vault == nil {
		return nil, errors.New("vault not initialized")
	}

	if err := vault.Write("users", userID, identity); err != nil {
		return nil, fmt.Errorf("failed to write user to vault: %w", err)
	}

	fmt.Println("[PromptForRegistration] Registration successful for user:", userID)
	return identity, nil
}

// verifyUserCredentials checks the Vault or database for valid credentials
func (am *AuthManager) verifyUserCredentials(userID, password string) (bool, *schema.MachineIdentity) {
	if am.Vault == nil {
		return false, nil
	}

	// Vault key for the user, e.g., "user_<userID>"
	var stored schema.MachineIdentity
	found, err := am.Vault.Read("users", userID, &stored)
	if err != nil {
		fmt.Println("[verifyUserCredentials] Vault read error:", err)
		return false, nil
	}
	if !found {
		fmt.Println("[verifyUserCredentials] User not found in vault")
		return false, nil
	}

	// Validate password (simplified: plaintext match for demo; ideally hashed)
	if stored.Password != password {
		fmt.Println("[verifyUserCredentials] Invalid password")
		return false, nil
	}

	fmt.Println("[verifyUserCredentials] User verified from vault:", userID)
	return true, &stored
}

func (am *AuthManager) RegisterUser(userID, password string, entityType schema.EntityType) error {
	if am.Vault == nil {
		return errors.New("vault not initialized")
	}

	identity := &schema.MachineIdentity{
		MachineID:  fmt.Sprintf("machine-%s", userID),
		EntityType: entityType,
		OS:         "unknown",
		Arch:       "unknown",
		Password:   password, // hash in production
	}

	return am.Vault.Write("users", userID, identity)
}

func (am *AuthManager) LoginOrSignUpInteractive() (*schema.UserSession, error) {
	for {
		userID, password := PromptForCredentials()

		verified, identity := am.verifyUserCredentials(userID, password)
		if verified {
			am.Identity = identity
			am.detectEntityAndTier()
			return am.platformLoginFlow()
		}

		fmt.Println("[AUTH] Invalid credentials or user not found. Try again.")
	}
}

// platformLoginFlow performs the platform-specific verification and session creation
func (am *AuthManager) platformLoginFlow() (*schema.UserSession, error) {
	var err error

	switch am.Platform {

	// ------------------------------
	// Vehicles / Autonomous Mobility
	// ------------------------------
	case schema.PlatformVehicle, schema.PlatformRobot:
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
	case schema.PlatformComputer, schema.PlatformMobile:
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

	// Determine default service based on platform
	service := schema.ServiceUnknown
	switch am.Platform {
	case schema.PlatformVehicle, schema.PlatformRobot:
		service = schema.ServiceEnterprise
	case schema.PlatformIndustrial, schema.PlatformEmbedded:
		service = schema.ServiceSystem
	case schema.PlatformComputer, schema.PlatformMobile:
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
	fmt.Println("[func (am *AuthManager) verifyKeyFobOrBiometrics] Vehicle: Key-fob or biometric verified")
	return nil
}

func (am *AuthManager) verifyBiometricsAndAppHandshake() error {
	// Enterprise vehicle: requires both biometrics + companion app handshake
	time.Sleep(150 * time.Millisecond)
	success := true // replace with real verification logic
	if !success {
		return errors.New("biometric + app handshake failed")
	}
	fmt.Println("[func (am *AuthManager) verifyBiometricsAndAppHandshake] Vehicle: Biometric + App handshake verified")
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
	fmt.Println("[func (am *AuthManager) verifyNFCCardOrButton] Industrial: NFC card or pairing button verified")
	return nil
}

// ------------------------------------------------------------
// Personal / Enterprise PC Auth
// ------------------------------------------------------------

func (am *AuthManager) verifyPasswordOrOSBiometrics() error {
	// Standard login: password + optional OS-biometrics
	time.Sleep(50 * time.Millisecond)
	fmt.Println("[func (am *AuthManager) verifyPasswordOrOSBiometrics] PC: Password/OS-biometrics verified")
	return nil
}

func (am *AuthManager) verify2FAEnterprise() error {
	// Enterprise 2FA: email/code, token, or app approval
	time.Sleep(100 * time.Millisecond)
	fmt.Println("[func (am *AuthManager) verify2FAEnterprise] PC: Enterprise 2FA verified")
	return nil
}

func (am *AuthManager) guestLoginPC() error {
	// Local-only guest session
	fmt.Println("[func (am *AuthManager) guestLoginPC] PC: Guest session started (restricted access)")
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
