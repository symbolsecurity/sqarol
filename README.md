# sqarol

A Go library and CLI tool designed to help cybersecurity analysts manually monitor lookalike domains. Given a legitimate domain name, it produces hundreds of plausible look-alike variations that attackers might register for phishing, brand impersonation, or credential harvesting. Each variation is annotated with the generation technique, Levenshtein edit distance, and a visual deceptiveness score.

For automated, continuous domain threat monitoring, check out [Symbol Security's Domain Threat Alerts](https://symbolsecurity.com/domain-threat-alerts).

## Installation

As a library:

```
go get github.com/symbolsecurity/sqarol
```

As a CLI tool:

```
go install github.com/symbolsecurity/sqarol/cmd@latest
```

## CLI

```
sqarol <domain>
sqarol -h | --help
```

Use `-h` or `--help` to print usage information:

```
$ sqarol --help
usage: sqarol <domain>

Generate typosquatting domain variations for the given domain.
Each variation is printed with its technique, variant, edit distance, and effectiveness score.
```

Example:

```
$ sqarol symbolsecurity.com
omission             ymbolsecurity.com                        dist=1   eff=0.69
omission             smbolsecurity.com                        dist=1   eff=0.69
transposition        ysmbolsecurity.com                       dist=2   eff=0.50
typo-trick           5ymb0l5ecurity.c0rn                     dist=9   eff=0.26
phonetic             zymbolzecurity.com                       dist=2   eff=0.50
...
```

Each line shows the technique name, the generated variant, the Levenshtein edit distance, and the visual deceptiveness score.

## Library Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/symbolsecurity/sqarol"
)

func main() {
    vars, err := sqarol.Generate("symbolsecurity.com")
    if err != nil {
        log.Fatal(err)
    }

    for _, v := range vars {
        fmt.Printf("%-20s %-40s dist=%d eff=%.2f\n", v.Kind, v.Variant, v.Distance, v.Effectiveness)
    }
}
```

Full URLs are also accepted; the hostname is extracted automatically:

```go
vars, err := sqarol.Generate("https://symbolsecurity.com/path")
// v.Original will be "symbolsecurity.com"
```

## API

### `Generate(domain string) ([]Variation, error)`

Normalizes the input domain and runs all fuzzing techniques against it, returning the generated variations.

### `Variation`

```go
type Variation struct {
    Original      string  `json:"original"`      // Normalized input domain
    Variant       string  `json:"variant"`        // Generated look-alike domain
    Kind          string  `json:"kind"`           // Algorithm that produced this variation
    Distance      int     `json:"distance"`       // Levenshtein distance from original
    Effectiveness float64 `json:"effectiveness"`  // Visual deceptiveness score (0.0 - 1.0)
}
```

## Techniques

| Technique | Description |
|---|---|
| `omission` | Removes one character at a time |
| `transposition` | Swaps adjacent character pairs |
| `vowel-swap` | Replaces each vowel with every other vowel |
| `hyphenation` | Inserts a hyphen between adjacent characters |
| `repetition` | Doubles each character one at a time |
| `extra-random-letter` | Inserts a random letter at each position |
| `prefix` | Prepends common phishing prefixes (`login`, `auth`, `account`, etc.) |
| `suffix` | Appends common phishing suffixes (`secure`, `verify`, `portal`, etc.) |
| `subdomain` | Inserts a dot at each position to simulate subdomains |
| `tld-swap` | Replaces the TLD with alternatives from a catalog of ~113 TLDs |
| `typo-trick` | Substitutes visually similar characters (`o`/`0`, `l`/`1`, `m`/`rn`, etc.) |
| `phonetic` | Substitutes phonetically similar groups (`f`/`ph`, `c`/`k`, etc.) |
| `inflect` | Generates plural/singular forms of the domain name |
| `complete-with-tld` | Splits name parts that end with a known TLD string |

## Project Structure

```
sqarol/
├── cmd/
│   └── main.go          # CLI entry point
├── generate.go          # Public API: Generate() and domain normalization
├── variation.go         # Variation type definition
├── generators/          # Domain variation algorithms
│   ├── all.go           # Fn type and All() registry
│   ├── omission.go
│   ├── transposition.go
│   ├── vowelswap.go
│   ├── hyphenation.go
│   ├── repetition.go
│   ├── extraletter.go
│   ├── prefix.go
│   ├── suffix.go
│   ├── subdomain.go
│   ├── tldswap.go
│   ├── tldcompletion.go
│   ├── typotrick.go
│   ├── phonetic.go
│   ├── inflect.go
│   └── tld.go           # TLD catalog
└── attributes/          # Variation scoring
    ├── levenshtein.go   # Edit distance computation
    └── effectiveness.go # Visual deceptiveness scoring
```
