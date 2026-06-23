# slide-forge

![Claude Code](https://img.shields.io/badge/Claude_Code-compatible-6B48FF?logo=anthropic&logoColor=white)
![MCP](https://img.shields.io/badge/MCP-slide--forge-0078D4?logo=amazonwebservices&logoColor=white)
![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)
![Maintained](https://img.shields.io/badge/maintained-yes-brightgreen)

MCP server that converts Markdown into standalone HTML slide presentations. Single portable `.html` output — no server required to view.

---

## Local development setup

### 1. Clone and verify it builds

```bash
git clone https://github.com/darknessnerd/slide-forge.git
cd slide-forge
go build ./...
go test -race ./...
```

### 2. Wire it into Claude Code

Add to your project's `.mcp.json` (or `~/.claude/mcp.json` for global use):

```json
{
  "mcpServers": {
    "slide-forge": {
      "command": "go",
      "args": ["run", "./cmd"],
      "cwd": "/absolute/path/to/slide-forge"
    }
  }
}
```

Claude Code spawns the process automatically on startup via stdio — no env vars needed.

**Prereq:** Go 1.21+.

### 3. Verify the tool is available

Open Claude Code in any project and ask:

```
List available MCP tools.
```

You should see `md_to_html_slides` in the list.

### 4. Try it

```
Convert README.md into slides using the slide-forge tool.
Use corporate theme. Save to README-slides.html.
```

Open the generated `.html` file in any browser — no server required.

---

## Configuration

Config is loaded from `config/config.yaml` (or `CONFIG_PATH` env var). All values support `${env:VAR=default}` substitution.

```yaml
# config/config.yaml
transport: ${env:MCP_TRANSPORT=stdio}
addr:      ${env:ADDR=:8080}
log_level: ${env:LOG_LEVEL=info}
```

| Variable | Default | Description |
|---|---|---|
| `MCP_TRANSPORT` | `stdio` | `stdio` for Claude Code local use, `http` for remote/public |
| `ADDR` | `:8080` | HTTP listen address (HTTP mode only) |
| `LOG_LEVEL` | `info` | `debug`, `info`, `warn`, `error` |
| `CONFIG_PATH` | `./config/config.yaml` | Path to config file |

---

## Running in HTTP mode

### 1. Start the server

```bash
MCP_TRANSPORT=http ADDR=:8080 go run ./cmd
```

Or with a pre-built binary:

```bash
go build -o slide-forge ./cmd
MCP_TRANSPORT=http ADDR=:8080 ./slide-forge
```

### 2. Verify it's up

```bash
curl http://localhost:8080/health
# → {"status":"ok"}

curl http://localhost:8080/ready
# → {"status":"ready"}
```

### 3. Call the tool directly (curl)

```bash
curl -X POST http://localhost:8080/mcp/ \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {
      "name": "md_to_html_slides",
      "arguments": {
        "markdown": "# Hello\n\n## Slide two\n\nBody text.",
        "theme": "dark",
        "transition_style": "slide"
      }
    }
  }'
```

Response contains the full standalone HTML in `result.content[0].text`. Write it to a file and open in a browser.

### 4. Connect Claude Code to a remote server

In `.mcp.json` (or `~/.claude/mcp.json`):

```json
{
  "mcpServers": {
    "slide-forge": {
      "type": "http",
      "url": "https://your-domain.com/mcp/"
    }
  }
}
```

Then use it the same way as local:

```
Convert README.md into slides using the slide-forge tool. Use corporate theme.
```

> **Note:** HTTP mode is not production-ready yet — see TODO section below.

---

## Tool: `md_to_html_slides`

| Parameter | Type | Required | Default | Description |
|---|---|---|---|---|
| `markdown` | string | Yes | — | Markdown to convert |
| `theme` | string | No | `light` | `light`, `dark`, `minimal`, `corporate` |
| `transition_style` | string | No | `fade` | `fade`, `slide`, `none` |
| `enable_keyboard_navigation` | bool | No | `true` | Arrow / Space key navigation |
| `include_progress_bar` | bool | No | `true` | Progress bar at top |
| `include_speaker_notes` | bool | No | `true` | Speaker notes panel |

Returns full HTML as a string. Write it to a `.html` file and open in any browser.

---

## Project layout

```
cmd/main.go                  ← entry point; transport selection, health endpoints
config/config.yaml           ← default config with ${env:VAR=default} placeholders
internal/
  config/                    ← AppConfig struct + YAML loader
  derr/                      ← DError: structured errors with Kind + Op + Msg
  domain/slide.go            ← value types (Slide, Theme, RenderRequest, …)
  service/parser.go          ← Markdown → []Slide
  service/renderer.go        ← []Slide → HTML string
  handler/mcp.go             ← MCP tool registration and request handling
```

---

## Development commands

```bash
# build
go build ./...

# tests
go test ./...

# tests with race detector (required before commit)
go test -race ./...

# coverage report
go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out

# run in stdio mode (default — used by Claude Code)
go run ./cmd

# run as HTTP server
MCP_TRANSPORT=http go run ./cmd

# run with debug logging
LOG_LEVEL=debug go run ./cmd
```

---

## TODO — before HTTP production deployment

- [ ] **Auth** — bearer token middleware in `internal/handler/` before `/mcp/` is exposed publicly
- [ ] **TLS** — put behind reverse proxy (nginx/caddy) or add `ListenAndServeTLS`; plain HTTP not acceptable for public endpoints
- [ ] **Body size limit** — wrap MCP handler with `http.MaxBytesReader` to prevent oversized payloads
- [ ] **Rate limiting** — add per-IP rate limiter middleware to prevent abuse
- [ ] **Test coverage** — currently 51.7%; CI floor is 80%; need tests for `internal/config`, `internal/derr`, HTTP handler paths
