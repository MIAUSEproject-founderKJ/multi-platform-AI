//MIAUSEproject-founderKJ/multi-platform-AI/cognition/memory/semantic_store.go

package memory

import (
	"encoding/json"

	"github.com/MIAUSEproject-founderKJ/multi-platform-AI/core/security"
)

type SemanticMemory struct {
	KnownLandmarks  map[string][]float64 `json:"landmarks"`
	UserPreferences map[string]string    `json:"prefs"`
}

type SemanticVault struct {
	Store *security.IsolatedVault
}

func (sv *SemanticVault) CommitExperience(mem SemanticMemory) error {
	data, _ := json.Marshal(mem)
	return sv.Store.WriteEncrypted("semantic_experience.bin", data)
}

func (sv *SemanticVault) RecallExperience() (*SemanticMemory, error) {
	data, err := sv.Store.ReadEncrypted("semantic_experience.bin")
	if err != nil {
		return nil, err
	}
	var mem SemanticMemory
	json.Unmarshal(data, &mem)
	return &mem, nil
}
