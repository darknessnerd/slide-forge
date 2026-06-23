# Claude Configuration — Team Shared

Do NOT add personal preferences here. Use `CLAUDE.local.md` (gitignored).

## What Claude Is

AI pair-programmer for this repo. Reads code, writes code, runs safe shell commands.
Not an autonomous agent — confirm before destructive or irreversible ops.

## Architecture

See `.claude/rules/01-architecture.md` for full details.

## Conventions

See `.claude/rules/02-conventions.md` for full details.

## Testing

See `.claude/rules/03-testing.md` for full details.

## Security

See `.claude/rules/04-security.md` for full details.

## What Claude Can Touch

Controlled via `.claude/settings.json`. Summary:
- Read: anything
- Write/Edit: source files, configs (non-secret)
- Run: build, test, lint commands
- Never: force-push, drop tables, write `.env` files

## MCP Connections

Defined in `.mcp.json`: GitHub, database, Datadog.
Credentials come from env vars — never hardcoded.
