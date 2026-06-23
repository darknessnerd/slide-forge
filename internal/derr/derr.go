package derr

// Kind classifies the category of a DError.
type Kind string

const (
	KindValidation Kind = "validation"
	KindInternal   Kind = "internal"
	KindNotFound   Kind = "not_found"
)

// DError is a structured domain error with a kind, operation label, and human-readable message.
type DError struct {
	Kind Kind
	Op   string
	Msg  string
	Err  error
}

func (e *DError) Error() string {
	if e.Err != nil {
		return e.Op + ": " + e.Msg + ": " + e.Err.Error()
	}
	return e.Op + ": " + e.Msg
}

func (e *DError) Unwrap() error { return e.Err }

// Validation returns a KindValidation DError for invalid input.
func Validation(op, msg string) *DError {
	return &DError{Kind: KindValidation, Op: op, Msg: msg}
}

// Internal returns a KindInternal DError wrapping an unexpected error.
func Internal(op, msg string, cause error) *DError {
	return &DError{Kind: KindInternal, Op: op, Msg: msg, Err: cause}
}

// NotFound returns a KindNotFound DError.
func NotFound(op, msg string) *DError {
	return &DError{Kind: KindNotFound, Op: op, Msg: msg}
}
