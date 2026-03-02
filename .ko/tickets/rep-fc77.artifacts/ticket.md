---
id: rep-fc77
status: open
deps: []
created: 2026-03-02T15:30:18Z
type: task
priority: 3
---
# Add --charge-code flag with per-run JSONL usage logging

Add a --charge-code <code> flag that tags each invocation for cost attribution. When set, append a JSONL line to $XDG_STATE_HOME/rep/<charge-code>.jsonl after each run with provider, model, token counts, cost (where available), and timestamp. One file per charge code, append-only — no read-modify-write. Rollup happens at retrieval time (see dependent ticket). Think department codes on a xerox machine: the caller decides granularity (e.g. 'exocortex' vs 'exocortex-2026-03-01' vs a UUID). Token/cost data availability varies by provider — Claude gives full JSON with total_cost_usd and breakdowns; Cursor and OpenCode are spottier (see rep-60e8 findings).

## Notes

**2026-03-02 15:49:47 UTC:** Question: Should we validate and sanitize the `--charge-code` string before using it in the file path to prevent path traversal attacks?
Answer: Good call. Let's go alphanumeric, dash, and underscore as our supported

**2026-03-02 15:49:47 UTC:** Question: When `--charge-code` is set with the Claude provider, should we extract the `result` text field from Claude's JSON output and write that to stdout, or pass the full JSON through to the caller?
Answer: Extract text result (Recommended)
Parse Claude JSON, write the `result` field to stdout — maintains consistent caller experience

**2026-03-02 15:49:47 UTC:** Question: For Cursor and OpenCode providers, is it acceptable to log sparse usage records (provider, model, timestamp only, no token counts), or should we attempt to parse their JSON output to extract token data?
Answer: Attempt JSON parsing
Try to extract token counts from Cursor/OpenCode JSON output for more complete records
