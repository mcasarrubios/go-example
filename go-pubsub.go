package main

import (
	"fmt"
	"time"

	"github.com/mcasarrubios/go-pubsub/pubsub"
)

const (
	eventName = "Awesome Event"
)

type awesomeEvent struct {
	name   string
	pubsub pubsub.Publisher
}

func (e *awesomeEvent) subscribe() chan interface{} {
	return e.pubsub.AsSubscriber().Subscribe(e.name)
}

func on(ch chan interface{}) {
	for msg := range ch {
		fmt.Println("On Event received", msg)
	}
}

func publish(e *awesomeEvent) {
	for i := 0; i < 5; i++ {
		time.Sleep(time.Millisecond * 500)
		e.pubsub.Publish(e.name, fmt.Sprintf("Message %d", i))
	}
}

func main() {
	myEvent := &awesomeEvent{eventName, pubsub.New()}
	ch := myEvent.subscribe()
	go on(ch)
	publish(myEvent)
}
