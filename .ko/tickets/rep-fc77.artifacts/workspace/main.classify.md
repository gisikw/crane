rep dev shell
  just test    — run tests
  just build   — build binary
Looking at this ticket:

**Title**: Add --charge-code flag with per-run JSONL usage logging

**Content Analysis**:
- Describes a specific feature to implement: `--charge-code` flag
- Specifies clear requirements:
  - Flag format and purpose
  - Output location and format (JSONL)
  - Data to capture (provider, model, tokens, cost, timestamp)
  - Implementation constraint (append-only)
- References related work (rollup at retrieval, rep-60e8 findings)

This ticket describes **what to build** with clear functional requirements and expected output (committed code). It's not asking "can we do this?" (research), not reporting unexpected behavior (bug), and not investigating "what's broken?"

```json
{"disposition": "route", "workflow": "task"}
```
