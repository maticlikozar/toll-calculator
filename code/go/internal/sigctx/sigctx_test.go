package sigctx

import (
	"syscall"
	"testing"
	"time"
)

func TestSignal_(t *testing.T) {
	t.Parallel()

	deadline := New()

	go func() {
		time.Sleep(10 * time.Millisecond)

		err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		if err != nil {
			panic(err)
		}
	}()

	// Blocks until go-routines sends interrupt signal.
	<-deadline.Done()
}
