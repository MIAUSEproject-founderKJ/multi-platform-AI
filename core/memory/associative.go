//core/memory/associative.go

package memory

import (
	"math"
	"sort"
)

// MemoryMatch pairs a memory with its "Distance" from the current situation
type MemoryMatch struct {
	Entry    CognitiveEntry
	Distance float64 // Closer to 1.0 means more similar
}

// RecallSimilar finds the top-N memories that match the provided vector.
// This allows the AI to react to patterns even if the exact IDs don't match.
func (cv *CognitiveVault) RecallSimilar(targetVector []float32, topN int) []MemoryMatch {
	cv.mu.RLock()
	defer cv.mu.RUnlock()

	var matches []MemoryMatch

	for _, entry := range cv.Storage {
		if len(entry.Vector) == 0 || len(entry.Vector) != len(targetVector) {
			continue
		}

		similarity := cosineSimilarity(targetVector, entry.Vector)
		matches = append(matches, MemoryMatch{
			Entry:    entry,
			Distance: similarity,
		})
	}

	// Sort by similarity descending (highest similarity first)
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Distance > matches[j].Distance
	})

	if len(matches) > topN {
		return matches[:topN]
	}
	return matches
}

// cosineSimilarity calculates the dot product divided by the magnitudes
func cosineSimilarity(a, b []float32) float64 {
	var dotProduct, magA, magB float64
	for i := range a {
		dotProduct += float64(a[i] * b[i])
		magA += float64(a[i] * a[i])
		magB += float64(b[i] * b[i])
	}
	if magA == 0 || magB == 0 {
		return 0
	}
	return dotProduct / (math.Sqrt(magA) * math.Sqrt(magB))
}