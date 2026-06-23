package service

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/your-org/your-repo/internal/domain"
)

func TestParse_EmptyInput(t *testing.T) {
	t.Parallel()
	_, err := Parse("")
	assert.ErrorIs(t, err, domain.ErrEmptyMarkdown)
}

func TestParse_WhitespaceOnly(t *testing.T) {
	t.Parallel()
	_, err := Parse("   \n\t  ")
	assert.ErrorIs(t, err, domain.ErrEmptyMarkdown)
}

func TestParse_TitleSlide(t *testing.T) {
	t.Parallel()
	p, err := Parse("# Hello World\n\nSome intro text.")
	require.NoError(t, err)
	require.Len(t, p.Slides, 1)
	assert.Equal(t, domain.SlideTypeTitle, p.Slides[0].Type)
	assert.Equal(t, "Hello World", p.Slides[0].Title)
	assert.Contains(t, p.Slides[0].Content, "intro text")
}

func TestParse_MultipleSlides(t *testing.T) {
	t.Parallel()
	md := `# Title

Intro.

## First Section

Content here.

## Second Section

More content.
`
	p, err := Parse(md)
	require.NoError(t, err)
	assert.Len(t, p.Slides, 3)
	assert.Equal(t, domain.SlideTypeTitle, p.Slides[0].Type)
	assert.Equal(t, "Title", p.Slides[0].Title)
}

func TestParse_SummarySlideDetection(t *testing.T) {
	t.Parallel()
	cases := []struct {
		title    string
		wantType domain.SlideType
	}{
		{"Summary", domain.SlideTypeSummary},
		{"Conclusion and Next Steps", domain.SlideTypeSummary},
		{"Recap", domain.SlideTypeSummary},
		{"Thank You", domain.SlideTypeSummary},
		{"Architecture Overview", domain.SlideTypeContent},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()
			md := "# First\n\nContent.\n\n## " + tc.title + "\n\nFinal."
			p, err := Parse(md)
			require.NoError(t, err)
			last := p.Slides[len(p.Slides)-1]
			assert.Equal(t, tc.wantType, last.Type, "title=%q", tc.title)
		})
	}
}

func TestParse_SpeakerNotes(t *testing.T) {
	t.Parallel()
	md := `## Demo Slide

Key point here.

Notes: Remember to show the live demo.
`
	p, err := Parse(md)
	require.NoError(t, err)
	require.Len(t, p.Slides, 1)
	assert.Equal(t, "Remember to show the live demo.", p.Slides[0].Notes)
	assert.NotContains(t, p.Slides[0].Content, "Notes:")
}

func TestParse_CodeBlock(t *testing.T) {
	t.Parallel()
	md := "## Code Example\n\n```go\nfmt.Println(\"hello\")\n```\n"
	p, err := Parse(md)
	require.NoError(t, err)
	require.Len(t, p.Slides, 1)
	assert.Contains(t, p.Slides[0].Content, "<code")
	assert.Contains(t, p.Slides[0].Content, "fmt.Println")
}

func TestParse_NoHeadings(t *testing.T) {
	t.Parallel()
	md := "Just a paragraph.\n\nAnother paragraph."
	p, err := Parse(md)
	require.NoError(t, err)
	require.Len(t, p.Slides, 1)
	assert.Equal(t, domain.SlideTypeContent, p.Slides[0].Type)
	assert.Contains(t, p.Slides[0].Content, "Just a paragraph")
}

func TestParse_SectionSlide(t *testing.T) {
	t.Parallel()
	// Single-word heading → section. Multi-word → content.
	md := "# Title\n\nIntro.\n\n## Demo\n\nContent.\n\n## Part Two\n\nMore."
	p, err := Parse(md)
	require.NoError(t, err)
	byTitle := map[string]domain.SlideType{}
	for _, s := range p.Slides {
		byTitle[s.Title] = s.Type
	}
	assert.Equal(t, domain.SlideTypeSection, byTitle["Demo"])
	assert.Equal(t, domain.SlideTypeContent, byTitle["Part Two"])
}

func TestParse_Table(t *testing.T) {
	t.Parallel()
	md := "## Data\n\n| Col A | Col B |\n|-------|-------|\n| 1     | 2     |\n"
	p, err := Parse(md)
	require.NoError(t, err)
	require.Len(t, p.Slides, 1)
	assert.Contains(t, p.Slides[0].Content, "<table")
}

func TestParse_SlideCountMatchesHeadings(t *testing.T) {
	t.Parallel()
	var sb strings.Builder
	for i := 1; i <= 10; i++ {
		sb.WriteString("## Slide ")
		sb.WriteString(strings.Repeat("x", i)) // unique titles
		sb.WriteString("\n\nContent.\n\n")
	}
	p, err := Parse(sb.String())
	require.NoError(t, err)
	assert.Len(t, p.Slides, 10)
}
