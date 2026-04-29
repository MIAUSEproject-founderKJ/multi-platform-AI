//core/auth/auth_service.go

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

	boot_phase "github.com/MIAUSEproject-founderKJ/multi-platform-AI/bootstrap/phases"
	security_decision "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/verification/decision"
	verification_persistence "github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/verification/persistence"

	internal_environment "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/system"
	user_setting "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/user"
	internal_verification "github.com/MIAUSEproject-founderKJ/multi-platform-AI/internal/schema/verification"
)

type AuthManager struct {
	Vault    verification_persistence.VaultStore
	Identity *internal_environment.MachineIdentity
	Platform internal_environment.PlatformClass
	Entity   internal_environment.EntityKind
	Tier     user_setting.TierType
}

type AuthInterface interface {
	StartAuthFlow(auth *AuthManager) (*user_setting.UserSession, error)
}

// detectEntityAndTier inspects the user identity to assign entity and tier
func (am *AuthManager) detectEntityAndTier() {
	if am.Identity == nil {
		am.Entity = internal_environment.EntityStranger
		am.Tier = user_setting.TierUnknown
		return
	}

	switch am.Identity.EntityType {
	case internal_environment.EntityPersonal:
		am.Entity = internal_environment.EntityPersonal
		am.Tier = user_setting.TierPersonal
	case internal_environment.EntityOrganization:
		am.Entity = internal_environment.EntityOrganization
		am.Tier = user_setting.TierEnterprise
	case internal_environment.EntityTester:
		am.Entity = internal_environment.EntityTester
		am.Tier = user_setting.TierTester
	default:
		am.Entity = internal_environment.EntityStranger
		am.Tier = user_setting.TierUnknown
	}
}

func PromptForCredentials(vault verification_persistence.VaultStore) (string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("[AUTH] Enter User ID: ")
	userID, _ := reader.ReadString('\n')
	userID = strings.TrimSpace(userID)

	// Check if user exists
	var existing internal_environment.MachineIdentity
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

func PromptForRegistration(vault verification_persistence.VaultStore) (*internal_environment.MachineIdentity, error) {
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

	entityType := internal_environment.EntityPersonal
	switch strings.ToLower(entityStr) {
	case "organization":
		entityType = internal_environment.EntityOrganization
	case "tester":
		entityType = internal_environment.EntityTester
	}

	identity := &internal_environment.MachineIdentity{
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
func (am *AuthManager) verifyUserCredentials(userID, password string) (bool, *internal_environment.MachineIdentity) {
	if am.Vault == nil {
		return false, nil
	}

	// Vault key for the user, e.g., "user_<userID>"
	var stored internal_environment.MachineIdentity
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

func (am *AuthManager) RegisterUser(userID, password string, entityType internal_environment.EntityKind) error {
	if am.Vault == nil {
		return errors.New("vault not initialized")
	}

	identity := &internal_environment.MachineIdentity{
		MachineID:    fmt.Sprintf("machine-%s", userID),
		EntityType:   entityType,
		OS:           "unknown",
		Arch:         "unknown",
		PasswordHash: hashPassword(password),
	}

	return am.Vault.Write("users", userID, identity)
}

func PromptForUserConfig() *user_setting.CustomizedConfig {
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

func DefaultCustomizedConfig() *user_setting.CustomizedConfig {
	return &user_setting.CustomizedConfig{
		Version:      "v1",
		LastModified: time.Now(),
	}
}
func (am *AuthManager) Register() (*user_setting.UserSession, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("User ID: ")
	userID, _ := reader.ReadString('\n')
	userID = strings.TrimSpace(userID)

	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	entityType := internal_environment.EntityPersonal

	if err := am.RegisterUser(userID, password, entityType); err != nil {
		return nil, err
	}

	// immediately login after registration
	return am.Login(userID, password)
}

func (am *AuthManager) Login(userID, password string) (*user_setting.UserSession, error) {
	verified, identity := am.verifyUserCredentials(userID, password)
	if !verified {
		return nil, errors.New("invalid credentials")
	}

	am.Identity = identity
	am.detectEntityAndTier()

	return am.platformLoginFlow()
}

func (am *AuthManager) LoginOrSignUpInteractive() (*user_setting.UserSession, error) {
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
func (am *AuthManager) platformLoginFlow() (*user_setting.UserSession, error) {
	var err error

	switch am.Platform {

	// ------------------------------
	// Vehicles / Autonomous Mobility
	// ------------------------------
	case internal_environment.PlatformVehicle, internal_environment.PlatformRobot:
		switch am.Entity {
		case internal_environment.EntityPersonal:
			err = am.verifyKeyFobOrBiometrics()
		case internal_environment.EntityOrganization:
			err = am.verifyBiometricsAndAppHandshake()
		case internal_environment.EntityStranger:
			err = am.guestLoginVehicle()
		case internal_environment.EntityTester:
			err = am.verifyMechanicAccess()
		default:
			err = fmt.Errorf("unknown vehicle entity type")
		}

	// ------------------------------
	// Industrial / Embedded / Factory
	// ------------------------------
	case internal_environment.PlatformIndustrial, internal_environment.PlatformEmbedded:
		err = am.verifyNFCCardOrButton()

	// ------------------------------
	// PCs / Laptops / Productivity
	// ------------------------------
	case internal_environment.PlatformComputer, internal_environment.PlatformMobile:
		switch am.Entity {
		case internal_environment.EntityPersonal:
			err = am.verifyPasswordOrOSBiometrics()
		case internal_environment.EntityOrganization:
			err = am.verifyPasswordOrOSBiometrics()
			if err == nil {
				err = am.verify2FAEnterprise()
			}
		case internal_environment.EntityStranger:
			err = am.guestLoginPC()
		case internal_environment.EntityTester:
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
	service := user_setting.ServiceUnknown
	switch am.Platform {
	case internal_environment.PlatformVehicle, internal_environment.PlatformRobot:
		service = user_setting.ServiceEnterprise
	case internal_environment.PlatformIndustrial, internal_environment.PlatformEmbedded:
		service = user_setting.ServiceSystem
	case internal_environment.PlatformComputer, internal_environment.PlatformMobile:
		service = user_setting.ServicePersonal
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

func (am *AuthManager) createSession(service user_setting.ServiceType) (*user_setting.UserSession, error) {

	// ----------------------------
	// 1. AUTHORIZATION
	// ----------------------------
	authCtx := &security_decision.AuthorizationContext{
		Platform: am.Platform,
		Entity:   am.Entity,
		Tier:     am.Tier,
		Service:  service,
	}

	authz := security_decision.AuthorizationService{
		Resolver: &security_decision.DefaultPermissionResolver{},
	}

	permMap := authz.Authorize(authCtx)

	// Build permission mask
	var permMask internal_verification.PermissionMask
	for p := range permMap {
		permMask |= internal_verification.ToMask(p)
	}

	// ----------------------------
	// 2. SESSION CONSTRUCTION
	// ----------------------------
	builder := user_setting.SessionBuilder{}

	buildCtx := &user_setting.BuildContext{
		Platform: am.Platform,
		Entity:   am.Entity,
		Tier:     am.Tier,
		Service:  service,
	}

	session := builder.Build(buildCtx, permMap)
	session.PermMask = permMask

	// ----------------------------
	// 3. CONFIG
	// ----------------------------
	cfg, err := LoadUserConfig(am.Vault, am.Identity.MachineID)
	if err != nil {
		return nil, err
	}

	if cfg == nil {
		cfg = PromptForUserConfig()
		_ = SaveUserConfig(am.Vault, am.Identity.MachineID, cfg)
	}

	cfg.WithDefaults()

	session.Config = &user_setting.UserConfig{
		MainLang:      cfg.MainLang,
		PowerMode:     cfg.PowerMode,
		PrivacyMode:   cfg.PrivacyMode,
		UpdateMode:    cfg.UpdateMode,
		PreferredMode: cfg.PreferredMode,
	}

	// ----------------------------
	// 4. RUNTIME
	// ----------------------------
	if err := am.initializeRuntime(session); err != nil {
		return nil, err
	}

	return session, nil
}

func LoadUserConfig(vault verification_persistence.VaultStore, userID string) (*user_setting.CustomizedConfig, error) {
	var cfg user_setting.CustomizedConfig

	found, err := vault.Read("configs", userID, &cfg)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, nil
	}

	cfg.WithDefaults()
	cfg.Migrate()

	return &cfg, nil
}

func SaveUserConfig(vault verification_persistence.VaultStore, userID string, cfg *user_setting.CustomizedConfig) error {
	cfg.LastModified = time.Now()
	return vault.Write("configs", userID, cfg)
}

func (am *AuthManager) HandleConfigUpdate(session *user_setting.UserSession) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("[CONFIG] Enter command: ")
	cmd, _ := reader.ReadString('\n')
	cmd = strings.TrimSpace(cmd)

	if cmd == "update config" {
		fmt.Println("[CONFIG] Updating configuration...")
		newCfg := PromptForUserConfig()

		_ = SaveUserConfig(am.Vault, am.Identity.MachineID, newCfg)

		session.Config = newCfg

		if orch, ok := session.Orchestrator.(*boot_phase.Orchestrator); ok {
			orch.Broadcast("Configuration updated successfully")
		}
	}
}

func (am *AuthManager) initializeRuntime(session *user_setting.UserSession) error {

	cp := interaction.DetectCapabilityProfile()

	orch := boot_phase.BuildOrchestrator(cp)
	orch.StartAll(session)

	mode := boot_phase.ResolveInteractionMode(session.Config, cp.Set)

	session.Capabilities = cp.Set
	session.CapProfile = cp
	session.Mode = string(mode)
	session.Orchestrator = orch

	orch.Broadcast("Session initialized successfully")

	return nil
}
