// sqarol is a command-line tool that generates domain typosquatting
// variations for a given domain name. It prints each variation with
// its generation technique, the look-alike domain, the Levenshtein
// edit distance, and a visual deceptiveness score.
//
// Usage:
//
//	sqarol <domain>
//	sqarol -h | --help
package main

import (
	"fmt"
	"os"

	"github.com/symbolsecurity/sqarol"
)

// main parses the command-line arguments and generates typosquatting
// variations for the provided domain.
func main() {
	if len(os.Args) < 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Fprintln(os.Stderr, "usage: sqarol <domain>")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Generate typosquatting domain variations for the given domain.")
		fmt.Fprintln(os.Stderr, "Each variation is printed with its technique, variant, edit distance, and effectiveness score.")
		os.Exit(0)
	}

	domain := os.Args[1]
	variations, err := sqarol.Generate(domain)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	for _, v := range variations {
		fmt.Printf("%-20s %-40s dist=%-3d eff=%.2f\n", v.Kind, v.Variant, v.Distance, v.Effectiveness)
	}
}
