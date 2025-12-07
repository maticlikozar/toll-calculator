package errlog

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"toll/internal/log"
)

func TestError_New(t *testing.T) {
	t.Parallel()

	err := New("message")

	traceErr := &TraceError{}

	if errors.As(err, &traceErr) {
		assert.Equal(
			t,
			"message",
			traceErr.Error(),
		)
	}

	if errors.As(err, &traceErr) {
		assert.Equal(
			t,
			[]string{"internal/errlog/errlog_test.go:17 TestError_New"},
			traceErr.Callers(),
		)
	}
}

func TestError_Errorf(t *testing.T) {
	t.Parallel()

	err := Errorf("%d message with %s", 1, "formatting")

	traceErr := &TraceError{}

	if errors.As(err, &traceErr) {
		assert.Equal(
			t,
			"1 message with formatting",
			traceErr.Error(),
		)
	}

	if errors.As(err, &traceErr) {
		assert.Equal(
			t,
			[]string{"internal/errlog/errlog_test.go:41 TestError_Errorf"},
			traceErr.Callers(),
		)
	}
}

func ExternalFunc() error {
	return New("error from external function")
}

func TestError_StackLog(t *testing.T) {
	t.Parallel()

	err := ExternalFunc()
	fields := StackLog(err)

	assert.Equal(
		t,
		"internal/errlog/errlog_test.go:70 TestError_StackLog",
		fields["log_loc"],
	)

	assert.Equal(
		t,
		[]string{"internal/errlog/errlog_test.go:63 ExternalFunc"},
		fields["err_loc"],
	)
}

func fn0() error {
	err := fn1()

	return fmt.Errorf("fn 0 error: %w", err)
}

func fn1() error {
	err := fn2()

	return Errorf("fn 1 error: %w", err)
}

func fn2() error {
	err := fn3()

	return Errorf("fn 2 error: %w", err)
}

func fn3() error {
	return errors.New("untraced error")
}

func TestError_StackLog_Unwrap(t *testing.T) {
	t.Parallel()

	err := fn0()
	fields := StackLog(err)

	assert.Equal(
		t,
		"internal/errlog/errlog_test.go:111 TestError_StackLog_Unwrap",
		fields["log_loc"],
	)

	assert.Equal(
		t,
		[]string{
			"internal/errlog/errlog_test.go:94 fn1",
			"internal/errlog/errlog_test.go:100 fn2",
		},
		fields["err_loc"],
	)

	assert.Equal(
		t,
		"fn 0 error: fn 1 error: fn 2 error: untraced error",
		err.Error(),
	)
}

func TestError_StackLog_OfNil(t *testing.T) {
	t.Parallel()

	fields := StackLog(nil)

	assert.Equal(
		t,
		"internal/errlog/errlog_test.go:138 TestError_StackLog_OfNil",
		fields["log_loc"],
	)

	assert.Nil(
		t,
		fields["err_loc"],
	)
}

func TestError_StackLog_Unwrap_withCustomFields(t *testing.T) {
	t.Parallel()

	err := fn1()

	customFields := log.Fields{
		"request_id": "request-id",
	}

	fields := StackLog(err, customFields)

	assert.Equal(
		t,
		"internal/errlog/errlog_test.go:161 TestError_StackLog_Unwrap_withCustomFields",
		fields["log_loc"],
	)

	assert.Equal(
		t,
		[]string{
			"internal/errlog/errlog_test.go:94 fn1",
			"internal/errlog/errlog_test.go:100 fn2",
		},
		fields["err_loc"],
	)

	assert.Equal(
		t,
		"request-id",
		fields["request_id"],
	)

	assert.Equal(
		t,
		"fn 1 error: fn 2 error: untraced error",
		err.Error(),
	)
}

var ErrComparison = errors.New("catch error with comparison")

func predefinedError() error {
	return Error(ErrComparison)
}

func TestError_Comparison(t *testing.T) {
	t.Parallel()

	err := predefinedError()
	assert.True(
		t,
		Is(err, ErrComparison),
	)
}

func TestError_ErrorNil(t *testing.T) {
	t.Parallel()

	err := Error(nil)
	assert.NoError(
		t,
		err,
	)
}

func TestError_Merge(t *testing.T) {
	t.Parallel()

	err1 := New("test")
	err2 := New("new error")

	errLevel1 := Merge(err1, err2)
	fields := StackLog(errLevel1)

	err3 := errors.New("test3")
	errLevel2 := Merge(errLevel1, err3)
	fields2 := StackLog(errLevel2)

	assert.Equal(
		t,
		"test3: new error: test",
		errLevel2.Error(),
	)

	assert.Equal(
		t,
		[]string{
			"internal/errlog/errlog_test.go:221 TestError_Merge",
			"internal/errlog/errlog_test.go:223 TestError_Merge",
			"internal/errlog/errlog_test.go:220 TestError_Merge",
		},
		fields["err_loc"],
	)

	assert.Equal(
		t,
		"internal/errlog/errlog_test.go:224 TestError_Merge",
		fields["log_loc"],
	)

	assert.Equal(
		t,
		[]string{
			"internal/errlog/errlog_test.go:227 TestError_Merge",
			"internal/errlog/errlog_test.go:221 TestError_Merge",
			"internal/errlog/errlog_test.go:223 TestError_Merge",
			"internal/errlog/errlog_test.go:220 TestError_Merge",
		},
		fields2["err_loc"],
	)
}

type ErrTooManyRequests struct {
	Err        error
	RetryAfter time.Duration
}

func (e *ErrTooManyRequests) Error() string {
	return fmt.Sprintf("%s (retry after %v)", e.Err.Error(), e.RetryAfter)
}

func (e *ErrTooManyRequests) Unwrap() error {
	return e.Err
}

func TestError_Test(t *testing.T) {
	t.Parallel()

	var tooManyErr *ErrTooManyRequests

	err := Error(&ErrTooManyRequests{
		Err: errors.New("test"),
	})

	if !errors.As(err, &tooManyErr) {
		t.Error("did not match ErrTooManyRequests")
	}
}

func TestError_MergedCompare(t *testing.T) {
	t.Parallel()

	err10 := errors.New("testing10")
	err11 := errors.New("testing11")

	mergedErr10 := Merge(err11, err10)

	assert.True(t, Is(mergedErr10, err10))
	assert.Equal(t, "testing10: testing11", mergedErr10.Error())

	err20 := Errorf("testing20")
	err21 := errors.New("testing21")

	mergedErr20 := Merge(err21, err20)
	mergedErr21 := Merge(err20, err21)

	assert.True(t, Is(mergedErr20, err20))
	assert.True(t, Is(mergedErr21, err21))
	assert.Equal(t, "testing20: testing21", mergedErr20.Error())
	assert.Equal(t, "testing21: testing20", mergedErr21.Error())

	err30 := errors.New("testing30")
	err31 := Errorf("testing31")

	mergedErr30 := Merge(err31, err30)
	mergedErr31 := Merge(err30, err31)

	assert.True(t, Is(mergedErr30, err30))
	assert.True(t, Is(mergedErr31, err31))
	assert.Equal(t, "testing30: testing31", mergedErr30.Error())
	assert.Equal(t, "testing31: testing30", mergedErr31.Error())

	err40 := Errorf("testing40")
	err41 := Errorf("testing41")

	mergedErr40 := Merge(err41, err40)
	mergedErr41 := Merge(err40, err41)

	assert.True(t, Is(mergedErr40, err40))
	assert.True(t, Is(mergedErr41, err41))
	assert.Equal(t, "testing40: testing41", mergedErr40.Error())
	assert.Equal(t, "testing41: testing40", mergedErr41.Error())
}
