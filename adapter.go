package main

import (
	"os/exec"
	"strings"
	"syscall"
)

// InvokeRequest is the provider-agnostic input for an agent invocation.
type InvokeRequest struct {
	Prompt          string
	Model           string
	SystemPrompt    string
	Dir             string
	AllowAll        bool
	AllowedTools    []string
	DisallowedTools []string
	CaptureUsage    bool // When true, adapters switch to JSON output mode for usage extraction
}

// Adapter builds an exec.Cmd for a specific provider.
type Adapter interface {
	Build(req InvokeRequest) *exec.Cmd
}

// GetAdapter returns the adapter for a provider name.
func GetAdapter(provider string) Adapter {
	switch provider {
	case "claude":
		return ClaudeAdapter{}
	case "cursor":
		return CursorAdapter{}
	case "opencode":
		return OpencodeAdapter{}
	default:
		return GenericAdapter{Binary: provider}
	}
}

// ClaudeAdapter builds commands for the Claude Code CLI.
type ClaudeAdapter struct{}

func (a ClaudeAdapter) Build(req InvokeRequest) *exec.Cmd {
	outputFormat := "text"
	if req.CaptureUsage {
		outputFormat = "json"
	}
	args := []string{"-p", "--output-format", outputFormat}
	if req.AllowAll {
		args = append(args, "--dangerously-skip-permissions")
	} else {
		// Tool filters are only compatible with permission mode (not skip-permissions)
		if len(req.AllowedTools) > 0 {
			args = append(args, "--allowed-tools")
			args = append(args, req.AllowedTools...)
		}
		if len(req.DisallowedTools) > 0 {
			args = append(args, "--disallowed-tools")
			args = append(args, req.DisallowedTools...)
		}
	}
	if req.Model != "" {
		args = append(args, "--model", req.Model)
	}
	if req.SystemPrompt != "" {
		args = append(args, "--append-system-prompt", req.SystemPrompt)
	}

	cmd := exec.Command("claude", args...)
	cmd.Stdin = strings.NewReader(req.Prompt)
	if req.Dir != "" {
		cmd.Dir = req.Dir
	}
	return cmd
}

// CursorAdapter builds commands for Cursor's agent CLI.
type CursorAdapter struct{}

func (a CursorAdapter) Build(req InvokeRequest) *exec.Cmd {
	fullPrompt := req.Prompt
	if req.SystemPrompt != "" {
		fullPrompt = req.SystemPrompt + "\n\n" + req.Prompt
	}

	outputFormat := "text"
	if req.CaptureUsage {
		outputFormat = "json"
	}
	args := []string{"-p", fullPrompt, "--output-format", outputFormat}
	if req.AllowAll {
		args = append(args, "--force")
	}
	if req.Model != "" {
		args = append(args, "--model", req.Model)
	}

	bin := resolveCursorBin()
	cmd := exec.Command(bin, args...)
	if req.Dir != "" {
		cmd.Dir = req.Dir
	}
	return cmd
}

func resolveCursorBin() string {
	if path, err := exec.LookPath("cursor-agent"); err == nil {
		return path
	}
	if path, err := exec.LookPath("agent"); err == nil {
		return path
	}
	return "cursor-agent"
}

// OpencodeAdapter builds commands for the OpenCode CLI.
type OpencodeAdapter struct{}

func (a OpencodeAdapter) Build(req InvokeRequest) *exec.Cmd {
	args := []string{"run"}
	if req.CaptureUsage {
		args = append(args, "--format", "json")
	}
	if req.Model != "" {
		args = append(args, "--model", req.Model)
	}
	args = append(args, req.Prompt)

	cmd := exec.Command("opencode", args...)
	if req.Dir != "" {
		cmd.Dir = req.Dir
	}
	return cmd
}

// GenericAdapter treats the provider name as a binary and passes claude-style flags.
type GenericAdapter struct {
	Binary string
}

func (a GenericAdapter) Build(req InvokeRequest) *exec.Cmd {
	args := []string{"-p", "--output-format", "text"}
	if req.Model != "" {
		args = append(args, "--model", req.Model)
	}
	if req.SystemPrompt != "" {
		args = append(args, "--append-system-prompt", req.SystemPrompt)
	}

	cmd := exec.Command(a.Binary, args...)
	cmd.Stdin = strings.NewReader(req.Prompt)
	if req.Dir != "" {
		cmd.Dir = req.Dir
	}
	return cmd
}

// exitCode extracts the exit code from an exec error.
func exitCode(err error) int {
	if exitErr, ok := err.(*exec.ExitError); ok {
		if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		}
	}
	return 1
}
