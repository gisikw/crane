package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RollupEntry holds aggregated usage for a single (provider, model) pair.
type RollupEntry struct {
	Provider         string  `json:"provider"`
	Model            string  `json:"model"`
	Invocations      int     `json:"invocations"`
	InputTokens      int     `json:"input_tokens"`
	OutputTokens     int     `json:"output_tokens"`
	CacheReadTokens  int     `json:"cache_read_tokens"`
	CacheWriteTokens int     `json:"cache_write_tokens"`
	CostUSD          float64 `json:"cost_usd"`
}

// RollupResult holds aggregated totals and a per-(provider, model) breakdown.
type RollupResult struct {
	ChargeCode             string        `json:"charge_code"`
	Invocations            int           `json:"invocations"`
	TotalInputTokens       int           `json:"total_input_tokens"`
	TotalOutputTokens      int           `json:"total_output_tokens"`
	TotalCacheReadTokens   int           `json:"total_cache_read_tokens"`
	TotalCacheWriteTokens  int           `json:"total_cache_write_tokens"`
	TotalCostUSD           float64       `json:"total_cost_usd"`
	Breakdown              []RollupEntry `json:"breakdown"`
}

// Rollup aggregates records by (provider, model), returning sorted breakdown and grand totals.
func Rollup(chargeCode string, records []UsageRecord) RollupResult {
	type key struct{ provider, model string }
	groups := make(map[key]*RollupEntry)
	var order []key

	result := RollupResult{ChargeCode: chargeCode}

	for _, r := range records {
		k := key{r.Provider, r.Model}
		if _, exists := groups[k]; !exists {
			groups[k] = &RollupEntry{Provider: r.Provider, Model: r.Model}
			order = append(order, k)
		}
		e := groups[k]
		e.Invocations++
		result.Invocations++

		if r.InputTokens != nil {
			e.InputTokens += *r.InputTokens
			result.TotalInputTokens += *r.InputTokens
		}
		if r.OutputTokens != nil {
			e.OutputTokens += *r.OutputTokens
			result.TotalOutputTokens += *r.OutputTokens
		}
		if r.CacheReadTokens != nil {
			e.CacheReadTokens += *r.CacheReadTokens
			result.TotalCacheReadTokens += *r.CacheReadTokens
		}
		if r.CacheWriteTokens != nil {
			e.CacheWriteTokens += *r.CacheWriteTokens
			result.TotalCacheWriteTokens += *r.CacheWriteTokens
		}
		if r.CostUSD != nil {
			e.CostUSD += *r.CostUSD
			result.TotalCostUSD += *r.CostUSD
		}
	}

	// Sort breakdown by provider then model for deterministic output.
	sort.Slice(order, func(i, j int) bool {
		if order[i].provider != order[j].provider {
			return order[i].provider < order[j].provider
		}
		return order[i].model < order[j].model
	})

	result.Breakdown = make([]RollupEntry, 0, len(order))
	for _, k := range order {
		result.Breakdown = append(result.Breakdown, *groups[k])
	}

	return result
}

// ReadChargeRecords reads all UsageRecord lines from stateDir()/<chargeCode>.jsonl.
// Returns an empty slice (not an error) if the file does not exist.
// Blank or unparseable lines are silently skipped.
func ReadChargeRecords(chargeCode string) ([]UsageRecord, error) {
	path := filepath.Join(stateDir(), chargeCode+".jsonl")
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []UsageRecord{}, nil
		}
		return nil, err
	}
	defer f.Close()

	var records []UsageRecord
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		var rec UsageRecord
		if err := json.Unmarshal([]byte(line), &rec); err != nil {
			continue
		}
		records = append(records, rec)
	}
	return records, scanner.Err()
}

// FormatRollup produces a human-readable summary table for the rollup result.
func FormatRollup(r RollupResult) string {
	var b strings.Builder

	fmt.Fprintf(&b, "Charge code: %s\n", r.ChargeCode)
	fmt.Fprintf(&b, "Invocations: %d\n", r.Invocations)
	fmt.Fprintf(&b, "Total cost:  $%.6f\n", r.TotalCostUSD)
	fmt.Fprintf(&b, "\n")

	if len(r.Breakdown) == 0 {
		fmt.Fprintf(&b, "No usage records found.\n")
		return b.String()
	}

	hasCache := r.TotalCacheReadTokens > 0 || r.TotalCacheWriteTokens > 0

	if hasCache {
		fmt.Fprintf(&b, "%-12s  %-30s  %6s  %8s  %8s  %9s  %10s  %10s\n",
			"PROVIDER", "MODEL", "RUNS", "INPUT", "OUTPUT", "CACHE-READ", "CACHE-WRITE", "COST")
	} else {
		fmt.Fprintf(&b, "%-12s  %-30s  %6s  %8s  %8s  %10s\n",
			"PROVIDER", "MODEL", "RUNS", "INPUT", "OUTPUT", "COST")
	}

	for _, e := range r.Breakdown {
		if hasCache {
			fmt.Fprintf(&b, "%-12s  %-30s  %6d  %8d  %8d  %9d  %11d  %10.6f\n",
				e.Provider, e.Model, e.Invocations,
				e.InputTokens, e.OutputTokens,
				e.CacheReadTokens, e.CacheWriteTokens,
				e.CostUSD)
		} else {
			fmt.Fprintf(&b, "%-12s  %-30s  %6d  %8d  %8d  %10.6f\n",
				e.Provider, e.Model, e.Invocations,
				e.InputTokens, e.OutputTokens,
				e.CostUSD)
		}
	}

	return b.String()
}

// cmdUsage implements the "rep usage <charge-code> [--json]" subcommand.
func cmdUsage(args []string) {
	var chargeCode string
	jsonOutput := false

	for _, arg := range args {
		switch {
		case arg == "--json":
			jsonOutput = true
		case strings.HasPrefix(arg, "-"):
			fmt.Fprintf(os.Stderr, "rep usage: unknown flag %q\n", arg)
			fmt.Fprintf(os.Stderr, "Usage: rep usage <charge-code> [--json]\n")
			os.Exit(1)
		default:
			if chargeCode != "" {
				fmt.Fprintf(os.Stderr, "rep usage: unexpected argument %q\n", arg)
				fmt.Fprintf(os.Stderr, "Usage: rep usage <charge-code> [--json]\n")
				os.Exit(1)
			}
			chargeCode = arg
		}
	}

	if chargeCode == "" {
		fmt.Fprintf(os.Stderr, "rep usage: charge-code required\n")
		fmt.Fprintf(os.Stderr, "Usage: rep usage <charge-code> [--json]\n")
		os.Exit(1)
	}

	if err := validateChargeCode(chargeCode); err != nil {
		fmt.Fprintf(os.Stderr, "rep usage: %v\n", err)
		os.Exit(1)
	}

	records, err := ReadChargeRecords(chargeCode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "rep usage: %v\n", err)
		os.Exit(1)
	}

	result := Rollup(chargeCode, records)

	if jsonOutput {
		data, err := json.Marshal(result)
		if err != nil {
			fmt.Fprintf(os.Stderr, "rep usage: marshal: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(data))
		return
	}

	fmt.Print(FormatRollup(result))
}
