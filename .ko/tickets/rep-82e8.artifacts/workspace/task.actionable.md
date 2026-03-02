rep dev shell
  just test    — run tests
  just build   — build binary
## Evaluation

I've reviewed the plan artifact. Here's my assessment:

### Completeness ✓
The plan has:
- **Clear goal**: Add `rep usage <charge-code> [--json]` subcommand with rollup summary
- **Context**: References existing `charge.go` functions, architectural constraints from INVARIANTS.md, and Go testing patterns
- **Approach**: Describes adding subcommand routing before `LoadConfig()` and creating `usage.go` with pure rollup logic
- **Ordered tasks**: 3 specific tasks with clear implementation details and verification steps
- **Verification**: Each task includes clear pass/fail criteria (build success, test coverage, etc.)

### Feasibility ✓
I verified the plan's assumptions against the actual codebase:
- ✓ `charge.go` has `UsageRecord`, `stateDir()`, `validateChargeCode()`, `AppendChargeRecord()`
- ✓ `main.go` exists and confirms: no subcommand routing currently (goes straight to `LoadConfig()` → `ParseArgs()`)
- ✓ INVARIANTS.md exists and the plan follows its constraints
- ✓ All structural assumptions match the codebase

### Open Questions ✓
The plan explicitly states "None" and provides reasonable justification:
- The ticket clearly specifies rollup aggregation (group by provider/model, sum tokens/cost)
- Output modes are explicit (human-readable table vs JSON)
- File path is specified
- Sort order (provider then model) is a natural, deterministic choice

```json
{"disposition": "continue"}
```
