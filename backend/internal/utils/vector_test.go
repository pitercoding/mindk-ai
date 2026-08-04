package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseVector(t *testing.T) {

	input := `[0.1,0.2,0.3]`

	vector, err := ParseVector(input)

	require.NoError(t, err)

	assert.Equal(
		t,
		[]float32{
			0.1,
			0.2,
			0.3,
		},
		vector,
	)
}

func TestCosineSimilarity(t *testing.T) {

	tests := []struct {
		name     string
		a        []float32
		b        []float32
		expected float32
	}{
		{
			name: "identical vectors",
			a: []float32{
				1,
				2,
				3,
			},
			b: []float32{
				1,
				2,
				3,
			},
			expected: 1,
		},
		{
			name: "different vectors",
			a: []float32{
				1,
				0,
			},
			b: []float32{
				0,
				1,
			},
			expected: 0,
		},
		{
			name:     "empty vectors",
			a:        []float32{},
			b:        []float32{},
			expected: 0,
		},
	}

	for _, tt := range tests {

		t.Run(
			tt.name,
			func(t *testing.T) {

				result := CosineSimilarity(
					tt.a,
					tt.b,
				)

				assert.InDelta(
					t,
					tt.expected,
					result,
					0.001,
				)
			},
		)
	}
}
