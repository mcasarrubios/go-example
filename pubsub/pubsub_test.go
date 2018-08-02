package pubsub

import (
	"reflect"
	"testing"
)

type publishTopic struct {
	topic    string
	messages []string
}

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

func TestSubscribe(t *testing.T) {
	publisher := New()
	subscriber := publisher.AsSubscriber()
	expected := []string{"Awesome message"}
	ch := subscriber.Subscribe("Awesome topic")
	publish(&publisher, publishTopic{"Awesome topic", expected})
	checkContents(t, ch, expected)
}

func TestMultipleSubscribers(t *testing.T) {
	publisher := New()
	subscriber := publisher.AsSubscriber()
	expected := []string{"Awesome message"}

	ch1 := subscriber.Subscribe("Awesome topic")
	ch2 := subscriber.Subscribe("Awesome topic")

	publish(&publisher, publishTopic{"Awesome topic", expected})
	checkContents(t, ch1, expected)
	checkContents(t, ch2, expected)
}

func TestMultipleTopics(t *testing.T) {
	publisher := New()
	subscriber := publisher.AsSubscriber()
	expected1 := []string{"Awesome message 1", "Awesome message 2"}
	expected2 := []string{"Awesome message 3", "Awesome message 4"}

	ch1 := subscriber.Subscribe("Awesome topic 1")
	ch2 := subscriber.Subscribe("Awesome topic 2")

	publish(&publisher,
		publishTopic{"Awesome topic 1", expected1},
		publishTopic{"Awesome topic 2", expected2})

	checkContents(t, ch1, expected1)
	checkContents(t, ch2, expected2)
}

func TestMultipleMessagess(t *testing.T) {
	publisher := New()
	subscriber := publisher.AsSubscriber()
	expected := []string{"Awesome message 1", "Awesome message 2", "Awesome message 3"}

	ch := subscriber.Subscribe("Awesome topic")

	publish(&publisher, publishTopic{"Awesome topic", expected})
	checkContents(t, ch, expected)
}

func publish(publisher *Publisher, topics ...publishTopic) {
	for _, item := range topics {
		for _, msg := range item.messages {
			(*publisher).Publish(item.topic, msg)
		}
	}

	(*publisher).Shutdown()
}

func checkContents(t *testing.T, ch <-chan interface{}, vals []string) {
	contents := []string{}

	for v := range ch {
		contents = append(contents, v.(string))
	}

	if !reflect.DeepEqual(contents, vals) {
		t.Fatalf("Invalid channel contents: Expected value %v, Current value %v", vals, contents)
	}
}
