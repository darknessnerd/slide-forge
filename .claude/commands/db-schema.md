<!--
  COMMAND: db-schema
  Trigger:  /db-schema
  How:      You type it. Claude does not trigger this automatically.
  Input:    None — connects via MCP postgres (requires DATABASE_URL env var, see .mcp.json).
  Customize: Add table filters, add ERD output, or point to a read replica URL
             by overriding DATABASE_URL in your CLAUDE.local.md.
-->

Fetch and display the current database schema using the MCP postgres connection.

- List all tables with columns, types, and constraints
- Highlight any missing indexes on foreign keys
- Flag columns whose names contain: `token`, `secret`, `password`, `key`, `credential`
  — these may be storing secrets and should be reviewed
