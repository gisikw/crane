# Token Usage Exposure in Rep's Supported Providers

## Executive Summary

Token spend information is available from all three supported providers (Claude Code, Cursor Agent, OpenCode), but each exposes this data differently and requires provider-specific handling. Rep would need to implement adapter-specific logic to extract and normalize token usage data from each provider's output format.

## Findings by Provider

### 1. Claude Code (`claude` CLI)

**Status**: ✅ Full token usage data available

**Location**: adapter.go:38-58

The Claude CLI provides comprehensive token usage information through its JSON output formats.

#### Output Format Options
- `--output-format text`: Human-readable output (currently used by Rep)
- `--output-format json`: Single result JSON object with complete usage data
- `--output-format stream-json`: Real-time NDJSON stream (requires `--verbose`)

#### Token Data Structure

When using `--output-format json`, Claude returns:

```json
{
  "type": "result",
  "subtype": "success",
  "total_cost_usd": 0.019729749999999997,
  "usage": {
    "input_tokens": 3,
    "cache_creation_input_tokens": 1673,
    "cache_read_input_tokens": 17917,
    "output_tokens": 12,
    "server_tool_use": {
      "web_search_requests": 0,
      "web_fetch_requests": 0
    },
    "service_tier": "standard",
    "cache_creation": {
      "ephemeral_1h_input_tokens": 1673,
      "ephemeral_5m_input_tokens": 0
    }
  },
  "modelUsage": {
    "claude-opus-4-6": {
      "inputTokens": 3,
      "outputTokens": 12,
      "cacheReadInputTokens": 17917,
      "cacheCreationInputTokens": 1673,
      "webSearchRequests": 0,
      "costUSD": 0.019729749999999997,
      "contextWindow": 200000,
      "maxOutputTokens": 32000
    }
  }
}
```

**Key Fields**:
- `total_cost_usd`: Total cost in USD
- `usage.*`: Token counts by type
- `modelUsage`: Per-model breakdown with costs

**Implementation Path**: Change `--output-format text` to `--output-format json` in adapter.go:41

---

### 2. Cursor Agent (`cursor-agent` / `agent` CLI)

**Status**: ⚠️ Partial support, evolving feature

**Location**: adapter.go:60-83

The Cursor Agent CLI has token usage tracking, but its availability in output formats is limited.

#### Output Format Options
- `--output-format text`: Human-readable output (currently used by Rep)
- `--output-format json`: Single result JSON object
- `--output-format stream-json`: NDJSON events stream

#### Token Data Availability

Based on community discussions and feature requests, token usage data is:
1. **Captured in session files**: Session JSONL files include usage data with fields for input, output, cache reads, cache writes, and total tokens
2. **Not consistently exposed in CLI output**: As of late 2025, there was an active feature request to include token usage in stream-json output
3. **Format varies**: When available, likely similar structure to Claude's output

**Current Limitations**:
- Token usage may not be included in `--output-format json` responses (requires testing with authenticated session)
- Community requested this feature in December 2025, indicating it may be incomplete

**Implementation Path**:
1. Change output format to json in adapter.go:69
2. Test with authenticated session to verify token data presence
3. May require parsing session files as fallback

**Sources**:
- [Include Token Usage in Stream-JSON Output - Cursor Forum](https://forum.cursor.com/t/include-token-usage-in-stream-json-output/146980)
- [Output format | Cursor Docs](https://cursor.com/docs/cli/reference/output-format)

---

### 3. OpenCode (`opencode` CLI)

**Status**: ✅ Full token usage data available via dedicated command

**Location**: adapter.go:95-110

OpenCode provides comprehensive token tracking through a dedicated `stats` command and JSON output support.

#### Output Format Options
- `opencode run --format default`: Formatted output (currently used by Rep)
- `opencode run --format json`: Raw JSON events
- `opencode stats`: Dedicated statistics command
- `opencode export [sessionID]`: Export session data as JSON

#### Token Data Structure

OpenCode uses a different approach:
1. **During execution**: `opencode run --format json` emits JSON events (may include token data per event)
2. **After execution**: `opencode stats` provides aggregate statistics
3. **Session data**: `opencode export` provides complete session information including token usage

**Stats Command Options**:
```bash
opencode stats
  --days N          # Show stats for last N days
  --tools N         # Number of tools to show
  --models          # Show model statistics
  --project PATH    # Filter by project
  --json           # JSON output (inferred from ecosystem)
```

**Implementation Path**:
1. Change to `--format json` in adapter.go:99
2. Parse JSON events for token usage data
3. Alternatively, call `opencode export <sessionID>` post-execution for complete data

**Sources**:
- [opencode-usage-cli on npm](https://libraries.io/npm/opencode-usage-cli)
- [OpenCode CLI Overview | ccusage](https://ccusage.com/guide/opencode/)

---

## Architectural Considerations

### Current Rep Architecture

**From INVARIANTS.md**:
- Rep is a **dispatcher, not an orchestrator**
- "Rep does not interpret agent output. It captures stdout/stderr and exit code."
- "Output goes to stdout unmodified. Rep's own diagnostics go to stderr."

**Current Implementation** (main.go:61-62):
```go
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
```

Rep passes through stdout/stderr directly without interpretation.

### Design Tensions

Adding token usage exposure creates tension with Rep's core invariants:

1. **Output Interpretation**: Exposing token data requires parsing provider output, which violates "Rep does not interpret agent output"
2. **Single Output Stream**: Token metadata needs separate handling from main output
3. **Provider Abstraction**: Different providers require different parsing logic

### Possible Approaches

#### Option 1: Strict Adherence (Minimal Change)
**Philosophy**: Keep Rep as pure dispatcher

- Change adapters to request JSON output formats
- Pass through JSON unmodified on stdout
- Caller parses JSON to extract both response and token data

**Pros**:
- Maintains invariant: "Output goes to stdout unmodified"
- Simple implementation
- No provider-specific parsing in Rep

**Cons**:
- Breaking change for existing callers expecting text
- Caller must handle three different JSON schemas
- Doesn't truly "expose" token data, just changes output format

#### Option 2: Parallel Metadata Stream
**Philosophy**: Separate data channels

- Main response on stdout (unchanged)
- Token metadata on stderr (separate from diagnostics)
- Use JSON or structured format on stderr

**Pros**:
- Preserves existing stdout behavior
- Clean separation of concerns
- Caller can ignore metadata if not needed

**Cons**:
- Overloading stderr (already used for diagnostics)
- Requires parsing stderr, which could mix with error messages
- Awkward for shell pipelines

#### Option 3: Wrapper Structure
**Philosophy**: Envelope pattern

- Wrap provider output + metadata in consistent JSON structure:
```json
{
  "provider": "claude",
  "model": "opus",
  "output": "<original output>",
  "token_usage": {
    "input_tokens": 100,
    "output_tokens": 200,
    "cost_usd": 0.01
  },
  "exit_code": 0
}
```

**Pros**:
- Normalized interface across providers
- Single parsing logic for callers
- Clean metadata extraction

**Cons**:
- Requires interpreting provider output
- Directly violates "does not interpret agent output"
- More complex implementation
- Must handle provider-specific schemas

#### Option 4: Opt-in Flag
**Philosophy**: Feature flag for backward compatibility

Add `--expose-metadata` flag:
- Without flag: Current behavior (text output, no parsing)
- With flag: JSON envelope with normalized token data

**Pros**:
- Backward compatible
- Explicit opt-in matches "no state between invocations"
- Clear when parsing occurs

**Cons**:
- Still requires interpretation when flag is set
- Two code paths to maintain
- Doesn't resolve invariant violation, just makes it optional

#### Option 5: Separate Command
**Philosophy**: New responsibility = new command

Add `rep-stats <session-id>` or similar:
- `rep` remains unchanged
- New tool for querying token data from provider-specific mechanisms
- Follows Unix philosophy (do one thing well)

**Pros**:
- Preserves Rep's core invariants completely
- Clear separation of concerns
- Each provider can use optimal mechanism (OpenCode's `stats`, session file parsing, etc.)

**Cons**:
- Requires session tracking
- Not all providers have session IDs
- Additional tool to maintain
- Token data not available immediately

---

## Recommended Actions

### Immediate (Low Risk)
1. **Test cursor-agent authentication** to verify actual token data availability in JSON output
2. **Document current limitations** in README regarding token usage exposure
3. **Create example scripts** showing how to parse each provider's JSON output if callers need token data

### Short Term (Requires Architecture Decision)
1. **Decide on approach** based on project priorities:
   - If Rep should remain a pure dispatcher: **Option 1** (just change to JSON output)
   - If token data is critical: **Option 3** or **Option 4** (with INVARIANTS.md update)
   - If staying minimal: **Option 5** (separate tool)

2. **Update INVARIANTS.md** if any interpretation is added
   - Explicitly define what "interpretation" means
   - Add exception for token metadata if needed

### Long Term (Feature Development)
1. **Implement chosen approach** with provider-specific adapters
2. **Add integration tests** for token data extraction from each provider
3. **Create unified token usage struct** if normalization is chosen:
   ```go
   type TokenUsage struct {
       Provider       string
       Model          string
       InputTokens    int
       OutputTokens   int
       CachedTokens   int
       TotalCost      float64
       Currency       string
   }
   ```

---

## Provider-Specific Implementation Notes

### Claude Adapter Changes
```go
// Current (adapter.go:41)
args := []string{"-p", "--output-format", "text"}

// Change to:
args := []string{"-p", "--output-format", "json"}
```

**Token Extraction**: Parse `total_cost_usd` and `usage` fields from JSON response

### Cursor Adapter Changes
```go
// Current (adapter.go:69)
args := []string{"-p", fullPrompt, "--output-format", "text"}

// Change to:
args := []string{"-p", fullPrompt, "--output-format", "json"}
```

**Token Extraction**:
- Check if response includes usage data
- May need fallback to session file parsing
- Format likely similar to Claude's structure

### OpenCode Adapter Changes
```go
// Current (adapter.go:99)
args := []string{"run"}
if req.Model != "" {
    args = append(args, "--model", req.Model)
}
args = append(args, req.Prompt)

// Change to:
args := []string{"run", "--format", "json"}
if req.Model != "" {
    args = append(args, "--model", req.Model)
}
args = append(args, req.Prompt)
```

**Token Extraction**: Parse JSON events or use `opencode export` post-execution

---

## External Resources

### Documentation
- [Cursor CLI Output Format Docs](https://cursor.com/docs/cli/reference/output-format)
- [OpenCode CLI Documentation](https://opencode.ai/docs/cli/)

### Community Tools
Several third-party tools track token usage across providers, showing demand for this feature:
- [tokscale](https://github.com/junhoyeo/tokscale) - Multi-provider token tracking
- [opencode-usage-cli](https://libraries.io/npm/opencode-usage-cli) - OpenCode-specific usage tracking
- [cursor-agent-mcp](https://github.com/sailay1996/cursor-agent-mcp) - Cost-efficient analysis

These tools demonstrate that users need normalized token tracking across providers, which Rep could provide.

---

## Summary

**Token usage data is available from all three providers**, but requires provider-specific extraction logic:

1. **Claude**: Fully supported via `--output-format json`, comprehensive data structure
2. **Cursor**: Partially supported, may require session file fallback, evolving feature
3. **OpenCode**: Fully supported via `--format json` or dedicated `stats` command

**The main architectural question is**: Should Rep interpret provider output to expose token data in a normalized way, or should it remain a pure dispatcher and let callers handle provider-specific JSON parsing?

This decision will shape the implementation approach and may require updating INVARIANTS.md to either:
- **Strengthen** the "no interpretation" rule and document that callers must parse JSON themselves, or
- **Relax** the rule to allow token metadata extraction as an exception to the dispatcher model
