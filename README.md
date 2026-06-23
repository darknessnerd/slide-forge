# slide-forge

![Claude Code](https://img.shields.io/badge/Claude_Code-compatible-6B48FF?logo=anthropic&logoColor=white)
![MCP](https://img.shields.io/badge/MCP-slide--forge-0078D4?logo=amazonwebservices&logoColor=white)
![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)
![Maintained](https://img.shields.io/badge/maintained-yes-brightgreen)

MCP server that converts Markdown into standalone HTML slide presentations. Single portable `.html` output — no server required to view.

---

## Quick start (local dev)

Claude Code spawns the server automatically via stdio. Add to your project's `.mcp.json`:

```json
{
  "mcpServers": {
    "slide-forge": {
      "command": "go",
      "args": ["run", "."],
      "cwd": "/path/to/slide-forge"
    }
  }
}
```

**Prereq:** Go 1.21+. No env vars needed.

---

## Deploying as a public HTTP server

```bash
MCP_TRANSPORT=http ADDR=:8080 go run .
```

The server exposes a single endpoint: `POST /mcp`

Clients connect via:

```json
{
  "mcpServers": {
    "slide-forge": {
      "type": "http",
      "url": "https://your-domain.com/mcp"
    }
  }
}
```

### Environment variables

| Variable | Default | Description |
|---|---|---|
| `MCP_TRANSPORT` | `stdio` | `stdio` for local Claude Code use, `http` for remote/public |
| `ADDR` | `:8080` | Listen address (HTTP mode only) |

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

### Example prompt

```
Convert README.md into slides using the slide-forge tool.
Use corporate theme. Save to README-slides.html.
```

---

## Project layout

```
main.go                   ← entry point; MCP_TRANSPORT selects stdio vs http
internal/
  domain/slide.go         ← value types (Slide, Theme, RenderRequest, …)
  service/parser.go       ← Markdown → []Slide
  service/renderer.go     ← []Slide → HTML string
  handler/mcp.go          ← MCP tool registration and request handling
```

---

## Development

```bash
# run tests
go test ./...

# run with race detector
go test -race ./...

# start locally in stdio mode (default)
go run .

# start as HTTP server
MCP_TRANSPORT=http go run .
```
