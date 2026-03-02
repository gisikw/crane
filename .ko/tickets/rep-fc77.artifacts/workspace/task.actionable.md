rep dev shell
  just test    — run tests
  just build   — build binary
The plan is well-structured with clear goals, identified files, and ordered tasks with verification steps. However, it has **3 open questions that require human input** before implementation can proceed:

```json
{
  "disposition": "needs_input",
  "plan_questions": [
    {
      "id": "q1",
      "question": "When `--charge-code` is set with the Claude provider, should we extract the `result` text field from Claude's JSON output and write that to stdout, or pass the full JSON through to the caller?",
      "context": "The ticket emphasizes 'output goes to stdout unmodified', but extracting text maintains consistency with the current caller experience. Passing JSON would give full fidelity but change output format unexpectedly.",
      "options": [
        {
          "label": "Extract text result (Recommended)",
          "value": "extract_text",
          "description": "Parse Claude JSON, write the `result` field to stdout — maintains consistent caller experience"
        },
        {
          "label": "Pass full JSON through",
          "value": "pass_json",
          "description": "Write Claude's full JSON output to stdout — preserves all fidelity but changes output format"
        }
      ]
    },
    {
      "id": "q2",
      "question": "For Cursor and OpenCode providers, is it acceptable to log sparse usage records (provider, model, timestamp only, no token counts), or should we attempt to parse their JSON output to extract token data?",
      "context": "The ticket notes these providers are 'spottier' on token data. The plan recommends sparse records as simpler and more honest about data availability.",
      "options": [
        {
          "label": "Sparse records only (Recommended)",
          "value": "sparse",
          "description": "Log only provider, model, and timestamp for these providers — matches ticket's assessment they are spottier"
        },
        {
          "label": "Attempt JSON parsing",
          "value": "parse_json",
          "description": "Try to extract token counts from Cursor/OpenCode JSON output for more complete records"
        }
      ]
    },
    {
      "id": "q3",
      "question": "Should we validate and sanitize the `--charge-code` string before using it in the file path to prevent path traversal attacks?",
      "context": "The ticket says callers decide granularity, but doesn't address safety. The plan recommends rejecting unsafe characters like `/`, `..`, and null bytes.",
      "options": [
        {
          "label": "Validate and reject unsafe chars (Recommended)",
          "value": "validate",
          "description": "Reject charge codes containing `/`, `..`, null bytes, and other path-unsafe characters"
        },
        {
          "label": "Trust caller input",
          "value": "no_validate",
          "description": "Accept any charge-code string as-is and let callers handle safety"
        }
      ]
    }
  ]
}
```
