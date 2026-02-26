# rep

A thin CLI dispatcher for headless AI agent invocations. Rep normalizes the interface across agent providers (Claude Code, Cursor, OpenCode, etc.), picks the best available one based on a preference list, and gets out of the way.

## Usage

```bash
rep "refactor the auth module to use JWT"
echo "fix the failing test in user_test.go" | rep
rep --provider claude --model opus "review this PR"
rep --provider opencode --model "ollama/qwen3:32b" "fix the build"
rep --dir /path/to/project "run the tests and fix failures"
```

### Flags

- `--provider <name>`: Force a specific provider (claude, cursor, opencode)
- `--model <model>`: Override model selection
- `--dir <path>`: Working directory for the agent
- `--system-prompt <text>`: Append system prompt
- `--no-permissions`: Skip permission checks (default)
- `--with-permissions`: Enable permission prompts
- `--allowed-tools <tools>`: Whitelist specific tools (Claude only, incompatible with `--no-permissions`)
- `--disallowed-tools <tools>`: Blacklist specific tools (Claude only, incompatible with `--no-permissions`)

Tool names can be comma or space-separated. Examples:
```bash
rep --with-permissions --allowed-tools "Bash,Edit,Read" "fix the tests"
rep --with-permissions --allowed-tools "Bash Edit Read" "fix the tests"
rep --with-permissions --disallowed-tools "Write" "review this code"
```

Note: `--allowed-tools` and `--disallowed-tools` only work with the Claude provider and require `--with-permissions` mode. They expose Claude CLI's fine-grained tool permission controls.

## Configuration

Rep reads `~/.config/rep/config.toml`:

```toml
# Provider preference order — first available wins
providers = ["claude", "cursor", "opencode"]

# Per-provider defaults
[claude]
model = "opus"
allow_all = true

[cursor]
allow_all = true

[opencode]
model = "ollama/qwen3-coder-next:latest"
```

## Provider Support

| Provider | Binary | Notes |
|----------|--------|-------|
| Claude Code | `claude` | Subscription-compliant invocation via official CLI |
| Cursor | `cursor-agent` / `agent` | Subscription-compliant invocation via official CLI |
| OpenCode | `opencode` | Local/remote models via Ollama, OpenRouter, etc. |

OpenCode handles tool calls and agentic workflows for models that don't have their own CLI harness. Configure the model endpoint (e.g. a remote Ollama server) in OpenCode's own config at `~/.config/opencode/config.json`.

## Install

```bash
just install   # builds and symlinks to ~/.local/bin
```

Or build manually:

```bash
just build     # produces ./rep
```
