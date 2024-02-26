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

func ExampleBus_Sink() {
	eventBus := bus.New[string]()

	subscription, cancel := eventBus.Sink()
	defer cancel()

	go produceEvents(eventBus)

	for event := range subscription {
		switch event {
		case "stop":
			return
		default:
			fmt.Println(event)
		}
	}

	// Output: Hello, World!
}

func produceEvents(eventBus *bus.Bus[string]) {
	eventBus.Publish("Hello, World!")
	eventBus.Publish("stop")
}
