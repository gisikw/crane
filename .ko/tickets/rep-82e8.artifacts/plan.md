## Goal
Add `rep usage <charge-code> [--json]` subcommand that reads the JSONL log file and displays a rolled-up summary grouped by provider/model.

## Context
- `charge.go` already defines `UsageRecord`, `stateDir()`, `validateChargeCode()`, and `AppendChargeRecord()`. The JSONL file lives at `stateDir()/<charge-code>.jsonl`.
- `main.go` currently has no subcommand routing — it always runs the agent invocation path through `ParseArgs()` and `LoadConfig()`. Subcommand detection must happen before those calls.
- INVARIANTS.md requires: decision logic is pure (data in, data out, no I/O); I/O is thin plumbing; no multi-purpose functions; files ≤500 lines.
- Testing pattern: table-driven tests in `*_test.go` files mirroring source files, using `t.TempDir()` and `t.Setenv()` for isolation.
- No external dependencies — only stdlib. The module has only one non-stdlib dep (BurntSushi/toml), which is for config parsing. Usage rollup needs none of that.

## Approach
Add subcommand detection at the top of `main()` by inspecting `os.Args[1]` before calling `LoadConfig()` or `ParseArgs()`. Create a new `usage.go` with pure rollup logic (`Rollup([]UsageRecord) RollupResult`) and I/O helpers (`ReadChargeRecords`, `cmdUsage`). Add `usage_test.go` covering the rollup logic and the JSONL reader.

## Tasks

1. **[main.go:main]** — Add subcommand routing at the very start of `main()`, before `LoadConfig()`. If `len(os.Args) >= 2 && os.Args[1] == "usage"`, call `cmdUsage(os.Args[2:])` and return. Also add `rep usage <charge-code>` to the help/usage error output.
   Verify: `go build ./...` passes; existing tests unbroken.

2. **[usage.go — new file]** — Create with the following pieces, in order:
   - `RollupEntry` struct: `Provider`, `Model` string fields; `Invocations` int; `InputTokens`, `OutputTokens`, `CacheReadTokens`, `CacheWriteTokens` int; `CostUSD` float64.
   - `RollupResult` struct: `ChargeCode` string; `Invocations` int; `TotalInputTokens`, `TotalOutputTokens`, `TotalCacheReadTokens`, `TotalCacheWriteTokens` int; `TotalCostUSD` float64; `Breakdown []RollupEntry` (ordered deterministically: sort by provider then model).
   - `Rollup(chargeCode string, records []UsageRecord) RollupResult` — pure function; iterates records, groups by (provider, model) key, accumulates per-group and grand totals; returns sorted breakdown.
   - `ReadChargeRecords(chargeCode string) ([]UsageRecord, error)` — opens `stateDir()/<chargeCode>.jsonl`, reads line-by-line with `bufio.Scanner`, JSON-unmarshals each line, skips blank/unparseable lines, returns slice. Returns empty slice (not error) if file doesn't exist.
   - `FormatRollup(r RollupResult) string` — builds human-readable table: header line with charge code, invocation count, total cost; then a row per breakdown entry showing provider, model, invocations, input/output tokens, cache tokens (omit cache columns if all zero), cost.
   - `cmdUsage(args []string)` — parse `args` for `<charge-code>` positional and optional `--json` flag; validate charge code; call `ReadChargeRecords`; call `Rollup`; if `--json`, `json.Marshal` and print; otherwise call `FormatRollup` and print. Exit 1 on error with message to stderr.
   Verify: `go build ./...` passes.

3. **[usage_test.go — new file]** — Tests covering:
   - `TestRollup_Empty`: empty records → zero-value result with empty breakdown.
   - `TestRollup_SingleRecord`: one record → breakdown with one entry, totals match.
   - `TestRollup_MultipleProviderModel`: records across two (provider, model) pairs → correct per-group and grand totals; breakdown sorted by provider then model.
   - `TestRollup_NilTokenFields`: records with nil token/cost pointers → treated as zero (no panic).
   - `TestReadChargeRecords_FileNotExist`: returns empty slice, no error.
   - `TestReadChargeRecords_ReadsAndParsesLines`: writes a temp JSONL file, reads it back, verifies fields.
   - `TestReadChargeRecords_SkipsBadLines`: JSONL with one bad line → skips it, returns the good records.
   Verify: `go test ./...` passes with all new and existing tests green.

## Open Questions
None. The ticket specifies the rollup aggregation, output modes, and file path clearly. The JSON schema for `--json` output maps directly to `RollupResult` with `json` struct tags. Sort order for breakdown (provider then model) is the obvious stable choice not explicitly stated but deterministic enough to implement without asking.
