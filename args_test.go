package main

import "testing"

func TestParseArgs_PromptFromPositional(t *testing.T) {
	opts, err := ParseArgs([]string{"do", "the", "thing"}, Config{})
	if err != nil {
		t.Fatal(err)
	}
	if opts.Prompt != "do the thing" {
		t.Errorf("prompt = %q, want %q", opts.Prompt, "do the thing")
	}
}

func TestParseArgs_Flags(t *testing.T) {
	opts, err := ParseArgs([]string{
		"--provider", "cursor",
		"--model", "gpt-4",
		"--dir", "/tmp/test",
		"--system-prompt", "be helpful",
		"do stuff",
	}, Config{})
	if err != nil {
		t.Fatal(err)
	}
	if opts.Provider != "cursor" {
		t.Errorf("provider = %q", opts.Provider)
	}
	if opts.Model != "gpt-4" {
		t.Errorf("model = %q", opts.Model)
	}
	if opts.Dir != "/tmp/test" {
		t.Errorf("dir = %q", opts.Dir)
	}
	if opts.SystemPrompt != "be helpful" {
		t.Errorf("system-prompt = %q", opts.SystemPrompt)
	}
	if opts.Prompt != "do stuff" {
		t.Errorf("prompt = %q", opts.Prompt)
	}
}

func TestParseArgs_AllowAllDefault(t *testing.T) {
	opts, err := ParseArgs([]string{"hello"}, Config{})
	if err != nil {
		t.Fatal(err)
	}
	if !opts.AllowAll {
		t.Error("expected AllowAll to default to true")
	}
}

func TestParseArgs_WithPermissions(t *testing.T) {
	opts, err := ParseArgs([]string{"--with-permissions", "hello"}, Config{})
	if err != nil {
		t.Fatal(err)
	}
	if opts.AllowAll {
		t.Error("expected AllowAll to be false with --with-permissions")
	}
}

func TestParseArgs_UnknownFlag(t *testing.T) {
	_, err := ParseArgs([]string{"--bogus"}, Config{})
	if err == nil {
		t.Error("expected error for unknown flag")
	}
}

func TestParseArgs_EmptyPrompt(t *testing.T) {
	opts, err := ParseArgs([]string{}, Config{})
	if err != nil {
		t.Fatal(err)
	}
	if opts.Prompt != "" {
		t.Errorf("expected empty prompt, got %q", opts.Prompt)
	}
}

func TestParseArgs_AllowedTools_Single(t *testing.T) {
	opts, err := ParseArgs([]string{"--allowed-tools", "Bash", "hello"}, Config{})
	if err != nil {
		t.Fatal(err)
	}
	if len(opts.AllowedTools) != 1 {
		t.Errorf("expected 1 allowed tool, got %d", len(opts.AllowedTools))
	}
	if opts.AllowedTools[0] != "Bash" {
		t.Errorf("allowed tool = %q, want Bash", opts.AllowedTools[0])
	}
}

func TestParseArgs_AllowedTools_Multiple_CommaSeparated(t *testing.T) {
	opts, err := ParseArgs([]string{"--allowed-tools", "Bash,Edit,Read", "hello"}, Config{})
	if err != nil {
		t.Fatal(err)
	}
	if len(opts.AllowedTools) != 3 {
		t.Errorf("expected 3 allowed tools, got %d", len(opts.AllowedTools))
	}
	expected := []string{"Bash", "Edit", "Read"}
	for i, tool := range expected {
		if opts.AllowedTools[i] != tool {
			t.Errorf("allowed tool[%d] = %q, want %q", i, opts.AllowedTools[i], tool)
		}
	}
}

func TestParseArgs_AllowedTools_Multiple_SpaceSeparated(t *testing.T) {
	opts, err := ParseArgs([]string{"--allowed-tools", "Bash Edit Read", "hello"}, Config{})
	if err != nil {
		t.Fatal(err)
	}
	if len(opts.AllowedTools) != 3 {
		t.Errorf("expected 3 allowed tools, got %d", len(opts.AllowedTools))
	}
	expected := []string{"Bash", "Edit", "Read"}
	for i, tool := range expected {
		if opts.AllowedTools[i] != tool {
			t.Errorf("allowed tool[%d] = %q, want %q", i, opts.AllowedTools[i], tool)
		}
	}
}

func TestParseArgs_DisallowedTools_Single(t *testing.T) {
	opts, err := ParseArgs([]string{"--disallowed-tools", "Bash", "hello"}, Config{})
	if err != nil {
		t.Fatal(err)
	}
	if len(opts.DisallowedTools) != 1 {
		t.Errorf("expected 1 disallowed tool, got %d", len(opts.DisallowedTools))
	}
	if opts.DisallowedTools[0] != "Bash" {
		t.Errorf("disallowed tool = %q, want Bash", opts.DisallowedTools[0])
	}
}

func TestParseArgs_ChargeCode(t *testing.T) {
	opts, err := ParseArgs([]string{"--charge-code", "exocortex", "do stuff"}, Config{})
	if err != nil {
		t.Fatal(err)
	}
	if opts.ChargeCode != "exocortex" {
		t.Errorf("charge code = %q, want exocortex", opts.ChargeCode)
	}
	if opts.Prompt != "do stuff" {
		t.Errorf("prompt = %q, want do stuff", opts.Prompt)
	}
}

func TestParseArgs_ChargeCodeMissingValue(t *testing.T) {
	_, err := ParseArgs([]string{"--charge-code"}, Config{})
	if err == nil {
		t.Error("expected error when --charge-code has no value")
	}
}

func TestParseArgs_ChargeCodeInvalidChars(t *testing.T) {
	_, err := ParseArgs([]string{"--charge-code", "foo/bar"}, Config{})
	if err == nil {
		t.Error("expected error for charge code with invalid characters")
	}
}

func TestParseArgs_DisallowedTools_Multiple(t *testing.T) {
	opts, err := ParseArgs([]string{"--disallowed-tools", "Bash,Edit", "hello"}, Config{})
	if err != nil {
		t.Fatal(err)
	}
	if len(opts.DisallowedTools) != 2 {
		t.Errorf("expected 2 disallowed tools, got %d", len(opts.DisallowedTools))
	}
	expected := []string{"Bash", "Edit"}
	for i, tool := range expected {
		if opts.DisallowedTools[i] != tool {
			t.Errorf("disallowed tool[%d] = %q, want %q", i, opts.DisallowedTools[i], tool)
		}
	}
}
