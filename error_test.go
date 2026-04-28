package serror_test

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/0xnshd/serror"
	"github.com/0xnshd/testingx"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_New(t *testing.T) {
	tests := []struct {
		name       string
		inputError error
		inputTrait serror.ErrorTrait
		inputCtx   map[string]any
		want       error
	}{
		{
			name:       "nil error panics with PanicNilError",
			inputError: nil,
			inputTrait: serror.ErrorTrait{},
			inputCtx:   map[string]any{},
			want:       &serror.ErrorRecord{},
		},
		{
			name:       "wraps error with empty trait and no context",
			inputError: errors.New("sample error"),
			inputTrait: serror.ErrorTrait{},
			inputCtx:   map[string]any{},
			want: &serror.ErrorRecord{
				Trace:   []string{"serror_test.Test_New.func1.1", "testingx.CheckPanic", "serror_test.Test_New.func1", "testing.tRunner", "runtime.goexit"},
				Err:     errors.New("sample error"),
				Trait:   serror.ErrorTrait{},
				Context: map[string]any{},
			},
		},
		{
			name:       "wraps error with trait and no context",
			inputError: errors.New("sample error"),
			inputTrait: serror.ErrorTrait{
				Code:  1,
				Trait: "SampleErrors",
			},
			inputCtx: map[string]any{},
			want: &serror.ErrorRecord{
				Trace: []string{"serror_test.Test_New.func1.1", "testingx.CheckPanic", "serror_test.Test_New.func1", "testing.tRunner", "runtime.goexit"},
				Err:   errors.New("sample error"),
				Trait: serror.ErrorTrait{
					Code:  1,
					Trait: "SampleErrors",
				},
				Context: map[string]any{},
			},
		},
		{
			name:       "wraps error with no trait and context",
			inputError: errors.New("sample error"),
			inputTrait: serror.ErrorTrait{},
			inputCtx: map[string]any{
				"a": 1,
				"b": "2",
			},
			want: &serror.ErrorRecord{
				Trace: []string{"serror_test.Test_New.func1.1", "testingx.CheckPanic", "serror_test.Test_New.func1", "testing.tRunner", "runtime.goexit"},
				Err:   errors.New("sample error"),
				Trait: serror.ErrorTrait{},
				Context: map[string]any{
					"a": 1,
					"b": "2",
				},
			},
		},
		{
			name:       "wraps error with trait and context",
			inputError: errors.New("sample error"),
			inputTrait: serror.ErrorTrait{
				Code:  1,
				Trait: "SampleErrors",
			},
			inputCtx: map[string]any{
				"a": 1,
				"b": "2",
			},
			want: &serror.ErrorRecord{
				Trace: []string{"serror_test.Test_New.func1.1", "testingx.CheckPanic", "serror_test.Test_New.func1", "testing.tRunner", "runtime.goexit"},
				Err:   errors.New("sample error"),
				Trait: serror.ErrorTrait{
					Code:  1,
					Trait: "SampleErrors",
				},
				Context: map[string]any{
					"a": 1,
					"b": "2",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotE error

			if testingx.CheckPanic(t, func() {
				gotE = serror.New(tt.inputError, tt.inputTrait, tt.inputCtx)
			}, serror.PanicNilError) {
				return
			}

			testingx.Check(t, gotE, tt.want,
				cmp.Comparer(func(x, y error) bool {
					return x.Error() == y.Error()
				}))
		})
	}
}

func Test_Wrap(t *testing.T) {
	tests := []struct {
		name       string
		inputCtx   map[string]any
		inputError error
	}{
		{
			name: "merges into existing context",
			inputCtx: map[string]any{
				"c": 3,
				"d": 4,
			},
			inputError: &serror.ErrorRecord{
				Context: map[string]any{
					"a": 1,
					"b": 2,
				},
			},
		},
		{
			name: "initializes nil context",
			inputCtx: map[string]any{
				"c": 3,
				"d": 4,
			},
			inputError: &serror.ErrorRecord{
				Context: nil,
			},
		},
		{
			name: "ignores non-ErrorRecord",
			inputCtx: map[string]any{
				"c": 3,
				"d": 4,
			},
			inputError: errors.New("sample error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serror.Wrap(tt.inputCtx, tt.inputError)

			containsMap := cmpopts.IgnoreMapEntries(func(k string, v any) bool {
				_, exists := tt.inputCtx[k]
				return !exists
			})

			inputRecord, ok := tt.inputError.(*serror.ErrorRecord)
			if !ok {
				return
			}

			testingx.Check(t, inputRecord.Context, tt.inputCtx, containsMap)
		})
	}
}

func Test_OfTrait(t *testing.T) {
	tests := []struct {
		name       string
		input      error
		inputTrait serror.ErrorTrait
		want       bool
	}{
		{
			name:       "non-ErrorRecord returns false",
			input:      errors.New("sample error"),
			inputTrait: serror.ErrorTrait{},
			want:       false,
		},
		{
			name:       "exact trait match returns true",
			input:      serror.New(errors.New("sample error"), serror.ErrorTrait{Code: 1, Trait: "Sample"}, map[string]any{}),
			inputTrait: serror.ErrorTrait{Code: 1, Trait: "Sample"},
			want:       true,
		},
		{
			name:       "mismatched code returns false",
			input:      serror.New(errors.New("sample error"), serror.ErrorTrait{Code: 2, Trait: "Sample"}, map[string]any{}),
			inputTrait: serror.ErrorTrait{Code: 1, Trait: "Sample"},
			want:       false,
		},
		{
			name:       "mismatched trait string returns false",
			input:      serror.New(errors.New("sample error"), serror.ErrorTrait{Code: 1, Trait: "Sample E"}, map[string]any{}),
			inputTrait: serror.ErrorTrait{Code: 1, Trait: "Sample"},
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := serror.OfTrait(tt.input, tt.inputTrait)

			if got != tt.want {
				t.Errorf("OfTrait = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_E(t *testing.T) {
	slogkeyError := "error"

	errR := &serror.ErrorRecord{
		Err: errors.New("sample error"),
	}
	err := errors.New("sample error")

	tests := []struct {
		name  string
		input error
		want  slog.Attr
	}{
		{
			name:  "nil error returns empty Attr",
			input: nil,
			want:  slog.Attr{},
		},
		{
			name:  "ErrorRecord returns slog.Attr with error key",
			input: errR,
			want:  slog.Any(slogkeyError, errR),
		},
		{
			name:  "plain error returns slog.Attr with error key",
			input: err,
			want:  slog.Any(slogkeyError, err),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := serror.E(tt.input)

			testingx.Check(t, got, tt.want)
		})
	}
}
