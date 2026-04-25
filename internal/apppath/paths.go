// internal/apppath/paths.go

package apppath

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	once sync.Once

	appName = "multi-platform-AI"

	// resolved paths (immutable after init)
	cachedRootDir    string
	cachedConfigDir  string
	cachedRuntimeDir string
	cachedLogDir     string
	cachedVaultDir   string

	// mode flags (deterministic)
	isPortableMode  bool
	isDevMode       bool
	isInstalledMode bool
)

// =========================
// Public API (stable access)
// =========================

func GetRootDir() string    { once.Do(initPaths); return cachedRootDir }
func GetConfigDir() string  { once.Do(initPaths); return cachedConfigDir }
func GetRuntimeDir() string { once.Do(initPaths); return cachedRuntimeDir }
func GetLogDir() string     { once.Do(initPaths); return cachedLogDir }
func GetVaultPath() string  { once.Do(initPaths); return cachedVaultDir }

func IsPortable() bool  { once.Do(initPaths); return isPortableMode }
func IsDev() bool       { once.Do(initPaths); return isDevMode }
func IsInstalled() bool { once.Do(initPaths); return isInstalledMode }

// =========================
// Initialization (deterministic)
// =========================

func initPaths() {
	// 1. Explicit override (HIGHEST PRIORITY, deterministic)
	if root := os.Getenv("AIOS_DATA_ROOT"); root != "" {
		resolveAsRoot(root, detectModeFromEnv())
		return
	}

	// 2. Portable mode via environment (NOT filesystem-based)
	if os.Getenv("AIOS_PORTABLE") == "true" {
		exeDir := mustExecutableDirFallback()
		resolveAsPortable(exeDir)
		return
	}

	// 3. Installed mode (system standard dirs)
	if runtime.GOOS != "" {
		if isSystemInstalled() {
			resolveAsInstalled()
			return
		}
	}

	// 4. Dev fallback (deterministic)
	exeDir := mustExecutableDirFallback()
	resolveAsDev(exeDir)
}

// =========================
// Mode resolution helpers
// =========================

func resolveAsRoot(root string, mode string) {
	cachedRootDir = filepath.Clean(root)

	switch mode {
	case "portable":
		isPortableMode = true
	case "installed":
		isInstalledMode = true
	default:
		isDevMode = true
	}

	cachedConfigDir = filepath.Join(cachedRootDir, "config")
	cachedRuntimeDir = filepath.Join(cachedRootDir, "runtime")
	cachedLogDir = filepath.Join(cachedRootDir, "logs")
	cachedVaultDir = filepath.Join(cachedRootDir, "vault")
}

func resolveAsPortable(root string) {
	isPortableMode = true
	cachedRootDir = filepath.Clean(root)

	cachedConfigDir = filepath.Join(root, "config")
	cachedRuntimeDir = filepath.Join(root, "runtime")
	cachedLogDir = filepath.Join(root, "logs")
	cachedVaultDir = filepath.Join(root, "vault")
}

func resolveAsInstalled() {
	isInstalledMode = true

	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		panic("failed to resolve user config dir")
	}

	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		panic("failed to resolve user cache dir")
	}

	root := filepath.Join(userConfigDir, appName)
	cachedRootDir = filepath.Clean(root)

	cachedConfigDir = root
	cachedRuntimeDir = filepath.Join(root, "runtime")
	cachedLogDir = filepath.Join(userCacheDir, appName, "logs")
	cachedVaultDir = filepath.Join(root, "vault")
}

func resolveAsDev(root string) {
	isDevMode = true
	cachedRootDir = filepath.Clean(root)

	cachedConfigDir = filepath.Join(root, "config")
	cachedRuntimeDir = filepath.Join(root, "runtime")
	cachedLogDir = filepath.Join(root, "runtime", "logs")
	cachedVaultDir = filepath.Join(root, "vault")
}

// =========================
// Environment detection
// =========================

func detectModeFromEnv() string {
	if os.Getenv("AIOS_PORTABLE") == "true" {
		return "portable"
	}
	if os.Getenv("AIOS_MODE") == "installed" {
		return "installed"
	}
	return "dev"
}

// =========================
// System detection (non-authoritative, advisory only)
// =========================

func isSystemInstalled() bool {
	// IMPORTANT:
	// This is NOT used as a primary switch (avoids nondeterminism)
	// Only used as a hint when no env override exists.

	exe, err := os.Executable()
	if err != nil {
		return false
	}

	exe = strings.ToLower(exe)

	switch runtime.GOOS {
	case "windows":
		programFiles := strings.ToLower(os.Getenv("ProgramFiles"))
		programFilesX86 := strings.ToLower(os.Getenv("ProgramFiles(x86)"))

		return strings.HasPrefix(exe, programFiles) ||
			strings.HasPrefix(exe, programFilesX86)

	case "darwin":
		return strings.Contains(exe, "/applications")

	case "linux":
		return strings.HasPrefix(exe, "/usr") ||
			strings.HasPrefix(exe, "/opt")

	default:
		return false
	}
}

// =========================
// Safe fallback
// =========================

func mustExecutableDirFallback() string {
	exe, err := os.Executable()
	if err != nil {
		// absolute fallback for deterministic behavior
		return "."
	}
	return filepath.Dir(exe)
}
