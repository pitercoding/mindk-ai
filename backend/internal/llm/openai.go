package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Represents the client. Needs to know which key to use
type OpenAIClient struct {
	APIKey string
}

// Constructor
func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		APIKey: apiKey,
	}
}

// Request model
type openAIRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

// Response model
type openAIResponse struct {
	Output []struct {
		Text    string `json:"text"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
	OutputText string `json:"output_text"`
}

func extractOpenAIText(result openAIResponse) string {
	var builder strings.Builder

	for _, output := range result.Output {
		if output.Text != "" {
			builder.WriteString(output.Text)
		}

		for _, content := range output.Content {
			if content.Text != "" {
				builder.WriteString(content.Text)
			}
		}
	}

	if builder.Len() > 0 {
		return builder.String()
	}

	return result.OutputText
}

// Implements the "Client" interface
func (c *OpenAIClient) Chat(prompt string) (string, error) {

	// 1. Building the request
	requestBody := openAIRequest{
		Model: "gpt-5-mini",
		Input: prompt,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.openai.com/v1/responses",
		bytes.NewBuffer(jsonBody),
	)

	if err != nil {
		return "", err
	}

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	req.Header.Set(
		"Authorization",
		"Bearer "+c.APIKey,
	)

	// 2. Send
	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	// 3. Read response body once
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read openai response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"openai returned status %d: %s",
			resp.StatusCode,
			string(bytes.TrimSpace(body)),
		)
	}

	// 4. Parse response body
	var result openAIResponse
	err = json.Unmarshal(body, &result)

	if err != nil {
		return "", err
	}

	answer := extractOpenAIText(result)
	if answer == "" {
		return "", fmt.Errorf("empty response from openai")
	}

	return answer, nil
}
