package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	appcontext "github.com/trankhanh040147/revcli/internal/context"
	"github.com/trankhanh040147/revcli/internal/filter"
	"github.com/trankhanh040147/revcli/internal/gemini"
	"github.com/trankhanh040147/revcli/internal/ui"
)

var (
	staged      bool
	model       string
	force       bool
	interactive bool
	baseBranch  string
)

// reviewCmd represents the review command
var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Review code changes using Gemini AI",
	Long: `Analyzes your git diff and provides an intelligent code review.
Uses Google's Gemini LLM to detect bugs, suggest optimizations,
and ensure idiomatic Go practices.

Examples:
  # Review staged changes with interactive chat
  revcli review --staged

  # Review changes against main branch
  revcli review --base main

  # Review all uncommitted changes with a specific model
  revcli review --model gemini-2.5-pro

  # Non-interactive mode (just print the review)
  revcli review --no-interactive

  # Skip secret detection check
  revcli review --force`,
	RunE: runReview,
}

func init() {
	rootCmd.AddCommand(reviewCmd)

	reviewCmd.Flags().BoolVar(&staged, "staged", false, "Review only staged changes (git diff --staged)")
	reviewCmd.Flags().StringVar(&baseBranch, "base", "", "Base branch/commit to compare against (e.g., main, develop, abc123)")
	reviewCmd.Flags().StringVar(&model, "model", "gemini-2.5-pro", "Gemini model to use (gemini-2.5-pro, gemini-1.5-flash, etc.)")
	reviewCmd.Flags().BoolVar(&force, "force", false, "Skip secret detection and proceed anyway")
	reviewCmd.Flags().BoolVar(&interactive, "interactive", true, "Enable interactive chat mode")
	reviewCmd.Flags().BoolVar(&interactive, "no-interactive", false, "Disable interactive chat mode")
}

func runReview(cmd *cobra.Command, args []string) error {
	// Check for API key
	apiKey := GetAPIKey()
	if apiKey == "" {
		return fmt.Errorf("GEMINI_API_KEY is required. Set it via environment variable or --api-key flag")
	}

	// Handle --no-interactive flag
	if cmd.Flags().Changed("no-interactive") {
		interactive = false
	}

	// Create context
	ctx := context.Background()

	// Validate mutually exclusive flags
	if staged && baseBranch != "" {
		return fmt.Errorf("cannot use --staged and --base together. Choose one")
	}

	// Step 1: Build the review context
	fmt.Println(ui.RenderTitle("🔍 Code Review"))
	fmt.Println()

	if baseBranch != "" {
		fmt.Printf("Comparing against: %s\n", baseBranch)
	} else if staged {
		fmt.Println("Reviewing staged changes...")
	} else {
		fmt.Println("Reviewing uncommitted changes...")
	}

	builder := appcontext.NewBuilder(staged, force, baseBranch)
	reviewCtx, err := builder.Build()
	if err != nil {
		// Check if it's a secrets error
		if reviewCtx != nil && len(reviewCtx.SecretsFound) > 0 {
			printSecretsWarning(reviewCtx.SecretsFound)
			return fmt.Errorf("review aborted due to potential secrets")
		}
		return fmt.Errorf("failed to build review context: %w", err)
	}

	// Check if there are changes to review
	if !reviewCtx.HasChanges() {
		fmt.Println(ui.RenderWarning("No changes detected. Make sure you have uncommitted changes."))
		return nil
	}

	// Print detailed summary with file list
	fmt.Println(ui.RenderSuccess("Changes collected!"))
	fmt.Println()
	fmt.Println(reviewCtx.DetailedSummary())
	fmt.Println()

	// Step 2: Initialize Gemini client
	fmt.Println("Connecting to Gemini API...")
	client, err := gemini.NewClient(ctx, apiKey, model)
	if err != nil {
		return fmt.Errorf("failed to create Gemini client: %w", err)
	}
	defer client.Close()

	fmt.Println(ui.RenderSuccess(fmt.Sprintf("Connected to %s", client.GetModelID())))
	fmt.Println()

	// Step 3: Run the review
	if interactive {
		// Interactive TUI mode
		return ui.Run(reviewCtx, client)
	}

	// Non-interactive mode
	return ui.RunSimple(ctx, reviewCtx, client)
}

// printSecretsWarning prints a warning about detected secrets
func printSecretsWarning(secrets []filter.SecretMatch) {
	fmt.Println(ui.RenderError("Potential secrets detected in your code!"))
	fmt.Println()

	for _, s := range secrets {
		fmt.Printf("  • %s (line %d): %s\n", s.FilePath, s.Line, s.Match)
	}

	fmt.Println()
	fmt.Println(ui.RenderWarning("Review aborted to prevent sending secrets to external API."))
	fmt.Println(ui.RenderHelp("Use --force to proceed anyway (not recommended)"))
	fmt.Println()

	os.Exit(1)
}
