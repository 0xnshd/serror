package serror

import (
	"log/slog"
	"strings"
)

const (
	slogkeyError      = "error"
	slogkeyErrorTrait = "trait"
	slogkeyErrorCode  = "code"
	slogkeyErrorTrace = "trace"
	slogkeyErrorCause = "cause"
)

type ErrorRecord struct {
	Trait   ErrorTrait
	Trace   []string
	Err     error
	Context map[string]any
}

func (e *ErrorRecord) Error() string {
	if e.Err == nil {
		panic(PanicNilError)
	}
	return e.Err.Error()
}

func (e *ErrorRecord) Attrs() []slog.Attr {
	attrs := []slog.Attr{}

	if e.Err == nil {
		panic(PanicNilError)
	}

	if e.Trait.Trait != "" {
		attrs = append(attrs, slog.String(slogkeyErrorTrait, e.Trait.Trait))
		attrs = append(attrs, slog.Int(slogkeyErrorCode, e.Trait.Code))
	}

	attrs = append(attrs, slog.String(slogkeyErrorTrace, strings.Join(e.Trace, " -> ")))
	attrs = append(attrs, slog.String(slogkeyErrorCause, e.Err.Error()))

	for k, v := range e.Context {
		attrs = append(attrs, slog.Any(k, v))
	}

	return attrs
}

func (e *ErrorRecord) LogValue() slog.Value {
	return slog.GroupValue(e.Attrs()...)
}
