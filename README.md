# sqarol

A Go library and CLI tool designed to help cybersecurity analysts manually monitor lookalike domains. Given a legitimate domain name, it produces hundreds of plausible look-alike variations that attackers might register for phishing, brand impersonation, or credential harvesting. Each variation is annotated with the generation technique, Levenshtein edit distance, and a visual deceptiveness score. The CLI can also check the most effective variations against live DNS and WHOIS data to determine which are registered, who owns them, and whether they appear to be parked.

For automated, continuous domain threat monitoring, check out [Symbol Security's Domain Threat Alerts](https://symbolsecurity.com/domain-threat-alerts).

![sqarol demo](demo.gif)

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
sqarol <command> [arguments]

Commands:
  generate <domain>            Generate typosquatting domain variations
  check <domain> [-n count]    Check status of top N most effective variations

Options:
  -n count   Number of top variations to check (default: 100)
  -h         Show help
```

### generate

Generates all typosquatting variations for a domain, sorted by effectiveness (highest first).

```
$ sqarol generate symbolsecurity.com
KIND                 VARIANT                            DISTANCE  EFFECTIVENESS
----                 -------                            --------  -------------
typo-trick           symbo1security.com                 1         0.95
phonetic             symbolsekurity.com                 1         0.79
vowel-swap           symbolsocurity.com                 1         0.79
...

Total: 467 variations
```

### check

Generates variations, takes the top N by effectiveness, and checks each one for registration status, WHOIS owner, A records, MX records, and parking detection. Checks run concurrently.

```
$ sqarol check symbolsecurity.com -n 5
Checking 5 variations...

#  VARIANT             REGISTERED  OWNER  A   MX  PARKED
-  -------             ----------  -----  -   --  ------
1  symbo1security.com  no          -      no  no  no
2  symbolsocurity.com  no          -      no  no  no
3  symbolsekurity.com  no          -      no  no  no
4  sympolsecurity.com  no          -      no  no  no
5  symbolsecuridy.com  no          -      no  no  no

Checked: 5/467 variations
```

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

Normalizes the input domain and runs all fuzzing techniques against it, returning the generated variations. The domain must be an ASCII hostname or a full URL (the hostname is extracted automatically). Internationalized (non-ASCII) domain names are not supported and will return an error.

### `Check(ctx context.Context, domain string) (*DomainCheck, error)`

Queries DNS and WHOIS to determine whether a domain is registered, who owns it, whether it has A and MX records, and whether it appears to be parked. Registration is determined by the presence of NS records (the most reliable signal). Parking detection uses known parking nameserver suffixes and IP address prefixes. WHOIS is queried separately to extract the registrant owner, with automatic referral following for richer results. Concurrent WHOIS connections are throttled to avoid rate-limiting.

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

check, err := sqarol.Check(ctx, "example.com")
// check.IsRegistered   -> true
// check.Owner          -> "Internet Assigned Numbers Authority"
// check.HasARecords    -> true
// check.HasMXRecords   -> false
// check.IsParked       -> false
```

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

### `DomainCheck`

```go
type DomainCheck struct {
    Domain       string `json:"domain"`                    // Queried domain name
    IsRegistered bool   `json:"is_registered"`             // Whether the domain is registered (via NS lookup)
    Owner        string `json:"owner,omitempty"`           // Registrant from WHOIS, if available
    HasARecords  bool   `json:"has_a_records"`             // Whether the domain has IPv4 address records
    HasMXRecords bool   `json:"has_mx_records"`            // Whether the domain has mail exchange records
    IsParked     bool   `json:"is_parked"`                 // Whether the domain appears to be parked
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

## License

This project is licensed under the [MIT License](LICENSE).
