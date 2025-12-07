package sigctx

import (
	"os"
	"os/signal"
	"sync"
)

var (
	once sync.Once

	d Ctx
	c chan os.Signal
)

func New() Ctx {
	once.Do(func() {
		c = make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		dc := make(chan struct{})
		d = dc

		go func() {
			select {
			case <-c:
				close(dc)
			case <-d.Done():
			}
		}()
	})

	return d
}
