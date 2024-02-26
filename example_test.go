package bus_test

import (
	"fmt"

	"github.com/maddiesch/go-bus"
)

func ExampleBus() {
	b := bus.New[string]()

	listener, cancel := b.Sink()
	defer cancel()

	go b.Publish("Hello, World!")

	message := <-listener

	fmt.Println(message)
	// Output: Hello, World!
}
