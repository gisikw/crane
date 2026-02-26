The plan has been written to `/home/dev/Projects/rep/.ko/tickets/rep-fa96.artifacts/plan.md`.

## Summary

I've investigated the codebase and created an implementation plan for adding `--allowed-tools` and `--disallowed-tools` flags. Here's what I found:

**Current State:**
- Rep is a CLI dispatcher that normalizes AI agent invocations across providers
- Currently supports binary permission control: `--no-permissions` (skip all) vs `--with-permissions` (prompt for each)
- Claude CLI already has fine-grained tool control via `--allowed-tools` and `--disallowed-tools`

**Planned Implementation:**
1. Add new flag parsing in `args.go` for tool allow/deny lists
2. Thread these through the `InvokeRequest` struct
3. Update `ClaudeAdapter` to pass these flags to the Claude CLI when `AllowAll=false`
4. Add comprehensive test coverage
5. Document the new flags in README

**Key Design Decision:**
These flags only work with the Claude provider (since Cursor and OpenCode don't have equivalent functionality) and are incompatible with `--no-permissions` mode. This aligns with rep's architecture where provider-specific features are handled in the adapters.

The plan is ready for your review at the artifact directory.
