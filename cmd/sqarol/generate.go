package main

import (
	"fmt"
	"os"
	"slices"
	"text/tabwriter"

	"github.com/symbolsecurity/sqarol"
)

// runGenerate generates all typosquatting variations for the given domain
// and prints them in a table sorted by effectiveness (descending).
func runGenerate(args []string) {
	if len(args) < 1 || args[0] == "-h" || args[0] == "--help" {
		fmt.Fprintln(os.Stderr, "usage: sqarol generate <domain>")
		os.Exit(0)
	}

	domain := args[0]
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

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "KIND\tVARIANT\tDISTANCE\tEFFECTIVENESS")
	fmt.Fprintln(w, "----\t-------\t--------\t-------------")
	for _, v := range variations {
		fmt.Fprintf(w, "%s\t%s\t%d\t%.2f\n", v.Kind, v.Variant, v.Distance, v.Effectiveness)
	}
	w.Flush()

	fmt.Fprintf(os.Stderr, "\nTotal: %d variations\n", len(variations))
}
