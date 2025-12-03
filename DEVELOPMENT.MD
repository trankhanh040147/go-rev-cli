# Development Roadmap

---

# v0.1 - MVP Release ✅

**Status:** Completed

**Features Implemented:**
- [x] Cobra CLI framework with `root` and `review` commands
- [x] Git diff extraction (`git diff` and `git diff --staged`)
- [x] File-scope context: reads full content of modified files
- [x] Gemini API client with streaming response support
- [x] Interactive TUI with Bubbletea
  - [x] State machine (Loading → Reviewing → Chatting)
  - [x] Markdown rendering with Glamour
  - [x] Follow-up chat mode
  - [x] Keyboard shortcuts (q: quit, Enter: chat, Esc: back)
- [x] Senior Go Engineer persona prompt
- [x] File filtering (vendor/, generated, tests, go.sum)
- [x] Secret detection (API keys, tokens, passwords, private keys)
- [x] Command flags: `--staged`, `--model`, `--force`, `--no-interactive`
- [x] Non-interactive mode for CI/scripts

---

# v0.2 - Enhanced Diff & Context ✅

**Status:** Completed

**Features Implemented:**
- [x] **Project rename:** `go-rev-cli` → `rev-cli`
  - Module path: `github.com/trankhanh040147/rev-cli`
  - Binary name: `rev-cli`
- [x] **Custom base branch/commit comparison**
  - `--base <branch>` - Compare against a branch (e.g., `main`, `develop`)
  - `--base <commit>` - Compare against a specific commit hash
  - MR-style diff using `git diff base...HEAD`
- [x] **Update default model** - Changed to `gemini-2.5-pro`
- [x] **Show context preview** - Display files/tokens being sent before review
  - File list with sizes
  - Total file count and size
  - Ignored files list
  - Token estimate
- [x] **Token usage display** - Show actual tokens used after review
  - Prompt tokens
  - Completion tokens
  - Total tokens

**Breaking Changes:**
- Binary renamed: `go-rev-cli` → `rev-cli`
- Module path changed: `github.com/trankhanh040147/go-rev-cli` → `github.com/trankhanh040147/rev-cli`
- Default model changed: `gemini-1.5-flash` → `gemini-2.5-pro`

---

# v0.3 - Interactive Setup & Customization

**Status:** Planned

**Features:**
- [ ] **Dry-run mode** - `--dry-run` to preview payload without calling API
- [ ] **Model selection** - Allow user to choose from available models
  - `--list-models` flag to show available options
- [ ] **Interactive setup wizard** - After `review` command, prompt user to:
  - Select target branch (from list of branches)
  - Enter custom prompt/instructions (optional)
  - Choose model from available list
  - Confirm before sending
- [ ] **Custom prompt templates**
  - User-defined prompts in `~/.config/rev-cli/prompts/`
  - `--prompt <name>` to use custom template
  - Built-in templates: `strict`, `security`, `performance`, `style`
- [ ] **Language-agnostic review**
  - Auto-detect language from file extensions
  - Language-specific review prompts
  - Support: Go, Python, TypeScript, Rust, Java, etc.
- [ ] **Config file support**
  - `~/.config/rev-cli/config.yaml`
  - Default model, base branch, custom ignore patterns
- [ ] **Output formats**
  - `--format json` for CI integration
  - `--format markdown` for saving to file

---

# v0.4 - CI/CD & Team Features

**Status:** Future

**Features:**
- [ ] **GitHub Action** - `rev-cli` as a GitHub Action
- [ ] **GitLab CI template** - Ready-to-use CI config
- [ ] **PR/MR comment mode** - Post review as PR comment
- [ ] **Diff from PR URL** - `rev-cli review --pr <url>`
- [ ] **Review history** - Cache previous reviews locally
- [ ] **Cost tracking** - Track API usage and estimated costs

---

# v1.0 - Production Ready

**Status:** Future

**Features:**
- [ ] **Multiple AI providers**
  - OpenAI GPT-4
  - Anthropic Claude
  - Local models (Ollama)
- [ ] **Team configuration**
  - Shared config via `.rev-cli.yaml` in repo
  - Team-specific prompts and rules
- [ ] **Pre-commit hook** - Auto-review before commit
- [ ] **VS Code extension** - Review from editor
- [ ] **Review annotations** - Inline comments on specific lines
- [ ] **Severity filtering** - `--min-severity warning`

---

# Ideas Backlog

> Add feature ideas here for future consideration

- Auto-fix suggestions (apply LLM suggestions automatically)
- Compare two branches directly (`rev-cli diff main feature-branch`)
- Review specific files only (`rev-cli review src/api.go`)
- Ignore patterns via `.revignore` file
- Statistics dashboard (reviews done, issues found)
- Multi-language support (i18n for prompts)
- Plugin system for custom analyzers
- Integration with SonarQube/CodeClimate

---
