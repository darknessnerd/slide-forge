package domain

import "errors"

// SlideType classifies the role of a slide in a presentation.
type SlideType string

const (
	SlideTypeTitle   SlideType = "title"
	SlideTypeContent SlideType = "content"
	SlideTypeCode    SlideType = "code"
	SlideTypeSection SlideType = "section"
	SlideTypeSummary SlideType = "summary"
)

// Theme controls the visual style of the generated presentation.
type Theme string

const (
	ThemeLight     Theme = "light"
	ThemeDark      Theme = "dark"
	ThemeMinimal   Theme = "minimal"
	ThemeCorporate Theme = "corporate"
)

// TransitionStyle controls slide transition animation.
type TransitionStyle string

const (
	TransitionFade  TransitionStyle = "fade"
	TransitionSlide TransitionStyle = "slide"
	TransitionNone  TransitionStyle = "none"
)

// Slide represents a single parsed slide.
type Slide struct {
	Type    SlideType
	Title   string
	Content string // HTML string
	Notes   string // speaker notes, plain text
}

// Presentation is the parsed, structured output ready for rendering.
type Presentation struct {
	Slides []Slide
}

// RenderRequest holds all inputs for one md_to_html_slides invocation.
type RenderRequest struct {
	Markdown                 string
	Theme                    Theme
	TransitionStyle          TransitionStyle
	EnableKeyboardNavigation bool
	IncludeProgressBar       bool
	IncludeSpeakerNotes      bool
}

// ErrEmptyMarkdown is returned when the input markdown is blank.
var ErrEmptyMarkdown = errors.New("markdown input is empty")
