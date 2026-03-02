# Invariants

These are architectural contracts. They're not aspirational — they're load-bearing.
Violations are bugs, not style issues. No grandfather clauses.

If an invariant no longer serves the project, remove it explicitly with rationale.
Don't just ignore it.

## Core Model

- Rep is a dispatcher, not an orchestrator. It picks a provider, invokes it, returns the result. No multi-step pipelines, no state between invocations, no retry logic.
- Provider selection is config-driven, not dynamic. A config file says "prefer cursor, fall back to claude, then opencode." Rep walks the list until one works.
- "Works" means the binary exists and is callable. Rep does not check token budgets, subscription status, or API quotas — that's a future concern and will be a separate mechanism if/when it matters.
- Rep does not interpret agent output. It captures stdout/stderr and exit code. The caller decides what to do with it. **Exception**: when `--charge-code` is set and the provider is Claude, Rep buffers stdout, parses the JSON response to extract usage metadata (written to JSONL), and writes the `result` text field to stdout. Cursor and OpenCode output is similarly buffered and parsed when `--charge-code` is set. This is the one sanctioned case where Rep interprets provider output.

## Interface

- Single binary, single command: `rep "do the thing"` — prompt on argv or stdin
- Flags for overrides: `--provider`, `--model`, `--dir`, `--system-prompt`
- Config file at `~/.config/rep/config.toml` for provider preference and defaults
- Exit code mirrors the agent's exit code. Rep's own errors are exit code 1 with a message on stderr.
- Output goes to stdout unmodified. Rep's own diagnostics go to stderr.

## Provider Adapters

- Each provider is a Go function that builds an exec.Cmd from a standard set of inputs (prompt, model, system prompt, working dir, allow-all flag)
- Adding a provider means adding one function and one config key. No interfaces, no registration, no plugins.
- Provider-specific flag translation is the adapter's job. The caller never sees provider details.

## Code Organization

- Decision logic is pure: functions take data, return decisions, no I/O
- I/O is plumbing: thin orchestrators that gather data -> call pure functions -> act
- No multi-purpose functions: separate decision from effect

## File Size

- 500 lines max per file (ergonomic, not aesthetic)
- Split along behavioral seams, not alphabetically
- Tests mirror source files
- No grab-bag utility files

## Secrets

- No hardcoded secrets, tokens, PII, or infrastructure-specific details
- Provider API keys come from environment variables, not config files
- Config file contains preferences and defaults only, never credentials

## Policy

- Decisions that shape code are explicit (here), not implicit
- No "look at how X does it" as policy — write it down or it doesn't exist
