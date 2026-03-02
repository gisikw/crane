rep dev shell
  just test    — run tests
  just build   — build binary
Plan written to `.ko/tickets/rep-82e8.artifacts/plan.md`. Here's a summary:

## Plan Summary

**Goal:** Add `rep usage <charge-code> [--json]` subcommand that reads the JSONL log and displays a rolled-up summary.

**Key findings from codebase:**
- `charge.go` already has `UsageRecord`, `stateDir()`, and validation — reuse directly.
- `main.go` has no subcommand routing yet — needs detection before `LoadConfig()`/`ParseArgs()`.
- INVARIANTS require pure decision functions separated from I/O plumbing.

**Three tasks:**

1. **`main.go`** — Add subcommand detection at the start of `main()`: if `os.Args[1] == "usage"`, dispatch to `cmdUsage(os.Args[2:])` and return early.

2. **`usage.go` (new)** — Contains:
   - `RollupEntry` / `RollupResult` structs
   - `Rollup(chargeCode string, records []UsageRecord) RollupResult` — pure aggregation, sorts breakdown by provider then model
   - `ReadChargeRecords(chargeCode string) ([]UsageRecord, error)` — reads JSONL, skips unparseable lines, returns empty (not error) on missing file
   - `FormatRollup(r RollupResult) string` — human-readable table
   - `cmdUsage(args []string)` — I/O orchestrator, handles `--json` flag

3. **`usage_test.go` (new)** — Tests for empty/single/multi-group rollup, nil pointer safety, file-not-exist, JSONL parsing, and bad-line skipping.

**No open questions** — all requirements are clearly specified.
