package image

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/trankhanh040147/revcli/internal/gemini"
)

// Generator wraps image generation functionality
type Generator struct {
	client *gemini.Client
	apiKey string
}

// NewGenerator creates a new image generator
func NewGenerator(client *gemini.Client, apiKey string) *Generator {
	return &Generator{
		client: client,
		apiKey: apiKey,
	}
}

// GenerateImageResult contains the result of image generation including metadata
type GenerateImageResult struct {
	Description string
	Usage       *gemini.TokenUsage
	Cost        *gemini.CostBreakdown
}

// GenerateImage generates an image from a prompt and saves it to the specified path
func (g *Generator) GenerateImage(ctx context.Context, prompt, aspectRatio, resolution, outputPath string) (*GenerateImageResult, error) {
	// Validate aspect ratio
	if !isValidAspectRatio(aspectRatio) {
		return nil, fmt.Errorf("invalid aspect ratio: %s. Valid options: %v", aspectRatio, ValidAspectRatios)
	}

	// Validate resolution
	if !isValidResolution(resolution) {
		return nil, fmt.Errorf("invalid resolution: %s. Valid options: %v", resolution, ValidResolutions)
	}

	// Generate image using Gemini client
	result, err := g.client.GenerateImage(ctx, g.apiKey, prompt, aspectRatio, resolution)
	if err != nil {
		return nil, fmt.Errorf("failed to generate image: %w", err)
	}

	// Calculate cost based on token usage and resolution
	if result.Usage != nil {
		hasImageOutput := len(result.ImageData) > 0
		hasTextOutput := result.Text != ""
		result.Cost = CalculateCost(
			result.Usage.PromptTokens,
			result.Usage.CompletionTokens,
			resolution,
			hasImageOutput,
			hasTextOutput,
		)
	}

	// Determine output path
	if outputPath == "" {
		outputPath = DefaultOutputFile
	}

	// Ensure output directory exists
	outputDir := filepath.Dir(outputPath)
	if outputDir != "." && outputDir != "" {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	// Save image to file
	if err := os.WriteFile(outputPath, result.ImageData, 0644); err != nil {
		return nil, fmt.Errorf("failed to save image: %w", err)
	}

	// Prepare result
	genResult := &GenerateImageResult{
		Usage: result.Usage,
		Cost:  result.Cost,
	}

	// Set description
	if result.Text != "" {
		genResult.Description = result.Text
	} else {
		genResult.Description = fmt.Sprintf("Image generated successfully and saved to %s", outputPath)
	}

	return genResult, nil
}

// isValidAspectRatio checks if the aspect ratio is valid
func isValidAspectRatio(ratio string) bool {
	for _, valid := range ValidAspectRatios {
		if ratio == valid {
			return true
		}
	}
	return false
}

// isValidResolution checks if the resolution is valid
func isValidResolution(res string) bool {
	for _, valid := range ValidResolutions {
		if res == valid {
			return true
		}
	}
	return false
}
