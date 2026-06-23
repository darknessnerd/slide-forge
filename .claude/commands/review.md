<!--
  COMMAND: review
  Trigger:  /review  or  /review path/to/file.go
  How:      You type it. Claude does not trigger this automatically.
  Input:    Optional — $ARGUMENTS is the file path. If omitted, Claude reviews the current diff.
  Customize: Change severity labels, add project-specific checks,
             or replace the rule file references with your own.
-->

Review the current diff or the file at `$ARGUMENTS` for:
- Correctness bugs
- Convention violations (see `.claude/rules/02-conventions.md`)
- Security issues (see `.claude/rules/04-security.md`)
- Missing tests

Format each finding as:
- **[file:line]** `severity` — problem → fix

Severities: `critical` / `warning` / `suggestion`

Only report findings you are confident about. Do not invent issues.
