# Implementation Summary: rep-fc77

## What Was Done

Added `--charge-code <code>` flag that logs per-invocation usage data to an append-only JSONL file at `$XDG_STATE_HOME/rep/<charge-code>.jsonl`.

**Files changed:**
- `args.go` — `ChargeCode string` added to `Options`; `--charge-code` flag parsed with `validateChargeCode` at parse time; usage text updated.
- `adapter.go` — `CaptureUsage bool` added to `InvokeRequest`; Claude and Cursor switch to `--output-format json`, OpenCode adds `--format json` when set.
- `main.go` — When `ChargeCode != ""`: buffers stdout, dispatches to `parseProviderOutput`, writes extracted text to stdout, appends `UsageRecord` to JSONL. JSONL write errors go to stderr without changing exit code.
- `INVARIANTS.md` — Documents the sanctioned charge-code exception to the "Rep does not interpret agent output" invariant, covering all three providers.

**New files:**
- `charge.go` — `UsageRecord` struct (pointer fields + omitempty for sparse data), `validateChargeCode` (`^[a-zA-Z0-9_-]+$`), `stateDir()` with XDG fallback, `AppendChargeRecord` (O_APPEND, no read-modify-write), `processClaudeOutput` / `processCursorOutput` / `processOpencodeOutput` with graceful fallback on parse failure.
- `charge_test.go` — Full coverage of all charge.go functions including fallback paths and both valid/invalid charge-code cases.

## Notable Decisions

1. **JSONL written even on provider failure** — The record is appended regardless of provider exit code (runErr is checked after JSONL write). This matches "after each run" semantics and records partial data (e.g. a killed run still consumed tokens).

2. **Cursor CaptureUsage switches to `--output-format json`** — The plan originally scoped this to Claude only but the ticket answers explicitly asked for JSON parsing on Cursor/OpenCode too. INVARIANTS.md was updated to reflect all three providers.

3. **`processOpencodeOutput` accumulates cost from first non-zero event only** — OpenCode may emit cost on multiple events; the implementation takes the first non-zero value to avoid double-counting.

## Fix Applied During Review

The plan's Task 2 verification step required `TestClaudeAdapter_CaptureUsage`, `TestCursorAdapter_CaptureUsage`, and `TestOpencodeAdapter_CaptureUsage` tests in `adapter_test.go`. These were omitted by the implementation stage. They were added during review; all tests pass.
