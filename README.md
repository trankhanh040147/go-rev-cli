# go-rev-cli

[![Go Report Card](https://goreportcard.com/badge/github.com/trankhanh040147/go-rev-cli)](https://goreportcard.com/report/github.com/trankhanh040147/go-rev-cli)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> **Gemini-powered code reviewer CLI for Go developers.**

**go-rev-cli** is a local command-line tool that acts as an intelligent peer reviewer. It reads your local git changes and uses Google's Gemini LLM to analyze your code for bugs, optimization opportunities, and idiomatic Go practices—all before you push a single commit.

## Features

- **Smart Context:** Analyzes `git diff` plus full file contents to understand exactly what you changed and where it fits.
- **Privacy-First:** Runs locally with built-in secret detection to prevent accidentally sending credentials to the LLM.
- **Interactive Chat:** Ask follow-up questions about the review in an interactive TUI.
- **Performance:** Focused on analyzing specific changes rather than the entire codebase to save tokens and time.
- **Gemini Integration:** Leverages the large context window and reasoning of Gemini 1.5 Flash/Pro.

## Prerequisites

Before using the tool, ensure you have the following installed:

- **Go** (version 1.21 or higher)
- **Git** installed and initialized in your project.
- A **Google Gemini API Key** (Get one [here](https://aistudio.google.com/)).

## Installation

You can install the tool directly using `go install`:

```bash
go install github.com/trankhanh040147/go-rev-cli@latest
```

Or build from source:

```bash
git clone https://github.com/trankhanh040147/go-rev-cli.git
cd go-rev-cli
go build -o go-rev-cli .
```

## Configuration

Set your Gemini API key as an environment variable:

```bash
export GEMINI_API_KEY="your-api-key-here"
```

Or pass it directly via the `--api-key` flag.

## Usage

### Basic Review

Review all uncommitted changes in your repository:

```bash
go-rev-cli review
```

### Review Staged Changes Only

Review only the changes you've staged for commit:

```bash
go-rev-cli review --staged
```

### Use a Specific Model

Choose between `gemini-1.5-flash` (faster, cheaper) or `gemini-1.5-pro` (more thorough):

```bash
go-rev-cli review --model gemini-1.5-pro
```

### Non-Interactive Mode

Get the review output without the interactive chat interface:

```bash
go-rev-cli review --no-interactive
```

### Skip Secret Detection

If you're confident there are no secrets in your code (use with caution):

```bash
go-rev-cli review --force
```

## Interactive Mode

When running in interactive mode (default), you can:

- **View the review:** The AI analysis is displayed in a scrollable viewport
- **Ask follow-up questions:** Press `Enter` to enter chat mode
- **Navigate:** Use arrow keys or scroll to navigate the review
- **Exit:** Press `q` to quit, `Esc` to exit chat mode

## What Gets Reviewed

The tool analyzes:
- All modified `.go` files
- The git diff showing exact changes
- Full file context for better understanding

The tool automatically filters out:
- `go.sum` and `go.mod` files
- `vendor/` directory
- Generated files (`*_generated.go`, `*.pb.go`)
- Test files (`*_test.go`)
- Mock files

## Security

The tool includes basic secret detection that scans for:
- API keys and tokens
- Passwords and secrets
- Private keys
- Database URLs with credentials
- Common credential patterns

If potential secrets are detected, the review is aborted unless `--force` is used.

## Review Focus Areas

The AI reviewer acts as a Senior Go Engineer and focuses on:

1. **Bug Detection** - Logic errors, nil pointer dereferences, race conditions
2. **Idiomatic Go Patterns** - Error handling, interface design, naming conventions
3. **Performance Optimizations** - Unnecessary allocations, inefficient loops
4. **Security Concerns** - Input validation, SQL injection risks
5. **Code Quality** - Readability, documentation, test coverage suggestions

## Example Output

```
🔍 Go Code Review

📋 Review Context:
   • Files to review: 3
   • Files ignored: 2
   • Estimated tokens: ~2,500

### Summary
The changes implement a new user authentication handler...

### Issues Found
🔴 **Critical**: Missing error check on line 45
🟠 **Warning**: Race condition in concurrent access
🟡 **Suggestion**: Consider using sync.Pool for better performance

### Code Suggestions
...
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) for details.
