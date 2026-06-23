package service

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"strings"

	"github.com/darknessnerd/slide-forge/internal/derr"
	"github.com/darknessnerd/slide-forge/internal/domain"
)

//go:embed templates/presentation.html.tmpl
var presentationTmpl string

// Renderer builds the final standalone HTML from a Presentation.
type Renderer struct{}

func NewRenderer() *Renderer { return &Renderer{} }

type slideJSON struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Notes   string `json:"notes"`
}

type tmplData struct {
	SlidesJSON              template.JS
	Theme                   domain.Theme
	TransitionStyle         domain.TransitionStyle
	EnableKeyboardNavigation bool
	IncludeProgressBar      bool
	IncludeSpeakerNotes     bool
	ThemeCSS                template.CSS
}

// Render generates a self-contained HTML presentation string.
func (r *Renderer) Render(p *domain.Presentation, req *domain.RenderRequest) (string, error) {
	slides := make([]slideJSON, len(p.Slides))
	for i, s := range p.Slides {
		slides[i] = slideJSON{
			Type:    string(s.Type),
			Title:   s.Title,
			Content: s.Content,
			Notes:   s.Notes,
		}
	}

	raw, err := json.Marshal(slides)
	if err != nil {
		return "", derr.Internal("renderer.Render", "marshal slides", err)
	}

	data := tmplData{
		SlidesJSON:               template.JS(raw),
		Theme:                    req.Theme,
		TransitionStyle:          req.TransitionStyle,
		EnableKeyboardNavigation: req.EnableKeyboardNavigation,
		IncludeProgressBar:       req.IncludeProgressBar,
		IncludeSpeakerNotes:      req.IncludeSpeakerNotes,
		ThemeCSS:                 template.CSS(themeCSS(req.Theme)),
	}

	t, err := template.New("presentation").Parse(presentationTmpl)
	if err != nil {
		return "", derr.Internal("renderer.Render", "template parse", err)
	}

	var buf bytes.Buffer
	if err = t.Execute(&buf, data); err != nil {
		return "", derr.Internal("renderer.Render", "template execute", err)
	}

	return buf.String(), nil
}

func themeCSS(t domain.Theme) string {
	switch t {
	case domain.ThemeDark:
		return darkCSS
	case domain.ThemeMinimal:
		return minimalCSS
	case domain.ThemeCorporate:
		return corporateCSS
	default:
		return lightCSS
	}
}

const lightCSS = `
:root {
  --bg: #ffffff;
  --surface: #f8f9fa;
  --text: #212529;
  --subtext: #6c757d;
  --accent: #4a6cf7;
  --accent-hover: #3451d1;
  --border: #dee2e6;
  --code-bg: #f1f3f5;
  --nav-bg: rgba(255,255,255,0.95);
  --progress-bg: #e9ecef;
  --progress-fill: #4a6cf7;
  --notes-bg: #fffbea;
  --notes-border: #ffc107;
  --slide-shadow: 0 4px 24px rgba(0,0,0,0.08);
}`

const darkCSS = `
:root {
  --bg: #0d1117;
  --surface: #161b22;
  --text: #e6edf3;
  --subtext: #8b949e;
  --accent: #58a6ff;
  --accent-hover: #79b8ff;
  --border: #30363d;
  --code-bg: #1c2128;
  --nav-bg: rgba(13,17,23,0.95);
  --progress-bg: #21262d;
  --progress-fill: #58a6ff;
  --notes-bg: #1c2128;
  --notes-border: #f0883e;
  --slide-shadow: 0 4px 24px rgba(0,0,0,0.4);
}`

const minimalCSS = `
:root {
  --bg: #fafafa;
  --surface: #ffffff;
  --text: #111111;
  --subtext: #555555;
  --accent: #111111;
  --accent-hover: #333333;
  --border: #eeeeee;
  --code-bg: #f5f5f5;
  --nav-bg: rgba(250,250,250,0.98);
  --progress-bg: #eeeeee;
  --progress-fill: #111111;
  --notes-bg: #f9f9f9;
  --notes-border: #cccccc;
  --slide-shadow: none;
}`

const corporateCSS = `
:root {
  --bg: #f0f2f5;
  --surface: #ffffff;
  --text: #1a1a2e;
  --subtext: #4a4a6a;
  --accent: #0052cc;
  --accent-hover: #0041a3;
  --border: #c1c7d0;
  --code-bg: #eaecf0;
  --nav-bg: rgba(240,242,245,0.97);
  --progress-bg: #c1c7d0;
  --progress-fill: #0052cc;
  --notes-bg: #e8f0fe;
  --notes-border: #0052cc;
  --slide-shadow: 0 2px 12px rgba(0,0,0,0.12);
}`

// GenerateHTML is the top-level function called by the MCP handler.
// It parses markdown then renders HTML.
func GenerateHTML(req *domain.RenderRequest) (string, error) {
	if err := validateTheme(req); err != nil {
		return "", err
	}

	presentation, err := Parse(req.Markdown)
	if err != nil {
		return "", err
	}

	r := NewRenderer()
	out, err := r.Render(presentation, req)
	if err != nil {
		return "", err
	}
	return out, nil
}

func validateTheme(req *domain.RenderRequest) error {
	if req.Theme == "" {
		req.Theme = domain.ThemeLight
	}
	if req.TransitionStyle == "" {
		req.TransitionStyle = domain.TransitionFade
	}
	switch req.Theme {
	case domain.ThemeLight, domain.ThemeDark, domain.ThemeMinimal, domain.ThemeCorporate:
	default:
		return derr.Validation("service.GenerateHTML", fmt.Sprintf("unknown theme %q", req.Theme))
	}
	switch req.TransitionStyle {
	case domain.TransitionFade, domain.TransitionSlide, domain.TransitionNone:
	default:
		return derr.Validation("service.GenerateHTML", fmt.Sprintf("unknown transition %q", req.TransitionStyle))
	}
	return nil
}

// joinBool returns "true" or "false" for Alpine.js data initialisation.
func joinBool(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// boolStr exported for template use.
var funcMap = template.FuncMap{
	"boolStr": func(b bool) string { return joinBool(b) },
	"themeClass": func(t domain.Theme) string {
		return "theme-" + strings.ReplaceAll(string(t), " ", "-")
	},
}
