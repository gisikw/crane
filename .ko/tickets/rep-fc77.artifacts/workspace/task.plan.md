rep dev shell
  just test    — run tests
  just build   — build binary
All done. Here's a summary of what was implemented:

**Files changed:**
- `args.go` — Added `ChargeCode string` to `Options`; added `--charge-code` flag parsing with validation (`^[a-zA-Z0-9_-]+$`)
- `adapter.go` — Added `CaptureUsage bool` to `InvokeRequest`; Claude and Cursor switch to `--output-format json`, OpenCode adds `--format json` when set
- `main.go` — When `--charge-code` is set: buffers stdout into a `bytes.Buffer`, calls `parseProviderOutput` to extract text (written to stdout) and usage fields, then appends to JSONL; adds `--charge-code` to usage text
- `INVARIANTS.md` — Documents the sanctioned charge-code exception to the output-interpretation invariant

**New files:**
- `charge.go` — `UsageRecord` struct, `validateChargeCode`, `stateDir`, `AppendChargeRecord`, and per-provider output parsers (`processClaudeOutput`, `processCursorOutput`, `processOpencodeOutput`)
- `charge_test.go` — Full coverage of all charge.go functions including fallback paths and validation

**Key design decisions applied from ticket answers:**
1. Claude stdout → extract `result` text field, write to stdout
2. Cursor/OpenCode → switch to JSON output mode and attempt token extraction; graceful fallback to raw output + sparse record
3. Charge code validation → `[a-zA-Z0-9_-]` only, reject at parse time
