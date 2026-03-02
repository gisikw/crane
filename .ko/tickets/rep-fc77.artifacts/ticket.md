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
