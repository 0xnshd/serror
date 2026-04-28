package serror

import (
	"runtime"
	"strings"
)

func getName(path string) string {
	// strip path
	parts := strings.Split(path, "/")
	name := parts[len(parts)-1]

	return name
}

func getTrace(skip int) []string {
	pcs := make([]uintptr, 8)
	n := runtime.Callers(skip, pcs)
	frames := runtime.CallersFrames(pcs[:n])

	trace := make([]string, 0, n)

	for {
		frame, more := frames.Next()

		fn := getName(frame.Function)
		trace = append(trace, fn)

		if !more {
			break
		}
	}

	return trace
}
