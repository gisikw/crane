Looking at this ticket, I need to classify what kind of work it requires.

The ticket states:
> "Identify how we might expose token spend in a response, for each of the underlying harnesses we support - this might require agent-runner-specific logic for each"

Key characteristics:
- It's asking to **"identify how"** — this is exploratory and investigative
- It's a question about **capability** ("might expose") rather than a directive to implement
- The output will be **findings and approach**, not committed code
- It acknowledges uncertainty ("might require agent-runner-specific logic") and needs investigation to determine feasibility and design

This is clearly a **research** ticket. The team needs to:
1. Investigate the current harnesses/runners supported
2. Understand how token spend information is available from each
3. Determine the best approach to expose this in responses
4. Document findings about what's possible and what architectural changes might be needed

Once research is complete, implementation would likely follow as a separate task ticket.

```json
{"disposition": "route", "workflow": "research"}
```
