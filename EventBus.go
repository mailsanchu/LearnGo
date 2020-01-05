package main

import (
	"fmt"
	"github.com/ispringteam/eventbus"
	"sync"
)

const (
	eventId  = "SanchuVarkey"
	min, max = 10, 59
)

type logEvent struct {
	duration string
}

func (e *logEvent) EventID() eventbus.EventID {
	return eventId
}

func main() {
	bus := eventbus.New()
	var wg sync.WaitGroup
	subscribe := bus.Subscribe(eventId, func(e eventbus.Event) {
		se := e.(*logEvent)
		fmt.Printf("%v\n", fmt.Sprintf("%03s", se.duration))
		wg.Done()
	})

	for i := 0; i <= 99; i++ {
		go func(j int) {
			wg.Add(1)
			bus.Publish(&logEvent{
				duration: fmt.Sprint(j),
			})

		}(i)
	}

	wg.Wait()
	bus.Unsubscribe(subscribe)

}
