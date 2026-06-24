# Copilot Instructions for slide-forge

- This repo builds Markdown into standalone HTML slide decks.
- Main entry point is `cmd/main.go`.
- Slide parsing lives in `internal/service/parser.go`.
- HTML rendering lives in `internal/service/renderer.go`.
- MCP tool name is `md_to_html_slides`.
- Config comes from `config/config.yaml` and supports `${env:VAR=default}`.
- Use existing `.claude/` rules, commands, skills, and hooks as repo guidance.
- Keep changes surgical and Go idiomatic.
