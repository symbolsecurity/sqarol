#!/bin/bash
#
# sqarol CLI demo script for asciinema
#
# Record with:
#   asciinema rec demo.cast --command="bash demo.sh"
#
# Play back with:
#   asciinema play demo.cast
#
# Prerequisites:
#   - asciinema installed (https://asciinema.org)
#   - sqarol installed and in $PATH
#     go install github.com/symbolsecurity/sqarol/cmd@latest
#

set -e

# --- Typing simulation helpers ---

TYPE_DELAY=0.04    # seconds between keystrokes
CMD_PAUSE=1.5      # pause before executing a typed command
READ_PAUSE=3       # pause to let the viewer read output
LONG_READ_PAUSE=5  # longer pause for dense output
PROMPT="$ "

# Simulate typing a command character by character
type_command() {
    local cmd="$1"
    printf "%s" "$PROMPT"
    for (( i=0; i<${#cmd}; i++ )); do
        printf "%s" "${cmd:$i:1}"
        sleep "$TYPE_DELAY"
    done
    sleep "$CMD_PAUSE"
    printf "\n"
}

# Type and execute a command
run() {
    type_command "$1"
    eval "$1"
}

# Print a comment line (dimmed)
comment() {
    printf "\033[2m# %s\033[0m\n" "$1"
    sleep 1
}

# --- Demo starts here ---

clear
sleep 0.5

comment "sqarol - domain typosquatting analysis CLI"
comment "https://github.com/symbolsecurity/sqarol"
echo
sleep 1

# Scene 1: Show help
comment "Let's start by looking at the available commands."
echo
run "sqarol -h"
sleep "$LONG_READ_PAUSE"

echo
# Scene 2: Generate command
comment "The 'generate' command produces look-alike domain variations."
comment "Each variant is scored by visual deceptiveness (effectiveness)."
echo
sleep 1
run "sqarol generate symbolsecurity.com"
sleep "$LONG_READ_PAUSE"

echo
# Scene 3: Check command
comment "The 'check' command verifies the top variations against live DNS and WHOIS."
comment "Let's check the 5 most effective variations."
echo
sleep 1
run "sqarol check symbolsecurity.com -n 5"
sleep "$LONG_READ_PAUSE"

echo
comment "Done! Learn more at github.com/symbolsecurity/sqarol"
comment "For automated, continuous domain threat monitoring check out:"
comment "https://symbolsecurity.com/domain-threat-alerts"
sleep 3
