package ai

import (
	"context"
	"fmt"
	"os"
	"time"

	"google.golang.org/genai"
)

// Client wraps the Gemini API client.
type Client struct {
	client *genai.Client
	model  string
}

// NewClient creates a new AI client.
func NewClient(apiKey, model string) (*Client, error) {
	if apiKey == "" {
		apiKey = os.Getenv("GOOGLE_API_KEY")
	}
	if apiKey == "" {
		apiKey = os.Getenv("GEMINI_API_KEY")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("no API key provided. Set GOOGLE_API_KEY env var or use --api-key flag")
	}
	if model == "" {
		model = "gemini-2.5-flash"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &Client{client: client, model: model}, nil
}

// GenerateContent sends a prompt to Gemini and returns the response text.
func (c *Client) GenerateContent(ctx context.Context, prompt string) (string, error) {
	result, err := c.client.Models.GenerateContent(
		ctx,
		c.model,
		[]*genai.Content{
			{Parts: []*genai.Part{{Text: prompt}}},
		},
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("gemini API error: %w", err)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from Gemini")
	}

	return result.Candidates[0].Content.Parts[0].Text, nil
}
