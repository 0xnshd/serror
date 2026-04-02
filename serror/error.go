// Package serror
package serror

import (
	"log/slog"
	"maps"
)

func New(err error, trait ErrorTrait, ctx map[string]any) ErrorRecord {
	if err == nil {
		panic(PanicNilError)
	}

	return ErrorRecord{
		Trace:   getTrace(3),
		Trait:   trait,
		Err:     err,
		Context: ctx,
	}
}

func Wrap(ctx map[string]any, errRecord *ErrorRecord) {
	if errRecord.Context == nil {
		errRecord.Context = map[string]any{}
	}
	maps.Copy(errRecord.Context, ctx)
}

func E(err *ErrorRecord) slog.Attr {
	if err == nil {
		return slog.Attr{}
	}

	return slog.Any(slogkeyError, err)
}
