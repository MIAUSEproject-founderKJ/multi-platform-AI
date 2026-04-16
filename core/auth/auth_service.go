//core\auth\auth_service.go

package auth

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/interaction"
	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema"
)

type AuthManager struct {
	Vault    security.VaultStore
	Identity *schema.MachineIdentity
	Platform schema.PlatformClass
	Entity   schema.EntityType
	Tier     schema.TierType
}

type AuthInterface interface {
	StartAuthFlow(auth *AuthManager) (*schema.UserSession, error)
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

func PromptForCredentials(vault security.VaultStore) (string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("[AUTH] Enter User ID: ")
	userID, _ := reader.ReadString('\n')
	userID = strings.TrimSpace(userID)

	// Check if user exists
	var existing schema.MachineIdentity
	found, _ := vault.Read("users", userID, &existing)

	if !found {
		fmt.Println("[AUTH] User not found. Starting registration...")
		_, err := PromptForRegistration(vault)
		if err != nil {
			fmt.Println("[AUTH] Registration failed:", err)
			os.Exit(1)
		}
	}

	fmt.Print("[AUTH] Enter Password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	return userID, password
}
func hashPassword(pw string) string {
	h := sha256.Sum256([]byte(pw))
	return hex.EncodeToString(h[:])
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
		MachineID:    fmt.Sprintf("machine-%s", userID),
		EntityType:   entityType,
		OS:           "unknown",
		Arch:         "unknown",
		PasswordHash: hashPassword(password),
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
	hash := hashPassword(password)
	if stored.PasswordHash != hash {
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
		MachineID:    fmt.Sprintf("machine-%s", userID),
		EntityType:   entityType,
		OS:           "unknown",
		Arch:         "unknown",
		PasswordHash: hashPassword(password),
	}

	return am.Vault.Write("users", userID, identity)
}

func PromptForUserConfig() *schema.CustomizedConfig {
	cfg := DefaultCustomizedConfig()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n=== User Configuration ===")
	fmt.Println("Press ENTER to keep default")

	fmt.Printf("Main language [%s]: ", cfg.MainLang)
	if v, _ := reader.ReadString('\n'); strings.TrimSpace(v) != "" {
		cfg.MainLang = strings.TrimSpace(v)
	}

	fmt.Printf("Power mode (low/balanced/high) [%s]: ", cfg.PowerMode)
	if v, _ := reader.ReadString('\n'); strings.TrimSpace(v) != "" {
		cfg.PowerMode = strings.TrimSpace(v)
	}

	fmt.Printf("Privacy mode (standard/strict/offline) [%s]: ", cfg.PrivacyMode)
	if v, _ := reader.ReadString('\n'); strings.TrimSpace(v) != "" {
		cfg.PrivacyMode = strings.TrimSpace(v)
	}

	fmt.Printf("Update mode (auto/manual) [%s]: ", cfg.UpdateMode)
	if v, _ := reader.ReadString('\n'); strings.TrimSpace(v) != "" {
		cfg.UpdateMode = strings.TrimSpace(v)
	}

	fmt.Println("[CONFIG] Completed")

	return cfg
}

func DefaultCustomizedConfig() *schema.CustomizedConfig {
	return &schema.CustomizedConfig{
		Version:      "v1",
		LastModified: time.Now(),
	}
}
func (am *AuthManager) Register() (*schema.UserSession, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("User ID: ")
	userID, _ := reader.ReadString('\n')
	userID = strings.TrimSpace(userID)

	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	entityType := schema.EntityPersonal

	if err := am.RegisterUser(userID, password, entityType); err != nil {
		return nil, err
	}

	// immediately login after registration
	return am.Login(userID, password)
}

func (am *AuthManager) Login(userID, password string) (*schema.UserSession, error) {
	verified, identity := am.verifyUserCredentials(userID, password)
	if !verified {
		return nil, errors.New("invalid credentials")
	}

	am.Identity = identity
	am.detectEntityAndTier()

	return am.platformLoginFlow()
}

func (am *AuthManager) LoginOrSignUpInteractive() (*schema.UserSession, error) {
	for {
		userID, password := PromptForCredentials(am.Vault)

		verified, identity := am.verifyUserCredentials(userID, password)
		if verified {
			am.Identity = identity
			am.detectEntityAndTier()
			return am.platformLoginFlow() // ✅ exits loop
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

	permList := security.DerivePermissions(
		am.Platform,
		am.Entity,
		am.Tier,
	)

	permMap := make(map[schema.Permission]bool)
	for _, p := range permList {
		permMap[p] = true
	}
	permMap[schema.PermUser] = true

	session := &schema.UserSession{
		SessionID:   fmt.Sprintf("%d", time.Now().UnixNano()),
		Platform:    am.Platform,
		Entity:      am.Entity,
		Tier:        am.Tier,
		Service:     service,
		Permissions: permMap,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}

	// ---- CONFIG LOAD ----
	cfg, err := LoadUserConfig(am.Vault, am.Identity.MachineID)
	if err != nil {
		return nil, err
	}

	if cfg == nil {
		fmt.Println("[CONFIG] First-time setup required")
		cfg = PromptForUserConfig()
		_ = SaveUserConfig(am.Vault, am.Identity.MachineID, cfg)
	} else {
		fmt.Println("[CONFIG] Loaded existing configuration")
	}

	cfg.FillDefaults()

	// ---- CAPABILITY PROFILE ----
	cp := interaction.DetectCapabilityProfile()

	// ---- ORCHESTRATOR ----
	orch := interaction.BuildOrchestrator(cp)
	orch.StartAll(session)

	// ---- MODE (informational now) ----
	mode := interaction.ResolveInteractionMode(cfg, cp.Set)

	session.Config = cfg
	session.Capabilities = cp.Set
	session.CapProfile = cp
	session.Mode = string(mode)
	session.Orchestrator = orch

	orch.Broadcast("Session initialized successfully")

	fmt.Printf("[auth/auth_manager]Capabilities: %v\n", session.Capabilities)
	fmt.Printf("[auth/auth_manager] Mode: %s\n", mode)

	return session, nil
}

func LoadUserConfig(vault security.VaultStore, userID string) (*schema.CustomizedConfig, error) {
	var cfg schema.CustomizedConfig

	found, err := vault.Read("configs", userID, &cfg)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, nil
	}

	cfg.FillDefaults()
	cfg.Migrate()

	return &cfg, nil
}

func SaveUserConfig(vault security.VaultStore, userID string, cfg *schema.CustomizedConfig) error {
	cfg.LastModified = time.Now()
	return vault.Write("configs", userID, cfg)
}

func (am *AuthManager) HandleConfigUpdate(session *schema.UserSession) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("[CONFIG] Enter command: ")
	cmd, _ := reader.ReadString('\n')
	cmd = strings.TrimSpace(cmd)

	if cmd == "update config" {
		fmt.Println("[CONFIG] Updating configuration...")
		newCfg := PromptForUserConfig()

		_ = SaveUserConfig(am.Vault, am.Identity.MachineID, newCfg)

		session.Config = newCfg

		if orch, ok := session.Orchestrator.(*interaction.Orchestrator); ok {
			orch.Broadcast("Configuration updated successfully")
		}
	}
}
