rep dev shell
  just test    — run tests
  just build   — build binary
The plan is written to `/home/dev/Projects/rep/.ko/tickets/rep-fc77.artifacts/plan.md`. Here's a summary:

## Plan Summary

**Goal**: Add `--charge-code <code>` that appends a JSONL usage record to `$XDG_STATE_HOME/rep/<code>.jsonl` after each run.

**Key design decision**: This feature requires relaxing the "Rep does not interpret agent output" invariant — specifically for Claude when `--charge-code` is set. The plan captures stdout, extracts the text `result` field for stdout (so callers still see text), and extracts token/cost fields for the JSONL record. Cursor/OpenCode get sparse records (provider, model, timestamp only — no token counts).

**7 tasks:**
1. `args.go` — Add `ChargeCode string` to `Options`, parse `--charge-code` flag
2. `adapter.go` — Add `CaptureUsage bool` to `InvokeRequest`; Claude adapter uses `--output-format json` when set
3. `charge.go` (new) — `UsageRecord` struct, `stateDir()`, `AppendChargeRecord()`, `processClaudeOutput()`
4. `charge_test.go` (new) — Unit tests for all charge.go functions
5. `main.go` — Wire charge-code: buffer stdout when set, parse usage, write text to stdout, append JSONL
6. `args_test.go` — Tests for `--charge-code` parsing
7. `INVARIANTS.md` — Document the sanctioned charge-code exception to the interpretation rule

**3 open questions flagged**: stdout format for Claude callers, whether to attempt Cursor/OpenCode token parsing, and charge-code string sanitization (path traversal).
