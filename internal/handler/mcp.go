package handler

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/darknessnerd/slide-forge/internal/domain"
	"github.com/darknessnerd/slide-forge/internal/service"
)

type toolArgs struct {
	Markdown                 string `json:"markdown"`
	Theme                    string `json:"theme,omitempty"                      jsonschema:"Visual theme: light (default), dark, minimal, corporate."`
	TransitionStyle          string `json:"transition_style,omitempty"           jsonschema:"Slide transition: fade (default), slide, none."`
	EnableKeyboardNavigation *bool  `json:"enable_keyboard_navigation,omitempty" jsonschema:"Enable arrow-key / space navigation (default true)."`
	IncludeProgressBar       *bool  `json:"include_progress_bar,omitempty"       jsonschema:"Show progress bar at the top (default true)."`
	IncludeSpeakerNotes      *bool  `json:"include_speaker_notes,omitempty"      jsonschema:"Show speaker notes panel below slides (default true)."`
}

// slideGenerator is the interface the handler needs from the service layer.
type slideGenerator interface {
	Generate(req *domain.RenderRequest) (html string, err error)
}

// NewMCPServer builds and returns an mcp.Server with the md_to_html_slides tool registered.
func NewMCPServer(gen slideGenerator) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{
		Name:    "slide-forge",
		Version: "1.0.0",
	}, nil)

	mcp.AddTool(s, &mcp.Tool{
		Name: "md_to_html_slides",
		Description: "Convert Markdown into a standalone HTML presentation (no server) " +
			"using Alpine.js for interactivity and PureCSS for styling. " +
			"Returns a single portable .html file content.",
	}, func(ctx context.Context, ss *mcp.ServerSession, req *mcp.CallToolParamsFor[toolArgs]) (*mcp.CallToolResult, error) {
		return handleGenerate(ctx, req.Arguments, gen)
	})

	return s
}

func handleGenerate(_ context.Context, args toolArgs, gen slideGenerator) (*mcp.CallToolResult, error) {
	if args.Markdown == "" {
		return errorResult("missing required param: markdown"), nil
	}

	domainReq := &domain.RenderRequest{
		Markdown:                 args.Markdown,
		Theme:                    domain.Theme(orDefault(args.Theme, string(domain.ThemeLight))),
		TransitionStyle:          domain.TransitionStyle(orDefault(args.TransitionStyle, string(domain.TransitionFade))),
		EnableKeyboardNavigation: boolVal(args.EnableKeyboardNavigation, true),
		IncludeProgressBar:       boolVal(args.IncludeProgressBar, true),
		IncludeSpeakerNotes:      boolVal(args.IncludeSpeakerNotes, true),
	}

	html, err := gen.Generate(domainReq)
	if err != nil {
		return errorResult("slide-forge: " + err.Error()), nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: html}},
	}, nil
}

func errorResult(msg string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{&mcp.TextContent{Text: msg}},
	}
}

func orDefault(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

func boolVal(p *bool, def bool) bool {
	if p == nil {
		return def
	}
	return *p
}

// generatorAdapter wraps service.GenerateHTML to satisfy slideGenerator.
type generatorAdapter struct{}

func (g *generatorAdapter) Generate(req *domain.RenderRequest) (string, error) {
	return service.GenerateHTML(req)
}

// NewGeneratorAdapter returns the default service-backed generator.
func NewGeneratorAdapter() slideGenerator { return &generatorAdapter{} }
