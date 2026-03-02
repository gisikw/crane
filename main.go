package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "rep: config error: %v\n", err)
		os.Exit(1)
	}

	opts, err := ParseArgs(os.Args[1:], cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "rep: %v\n", err)
		fmt.Fprintf(os.Stderr, "Usage: rep [flags] [prompt]\n")
		fmt.Fprintf(os.Stderr, "  --provider <name>    Force a specific provider\n")
		fmt.Fprintf(os.Stderr, "  --model <model>      Override model selection\n")
		fmt.Fprintf(os.Stderr, "  --dir <path>         Working directory for the agent\n")
		fmt.Fprintf(os.Stderr, "  --system-prompt <s>  Append system prompt\n")
		fmt.Fprintf(os.Stderr, "  --no-permissions     Skip permission checks (default: true)\n")
		fmt.Fprintf(os.Stderr, "  --charge-code <code> Tag invocation for cost attribution\n")
		os.Exit(1)
	}

	// Read prompt from stdin if not provided as args
	if opts.Prompt == "" {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "rep: failed to read stdin: %v\n", err)
			os.Exit(1)
		}
		opts.Prompt = strings.TrimSpace(string(data))
	}

	if opts.Prompt == "" {
		fmt.Fprintf(os.Stderr, "rep: no prompt provided\n")
		os.Exit(1)
	}

	// Resolve provider
	provider, err := ResolveProvider(opts.Provider, cfg.Providers)
	if err != nil {
		fmt.Fprintf(os.Stderr, "rep: %v\n", err)
		os.Exit(1)
	}

	// Build and run the command
	adapter := GetAdapter(provider)
	cmd := adapter.Build(InvokeRequest{
		Prompt:          opts.Prompt,
		Model:           firstNonEmpty(opts.Model, cfg.ProviderConfig(provider).Model),
		SystemPrompt:    opts.SystemPrompt,
		Dir:             opts.Dir,
		AllowAll:        opts.AllowAll,
		AllowedTools:    opts.AllowedTools,
		DisallowedTools: opts.DisallowedTools,
		CaptureUsage:    opts.ChargeCode != "",
	})

	cmd.Stderr = os.Stderr

	var runErr error
	if opts.ChargeCode != "" {
		var buf bytes.Buffer
		cmd.Stdout = &buf
		runErr = cmd.Run()

		stdout, rec, _ := parseProviderOutput(provider, buf.Bytes())
		os.Stdout.Write(stdout) //nolint:errcheck

		rec.Timestamp = time.Now().UTC().Format(time.RFC3339)
		rec.Provider = provider
		rec.Model = firstNonEmpty(opts.Model, cfg.ProviderConfig(provider).Model)
		rec.ChargeCode = opts.ChargeCode

		if err := AppendChargeRecord(opts.ChargeCode, rec); err != nil {
			fmt.Fprintf(os.Stderr, "rep: charge record: %v\n", err)
		}
	} else {
		cmd.Stdout = os.Stdout
		runErr = cmd.Run()
	}

	if runErr != nil {
		if exitErr, ok := runErr.(*os.PathError); ok {
			fmt.Fprintf(os.Stderr, "rep: provider %q not found: %v\n", provider, exitErr)
			os.Exit(1)
		}
		// Mirror the agent's exit code
		os.Exit(exitCode(runErr))
	}
}

// parseProviderOutput dispatches to the appropriate output parser based on provider.
func parseProviderOutput(provider string, output []byte) ([]byte, UsageRecord, error) {
	switch provider {
	case "claude":
		return processClaudeOutput(output)
	case "cursor":
		return processCursorOutput(output)
	case "opencode":
		return processOpencodeOutput(output)
	default:
		return output, UsageRecord{}, nil
	}
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
