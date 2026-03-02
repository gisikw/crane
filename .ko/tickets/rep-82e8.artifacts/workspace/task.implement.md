rep dev shell
  just test    — run tests
  just build   — build binary
All three tasks complete. Here's the summary:

**Task 1 — `main.go`:** Added subcommand detection at the top of `main()`. If `os.Args[1] == "usage"`, dispatches to `cmdUsage(os.Args[2:])` and returns. Also added `rep usage <charge-code> [--json]` to the help output.

**Task 2 — `usage.go` (new):**
- `RollupEntry` / `RollupResult` structs with JSON tags
- `Rollup(chargeCode string, records []UsageRecord) RollupResult` — pure aggregation, groups by (provider, model), nil pointer fields treated as zero, breakdown sorted by provider then model
- `ReadChargeRecords(chargeCode string)` — reads JSONL via `bufio.Scanner`, skips blank/unparseable lines, returns empty slice (no error) on missing file
- `FormatRollup(r RollupResult) string` — human-readable table; omits cache columns when all cache totals are zero
- `cmdUsage(args []string)` — parses positional charge-code and `--json` flag, validates, reads, rolls up, and prints

**Task 3 — `usage_test.go` (new):** Seven tests covering empty rollup, single record, multi-(provider, model) with sort verification, nil pointer safety, file-not-exist, JSONL round-trip, and bad-line skipping. All pass.

No deviations from the plan.
