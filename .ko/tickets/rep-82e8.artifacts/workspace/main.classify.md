rep dev shell
  just test    — run tests
  just build   — build binary
Looking at this ticket, I need to classify the type of work:

**Ticket analysis:**
- The ticket starts with "Add a 'rep usage <charge-code>' subcommand"
- It specifies exactly what to build: a new subcommand with defined behavior
- It includes concrete requirements: read JSONL files, aggregate data, display summaries, support a `--json` flag
- The expected output is implemented code

**Classification:** This is a **task**

The ticket clearly describes *what to build* (a new feature), not *what's broken* (which would be a bug) or *whether something is possible* (which would be research). It has specific functional requirements and the expected deliverable is committed code.

```json
{"disposition": "route", "workflow": "task"}
```
