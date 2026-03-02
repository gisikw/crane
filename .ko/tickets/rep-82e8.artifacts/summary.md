# After-Action Summary — rep-82e8

## What Was Done

Added `rep usage <charge-code> [--json]` subcommand per the ticket spec.

**main.go**: Added early subcommand dispatch at the top of `main()` before `LoadConfig()`/`ParseArgs()` — if `os.Args[1] == "usage"`, routes to `cmdUsage(os.Args[2:])` and returns. Also added the usage subcommand to the help output.

**usage.go** (new, 219 lines): Implements:
- `RollupEntry` / `RollupResult` structs with JSON tags
- `Rollup(chargeCode, records)` — pure function, groups by (provider, model), accumulates per-group and grand totals, returns deterministically sorted breakdown (provider then model)
- `ReadChargeRecords(chargeCode)` — opens `stateDir()/<chargeCode>.jsonl`, reads line-by-line with `bufio.Scanner`, skips blank/unparseable lines, returns empty slice (not error) when file absent
- `FormatRollup(r)` — human-readable table; omits cache columns when all cache tokens are zero
- `cmdUsage(args)` — arg parsing, charge-code validation, I/O orchestration, `--json` flag support

**usage_test.go** (new, 214 lines): All 7 required tests implemented and passing:
- `TestRollup_Empty`, `TestRollup_SingleRecord`, `TestRollup_MultipleProviderModel`, `TestRollup_NilTokenFields`
- `TestReadChargeRecords_FileNotExist`, `TestReadChargeRecords_ReadsAndParsesLines`, `TestReadChargeRecords_SkipsBadLines`

## Notable Decisions

- **`FormatRollup` is pure** (returns string, no I/O), consistent with the INVARIANTS.md "decision logic is pure" rule. `cmdUsage` is the thin I/O orchestrator that calls it.
- Cache columns are omitted from the human-readable table when all cache token counts are zero, keeping output clean for the common (non-Claude) case.
- `--json` flag is detected with simple linear arg scan rather than `flag.FlagSet`, which is appropriate for a subcommand with one positional and one flag.

## Invariant Compliance

- All files ≤500 lines ✓
- Decision logic pure, I/O is thin plumbing ✓
- No multi-purpose functions ✓
- No new dependencies ✓
- No secrets ✓

## Build / Test Status

`go build ./...` and `go test ./...` both pass cleanly.
