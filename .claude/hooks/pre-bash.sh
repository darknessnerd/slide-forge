#!/usr/bin/env bash
# Hook: pre-bash
# Event:  PreToolUse (Bash)
# Fires:  before every Bash tool call Claude makes
# Effect: if a blocked pattern matches → exit 2 → Claude never runs the command
#
# How Claude Code calls this:
#   - passes a JSON object via stdin with the command at .tool_input.command

COMMAND=$(jq -r '.tool_input.command' < /dev/stdin)

# Patterns blocked unconditionally — regardless of settings.json allow list.
# settings.json is the first gate; this script is the hard safety net.
# Add patterns your team considers unrecoverable (data loss, secret leaks, etc.)
BLOCKED_PATTERNS=(
  "git push --force"       # history rewrite on shared branches
  "git reset --hard"       # discards uncommitted work without warning
  "rm -rf /"               # filesystem wipe
  "DROP TABLE"             # destructive DDL — use migrations
  "DELETE FROM.*WHERE.*1=1" # full-table delete with no real filter
  "> .env"                 # overwrite / create .env — secrets must never be written by Claude
)

for pattern in "${BLOCKED_PATTERNS[@]}"; do
  # -q: silent, -i: case-insensitive
  if echo "$COMMAND" | grep -qi "$pattern"; then
    echo "BLOCKED: command matches forbidden pattern '$pattern'" >&2
    exit 2  # exit 2 → Claude Code blocks the tool call and reports the error
  fi
done

exit 0  # allow
