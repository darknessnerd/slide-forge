# Graph Report - .  (2026-06-22)

## Corpus Check
- Corpus is ~11,405 words - fits in a single context window. You may not need a graph.

## Summary
- 122 nodes · 131 edges · 13 communities (12 shown, 1 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 9 edges (avg confidence: 0.81)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_CI & Quality Rules|CI & Quality Rules]]
- [[_COMMUNITY_Architecture & Commands|Architecture & Commands]]
- [[_COMMUNITY_Scaffold Overview & Security|Scaffold Overview & Security]]
- [[_COMMUNITY_Claude Skills Package|Claude Skills Package]]
- [[_COMMUNITY_MCP Servers & Env Vars|MCP Servers & Env Vars]]
- [[_COMMUNITY_Skill Installer Script|Skill Installer Script]]
- [[_COMMUNITY_Safety Hooks|Safety Hooks]]
- [[_COMMUNITY_NPM Package Config|NPM Package Config]]
- [[_COMMUNITY_Caveman Installer|Caveman Installer]]
- [[_COMMUNITY_SOLID Principles|SOLID Principles]]
- [[_COMMUNITY_Graphify Guide|Graphify Guide]]
- [[_COMMUNITY_Go Module|Go Module]]

## God Nodes (most connected - your core abstractions)
1. `README: claude-baseline-go scaffold overview` - 10 edges
2. `Testing Rules` - 7 edges
3. `Security Rules` - 7 edges
4. `datadog` - 6 edges
5. `01-architecture.md (Go layer architecture rules)` - 6 edges
6. `Skill: solid-principles` - 6 edges
7. `02-conventions.md (Go conventions)` - 5 edges
8. `.golangci.yml linter configuration` - 5 edges
9. `github` - 4 edges
10. `postgres` - 4 edges

## Surprising Connections (you probably didn't know these)
- `Observability interfaces in domain (Logger, MetricsRecorder) injected via constructor` --conceptually_related_to--> `MCP Datadog Server (@datadog/mcp-server-datadog)`  [INFERRED]
  .claude/rules/01-architecture.md → .mcp.json
- `/standup command` --references--> `MCP GitHub Server (@modelcontextprotocol/server-github)`  [EXTRACTED]
  .claude/commands/standup.md → .mcp.json
- `/db-schema command` --references--> `MCP Postgres Server (@modelcontextprotocol/server-postgres)`  [EXTRACTED]
  .claude/commands/db-schema.md → .mcp.json
- `gosec linter (security checks, hardcoded creds, weak crypto, SQLi)` --conceptually_related_to--> `Security Rules`  [INFERRED]
  .golangci.yml → .claude/rules/04-security.md
- `Vulnerability Classes (SQLi, XSS, SSRF, cmd injection, path traversal, open redirect)` --conceptually_related_to--> `gosec linter (security checks, hardcoded creds, weak crypto, SQLi)`  [INFERRED]
  .claude/rules/04-security.md → .golangci.yml

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Two-layer command enforcement: settings.json deny list + pre-bash.sh hard safety net** — settings_json, hooks_pre_bash, hooks_pre_bash_blockedpatterns [EXTRACTED 1.00]
- **MCP server trio: GitHub, Postgres, Datadog — all env-var credentials, no hardcoding** — mcp_json, mcp_github_server, mcp_postgres_server, mcp_datadog_server [EXTRACTED 1.00]
- **Skill install pipeline: package.json → install.js → team skills + caveman-installer → .claude/skills/** — package_json, claude_skills_install_js, caveman_installer, claude_skills_target_dir [EXTRACTED 1.00]
- **Security enforcement triad: security rules, gosec linter, pre-commit check** — rules_04_security_security_rules, golangci_yml_gosec, skills_pre_commit_check_skill [INFERRED 0.85]
- **CI quality gates enforced by testing rules, CI rules, and linter config** — rules_03_testing_testing_rules, rules_05_ci_ci_rules, golangci_yml_linter_config [INFERRED 0.85]
- **Auto-triggering skills cluster** — skills_on_new_file_skill, skills_pre_commit_check_skill, skills_explain_error_skill, skills_c4_architecture_skill, skills_solid_principles_skill, skills_frontend_design_skill [EXTRACTED 1.00]

## Communities (13 total, 1 thin omitted)

### Community 0 - "CI & Quality Rules"
Cohesion: 0.12
Nodes (21): CLAUDE.md project configuration (team shared), errcheck linter (unchecked errors, type assertions, blank), gofumpt linter (strict formatting, superset of gofmt), gosec linter (security checks, hardcoded creds, weak crypto, SQLi), .golangci.yml linter configuration, Coverage Expectations (80% floor, 100% domain/service exports), LLM-as-judge grading for agent output tests, What Must NOT Be Mocked (DB, context, sentinel errors, Logger) (+13 more)

### Community 1 - "Architecture & Commands"
Cohesion: 0.13
Nodes (17): All implementation packages must live under /internal/, Layer dependency direction: handler → service → domain ← repository, Observability interfaces in domain (Logger, MetricsRecorder) injected via constructor, /db-schema command, /review command, /standup command, Forbidden patterns: global mutable state, interface{}/any, init() with I/O, type switches for behavior, Interfaces belong in consuming package, not implementing package (+9 more)

### Community 2 - "Scaffold Overview & Security"
Cohesion: 0.15
Nodes (14): Minimal Fork Path (Stage 1: no DB, Stage 2: persistence, Stage 3: scale), README: claude-baseline-go scaffold overview, Skill distribution via GitHub Packages (npm versioned packages), Secrets Policy (no hardcode, no .env write, env vars only), C4 Model (4 levels: Context, Container, Component, Code), Mermaid diagram templates (C4 levels), Skill: c4-architecture, Skill: explain-error (+6 more)

### Community 3 - "Claude Skills Package"
Cohesion: 0.15
Nodes (12): bin, claude-skills, description, files, keywords, license, name, publishConfig (+4 more)

### Community 4 - "MCP Servers & Env Vars"
Cohesion: 0.18
Nodes (12): DD_API_KEY, DD_APP_KEY, DD_SITE, GITHUB_PERSONAL_ACCESS_TOKEN, POSTGRES_CONNECTION_STRING, npx, datadog, github (+4 more)

### Community 5 - "Skill Installer Script"
Cohesion: 0.28
Nodes (7): copySkill(), dryRun, fs, installCavemanSkills(), installFlatSkills(), path, targetDir

### Community 6 - "Safety Hooks"
Cohesion: 0.29
Nodes (6): post-tool-use.sh script, pre-bash.sh script, Tool Audit Log (.claude/logs/tool-audit.log), BlockedPatterns (pre-bash hard safety net), settings.json (Claude Code permissions + hooks), settings.local.json (personal permission overrides)

### Community 7 - "NPM Package Config"
Cohesion: 0.25
Nodes (7): devDependencies, @team/claude-skills, name, private, scripts, setup:claude, workspaces

### Community 8 - "Caveman Installer"
Cohesion: 0.38
Nodes (7): caveman-installer (JuliusBrussee/caveman), install.js (claude-skills CLI), dependencies, caveman-installer, @team/claude-skills package.json, .claude/skills/ (skill install target directory), package.json (root workspace)

### Community 9 - "SOLID Principles"
Cohesion: 0.40
Nodes (6): DIP (Dependency Inversion Principle) Go mapping, ISP (Interface Segregation Principle) Go mapping, LSP (Liskov Substitution Principle) Go mapping, OCP (Open/Closed Principle) Go mapping, Skill: solid-principles, SRP (Single Responsibility Principle) Go mapping

### Community 10 - "Graphify Guide"
Cohesion: 0.67
Nodes (3): Graphify commands (query, path, explain, watch, wiki, svg), Graphify team guide, Knowledge Graph (persistent, queryable, tagged edges)

## Knowledge Gaps
- **55 isolated node(s):** `post-tool-use.sh script`, `pre-bash.sh script`, `@modelcontextprotocol/server-github`, `GITHUB_PERSONAL_ACCESS_TOKEN`, `@modelcontextprotocol/server-postgres` (+50 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **1 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `README: claude-baseline-go scaffold overview` connect `Scaffold Overview & Security` to `CI & Quality Rules`, `SOLID Principles`?**
  _High betweenness centrality (0.073) - this node is a cross-community bridge._
- **Why does `CLAUDE.md project configuration (team shared)` connect `CI & Quality Rules` to `Scaffold Overview & Security`?**
  _High betweenness centrality (0.038) - this node is a cross-community bridge._
- **Are the 3 inferred relationships involving `01-architecture.md (Go layer architecture rules)` (e.g. with `/db-schema command` and `/review command`) actually correct?**
  _`01-architecture.md (Go layer architecture rules)` has 3 INFERRED edges - model-reasoned connections that need verification._
- **What connects `post-tool-use.sh script`, `pre-bash.sh script`, `@modelcontextprotocol/server-github` to the rest of the system?**
  _58 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `CI & Quality Rules` be split into smaller, more focused modules?**
  _Cohesion score 0.11904761904761904 - nodes in this community are weakly interconnected._
- **Should `Architecture & Commands` be split into smaller, more focused modules?**
  _Cohesion score 0.1323529411764706 - nodes in this community are weakly interconnected._