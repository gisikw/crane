package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRollup_Empty(t *testing.T) {
	result := Rollup("proj", []UsageRecord{})
	if result.ChargeCode != "proj" {
		t.Errorf("ChargeCode = %q, want %q", result.ChargeCode, "proj")
	}
	if result.Invocations != 0 {
		t.Errorf("Invocations = %d, want 0", result.Invocations)
	}
	if result.TotalCostUSD != 0 {
		t.Errorf("TotalCostUSD = %f, want 0", result.TotalCostUSD)
	}
	if len(result.Breakdown) != 0 {
		t.Errorf("Breakdown len = %d, want 0", len(result.Breakdown))
	}
}

func TestRollup_SingleRecord(t *testing.T) {
	in := 100
	out := 50
	cost := 0.001
	records := []UsageRecord{
		{
			Provider:     "claude",
			Model:        "claude-3-5-sonnet",
			InputTokens:  &in,
			OutputTokens: &out,
			CostUSD:      &cost,
		},
	}
	result := Rollup("test", records)
	if result.Invocations != 1 {
		t.Errorf("Invocations = %d, want 1", result.Invocations)
	}
	if result.TotalInputTokens != 100 {
		t.Errorf("TotalInputTokens = %d, want 100", result.TotalInputTokens)
	}
	if result.TotalOutputTokens != 50 {
		t.Errorf("TotalOutputTokens = %d, want 50", result.TotalOutputTokens)
	}
	if result.TotalCostUSD != 0.001 {
		t.Errorf("TotalCostUSD = %f, want 0.001", result.TotalCostUSD)
	}
	if len(result.Breakdown) != 1 {
		t.Fatalf("Breakdown len = %d, want 1", len(result.Breakdown))
	}
	e := result.Breakdown[0]
	if e.Provider != "claude" || e.Model != "claude-3-5-sonnet" {
		t.Errorf("Breakdown[0] = {%s, %s}, want {claude, claude-3-5-sonnet}", e.Provider, e.Model)
	}
	if e.Invocations != 1 || e.InputTokens != 100 || e.OutputTokens != 50 || e.CostUSD != 0.001 {
		t.Errorf("Breakdown[0] fields mismatch: %+v", e)
	}
}

func TestRollup_MultipleProviderModel(t *testing.T) {
	in1, out1, cost1 := 100, 50, 0.001
	in2, out2, cost2 := 200, 80, 0.002
	in3, out3, cost3 := 300, 90, 0.003

	records := []UsageRecord{
		{Provider: "openai", Model: "gpt-4", InputTokens: &in1, OutputTokens: &out1, CostUSD: &cost1},
		{Provider: "claude", Model: "claude-3-5-sonnet", InputTokens: &in2, OutputTokens: &out2, CostUSD: &cost2},
		{Provider: "openai", Model: "gpt-4", InputTokens: &in3, OutputTokens: &out3, CostUSD: &cost3},
	}

	result := Rollup("multi", records)

	if result.Invocations != 3 {
		t.Errorf("Invocations = %d, want 3", result.Invocations)
	}
	if result.TotalInputTokens != 600 {
		t.Errorf("TotalInputTokens = %d, want 600", result.TotalInputTokens)
	}
	if result.TotalOutputTokens != 220 {
		t.Errorf("TotalOutputTokens = %d, want 220", result.TotalOutputTokens)
	}
	if result.TotalCostUSD != 0.006 {
		t.Errorf("TotalCostUSD = %f, want 0.006", result.TotalCostUSD)
	}
	if len(result.Breakdown) != 2 {
		t.Fatalf("Breakdown len = %d, want 2", len(result.Breakdown))
	}

	// Sorted by provider then model: claude < openai
	if result.Breakdown[0].Provider != "claude" {
		t.Errorf("Breakdown[0].Provider = %q, want %q", result.Breakdown[0].Provider, "claude")
	}
	if result.Breakdown[1].Provider != "openai" {
		t.Errorf("Breakdown[1].Provider = %q, want %q", result.Breakdown[1].Provider, "openai")
	}

	openai := result.Breakdown[1]
	if openai.Invocations != 2 || openai.InputTokens != 400 || openai.OutputTokens != 140 {
		t.Errorf("openai entry mismatch: %+v", openai)
	}
}

func TestRollup_NilTokenFields(t *testing.T) {
	records := []UsageRecord{
		{Provider: "claude", Model: "claude-3-5-sonnet"},
	}
	// Should not panic; nil pointer fields treated as zero.
	result := Rollup("niltest", records)
	if result.Invocations != 1 {
		t.Errorf("Invocations = %d, want 1", result.Invocations)
	}
	if result.TotalInputTokens != 0 {
		t.Errorf("TotalInputTokens = %d, want 0", result.TotalInputTokens)
	}
	if result.TotalCostUSD != 0 {
		t.Errorf("TotalCostUSD = %f, want 0", result.TotalCostUSD)
	}
}

func TestReadChargeRecords_FileNotExist(t *testing.T) {
	t.Setenv("XDG_STATE_HOME", t.TempDir())
	records, err := ReadChargeRecords("nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 0 {
		t.Errorf("expected empty slice, got %d records", len(records))
	}
}

func TestReadChargeRecords_ReadsAndParsesLines(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_STATE_HOME", dir)

	stateD := filepath.Join(dir, "rep")
	if err := os.MkdirAll(stateD, 0700); err != nil {
		t.Fatal(err)
	}

	in1, out1 := 100, 50
	cost1 := 0.001
	in2, out2 := 200, 80
	cost2 := 0.002

	r1 := UsageRecord{Provider: "claude", Model: "claude-3-5-sonnet", InputTokens: &in1, OutputTokens: &out1, CostUSD: &cost1}
	r2 := UsageRecord{Provider: "openai", Model: "gpt-4", InputTokens: &in2, OutputTokens: &out2, CostUSD: &cost2}

	if err := AppendChargeRecord("myproject", r1); err != nil {
		t.Fatal(err)
	}
	if err := AppendChargeRecord("myproject", r2); err != nil {
		t.Fatal(err)
	}

	records, err := ReadChargeRecords("myproject")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	if records[0].Provider != "claude" {
		t.Errorf("records[0].Provider = %q, want %q", records[0].Provider, "claude")
	}
	if records[1].Provider != "openai" {
		t.Errorf("records[1].Provider = %q, want %q", records[1].Provider, "openai")
	}
	if records[0].InputTokens == nil || *records[0].InputTokens != 100 {
		t.Errorf("records[0].InputTokens mismatch")
	}
}

func TestReadChargeRecords_SkipsBadLines(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_STATE_HOME", dir)

	stateD := filepath.Join(dir, "rep")
	if err := os.MkdirAll(stateD, 0700); err != nil {
		t.Fatal(err)
	}

	in := 42
	good := UsageRecord{Provider: "claude", Model: "claude-3-5-sonnet", InputTokens: &in}
	if err := AppendChargeRecord("skiptest", good); err != nil {
		t.Fatal(err)
	}

	// Append a bad line directly.
	path := filepath.Join(stateD, "skiptest.jsonl")
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString("not valid json\n"); err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()

	if err := AppendChargeRecord("skiptest", good); err != nil {
		t.Fatal(err)
	}

	records, err := ReadChargeRecords("skiptest")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 2 {
		t.Errorf("expected 2 records (bad line skipped), got %d", len(records))
	}
}
