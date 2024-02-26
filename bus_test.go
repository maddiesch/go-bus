package bus_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/maddiesch/go-bus"
	"github.com/stretchr/testify/assert"
)

func TestBus(t *testing.T) {
	b := bus.New[int]()
	b.SetBufferSize(1)

	listener, cancel := b.Sink()
	defer cancel()

	go b.Publish(100)

	e := receiveOnce(t, listener)

	assert.Equal(t, 100, e)
}

func TestBusListen(t *testing.T) {
	b := bus.New[int32]()
	b.SetBufferSize(5)

	var (
		l1 atomic.Int32
		l2 atomic.Int32

		wg1 sync.WaitGroup
		wg2 sync.WaitGroup
	)

	wg1.Add(1)
	wg2.Add(2)

	cancel1 := bus.Listen(b, func(v int32) {
		l1.Add(v)
		wg1.Done()
	})
	cancel2 := bus.Listen(b, func(v int32) {
		l2.Add(v)
		wg2.Done()
	})

	b.Publish(1)
	cancel1()
	b.Publish(2)
	cancel2()
	b.Publish(3)

	wg1.Wait()
	wg2.Wait()

	assert.Equal(t, int32(1), l1.Load())
	assert.Equal(t, int32(3), l2.Load())
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
