package main

import (
	"fmt"
	"io"
	"os"
	"strings"
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
	})

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*os.PathError); ok {
			fmt.Fprintf(os.Stderr, "rep: provider %q not found: %v\n", provider, exitErr)
			os.Exit(1)
		}
		// Mirror the agent's exit code
		os.Exit(exitCode(err))
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
