// sqarol is a command-line tool for domain typosquatting analysis.
//
// Usage:
//
//	sqarol generate <domain>
//	sqarol check <domain> [-n count]
//	sqarol -h | --help
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		printUsage()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "generate":
		runGenerate(os.Args[2:])
	case "check":
		runCheck(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "error: unknown command %q\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "usage: sqarol <command> [arguments]")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Commands:")
	fmt.Fprintln(os.Stderr, "  generate <domain>            Generate typosquatting domain variations")
	fmt.Fprintln(os.Stderr, "  check <domain> [-n count]    Check status of top N most effective variations")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Options:")
	fmt.Fprintln(os.Stderr, "  -n count   Number of top variations to check (default: 100)")
	fmt.Fprintln(os.Stderr, "  -h         Show this help message")
}
