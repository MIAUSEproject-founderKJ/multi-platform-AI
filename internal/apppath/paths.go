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
	once           sync.Once
	isPortable     bool
	appName        = "AIofSpeech"
	// Path Caches
	cachedConfigDir string
	cachedRuntimeDir string
	cachedLogDir    string
)

// Public Getters
func GetDataDir() string    { once.Do(initPaths); return cachedRuntimeDir }
func GetConfigDir() string  { once.Do(initPaths); return cachedConfigDir }
func GetLogDir() string     { once.Do(initPaths); return cachedLogDir }
func IsPortable() bool      { once.Do(initPaths); return isPortable }

func initPaths() {
	// 1. Allow Environment Overrides (DevOps/Docker Inclusivity)
	if envDir := os.Getenv("AIOS_DATA_ROOT"); envDir != "" {
		setPaths(envDir, false)
		return
	}

	exePath, err := os.Executable()
	if err != nil {
		setPaths(".", false) // Fallback to current working dir
		return
	}
	appDir := filepath.Dir(exePath)

	// 2. Portable Mode Check (The ".portable.done" strategy)
	if _, err := os.Stat(filepath.Join(appDir, ".portable.done")); err == nil {
		isPortable = true
		setPaths(appDir, true)
		return
	}

	// 3. System-Specific Logic (The "Inclusive" expansion)
	if isInstalled(appDir) {
		userConfig, _ := os.UserConfigDir()      // Linux: ~/.config, Win: AppData, Mac: AppSupport
		userCache, _ := os.UserCacheDir()        // Linux: ~/.cache, Win: LocalAppData, Mac: Caches
		
		cachedConfigDir = filepath.Join(userConfig, appName)
		cachedRuntimeDir = filepath.Join(userConfig, appName, "runtime")
		cachedLogDir = filepath.Join(userCache, appName, "logs")
	} else {
		// Fallback for Dev environments
		setPaths(appDir, true)
	}
}

// setPaths helper to maintain structure consistency
func setPaths(root string, isLocal bool) {
	if isLocal {
		cachedConfigDir = filepath.Join(root, "configs")
		cachedRuntimeDir = filepath.Join(root, "runtime")
		cachedLogDir = filepath.Join(root, "runtime", "logs")
	} else {
		cachedConfigDir = root
		cachedRuntimeDir = filepath.Join(root, "runtime")
		cachedLogDir = filepath.Join(root, "logs")
	}
}

// isInstalled determines if the app is running from a protected system location
func isInstalled(path string) bool {
	path = strings.ToLower(path)
	switch runtime.GOOS {
	case "windows":
		pf := strings.ToLower(os.Getenv("ProgramFiles"))
		pf86 := strings.ToLower(os.Getenv("ProgramFiles(x86)"))
		return strings.HasPrefix(path, pf) || strings.HasPrefix(path, pf86)
	case "darwin": // macOS Inclusion
		return strings.HasPrefix(path, "/applications") || strings.HasPrefix(path, "/library")
	case "linux":
		return strings.HasPrefix(path, "/usr") || strings.HasPrefix(path, "/bin") || 
			   strings.HasPrefix(path, "/opt") || strings.HasPrefix(path, "/var")
	default:
		return false
	}
}