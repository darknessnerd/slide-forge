package service

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/your-org/your-repo/internal/domain"
)

func TestGenerateHTML_EmptyMarkdown(t *testing.T) {
	t.Parallel()
	_, err := GenerateHTML(&domain.RenderRequest{Markdown: ""})
	assert.ErrorIs(t, err, domain.ErrEmptyMarkdown)
}

func TestGenerateHTML_InvalidTheme(t *testing.T) {
	t.Parallel()
	_, err := GenerateHTML(&domain.RenderRequest{
		Markdown: "# Hello",
		Theme:    "neon",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown theme")
}

func TestGenerateHTML_InvalidTransition(t *testing.T) {
	t.Parallel()
	_, err := GenerateHTML(&domain.RenderRequest{
		Markdown:        "# Hello",
		TransitionStyle: "zoom",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown transition")
}

func TestGenerateHTML_DefaultsApplied(t *testing.T) {
	t.Parallel()
	req := &domain.RenderRequest{Markdown: "# Hello\n\nWorld."}
	html, err := GenerateHTML(req)
	require.NoError(t, err)
	// defaults should be filled in
	assert.Equal(t, domain.ThemeLight, req.Theme)
	assert.Equal(t, domain.TransitionFade, req.TransitionStyle)
	assert.Contains(t, html, "<!doctype html>")
}

func TestGenerateHTML_ContainsAlpineAndPureCSS(t *testing.T) {
	t.Parallel()
	html, err := GenerateHTML(&domain.RenderRequest{
		Markdown: "# Test\n\nContent.",
		Theme:    domain.ThemeLight,
	})
	require.NoError(t, err)
	assert.Contains(t, html, "alpinejs")
	assert.Contains(t, html, "purecss")
}

func TestGenerateHTML_ContainsSlideData(t *testing.T) {
	t.Parallel()
	html, err := GenerateHTML(&domain.RenderRequest{
		Markdown: "# My Title\n\nSlide content here.",
		Theme:    domain.ThemeDark,
	})
	require.NoError(t, err)
	assert.Contains(t, html, "My Title")
}

func TestGenerateHTML_ProgressBarToggle(t *testing.T) {
	t.Parallel()
	withBar, err := GenerateHTML(&domain.RenderRequest{
		Markdown:           "# Hello",
		Theme:              domain.ThemeLight,
		IncludeProgressBar: true,
	})
	require.NoError(t, err)

	withoutBar, err := GenerateHTML(&domain.RenderRequest{
		Markdown:           "# Hello",
		Theme:              domain.ThemeLight,
		IncludeProgressBar: false,
	})
	require.NoError(t, err)

	// With bar: the div#progress-bar-wrap element must be present.
	assert.Contains(t, withBar, `id="progress-bar-wrap"`)
	// Without bar: the element must be absent (CSS class name may still appear in style block).
	assert.NotContains(t, withoutBar, `id="progress-bar-wrap"`)
}

func TestGenerateHTML_SpeakerNotesToggle(t *testing.T) {
	t.Parallel()
	withNotes, err := GenerateHTML(&domain.RenderRequest{
		Markdown:            "# Hello",
		Theme:               domain.ThemeLight,
		IncludeSpeakerNotes: true,
	})
	require.NoError(t, err)

	withoutNotes, err := GenerateHTML(&domain.RenderRequest{
		Markdown:            "# Hello",
		Theme:               domain.ThemeLight,
		IncludeSpeakerNotes: false,
	})
	require.NoError(t, err)

	// With notes: the div#notes-panel element must be present.
	assert.Contains(t, withNotes, `id="notes-panel"`)
	// Without notes: the element must be absent.
	assert.NotContains(t, withoutNotes, `id="notes-panel"`)
}

func TestGenerateHTML_KeyboardNavigationToggle(t *testing.T) {
	t.Parallel()
	withKb, err := GenerateHTML(&domain.RenderRequest{
		Markdown:                 "# Hello",
		Theme:                    domain.ThemeLight,
		EnableKeyboardNavigation: true,
	})
	require.NoError(t, err)

	withoutKb, err := GenerateHTML(&domain.RenderRequest{
		Markdown:                 "# Hello",
		Theme:                    domain.ThemeLight,
		EnableKeyboardNavigation: false,
	})
	require.NoError(t, err)

	assert.Contains(t, withKb, "ArrowRight")
	assert.NotContains(t, withoutKb, "ArrowRight")
}

func TestGenerateHTML_AllThemes(t *testing.T) {
	t.Parallel()
	themes := []domain.Theme{
		domain.ThemeLight,
		domain.ThemeDark,
		domain.ThemeMinimal,
		domain.ThemeCorporate,
	}
	for _, theme := range themes {
		theme := theme
		t.Run(string(theme), func(t *testing.T) {
			t.Parallel()
			html, err := GenerateHTML(&domain.RenderRequest{
				Markdown: "# Hello\n\nContent.",
				Theme:    theme,
			})
			require.NoError(t, err)
			assert.True(t, strings.HasPrefix(strings.TrimSpace(html), "<!doctype html>"),
				"theme=%q: expected HTML document", theme)
		})
	}
}

func TestGenerateHTML_AllTransitions(t *testing.T) {
	t.Parallel()
	transitions := []domain.TransitionStyle{
		domain.TransitionFade,
		domain.TransitionSlide,
		domain.TransitionNone,
	}
	for _, tr := range transitions {
		tr := tr
		t.Run(string(tr), func(t *testing.T) {
			t.Parallel()
			html, err := GenerateHTML(&domain.RenderRequest{
				Markdown:        "# Hello\n\nContent.",
				TransitionStyle: tr,
			})
			require.NoError(t, err)
			assert.Contains(t, html, "<!doctype html>")
		})
	}
}

func TestGenerateHTML_StandaloneNoServer(t *testing.T) {
	t.Parallel()
	html, err := GenerateHTML(&domain.RenderRequest{
		Markdown: "# Test",
		Theme:    domain.ThemeLight,
	})
	require.NoError(t, err)
	// Must be self-contained — no server-side template tags
	assert.NotContains(t, html, "{{")
	assert.NotContains(t, html, "}}")
	// CDN links present
	assert.Contains(t, html, "cdn.jsdelivr.net")
}
