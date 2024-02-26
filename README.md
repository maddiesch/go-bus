# go-bus - Generic Event Bus

[![tag](https://img.shields.io/github/v/tag/maddiesch/go-bus.svg)](https://github.com/maddiesch/go-bus/releases)
[![codecov](https://codecov.io/gh/maddiesch/go-bus/graph/badge.svg?token=1PZ250SBC7)](https://codecov.io/gh/maddiesch/go-bus)
[![GoDoc](https://godoc.org/github.com/maddiesch/go-bus?status.svg)](https://pkg.go.dev/github.com/maddiesch/go-bus)
![Build Status](https://github.com/maddiesch/go-bus/actions/workflows/ci.yml/badge.svg)
[![License](https://img.shields.io/github/license/maddiesch/go-bus)](./LICENSE)

A multi-producer multi-consumer generic event bus.

## 🔨 Usage

### Get

`go get github.com/maddiesch/go-bus`

### Import

```go
import (
  "github.com/maddiesch/go-bus"
)
```

### Use

```go
func main() {
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
}

func produceEvents(eventBus *bus.Bus[string]) {
  eventBus.Publish("Hello, World!")
  eventBus.Publish("stop")
}
```
