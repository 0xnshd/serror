package serror

import "testing"

func Test_getName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"github.com/user/repo/pkg.Function", "pkg.Function"},
		{"github.com/user/repo/pkg.(*Service).Do", "pkg.(*Service).Do"},
		{"github.com/user/repo/pkg.Service.Do", "pkg.Service.Do"},
		{"main.main", "main.main"},
		{"net/http.(*Server).Serve", "http.(*Server).Serve"},
		{"runtime.goexit", "runtime.goexit"},
		{"github.com/user/repo/pkg.Function.func1", "pkg.Function.func1"},
		{"github.com/user/repo/db.(*Service).Do.func2", "db.(*Service).Do.func2"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := getName(tt.input)

			if got != tt.want {
				t.Errorf("getName(%s) = %s, want %s", tt.input, got, tt.want)
			}
		})
	}
}

func Test_getTrace(t *testing.T) {
	type traceFunc func() []string
	tests := []struct {
		f    traceFunc
		want []string
	}{
		{
			func() []string {
				return getTrace(2)
			},
			[]string{"serror.Test_getTrace.func1", "serror.Test_getTrace"},
		},
		{
			func() []string {
				return func() []string {
					return getTrace(2)
				}()
			},
			[]string{"serror.Test_getTrace.func2.1", "serror.Test_getTrace.func2", "serror.Test_getTrace"},
		},
		{
			func() []string {
				return func() []string {
					return func() []string {
						return getTrace(2)
					}()
				}()
			},
			[]string{"serror.Test_getTrace.func3.Test_getTrace.func3.1.2", "serror.Test_getTrace.func3.1", "serror.Test_getTrace.func3", "serror.Test_getTrace"},
		},
		{
			func() []string {
				return func() []string {
					return func() []string {
						return getTrace(3)
					}()
				}()
			},
			[]string{"serror.Test_getTrace.func4.1", "serror.Test_getTrace.func4", "serror.Test_getTrace"},
		},
	}

	for _, tt := range tests {

		tt.want = append(tt.want, "testing.tRunner", "runtime.goexit")

		got := tt.f()

		if len(got) != len(tt.want) {
			t.Errorf("Trace = %v, want %v", got, tt.want)
			continue
		}

		for i := range len(got) - 1 {
			if got[i] != tt.want[i] {
				t.Errorf("Trace = %v, want %v", got, tt.want)
			}
		}
	}
}
