package errlog

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"toll/internal/log"
)

// TraceError is an error that has a message and a caller that created new error.
type TraceError struct {
	merged  *MergedError
	err     error
	callers []string
}

type MergedError struct {
	original error
}

// rightSubstring consumes val from right to left until it encounters
// separator and then returns consumed characters.
// If separator is not encountered the whole unmodified val is returned.
func rightSubstring(val string, separator rune) string {
	end := val

	for i := len(val) - 1; i > 0; i-- {
		if rune(val[i]) == separator {
			end = val[i+1:]

			break
		}
	}

	return end
}

// callerInfo returns formatted string containing caller location.
func callerInfo() string {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return "unable to get function caller from runtime"
	}

	details := runtime.FuncForPC(pc)

	file, line := details.FileLine(pc)
	functionName := details.Name()

	// Remove file path before /go/ to leave only the
	// relevant path from go project root.
	if idx := strings.Index(file, "/go/"); idx != -1 {
		idx = idx + len("/go/")
		file = file[idx:]
	}
	// Remove package name from functionName where format is "package.Function".
	functionName = rightSubstring(functionName, '.')

	return file + ":" + strconv.Itoa(line) + " " + functionName
}

// New returns an error with the supplied message.
// New also records the the caller at the point it was called.
func New(message string) error {
	return &TraceError{
		err:     errors.New(message),
		callers: []string{callerInfo()},
	}
}

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
// Errorf also records the caller at the point it was called.
func Errorf(format string, args ...interface{}) error {
	return &TraceError{
		err:     fmt.Errorf(format, args...),
		callers: []string{callerInfo()},
	}
}

func callersInfo(err error) []string {
	callersInfo := []string{}

	for {
		if err == nil {
			return callersInfo
		}

		if traceErr := traceError(err); traceErr != nil {
			if traceErr.merged == nil {
				callersInfo = append(callersInfo, traceErr.callers...)
			} else {
				return traceErr.Callers()
			}
		}

		err = errors.Unwrap(err)
	}
}

// Merge combines two errors with their respective callers.
func Merge(err error, errNew error) error {
	callers := callersInfo(errNew)
	callers = append(callers, callerInfo())
	callers = append(callers, callersInfo(err)...)

	mergedErrorMsg := strings.Join([]string{errNew.Error(), err.Error()}, ": ")

	return &TraceError{
		err:     errors.New(mergedErrorMsg),
		callers: callers,
		merged: &MergedError{
			original: errNew,
		},
	}
}

// Error wraps existing error into the TraceError.
func Error(err error) error {
	if err == nil {
		return nil
	}

	return &TraceError{
		err:     err,
		callers: []string{callerInfo()},
	}
}

// Is reports whether any error in err's chain matches target.
func Is(err error, target error) bool {
	if traceErr := traceError(err); traceErr != nil {
		if traceErr.merged != nil {
			return errors.Is(traceErr.merged.original, target)
		}
	}

	return errors.Is(err, target)
}

// As finds the first error in err's chain that matches target, and if one is found, sets
// target to that error value and returns true. Otherwise, it returns false.
func (t *TraceError) As(target any) bool {
	return errors.As(t.err, target)
}

// Error returns error message.
func (t *TraceError) Error() string {
	return t.err.Error()
}

// Caller returns the location of the caller that created this error.
func (t *TraceError) Callers() []string {
	return t.callers
}

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
func (t *TraceError) Unwrap() error {
	return t.err
}

// Is compares targer error to the wrapped error in the TraceError.
func (t *TraceError) Is(target error) bool {
	return errors.Is(t.err, target)
}

func traceError(err error) *TraceError {
	// Only try to cast provided error and ignore the wrapped error.
	// Linter warning making sure that wrapped errors are checked needs to be disabled.
	if traceErr, ok := err.(*TraceError); ok { //nolint:errorlint
		return traceErr
	}

	return nil
}

func findTraceError(err error) *TraceError {
	for {
		if err == nil {
			return nil
		}

		if traceErr := traceError(err); traceErr != nil {
			return traceErr
		}

		err = errors.Unwrap(err)
	}
}

// StackLog returns log fields containing the log location and
// where error originated from if error is of type TraceError.
// Custom fields are appended to original error fields when provided.
func StackLog(err error, customFields ...log.Fields) log.Fields {
	fields := log.Fields{}

	for _, cf := range customFields {
		for label, field := range cf {
			fields[label] = field
		}
	}

	fields["log_loc"] = callerInfo()

	traceErr := findTraceError(err)
	if traceErr == nil {
		return fields
	}

	fields["err_loc"] = callersInfo(traceErr)

	return fields
}
