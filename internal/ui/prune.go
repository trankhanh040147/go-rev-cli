package ui

import (
	"context"
	"fmt"
	"strings"

	"github.com/trankhanh040147/revcli/internal/gemini"
)

// PruneFile summarizes a file using Gemini Flash model
// Returns the summary string and any error
// flashClient must be a non-nil client instance (typically gemini-2.5-flash)
func PruneFile(ctx context.Context, flashClient *gemini.Client, filePath, content string) (string, error) {
	// Check if context is already cancelled before starting
	select {
	case <-ctx.Done():
		return "", fmt.Errorf("pruning cancelled: %w", ctx.Err())
	default:
	}

	if flashClient == nil {
		return "", fmt.Errorf("flash client is nil")
	}

	// Build summarization prompt
	prompt := fmt.Sprintf(PruneFilePromptTemplate, filePath, content)

	// Use a simple system prompt for summarization
	systemPrompt := "You are a code summarization assistant. Provide concise, one-sentence summaries of code files."

	// Generate summary (GenerateContent should respect context cancellation)
	summary, err := flashClient.GenerateContent(ctx, systemPrompt, prompt)
	if err != nil {
		// Check if error is due to context cancellation
		if ctx.Err() != nil {
			return "", fmt.Errorf("pruning cancelled: %w", ctx.Err())
		}
		return "", fmt.Errorf("failed to generate summary: %w", err)
	}

	// Clean up summary (remove quotes, trim whitespace)
	summary = strings.TrimSpace(summary)
	summary = strings.Trim(summary, `"`)

	return summary, nil
}
