package bus_test

import (
	"sync"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/maddiesch/bus"
)

func TestBus(t *testing.T) {
	b := bus.New[int]()

	listener, cancel := b.Sink()
	defer cancel()

	var waiter sync.WaitGroup

	waiter.Add(1)
	go func() {
		defer waiter.Done()

		e := <-listener

		spew.Dump(e)
	}()

	b.Publish(100)

	waiter.Wait()
}
