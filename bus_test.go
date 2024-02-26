package bus_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/maddiesch/go-bus"
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

func TestBusListen(t *testing.T) {
	var (
		l1 atomic.Int32
		l2 atomic.Int32
	)

	b := bus.New[int32]()
	b.SetBufferSize(0)

	cancel1 := b.Listen(func(i int32) {
		l1.Add(i)
	})
	cancel2 := b.Listen(func(i int32) {
		l2.Add(i)
	})

	b.Publish(1)
	cancel1()
	b.Publish(1)
	cancel2()
	b.Publish(1)

	assert.Equal(t, int32(1), l1.Load())
	assert.Equal(t, int32(2), l2.Load())
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
