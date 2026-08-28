package llm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractOpenAIText(t *testing.T) {
	tests := []struct {
		name     string
		response openAIResponse
		want     string
	}{
		{
			name: "output text field",
			response: openAIResponse{
				Output: []struct {
					Text    string `json:"text"`
					Content []struct {
						Type string `json:"type"`
						Text string `json:"text"`
					} `json:"content"`
				}{
					{Text: "Hello from output.text"},
				},
			},
			want: "Hello from output.text",
		},
		{
			name: "output content text field",
			response: openAIResponse{
				Output: []struct {
					Text    string `json:"text"`
					Content []struct {
						Type string `json:"type"`
						Text string `json:"text"`
					} `json:"content"`
				}{
					{
						Content: []struct {
							Type string `json:"type"`
							Text string `json:"text"`
						}{
							{Type: "output_text", Text: "Hello from output.content[].text"},
						},
					},
				},
			},
			want: "Hello from output.content[].text",
		},
		{
			name: "falls back to output_text when output has no usable text",
			response: openAIResponse{
				OutputText: "Hello from output_text fallback",
			},
			want: "Hello from output_text fallback",
		},
		{
			name:     "no usable text anywhere returns empty string",
			response: openAIResponse{},
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractOpenAIText(tt.response)
			assert.Equal(t, tt.want, got)
		})
	}
}
