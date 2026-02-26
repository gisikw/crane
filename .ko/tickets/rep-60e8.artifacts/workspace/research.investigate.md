Perfect! Let me create a brief summary for the user.

## Investigation Complete

I've thoroughly investigated how token spend can be exposed across Rep's three supported providers. The findings are documented in `/home/dev/Projects/rep/.ko/tickets/rep-60e8.artifacts/findings.md`.

### Key Findings:

**All three providers support token usage data, but with different mechanisms:**

1. **Claude Code** ✅ - Fully supported via `--output-format json`
   - Provides comprehensive token counts, cache usage, and USD cost
   - Returns structured JSON with `total_cost_usd`, `usage`, and per-model breakdown

2. **Cursor Agent** ⚠️ - Partially supported
   - Token data captured in session files
   - JSON output support exists but token inclusion is evolving (community feature request active as of Dec 2025)
   - May require session file parsing as fallback

3. **OpenCode** ✅ - Fully supported via multiple methods
   - `--format json` flag for event-based output
   - Dedicated `opencode stats` command for aggregate statistics
   - `opencode export` for complete session data

### The Core Architectural Question:

This investigation reveals a **design tension with Rep's invariants**. The INVARIANTS.md states "Rep does not interpret agent output" and "Output goes to stdout unmodified." Exposing token data requires parsing provider output, which conflicts with Rep's philosophy as a pure dispatcher.

I've outlined **5 different architectural approaches** in the findings document, ranging from minimal changes (just switch to JSON output, let callers parse) to full normalization (Rep parses and wraps with unified token metadata).

### Recommended Next Steps:

1. Decide which architectural approach aligns with Rep's goals
2. Test cursor-agent's actual JSON output with authentication
3. Either update INVARIANTS.md to allow token metadata extraction, or keep Rep as a pure dispatcher and document that callers must handle provider-specific JSON parsing

The detailed findings include specific code locations, JSON structure examples, implementation paths for each adapter, and a comparative analysis of all five architectural options.
