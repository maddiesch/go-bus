package bus_test

import (
	"testing"
	"time"

	"github.com/maddiesch/bus"
	"github.com/stretchr/testify/assert"
)

func TestBus(t *testing.T) {
	b := bus.New[int]()

	listener, cancel := b.Sink()
	defer cancel()

	go b.Publish(100)

	e := receiveOnce(t, listener)

	assert.Equal(t, 100, e)
}

func receiveOnce[V any](t *testing.T, ch <-chan V) V {
	select {
	case e := <-ch:
		return e
	case <-time.After(time.Millisecond * 100):
		t.Log("failed to receive withing specified timeout")
		t.FailNow()
		panic("should have failed immediately")
	}
}
