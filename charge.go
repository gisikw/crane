package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var validChargeCodeRe = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// UsageRecord holds per-invocation usage data written as a JSONL line.
type UsageRecord struct {
	Timestamp        string   `json:"timestamp"`
	Provider         string   `json:"provider"`
	Model            string   `json:"model"`
	ChargeCode       string   `json:"charge_code"`
	InputTokens      *int     `json:"input_tokens,omitempty"`
	OutputTokens     *int     `json:"output_tokens,omitempty"`
	CacheReadTokens  *int     `json:"cache_read_tokens,omitempty"`
	CacheWriteTokens *int     `json:"cache_write_tokens,omitempty"`
	CostUSD          *float64 `json:"cost_usd,omitempty"`
}

// validateChargeCode returns an error if code contains characters outside [a-zA-Z0-9_-].
func validateChargeCode(code string) error {
	if !validChargeCodeRe.MatchString(code) {
		return fmt.Errorf("--charge-code %q: only letters, digits, dashes, and underscores are allowed", code)
	}
	return nil
}

// stateDir returns the rep state directory ($XDG_STATE_HOME/rep or ~/.local/state/rep).
func stateDir() string {
	if base := os.Getenv("XDG_STATE_HOME"); base != "" {
		return filepath.Join(base, "rep")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "state", "rep")
}

// AppendChargeRecord appends rec as a JSONL line to stateDir()/<chargeCode>.jsonl.
// The file is opened with O_APPEND so no read-modify-write occurs.
func AppendChargeRecord(chargeCode string, rec UsageRecord) error {
	dir := stateDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("create state dir: %w", err)
	}
	path := filepath.Join(dir, chargeCode+".jsonl")
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	data, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(f, "%s\n", data)
	return err
}

// claudeJSONOutput matches Claude's --output-format json schema.
type claudeJSONOutput struct {
	Result    string  `json:"result"`
	TotalCost float64 `json:"total_cost_usd"`
	Usage     struct {
		InputTokens              int `json:"input_tokens"`
		OutputTokens             int `json:"output_tokens"`
		CacheReadInputTokens     int `json:"cache_read_input_tokens"`
		CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
	} `json:"usage"`
}

// processClaudeOutput parses Claude's JSON output, returns the text result as stdout
// and populates a UsageRecord with token/cost fields. Falls back to raw output on
// parse failure.
func processClaudeOutput(output []byte) ([]byte, UsageRecord, error) {
	var out claudeJSONOutput
	if err := json.Unmarshal(output, &out); err != nil {
		return output, UsageRecord{}, nil
	}
	rec := UsageRecord{
		InputTokens:      intPtr(out.Usage.InputTokens),
		OutputTokens:     intPtr(out.Usage.OutputTokens),
		CacheReadTokens:  intPtr(out.Usage.CacheReadInputTokens),
		CacheWriteTokens: intPtr(out.Usage.CacheCreationInputTokens),
	}
	if out.TotalCost != 0 {
		rec.CostUSD = float64Ptr(out.TotalCost)
	}
	return []byte(out.Result), rec, nil
}

// cursorJSONOutput tries common field names from Cursor's --output-format json schema.
type cursorJSONOutput struct {
	Result       string  `json:"result"`
	Text         string  `json:"text"`
	CostUSD      float64 `json:"cost_usd"`
	Usage        struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// processCursorOutput attempts to parse Cursor's JSON output for token data.
// Falls back to raw output + empty record on parse failure.
func processCursorOutput(output []byte) ([]byte, UsageRecord, error) {
	var out cursorJSONOutput
	if err := json.Unmarshal(output, &out); err != nil {
		return output, UsageRecord{}, nil
	}
	text := out.Result
	if text == "" {
		text = out.Text
	}
	stdout := output
	if text != "" {
		stdout = []byte(text)
	}
	rec := UsageRecord{}
	in := out.Usage.InputTokens
	if in == 0 {
		in = out.InputTokens
	}
	outToks := out.Usage.OutputTokens
	if outToks == 0 {
		outToks = out.OutputTokens
	}
	if in > 0 {
		rec.InputTokens = intPtr(in)
	}
	if outToks > 0 {
		rec.OutputTokens = intPtr(outToks)
	}
	if out.CostUSD != 0 {
		rec.CostUSD = float64Ptr(out.CostUSD)
	}
	return stdout, rec, nil
}

// opencodeEvent represents a single event from OpenCode's --format json stream.
type opencodeEvent struct {
	Type         string  `json:"type"`
	Content      string  `json:"content"`
	Text         string  `json:"text"`
	CostUSD      float64 `json:"cost_usd"`
	Usage        struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// processOpencodeOutput parses OpenCode's newline-delimited JSON event stream,
// concatenates text content for stdout, and accumulates token counts. Falls back
// to raw output + empty record if no lines parse successfully.
func processOpencodeOutput(output []byte) ([]byte, UsageRecord, error) {
	lines := strings.Split(strings.TrimRight(string(output), "\n"), "\n")
	var textParts []string
	rec := UsageRecord{}
	parsed := false
	var totalIn, totalOut int

	for _, line := range lines {
		if line == "" {
			continue
		}
		var event opencodeEvent
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			continue
		}
		parsed = true
		if event.Content != "" {
			textParts = append(textParts, event.Content)
		} else if event.Text != "" {
			textParts = append(textParts, event.Text)
		}
		in := event.Usage.InputTokens
		if in == 0 {
			in = event.InputTokens
		}
		out := event.Usage.OutputTokens
		if out == 0 {
			out = event.OutputTokens
		}
		totalIn += in
		totalOut += out
		if event.CostUSD != 0 && rec.CostUSD == nil {
			rec.CostUSD = float64Ptr(event.CostUSD)
		}
	}

	if !parsed {
		return output, UsageRecord{}, nil
	}
	if totalIn > 0 {
		rec.InputTokens = intPtr(totalIn)
	}
	if totalOut > 0 {
		rec.OutputTokens = intPtr(totalOut)
	}
	stdout := output
	if len(textParts) > 0 {
		stdout = []byte(strings.Join(textParts, ""))
	}
	return stdout, rec, nil
}

func intPtr(v int) *int {
	return &v
}

func float64Ptr(v float64) *float64 {
	return &v
}
