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
		want       error
	}{
		{
			inputError: nil,
			inputTrait: ErrorTrait{},
			inputCtx:   map[string]any{},
			want:       &ErrorRecord{},
		},
		{
			inputError: errors.New("sample error"),
			inputTrait: ErrorTrait{},
			inputCtx:   map[string]any{},
			want: &ErrorRecord{
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
			want: &ErrorRecord{
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
			want: &ErrorRecord{
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
			want: &ErrorRecord{
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
		var gotE error
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
			gotE = New(tt.inputError, tt.inputTrait, tt.inputCtx)
		}()

		if _panic {
			continue
		}

		if gotE.Error() != tt.want.Error() {
			t.Errorf(_msg, gotE, tt.want)
			continue
		}

		got := gotE.(*ErrorRecord)
		want := tt.want.(*ErrorRecord)

		want.Trace = append(want.Trace, "testing.tRunner", "runtime.goexit")
		if !reflect.DeepEqual(got.Trace, want.Trace) {
			t.Errorf(_msg, got, want)
			continue
		}

		if !reflect.DeepEqual(got.Trait, want.Trait) {
			t.Errorf(_msg, got, want)
			continue
		}

		if !reflect.DeepEqual(got.Context, want.Context) {
			t.Errorf(_msg, got, want)
			continue
		}
	}
}

func Test_Wrap(t *testing.T) {
	tests := []struct {
		inputCtx   map[string]any
		inputError error
	}{
		{
			inputCtx: map[string]any{
				"c": 3,
				"d": 4,
			},
			inputError: &ErrorRecord{
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
			inputError: &ErrorRecord{
				Context: nil,
			},
		},
		{
			inputCtx: map[string]any{
				"c": 3,
				"d": 4,
			},
			inputError: errors.New("sample error"),
		},
	}

	for _, tt := range tests {
		Wrap(tt.inputCtx, tt.inputError)

		_msg := "Wrap : %v, want %v"

		inputErrorRecord, ok := tt.inputError.(*ErrorRecord)
		if !ok {
			continue
		}

		for k, v := range tt.inputCtx {
			v2, ok := inputErrorRecord.Context[k]

			if !ok {
				t.Errorf(_msg, tt.inputCtx, inputErrorRecord)
			}

			if v2 != v {
				t.Errorf(_msg, tt.inputCtx, inputErrorRecord)
			}
		}
	}
}

func Test_E(t *testing.T) {
	errR := &ErrorRecord{
		Err: errors.New("sample error"),
	}
	err := errors.New("sample error")

	tests := []struct {
		input error
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
		{
			input: err,
			want:  slog.Any(slogkeyError, err),
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

func Test_OfTrait(t *testing.T) {
	tests := []struct {
		input      error
		inputTrait ErrorTrait
		want       bool
	}{
		{
			errors.New("sample error"),
			ErrorTrait{},
			false,
		},
		{
			New(errors.New("sample error"), ErrorTrait{Code: 1, Trait: "Sample"}, map[string]any{}),
			ErrorTrait{Code: 1, Trait: "Sample"},
			true,
		},
		{
			New(errors.New("sample error"), ErrorTrait{Code: 2, Trait: "Sample"}, map[string]any{}),
			ErrorTrait{Code: 1, Trait: "Sample"},
			false,
		},
		{
			New(errors.New("sample error"), ErrorTrait{Code: 1, Trait: "Sample E"}, map[string]any{}),
			ErrorTrait{Code: 1, Trait: "Sample"},
			false,
		},
	}

	for _, tt := range tests {
		got := OfTrait(tt.input, tt.inputTrait)

		if got != tt.want {
			t.Errorf("OfTrait = %v, want %v", got, tt.want)
		}
	}
}
