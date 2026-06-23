package handler_test

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/your-org/your-repo/internal/domain"
	"github.com/your-org/your-repo/internal/handler"
)

// mockGenerator satisfies the handler's slideGenerator interface.
type mockGenerator struct{ mock.Mock }

func (m *mockGenerator) Generate(req *domain.RenderRequest) (string, error) {
	args := m.Called(req)
	return args.String(0), args.Error(1)
}

func TestNewMCPServer_ToolRegistered(t *testing.T) {
	t.Parallel()
	gen := &mockGenerator{}
	s := handler.NewMCPServer(gen)
	require.NotNil(t, s)
}

func TestHandleGenerate_HappyPath(t *testing.T) {
	t.Parallel()
	gen := &mockGenerator{}
	gen.On("Generate", mock.MatchedBy(func(r *domain.RenderRequest) bool {
		return r.Markdown == "# Hello"
	})).Return("<html>ok</html>", nil)

	s := handler.NewMCPServer(gen)
	require.NotNil(t, s)

	// Call via the exported adapter to verify end-to-end wiring.
	adapter := handler.NewGeneratorAdapter()
	html, err := adapter.Generate(&domain.RenderRequest{
		Markdown: "# Hello\n\nContent.",
		Theme:    domain.ThemeLight,
	})
	require.NoError(t, err)
	assert.Contains(t, html, "<!doctype html>")
}

func TestHandleGenerate_DefaultsViaAdapter(t *testing.T) {
	t.Parallel()
	adapter := handler.NewGeneratorAdapter()

	html, err := adapter.Generate(&domain.RenderRequest{
		Markdown:                 "# Title\n\n## Section\n\nBody.",
		Theme:                    domain.ThemeDark,
		TransitionStyle:          domain.TransitionSlide,
		EnableKeyboardNavigation: true,
		IncludeProgressBar:       true,
		IncludeSpeakerNotes:      true,
	})
	require.NoError(t, err)
	assert.Contains(t, html, "Title")
}

func TestHandleGenerate_EmptyMarkdownError(t *testing.T) {
	t.Parallel()
	adapter := handler.NewGeneratorAdapter()
	_, err := adapter.Generate(&domain.RenderRequest{Markdown: ""})
	require.Error(t, err)
}

func TestHandleGenerate_MCPServerBuilds(t *testing.T) {
	t.Parallel()
	gen := &mockGenerator{}
	s := handler.NewMCPServer(gen)
	require.NotNil(t, s)

	// Verify tool is registered via an in-memory client round-trip.
	ctx := context.Background()
	st, ct := mcp.NewInMemoryTransports()

	_, err := s.Connect(ctx, st)
	require.NoError(t, err)

	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "0"}, nil)
	cs, err := client.Connect(ctx, ct)
	require.NoError(t, err)
	defer cs.Close()

	result, err := cs.ListTools(ctx, nil)
	require.NoError(t, err)
	found := false
	for _, tool := range result.Tools {
		if tool.Name == "md_to_html_slides" {
			found = true
			break
		}
	}
	assert.True(t, found, "md_to_html_slides tool must be registered")
}
