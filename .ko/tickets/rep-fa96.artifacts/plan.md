## Goal
Add support for `--allowed-tools` and `--disallowed-tools` flags to enable selective tool permissions for Claude Code invocations.

## Context
The codebase is a thin CLI dispatcher that normalizes invocations across different AI agent providers (Claude, Cursor, OpenCode). The architecture follows clear separation:

- `args.go`: CLI argument parsing into an `Options` struct
- `config.go`: TOML configuration loading with per-provider defaults
- `adapter.go`: Provider-specific command builders implementing the `Adapter` interface
- `main.go`: Orchestration that glues config, args, and adapter together

Currently, rep supports binary permission control via `--no-permissions` (AllowAll=true, default) and `--with-permissions` (AllowAll=false). This maps to Claude's `--dangerously-skip-permissions` flag.

Claude CLI supports fine-grained tool control:
- `--allowed-tools <tools...>`: whitelist (e.g., "Bash(git:*) Edit")
- `--disallowed-tools <tools...>`: blacklist (e.g., "Bash(git:*) Edit")

These are mutually exclusive with `--dangerously-skip-permissions`.

The existing test pattern uses table-driven tests in `*_test.go` files mirroring source files. Tests verify command construction by inspecting `cmd.Args`.

## Approach
Add new CLI flags `--allowed-tools` and `--disallowed-tools` that accept a comma or space-separated list of tool names. Store these in `Options` as string slices. Thread them through `InvokeRequest` to the Claude adapter. The Claude adapter builds the command with these flags when provided, but only if AllowAll is false (they're incompatible with skip-permissions). For other providers (Cursor, OpenCode), ignore these flags since they don't have equivalent functionality.

## Tasks
1. [args.go:Options] — Add `AllowedTools []string` and `DisallowedTools []string` fields to Options struct.
   Verify: code compiles.

2. [args.go:ParseArgs] — Add flag parsing for `--allowed-tools` and `--disallowed-tools`. Accept comma or space-separated values. Store as slice of strings.
   Verify: `go test ./args_test.go` passes.

3. [args_test.go] — Add test cases for parsing `--allowed-tools` and `--disallowed-tools` flags. Verify both single and multiple tool names are parsed correctly.
   Verify: new test passes, existing tests unbroken.

4. [adapter.go:InvokeRequest] — Add `AllowedTools []string` and `DisallowedTools []string` fields to InvokeRequest struct.
   Verify: code compiles.

5. [adapter.go:ClaudeAdapter.Build] — When AllowedTools is non-empty and AllowAll is false, append `--allowed-tools` flag followed by the tool list. When DisallowedTools is non-empty and AllowAll is false, append `--disallowed-tools` flag followed by the tool list. Do not add these flags when AllowAll is true (incompatible with skip-permissions).
   Verify: code compiles.

6. [adapter_test.go] — Add test case verifying ClaudeAdapter includes `--allowed-tools` with correct values when AllowAll=false and AllowedTools is set. Add test case verifying `--disallowed-tools` similarly. Add test case verifying neither flag appears when AllowAll=true.
   Verify: `go test ./adapter_test.go` passes.

7. [main.go] — Thread `opts.AllowedTools` and `opts.DisallowedTools` into the InvokeRequest construction.
   Verify: `go test ./...` passes, `go build` succeeds.

8. [README.md] — Add documentation for `--allowed-tools` and `--disallowed-tools` flags in the Usage section. Note they only work with Claude provider and are incompatible with `--no-permissions`.
   Verify: manually review documentation clarity.

## Open Questions
None. The implementation is straightforward — it's exposing existing Claude CLI functionality through rep's interface. The flags are provider-specific (Claude only) which aligns with rep's design of provider-specific flag translation in adapters. The mutual exclusivity with `--dangerously-skip-permissions` is handled by conditional logic in the adapter.
