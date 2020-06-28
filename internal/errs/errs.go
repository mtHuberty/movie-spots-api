package errs

import (
	"fmt"
	"net/http"
	"runtime"

	log "github.com/sirupsen/logrus"
)

type stackTraceNode struct {
	Line int
	File string
}

type tracedError struct {
	Err        error
	stackTrace []*stackTraceNode
}

func (te tracedError) Error() string {
	return te.Err.Error()
}

func wrap(pkg string, fn string, err error) tracedError {
	switch e := err.(type) {
	case tracedError:
		e.Err = fmt.Errorf("%s.%s: %w", pkg, fn, err)
		return e
	default:
		return tracedError{
			Err:        fmt.Errorf("%s.%s: %w", pkg, fn, err),
			stackTrace: []*stackTraceNode{},
		}
	}
}

func trace(err error) tracedError {
	// 2 causes runtime.Caller to skip this func and it's calling func (for now, the WrapTrace func below)
	if _, fl, ln, ok := runtime.Caller(2); ok {
		switch e := err.(type) {
		case tracedError:
			return tracedError{
				Err: e.Err,
				stackTrace: append(e.stackTrace, &stackTraceNode{
					Line: ln,
					File: fl,
				}),
			}
		default:
			return tracedError{
				Err: err,
				stackTrace: []*stackTraceNode{{
					Line: ln,
					File: fl,
				}},
			}
		}
	}
	return tracedError{
		Err:        err,
		stackTrace: []*stackTraceNode{},
	}
}

// WrapTrace calls wrap() and trace(), adding a node to the callstack and wrapping the err. pkg should be the package name and fn should be the function name.
func WrapTrace(pkg string, fn string, err error) tracedError {
	return wrap(pkg, fn, trace(err))
}

func (te tracedError) GetStackTrace() string {
	if len(te.stackTrace) > 0 {
		var str string
		for _, val := range te.stackTrace {
			str += fmt.Sprintf("%s:%d\n", val.File, val.Line)
		}
		return str
	}
	return "No stacktrace available"
}

func HandleErrorResponse(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case tracedError:
		logEntry := log.WithFields(log.Fields{
			"stackTrace": e.GetStackTrace(),
		})
		logEntry.WithError(e).Error("tracedError")
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	default:
		log.WithError(err).Error("Error")
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
