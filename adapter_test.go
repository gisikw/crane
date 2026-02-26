package main

import "testing"

func TestClaudeAdapter_Build(t *testing.T) {
	a := ClaudeAdapter{}
	cmd := a.Build(InvokeRequest{
		Prompt:       "hello",
		Model:        "opus",
		SystemPrompt: "be nice",
		Dir:          "/tmp",
		AllowAll:     true,
	})

	if cmd.Path == "" {
		t.Fatal("expected non-empty path")
	}
	args := cmd.Args[1:] // skip binary name
	assertContains(t, args, "--dangerously-skip-permissions")
	assertContains(t, args, "--model")
	assertContains(t, args, "opus")
	assertContains(t, args, "--append-system-prompt")
	assertContains(t, args, "be nice")
	if cmd.Dir != "/tmp" {
		t.Errorf("dir = %q, want /tmp", cmd.Dir)
	}
}

func TestClaudeAdapter_NoPermissions(t *testing.T) {
	a := ClaudeAdapter{}
	cmd := a.Build(InvokeRequest{
		Prompt:   "hello",
		AllowAll: false,
	})
	args := cmd.Args[1:]
	for _, arg := range args {
		if arg == "--dangerously-skip-permissions" {
			t.Error("should not include --dangerously-skip-permissions when AllowAll is false")
		}
	}
}

func TestCursorAdapter_InlinesSystemPrompt(t *testing.T) {
	a := CursorAdapter{}
	cmd := a.Build(InvokeRequest{
		Prompt:       "do stuff",
		SystemPrompt: "context here",
	})
	args := cmd.Args[1:]
	// System prompt should be inlined into -p value, not a separate flag
	for _, arg := range args {
		if arg == "--append-system-prompt" {
			t.Error("cursor should not use --append-system-prompt")
		}
	}
	// The -p value should contain both system prompt and prompt
	if len(args) > 1 && args[0] == "-p" {
		if args[1] != "context here\n\ndo stuff" {
			t.Errorf("prompt = %q, want system prompt inlined", args[1])
		}
	}
}

func TestOpencodeAdapter_Build(t *testing.T) {
	a := OpencodeAdapter{}
	cmd := a.Build(InvokeRequest{
		Prompt: "hello world",
		Model:  "ollama/qwen2.5-coder:32b",
		Dir:    "/tmp",
	})

	args := cmd.Args[1:]
	if args[0] != "run" {
		t.Errorf("first arg = %q, want 'run'", args[0])
	}
	assertContains(t, args, "--model")
	assertContains(t, args, "ollama/qwen2.5-coder:32b")
	// Prompt should be the last arg
	if args[len(args)-1] != "hello world" {
		t.Errorf("last arg = %q, want prompt", args[len(args)-1])
	}
	if cmd.Dir != "/tmp" {
		t.Errorf("dir = %q, want /tmp", cmd.Dir)
	}
}

func TestOpencodeAdapter_NoModel(t *testing.T) {
	a := OpencodeAdapter{}
	cmd := a.Build(InvokeRequest{
		Prompt: "hello",
	})

	args := cmd.Args[1:]
	for _, arg := range args {
		if arg == "--model" {
			t.Error("should not include --model when none specified")
		}
	}
	// Should be: run hello
	if len(args) != 2 || args[0] != "run" || args[1] != "hello" {
		t.Errorf("args = %v, want [run hello]", args)
	}
}

func TestGetAdapter_Defaults(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"claude", "main.ClaudeAdapter"},
		{"cursor", "main.CursorAdapter"},
		{"opencode", "main.OpencodeAdapter"},
		{"unknown", "main.GenericAdapter"},
	}
	for _, tt := range tests {
		a := GetAdapter(tt.name)
		got := typeName(a)
		if got != tt.expected {
			t.Errorf("GetAdapter(%q) = %s, want %s", tt.name, got, tt.expected)
		}
	}
}

func TestClaudeAdapter_AllowedTools(t *testing.T) {
	a := ClaudeAdapter{}
	cmd := a.Build(InvokeRequest{
		Prompt:       "hello",
		AllowAll:     false,
		AllowedTools: []string{"Bash", "Edit", "Read"},
	})
	args := cmd.Args[1:]
	assertContains(t, args, "--allowed-tools")
	assertContains(t, args, "Bash")
	assertContains(t, args, "Edit")
	assertContains(t, args, "Read")
	// Should not include skip-permissions
	for _, arg := range args {
		if arg == "--dangerously-skip-permissions" {
			t.Error("should not include --dangerously-skip-permissions when AllowAll is false")
		}
	}
}

func TestClaudeAdapter_DisallowedTools(t *testing.T) {
	a := ClaudeAdapter{}
	cmd := a.Build(InvokeRequest{
		Prompt:          "hello",
		AllowAll:        false,
		DisallowedTools: []string{"Bash", "Write"},
	})
	args := cmd.Args[1:]
	assertContains(t, args, "--disallowed-tools")
	assertContains(t, args, "Bash")
	assertContains(t, args, "Write")
}

func TestClaudeAdapter_ToolsIgnoredWithAllowAll(t *testing.T) {
	a := ClaudeAdapter{}
	cmd := a.Build(InvokeRequest{
		Prompt:       "hello",
		AllowAll:     true,
		AllowedTools: []string{"Bash", "Edit"},
	})
	args := cmd.Args[1:]
	// Should include skip-permissions
	assertContains(t, args, "--dangerously-skip-permissions")
	// Should not include tool flags
	for _, arg := range args {
		if arg == "--allowed-tools" || arg == "--disallowed-tools" {
			t.Errorf("should not include tool flags when AllowAll is true, got %q", arg)
		}
	}
}

func TestExitCode(t *testing.T) {
	code := exitCode(nil)
	if code != 1 {
		t.Errorf("exitCode(nil) = %d, want 1", code)
	}
}

func assertContains(t *testing.T, args []string, val string) {
	t.Helper()
	for _, a := range args {
		if a == val {
			return
		}
	}
	t.Errorf("args %v missing %q", args, val)
}

func typeName(v interface{}) string {
	switch v.(type) {
	case ClaudeAdapter:
		return "main.ClaudeAdapter"
	case CursorAdapter:
		return "main.CursorAdapter"
	case OpencodeAdapter:
		return "main.OpencodeAdapter"
	case GenericAdapter:
		return "main.GenericAdapter"
	default:
		return "unknown"
	}
}
