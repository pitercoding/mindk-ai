package utils

import (
	"encoding/json"
	"math"
)

// ParseVector converts a JSON vector string into float32 slice.
func ParseVector(data string) ([]float32, error) {

	var vector []float32

	err := json.Unmarshal([]byte(data), &vector)
	if err != nil {
		return nil, err
	}

	return vector, nil
}

// CosineSimilarity calculates the similarity between two vectors.
func CosineSimilarity(a []float32, b []float32) float32 {

	if len(a) != len(b) || len(a) == 0 {
		return 0
	}

	var dotProduct float32
	var normA float32
	var normB float32

	for i := range a {

		dotProduct += a[i] * b[i]

		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
}
