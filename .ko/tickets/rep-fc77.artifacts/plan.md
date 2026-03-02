## Goal

Add `--charge-code <code>` flag that logs provider, model, token counts, cost, and timestamp to an append-only JSONL file per charge code after each invocation.

## Context

**Key files:**
- `args.go`: `Options` struct + `ParseArgs`; new field goes here
- `adapter.go`: `InvokeRequest` struct, `Adapter` interface, per-provider adapters; `Build` returns `*exec.Cmd`
- `main.go`: orchestrates args → provider resolution → adapter.Build → `cmd.Run` with `cmd.Stdout = os.Stdout`
- `INVARIANTS.md`: "Rep does not interpret agent output" and "Output goes to stdout unmodified" — this feature relaxes the interpretation invariant specifically for Claude when `--charge-code` is set

**From rep-60e8 findings:**
- Claude with `--output-format json` gives a JSON object on stdout with `result` (text), `total_cost_usd`, and `usage.{input_tokens, output_tokens, cache_read_input_tokens, cache_creation_input_tokens}`
- Cursor: token data not reliably present in CLI JSON output (active community feature request)
- OpenCode: `--format json` emits events; usable but complex; treat as spottier for now

**Patterns:**
- Tests mirror source files; `assertContains` helper in adapter_test.go; table-driven style in provider_test.go
- XDG_STATE_HOME defaults to `~/.local/state`; state files go to `$XDG_STATE_HOME/rep/`
- File-size limit 500 lines/file; new behavioral seam → new file

## Approach

Add `ChargeCode string` to `Options` and `CaptureUsage bool` to `InvokeRequest`. When `--charge-code` is set, mark the request with `CaptureUsage=true` so the Claude adapter switches to `--output-format json`. In `main.go`, when a charge code is present, buffer the provider's stdout: for Claude, parse the JSON to extract the text result (written to stdout) and usage metadata (written to JSONL); for Cursor/OpenCode, pass stdout through unchanged and log a sparse record (provider, model, timestamp, no token counts). A new `charge.go` owns JSONL writing, state-dir resolution, and Claude output parsing.

## Tasks

1. **[args.go]** — Add `ChargeCode string` to `Options` struct. Add a `--charge-code` case to `ParseArgs` (same pattern as other string flags: check `i+1 < len(args)`, advance `i`, assign). Update the usage text printed in `main.go` to include `--charge-code <code>`.
   Verify: `go test ./...` passes; `TestParseArgs_Flags` still passes with the new field absent.

2. **[adapter.go:InvokeRequest]** — Add `CaptureUsage bool` to `InvokeRequest`. In `ClaudeAdapter.Build`, when `req.CaptureUsage` is true, use `--output-format json` instead of `--output-format text`.  All other adapters ignore `CaptureUsage` (they continue using their current text format).
   Verify: existing adapter tests pass; new `TestClaudeAdapter_CaptureUsage` test confirms the flag appears in args when set.

3. **[charge.go]** (new file) — Define:
   - `UsageRecord` struct with JSON tags: `Timestamp string`, `Provider string`, `Model string`, `ChargeCode string`, `InputTokens *int`, `OutputTokens *int`, `CacheReadTokens *int`, `CacheWriteTokens *int`, `CostUSD *float64`. Use pointer fields with `omitempty` so absent data is omitted from JSON.
   - `stateDir() string` — returns `$XDG_STATE_HOME/rep` if env var is set, else `~/.local/state/rep`.
   - `AppendChargeRecord(chargeCode string, rec UsageRecord) error` — opens `stateDir()/<chargeCode>.jsonl` with `O_APPEND|O_CREATE|O_WRONLY|0600`, creates parent dir if needed (`os.MkdirAll`), marshals `rec` to JSON, writes line + newline. No read-modify-write.
   - `processClaudeOutput(output []byte) (stdout []byte, rec UsageRecord, err error)` — unmarshals JSON, extracts `result` string for stdout (falls back to raw output on parse failure), fills `UsageRecord` fields from `usage` and `total_cost_usd`.
   Verify: unit tests in `charge_test.go` cover each function.

4. **[charge_test.go]** (new file) — Tests:
   - `TestStateDir_Default`: unsets `XDG_STATE_HOME`, verifies path ends in `/.local/state/rep`
   - `TestStateDir_Custom`: sets `XDG_STATE_HOME=/tmp/foo`, verifies path is `/tmp/foo/rep`
   - `TestAppendChargeRecord_CreatesAndAppends`: writes two records to a temp dir, reads file back, confirms two valid JSON lines with expected fields
   - `TestProcessClaudeOutput_Full`: feeds a realistic Claude JSON blob, asserts stdout == result text and all usage fields populated
   - `TestProcessClaudeOutput_Fallback`: feeds invalid JSON, asserts stdout == raw input and `CostUSD == nil`
   Verify: `go test ./...` passes.

5. **[main.go]** — Wire charge-code into the run loop:
   - Pass `ChargeCode: opts.ChargeCode` and `CaptureUsage: opts.ChargeCode != ""` into `InvokeRequest`.
   - When `opts.ChargeCode != ""`: set `cmd.Stdout` to a `bytes.Buffer` instead of `os.Stdout`. After `cmd.Run`, call `processClaudeOutput` for the `claude` provider (and a passthrough for others), write the resulting stdout bytes to `os.Stdout`, then call `AppendChargeRecord` with the extracted `UsageRecord` (provider, model, timestamp via `time.Now().UTC().Format(time.RFC3339)`, plus any token/cost fields). Log JSONL write errors to stderr but do not change exit code.
   - When `opts.ChargeCode == ""`: existing path unchanged (`cmd.Stdout = os.Stdout`).
   Verify: `go build ./...` succeeds; manual smoke test with `--charge-code test` produces a JSONL file in state dir.

6. **[args_test.go]** — Add `TestParseArgs_ChargeCode`: parse `--charge-code exocortex` and verify `opts.ChargeCode == "exocortex"`. Add `TestParseArgs_ChargeCodeMissingValue`: verify error returned when `--charge-code` has no value.
   Verify: new tests pass.

7. **[INVARIANTS.md]** — Under "Core Model", add a sentence explicitly documenting the charge-code exception: when `--charge-code` is set and the provider is Claude, Rep buffers stdout, extracts usage metadata for JSONL logging, and writes the text result to stdout — this is the one sanctioned case where Rep interprets provider output.
   Verify: document reads cleanly; no other invariants contradict.

## Open Questions

1. **Claude stdout format with `--charge-code`**: callers currently get text on stdout. When `--charge-code` is set and provider is Claude, the plan extracts the `result` text field and writes that to stdout — so the caller still sees text. Is this correct, or should the full Claude JSON be passed through? The ticket says "output goes to stdout unmodified" in spirit, but extracting text is the only way to keep the caller experience consistent.

2. **Cursor/OpenCode token data**: the plan logs sparse records (no token counts) for these providers. Is that acceptable, or should we attempt to parse their JSON output too? The ticket says "Cursor and OpenCode are spottier" which suggests sparse records are fine, but confirm.

3. **Charge-code validation**: should we validate the charge-code string (e.g., reject slashes, null bytes, path traversal characters) before constructing the file path? The ticket says callers decide granularity but doesn't address sanitization. Recommend yes — at minimum reject strings containing `/` or `..`.
