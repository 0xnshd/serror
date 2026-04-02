// Package serror
package serror

import (
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
