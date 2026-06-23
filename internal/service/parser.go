package service

import (
	"bytes"
	"strings"

	"github.com/your-org/your-repo/internal/derr"
	"github.com/your-org/your-repo/internal/domain"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

// Parse converts raw Markdown into a Presentation.
// Slide boundaries are determined by ATX headings (# and ##).
// A notes section is everything after a "Notes:" paragraph inside a slide.
func Parse(markdown string) (*domain.Presentation, error) {
	if strings.TrimSpace(markdown) == "" {
		return nil, domain.ErrEmptyMarkdown
	}

	src := []byte(markdown)
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)
	reader := text.NewReader(src)
	doc := md.Parser().Parse(reader)

	var slides []domain.Slide
	var currentTitle string
	var currentType domain.SlideType
	var contentNodes []ast.Node
	var inDoc bool

	flush := func() {
		if !inDoc {
			return
		}
		content, notes := renderNodes(src, contentNodes, md)
		if currentTitle == "" && content == "" {
			return
		}
		slides = append(slides, domain.Slide{
			Type:    currentType,
			Title:   currentTitle,
			Content: content,
			Notes:   notes,
		})
		contentNodes = nil
	}

	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		// Only process top-level children of the document.
		if n.Parent() != doc {
			return ast.WalkContinue, nil
		}

		heading, ok := n.(*ast.Heading)
		if !ok {
			contentNodes = append(contentNodes, n)
			return ast.WalkSkipChildren, nil
		}

		flush()
		inDoc = true
		currentTitle = extractText(n, src)

		if heading.Level == 1 {
			currentType = domain.SlideTypeTitle
		} else if heading.Level == 2 {
			if isSection(currentTitle) {
				currentType = domain.SlideTypeSection
			} else {
				currentType = domain.SlideTypeContent
			}
		} else {
			currentType = domain.SlideTypeContent
		}
		contentNodes = nil
		return ast.WalkSkipChildren, nil
	})
	if err != nil {
		return nil, derr.Internal("parser.Parse", "ast walk", err)
	}
	flush()

	if len(slides) == 0 {
		// No headings: treat entire markdown as a single content slide.
		var buf bytes.Buffer
		if renderErr := md.Renderer().Render(&buf, src, doc); renderErr != nil {
			return nil, derr.Internal("parser.Parse", "render fallback", renderErr)
		}
		slides = append(slides, domain.Slide{
			Type:    domain.SlideTypeContent,
			Content: buf.String(),
		})
	} else {
		markLastAsSummary(slides)
	}

	return &domain.Presentation{Slides: slides}, nil
}

// isSection heuristic: single-word titles are section breaks (e.g. "Demo", "Q&A").
// Multi-word titles are content slides regardless of length.
func isSection(title string) bool {
	return len(strings.Fields(title)) == 1
}

// markLastAsSummary relabels the final slide type when it looks like a summary/conclusion.
func markLastAsSummary(slides []domain.Slide) {
	if len(slides) < 2 {
		return
	}
	last := &slides[len(slides)-1]
	low := strings.ToLower(last.Title)
	for _, kw := range []string{"summary", "conclusion", "recap", "takeaway", "next step", "thank"} {
		if strings.Contains(low, kw) {
			last.Type = domain.SlideTypeSummary
			return
		}
	}
}

// extractText returns the plain-text content of an AST node.
func extractText(n ast.Node, src []byte) string {
	var sb strings.Builder
	_ = ast.Walk(n, func(child ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if t, ok := child.(*ast.Text); ok {
			sb.Write(t.Segment.Value(src))
		}
		return ast.WalkContinue, nil
	})
	return sb.String()
}

// renderNodes renders a list of sibling AST nodes into HTML, extracting speaker notes.
// A paragraph whose first word is "Notes:" is treated as speaker notes.
func renderNodes(src []byte, nodes []ast.Node, md goldmark.Markdown) (content, notes string) {
	md2 := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)

	var htmlBuf, notesBuf bytes.Buffer
	for _, node := range nodes {
		// Check if this paragraph is a notes paragraph.
		if para, ok := node.(*ast.Paragraph); ok {
			raw := extractText(para, src)
			if strings.HasPrefix(strings.TrimSpace(raw), "Notes:") {
				notesBuf.WriteString(strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(raw), "Notes:")))
				notesBuf.WriteString(" ")
				continue
			}
		}

		// Wrap node in a minimal document for rendering.
		tmpDoc := ast.NewDocument()
		tmpDoc.AppendChild(tmpDoc, cloneNode(node))
		var buf bytes.Buffer
		_ = md2.Renderer().Render(&buf, src, tmpDoc)
		htmlBuf.Write(buf.Bytes())
	}

	_ = md
	return strings.TrimSpace(htmlBuf.String()), strings.TrimSpace(notesBuf.String())
}

// cloneNode shallow-clones an AST node so it can be re-parented without mutation.
func cloneNode(n ast.Node) ast.Node {
	return n
}
