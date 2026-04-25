//internal/schema/security/config.go

package schema_security

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"
)

type CustomizedConfig struct {
	Version      string
	LastModified time.Time

	MainLang    string
	PowerMode   string
	PrivacyMode string
	UpdateMode  string

	PreferredMode string // IMPORTANT for runtime override
}

func (c *CustomizedConfig) FillDefaults() {
	if c.MainLang == "" {
		c.MainLang = "en"
	}
	if c.PowerMode == "" {
		c.PowerMode = "balanced"
	}
	if c.PrivacyMode == "" {
		c.PrivacyMode = "standard"
	}
	if c.UpdateMode == "" {
		c.UpdateMode = "auto"
	}
	if c.PreferredMode == "" {
		c.PreferredMode = "auto"
	}
}

func (c *CustomizedConfig) Hash() string {
	type stable struct {
		Version       string
		MainLang      string
		PowerMode     string
		PrivacyMode   string
		UpdateMode    string
		PreferredMode string
	}

	s := stable{
		Version:       c.Version,
		MainLang:      c.MainLang,
		PowerMode:     c.PowerMode,
		PrivacyMode:   c.PrivacyMode,
		UpdateMode:    c.UpdateMode,
		PreferredMode: c.PreferredMode,
	}

	b, _ := json.Marshal(s)
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
