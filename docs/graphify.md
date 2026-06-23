# Graphify — Team Guide

## What It Is

Graphify turns a codebase (or any folder of code, docs, PDFs, images) into a **persistent, queryable knowledge graph**. Instead of grepping files or reading them raw, Claude navigates the graph — 71× fewer tokens per query on a typical mixed corpus.

Every relationship is tagged: `EXTRACTED` (found directly), `INFERRED` (deduced), or `AMBIGUOUS` (uncertain). You see exactly what was found vs. guessed.

## Why Use It

| Without graphify | With graphify |
|-----------------|---------------|
| Claude reads 40 files to answer "what calls the auth middleware?" | Claude queries the graph: 1 traversal, 1 answer |
| Context fills up on large repos | Graph persists across sessions — no re-scan |
| New team members grep to understand structure | `/graphify query "how does X work?"` gives a guided tour |
| Agent workflows re-read files on every task | Parallel agents share one up-to-date graph |

## When to Use It

- **Onboarding** — new team member runs `/graphify .` once, then asks questions
- **Code review** — trace data flow: `/graphify path "HTTPHandler" "Database"`
- **Architecture docs** — open `graph.html` in a browser, screenshot for slides
- **Refactoring** — find all callers/dependencies before touching a package
- **Cross-repo work** — `/graphify <url1> <url2>` merges multiple repos into one graph
- **Agent workflows** — `--watch` keeps graph current as multiple agents write code simultaneously

**Not worth it:** repos under ~10 files. The graph pays off at 20+ files.

## Install (one-time, per machine)

```bash
# Recommended
uv tool install graphifyy && graphify install

# Alternatively (macOS externally-managed Python)
pipx install graphifyy && graphify install

# Windows: after pip install, add to PATH:
# %APPDATA%\Python\Python3xx\Scripts
```

`graphify install` wires the skill into Claude Code (`~/.claude/skills/graphify/`).  
No API key needed for code — AST extraction runs fully offline.  
Docs/images/PDFs use your Claude session (already active in Claude Code).

## Quick Start

```bash
# 1. Build the graph for this repo
/graphify .

# Outputs in graphify-out/:
#   graph.html       ← open in browser for the slide
#   GRAPH_REPORT.md  ← god nodes, surprising connections, suggested questions
#   graph.json       ← persists across sessions

# 2. Ask questions
/graphify query "how does the handler layer talk to the service layer?"
/graphify path "HTTPHandler" "Repository"
/graphify explain "domain.ErrNotFound"

# 3. Update after code changes (fast — only re-extracts changed files)
/graphify . --update
```

## Team Workflow

```bash
# Commit the graph so teammates start with it already built
echo "graphify-out/cache/" >> .gitignore   # exclude cache, keep outputs
git add graphify-out/graph.json graphify-out/GRAPH_REPORT.md
git commit -m "chore: add initial knowledge graph"

# Install the git hook so the graph auto-rebuilds on every commit
graphify hook install

# Exclude files you don't want graphed (mirrors .gitignore syntax)
echo "node_modules/" >> .graphifyignore
echo "graphify-out/" >> .graphifyignore
echo "vendor/" >> .graphifyignore
```

## Key Commands

| Command | What it does |
|---------|-------------|
| `/graphify .` | Full build on current directory |
| `/graphify . --update` | Re-extract only changed files |
| `/graphify . --mode deep` | Richer extraction (more inferred edges, slower) |
| `/graphify query "<question>"` | BFS traversal — broad context |
| `/graphify path "A" "B"` | Shortest path between two concepts |
| `/graphify explain "<concept>"` | Plain-language node explanation |
| `/graphify . --watch` | Auto-rebuild on file changes |
| `/graphify . --wiki` | Generate agent-crawlable markdown wiki |
| `/graphify . --svg` | Export SVG (embeds in Notion, GitHub) |
| `graphify hook install` | Post-commit auto-rebuild |

## Outputs

```
graphify-out/
  graph.html        ← interactive: click, filter, search nodes
  graph.json        ← persistent graph data (reused across sessions)
  GRAPH_REPORT.md   ← god nodes, surprising connections, suggested questions
  wiki/             ← (with --wiki) one .md per community cluster
  cache/            ← SHA256 cache, skip to .gitignore
```

### What to Commit vs Ignore

| File | Commit? | Why |
|------|---------|-----|
| `graphify-out/graph.json` | ✅ Yes | Persists across sessions — teammates query without rebuilding |
| `graphify-out/GRAPH_REPORT.md` | ✅ Yes | Human-readable audit trail, useful in PRs |
| `graphify-out/graph.html` | ✅ Yes | Shareable for slides — no server needed |
| `graphify-out/cache/` | ❌ Ignore | Per-file SHA256 cache, reconstructed automatically |
| `graphify-out/cost.json` | ⚠️ Optional | Tracks team token spend across runs |
| `graphify-out/.graphify_*` | ❌ Ignore | Temp files, deleted after each run |

Add to `.gitignore`:
```
graphify-out/cache/
```

## Privacy

- **Code** — processed 100% locally via tree-sitter. No API calls.
- **Docs/PDFs/images** — use your active Claude Code session (same model you're already talking to).
- **Query logs** — written to `~/.cache/graphify-queries.log`. Disable: `export GRAPHIFY_QUERY_LOG_DISABLE=1`

## Running on This Scaffold (for slides)

```bash
# From repo root — builds a graph of the scaffold's own structure
/graphify .

# Then open in browser
open graphify-out/graph.html
```

The graph will show:
- `cmd/` → `internal/handler` → `internal/service` → `internal/domain` ← `internal/repository` dependency chain
- `.claude/rules/` as a documentation cluster
- Cross-links between conventions, security rules, and testing rules
