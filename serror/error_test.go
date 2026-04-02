package serror

import (
	"errors"
	"log/slog"
	"reflect"
	"testing"
)

func Test_New(t *testing.T) {
	tests := []struct {
		inputError error
		inputTrait ErrorTrait
		inputCtx   map[string]any
		want       ErrorRecord
	}{
		{
			inputError: nil,
			inputTrait: ErrorTrait{},
			inputCtx:   map[string]any{},
			want:       ErrorRecord{},
		},
		{
			inputError: errors.New("sample error"),
			inputTrait: ErrorTrait{},
			inputCtx:   map[string]any{},
			want: ErrorRecord{
				Trace:   []string{"serror.Test_New.func1", "serror.Test_New"},
				Err:     errors.New("sample error"),
				Trait:   ErrorTrait{},
				Context: map[string]any{},
			},
		},
		{
			inputError: errors.New("sample error"),
			inputTrait: ErrorTrait{
				Code:  1,
				Trait: "SampleErrors",
			},
			inputCtx: map[string]any{},
			want: ErrorRecord{
				Trace: []string{"serror.Test_New.func1", "serror.Test_New"},
				Err:   errors.New("sample error"),
				Trait: ErrorTrait{
					Code:  1,
					Trait: "SampleErrors",
				},
				Context: map[string]any{},
			},
		},
		{
			inputError: errors.New("sample error"),
			inputTrait: ErrorTrait{},
			inputCtx: map[string]any{
				"a": 1,
				"b": "2",
			},
			want: ErrorRecord{
				Trace: []string{"serror.Test_New.func1", "serror.Test_New"},
				Err:   errors.New("sample error"),
				Trait: ErrorTrait{},
				Context: map[string]any{
					"a": 1,
					"b": "2",
				},
			},
		},
		{
			inputError: errors.New("sample error"),
			inputTrait: ErrorTrait{
				Code:  1,
				Trait: "SampleErrors",
			},
			inputCtx: map[string]any{
				"a": 1,
				"b": "2",
			},
			want: ErrorRecord{
				Trace: []string{"serror.Test_New.func1", "serror.Test_New"},
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
		},
	}

	_msg := "New = %v, want %v"

	for _, tt := range tests {
		var got ErrorRecord
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
			got = New(tt.inputError, tt.inputTrait, tt.inputCtx)
		}()

		if _panic {
			continue
		}

		if got.Error() != tt.want.Error() {
			t.Errorf(_msg, got, tt.want)
			continue
		}

		tt.want.Trace = append(tt.want.Trace, "testing.tRunner", "runtime.goexit")
		if !reflect.DeepEqual(got.Trace, tt.want.Trace) {
			t.Errorf(_msg, got, tt.want)
			continue
		}

		if !reflect.DeepEqual(got.Trait, tt.want.Trait) {
			t.Errorf(_msg, got, tt.want)
			continue
		}

		if !reflect.DeepEqual(got.Context, tt.want.Context) {
			t.Errorf(_msg, got, tt.want)
			continue
		}
	}
}

func Test_Wrap(t *testing.T) {
	tests := []struct {
		inputCtx         map[string]any
		inputErrorRecord ErrorRecord
	}{
		{
			inputCtx: map[string]any{
				"c": 3,
				"d": 4,
			},
			inputErrorRecord: ErrorRecord{
				Context: map[string]any{
					"a": 1,
					"b": 2,
				},
			},
		},
		{
			inputCtx: map[string]any{
				"c": 3,
				"d": 4,
			},
			inputErrorRecord: ErrorRecord{
				Context: nil,
			},
		},
	}

	for _, tt := range tests {
		Wrap(tt.inputCtx, &tt.inputErrorRecord)

		_msg := "Wrap : %v, want %v"
		for k, v := range tt.inputCtx {
			v2, ok := tt.inputErrorRecord.Context[k]

			if !ok {
				t.Errorf(_msg, tt.inputCtx, tt.inputErrorRecord)
			}

			if v2 != v {
				t.Errorf(_msg, tt.inputCtx, tt.inputErrorRecord)
			}
		}
	}
}

func Test_E(t *testing.T) {
	errR := &ErrorRecord{
		Err: errors.New("sample error"),
	}
	tests := []struct {
		input *ErrorRecord
		want  slog.Attr
	}{
		{
			input: nil,
			want:  slog.Attr{},
		},
		{
			input: errR,
			want:  slog.Any(slogkeyError, errR),
		},
	}

	for _, tt := range tests {
		got := E(tt.input)

		_msg := "E = %v, want %v"

		if got.Key != tt.want.Key {
			t.Errorf(_msg, got, tt.want)
		}

		if got.Value.Any() != tt.want.Value.Any() {
			t.Errorf(_msg, got, tt.want)
		}
	}
}
