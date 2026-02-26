# Implementation Summary: Selective Tool Permissions

## What Was Done

Successfully implemented `--allowed-tools` and `--disallowed-tools` CLI flags to enable fine-grained tool permission control when using Claude Code through rep. This exposes Claude CLI's native tool filtering capabilities while maintaining provider independence.

## Changes Made

### 1. Core Data Structures
- **args.go**: Added `AllowedTools` and `DisallowedTools` fields to `Options` struct
- **adapter.go**: Added matching fields to `InvokeRequest` struct

### 2. Flag Parsing (args.go)
- Implemented `--allowed-tools` and `--disallowed-tools` flag parsing
- Created `parseToolList()` helper function that handles both comma-separated ("Bash,Edit,Read") and space-separated ("Bash Edit Read") tool names
- Proper error handling for missing flag values

### 3. Claude Adapter Integration (adapter.go)
- Modified `ClaudeAdapter.Build()` to conditionally include tool flags
- **Key logic**: Tool flags are only added when `AllowAll=false` (permission mode), ensuring mutual exclusivity with `--dangerously-skip-permissions`
- Tool lists are passed directly to Claude CLI without transformation

### 4. Main Orchestration (main.go)
- Threaded `AllowedTools` and `DisallowedTools` from `Options` to `InvokeRequest`

### 5. Test Coverage (args_test.go, adapter_test.go)
- Added 5 args tests covering single/multiple tool parsing for both flags
- Added 3 adapter tests verifying:
  - `--allowed-tools` passes through correctly
  - `--disallowed-tools` passes through correctly
  - Tool flags are ignored when `AllowAll=true` (no conflict with skip-permissions)

### 6. Documentation (README.md)
- Added comprehensive "Flags" section documenting all CLI options
- Included usage examples for tool filtering
- Clearly noted Claude-only limitation and `--no-permissions` incompatibility

## Implementation Decisions

1. **Tool List Parsing**: Chose to support both comma and space-separated formats for flexibility. The `parseToolList()` function handles both, making the CLI more forgiving.

2. **Conditional Logic**: Placed mutual exclusivity logic in the adapter's `Build()` method where command construction happens. This keeps the concern localized and provider-specific.

3. **Pass-Through Design**: Tool names are passed directly to Claude CLI without validation or transformation. This maintains rep's role as a dispatcher, not an interpreter, and ensures forward compatibility as Claude adds new tools.

4. **Provider Isolation**: Only the Claude adapter handles these flags. Cursor and OpenCode adapters ignore them, consistent with rep's provider-specific adaptation strategy.

## Verification

- ✅ All 24 tests pass (5 new args tests, 3 new adapter tests)
- ✅ Code compiles successfully
- ✅ No INVARIANTS.md violations:
  - Remains a dispatcher (no multi-step logic)
  - Decision logic is pure (`parseToolList()`)
  - Provider-specific logic isolated in adapter
  - Single responsibility maintained
- ✅ Plan compliance: All 8 tasks completed exactly as specified
- ✅ No scope creep: Implementation strictly within ticket boundaries

## Notes for Future Readers

- The flags only affect Claude provider because Cursor and OpenCode don't have equivalent functionality
- The mutual exclusivity with `--no-permissions` is enforced at command-build time, not at flag-parse time, which keeps flag parsing simple and provider-agnostic
- Tool name validation is delegated to Claude CLI—rep doesn't maintain a tool registry
- Format flexibility (comma vs space) was intentional to reduce user friction
