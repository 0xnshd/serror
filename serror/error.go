// Package serror
package serror

import (
	"log/slog"
	"maps"
)

func New(err error, trait ErrorTrait, ctx map[string]any) error {
	if err == nil {
		panic(PanicNilError)
	}

	return &ErrorRecord{
		Trace:   getTrace(3),
		Trait:   trait,
		Err:     err,
		Context: ctx,
	}
}

func Wrap(ctx map[string]any, err error) {
	errRecord, ok := err.(*ErrorRecord)
	if !ok {
		return
	}

	if errRecord.Context == nil {
		errRecord.Context = map[string]any{}
	}
	maps.Copy(errRecord.Context, ctx)
}

func E(err error) slog.Attr {
	if err == nil {
		return slog.Attr{}
	}

	e, ok := err.(*ErrorRecord)
	if !ok {
		return slog.Any(slogkeyError, err)
	}

	return slog.Any(slogkeyError, e)
}

func OfTrait(err error, trait ErrorTrait) bool {
	e, ok := err.(*ErrorRecord)

	if !ok {
		return false
	}

	return (e.Trait.Code == trait.Code) && (e.Trait.Trait == trait.Trait)
}
