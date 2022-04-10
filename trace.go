package whoops

import (
	"fmt"
	"runtime"
	"strings"
)

const (
	errFieldStackTrace Field[[]runtime.Frame] = "stacktrace"
)

type traceError struct {
	error
	stack []any
}

func Trace(err error) error {
	return trace(err, 3)
}

func trace(err error, skip int) error {
	if err == nil {
		return nil
	}
	var frames []runtime.Frame
	for i := skip; i < 1000; i++ {
		frame, ok := getFrame(i)
		if !ok {
			break
		}
		frames = append(frames, frame)
	}

	return Compose(
		err,
		ComposeEnrich(
			errFieldStackTrace.Val(frames),
		),
		ComposeCustomMessage(func(err error) string {
			var sb strings.Builder
			sb.WriteString(err.Error())

			stack := FormatStacktrace(err)
			if len(stack) > 0 {
				sb.WriteString("\nStack trace:")
				for _, s := range stack {
					sb.WriteByte('\n')
					sb.WriteByte('\t')
					sb.WriteString(s)
				}
			}
			return sb.String()
		}),
	)
}

func (s String) Trace() error {
	return trace(s, 3)
}

func (s formattedError) Trace() error {
	return trace(s, 3)
}

func (s Group) Trace() error {
	return trace(s, 3)
}
func (s wrapper) Trace() error {
	return trace(s, 3)
}

func getFrame(calldepth int) (runtime.Frame, bool) {
	pc, file, line, ok := runtime.Caller(calldepth)
	if !ok {
		return runtime.Frame{}, false
	}

	frame := runtime.Frame{
		PC:   pc,
		File: file,
		Line: line,
	}

	funcForPc := runtime.FuncForPC(pc)
	if funcForPc != nil {
		frame.Func = funcForPc
		frame.Function = funcForPc.Name()
		frame.Entry = funcForPc.Entry()
	}

	return frame, true
}

func FormatStacktrace(err error) (res []string) {
	frames, ok := errFieldStackTrace.GetFrom(err)
	if !ok {
		return nil
	}
	for _, frame := range frames {
		res = append(res, fmt.Sprintf("%s:%d\t%s", frame.File, frame.Line, frame.Function))
	}
	return
}
