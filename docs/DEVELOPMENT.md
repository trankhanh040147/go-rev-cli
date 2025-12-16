# Development Roadmap

## Design Principles & Coding Standards

> **Reference:** All design principles, coding standards, and implementation guidelines are defined in [`.cursor/rules/rules.mdc`](../.cursor/rules/rules.mdc).

### How To Apply These Rules

Automatically loads rules from the `.cursor/rules/` directory. The `rules.mdc` file includes `alwaysApply: true` in its frontmatter, which ensures:

- **Automatic Application:** Rules are always active during coding sessions
- **Context Awareness:** Understands project-specific patterns (Vim navigation, TUI-first UX, Go conventions)
- **Consistency:** All code suggestions follow the defined principles without manual reminders

## Bug Fix Protocol

1. **Global Fix:** Search codebase (`rg`/`fd`) for similar patterns/implementations. Fix **all** occurrences, not just the reported one.
2. **Documentation:**
    - Update "Known Bugs" table (Status: Fixed).
    - Update coding standards in `.cursor/rules/rules.mdc` if the bug reflects a common anti-pattern.
3. **Testing:** Verify edge cases: Interactive, Piped (`|`), Redirected (`<`), and Non-interactive modes.
> **Reference:** Bug Fix Protocol are defined in [`.cursor/rules/rules.mdc`](../.cursor/rules/rules.mdc).

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
- Default model changed: `gemini-1.5-flash` → `gemini-2.5-pro`

---

# v0.3.0 - Short Flags & Preset Management ✅

**Status:** Completed

**Features Implemented:**

### Short Flag Aliases ✅
- [x] Short aliases for all flags (`-s`, `-b`, `-m`, `-f`, `-i`, `-I`, `-k`, `-p`)
- [x] Version flag (`--version`, `-v`)

### Vim-Style Keybindings ✅
- [x] Navigation: `j/k`, `g/G`, `Ctrl+d/u/f/b`
- [x] Search: `/`, `n/N`, `Tab` toggle
- [x] Help overlay: `?` key

### Yank to Clipboard ✅
- [x] `y` - Yank entire review + chat history
- [x] `Y` - Yank only last response
- [x] `yb` - Yank code block
- [x] Visual feedback (toast notification)

### Review Presets ✅
- [x] `--preset <name>` / `-p` flag
- [x] Built-in presets: `quick`, `strict`, `security`, `performance`, `logic`, `style`, `typo`, `naming`
- [x] Custom presets in `~/.config/revcli/presets/*.yaml`
- [x] Default preset support via config
- [x] Preset replace mode (`--preset-replace` / `-R`)

### Preset Management Commands ✅
- [x] `preset list` - List all presets
- [x] `preset create` - Create custom preset
- [x] `preset edit` - Edit custom preset (external editor)
- [x] `preset delete` - Delete custom preset
- [x] `preset show` - Show preset details
- [x] `preset open` - Open preset file/directory
- [x] `preset path` - Show preset path
- [x] `preset default` - Set/show default preset
- [x] `preset system` - Manage system prompt (`show/edit/reset`)

---


# v0.3.1 - TUI Refactor & Code Block Removal ✅

**Status:** Completed

**Features:**

### TUI Refactoring
- [x] Replace `msg.String()` key comparisons with `key.Matches()` using centralized `KeyMap` structs
- [x] Decompose monolithic `Update` function into state-specific handlers (`updateKeyMsgReviewing`, `updateKeyMsgChatting`, etc.)
- [x] Decompose monolithic `View` function into state-specific renderers (`viewLoading`, `viewMain`, `viewError`)
- [x] Centralize yank chord state reset logic

### Code Block Feature Removal
- [x] Remove code block navigation (`[`, `]`) and `yb` yank functionality (deferred to v0.6)
- [x] Update all documentation to reflect removal
- [x] Add in-code comments explaining rationale for removal

### Documentation Updates
- [x] Update Coding Styles with TUI key-handling and feature-removal guidelines
- [x] Update help text and footer to remove code block references

---

# v0.3.2 - Prompt Memory

**Status**: Raw ideas, need to review and discuss

**Features**: Reviews Interaction

### Prompt First
- [ ] Prompt first before start the conversation (optional)

### Reviews Interaction
- [ ] System prompt will make sure when review, the content will be break into reviews, for example:
  - [ ] Can navigate and interactive with reviews:
  - [ ] Can "Ignore from context" --> Condense the review and add to current context ignore section
  - [ ] Can 'Ignore from system prompt' --> Add to system prompt ignore section
  - [ ] Can 'Ignore from preset' --> Add to a preset's ignore section

---

# v0.3.3 - Chat Enhancements

**Status:** Planned

**Features:**

### Chat/Request Management (In Testing)
- [ ] `Ctrl+X` cancels streaming requests
- [ ] Prompt history navigation (`Ctrl+P`/`Ctrl+N`)
- [ ] Request cancellation feedback


# v0.3.4 - Extend reading
- Able to read all project for context, then combine with git diff 

# v0.4 - Panes & Export (Lazy-git Style)

**Status:** Planned

**Features:**

### Setting Management
- [ ] Can change default setting (new subcommand)

### Panes Management Mode
- [ ] Multi-pane layout inspired by lazy-git/lazy-docker
- [ ] Panes:
  - Reviews pane (list of reviews in session)
  - Conversation pane (current chat)
  - Config pane (model, API key, style)
- [ ] `Tab` to switch between panes
- [ ] `1/2/3` to jump to specific pane

### Review Actions
- [ ] `a` - Accept/apply suggestion
- [ ] `x` - Reject/ignore suggestion
- [ ] Add to ignore list (global/conversation)
- [ ] Navigate through suggestions with `[` and `]`

### Export & Save
- [ ] `e` - Export current review to file
- [ ] `E` - Export entire conversation
- [ ] Auto-save conversations to `~/.local/share/revcli/`
- [ ] `--format json|markdown` output formats

### Config Management
- [ ] `~/.config/revcli/config.yaml` support
- [ ] Settings: default model, base branch, ignore patterns
- [ ] In-app config editing via config pane

---

# v0.5 - Power User Features

**Status:** Future

**Features:**

### Token Rotation
- [ ] Support multiple API keys
- [ ] Round-robin rotation between keys
- [ ] Auto-switch on rate limit
- [ ] Key usage tracking per key

### Blacklist & Filters
- [ ] Blacklist review styles (e.g., "don't suggest X")
- [ ] Global vs conversation-level blacklist
- [ ] `--min-severity` flag to filter output

### Dry-Run & Preview
- [ ] `--dry-run` / `-n` - Preview payload without API call
- [ ] `--list-models` - Show available models
- [ ] Token cost estimation

# v0.6 - Code Block Management (Deferred)

**Status:** Deferred

**Features:**

### Code Block Highlighting & Navigation
- [ ] Code block detection in review/chat responses
- [ ] Visual highlighting with purple border
- [ ] Navigate with `[` / `]` keys
- [ ] Contextual hints and block indicators
- [ ] `yb` yanks highlighted block
- [ ] Code block index indicator (e.g., "Block 2/5")
- [ ] Jump to specific block with number prefix (e.g., `2]` jumps to block 2)

### Code Block Folding
- [ ] `zc` - Fold/collapse current code block
- [ ] `zo` - Unfold/expand current code block
- [ ] `za` - Toggle fold state
- [ ] `zM` - Fold all code blocks
- [ ] `zR` - Unfold all code blocks
- [ ] Collapsed indicator showing language and line count

---

# v1.0 - Production Ready

**Status:** Future

**Features:**

### Multiple AI Providers
- [ ] OpenAI GPT-4
- [ ] Anthropic Claude
- [ ] Local models (Ollama)
- [ ] `--provider` flag

### Build Mode
- [ ] `revcli build docs` - Generate documentation
- [ ] `revcli build postman` - Generate Postman collections
- [ ] Interactive file/folder selection with Vim navigation
- [ ] Read from controller, serializers, routers
- [ ] After implemented `build mode`, bring file/folder selection feature in `review mode`: add option to include other files than git diff. Interactive files/folders selection like `build mode`

### Team Features
- [ ] Shared config via `.revcli.yaml` in repo
- [ ] Team-specific prompts and rules
- [ ] Pre-commit hook integration

### Advanced UI
- [ ] VS Code extension
- [ ] Review annotations (inline comments)
- [ ] Diff viewer with syntax highlighting

---

# v2.0 - Future Vision

**Status:** Ideas

**Features:**

### Interview Mode
- [ ] `revcli interview` - Practice coding interviews
- [ ] Algorithm questions with hints
- [ ] Code review practice

### Auto-Fix
- [ ] Apply LLM suggestions automatically
- [ ] `--auto-fix` flag for non-breaking changes
- [ ] Git commit integration

### Integrations
- [ ] GitHub Action
- [ ] GitLab CI template
- [ ] PR/MR comment posting
- [ ] SonarQube/CodeClimate integration

---

# Ideas Backlog

> Raw ideas for future consideration

**Presets**
- Remove `built-in` type, built-in treated as custom presets

**Uncategorized**
- Ask user for MR intention/Summary MR intention based on diff change to verify business logic
- Make the base prompt more generic/neutral (Not just Go reviewer)
- Compare two branches directly (`revcli diff main feature-branch`)
- Review specific files only (`revcli review src/api.go`)
- Ignore patterns via `.revignore` file
- Statistics dashboard (reviews done, issues found)
- Multi-language support (i18n for prompts)
- Plugin system for custom analyzers

---

# Known Bugs

> Track and fix these issues

| Bug                                                                       | Status | Notes                                                                                                                                                                            |                                         |
| ------------------------------------------------------------------------- | ------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------- |
| Navigation issues after reviews                                           | Fixed  | No longer auto-scrolls to bottom; users read from top                                                                                                                            |                                         |
| Redundant spaces below terminal                                           | Fixed  | Dynamic viewport height calculation based on UI state                                                                                                                            |                                         |
| Yank only copies initial review                                           | Fixed  | Now yanks full content including chat history                                                                                                                                    |                                         |
| Code block navigation removed                                             | Fixed    | Code block navigation (`[`, `]`, `yb`) removed in v0.3.1; deferred to v0.6 for complexity/UX reasons                                                                          |                                         |
| Panic when using --interactive flag                                       | Fixed  | Added nil checks for renderer fallback                                                                                                                                           |                                         |
| Can't type `?` in chat mode                                               | Fixed  | `?` now only triggers help in reviewing mode, passes through in chat                                                                                                             |                                         |
| Can't press Enter for newline in chat                                     | Fixed  | Changed to `Alt+Enter` to send; Enter creates newlines                                                                                                                           |                                         |
| Textarea has white/highlighted background                                 | Fixed  | Custom textarea styling with rounded borders                                                                                                                                     |                                         |
| Preset create fails with multi-word descriptions                          | Fixed  | `fmt.Scanln()` only reads first word; replaced with `bufio.Reader` for full-line input                                                                                           |                                         |
| Multiple stdin readers cause data loss                                    | Fixed  | Creating multiple `bufio.NewReader(os.Stdin)` instances causes buffered data loss when input is piped. Fixed by creating a single reader and reusing it.                         |                                         |
| IS01: Users can't edit custom presets                                     | Fixed  | Added `preset edit` command to allow editing custom presets interactively                                                                                                        |                                         |
| IS02: Missing feature to edit preset in command line or manually          | Fixed  | Added `preset edit` command and `preset open` command for manual editing                                                                                                         |                                         |
| IS03: Missing feature to open preset folder/file                          | Fixed  | Added `preset open` (opens in editor/file manager) and `preset path` (shows path) commands                                                                                       |                                         |
| IS04: Missing feature to set default preset                               | Fixed  | Added `preset default` command and config.yaml support for default preset                                                                                                        |                                         |
| IS05: Preset gets appended to system prompt, should optionally replace it | Fixed  | Added `replace` field to preset YAML and `--preset-replace` flag to review command                                                                                               |                                         |
| Flag redefined panic in preset command                                    | Fixed  | Duplicate flag definitions in `init()` function caused panic. Fixed by removing duplicate flag registrations. Always check for duplicate flag definitions when adding new flags. |                                         |
| IS06: No helper found when run `rv review -h`                             | Fixed  | Added `--preset-replace` flag to help examples in command Long description                                                                                                       |                                         |
| IS07: Flag is too long, not great for typing                              | Fixed  | Added short alias `-R` for `--preset-replace` flag using `BoolVarP`                                                                                                              |                                         |
| IS08: Missing feature: edit system prompt                                 | Fixed  | Added `preset system` command with `show/edit/reset` subcommands. System prompt can be customized via `~/.config/revcli/presets/system.yaml`                                     |                                         |
| IS09: Fail to edit when enter a different name                            | Fixed  | Disabled name editing in `preset edit` command. Name is now read-only to prevent data loss.                                                                                      |                                         |
| IS10: Not enter default value when editing                                | Fixed  | Improved default value prompts with clearer instructions and explicit current value display.                                                                                     |                                         |
| IS11: Can not move cursor up and down when edit prompt                    | Fixed  | Replaced line-by-line stdin input with external editor (`$EDITOR` or `vi` fallback) for multiline prompt editing.                                                                |                                         |
| IS12: Preset files when saved got `\n` instead of breaks                  | Fixed  | Added custom `MarshalYAML()` method to `Preset` struct to use literal block scalars (`                                                                                           | `) for multiline prompts in YAML files. |

---

