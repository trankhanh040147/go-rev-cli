package gemini

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// Client wraps the Gemini API client
type Client struct {
	client    *genai.Client
	model     *genai.GenerativeModel
	chat      *genai.ChatSession
	modelID   string
	lastUsage *TokenUsage
}

// TokenUsage contains token usage information from a response
type TokenUsage struct {
	PromptTokens     int32
	CompletionTokens int32
	TotalTokens      int32
}

// StreamCallback is called for each chunk of streamed response
type StreamCallback func(chunk string)

// NewClient creates a new Gemini client
func NewClient(ctx context.Context, apiKey, modelID string) (*Client, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	model := client.GenerativeModel(modelID)

	// Configure the model for code review
	model.SetTemperature(0.3) // Lower temperature for more focused responses
	model.SetTopP(0.95)
	model.SetTopK(40)

	// Set safety settings to be less restrictive for code content
	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockOnlyHigh,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockOnlyHigh,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockOnlyHigh,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockOnlyHigh,
		},
	}

	return &Client{
		client:  client,
		model:   model,
		modelID: modelID,
	}, nil
}

// StartChat initializes a chat session with the system prompt
func (c *Client) StartChat(systemPrompt string) {
	c.model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemPrompt)},
	}
	c.chat = c.model.StartChat()
}

// SendMessage sends a message and returns the full response
func (c *Client) SendMessage(ctx context.Context, message string) (string, error) {
	if c.chat == nil {
		return "", fmt.Errorf("chat session not initialized, call StartChat first")
	}

	resp, err := c.chat.SendMessage(ctx, genai.Text(message))
	if err != nil {
		return "", fmt.Errorf("failed to send message: %w", err)
	}

	return extractText(resp), nil
}

// SendMessageStream sends a message and streams the response
func (c *Client) SendMessageStream(ctx context.Context, message string, callback StreamCallback) (string, error) {
	if c.chat == nil {
		return "", fmt.Errorf("chat session not initialized, call StartChat first")
	}

	iter := c.chat.SendMessageStream(ctx, genai.Text(message))

	var fullResponse string
	var lastResp *genai.GenerateContentResponse
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fullResponse, fmt.Errorf("stream error: %w", err)
		}

		lastResp = resp
		chunk := extractText(resp)
		fullResponse += chunk
		if callback != nil {
			callback(chunk)
		}
	}

	// Extract token usage from the last response
	if lastResp != nil {
		c.lastUsage = extractUsage(lastResp)
	}

	return fullResponse, nil
}

// GenerateContent sends a one-off generation request (no chat history)
func (c *Client) GenerateContent(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	// Create a temporary model with system instruction
	model := c.client.GenerativeModel(c.modelID)
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemPrompt)},
	}

	resp, err := model.GenerateContent(ctx, genai.Text(userPrompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	return extractText(resp), nil
}

// GenerateContentStream sends a one-off generation request with streaming
func (c *Client) GenerateContentStream(ctx context.Context, systemPrompt, userPrompt string, callback StreamCallback) (string, error) {
	// Create a temporary model with system instruction
	model := c.client.GenerativeModel(c.modelID)
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemPrompt)},
	}
	model.SetTemperature(0.3)

	iter := model.GenerateContentStream(ctx, genai.Text(userPrompt))

	var fullResponse string
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fullResponse, fmt.Errorf("stream error: %w", err)
		}

		chunk := extractText(resp)
		fullResponse += chunk
		if callback != nil {
			callback(chunk)
		}
	}

	return fullResponse, nil
}

// Close closes the client connection
func (c *Client) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

// extractText extracts text content from a GenerateContentResponse
func extractText(resp *genai.GenerateContentResponse) string {
	if resp == nil || len(resp.Candidates) == 0 {
		return ""
	}

	var text string
	for _, candidate := range resp.Candidates {
		if candidate.Content != nil {
			for _, part := range candidate.Content.Parts {
				if t, ok := part.(genai.Text); ok {
					text += string(t)
				}
			}
		}
	}

	return text
}

// StreamWriter wraps an io.Writer for streaming responses
type StreamWriter struct {
	Writer io.Writer
}

// Write implements StreamCallback for writing to an io.Writer
func (sw *StreamWriter) Write(chunk string) {
	sw.Writer.Write([]byte(chunk))
}

// GetModelID returns the current model ID
func (c *Client) GetModelID() string {
	return c.modelID
}

// GetLastUsage returns the token usage from the last API call
func (c *Client) GetLastUsage() *TokenUsage {
	return c.lastUsage
}

// extractUsage extracts token usage from a response
func extractUsage(resp *genai.GenerateContentResponse) *TokenUsage {
	if resp == nil || resp.UsageMetadata == nil {
		return nil
	}

	return &TokenUsage{
		PromptTokens:     resp.UsageMetadata.PromptTokenCount,
		CompletionTokens: resp.UsageMetadata.CandidatesTokenCount,
		TotalTokens:      resp.UsageMetadata.TotalTokenCount,
	}
}

// FormatUsage returns a formatted string of token usage
func (u *TokenUsage) FormatUsage() string {
	if u == nil {
		return "Token usage not available"
	}
	return fmt.Sprintf("Tokens: %d prompt + %d completion = %d total",
		u.PromptTokens, u.CompletionTokens, u.TotalTokens)
}

// CostBreakdown contains detailed cost information
type CostBreakdown struct {
	InputCost  float64 // Cost for input tokens (USD)
	OutputCost float64 // Cost for output (text or image) (USD)
	TotalCost  float64 // Total cost (USD)
	InputCostVND  float64 // Cost for input tokens (VND)
	OutputCostVND float64 // Cost for output (text or image) (VND)
	TotalCostVND  float64 // Total cost (VND)
}

// ImageGenerationResult contains the result of an image generation request
type ImageGenerationResult struct {
	ImageData []byte
	Text      string
	Usage     *TokenUsage
	Cost      *CostBreakdown
}

// imageGenerationRequest represents the API request structure for image generation
type imageGenerationRequest struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
	GenerationConfig struct {
		ResponseModalities []string `json:"responseModalities"`
		ImageConfig        struct {
			AspectRatio string `json:"aspectRatio"`
			ImageSize   string `json:"imageSize"`
		} `json:"imageConfig"`
	} `json:"generationConfig"`
}

// imageGenerationResponse represents the API response structure for image generation
type imageGenerationResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text       *string `json:"text"`
				InlineData *struct {
					MimeType string `json:"mimeType"`
					Data     string `json:"data"`
				} `json:"inlineData"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
	UsageMetadata *struct {
		PromptTokenCount     int32 `json:"promptTokenCount"`
		CandidatesTokenCount int32 `json:"candidatesTokenCount"`
		TotalTokenCount      int32 `json:"totalTokenCount"`
	} `json:"usageMetadata"`
}

// GenerateImage generates an image using the gemini-3-pro-image-preview model via REST API
func (c *Client) GenerateImage(ctx context.Context, apiKey, prompt, aspectRatio, resolution string) (*ImageGenerationResult, error) {
	// Construct the API request
	reqBody := imageGenerationRequest{
		Contents: []struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			{
				Parts: []struct {
					Text string `json:"text"`
				}{
					{Text: prompt},
				},
			},
		},
	}
	reqBody.GenerationConfig.ResponseModalities = []string{"TEXT", "IMAGE"}
	reqBody.GenerationConfig.ImageConfig.AspectRatio = aspectRatio
	reqBody.GenerationConfig.ImageConfig.ImageSize = resolution

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make HTTP request to Gemini API
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", c.modelID, apiKey)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response
	var apiResp imageGenerationResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Extract token usage
	result := &ImageGenerationResult{}
	if apiResp.UsageMetadata != nil {
		result.Usage = &TokenUsage{
			PromptTokens:     apiResp.UsageMetadata.PromptTokenCount,
			CompletionTokens: apiResp.UsageMetadata.CandidatesTokenCount,
			TotalTokens:      apiResp.UsageMetadata.TotalTokenCount,
		}
		c.lastUsage = result.Usage
	}

	// Extract image data and text from response
	if len(apiResp.Candidates) == 0 {
		return nil, fmt.Errorf("no candidates in response")
	}

	imageFound := false
	for _, candidate := range apiResp.Candidates {
		for _, part := range candidate.Content.Parts {
			// Extract text
			if part.Text != nil {
				result.Text += *part.Text
			}
			// Extract image data - take the first valid image found
			// Note: Gemini API typically returns one image per request, but we handle
			// multiple by taking the first valid one to avoid overwriting
			if !imageFound && part.InlineData != nil {
				// Decode base64 image data
				decoded, err := base64.StdEncoding.DecodeString(part.InlineData.Data)
				if err == nil && len(decoded) > 0 {
					result.ImageData = decoded
					imageFound = true
				}
			}
		}
	}

	if !imageFound || len(result.ImageData) == 0 {
		return nil, fmt.Errorf("no image data found in response")
	}

	return result, nil
}
