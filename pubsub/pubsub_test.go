package pubsub

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewReturnsPublisher(t *testing.T) {
	ps := New()

	if _, ok := ps.(Publisher); !ok {
		t.Fatalf("New returns an invalid type")
	}
}

func TestAsSubscriberReturnsSuscriber(t *testing.T) {
	ps := New().AsSubscriber()
	if _, ok := ps.(Subscriber); !ok {
		t.Fatalf("AsSubscriber fails")
	}
}

func TestPublish(t *testing.T) {
	publisher := New()
	subscriber := publisher.AsSubscriber()
	ch := subscriber.Subscribe("Awesome topic")
	publisher.Publish("Awesome topic", "Awesome message")
	fmt.Printf("--->B0 ")
	checkContents(t, ch, []string{"Awesome message"})
}

func checkContents(t *testing.T, ch chan interface{}, vals []string) {
	contents := []string{}
	fmt.Printf("--->B1")
	for v := range ch {
		contents = append(contents, v.(string))
	}

	if !reflect.DeepEqual(contents, vals) {
		t.Fatalf("Invalid channel contents: Expected value %v, Current value %v", vals, contents)
	}
}
