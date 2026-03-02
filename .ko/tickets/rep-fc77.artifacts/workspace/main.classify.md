rep dev shell
  just test    — run tests
  just build   — build binary
Looking at this ticket:

**Title**: "Add --charge-code flag with per-run JSONL usage logging"

**What it asks for**: Implement a new feature with specific requirements:
- Add a `--charge-code` CLI flag
- Log usage data (provider, model, tokens, cost, timestamp) to JSONL files
- Validate/sanitize the charge code string
- Extract and format output appropriately for each provider
- Append-only, one file per charge code

**Decision status**: The design decisions are already made (shown in the Notes section):
- Path traversal validation: alphanumeric, dash, underscore only
- Claude output: extract text result to stdout
- Cursor/OpenCode: attempt JSON parsing for token counts

This is clearly asking for **implementation of a specified feature**, not investigation or diagnosis.

```json
{"disposition": "route", "workflow": "task"}
```
