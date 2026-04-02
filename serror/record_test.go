package serror

import (
	"errors"
	"log/slog"
	"sort"
	"testing"
)

func sortAttrs(attrs []slog.Attr) {
	sort.Slice(attrs, func(i, j int) bool {
		return attrs[i].Key < attrs[j].Key
	})
}

func Test_ErrorRecord_Error(t *testing.T) {
	tests := []struct {
		input ErrorRecord
		want  string
	}{
		{ErrorRecord{}, ""},
		{ErrorRecord{Err: errors.New("sample error")}, "sample error"},
	}

	for _, tt := range tests {
		var got string
		_panic := false

		func() {
			defer func() {
				if r := recover(); r != nil {
					if r != PanicNilError {
						panic(r)
					}
					_panic = true
				}
			}()
			got = tt.input.Error()
		}()

		if _panic {
			continue
		}

		if got != tt.want {
			t.Errorf("ErrorRecord.Error() = %s, want = %s", got, tt.want)
		}
	}
}

func Test_ErrorRecord_Attrs(t *testing.T) {
	tests := []struct {
		input ErrorRecord
		want  []slog.Attr
	}{
		{
			ErrorRecord{},
			[]slog.Attr{},
		},
		{
			ErrorRecord{
				Trace: []string{"main"},
				Err:   errors.New("sample error"),
			},
			[]slog.Attr{
				slog.String(slogkeyErrorTrace, "main"),
				slog.String(slogkeyErrorCause, "sample error"),
			},
		},
		{
			ErrorRecord{
				Trace: []string{"main", "func1"},
				Err:   errors.New("sample error"),
			},
			[]slog.Attr{
				slog.String(slogkeyErrorTrace, "main -> func1"),
				slog.String(slogkeyErrorCause, "sample error"),
			},
		},
		{
			ErrorRecord{
				Trace: []string{"main", "func1"},
				Err:   errors.New("sample error"),
				Trait: ErrorTrait{
					Code:  1,
					Trait: "SampleErrors",
				},
			},
			[]slog.Attr{
				slog.String(slogkeyErrorTrace, "main -> func1"),
				slog.String(slogkeyErrorCause, "sample error"),
				slog.String(slogkeyErrorTrait, "SampleErrors"),
				slog.Int(slogkeyErrorCode, 1),
			},
		},
		{
			ErrorRecord{
				Trace: []string{"main", "func1"},
				Err:   errors.New("sample error"),
				Context: map[string]any{
					"a": 1,
					"b": "2",
				},
			},
			[]slog.Attr{
				slog.String(slogkeyErrorTrace, "main -> func1"),
				slog.String(slogkeyErrorCause, "sample error"),
				slog.Any("a", 1),
				slog.Any("b", "2"),
			},
		},
		{
			ErrorRecord{
				Trace: []string{"main", "func1"},
				Err:   errors.New("sample error"),
				Trait: ErrorTrait{
					Code:  1,
					Trait: "SampleErrors",
				},
				Context: map[string]any{
					"a": 1,
					"b": "2",
				},
			},
			[]slog.Attr{
				slog.String(slogkeyErrorTrace, "main -> func1"),
				slog.String(slogkeyErrorCause, "sample error"),
				slog.String(slogkeyErrorTrait, "SampleErrors"),
				slog.Int(slogkeyErrorCode, 1),
				slog.Any("a", 1),
				slog.Any("b", "2"),
			},
		},
	}

	_msg := "ErrorRecord.Attrs() = %v, want = %v"

	for _, tt := range tests {
		var got []slog.Attr
		_panic := false

		func() {
			defer func() {
				if r := recover(); r != nil {
					if r != PanicNilError {
						panic(r)
					}
					_panic = true
				}
			}()

			got = tt.input.Attrs()
		}()

		if _panic {
			continue
		}

		if len(got) != len(tt.want) {
			t.Errorf(_msg, got, tt.want)
			continue
		}

		sortAttrs(got)
		sortAttrs(tt.want)

		for i := range len(got) {
			if !got[i].Equal(tt.want[i]) {
				t.Errorf(_msg, got, tt.want)
				break
			}
		}
	}
}

func Test_ErrorRecord_LogValue(t *testing.T) {
	tests := []struct {
		input ErrorRecord
		want  slog.Value
	}{
		{
			input: ErrorRecord{
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
		got := tt.input.LogValue()

		_msg := "ErrorRecord.LogValue() = %v, want %v"

		if got.Kind() != tt.want.Kind() {
			t.Errorf(_msg, got, tt.want)
		}

		gotAttrs := got.Group()
		wantAttrs := tt.want.Group()

		sortAttrs(gotAttrs)
		sortAttrs(wantAttrs)

		for i := range len(gotAttrs) {
			if !gotAttrs[i].Equal(wantAttrs[i]) {
				t.Errorf(_msg, got, tt.want)
			}
		}
	}
}
