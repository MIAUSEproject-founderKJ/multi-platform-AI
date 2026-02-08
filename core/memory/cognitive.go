//core/memory/cognitive.go

/*
While the IsolatedVault (which we built earlier) stores "static" secrets like hardware hashes and encryption keys, the CognitiveVault stores "dynamic" intelligence—learned behaviors, environmental maps, and pattern recognition data.

Think of it as the difference between a person's Social Security Number (IsolatedVault) and their Skills/Memories (CognitiveVault).
*/

package memory

import (
	"sync"
	"time"
)

// CognitiveEntry represents a single "Memory Fragment"
type CognitiveEntry struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`       // "mapping", "behavior", "threat"
	Data      interface{} `json:"data"`
	Weight    float64   `json:"weight"`     // How "important" this memory is
	Timestamp time.Time `json:"timestamp"`
	Vector    []float32 `json:"vector,omitempty"` // For semantic search
}

type CognitiveVault struct {
	mu       sync.RWMutex
	Storage  map[string]CognitiveEntry
	BasePath string
}

// Store persists an experience into the AI's long-term memory
func (cv *CognitiveVault) Store(id string, entry CognitiveEntry) {
	cv.mu.Lock()
	defer cv.mu.Unlock()
	
	entry.Timestamp = time.Now()
	cv.Storage[id] = entry
	// Logic to flush to encrypted disk storage would go here
}

// Recall retrieves a specific memory by ID
func (cv *CognitiveVault) Recall(id string) (CognitiveEntry, bool) {
	cv.mu.RLock()
	defer cv.mu.RUnlock()
	
	entry, exists := cv.Storage[id]
	return entry, exists
}