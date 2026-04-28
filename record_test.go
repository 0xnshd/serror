package serror_test

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/0xnshd/serror"
	"github.com/0xnshd/testingx"
	"github.com/google/go-cmp/cmp/cmpopts"
)

const (
	slogkeyErrorTrait = "trait"
	slogkeyErrorCode  = "code"
	slogkeyErrorTrace = "trace"
	slogkeyErrorCause = "cause"
)

var sortAttrs func(a, b slog.Attr) bool = func(a, b slog.Attr) bool {
	return a.Key > b.Key
}

func Test_ErrorRecord_Error(t *testing.T) {
	tests := []struct {
		name  string
		input serror.ErrorRecord
		want  string
	}{
		{
			name:  "empty ErrorRecord returns empty string",
			input: serror.ErrorRecord{},
			want:  "",
		},
		{
			name:  "ErrorRecord with error returns error message",
			input: serror.ErrorRecord{Err: errors.New("sample error")},
			want:  "sample error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			if testingx.CheckPanic(t, func() { got = tt.input.Error() }, serror.PanicNilError) {
				return
			}

			testingx.Check(t, got, tt.want)
		})
	}
}

func Test_ErrorRecord_Attrs(t *testing.T) {
	tests := []struct {
		name  string
		input serror.ErrorRecord
		want  []slog.Attr
	}{
		{
			name:  "empty ErrorRecord returns empty attrs",
			input: serror.ErrorRecord{},
			want:  []slog.Attr{},
		},
		{
			name: "single trace entry and error",
			input: serror.ErrorRecord{
				Trace: []string{"main"},
				Err:   errors.New("sample error"),
			},
			want: []slog.Attr{
				slog.String(slogkeyErrorTrace, "main"),
				slog.String(slogkeyErrorCause, "sample error"),
			},
		},
		{
			name: "multiple trace entries joined with arrow",
			input: serror.ErrorRecord{
				Trace: []string{"main", "func1"},
				Err:   errors.New("sample error"),
			},
			want: []slog.Attr{
				slog.String(slogkeyErrorTrace, "main -> func1"),
				slog.String(slogkeyErrorCause, "sample error"),
			},
		},
		{
			name: "with trait and no context",
			input: serror.ErrorRecord{
				Trace: []string{"main", "func1"},
				Err:   errors.New("sample error"),
				Trait: serror.ErrorTrait{
					Code:  1,
					Trait: "SampleErrors",
				},
			},
			want: []slog.Attr{
				slog.String(slogkeyErrorTrace, "main -> func1"),
				slog.String(slogkeyErrorCause, "sample error"),
				slog.String(slogkeyErrorTrait, "SampleErrors"),
				slog.Int(slogkeyErrorCode, 1),
			},
		},
		{
			name: "with context and no trait",
			input: serror.ErrorRecord{
				Trace: []string{"main", "func1"},
				Err:   errors.New("sample error"),
				Context: map[string]any{
					"a": 1,
					"b": "2",
				},
			},
			want: []slog.Attr{
				slog.String(slogkeyErrorTrace, "main -> func1"),
				slog.String(slogkeyErrorCause, "sample error"),
				slog.Any("a", 1),
				slog.Any("b", "2"),
			},
		},
		{
			name: "with trait and context",
			input: serror.ErrorRecord{
				Trace: []string{"main", "func1"},
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
			want: []slog.Attr{
				slog.String(slogkeyErrorTrace, "main -> func1"),
				slog.String(slogkeyErrorCause, "sample error"),
				slog.String(slogkeyErrorTrait, "SampleErrors"),
				slog.Int(slogkeyErrorCode, 1),
				slog.Any("a", 1),
				slog.Any("b", "2"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []slog.Attr
			if testingx.CheckPanic(t, func() {
				got = tt.input.Attrs()
			}, serror.PanicNilError) {
				return
			}

			sortOpt := cmpopts.SortSlices(sortAttrs)

			testingx.Check(t, got, tt.want, sortOpt)
		})
	}
}

func Test_ErrorRecord_LogValue(t *testing.T) {
	tests := []struct {
		name  string
		input serror.ErrorRecord
		want  slog.Value
	}{
		{
			name: "returns GroupValue with cause and trace attrs",
			input: serror.ErrorRecord{
				Trace: []string{"main"},
				Err:   errors.New("sample error"),
			},
			want: slog.GroupValue([]slog.Attr{
				slog.String(slogkeyErrorCause, "sample error"),
				slog.String(slogkeyErrorTrace, "main"),
			}...),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.LogValue()
			testingx.Check(t, got.Kind(), slog.KindGroup)

			gotAttrs := got.Group()
			wantAttrs := tt.want.Group()
			sortOpt := cmpopts.SortSlices(sortAttrs)

			testingx.Check(t, gotAttrs, wantAttrs, sortOpt)
		})
	}
}
