// Package bus implements a simple multi-publisher multi-listener message bus.
package bus

import (
	"sync"
)

// Function signature for the cancelation returned with the listener channel, called when you no longer want to receive messages from the subscription
type Canceler func()

// Bus is the data structure that manages publishing events to the subscribed listeners
type Bus[E any] struct {
	listeners map[uint64]chan E
	mu        sync.RWMutex
	lid       uint64
	chanSize  int
}

// Create a new Bus instance
func New[E any]() *Bus[E] {
	return &Bus[E]{
		listeners: make(map[uint64]chan E),
		lid:       1,
		chanSize:  1,
	}
}

// Set the channel buffer size for the bus. This will only affect new subscriptions.
func (b *Bus[E]) SetBufferSize(size int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.chanSize = size
}

// Add a subscription to the bus
func (b *Bus[E]) Sink() (<-chan E, Canceler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	id := b.lid + 1
	b.lid = id

	ch := make(chan E, b.chanSize)

	b.listeners[id] = ch

	return ch, func() {
		b.mu.Lock()
		defer b.mu.Unlock()

		if c, ok := b.listeners[id]; ok {
			delete(b.listeners, id)
			close(c)
		}
	}
}

// Publish a message to the bus
func (b *Bus[E]) Publish(e E) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, c := range b.listeners {
		c <- e
	}
}

// Listen provides a convenient way to listen for events on a background go routine
func (b *Bus[E]) Listen(fn func(E)) Canceler {
	ch, cancel := b.Sink()
	stopChan := make(chan struct{})
	waitChan := make(chan struct{})

	go func() {
		defer func() {
			cancel()
			close(stopChan)
			close(waitChan)
		}()

		for {
			select {
			case e := <-ch:
				fn(e)
			case <-stopChan:
				return
			}
		}
	}()

	return func() {
		stopChan <- struct{}{}
		cancel()
		<-waitChan
	}
}
