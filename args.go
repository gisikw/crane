package main

import (
	"fmt"
	"strings"
)

// Options holds the parsed command-line options.
type Options struct {
	Prompt          string
	Provider        string // Empty means "use config preference order"
	Model           string
	Dir             string
	SystemPrompt    string
	AllowAll        bool
	AllowedTools    []string
	DisallowedTools []string
	ChargeCode      string // Tag invocation for cost attribution; empty means no logging
}

// ParseArgs parses CLI arguments into Options.
func ParseArgs(args []string, cfg Config) (Options, error) {
	opts := Options{
		AllowAll: true, // Default: skip permissions
	}

	var positional []string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--provider":
			if i+1 >= len(args) {
				return opts, fmt.Errorf("--provider requires a value")
			}
			i++
			opts.Provider = args[i]
		case "--model":
			if i+1 >= len(args) {
				return opts, fmt.Errorf("--model requires a value")
			}
			i++
			opts.Model = args[i]
		case "--dir":
			if i+1 >= len(args) {
				return opts, fmt.Errorf("--dir requires a value")
			}
			i++
			opts.Dir = args[i]
		case "--system-prompt":
			if i+1 >= len(args) {
				return opts, fmt.Errorf("--system-prompt requires a value")
			}
			i++
			opts.SystemPrompt = args[i]
		case "--allowed-tools":
			if i+1 >= len(args) {
				return opts, fmt.Errorf("--allowed-tools requires a value")
			}
			i++
			opts.AllowedTools = parseToolList(args[i])
		case "--disallowed-tools":
			if i+1 >= len(args) {
				return opts, fmt.Errorf("--disallowed-tools requires a value")
			}
			i++
			opts.DisallowedTools = parseToolList(args[i])
		case "--no-permissions":
			opts.AllowAll = true
		case "--with-permissions":
			opts.AllowAll = false
		case "--charge-code":
			if i+1 >= len(args) {
				return opts, fmt.Errorf("--charge-code requires a value")
			}
			i++
			if err := validateChargeCode(args[i]); err != nil {
				return opts, err
			}
			opts.ChargeCode = args[i]
		case "--help", "-h":
			return opts, fmt.Errorf("help requested")
		default:
			if strings.HasPrefix(args[i], "--") {
				return opts, fmt.Errorf("unknown flag: %s", args[i])
			}
			positional = append(positional, args[i])
		}
	}

	opts.Prompt = strings.Join(positional, " ")
	return opts, nil
}

// parseToolList splits a comma or space-separated list of tool names.
func parseToolList(s string) []string {
	var tools []string
	// First split by comma
	parts := strings.Split(s, ",")
	for _, part := range parts {
		// Then split each part by space
		for _, tool := range strings.Fields(part) {
			if tool != "" {
				tools = append(tools, tool)
			}
		}
	}
	return tools
}
