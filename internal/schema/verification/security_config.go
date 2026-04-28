//internal/schema/verification/config.go

package internal_verification

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

func (c CustomizedConfig) WithDefaults() CustomizedConfig {
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

	return c
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

	b, err := json.Marshal(s)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
