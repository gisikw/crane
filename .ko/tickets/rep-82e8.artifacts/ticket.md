---
id: rep-82e8
status: open
deps: [rep-fc77]
created: 2026-03-02T15:30:37Z
type: task
priority: 3
---
# Add rep usage <code> subcommand for retrieving accumulated usage

Add a 'rep usage <charge-code>' subcommand that reads $XDG_STATE_HOME/rep/<charge-code>.jsonl, rolls up the per-run lines, and displays a human-readable summary (provider/model breakdown, total tokens, total cost, invocation count). Add a --json flag to output the rolled-up data as JSON instead. Rollup aggregation: group by provider then model, sum tokens/cost/invocations.
