package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStateDir_Default(t *testing.T) {
	t.Setenv("XDG_STATE_HOME", "")
	dir := stateDir()
	if !strings.HasSuffix(dir, "/.local/state/rep") {
		t.Errorf("stateDir() = %q, want suffix /.local/state/rep", dir)
	}
}

func TestStateDir_Custom(t *testing.T) {
	t.Setenv("XDG_STATE_HOME", "/tmp/foo")
	dir := stateDir()
	if dir != "/tmp/foo/rep" {
		t.Errorf("stateDir() = %q, want /tmp/foo/rep", dir)
	}
}

func TestAppendChargeRecord_CreatesAndAppends(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_STATE_HOME", tmp)

	rec1 := UsageRecord{
		Timestamp:  "2026-03-02T00:00:00Z",
		Provider:   "claude",
		Model:      "opus",
		ChargeCode: "test",
	}
	in := 10
	rec1.InputTokens = &in

	if err := AppendChargeRecord("test", rec1); err != nil {
		t.Fatal(err)
	}

	rec2 := UsageRecord{
		Timestamp:  "2026-03-02T00:01:00Z",
		Provider:   "claude",
		Model:      "opus",
		ChargeCode: "test",
	}
	out := 20
	rec2.OutputTokens = &out

	if err := AppendChargeRecord("test", rec2); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(tmp, "rep", "test.jsonl"))
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}

	var parsed1 UsageRecord
	if err := json.Unmarshal([]byte(lines[0]), &parsed1); err != nil {
		t.Fatalf("line 1 parse error: %v", err)
	}
	if parsed1.Provider != "claude" {
		t.Errorf("line 1 provider = %q, want claude", parsed1.Provider)
	}
	if parsed1.InputTokens == nil || *parsed1.InputTokens != 10 {
		t.Errorf("line 1 input_tokens = %v, want 10", parsed1.InputTokens)
	}
	if parsed1.OutputTokens != nil {
		t.Error("line 1 output_tokens should be absent")
	}

	var parsed2 UsageRecord
	if err := json.Unmarshal([]byte(lines[1]), &parsed2); err != nil {
		t.Fatalf("line 2 parse error: %v", err)
	}
	if parsed2.OutputTokens == nil || *parsed2.OutputTokens != 20 {
		t.Errorf("line 2 output_tokens = %v, want 20", parsed2.OutputTokens)
	}
}

func TestProcessClaudeOutput_Full(t *testing.T) {
	input := []byte(`{
		"result": "The answer is 42.",
		"total_cost_usd": 0.001234,
		"usage": {
			"input_tokens": 100,
			"output_tokens": 50,
			"cache_read_input_tokens": 10,
			"cache_creation_input_tokens": 5
		}
	}`)

	stdout, rec, err := processClaudeOutput(input)
	if err != nil {
		t.Fatal(err)
	}
	if string(stdout) != "The answer is 42." {
		t.Errorf("stdout = %q, want %q", stdout, "The answer is 42.")
	}
	if rec.InputTokens == nil || *rec.InputTokens != 100 {
		t.Errorf("input_tokens = %v, want 100", rec.InputTokens)
	}
	if rec.OutputTokens == nil || *rec.OutputTokens != 50 {
		t.Errorf("output_tokens = %v, want 50", rec.OutputTokens)
	}
	if rec.CacheReadTokens == nil || *rec.CacheReadTokens != 10 {
		t.Errorf("cache_read_tokens = %v, want 10", rec.CacheReadTokens)
	}
	if rec.CacheWriteTokens == nil || *rec.CacheWriteTokens != 5 {
		t.Errorf("cache_write_tokens = %v, want 5", rec.CacheWriteTokens)
	}
	if rec.CostUSD == nil || *rec.CostUSD != 0.001234 {
		t.Errorf("cost_usd = %v, want 0.001234", rec.CostUSD)
	}
}

func TestProcessClaudeOutput_Fallback(t *testing.T) {
	raw := []byte("not json at all")
	stdout, rec, err := processClaudeOutput(raw)
	if err != nil {
		t.Fatal(err)
	}
	if string(stdout) != "not json at all" {
		t.Errorf("stdout = %q, want raw input", stdout)
	}
	if rec.CostUSD != nil {
		t.Error("CostUSD should be nil on fallback")
	}
	if rec.InputTokens != nil {
		t.Error("InputTokens should be nil on fallback")
	}
}

func TestProcessCursorOutput_Full(t *testing.T) {
	input := []byte(`{"result": "Done!", "usage": {"input_tokens": 80, "output_tokens": 30}}`)
	stdout, rec, err := processCursorOutput(input)
	if err != nil {
		t.Fatal(err)
	}
	if string(stdout) != "Done!" {
		t.Errorf("stdout = %q, want Done!", stdout)
	}
	if rec.InputTokens == nil || *rec.InputTokens != 80 {
		t.Errorf("input_tokens = %v, want 80", rec.InputTokens)
	}
	if rec.OutputTokens == nil || *rec.OutputTokens != 30 {
		t.Errorf("output_tokens = %v, want 30", rec.OutputTokens)
	}
}

func TestProcessCursorOutput_FlatTokens(t *testing.T) {
	input := []byte(`{"text": "Done!", "input_tokens": 40, "output_tokens": 15}`)
	stdout, rec, err := processCursorOutput(input)
	if err != nil {
		t.Fatal(err)
	}
	if string(stdout) != "Done!" {
		t.Errorf("stdout = %q, want Done!", stdout)
	}
	if rec.InputTokens == nil || *rec.InputTokens != 40 {
		t.Errorf("input_tokens = %v, want 40", rec.InputTokens)
	}
}

func TestProcessCursorOutput_Fallback(t *testing.T) {
	raw := []byte("plain text output")
	stdout, rec, err := processCursorOutput(raw)
	if err != nil {
		t.Fatal(err)
	}
	if string(stdout) != "plain text output" {
		t.Errorf("stdout = %q, want raw input", stdout)
	}
	if rec.InputTokens != nil {
		t.Error("InputTokens should be nil on fallback")
	}
}

func TestProcessOpencodeOutput_Events(t *testing.T) {
	input := []byte(`{"type":"text","content":"Hello "}
{"type":"text","content":"world!"}
{"type":"done","usage":{"input_tokens":60,"output_tokens":20}}
`)
	stdout, rec, err := processOpencodeOutput(input)
	if err != nil {
		t.Fatal(err)
	}
	if string(stdout) != "Hello world!" {
		t.Errorf("stdout = %q, want %q", stdout, "Hello world!")
	}
	if rec.InputTokens == nil || *rec.InputTokens != 60 {
		t.Errorf("input_tokens = %v, want 60", rec.InputTokens)
	}
	if rec.OutputTokens == nil || *rec.OutputTokens != 20 {
		t.Errorf("output_tokens = %v, want 20", rec.OutputTokens)
	}
}

func TestProcessOpencodeOutput_Fallback(t *testing.T) {
	raw := []byte("not\x00valid\x00json\x00events")
	stdout, rec, err := processOpencodeOutput(raw)
	if err != nil {
		t.Fatal(err)
	}
	if string(stdout) != string(raw) {
		t.Errorf("stdout = %q, want raw input", stdout)
	}
	if rec.InputTokens != nil {
		t.Error("InputTokens should be nil on fallback")
	}
}

func TestValidateChargeCode_Valid(t *testing.T) {
	cases := []string{"exocortex", "dept-2026", "a_b", "ABC123", "x-y_z"}
	for _, c := range cases {
		if err := validateChargeCode(c); err != nil {
			t.Errorf("validateChargeCode(%q) = %v, want nil", c, err)
		}
	}
}

func TestValidateChargeCode_Invalid(t *testing.T) {
	cases := []string{"foo/bar", "../etc", "foo bar", "a\x00b", "dept:2026"}
	for _, c := range cases {
		if err := validateChargeCode(c); err == nil {
			t.Errorf("validateChargeCode(%q) = nil, want error", c)
		}
	}
}
