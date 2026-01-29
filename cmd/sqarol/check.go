package main

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strconv"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/symbolsecurity/sqarol"
)

// runCheck generates variations, picks the top N by effectiveness, and
// checks each one for registration status, owner, A and MX records.
func runCheck(args []string) {
	if len(args) < 1 || args[0] == "-h" || args[0] == "--help" {
		fmt.Fprintln(os.Stderr, "usage: sqarol check <domain> [-n count]")
		os.Exit(0)
	}

	domain := args[0]
	n := 100

	// Parse -n flag from remaining args.
	for i := 1; i < len(args); i++ {
		if args[i] == "-n" && i+1 < len(args) {
			parsed, err := strconv.Atoi(args[i+1])
			if err != nil || parsed < 1 {
				fmt.Fprintln(os.Stderr, "error: -n requires a positive integer")
				os.Exit(1)
			}
			n = parsed
			i++ // skip the value
		}
	}

	variations, err := sqarol.Generate(domain)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	slices.SortFunc(variations, func(a, b sqarol.Variation) int {
		if a.Effectiveness > b.Effectiveness {
			return -1
		}
		if a.Effectiveness < b.Effectiveness {
			return 1
		}
		return 0
	})

	if n > len(variations) {
		n = len(variations)
	}
	top := variations[:n]

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	type result struct {
		variation sqarol.Variation
		status    *sqarol.DomainCheck
		err       error
	}

	fmt.Fprintf(os.Stderr, "Checking %d variations...\n\n", n)

	results := make([]result, len(top))
	var wg sync.WaitGroup
	for i, v := range top {
		wg.Add(1)
		go func(i int, v sqarol.Variation) {
			defer wg.Done()
			status, err := sqarol.Check(ctx, v.Variant)
			results[i] = result{variation: v, status: status, err: err}
		}(i, v)
	}
	wg.Wait()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "#\tVARIANT\tREGISTERED\tOWNER\tA\tMX\tPARKED")
	fmt.Fprintln(w, "-\t-------\t----------\t-----\t-\t--\t------")
	for i, r := range results {
		registered := "no"
		owner := "-"
		hasA := "no"
		hasMX := "no"
		parked := "no"

		if r.err != nil {
			registered = "error"
		} else if r.status != nil {
			if r.status.IsRegistered {
				registered = "yes"
			}
			if r.status.Owner != "" {
				owner = r.status.Owner
			}
			if r.status.HasARecords {
				hasA = "yes"
			}
			if r.status.HasMXRecords {
				hasMX = "yes"
			}
			if r.status.IsParked {
				parked = "yes"
			}
		}

		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\t%s\n",
			i+1,
			r.variation.Variant,
			registered,
			owner,
			hasA,
			hasMX,
			parked,
		)
	}
	w.Flush()

	fmt.Fprintf(os.Stderr, "\nChecked: %d/%d variations\n", n, len(variations))
}
