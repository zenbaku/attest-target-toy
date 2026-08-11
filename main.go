// Command toy applies a per-key budget to a list of items and logs the outcome.
//
// It exists to be a target for the attest runner: it builds, runs, emits
// structured log lines, and has behaviour that can be broken on demand.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"toy/budget"
)

func main() {
	input := flag.String("input", "", "path to a JSON array of {key, cost} items")
	perKey := flag.Int("budget", 10, "per-key budget")
	flag.Parse()

	if *input == "" {
		fmt.Fprintln(os.Stderr, "[ERROR] toy: --input is required")
		os.Exit(2)
	}

	raw, err := os.ReadFile(*input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] toy: cannot read input: %v\n", err)
		os.Exit(1)
	}

	var items []budget.Item
	if err := json.Unmarshal(raw, &items); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] toy: cannot parse input: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("[INFO] toy: run_start items=%d budget=%d\n", len(items), *perKey)

	results := budget.Apply(items, *perKey)
	for _, r := range results {
		fmt.Printf("[INFO] toy: item_done key=%s cost=%d status=%s spent=%d\n",
			r.Item.Key, r.Item.Cost, r.Status, r.Spent)
	}

	s := budget.Summarize(results)
	fmt.Printf("[INFO] toy: run_done keys=%d processed=%d capped=%d invalid=%d\n",
		s.Keys, s.Processed, s.Capped, s.Invalid)
}
