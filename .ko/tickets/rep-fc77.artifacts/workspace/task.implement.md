rep dev shell
  just test    — run tests
  just build   — build binary
All tests pass and the build succeeds. The implementation was already completed by a prior build attempt. Here's a summary of what exists across all tasks:

**Task 1 — args.go**: `ChargeCode string` added to `Options`; `--charge-code` flag parsing with `validateChargeCode` call at parse time.

**Task 2 — adapter.go**: `CaptureUsage bool` added to `InvokeRequest`; `ClaudeAdapter` and `CursorAdapter` switch to `--output-format json`, `OpencodeAdapter` adds `--format json` when `CaptureUsage` is true.

**Task 3 — charge.go** (new): `UsageRecord` struct with pointer fields + `omitempty`; `validateChargeCode` with `^[a-zA-Z0-9_-]+$` regex; `stateDir()` with XDG fallback; `AppendChargeRecord` with `O_APPEND|O_CREATE|O_WRONLY`; `processClaudeOutput`, `processCursorOutput`, `processOpencodeOutput` with graceful fallback on parse failure.

**Task 4 — charge_test.go** (new): Full coverage — `TestStateDir_Default/Custom`, `TestAppendChargeRecord_CreatesAndAppends`, `TestProcessClaudeOutput_Full/Fallback`, `TestProcessCursorOutput_Full/FlatTokens/Fallback`, `TestProcessOpencodeOutput_Events/Fallback`, `TestValidateChargeCode_Valid/Invalid`.

**Task 5 — main.go**: When `ChargeCode != ""`: buffers stdout into `bytes.Buffer`, calls `parseProviderOutput` to dispatch to per-provider parser, writes text to stdout, appends `UsageRecord` to JSONL; otherwise existing path unchanged. Usage text includes `--charge-code <code>`.

**Task 6 — args_test.go**: `TestParseArgs_ChargeCode`, `TestParseArgs_ChargeCodeMissingValue`, `TestParseArgs_ChargeCodeInvalidChars` all present and passing.

**Task 7 — INVARIANTS.md**: Charge-code exception documented in the "Core Model" section.

No deviations from the plan.
